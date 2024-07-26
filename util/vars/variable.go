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
	// ReferenceSelf = "[variable-ref]"
	// OnNotFoundTag                      = "onf"
	// KeepReferenceOnNotFoundOptionValue = "keep-ref"
	// OnNotFoundKeepVariableOption = OnNotFoundTag + "=" + KeepReferenceOnNotFoundOptionValue
	DeferOption = "defer"
)

type VariableOpts struct {
	Rotate       bool
	Quoted       bool
	FormatType   string
	Format       string
	MaxLength    int
	PadChar      string
	DefaultValue interface{}
}

type Variable struct {
	Prefix   VariablePrefix
	Name     string
	Deferred bool
	tags     []string
}

func ParseVariable(n string) (Variable, error) {
	tags := strings.Split(n, ",")

	nm := tags[0]
	pfix := getPrefix(nm, VariablePrefixNotSpecified)
	if pfix != VariablePrefixNotSpecified {
		nm = strings.TrimPrefix(nm, string(pfix))
		if pfix == VariablePrefixDollarSquareBracket {
			nm = strings.TrimPrefix(nm, "\"")
			nm = strings.TrimSuffix(nm, "\"]")
		}
	} else {
		// handle the special case of $ to match the whole object. The case is put back in the VariablePrefixDollarDot case that is handled in the JsonPathName method
		if nm == "$" {
			pfix = VariablePrefixDollarDot
			nm = ""
		}
	}
	v := Variable{Name: nm, Prefix: pfix}

	for _, t := range tags[1:] {
		if t == DeferOption {
			v.Deferred = true
		} else {
			v.tags = append(v.tags, t)
		}
	}

	return v, nil
}

func (vr Variable) JsonPathName() string {
	if vr.Prefix == VariablePrefixDollarDot {
		// Handle the case of $ as the whole object.
		if vr.Name == "" {
			return "$"
		}

		return string(vr.Prefix) + vr.Name
	}

	if vr.Prefix == VariablePrefixDollarSquareBracket {
		return string(vr.Prefix) + "\"" + vr.Name + "\"]"
	}
	return ""
}

func (vr Variable) Raw() string {
	var s []string
	if vr.Prefix != VariablePrefixNotSpecified {
		s = append(s, string(vr.Prefix)+vr.Name)
	} else {
		s = append(s, vr.Name)
	}
	s = append(s, vr.tags...)
	return strings.Join(s, ",")
}

// ToString introduced skipOpts to not interpret unknown properties as sprintf or time layout.
// unknown opts are deprecated.
func (vr Variable) ToString(v interface{}, jsonEscape bool, skipOpts bool) (string, error) {

	opts := vr.getOpts(v, skipOpts)
	if v == nil {
		v = opts.DefaultValue
	}

	var res string
	var b []byte
	var err error
	switch opts.FormatType {
	case FormatTypeTimeLayout:
		res = v.(time.Time).Format(opts.Format)
	case FormatTypeSprintf:
		res = fmt.Sprintf(opts.Format, v)
	case FormatTypeMapJson:
		b, err = json.Marshal(v)
		if err == nil {
			res = string(b)
		}
	case FormatTypeArrayJson:
		b, err = json.Marshal(v)
		if err == nil {
			res = string(b)
		}

	case FormatTypeOnTrue:
		if v != nil {
			res = opts.Format
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

type VariablePrefix string

const (
	// VariablePrefixDollarDot and VariablePrefixDollar jsonpath expressions are of type $. to select properties or simply $ to match the whole object.

	// VariablePrefixDollar              VariablePrefix = "$"

	VariablePrefixNotSpecified        VariablePrefix = "not-present"
	VariablePrefixEnv                 VariablePrefix = "env:"
	VariablePrefixDollarDot           VariablePrefix = "$."
	VariablePrefixVColon              VariablePrefix = "v:"
	VariablePrefixHColon              VariablePrefix = "h:"
	VariablePrefixDollarSquareBracket VariablePrefix = "$["

	FormatOptLen        = "len="
	FormatOptPad        = "pad="
	FormatOptOnf        = "onf="
	FormatOptOnt        = "ont="
	FormatOptSprintf    = "sprf="
	FormatOptTimeLayout = "tml="
	FormatOptRotate     = "rotate"
	FormatOptQuoted     = "quoted"
	FormatOptPadChar    = "pad"

	FormatTypeTimeLayout = "time-layout"
	FormatTypeSprintf    = "sprintf"
	FormatTypeSprint     = "sprint"
	FormatTypeMapJson    = "map-json"
	FormatTypeArrayJson  = "array-json"
	FormatTypeOnFalse    = "on-false"
	FormatTypeOnTrue     = "on-true"
)

var optsMap = map[string]struct{}{
	FormatOptLen:        struct{}{},
	FormatOptPad:        struct{}{},
	FormatOptOnf:        struct{}{},
	FormatOptOnt:        struct{}{},
	FormatOptSprintf:    struct{}{},
	FormatOptTimeLayout: struct{}{},
	FormatOptRotate:     struct{}{},
	FormatOptQuoted:     struct{}{},
	FormatOptPadChar:    struct{}{},
}

func resolveFormatOption(s string) string {
	const semLogContext = "variable::resolve-format-opts"

	if ndx := strings.Index(s, "="); ndx >= 0 {
		s = s[:ndx+1]
	}

	_, ok := optsMap[s]
	if !ok {
		log.Warn().Str("opt", s).Msg(semLogContext + " format option not found")
		return ""
	}

	return s
}

// knownPrefixes VariablePrefixDollarDot and VariablePrefixDollar jsonpath expressions are of type $. to select properties or simply $ to match the whole object.
// it's crucial that in the knownPrefixes array the less specific VariablePrefixDollar be put after the more specific VariablePrefixDollarDot in order not to have a sort of catch all
// effect
var knownPrefixes = []VariablePrefix{VariablePrefixEnv, VariablePrefixDollarDot, VariablePrefixVColon, VariablePrefixHColon, VariablePrefixDollarSquareBracket}

func getPrefix(s string, defaultPrefix VariablePrefix) VariablePrefix {
	for _, pfix := range knownPrefixes {
		if strings.HasPrefix(s, string(pfix)) {
			return pfix
		}
	}

	return defaultPrefix
}

func (vr Variable) getOpts(value interface{}, skipOpts bool) VariableOpts {

	const semLogContext = "variable-name::get-opts"

	opts := VariableOpts{
		Rotate:       false,
		Quoted:       false,
		FormatType:   "",
		Format:       "",
		MaxLength:    0,
		PadChar:      "",
		DefaultValue: "",
	}

	if !skipOpts {
		onFalsePresent := false
		for i := 0; i < len(vr.tags); i++ {

			formatOption := resolveFormatOption(vr.tags[i])

			switch formatOption {
			case FormatOptRotate:
				opts.Rotate = true
			case FormatOptQuoted:
				opts.Quoted = true
			case FormatOptPadChar:
				opts.PadChar = "0"
			case FormatOptLen:

				v, err := strconv.Atoi(strings.TrimPrefix(vr.tags[i], FormatOptLen))
				if err != nil {
					log.Error().Err(err).Msg(semLogContext + " invalid variable tag")
				} else {
					opts.MaxLength = v
				}

			case FormatOptPad:
				v := strings.TrimPrefix(vr.tags[i], FormatOptPad)
				switch v {
				case "blnk":
					opts.PadChar = " "
				case "":
					log.Warn().Msg(semLogContext + " no pad char provided")
				default:
					opts.PadChar = v[0:1]
				}

			case FormatOptOnf:
				if value == nil {
					v := strings.TrimPrefix(vr.tags[i], FormatOptOnf)
					switch v {
					case "now":
						opts.DefaultValue = time.Now()
					case DeferOption:
						// handled in advance.
						// opts.KeepVariableReferenceONF = true
					default:
						opts.DefaultValue = fmt.Sprint(v)
					}

					value = opts.DefaultValue
					onFalsePresent = true
				}

			case FormatOptSprintf:
				v := strings.TrimPrefix(vr.tags[i], FormatOptSprintf)
				opts.Format = "%" + v
				opts.FormatType = FormatTypeSprintf
			case FormatOptTimeLayout:
				v := strings.TrimPrefix(vr.tags[i], FormatOptTimeLayout)
				opts.Format = v
				opts.FormatType = FormatTypeTimeLayout
			case FormatOptOnt:
				if value != nil && !onFalsePresent {
					v := strings.TrimPrefix(vr.tags[i], FormatOptOnt)
					opts.Format = v
					opts.FormatType = FormatTypeOnTrue
				}

			default:
				switch value.(type) {
				case time.Time:
					opts.Format = vr.tags[i]
					opts.FormatType = FormatTypeTimeLayout
				default:
					opts.Format = "%" + vr.tags[i]
					opts.FormatType = FormatTypeSprintf
				}
			}
		}
	}

	if opts.FormatType == "" {
		opts.FormatType = FormatTypeSprint
		switch value.(type) {
		case float64, float32:
			opts.Format = "%f"
			opts.FormatType = FormatTypeSprintf
		case map[string]interface{}:
			opts.FormatType = FormatTypeMapJson
		case []interface{}:
			opts.FormatType = FormatTypeArrayJson
		default:
			opts.Format = "%v"
			opts.FormatType = FormatTypeSprintf
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
