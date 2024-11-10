package fileutil_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	SymphonyPattern        = "tpm-symphony-openapi\\.(yml|yaml)"
	FolderPath             = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples"
	FolderPath001          = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples/example-001"
	FolderPath001SendNotif = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples/example-001/spm-send-notification"
)

func TestFindFiles(t *testing.T) {

	t.Log("list only files with subfolders filtered by pattern")
	fs, err := fileutil.FindFiles(FolderPath, fileutil.WithFindOptionNavigateSubDirs(), fileutil.WithFindFileType(fileutil.FileTypeFile), fileutil.WithFindOptionIncludeList([]string{SymphonyPattern}))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

	t.Log("list only directories no subfolders")
	fs, err = fileutil.FindFiles(FolderPath001, fileutil.WithFindFileType(fileutil.FileTypeDir))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

	t.Log("list only files with exclusions")
	fs, err = fileutil.FindFiles(FolderPath001SendNotif, fileutil.WithFindFileType(fileutil.FileTypeFile), fileutil.WithFindOptionIgnoreList([]string{"^\\."}))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

}
