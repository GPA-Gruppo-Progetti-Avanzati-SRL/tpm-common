package fixedlengthwriter_test

import (
	"encoding/json"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile/fixedlengthwriter"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWriter(t *testing.T) {

	cfg := fixedlengthwriter.Config{
		Fields: []textfile.FixedLengthFieldInfo{
			{
				Id:     "f1",
				Name:   "field01",
				Length: 10,
			},
			{
				Id:     "f2",
				Name:   "field02",
				Length: 5,
			},
		},
	}

	b, err := json.Marshal(cfg)
	require.NoError(t, err)

	t.Log(string(b))

	w, err := fixedlengthwriter.NewWriter(cfg, fixedlengthwriter.WithIoWriter(os.Stdout))
	require.NoError(t, err)
	defer w.Close(true)

	r := w.NewRecord()
	_ = r.Set("f2", "ciao on field 2")
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r = w.NewRecord()
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r = w.NewRecord()
	_ = r.Set("f1", "hello on field 1")
	_ = r.Set("f3", "doesn't exists")
	err = w.WriteRecord(r)
	require.NoError(t, err)
}
