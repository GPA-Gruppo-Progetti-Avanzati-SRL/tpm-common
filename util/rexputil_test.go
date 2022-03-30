package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"regexp"
	"testing"
)

var ServersUrlPattern = regexp.MustCompile(`^(?:http|https)://[0-9a-zA-Z\.]*(?:\:[0-9]{2,4})?(.*)`)

func TestExtractCapturedGroupIfMatch(t *testing.T) {

	sarr := []string{"http://localhost:8080/bpap-servizi-pp", "/api/v1"}
	for _, s := range sarr {
		sub := util.ExtractCapturedGroupIfMatch(ServersUrlPattern, s)
		t.Log(sub)
	}

}
