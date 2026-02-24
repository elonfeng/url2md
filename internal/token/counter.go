package token

import (
	"strings"
	"unicode"
)

// Estimate returns an approximate token count for text, accounting for CJK characters.
func Estimate(text string) int {
	words := len(strings.Fields(text))
	cjk := 0
	for _, r := range text {
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hiragana, r) {
			cjk++
		}
	}
	return int(float64(words)*1.3) + int(float64(cjk)*1.5)
}
