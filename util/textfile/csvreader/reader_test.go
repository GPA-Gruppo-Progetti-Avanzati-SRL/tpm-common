package csvreader_test

import (
	"bytes"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile/csvreader"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var buffers []bytes.Buffer

	var buffer bytes.Buffer
	buffer.WriteString("campaign;ssn;first-name;last-name;email")
	buffer.WriteString("\r\n")
	for i := 0; i < 1; i++ {
		buffer.WriteString(`BPMIFI;SSN-CODE;Theodore Jr;Smith;ted.smith@gmail.com`)
		buffer.WriteString("\r\n")
	}
	buffers = append(buffers, buffer)

	var buffer2 bytes.Buffer
	for i := 0; i < 1; i++ {
		buffer2.WriteString(`BPMIFI;SSN-CODE;Theodore Jr;Smith;ted.smith@gmail.com`)
		buffer2.WriteString("\r\n")
	}
	buffers = append(buffers, buffer2)

	cfgs := []csvreader.Config{
		{HeaderLine: true, Separator: ';', Fields: []textfile.FieldInfo{{Name: "campaign", Validation: "required"}, {Name: "email", Validation: "email"}}},
		{HeaderLine: false, Separator: ';'},
	}

	for i, cfg := range cfgs {
		r, err := csvreader.NewCSVFileReader(cfg, csvreader.WithFields(cfg.Fields), csvreader.WithIoReader(&buffers[i]))
		require.NoError(t, err)

		parsed, err := r.Read()
		for err == nil {
			t.Log("parsed-line: ", parsed)

			parsed, err = r.Read()
		}

		require.Equal(t, io.EOF, err, "got: "+err.Error())
	}

}
