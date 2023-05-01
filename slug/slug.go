package slug

import (
	"strings"
	"unicode"
)

func Make(s string) string {
	r := strings.TrimSpace(s)
	r = strings.ToLower(r)

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
