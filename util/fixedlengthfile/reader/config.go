package reader

import (
	"fmt"
	"io"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
)

type EmptyLinesMode string

const (
	EmptyLinesModeErr  = "err"
	EmptyLinesModeSkip = "skip"
	EmptyLinesModeKeep = "keep"

	DiscriminatorModePrefix = "prefix"
)

type Config struct {
	FileName       string                                        `yaml:"filename,omitempty" mapstructure:"filename,omitempty" json:"filename,omitempty"`
	EmptyLinesMode EmptyLinesMode                                `yaml:"empty-lines,omitempty" mapstructure:"empty-lines,omitempty" json:"empty-lines,omitempty"`
	Discriminator  string                                        `yaml:"line-discriminator,omitempty" mapstructure:"line-discriminator,omitempty" json:"line-discriminator,omitempty"`
	Records        []fixedlengthfile.FixedLengthRecordDefinition `yaml:"records,omitempty" mapstructure:"records,omitempty" json:"records,omitempty"`
	ioReader       io.Reader
}

type Option func(cfg *Config)

func WithIoReader(writer io.Reader) Option {
	return func(cfg *Config) {
		cfg.ioReader = writer
	}
}

func WithFilename(fn string) Option {
	return func(cfg *Config) {
		cfg.FileName = fn
	}
}

func WithRecord(r fixedlengthfile.FixedLengthRecordDefinition) Option {
	return func(cfg *Config) {
		cfg.Records = append(cfg.Records, r)
	}
}

func (c Config) FindRecordDefinitionById(id string) (fixedlengthfile.FixedLengthRecordDefinition, error) {
	for _, r := range c.Records {
		if r.Id == id {
			return r, nil
		}
	}

	return fixedlengthfile.FixedLengthRecordDefinition{}, fmt.Errorf("cannot find record by id: %s", id)
}
