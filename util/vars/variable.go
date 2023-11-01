package varResolver

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/mangling"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const (
	ReferenceSelf                      = "[variable-ref]"
	KeepReferenceOnNotFoundOptionValue = "keep-ref"
)

type VariableOpts struct {
	Rotate                   bool
	Quoted                   bool
	FormatType               string
	Format                   string
	MaxLength                int
	PadChar                  string
	DefaultValue             interface{}
	KeepVariableReferenceONF bool
}

type Variable struct {
	Name string
	tags []string
}

func ParseVariable(n string) (Variable, error) {
	tags := strings.Split(n, ",")
	v := Variable{Name: tags[0], tags: tags[1:]}
	return v, nil
}

func (vr Variable) Raw() string {
	s := []string{vr.Name}
	s = append(s, vr.tags...)
	return strings.Join(s, ",")
}

func (vr Variable) ToString(v interface{}, jsonEscape bool) (string, error) {

	opts := vr.getOpts(v)
	if v == nil {
		if opts.KeepVariableReferenceONF {
			return ReferenceSelf, nil
		}
		v = opts.DefaultValue
	}

	var res string
	var b []byte
	var err error
	switch opts.FormatType {
	case "time-layout":
		res = v.(time.Time).Format(opts.Format)
	case "sprintf":
		res = fmt.Sprintf(opts.Format, v)
	case "map-json":
		b, err = json.Marshal(v)
		if err == nil {
			res = string(b)
		}
	case "array-json":
		b, err = json.Marshal(v)
		if err == nil {
			res = string(b)
		}
	default:
		res = fmt.Sprint(v)
	}

	if opts.Rotate {
		res = mangling.AlphabetRot(res, true)
	}

	if opts.PadChar != "" {
		res, _ = util.Pad2Length(res, opts.MaxLength, opts.PadChar)
	}

	if opts.MaxLength != 0 {
		res, _ = util.ToMaxLength(res, opts.MaxLength)
	}

	if opts.Quoted {
		res = fmt.Sprintf("\"%s\"", res)
	}

	if jsonEscape {
		res = util.JSONEscape(res)
	}

	return res, nil
}

func (vr Variable) getOpts(value interface{}) VariableOpts {

	const semLogContext = "variable-name::get-opts"

	opts := VariableOpts{
		Rotate:                   false,
		Quoted:                   false,
		FormatType:               "",
		Format:                   "",
		MaxLength:                0,
		PadChar:                  "",
		DefaultValue:             "",
		KeepVariableReferenceONF: false,
	}

	for i := 0; i < len(vr.tags); i++ {
		switch vr.tags[i] {
		case "rotate":
			opts.Rotate = true
		case "quoted":
			opts.Quoted = true
		case "pad":
			opts.PadChar = "0"
		default:
			resolved := false
			if strings.HasPrefix(vr.tags[i], "len=") {
				resolved = true
				v, err := strconv.Atoi(strings.TrimPrefix(vr.tags[i], "len="))
				if err != nil {
					log.Error().Err(err).Msg(semLogContext + " invalid variable tag")
				} else {
					opts.MaxLength = v
				}
			}

			if !resolved && strings.HasPrefix(vr.tags[i], "pad=") {
				resolved = true
				v := strings.TrimPrefix(vr.tags[i], "pad=")
				switch v {
				case "blnk":
					opts.PadChar = " "
				case "":
					log.Warn().Msg(semLogContext + " no pad char provided")
				default:
					opts.PadChar = v[0:1]
				}
			}

			if !resolved && strings.HasPrefix(vr.tags[i], "onf=") {
				resolved = true
				if value == nil {
					v := strings.TrimPrefix(vr.tags[i], "onf=")
					switch v {
					case "now":
						opts.DefaultValue = time.Now()
					case KeepReferenceOnNotFoundOptionValue:
						opts.KeepVariableReferenceONF = true
					default:
						opts.DefaultValue = fmt.Sprint(v)
					}

					value = opts.DefaultValue
				}
			}

			if !resolved {
				switch value.(type) {
				case time.Time:
					opts.Format = vr.tags[i]
					opts.FormatType = "time-layout"
				default:
					opts.Format = "%" + vr.tags[i]
					opts.FormatType = "sprintf"
				}
			}
		}
	}

	if opts.FormatType == "" {
		opts.FormatType = "sprint"
		switch value.(type) {
		case float64, float32:
			opts.Format = "%f"
			opts.FormatType = "sprintf"
		case map[string]interface{}:
			opts.FormatType = "map-json"
		case []interface{}:
			opts.FormatType = "array-json"
		default:
			opts.Format = "%v"
			opts.FormatType = "sprintf"
		}
	}

	return opts
}

/*
func resolveFormatOptions(v interface{}, tags []string) VarReferenceFormatOpts {

	const semLogContext = "common-util-vars::simple-map-resolver"
	var ok bool

	opts := VarReferenceFormatOpts{
		rotate:     false,
		quoted:     false,
		formatType: "sprint",
		format:     "",
		maxLength:  0,
		padChar:    "",
	}

	for i := 1; i < len(tags); i++ {
		switch tags[i] {
		case "rotate":
			opts.rotate = true
		case "quoted":
			opts.quoted = true
		case "pad":
			opts.padChar = "0"
		default:
			resolved := false
			if strings.HasPrefix(tags[i], "len=") {
				resolved = true
				v, err := strconv.Atoi(strings.TrimPrefix(tags[i], "len="))
				if err != nil {
					log.Error().Err(err).Msg(semLogContext + " invalid variable tag")
				} else {
					opts.maxLength = v
				}
			}

			if !resolved && strings.HasPrefix(tags[i], "pad=") {
				resolved = true
				v := strings.TrimPrefix(tags[i], "pad=")
				if len(v) > 0 {
					opts.padChar = v[0:1]
				} else {
					log.Warn().Msg(semLogContext + " no pad char provided")
				}
			}

			if !resolved {
				if _, ok = v.(time.Time); ok {
					opts.format = tags[i]
					opts.formatType = "time-layout"
				} else {
					opts.format = "%" + tags[i]
					opts.formatType = "sprintf"
				}
			}
		}
	}

	return opts
}
*/
