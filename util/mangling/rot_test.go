package mangling_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/mangling"
	"testing"
)

func TestRotation(t *testing.T) {

	arr := []string{
		"MPRNCR07M53H501",
		"2023013100000086",
		"amnz0459",
		"0uvZ",
		"0HIZ",
		"KIKLIJLJIIIIIIQP",
	}

	t.Log("rot-13")
	for _, s := range arr {
		s1 := mangling.Rot13(s)
		t.Log(s, s1, mangling.Rot13(s1))
	}

	t.Log("alphabet-mixed")
	for _, s := range arr {
		s1 := mangling.AlphabetRot(s, false)
		t.Log(s, s1, mangling.AlphabetRot(s1, false))
	}

	t.Log("alphabet-upper")
	for _, s := range arr {
		s1 := mangling.AlphabetRot(s, true)
		t.Log(s, s1, mangling.AlphabetRot(s1, true))
	}

	/* alphabet generation
	var sb strings.Builder
	for c := 'a'; c <= 'z'; c++ {
		sb.WriteRune(c)
	}

	t.Log(sb.String())
	*/
}
