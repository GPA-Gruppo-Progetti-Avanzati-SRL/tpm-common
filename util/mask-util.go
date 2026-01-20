package util

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type MaskDirection string
type MaskLengthMode string

const (
	MaskAtBeginning MaskDirection = "at-beginning"
	MaskAtEnd       MaskDirection = "at-end"

	MaskLengthModeMask       MaskLengthMode = "mask"
	MaskKeepSizeCharsInClear MaskLengthMode = "keep-size-clear"
)

type MaskOptions struct {
	Size      int
	Direction MaskDirection
	SizeMode  MaskLengthMode
	Char      rune
}

func MaskValue(s string, opts MaskOptions) string {
	const semLogContext = "util::mask-value"

	if opts.Size <= 0 {
		log.Warn().Int("mask-length", opts.Size).Msg(semLogContext + " - invalid length")
		return s
	}

	if opts.Direction == "" {
		opts.Direction = MaskAtBeginning
	}

	if opts.SizeMode == "" {
		opts.SizeMode = MaskLengthModeMask
	}

	if len(s) == 0 || (len(s) <= opts.Size && opts.SizeMode == MaskKeepSizeCharsInClear) {
		return s
	}

	if opts.SizeMode == MaskKeepSizeCharsInClear {
		opts.Size = len(s) - opts.Size
		opts.SizeMode = MaskLengthModeMask
	}

	if opts.Direction == MaskAtEnd {
		return rightMaskValue(s, opts.Char, opts.Size)
	}

	return leftMaskValue(s, opts.Char, opts.Size)
}

func leftMaskValue(s string, maskingChar rune, maskLength int) string {

	if maskLength > len(s) {
		maskLength = len(s)
	}

	ms := strings.Repeat(string(maskingChar), maskLength)
	if len(s) <= maskLength {
		return ms
	}

	s = ms + s[maskLength:]
	return s
}

func rightMaskValue(s string, maskingChar rune, maskLength int) string {
	if maskLength > len(s) {
		maskLength = len(s)
	}

	ms := strings.Repeat(string(maskingChar), maskLength)
	if len(s) <= maskLength {
		return ms
	}

	s = s[:len(s)-maskLength] + ms
	return s
}
