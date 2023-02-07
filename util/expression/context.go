package expression

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	varResolver "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type Option func(r *Context) error

type Context struct {
	vars    map[string]interface{}
	input   map[string]interface{}
	headers NameValuePairs
}

type NameValuePair struct {
	Name    string `json:"name"`              // Name of the pair.
	Value   string `json:"value"`             // Value of the pair.
	Comment string `json:"comment,omitempty"` // A comment provided by the user or the application.
}

type NameValuePairs []NameValuePair

func (nvs NameValuePairs) GetFirst(n string) NameValuePair {

	n = strings.ToLower(n)
	for _, nv := range nvs {
		if strings.ToLower(nv.Name) == n {
			return nv
		}
	}
	return NameValuePair{}
}

func WithHeaders(h []NameValuePair) Option {
	return func(r *Context) error {
		r.headers = h
		return nil
	}
}

func WithVars(m map[string]interface{}) Option {
	return func(r *Context) error {

		if len(r.vars) == 0 {
			r.vars = make(map[string]interface{})
		}

		for n, i := range m {
			r.vars[n] = i
		}

		for n, i := range BuiltinFuncMap() {
			r.vars[n] = i
		}

		return nil
	}
}

func WithFuncMap(fm map[string]interface{}) Option {
	return func(r *Context) error {

		if len(r.vars) == 0 {
			r.vars = make(map[string]interface{})
		}

		for n, i := range fm {
			r.vars[n] = i
		}

		return nil
	}
}

func WithJsonInput(aBody []byte) Option {
	return func(r *Context) error {
		if aBody != nil {
			v := make(map[string]interface{})
			err := json.Unmarshal(aBody, &v)
			if err == nil {
				r.input = v
			} else {
				return err
			}
		}

		return nil
	}
}

func WithMapInput(aBody map[string]interface{}) Option {
	return func(r *Context) error {

		r.input = make(map[string]interface{})
		for n, i := range aBody {
			r.input[n] = i
		}

		return nil
	}
}

func NewContext(opts ...Option) (*Context, error) {
	pvr := &Context{}

	for _, o := range opts {
		err := o(pvr)
		if err != nil {
			return pvr, err
		}
	}

	return pvr, nil
}

type EvaluationMode string

const (
	ExactlyOne   EvaluationMode = "exactly-one"
	AtLeastOne   EvaluationMode = "at-least-one"
	AllMustMatch EvaluationMode = "all-must-match"
)

func (pvr *Context) Add(n string, v interface{}) error {
	if pvr.vars == nil {
		pvr.vars = make(map[string]interface{})
	}

	pvr.vars[n] = v
	return nil
}

func (pvr *Context) EvalOne(v string) (interface{}, error) {

	if v == "" {
		return "", nil
	}

	var err error

	v, err = varResolver.ResolveVariables(v, varResolver.AnyVariableReference, pvr.resolveVar, true)
	if err != nil {
		return "", err
	}

	isExpr := IsExpression(v)
	if isExpr {
		return gval.Evaluate(v, pvr)
	}

	return v, nil
}

func (pvr *Context) BoolEvalMany(varExpressions []string, mode EvaluationMode) (bool, error) {

	if len(varExpressions) == 0 {
		return false, nil
	}

	foundNdx := -1
	for ndx, v := range varExpressions {

		// The empty expression evaluates to true.
		boolVal, err := pvr.BoolEvalOne(v)
		if err != nil {
			return false, err
		}

		if boolVal {
			switch mode {
			case ExactlyOne:
				if foundNdx >= 0 {
					log.Trace().Msgf("expression (%s) at  %d and expression (%s) at %d both evaluate and violate the %s mode",
						varExpressions[foundNdx], foundNdx,
						varExpressions[ndx], ndx,
						mode)
					return false, nil
				}

				foundNdx = ndx
			case AtLeastOne:
				return true, nil
			case AllMustMatch:
				foundNdx = 0
			}
		} else if mode == AllMustMatch {

			log.Trace().Msgf("expression (%s) at  %d evaluate to false and violate the %s mode",
				varExpressions[ndx], ndx,
				mode)
			return false, nil
		}
	}

	if (foundNdx >= 0 && mode == ExactlyOne) || mode == AllMustMatch {
		return true, nil
	}

	return false, nil
}

func (pvr *Context) BoolEvalOne(v string) (bool, error) {

	// The empty expression evaluates to true.
	if v == "" {
		return true, nil
	}

	var err error
	isExpr := IsExpression(v)
	if isExpr {
		v, err = varResolver.ResolveVariables(v, varResolver.AnyVariableReference, pvr.resolveVar, true)
		if err != nil {
			return false, err
		}
	}

	exprValue, err := gval.Evaluate(v, pvr)
	if err != nil {
		return false, err
	}

	boolVal := true
	ok := false
	if boolVal, ok = exprValue.(bool); !ok {
		return false, fmt.Errorf("expression %s is not a boolean expression", v)
	}

	return boolVal, nil
}

var resolverTypePrefix = []string{"$.", "$[", "h:", "p:", "v:"}

func (pvr *Context) resolveVar(_ string, s string) string {

	doEscape := false
	if strings.HasPrefix(s, "!") {
		doEscape = true
		s = strings.TrimPrefix(s, "!")
	}

	pfix, err := pvr.getPrefix(s)
	if err != nil {
		return ""
	}

	switch pfix {
	case "$[":
		fallthrough
	case "$.":
		var v interface{}
		v, err = jsonpath.Get(s, pvr.input)
		// log.Trace().Str("path-name", s).Interface("value", v).Msg("evaluation of var")
		if err == nil {
			s, err = pvr.resolveJsonPathExpr(v)
			if err == nil {
				return pvr.jsonEscape(s, doEscape)
			}
		}

	case "h:":
		s = pvr.headers.GetFirst(s[2:]).Value
		return pvr.jsonEscape(s, doEscape)

	case "v:":
		v, ok := pvr.vars[s[2:]]
		if ok {
			s = fmt.Sprintf("%v", v)
			return pvr.jsonEscape(s, doEscape)
		}

	default:
		v, ok := os.LookupEnv(s)
		if ok {
			return pvr.jsonEscape(v, doEscape)
		}
	}

	log.Info().Str("var-name", s).Msg("could not resolve variable")
	return ""
}

func (pvr *Context) jsonEscape(s string, doEscape bool) string {
	if doEscape {
		s = util.JSONEscape(s)
	}
	return s
}

func (pvr *Context) resolveJsonPathExpr(v interface{}) (string, error) {

	var s string
	var err error
	if v != nil {
		var b []byte
		switch v.(type) {
		case float64, float32:
			s = fmt.Sprintf("%f", v)
		case map[string]interface{}:
			b, err = json.Marshal(v)
			if err == nil {
				s = string(b)
			}
		case []interface{}:
			b, err = json.Marshal(v)
			if err == nil {
				s = string(b)
			}
		default:
			s = fmt.Sprintf("%v", v)
		}
	}

	return s, err
}

func (pvr *Context) getPrefix(s string) (string, error) {

	matchedPrefix := "env"

	for _, pfix := range resolverTypePrefix {
		if strings.HasPrefix(s, pfix) {
			matchedPrefix = pfix
			break
		}
	}

	isValid := false
	switch matchedPrefix {
	case "$[":
		fallthrough
	case "$.":
		if pvr.input != nil {
			isValid = true
		}

	case "v:":
		if pvr.vars != nil {
			isValid = true
		}

	case "h:":
		if pvr.headers != nil {
			isValid = true
		}

	case "env":
		isValid = true
	}

	if !isValid {
		return matchedPrefix, fmt.Errorf("found prefix but resover doesn't have data for resolving")
	}

	return matchedPrefix, nil
}
