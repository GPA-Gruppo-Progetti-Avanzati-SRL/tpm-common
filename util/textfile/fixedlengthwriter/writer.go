package fixedlengthwriter

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strings"
)

type Writer interface {
	Close(removeFile bool)
	WriteMap(map[string]interface{}) error
	WriteRecord(Record) error
	NewRecord() Record
	NewHeadRecord() Record
	NewTailRecord() Record
	Filename() string
}

type Record struct {
	csvRecord []string
	fields    []textfile.FixedLengthFieldInfo
	fieldMap  map[string]int
}

func newRecord(fields []textfile.FixedLengthFieldInfo, fieldMap map[string]int) Record {
	return Record{csvRecord: make([]string, len(fields), len(fields)), fields: fields, fieldMap: fieldMap}
}

func computeFieldMap(fields []textfile.FixedLengthFieldInfo) map[string]int {
	fieldMap := make(map[string]int)
	for i, f := range fields {
		fId := f.Id
		if fId == "" {
			fId = f.Name
		}
		fieldMap[fId] = i
	}

	return fieldMap
}

func (r *Record) String() string {
	var sb strings.Builder
	for i := 0; i < len(r.csvRecord); i++ {

		if len(r.csvRecord[i]) < r.fields[i].Length {
			s, _ := util.ToFixedLength(r.csvRecord[i], r.fields[i].Length)
			sb.WriteString(s)
		} else {
			sb.WriteString(r.csvRecord[i])
		}
	}

	return sb.String()
}

func (r *Record) Set(fieldId string, fieldValue interface{}) error {

	const semLogContext = "fixed-length-writer::set-field"

	var s string
	if fieldValue != nil && !(reflect.ValueOf(fieldValue).Kind() == reflect.Ptr && reflect.ValueOf(fieldValue).IsNil()) {
		s = fmt.Sprint(fieldValue)
	}

	if fIndex, ok := r.fieldMap[fieldId]; ok {
		f := r.fields[fIndex]
		s, _ = util.ToFixedLength(s, f.Length)
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
	cfg          Config
	ioWriter     *bufio.Writer
	headFieldMap map[string]int
	fieldMap     map[string]int
	tailFieldMap map[string]int
	osFile       *os.File
	lineNumber   int

	logger util.GeometricTraceLogger
}

func NewWriter(cfg Config, opts ...Option) (Writer, error) {

	const semLogContext = "fixed-length-writer::new"
	var err error

	config := cfg

	for _, o := range opts {
		o(&config)
	}

	if len(config.Fields) == 0 {
		log.Info().Msg(semLogContext + " fields have not been provided")
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

	r.fieldMap = computeFieldMap(config.Fields)
	if len(config.HeadFields) > 0 {
		r.headFieldMap = computeFieldMap(config.HeadFields)
	}
	if len(config.TailFields) > 0 {
		r.tailFieldMap = computeFieldMap(config.TailFields)
	}

	log.Info().Int("number-of-fields", len(config.Fields)).Int("number-of-head-fields", len(config.HeadFields)).Int("number-of-tail-fields", len(config.TailFields)).Msg(semLogContext)

	if config.ioWriter != nil {
		r.ioWriter = bufio.NewWriter(config.ioWriter)
	} else {
		r.osFile, err = os.Create(config.FileName)
		if err != nil {
			return nil, err
		}

		r.ioWriter = bufio.NewWriter(r.osFile)
	}

	return r, nil
}

func (w *writerImpl) Close(removeFile bool) {

	const semLogContext = "fixed-length-writer::close"
	log.Info().Str("filename", w.cfg.FileName).Bool("remove-file", removeFile).Msg(semLogContext)

	if w.ioWriter != nil {
		_ = w.ioWriter.Flush()
		w.ioWriter = nil
	}

	if w.osFile != nil {
		_ = w.osFile.Close()
		w.osFile = nil
	}

	if removeFile && w.cfg.FileName != "" {
		_ = os.Remove(w.cfg.FileName)
	}
}

func (w *writerImpl) Filename() string {
	return w.cfg.FileName
}

func (w *writerImpl) NewHeadRecord() Record {
	return newRecord(w.cfg.HeadFields, w.headFieldMap)
}

func (w *writerImpl) NewRecord() Record {
	return newRecord(w.cfg.Fields, w.fieldMap)
}

func (w *writerImpl) NewTailRecord() Record {
	return newRecord(w.cfg.TailFields, w.tailFieldMap)
}

func (w *writerImpl) WriteRecord(rec Record) error {
	_, err := w.ioWriter.WriteString(rec.String())
	if err != nil {
		return err
	}
	_, err = w.ioWriter.WriteRune('\n')
	return err
}

func (w *writerImpl) WriteMap(m map[string]interface{}) error {
	panic(errors.New("not implemented record"))
}

//func checkFieldInfo(fields []textfile.FixedLengthFieldInfo) (map[string]int, error) {
//	const semLogContext = "fixed-length-writer::new"
//	fieldMap := make(map[string]int)
//	recordLength := -1
//	for i, f := range fields {
//
//		if f.Offset != (recordLength + 1) {
//			err := errors.New("fields have to be provided sorted by offset")
//			log.Error().Err(err).Msg(semLogContext)
//			return nil, err
//		}
//
//		if f.Length <= 0 {
//			err := errors.New("fields have to be length greater than zero")
//			log.Error().Err(err).Msg(semLogContext)
//			return nil, err
//		}
//
//		fId := f.Id
//		if fId == "" {
//			fId = f.Name
//		}
//		fieldMap[fId] = i
//
//		recordLength += f.Length
//	}
//
//	return fieldMap, nil
//}
