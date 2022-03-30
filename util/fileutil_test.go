package util_test

import (
	"github.com/mario-imperato/tpm-common/util"
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
	fs, err := util.FindFiles(FolderPath, util.WithFindOptionNavigateSubDirs(), util.WithFindFileType(util.FileTypeFile), util.WithFindOptionIncludeList([]string{SymphonyPattern}))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

	t.Log("list only directories no subfolders")
	fs, err = util.FindFiles(FolderPath001, util.WithFindFileType(util.FileTypeDir))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

	t.Log("list only files with exclusions")
	fs, err = util.FindFiles(FolderPath001SendNotif, util.WithFindFileType(util.FileTypeFile), util.WithFindOptionIgnoreList([]string{"^\\."}))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}

}
