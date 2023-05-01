package slug

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func Make(s string) string {
	r := strings.TrimSpace(s)
	r = strings.ToLower(r)

	b := make([]byte, len(r))
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, err := t.Transform(b, []byte(s), true)
	if err != nil {
		return ""
	}

	r = string(b)

	for _, c := range r {
		if unicode.IsNumber(c) || unicode.IsLetter(c) {
			continue
		}
		r = strings.Replace(r, string(c), "-", -1)
	}

	for strings.Contains(r, "--") {
		r = strings.Replace(r, "--", "-", -1)
	}

	r = strings.Trim(r, "-")

	return r
}
