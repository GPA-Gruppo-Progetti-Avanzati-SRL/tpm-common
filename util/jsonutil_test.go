package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"testing"
)

func TestJSONEscape(t *testing.T) {
	s := "My escaped var \"hello\""
	t.Log(util.JSONEscape(s, false))

	s = "{ \"name\": \"hello\"}"
	t.Log(util.JSONEscape(s, false))

	s = "Hello World"
	t.Log(util.JSONEscape(s, false))

}
