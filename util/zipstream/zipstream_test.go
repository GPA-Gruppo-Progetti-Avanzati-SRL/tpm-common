package zipstream_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/zipstream"
	"github.com/stretchr/testify/require"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

var readmeFileContent = []byte(`This is the readme file`)

func TestZipStream(t *testing.T) {

	mzs, err := zipstream.NewZipStream(zipstream.MemZipStream, "")
	require.NoError(t, err)

	err = mzs.Add("readme.txt", readmeFileContent)
	require.NoError(t, err)

	err = mzs.CloseAndSave("/tmp/memstream.zip")
	require.NoError(t, err)

	defer os.Remove("/tmp/memstream.zip")

	err = mzs.Dispose()
	require.NoError(t, err)

	fzs, err := zipstream.NewZipStream(zipstream.DiskZipStream, "/tmp")
	require.NoError(t, err)

	err = fzs.Add("readme.txt", readmeFileContent)
	require.NoError(t, err)

	err = fzs.Close()
	require.NoError(t, err)

	err = mzs.CloseAndSave("/tmp/filestream.zip")
	require.NoError(t, err)

	defer os.Remove("/tmp/filestream.zip")

	err = fzs.Dispose()
	require.NoError(t, err)
}

func TestZipStream2(t *testing.T) {

	zs, err := zipstream.NewZipStream(zipstream.MemZipStream, "")
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ioutil.ReadFile("/tmp/in_gect_ch7.pdf")
	if err != nil {
		t.Fatal(err)
	}

	err = zs.Add("file1.pdf", buf)
	if err != nil {
		t.Fatal(err)
	}

	err = zs.Add("file2.pdf", buf)
	if err != nil {
		t.Fatal(err)
	}

	err = zs.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("my-zipped.zip", zs.Bytes(), fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

}
