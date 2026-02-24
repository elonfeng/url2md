package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elonfeng/url2md/pkg/converter"
	"github.com/elonfeng/url2md/pkg/converter/filetype"
	"github.com/elonfeng/url2md/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var (
		method        string
		retainImages  bool
		retainLinks   bool
		frontmatter   bool
		enableBrowser bool
		timeout       int
		output        string
	)

	root := &cobra.Command{
		Use:   "url2md [url]",
		Short: "Convert web pages to clean Markdown",
		Long:  "url2md fetches a URL and converts it to clean, LLM-friendly Markdown using a three-layer fallback pipeline.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			opts := &converter.Options{
				Method:        method,
				RetainImages:  retainImages,
				RetainLinks:   retainLinks,
				Frontmatter:   frontmatter,
				EnableBrowser: enableBrowser,
				Timeout:       time.Duration(timeout) * time.Second,
				UserAgent:     "url2md/1.0",
				Vision:        visionFromEnv(),
			}

			ctx := context.Background()
			c := converter.New()
			result, err := c.Convert(ctx, url, opts)
			if err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}

			// print header info to stderr
			fmt.Fprintf(os.Stderr, "Title:  %s\n", result.Title)
			fmt.Fprintf(os.Stderr, "Method: %s\n", result.Method)
			fmt.Fprintf(os.Stderr, "Tokens: ~%d\n", result.TokenCount)
			fmt.Fprintf(os.Stderr, "Fetch:  %s\n", result.FetchTime.Round(time.Millisecond))

			if output != "" {
				return os.WriteFile(output, []byte(result.Markdown), 0644)
			}
			fmt.Println(result.Markdown)
			return nil
		},
	}

	root.Flags().StringVarP(&method, "method", "m", "auto", "Conversion method: auto, negotiate, static, browser")
	root.Flags().BoolVar(&retainImages, "images", false, "Retain images in output")
	root.Flags().BoolVar(&retainLinks, "links", true, "Retain links in output")
	root.Flags().BoolVar(&frontmatter, "frontmatter", true, "Prepend YAML frontmatter (title, description, image)")
	root.Flags().BoolVar(&enableBrowser, "browser", false, "Enable headless Chrome fallback")
	root.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds")
	root.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")

	root.AddCommand(serveCmd())
	root.AddCommand(batchCmd())

	return root
}

func serveCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := server.New(port)
			return srv.ListenAndServe()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
	return cmd
}

func batchCmd() *cobra.Command {
	var (
		method        string
		retainImages  bool
		frontmatter   bool
		enableBrowser bool
		timeout       int
	)

	cmd := &cobra.Command{
		Use:   "batch [url1] [url2] ...",
		Short: "Convert multiple URLs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &converter.Options{
				Method:        method,
				RetainImages:  retainImages,
				Frontmatter:   frontmatter,
				RetainLinks:   true,
				EnableBrowser: enableBrowser,
				Timeout:       time.Duration(timeout) * time.Second,
				UserAgent:     "url2md/1.0",
				Vision:        visionFromEnv(),
			}

			ctx := context.Background()
			c := converter.New()

			for _, url := range args {
				if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
					url = "https://" + url
				}

				result, err := c.Convert(ctx, url, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[FAIL] %s: %v\n", url, err)
					continue
				}

				fmt.Fprintf(os.Stderr, "[OK] %s → %d tokens (%s)\n", url, result.TokenCount, result.Method)
				fmt.Println(result.Markdown)
				fmt.Print("\n---\n\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&method, "method", "m", "auto", "Conversion method")
	cmd.Flags().BoolVar(&retainImages, "images", false, "Retain images")
	cmd.Flags().BoolVar(&frontmatter, "frontmatter", true, "Prepend YAML frontmatter")
	cmd.Flags().BoolVar(&enableBrowser, "browser", false, "Enable browser fallback")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds")

	return cmd
}

// visionFromEnv reads Cloudflare Workers AI credentials from environment variables.
func visionFromEnv() *filetype.VisionConfig {
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if accountID == "" || apiToken == "" {
		return nil
	}
	return &filetype.VisionConfig{
		AccountID: accountID,
		APIToken:  apiToken,
	}
}
