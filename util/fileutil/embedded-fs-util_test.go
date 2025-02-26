package fileutil_test

import (
	"embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed embedded-files/*
var templates embed.FS

func TestFindEmbeddedFiles(t *testing.T) {
	files, err := fileutil.FindEmbeddedFiles(templates,
		"embedded-files",
		fileutil.WithFindOptionNavigateSubDirs(),
		fileutil.WithFindOptionExcludeRootFolderInNames(),
		fileutil.WithFindOptionPreloadContent(),
		// util.WithFindOptionIgnoreList([]string{"\\.txt$"}),                 // exclude text files....
		// util.WithFindOptionIncludeList([]string{"\\.template", "sub-dir"}), // include template files... the point is that dirs are excluded...
	)
	require.NoError(t, err)

	for _, f := range files {
		t.Log(f.Path, f.Info.Name(), f.Info.IsDir(), len(f.Content))
	}

	files, err = fileutil.FindEmbeddedFiles(templates,
		"embedded-files/sub-dir",
		fileutil.WithFindOptionNavigateSubDirs(),
		// fileutil.WithFindOptionExcludeRootFolderInNames(),
		fileutil.WithFindOptionTrimRootFolderFromNames(),
		fileutil.WithFindOptionPreloadContent(),
		// util.WithFindOptionIgnoreList([]string{"\\.txt$"}),                 // exclude text files....
		// util.WithFindOptionIncludeList([]string{"\\.template", "sub-dir"}), // include template files... the point is that dirs are excluded...
	)

	require.NoError(t, err)

	for _, f := range files {
		t.Log(f.Path, f.Info.Name(), f.Info.IsDir(), len(f.Content))
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
