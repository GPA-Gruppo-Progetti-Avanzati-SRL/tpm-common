package csvwriter

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
)

type Writer interface {
	Close(removeFile bool)
	WriteMap(map[string]interface{}) error
	WriteRecord(Record) error
	NewRecord() Record
	Filename() string
}

type Record struct {
	csvRecord []string
	fieldMap  map[string]int
}

func NewRecord(numFields int, fieldMap map[string]int) Record {
	return Record{csvRecord: make([]string, numFields, numFields), fieldMap: fieldMap}
}

func (r *Record) Set(fieldId string, fieldValue interface{}) error {

	const semLogContext = "csv-writer::set-field"

	var s string
	if fieldValue != nil && !(reflect.ValueOf(fieldValue).Kind() == reflect.Ptr && reflect.ValueOf(fieldValue).IsNil()) {
		s = fmt.Sprint(fieldValue)
	}

	if fIndex, ok := r.fieldMap[fieldId]; ok {
		r.csvRecord[fIndex] = s
	} else {
		log.Error().Str("field-id", fieldId).Msg(semLogContext + " field not found")
	}

	// At the moment is forgiving.... but the log is in error...
	return nil
}

func (r *Record) Fields() []string {
	return r.csvRecord
}

type writerImpl struct {
	cfg       Config
	csvWriter *csv.Writer
	fieldMap  map[string]int

	osFile     *os.File
	lineNumber int

	logger util.GeometricTraceLogger
}

func NewWriter(cfg Config, opts ...Option) (Writer, error) {

	const semLogContext = "csv-writer::new"
	var err error

	config := cfg
	if config.Separator == "" {
		config.Separator = ";"
	}

	for _, o := range opts {
		o(&config)
	}

	if len(config.Fields) == 0 {
		log.Info().Msg(semLogContext + " file has header line, fields have not been provided")
		return nil, errors.New(semLogContext + " Fields configuration have not been provided")
	}

	if config.ioWriter == nil && config.FileName == "" {
		err = errors.New("please provide a writer or filename")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	r := &writerImpl{
		cfg:    config,
		logger: util.GeometricTraceLogger{},
	}

	if config.ioWriter != nil {
		r.csvWriter = csv.NewWriter(config.ioWriter)
	} else {
		r.osFile, err = os.Create(config.FileName)
		if err != nil {
			return nil, err
		}

		r.csvWriter = csv.NewWriter(bufio.NewWriter(r.osFile))
	}

	r.csvWriter.Comma = rune(config.Separator[0])

	r.fieldMap = make(map[string]int)
	for i, f := range config.Fields {
		fId := f.Id
		if fId == "" {
			fId = f.Name
		}
		r.fieldMap[fId] = i
	}

	if cfg.HeaderLine {
		var record []string
		for _, f := range config.Fields {
			record = append(record, f.Name)
		}

		err = r.csvWriter.Write(record)
		if err != nil {
			return nil, err
		}
	}

	log.Info().Int("number-of-fields", len(config.Fields)).Msg(semLogContext)
	return r, nil
}

func (w *writerImpl) Close(removeFile bool) {

	const semLogContext = "csv-writer::close"
	log.Info().Str("filename", w.cfg.FileName).Bool("remove-file", removeFile).Msg(semLogContext)

	if w.csvWriter != nil {
		w.csvWriter.Flush()
		w.csvWriter = nil
	}

	if removeFile && w.cfg.FileName != "" {
		os.Remove(w.cfg.FileName)
	}
}

func (w *writerImpl) Filename() string {
	return w.cfg.FileName
}

func (w *writerImpl) NewRecord() Record {
	return NewRecord(len(w.cfg.Fields), w.fieldMap)
}

func (w *writerImpl) WriteRecord(rec Record) error {
	return w.csvWriter.Write(rec.csvRecord)
}

func (w *writerImpl) WriteMap(m map[string]interface{}) error {
	panic(errors.New("not implemented record"))
}
