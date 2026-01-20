package util_test

import (
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
)

func TestMask(t *testing.T) {

	sarr := []string{
		"ABCDEFGHIJKLMNOP",
		"ABCDEFGHIJKL",
		"ABCDEFGHI",
		"",
	}

	for _, s := range sarr {
		s = util.MaskValue(s, util.MaskOptions{Char: '*', Size: 12})
		t.Log(s)
	}

	for _, s := range sarr {
		s = util.MaskValue(s, util.MaskOptions{Char: '*', Size: 12, Direction: util.MaskAtEnd})
		t.Log(s)
	}

	for _, s := range sarr {
		s = util.MaskValue(s, util.MaskOptions{Char: '*', Size: 4, Direction: util.MaskAtEnd, SizeMode: util.MaskKeepSizeCharsInClear})
		t.Log(s)
	}

	for _, s := range sarr {
		s = util.MaskValue(s, util.MaskOptions{Char: '*', Size: 4, Direction: util.MaskAtBeginning, SizeMode: util.MaskKeepSizeCharsInClear})
		t.Log(s)
	}
}
