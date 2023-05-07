package slug

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func Make(s string) string {
	r := strings.ToLower(strings.TrimSpace(s))
	r, _, err := transform.String(normalizer, r)
	if err != nil {
		return ""
	}

	for _, c := range r {
		if unicode.IsNumber(c) || unicode.IsLetter(c) {
			continue
		}
		r = strings.ReplaceAll(r, string(c), "-")
	}

	for strings.Contains(r, "--") {
		r = strings.ReplaceAll(r, "--", "-")
	}

	r = strings.Trim(r, "-")

	return r
}
