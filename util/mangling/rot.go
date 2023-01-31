package mangling

import "strings"

var alphabetMixed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var alphabetUpper = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func AlphabetRot(s string, onlyUppercase bool) string {
	var sb strings.Builder

	var alphabet string
	var rotationParams struct {
		midPoint   int
		upperPoint int
		offset     int
	}
	if onlyUppercase {
		alphabet = alphabetUpper
		s = strings.ToUpper(s)
		rotationParams.midPoint = 17
		rotationParams.upperPoint = 35
		rotationParams.offset = 18
	} else {
		alphabet = alphabetMixed
		rotationParams.midPoint = 30
		rotationParams.upperPoint = 61
		rotationParams.offset = 31
	}

	for _, c := range s {
		ndx := strings.Index(alphabet, string(c))
		if ndx >= 0 && ndx <= rotationParams.midPoint {
			sb.WriteRune(rune(alphabet[ndx+rotationParams.offset]))
		} else if ndx > rotationParams.midPoint && ndx <= rotationParams.upperPoint {
			sb.WriteRune(rune(alphabet[ndx-rotationParams.offset]))
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

func Rot13(s string) string {

	var sb strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'm') || (c >= 'A' && c <= 'M') {
			sb.WriteRune(c + 13)
		} else if (c > 'm' && c <= 'z') || (c > 'M' && c <= 'Z') {
			sb.WriteRune(c - 13)
			//} else if c >= '0' && c <= '4' {
			//	sb.WriteRune(c + 5)
			//} else if c > '4' && c <= '9' {
			//	sb.WriteRune(c - 5)
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
