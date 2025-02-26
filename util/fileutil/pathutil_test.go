package fileutil_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/stretchr/testify/require"
	"path/filepath"
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
		ph := fileutil.ListPathHierarchy(p, true)
		t.Log(len(ph), ph)

		t.Logf("[%d] DownWard - path: [%s]", i, p)
		ph = fileutil.ListPathHierarchy(p, false)
		t.Log(len(ph), ph)
	}

	ph := fileutil.FindGoModFolder(".")
	require.NotEmpty(t, ph, "could not find go.mod")
	t.Log("go mod folder:", ph)

	ps1 := []string{
		"./relPathWithDot/cpxsequence",
		"~/util/ghostscript",
		"~/Applications",
	}

	for _, p := range ps1 {
		rp, ok := fileutil.ResolvePath(p)
		require.NotEmpty(t, rp, "could not find ", p)
		t.Log(p, " resolved -->", rp, ok)
	}
}

func TestPathRel(t *testing.T) {
	s1arr := []string{
		"docs/ia/arricchimenti",
		"docs/ia/arricchimenti/mongo-db",
		"docs/ia/arricchimenti/pipeline-listini",
		"docs/ia/arricchimenti/pipeline-tana",
		"docs/ia/arricchimenti/pipeline-tcol",
		"docs/ia/arricchimenti/pipeline-trap",
		"docs/ia/ce",
		"docs/ia/lm",
		"docs/ia/lm/simulazione",
		"docs/ia/lm/sottoscrizione-old",
	}

	s2arr := []string{
		"docs/ia/arricchimenti",
		"docs/ia/arricchimenti/mongo-db",
		"docs/ia/arricchimenti/pipeline-listini",
		"docs/ia",
		"docs/ia/ce",
		"docs/ia/lm",
		"docs/ia/lm/simulazione",
	}

	for i, p := range s2arr {
		t.Logf("[%d] path: %s", i, p)
		for j, nodeLoc := range s1arr {
			rel, err := filepath.Rel(nodeLoc, p)
			require.NoError(t, err)
			t.Logf("[%d] node-loc: %s, rel: %s", j, nodeLoc, rel)
		}

	}
}
