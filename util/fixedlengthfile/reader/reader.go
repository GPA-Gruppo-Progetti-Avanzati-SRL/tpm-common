package reader

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"os"
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

func (pr *Record) parse() error {
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
		err = r.ValidateLineLength(w.lineNumber, string(l))
		if err != nil {
			return Record{RecordId: ErrRecordId, LineNo: w.lineNumber}, err
		}

		pr := Record{
			RecordId: rId,
			LineNo:   w.lineNumber,
			Fields:   nil,
			fieldMap: r.FieldMap,
		}

		err = pr.parse()
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
