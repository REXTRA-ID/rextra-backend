package utils

import (
	"strings"
)

func ToSlug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder

	spaceFound := false
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			if spaceFound {
				b.WriteRune('-')
				spaceFound = false
			}
			b.WriteRune(r)
		case r == ' ':
			spaceFound = true
		}
	}

	return b.String()
}
