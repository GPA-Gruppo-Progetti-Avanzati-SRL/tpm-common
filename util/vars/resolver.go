package varResolver

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type VariableReference struct {
	RefType VariableReferenceType
	Match   string
	VarName string
}

type VariableReferenceType string

const (
	AnyVariableReference  VariableReferenceType = "any"
	nullVariableReference VariableReferenceType = "null"

	PercentVariableReference       VariableReferenceType = "percent"
	PercentVariableReferencePrefix string                = "<%="
	PercentVariableReferenceSuffix string                = "%>"

	DashVariableReference       VariableReferenceType = "dash"
	DashVariableReferencePrefix string                = "<#="
	DashVariableReferenceSuffix string                = "#>"

	DollarVariableReference       VariableReferenceType = "dollar"
	DollarVariableReferencePrefix string                = "${"
	DollarVariableReferenceSuffix string                = "}"

	SimpleVariableReference       VariableReferenceType = "simple"
	SimpleVariableReferencePrefix string                = "{"
	SimpleVariableReferenceSuffix string                = "}"

	SuffixErrorMessage = "suffix %s doesn't match prefix %s"
)

// VariableReferencePatternRegexp sort of strict mode with names of vars starting with letters followed by letters, digits and the chars ':', '_', '-'
var VariableReferencePatternRegexp = regexp.MustCompile("((?:<[%#]=)|(?:\\$\\{)|{)([a-zA-Z][\\:a-zA-Z0-9_\\-,]*)([%#]>|})")

// VariableReferencePatternRegexpExt sort of extended mode with names of vars starting with letters or the dollar sign followed by more possible chars.
// Tried to include symbols from https://goessner.net/articles/JsonPath/
var VariableReferencePatternRegexpExt = regexp.MustCompile("((?:<[%#]=)|(?:\\$\\{)|{)(!?[$a-zA-Z][:,@'$\\.\\\"\\[\\]a-zA-Z0-9_\\-]*)([%#]>|})")

type PrefixSuffixTypeMapping struct {
	Type   VariableReferenceType
	Prefix string
	Suffix string
}

var PrefixMap = map[string]PrefixSuffixTypeMapping{
	PercentVariableReferencePrefix: {
		PercentVariableReference, PercentVariableReferencePrefix, PercentVariableReferenceSuffix,
	},
	DashVariableReferencePrefix: {
		DashVariableReference, DashVariableReferencePrefix, DashVariableReferenceSuffix,
	},
	DollarVariableReferencePrefix: {
		DollarVariableReference, DollarVariableReferencePrefix, DollarVariableReferenceSuffix,
	},
	SimpleVariableReferencePrefix: {
		SimpleVariableReference, SimpleVariableReferencePrefix, SimpleVariableReferenceSuffix,
	},
}

func FindVariableReferences(s string, ofType VariableReferenceType) ([]VariableReference, error) {
	matches := VariableReferencePatternRegexpExt.FindAllSubmatch([]byte(s), -1)

	var resp []VariableReference
	for _, m := range matches {

		pfix := string(m[1])
		varname := string(m[2])
		sfix := string(m[3])

		refType, ok := PrefixMap[pfix]
		if !ok {
			return nil, fmt.Errorf("cannot find a match for prefix %s", pfix)
		}

		requiredSuffix := refType.Suffix
		if sfix != requiredSuffix {
			return nil, fmt.Errorf("suffix %s doesn't match prefix %s", sfix, pfix)
		}

		if refType.Type == ofType || ofType == AnyVariableReference {
			resp = append(resp, VariableReference{RefType: refType.Type, Match: string(m[0]), VarName: varname})
		}

	}

	return resp, nil
}

type VariableResolverFunc func(current string, s string) string

func ResolveVariables(s string, ofType VariableReferenceType, aResolver VariableResolverFunc, trimResult bool) (string, error) {

	if s == "" {
		return s, nil
	}

	vars, err := FindVariableReferences(s, ofType)
	if err != nil || len(vars) == 0 {
		return s, err
	}

	for _, v := range vars {
		s = strings.ReplaceAll(s, v.Match, aResolver(s, v.VarName))
	}

	return strings.TrimSpace(s), nil
}

func SimpleMapResolver(m map[string]interface{}, onVarNotFound string) func(a, s string) string {
	return func(a, s string) string {

		tags := strings.Split(s, ",")
		var v interface{}
		var ok bool
		if v, ok = m[tags[0]]; !ok {
			v = onVarNotFound
		}

		if f, ok := v.(func(s string) string); ok {
			v = f(a)
		}

		opts := struct {
			quoted     bool
			formatType string
			format     string
		}{
			quoted:     false,
			formatType: "sprint",
			format:     "",
		}

		for i := 1; i < len(tags); i++ {
			switch tags[i] {
			case "quoted":
				opts.quoted = true
			default:
				if _, ok = v.(time.Time); ok {
					opts.format = tags[i]
					opts.formatType = "time-layout"
				} else {
					opts.format = "%" + tags[i]
					opts.formatType = "sprintf"
				}
			}
		}

		var res string
		switch opts.formatType {
		case "time-layout":
			res = v.(time.Time).Format(opts.format)
		case "sprintf":
			res = fmt.Sprintf(opts.format, v)
		default:
			res = fmt.Sprint(v)
		}
		if opts.format != "" {

		} else {

		}

		if opts.quoted {
			res = fmt.Sprintf("\"%s\"", res)
		}

		return res
	}
}
