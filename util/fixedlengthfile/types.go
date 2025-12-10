package fixedlengthfile

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
)

type FieldAlignment string

const (
	AlignmentLeft  FieldAlignment = "left"
	AlignmentRight                = "right"
)

type FieldFormat struct {
	PadCharacter string         `yaml:"pad-character,omitempty" mapstructure:"pad-character,omitempty" json:"pad-character,omitempty"`
	Alignment    FieldAlignment `yaml:"alignment,omitempty" mapstructure:"alignment,omitempty" mapstructure:"alignment,omitempty" json:"alignment,omitempty"`
	Trim         bool           `yaml:"trim,omitempty" mapstructure:"trim,omitempty" json:"trim,omitempty"`
	SubLength    int            `yaml:"sub-length,omitempty" mapstructure:"sub-length,omitempty" json:"sub-length,omitempty"`
}

// TrimPrefixPadding Convenience method to manage the deletion of the field UnPadPrefix
func (ff FieldFormat) TrimPrefixPadding() bool {
	if ff.PadCharacter == "0" && ff.Alignment == AlignmentRight {
		return true
	}

	return false
}

type FixedLengthFieldDefinition struct {
	Id     string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Name   string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Offset int    `yaml:"offset,omitempty" mapstructure:"offset,omitempty" json:"offset,omitempty"`
	Length int    `yaml:"length,omitempty" mapstructure:"length,omitempty" json:"length,omitempty"`
	Help   string `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	Index  int    `yaml:"index,omitempty" mapstructure:"index,omitempty" json:"index,omitempty"`
	// Trim   bool   `yaml:"trim,omitempty" mapstructure:"trim,omitempty" json:"trim,omitempty"`
	Type string `yaml:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	// UnPadPrefix string `yaml:"unpad-prefix,omitempty" mapstructure:"unpad-prefix,omitempty" json:"unpad-prefix,omitempty"`
	Drop     bool        `yaml:"drop,omitempty" mapstructure:"drop,omitempty" json:"drop,omitempty"`
	Disabled bool        `yaml:"disabled,omitempty" mapstructure:"disabled,omitempty" json:"disabled,omitempty"`
	Format   FieldFormat `yaml:"format,omitempty" mapstructure:"format,omitempty" json:"format,omitempty"`
	WarnOn   bool        `yaml:"warn-on,omitempty" mapstructure:"warn-on,omitempty" json:"warn-on,omitempty"`
}

func (fd FixedLengthFieldDefinition) Sprintf(value string) string {
	s := value
	if fd.Format.Trim {
		s = strings.TrimSpace(s)
	}

	l := fd.Length
	if fd.Format.Alignment == AlignmentRight {
		l = -l
	}

	pad := fd.Format.PadCharacter
	if pad == "" {
		pad = " "
	}

	var truncated bool
	s, truncated = util.ToFixedLength(s, false, l, pad)
	if truncated {
		log.Warn().Msg("fixed-length-file: truncated value")
	}

	return s
}

func (fd FixedLengthFieldDefinition) Sscanf(value string) string {

	if fd.Format.SubLength > 0 && len(value) > fd.Format.SubLength {
		if fd.Format.Alignment == AlignmentRight {
			value, _ = util.ToMaxLength(value, false, -fd.Format.SubLength)
		} else {
			value, _ = util.ToMaxLength(value, false, fd.Format.SubLength)
		}
	}

	if fd.Format.Trim {
		value = strings.TrimSpace(value)
		if fd.Type == FixedLengthFieldNumeric {
			value = util.TrimPrefixCharacters(value, false, "0")
		}
	}

	return value
}

type FixedLengthRecordMode string

const (
	FixedLengthRecordModeExact   = "exact"
	FixedLengthRecordModeAtLeast = "at-least"
	FixedLengthRecordModeAny     = "any"

	FixedLengthFieldAlpha   = "alpha"
	FixedLengthFieldNumeric = "numeric"
)

type FixedLengthRecordDefinition struct {
	Id                  string                       `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	PrefixDiscriminator string                       `yaml:"prefix-discriminator,omitempty" mapstructure:"prefix-discriminator,omitempty" json:"prefix-discriminator,omitempty"`
	LengthMode          FixedLengthRecordMode        `yaml:"length-mode,omitempty" mapstructure:"length-mode,omitempty" json:"length-mode,omitempty"`
	Fields              []FixedLengthFieldDefinition `yaml:"fields,omitempty" mapstructure:"fields,omitempty" json:"fields,omitempty"`
	Len                 int                          `yaml:"len,omitempty" mapstructure:"len,omitempty" json:"len,omitempty"`
	FieldMap            map[string]int
}

func (r *FixedLengthRecordDefinition) NumOfDroppedFields() int {
	numDropped := 0
	for _, f := range r.Fields {
		if f.Drop {
			numDropped++
		}
	}

	return numDropped
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

	fndx := 0
	for _, f := range r.Fields {

		if f.Drop {
			continue
		}

		fId := f.Id
		if fId == "" {
			fId = f.Name
		}

		fieldMap[fId] = fndx
		fndx++
	}

	r.FieldMap = fieldMap
	return
}

func (r *FixedLengthRecordDefinition) ValidateLineLength(lineno int, l []byte) error {
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
