package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
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

	ph := util.FindGoModFolder(".")
	require.NotEmpty(t, ph, "could not find go.mod")
	t.Log("go mod folder:", ph)

	ps1 := []string{
		"./relPathWithDot/cpxsequence",
		"~/util/ghostscript",
		"~/Applications",
	}

	for _, p := range ps1 {
		rp, err := util.ResolveFolder(p)
		require.NoError(t, err)
		require.NotEmpty(t, rp, "could not find ", p)
		t.Log(p, " resolved -->", rp)
	}
}
