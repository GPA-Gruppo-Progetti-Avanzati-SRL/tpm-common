package util_test

import (
	"github.com/mario-imperato/tpm-common/util"
	"testing"
)

func TestPathUtil(t *testing.T) {

	ps := []string{
		"./relPathWithDot/cpxsequence",
		"relPathWithOutDot/cpxsequence",
		"/absPath/cpxsequence",
		"/absPathWithTrailingSlash/cpxsequence/",
		"/",
		"",
	}

	for i, p := range ps {
		t.Logf("[%d] UpWard - path: [%s]", i, p)
		ph := util.ListPathHierarchy(p, true)
		t.Log(len(ph), ph)

		t.Logf("[%d] DownWard - path: [%s]", i, p)
		ph = util.ListPathHierarchy(p, false)
		t.Log(len(ph), ph)
	}

}
