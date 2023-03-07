package csvreader

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/textfile"
	"io"
	"strings"
)

type Config struct {
	HeaderLine bool
	Separator  rune
	FileName   string
	Fields     []textfile.FieldInfo
	ioReader   io.Reader
}

type Option func(cfg *Config)

func WithHeaderLine(b bool) Option {
	return func(cfg *Config) {
		cfg.HeaderLine = b
	}
}

func WithSeparator(aSep rune) Option {
	return func(cfg *Config) {
		cfg.Separator = aSep
	}
}

func WithIoReader(reader io.Reader) Option {
	return func(cfg *Config) {
		cfg.ioReader = reader
	}
}

func WithFilename(fn string) Option {
	return func(cfg *Config) {
		cfg.FileName = fn
	}
}

func WithFields(fi []textfile.FieldInfo) Option {
	return func(cfg *Config) {
		cfg.Fields = fi
	}
}

func (c *Config) AdjustFieldIndexes(fs []string) {

	for i, f := range c.Fields {
		if len(fs) == 0 {
			c.Fields[i].Index = i
		} else {
			for j := range fs {
				if strings.ToLower(fs[j]) == strings.ToLower(f.Name) {
					c.Fields[i].Index = j
					break
				}
			}
		}
	}

}
