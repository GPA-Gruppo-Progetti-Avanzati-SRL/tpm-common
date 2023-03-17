package csvwriter_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile/csvwriter"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWriter(t *testing.T) {

	cfg := csvwriter.Config{
		HeaderLine: true,
		Separator:  "|",
		FileName:   "",
		Fields: []textfile.FieldInfo{
			{
				Id:   "f1",
				Name: "field01",
			},
			{
				Id:   "f2",
				Name: "field02",
			},
		},
	}

	w, err := csvwriter.NewWriter(cfg, csvwriter.WithIoWriter(os.Stdout))
	require.NoError(t, err)
	defer w.Close()

	r := w.NewRecord()
	r.Set("f2", "ciao on field 2")
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r = w.NewRecord()
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r = w.NewRecord()
	r.Set("f1", "hello on field 1")
	r.Set("f3", "doesn't exists")
	err = w.WriteRecord(r)
	require.NoError(t, err)
}
