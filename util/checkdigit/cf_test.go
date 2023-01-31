package checkdigit_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/checkdigit"
	"testing"
)

func TestCFCheckDigit(t *testing.T) {

	arr := []string{
		"MPRMLS62S21G337",
		"BBBTTT20H12X122",
		"MPRNCR07M53H501",
		"SCHDTL63T49H501",
		"MPR MLS 62S21 G337 ",
	}

	for _, s := range arr {
		t.Log(s, checkdigit.ComputeCFCheckDigit(s))
	}
}
