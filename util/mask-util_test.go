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
		s = util.MaskValue(s, '*', 12)
		t.Log(s)
	}

}
