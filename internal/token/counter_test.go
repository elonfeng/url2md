package token

import "testing"

func TestEstimate_EnglishOnly(t *testing.T) {
	text := "Hello world this is a test"
	count := Estimate(text)
	// 6 words * 1.3 = 7
	if count < 5 || count > 15 {
		t.Errorf("unexpected token count for English text: %d", count)
	}
}

func TestEstimate_CJK(t *testing.T) {
	text := "这是一个测试"
	count := Estimate(text)
	// 1 word (fields splits on whitespace) * 1.3 + 6 CJK * 1.5 = 1 + 9 = 10
	if count < 5 {
		t.Errorf("expected higher token count for CJK text: %d", count)
	}
}

func TestEstimate_Mixed(t *testing.T) {
	text := "Hello 你好 World 世界"
	count := Estimate(text)
	// 4 words * 1.3 + 4 CJK * 1.5 = 5 + 6 = 11
	if count < 8 {
		t.Errorf("expected higher token count for mixed text: %d", count)
	}
}

func TestEstimate_Empty(t *testing.T) {
	count := Estimate("")
	if count != 0 {
		t.Errorf("expected 0 tokens for empty text, got %d", count)
	}
}

func TestEstimate_Japanese(t *testing.T) {
	text := "テスト テスト"
	count := Estimate(text)
	if count < 5 {
		t.Errorf("expected higher token count for Japanese: %d", count)
	}
}
