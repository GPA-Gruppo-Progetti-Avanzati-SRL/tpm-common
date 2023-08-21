package fixedlengthfile

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type FixedLengthFieldDefinition struct {
	Id     string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Name   string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Offset int    `yaml:"offset,omitempty" mapstructure:"offset,omitempty" json:"offset,omitempty"`
	Length int    `yaml:"length,omitempty" mapstructure:"length,omitempty" json:"length,omitempty"`
	Help   string `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	Index  int    `yaml:"index,omitempty" mapstructure:"index,omitempty" json:"index,omitempty"`
}

type FixedLengthRecordMode string

const (
	FixedLengthRecordModeExact   = "exact"
	FixedLengthRecordModeAtLeast = "at-least"
	FixedLengthRecordModeAny     = "any"
)

type FixedLengthRecordDefinition struct {
	Id                  string                       `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Fields              []FixedLengthFieldDefinition `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	Len                 int                          `yaml:"len,omitempty" mapstructure:"len,omitempty" json:"len,omitempty"`
	LengthMode          FixedLengthRecordMode        `yaml:"length-mode,omitempty" mapstructure:"length-mode,omitempty" json:"length-mode,omitempty"`
	PrefixDiscriminator string                       `yaml:"prefix-discriminator,omitempty" mapstructure:"prefix-discriminator,omitempty" json:"prefix-discriminator,omitempty"`
	FieldMap            map[string]int
}

func (r *FixedLengthRecordDefinition) AdjustFieldInfoIndex() error {
	const semLogContext = "fixed-length-record::adjust-field-indexes"

	recordLength := -1
	for i, f := range r.Fields {

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
			return err
		}

		r.Fields[i].Offset = recordLength + 1
		r.Fields[i].Index = i
		recordLength += f.Length
	}

	// recordLength started from -1 to compute correct offsets...
	r.Len = recordLength + 1
	return nil
}

func (r *FixedLengthRecordDefinition) ComputeFieldMap() {
	fieldMap := make(map[string]int)
	for i, f := range r.Fields {
		fId := f.Id
		if fId == "" {
			fId = f.Name
		}
		fieldMap[fId] = i
	}

	r.FieldMap = fieldMap
	return
}

func (r *FixedLengthRecordDefinition) ValidateLineLength(lineno int, l string) error {
	switch r.LengthMode {
	case FixedLengthRecordModeAtLeast:
		if len(l) < r.Len {
			err := fmt.Errorf("line length expected (%d) greater than actual (%d) for line at %d", r.Len, len(l), lineno)
			return err
		}
	case FixedLengthRecordModeAny:
	default:
		if len(l) != r.Len {
			err := fmt.Errorf("line length expected (%d) different than actual (%d) for line at %d", r.Len, len(l), lineno)
			return err
		}
	}

	return nil
}
