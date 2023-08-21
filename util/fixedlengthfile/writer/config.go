package writer

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
	"github.com/rs/zerolog/log"
	"io"
)

type Config struct {
	FileName   string                                       `yaml:"filename,omitempty" mapstructure:"filename,omitempty" json:"filename,omitempty"`
	HeadFields []fixedlengthfile.FixedLengthFieldDefinition `yaml:"h-fields,omitempty" mapstructure:"h-fields,omitempty" json:"h-fields,omitempty"`
	Fields     []fixedlengthfile.FixedLengthFieldDefinition `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	TailFields []fixedlengthfile.FixedLengthFieldDefinition `yaml:"t-fields,omitempty" mapstructure:"t-fields,omitempty" json:"t-fields,omitempty"`
	ioWriter   io.Writer
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

func WithFields(fi []fixedlengthfile.FixedLengthFieldDefinition) Option {
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
}

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
