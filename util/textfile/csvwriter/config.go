package csvwriter

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"io"
)

type Config struct {
	HeaderLine bool                    `yaml:"header-line,omitempty" mapstructure:"header-line,omitempty" json:"header-line,omitempty"`
	Separator  string                  `yaml:"separator,omitempty" mapstructure:"separator,omitempty" json:"separator,omitempty"`
	FileName   string                  `yaml:"filename,omitempty" mapstructure:"filename,omitempty" json:"filename,omitempty"`
	Fields     []textfile.CSVFieldInfo `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	ioWriter   io.Writer
}

type Option func(cfg *Config)

func WithHeaderLine(b bool) Option {
	return func(cfg *Config) {
		cfg.HeaderLine = b
	}
}

func WithSeparator(aSep string) Option {
	return func(cfg *Config) {
		cfg.Separator = aSep
	}
}

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

func WithFields(fi []textfile.CSVFieldInfo) Option {
	return func(cfg *Config) {
		cfg.Fields = fi
	}
}
