package fileutil_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
const (

	SymphonyPattern        = "tpm-symphony-openapi\\.(yml|yaml)"
	FolderPath             = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples"
	FolderPath001          = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples/example-001"
	FolderPath001SendNotif = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples/example-001/spm-send-notification"

)
*/

const (
	SymphonyPattern = "tpm-symphony-openapi\\.(yml|yaml)"
	FolderPath      = "/Users/marioa.imperato/projects/tpm/tpm-common/test"
	FolderPath001   = "/Users/marioa.imperato/projects/tpm/tpm-common/test/test-sub-folder-1"
	FolderPath002   = "/Users/marioa.imperato/projects/tpm/tpm-common/test/test-sub-folder-2"
	//FolderPath001SendNotif = "/Users/marioa.imperato/projects/tpm/tpm-symphony/examples/example-001/spm-send-notification"
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
	fs, err = fileutil.FindFiles(FolderPath002, fileutil.WithFindFileType(fileutil.FileTypeFile), fileutil.WithFindOptionIgnoreList([]string{"^\\."}))
	require.NoError(t, err)

	for _, f := range fs {
		t.Log(f)
	}
}

func TestCopyFolder(t *testing.T) {
	src := "~/test"
	dst := "~/test/out-copy"

	_, err := fileutil.CopyFolder(dst, src, fileutil.WithCopyOptionCreateIfMissing(), fileutil.WithCopyOptionIncludeSubFolder())
	require.NoError(t, err)
}
