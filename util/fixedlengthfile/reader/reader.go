package reader

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type Reader interface {
	Close()
	Filename() string
	Read() (Record, error)
	LineNumber() int
}

const (
	EmptyRecordId = "--empty--"
	EofRecordId   = "--eof--"
	ErrRecordId   = "--err--"
)

var EofRecord = Record{
	RecordId: EofRecordId,
}

type Record struct {
	RecordId string   `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	LineNo   int      `yaml:"line-no,omitempty" mapstructure:"line-no,omitempty" json:"line-no,omitempty"`
	Fields   []string `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	fieldMap map[string]int
}

func (r *Record) String() string {

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%d] %s: num-fields: %d", r.LineNo, r.RecordId, len(r.Fields)))
	for k, v := range r.fieldMap {
		sb.WriteString(fmt.Sprintf(" (%s: %s)", k, r.Fields[v]))
	}
	return sb.String()
}

/*
 * Get property methods. A number of variants with the ability to provide a default-value and a simple value mapping.
 */

// KeyValue basic type to support mapping of property values
type KeyValue struct {
	Key   string
	Value string
}

const (
	GetPropertyOtherwiseMappingKey = "--otherwise--"
)

type GerPropertyOptions struct {
	defaultValue        string
	defaultValueOnEmpty bool
	valueMappings       []KeyValue
}

type GetPropertyOption func(opts *GerPropertyOptions)

func WithDefaultValueOnEmpty(b bool) GetPropertyOption {
	return func(opts *GerPropertyOptions) {
		opts.defaultValueOnEmpty = b
	}
}

func WithDefaultValue(v string) GetPropertyOption {
	return func(opts *GerPropertyOptions) {
		opts.defaultValue = v
	}
}

func WithValueMappings(m []KeyValue) GetPropertyOption {
	return func(opts *GerPropertyOptions) {
		opts.valueMappings = m
	}
}

func (r *Record) Get(n string, opts ...GetPropertyOption) string {
	v, _ := r.GetWithIndicator(n, opts...)
	return v
}

func (r *Record) GetWithIndicator(n string, opts ...GetPropertyOption) (string, bool) {
	options := GerPropertyOptions{defaultValueOnEmpty: true}
	for _, o := range opts {
		o(&options)
	}

	v := options.defaultValue
	mappingNdx, ok := r.fieldMap[n]
	if ok {
		v = r.Fields[mappingNdx]
		if v == "" && options.defaultValueOnEmpty {
			v = options.defaultValue
		}
	}

	// The mapping gets into way only in case of resolved values... doesn't get into way for the default value that should be a value already mapped.
	if ok && len(options.valueMappings) > 0 {
		for _, kv := range options.valueMappings {
			kvKey := strings.ToLower(kv.Key)
			if kvKey == GetPropertyOtherwiseMappingKey || strings.ToLower(v) == kvKey {
				v = kv.Value
				break
			}
		}
	}

	return v, ok
}

func (pr *Record) parse(l []byte, definition fixedlengthfile.FixedLengthRecordDefinition) error {

	pr.Fields = make([]string, len(definition.Fields)-definition.NumOfDroppedFields(), len(definition.Fields)-definition.NumOfDroppedFields())
	fndx := 0
	for _, f := range definition.Fields {

		if f.Drop {
			continue
		}

		if f.Offset >= len(l) {
			return nil
		}

		lenField := f.Length
		if (f.Offset + lenField) > len(l) {
			lenField = len(l) - f.Offset
		}

		res := string(l[f.Offset : f.Offset+lenField])
		if f.Trim {
			res = strings.TrimSpace(res)
		}
		pr.Fields[fndx] = res
		fndx++
	}

	return nil
}

func (pr *Record) IsEmpty() bool {
	return pr.RecordId == "" || pr.RecordId == EmptyRecordId
}

type readerImpl struct {
	cfg           Config
	ioReader      *bufio.Reader
	fieldMap      map[string]int
	osFile        *os.File
	lineNumber    int
	discriminator Discriminator
	logger        util.GeometricTraceLogger
}

func (r *readerImpl) LineNumber() int {
	return r.lineNumber
}

func NewReader(cfg Config, opts ...Option) (Reader, error) {

	const semLogContext = "fixed-length-reader::new"
	var err error

	config := cfg

	for _, o := range opts {
		o(&config)
	}

	if len(config.Records) == 0 {
		log.Info().Msg(semLogContext + " records defs have not been provided")
		return nil, errors.New(semLogContext + " Fields configuration have not been provided")
	}

	if config.ioReader == nil && config.FileName == "" {
		err = errors.New("please provide a reader or filename")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	r := &readerImpl{
		cfg:    config,
		logger: util.GeometricTraceLogger{},
	}

	if config.Discriminator == "prefix" {
		r.discriminator = DiscriminatorFunc(PrefixDiscriminator)
	}

	for i := 0; i < len(cfg.Records); i++ {
		err = cfg.Records[i].AdjustFieldInfoIndex()
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		cfg.Records[i].ComputeFieldMap()
		log.Info().Str("rec-id", cfg.Records[i].Id).Int("number-of-fields", len(cfg.Records[i].Fields)).Msg(semLogContext)
	}

	if config.ioReader != nil {
		r.ioReader = bufio.NewReader(config.ioReader)
	} else {
		r.osFile, err = os.Create(config.FileName)
		if err != nil {
			return nil, err
		}

		r.ioReader = bufio.NewReader(r.osFile)
	}

	return r, nil
}

func (w *readerImpl) Close() {

	const semLogContext = "fixed-length-reader::close"
	log.Info().Str("filename", w.cfg.FileName).Msg(semLogContext)

	if w.ioReader != nil {
		w.ioReader = nil
	}

	if w.osFile != nil {
		_ = w.osFile.Close()
		w.osFile = nil
	}

}

func (w *readerImpl) Filename() string {
	return w.cfg.FileName
}

func (w *readerImpl) Read() (Record, error) {

	r, err := w.read()
	if err == nil {
		switch w.cfg.EmptyLinesMode {
		case EmptyLinesModeSkip:
			for r.IsEmpty() && err == nil {
				r, err = w.read()
			}
		case EmptyLinesModeKeep:
		default:
			if r.IsEmpty() {
				err = fmt.Errorf("empty line found at line %d", w.lineNumber)
			}
		}
	}

	return r, err
}

func (w *readerImpl) read() (Record, error) {

	l, _, err := w.ioReader.ReadLine()
	if err == nil {
		w.lineNumber++

		rId, err := w.discriminateLine(w.lineNumber, string(l))
		if err != nil {
			return Record{RecordId: ErrRecordId, LineNo: w.lineNumber}, err
		}

		// Handling of empty lines is done in the caller... The empty lines are not validated against len. It' kind of specific case.
		if rId == EmptyRecordId {
			return Record{RecordId: EmptyRecordId, LineNo: w.lineNumber}, nil
		}

		r, _ := w.cfg.FindRecordDefinitionById(rId)
		err = r.ValidateLineLength(w.lineNumber, l)
		if err != nil {
			return Record{RecordId: ErrRecordId, LineNo: w.lineNumber}, err
		}

		pr := Record{
			RecordId: rId,
			LineNo:   w.lineNumber,
			Fields:   nil,
			fieldMap: r.FieldMap,
		}

		err = pr.parse(l, r)
		if err != nil {
			return Record{RecordId: ErrRecordId, LineNo: w.lineNumber}, err
		}

		return pr, nil
	}

	return EofRecord, err
}

func (w *readerImpl) discriminateLine(lineno int, l string) (string, error) {

	var err error
	if l == "" {
		return EmptyRecordId, nil
	}

	rId := ErrRecordId
	if w.discriminator != nil {
		rId, err = w.discriminator.DiscriminateLine(lineno, l, w.cfg.Records)
	} else {
		if len(w.cfg.Records) != 1 {
			err = fmt.Errorf("multiple record definitions but no discriminator defined for record at %d line", lineno)
		} else {
			rId = w.cfg.Records[0].Id
		}
	}

	/*
		if err != nil {
			return rId, err
		}

		switch w.cfg.Records[selectedRecordNdx].LengthMode {
		case textfile.FixedLengthRecordModeAtLeast:
			if len(l) < w.cfg.Records[selectedRecordNdx].Len {
				err := fmt.Errorf("line length expected (%d) greater than actual (%d) for line at %d", w.cfg.Records[selectedRecordNdx].Len, len(l), lineno)
				return ErrRecordId, err
			}
		case textfile.FixedLengthRecordModeAny:
		default:
			if len(l) != w.cfg.Records[selectedRecordNdx].Len {
				err := fmt.Errorf("line length expected (%d) different than actual (%d) for line at %d", w.cfg.Records[selectedRecordNdx].Len, len(l), lineno)
				return ErrRecordId, err
			}
		}
	*/

	return rId, err
}
