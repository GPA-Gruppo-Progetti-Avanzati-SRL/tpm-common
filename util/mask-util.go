package util

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func MaskValue(s string, maskingChar rune, maskLength int) string {
	const semLogContext = "util::mask-value"

	if maskLength <= 0 {
		log.Warn().Int("mask-length", maskLength).Msg(semLogContext + " - invalid length")
		return s
	}

	if len(s) == 0 {
		return s
	}

	ms := strings.Repeat(string(maskingChar), maskLength)

	if len(s) <= maskLength {
		maskLength = len(s)
		return ms
	}

	s = ms + s[maskLength:]
	return s
}
