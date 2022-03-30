package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"testing"
)

func TestMida(t *testing.T) {
	id := "pkg_-0202"
	platform := "GN"

	m, err := util.ComputeMida(platform, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mida: %s", m)
}
