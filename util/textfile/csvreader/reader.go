package csvreader

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

type Reader interface {
	Close(removeFile bool)
	Read() (map[string]interface{}, error)
	Filename() string
}

type readerImpl struct {
	cfg       Config
	csvReader *csv.Reader

	osFile     *os.File
	lineNumber int
	isEOF      bool

	logger util.GeometricTraceLogger
}

func NewReader(cfg Config, opts ...Option) (Reader, error) {

	const semLogContext = "csv-reader::new"
	var err error

	config := cfg
	if config.Separator == "" {
		config.Separator = ";"
	}

	for _, o := range opts {
		o(&config)
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

	if config.ioReader != nil {
		r.csvReader = csv.NewReader(config.ioReader)
	} else {
		r.osFile, err = os.Open(config.FileName)
		if err != nil {
			return nil, err
		}

		r.csvReader = csv.NewReader(bufio.NewReader(r.osFile))
	}

	r.csvReader.Comma = rune(config.Separator[0])
	if cfg.HeaderLine {
		fieldNames, err := r.csvReader.Read()
		if err != nil {
			if err == io.EOF {
				log.Trace().Msg(semLogContext + " file empty on reading header line")
				r.isEOF = true
			} else {
				return nil, err
			}
		}

		if len(config.Fields) != 0 {
			log.Info().Msg(semLogContext + " file has header line, field names will be taken from first line")
			r.cfg.AdjustFieldIndexes(fieldNames)
		} else {
			for i, s := range fieldNames {
				r.cfg.Fields = append(r.cfg.Fields, textfile.FieldInfo{Name: s, Index: i})
			}
		}

	} else {
		if len(r.cfg.Fields) == 0 {
			log.Warn().Msg(semLogContext + " please provide field names definition, number of fields not known up front, no header line found, using synth names")
		} else {
			// In this case field indexs are numbered according to fields.... sequential.... At the moment I set also the number of expected fields....
			r.cfg.AdjustFieldIndexes(nil)
			r.csvReader.FieldsPerRecord = len(r.cfg.Fields)
		}
	}

	log.Info().Int("number-of-fields", r.csvReader.FieldsPerRecord).Msg(semLogContext)
	return r, nil
}

func (r *readerImpl) Close(removeFile bool) {
	const semLogContext = "csv-reader::close"
	log.Info().Msg(semLogContext)

	if r.osFile != nil {
		r.osFile.Close()
	}

	if removeFile && r.cfg.FileName != "" {
		os.Remove(r.cfg.FileName)
	}
}

func (w *readerImpl) Filename() string {
	return w.cfg.FileName
}

func (r *readerImpl) Read() (map[string]interface{}, error) {

	const semLogContext = "csv-reader::read"

	validate := validator.New()

	if r.isEOF {
		return nil, io.EOF
	}

	record := make(map[string]interface{})

	fields, err := r.csvReader.Read()
	if err == io.EOF {
		return nil, io.EOF
	}

	// might check for ErrFieldCount... and returning what I got....

	if err != nil {
		return nil, err
	}

	r.lineNumber++

	if r.logger.IsEnabled() {
		r.logger.LogEvent(log.Trace().Int("line-number", r.lineNumber), semLogContext)
	}

	if len(r.cfg.Fields) > 0 {
		var firstErr error
		for i := range r.cfg.Fields {
			if r.cfg.Fields[i].Index >= 0 && r.cfg.Fields[i].Index < len(fields) {
				fieldName := r.cfg.Fields[i].Name
				fieldId := r.cfg.Fields[i].Id
				if fieldId == "" {
					fieldId = fieldName
				}
				fieldValue := fields[r.cfg.Fields[i].Index]
				record[fieldId] = fieldValue
				err = validateField(validate, fieldId, fieldValue, r.cfg.Fields[i].Validation, r.cfg.Fields[i].Help)
				if err != nil {
					log.Error().Err(err).Msg(semLogContext)
					if firstErr == nil {
						firstErr = err
					}
				}
			} else {
				err = fmt.Errorf("field %s not found", r.cfg.Fields[i].Name)
				log.Error().Err(err).Msg(semLogContext)
				if firstErr == nil {
					firstErr = err
				}
			}
		}

		if firstErr != nil {
			return record, firstErr
		}
	} else {
		for i, s := range fields {
			fName := fmt.Sprintf("f%d", i)
			record[fName] = s
		}
	}

	return record, nil
}

func validateField(validate *validator.Validate, fieldName, fieldValue string, rules string, help string) error {
	if rules != "" {
		errs := validate.Var(fieldValue, rules)
		if errs != nil {
			var err error
			switch verr := errs.(type) {
			case validator.ValidationErrors:
				if ferr, ok := verr[0].(validator.FieldError); ok {
					if help != "" {
						err = errors.New(help)
					} else {
						err = fmt.Errorf("property %s of value %s cannot be validated against rule: %s", fieldName, fmt.Sprint(fieldValue), ferr.Tag())
					}
				} else {
					err = fmt.Errorf("property %s of value %s cannot be validated against rule: %s", fieldName, fmt.Sprint(fieldValue), rules)
				}
			case *validator.InvalidValidationError:
				err = fmt.Errorf("property %s of value %s cannot be validated against an INVALID rule: %s", fieldName, fmt.Sprint(fieldValue), rules)
			default:
				err = fmt.Errorf("property %s of value %s cannot be validated (with unknown error) against rule: %s", fieldName, fmt.Sprint(fieldValue), rules)
			}

			return err
		}
	}

	return nil
}
