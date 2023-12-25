package util_test

import (
	"embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed embedded-files/*
var templates embed.FS

func TestFindEmbeddedFiles(t *testing.T) {
	files, err := util.FindEmbeddedFiles(templates,
		"embedded-files",
		util.WithFindOptionNavigateSubDirs(),
		util.WithExcludeRootFolderInNames(),
		// util.WithFindOptionIgnoreList([]string{"\\.txt$"}),                 // exclude text files....
		// util.WithFindOptionIncludeList([]string{"\\.template", "sub-dir"}), // include template files... the point is that dirs are excluded...
	)
	require.NoError(t, err)

	for _, f := range files {
		t.Log(f.Path, f.Info.Name(), f.Info.IsDir())
	}

	//err = util.WalkEmbeddedFS(templates,
	//	"embedded-files",
	//	func(n string, info fs.FileInfo, err error) error {
	//		t.Log(n)
	//		return nil
	//	})
	//
	//require.NoError(t, err)
}
