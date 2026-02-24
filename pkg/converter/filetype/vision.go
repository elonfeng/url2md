package filetype

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// VisionConfig holds configuration for Cloudflare Workers AI vision.
type VisionConfig struct {
	AccountID string
	APIToken  string
}

// IsConfigured returns true if the vision provider is configured.
func (v *VisionConfig) IsConfigured() bool {
	return v.AccountID != "" && v.APIToken != ""
}

// cfAIRequest is the request body for Cloudflare Workers AI.
type cfAIRequest struct {
	Messages []cfMessage `json:"messages"`
}

type cfMessage struct {
	Role    string    `json:"role"`
	Content []cfBlock `json:"content"`
}

type cfBlock struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *cfImage  `json:"image_url,omitempty"`
}

type cfImage struct {
	URL string `json:"url"`
}

type cfAIResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success bool     `json:"success"`
	Errors  []cfError `json:"errors"`
}

type cfError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type cfAgreeRequest struct {
	Prompt string `json:"prompt"`
}

// DescribeImage sends the image to Cloudflare Workers AI vision and returns a text description.
func DescribeImage(ctx context.Context, cfg *VisionConfig, data []byte, contentType string) (string, error) {
	if !cfg.IsConfigured() {
		return "", fmt.Errorf("vision not configured")
	}

	// Build base64 data URL
	if contentType == "" {
		contentType = "image/png"
	}
	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data))

	reqBody := cfAIRequest{
		Messages: []cfMessage{
			{
				Role: "user",
				Content: []cfBlock{
					{
						Type:     "image_url",
						ImageURL: &cfImage{URL: dataURL},
					},
					{
						Type: "text",
						Text: "Describe this image in detail. Focus on the main content, colors, layout, and any text visible in the image. Output a concise description suitable for a markdown document.",
					},
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/@cf/meta/llama-3.2-11b-vision-instruct", cfg.AccountID)

	desc, err := callVisionAPI(ctx, cfg, apiURL, body)
	if err != nil && isModelAgreementError(err) {
		// Auto-agree to model license and retry
		if agreeErr := acceptModelAgreement(ctx, cfg, apiURL); agreeErr != nil {
			return "", fmt.Errorf("accept model agreement: %w", agreeErr)
		}
		return callVisionAPI(ctx, cfg, apiURL, body)
	}
	return desc, err
}

// isModelAgreementError checks if the error is a model license agreement requirement.
func isModelAgreementError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Model Agreement")
}

// acceptModelAgreement sends "agree" to accept the model's license agreement.
func acceptModelAgreement(ctx context.Context, cfg *VisionConfig, apiURL string) error {
	agreeBody, _ := json.Marshal(cfAgreeRequest{Prompt: "agree"})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(agreeBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

// callVisionAPI sends a request and parses the response.
func callVisionAPI(ctx context.Context, cfg *VisionConfig, apiURL string, body []byte) (string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.APIToken)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("vision API call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vision API HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var aiResp cfAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if !aiResp.Success {
		msg := "unknown error"
		if len(aiResp.Errors) > 0 {
			msg = aiResp.Errors[0].Message
		}
		return "", fmt.Errorf("vision API error: %s", msg)
	}

	return aiResp.Result.Response, nil
}
