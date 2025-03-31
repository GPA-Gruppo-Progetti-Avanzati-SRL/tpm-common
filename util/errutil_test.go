package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"testing"
)

func TestErrorRandomizer(t *testing.T) {

	rnd := util.ErrorRandomizer(util.NewErrorRandomizer(5000, 1))
	for i := 0; i < 5000; i++ {
		if rnd.GenerateRandomError() != nil {
			t.Log("Error from GenerateRandomError at i:", i)
		}
	}

}
