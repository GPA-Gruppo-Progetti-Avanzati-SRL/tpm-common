package reader_test

import (
	"bytes"
	_ "embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile/reader"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

//go:embed cbi-rnd-example.txt
var example []byte

func TestNewReader(t *testing.T) {

	cfg := reader.Config{
		FileName:       "",
		Discriminator:  "prefix",
		EmptyLinesMode: reader.EmptyLinesModeKeep,
		Records: []fixedlengthfile.FixedLengthRecordDefinition{
			reader.RHDefinition,
			reader.RHEFDefinition,
			reader.RH61Definition,
			reader.RH62Definition,
			reader.RH63Definition_KKK,
			reader.RH63Definition_YYY,
			reader.RH63Definition_YY2,
			reader.RH63Definition_ZZ1,
			reader.RH63Definition_ZZ2,
			reader.RH63Definition_ZZ3,
			reader.RH63Definition_ID1,
			reader.RH63Definition_RI1,
			reader.RH63Definition_RI2,
			reader.RH63Definition_Else,
			reader.RH64Definition,
			reader.RH65Definition,
		},
	}

	strReader := bytes.NewReader(example)
	reader, err := reader.NewReader(cfg, reader.WithIoReader(strReader))
	require.NoError(t, err)

	r, err := reader.Read()
	require.NoError(t, err)

	for err == nil {
		log.Info().Str("r", r.String()).Send()
		r, err = reader.Read()
	}

	if err != io.EOF {
		t.Log(err)
	}

	reader.Close()
}
