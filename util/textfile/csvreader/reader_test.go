package csvreader_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile/csvreader"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
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

	var buffer3 bytes.Buffer
	for i := 0; i < 1; i++ {
		buffer3.WriteString(`BPMIFI|SSN-CODE-2|Theodore Jr|Smith;ted.smith@gmail.com|`)
		buffer3.WriteString("\r\n")
	}
	buffers = append(buffers, buffer3)

	cfgs := []csvreader.Config{
		{HeaderLine: true, Separator: ";", Fields: []textfile.CSVFieldInfo{{Name: "campaign", Validation: "required"}, {Name: "email", Validation: "email"}}},
		{HeaderLine: false, Separator: ";"},
		{HeaderLine: false, Separator: "|"},
	}

	for i, cfg := range cfgs {
		r, err := csvreader.NewReader(cfg, csvreader.WithFields(cfg.Fields), csvreader.WithIoReader(&buffers[i]))
		require.NoError(t, err)

		parsed, err := r.Read()
		for err == nil {
			t.Log("parsed-line: ", parsed)

			parsed, err = r.Read()
		}

		require.Equal(t, io.EOF, err, "got: "+err.Error())
	}
}

func TestValidate(t *testing.T) {
	validate := validator.New()
	record := map[string]interface{}{
		"campaign": "12345",
		"natura":   "PG",
	}

	rules := map[string]interface{}{
		"campaign": "required",
		"natura":   "required,lte=2,oneof=PF PG",
		"descr":    `required_if=natura ZIC`,
	}

	resp := validate.ValidateMap(record, rules)
	t.Log(resp)

}
