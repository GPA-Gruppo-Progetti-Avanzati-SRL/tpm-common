package mangling

import "strings"

func Rot(s string) string {

	var sb strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'm') || (c >= 'A' && c <= 'M') {
			sb.WriteRune(c + 13)
		} else if (c > 'm' && c <= 'z') || (c > 'M' && c <= 'Z') {
			sb.WriteRune(c - 13)
		} else if c >= '0' && c <= '4' {
			sb.WriteRune(c + 5)
		} else if c > '4' && c <= '9' {
			sb.WriteRune(c - 5)
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
