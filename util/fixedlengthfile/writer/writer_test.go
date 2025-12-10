package writer_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile/writer"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {

	cfg := writer.Config{
		ForgiveOnMissingField: true,
		Records: []fixedlengthfile.FixedLengthRecordDefinition{{
			Id: "r1",
			Fields: []fixedlengthfile.FixedLengthFieldDefinition{
				{
					Id:     "f1",
					Name:   "std field",
					Length: 20,
					Format: fixedlengthfile.FieldFormat{
						PadCharacter: "",
						Alignment:    fixedlengthfile.AlignmentLeft,
						Trim:         true,
					},
				},
				{
					Id:     "f2",
					Name:   "short field",
					Length: 5,
					Format: fixedlengthfile.FieldFormat{
						PadCharacter: "0",
						Alignment:    fixedlengthfile.AlignmentRight,
						Trim:         true,
					},
				},
				{
					Id:       "f3",
					Name:     "disabled field",
					Length:   5,
					Disabled: true,
					Format: fixedlengthfile.FieldFormat{
						PadCharacter: "",
						Alignment:    fixedlengthfile.AlignmentLeft,
						Trim:         true,
					},
				},
			},
		},
		},
	}

	b, err := json.Marshal(cfg)
	require.NoError(t, err)

	t.Log(string(b))

	w, err := writer.NewWriter(cfg, writer.WithIoWriter(os.Stdout))
	require.NoError(t, err)
	defer w.Close(true)

	r, err := w.NewRecord("r1")
	require.NoError(t, err)
	_ = r.Set("f2", "ciao on short field 2")
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r, err = w.NewRecord("r1")
	require.NoError(t, err)
	_ = r.Set("f1", "ciao on std field 1")
	err = w.WriteRecord(r)
	require.NoError(t, err)

	r, err = w.NewRecord("r1")
	require.NoError(t, err)
	_ = r.Set("f1", "hello on field 1")
	_ = r.Set("f3", "is disabled!")
	_ = r.Set("f4", "doesn't exists")
	err = w.WriteRecord(r)
	require.NoError(t, err)
}
