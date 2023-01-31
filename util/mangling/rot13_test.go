package mangling_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/mangling"
	"testing"
)

func TestRot13(t *testing.T) {

	arr := []string{
		"MPRMLS62S21G337",
		"BBBTTT20H12X122",
		"MPRNCR07M53H501",
		"SCHDTL63T49H501",
		"MPR MLS 62S21 G337 ",
		"amnz0459",
	}

	for _, s := range arr {
		s1 := mangling.Rot(s)
		t.Log(s, s1, mangling.Rot(s1))
	}
}
