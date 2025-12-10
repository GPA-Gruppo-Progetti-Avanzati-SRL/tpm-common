package writer

import (
	"errors"
	"io"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"github.com/rs/zerolog/log"
)

type RecordConfig struct {
	Key      string                                       `yaml:"key,omitempty" mapstructure:"key,omitempty" json:"key,omitempty"`
	Fields   []fixedlengthfile.FixedLengthFieldDefinition `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	fieldMap map[string]int
}

type Config struct {
	FileName              string         `yaml:"filename,omitempty" mapstructure:"filename,omitempty" json:"filename,omitempty"`
	ForgiveOnMissingField bool           `yaml:"forgive-on-missing-fields,omitempty" mapstructure:"forgive-on-missing-fields,omitempty" json:"forgive-on-missing-fields,omitempty"`
	Records               []RecordConfig `yaml:"records,omitempty" mapstructure:"records,omitempty" json:"records,omitempty"`
	ioWriter              io.Writer

	// HeadFields            []fixedlengthfile.FixedLengthFieldDefinition `yaml:"h-fields,omitempty" mapstructure:"h-fields,omitempty" json:"h-fields,omitempty"`
	// Fields                []fixedlengthfile.FixedLengthFieldDefinition `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	// TailFields            []fixedlengthfile.FixedLengthFieldDefinition `yaml:"t-fields,omitempty" mapstructure:"t-fields,omitempty" json:"t-fields,omitempty"`

}

type Option func(cfg *Config)

func WithIoWriter(writer io.Writer) Option {
	return func(cfg *Config) {
		cfg.ioWriter = writer
	}
}

func WithFilename(fn string) Option {
	return func(cfg *Config) {
		cfg.FileName = fn
	}
}

func WithRecord(recCfg RecordConfig) Option {
	return func(cfg *Config) {
		flds, _ := adjustFieldInfoIndex(recCfg.Fields)
		cfg.Records = append(cfg.Records, RecordConfig{Key: recCfg.Key, Fields: flds})
	}
}

/*func WithFields(fi []fixedlengthfile.FixedLengthFieldDefinition) Option {
	return func(cfg *Config) {
		cfg.Fields, _ = adjustFieldInfoIndex(fi)
	}
}

func WithHeadFields(fi []fixedlengthfile.FixedLengthFieldDefinition) Option {
	return func(cfg *Config) {
		cfg.HeadFields, _ = adjustFieldInfoIndex(fi)
	}
}

func WithTailFields(fi []fixedlengthfile.FixedLengthFieldDefinition) Option {
	return func(cfg *Config) {
		cfg.TailFields, _ = adjustFieldInfoIndex(fi)
	}
}*/

func adjustFieldInfoIndex(fields []fixedlengthfile.FixedLengthFieldDefinition) ([]fixedlengthfile.FixedLengthFieldDefinition, error) {
	const semLogContext = "fixed-length-writer::adjust-field-indexes"

	recordLength := -1
	for i, f := range fields {

		/*
			if f.Offset != (recordLength + 1) {
				err := errors.New("fields have to be provided sorted by offset")
				log.Error().Err(err).Msg(semLogContext)
				return fields, err
			}
		*/

		if f.Length <= 0 {
			err := errors.New("fields have to be length greater than zero")
			log.Error().Err(err).Msg(semLogContext)
			return fields, err
		}

		fields[i].Offset = recordLength + 1
		fields[i].Index = i
		recordLength += f.Length
	}

	return fields, nil
}

func (cfg Config) ResolveConfig(opts ...Option) (*Config, error) {
	const semLogContext = "fixed-length-writer::resolve-config"
	var err error

	config := cfg
	for _, o := range opts {
		o(&config)
	}

	var adjustedRecords []RecordConfig
	for _, rec := range config.Records {
		var activeFields []fixedlengthfile.FixedLengthFieldDefinition
		for _, field := range rec.Fields {
			if !field.Disabled {
				activeFields = append(activeFields, field)
			}
		}

		if len(activeFields) == 0 {
			err = errors.New("no active fields configured")
			log.Error().Err(err).Msg(semLogContext)
			return nil, err
		}

		log.Info().Str("key", rec.Key).Int("number-of-fields", len(activeFields)).Msg(semLogContext)

		adjustedRecords = append(adjustedRecords, RecordConfig{
			Key:      rec.Key,
			Fields:   activeFields,
			fieldMap: computeFieldMap(activeFields),
		})
	}

	config.Records = adjustedRecords
	if config.ioWriter == nil && config.FileName == "" {
		err = errors.New("please provide a writer or filename")
		log.Error().Err(err).Msg(semLogContext)
		return nil, err
	}

	return &config, nil
}
