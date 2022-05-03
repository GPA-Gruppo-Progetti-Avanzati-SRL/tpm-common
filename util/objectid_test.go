package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"testing"
)

func TestNewObjectId(t *testing.T) {
	oid := util.NewObjectId()
	t.Log(oid)
}
