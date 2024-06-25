package expression

import (
	"encoding/json"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression/funcs"
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
	gvals   []gval.Language
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

func WithGValFunctions(fm []gval.Language) Option {
	return func(r *Context) error {

		if len(r.gvals) == 0 {
			r.gvals = make([]gval.Language, 0)
		}

		for _, f := range fm {
			r.gvals = append(r.gvals, f)
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
	pvr.vars = make(map[string]interface{})
	for n, i := range funcs.Builtins() {
		pvr.vars[n] = i
	}

	pvr.gvals = make([]gval.Language, 0)
	for _, f := range funcs.GValFunctions() {
		pvr.gvals = append(pvr.gvals, f)
	}

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

func (pvr *Context) SetInput(n string, v interface{}) error {
	if pvr.input == nil {
		pvr.input = make(map[string]interface{})
	}

	pvr.input[n] = v
	return nil
}

func (pvr *Context) SetVar(n string, v interface{}) error {
	if pvr.vars == nil {
		pvr.vars = make(map[string]interface{})
	}

	pvr.vars[n] = v
	return nil
}

// Add Deprecated
func (pvr *Context) Add(n string, v interface{}) error {
	return pvr.SetVar(n, v)
}

func (pvr *Context) EvalOne(v string) (interface{}, error) {

	if v == "" {
		return "", nil
	}

	var err error
	var deferred bool
	v, deferred, err = varResolver.ResolveVariables(v, varResolver.AnyVariableReference, pvr.resolveVar, true)
	if err != nil {
		return "", err
	}

	isExpr := false
	if !deferred {
		v, isExpr = funcs.IsExpression(v)

		if isExpr {
			return gval.Evaluate(v, pvr.vars, pvr.gvals...)
		}
	}

	return v, nil
}

func (pvr *Context) BoolEvalMany(varExpressions []string, mode EvaluationMode) (bool, int, error) {

	if len(varExpressions) == 0 {
		return false, -1, nil
	}

	foundNdx := -1
	for ndx, v := range varExpressions {

		// The empty expression evaluates to true.
		boolVal, err := pvr.BoolEvalOne(v)
		if err != nil {
			return false, -1, err
		}

		if boolVal {
			switch mode {
			case ExactlyOne:
				if foundNdx >= 0 {
					log.Trace().Msgf("expression (%s) at  %d and expression (%s) at %d both evaluate and violate the %s mode",
						varExpressions[foundNdx], foundNdx,
						varExpressions[ndx], ndx,
						mode)
					return false, -1, nil
				}

				foundNdx = ndx
			case AtLeastOne:
				return true, -1, nil
			case AllMustMatch:
				foundNdx = 0
			}
		} else if mode == AllMustMatch {

			log.Trace().Msgf("expression (%s) at  %d evaluate to false and violate the %s mode",
				varExpressions[ndx], ndx,
				mode)
			return false, ndx, nil
		}
	}

	if (foundNdx >= 0 && mode == ExactlyOne) || mode == AllMustMatch {
		return true, -1, nil
	}

	return false, -1, nil
}

func (pvr *Context) BoolEvalOne(v string) (bool, error) {

	const semLogContext = "expression-ctx::bool-eval-one"
	// The empty expression evaluates to true.
	if v == "" {
		return true, nil
	}

	var err error
	var deferred bool

	// Current formulation seems to be wrong.... variables are resolved only if it's an expression...
	// if isExpr {
	v, deferred, err = varResolver.ResolveVariables(v, varResolver.AnyVariableReference, pvr.resolveVar, true)
	if err != nil {
		log.Error().Err(err).Str("expr", v).Msg(semLogContext)
		return false, err
	}

	if deferred {
		err = fmt.Errorf("expression deferred: %s", v)
		log.Error().Err(err).Str("expr", v).Msg(semLogContext)
		return false, err
	}

	//}

	/* Need to risk the evaluation anyway.... the expression check might fail to recognize expression such as 'true' or 'false'
	v, isExpr := IsExpression(v)
	if !isExpr {
		err := fmt.Errorf("expression %s seems not to be evaluable", v)
		log.Error().Err(err).Str("expr", v).Msg(semLogContext)
		return false, err
	}
	*/

	exprValue, err := gval.Evaluate(v, pvr)
	if err != nil {
		log.Error().Err(err).Str("expr", v).Msg(semLogContext)
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

func (pvr *Context) resolveVar(_ string, s string) (string, bool) {

	const semLogContext = "expr-context::resolve-var"

	var err error
	doEscape := false
	if strings.HasPrefix(s, "!") {
		doEscape = true
		s = strings.TrimPrefix(s, "!")
	}

	variable, _ := varResolver.ParseVariable(s)
	if variable.Deferred {
		return variable.Raw(), variable.Deferred
	}

	pfix, err := pvr.validatePrefix(variable.Prefix)
	if err != nil {
		return "", variable.Deferred
	}

	var varValue interface{}

	switch pfix {
	case varResolver.VariablePrefixDollarSquareBracket:
		fallthrough
	case varResolver.VariablePrefixDollarDot:

		varValue, err = jsonpath.Get(variable.Name, pvr.input)
		// log.Trace().Str("path-name", s).Interface("value", v).Msg("evaluation of var")
		/*
			if err == nil {
				s, err = pvr.jsonPathValueToString(varValue)
				if err == nil {
					return pvr.jsonEscape(s, doEscape)
				}
			}
		*/

	case varResolver.VariablePrefixHColon:
		varValue = pvr.headers.GetFirst(variable.Name).Value
		// return pvr.jsonEscape(s, doEscape)

	case varResolver.VariablePrefixVColon:
		varValue, _ = pvr.vars[variable.Name]
		/*
			if ok {
				s = fmt.Sprintf("%v", varValue)
				return pvr.jsonEscape(s, doEscape)
			}
		*/

	default:
		varValue, _ = os.LookupEnv(variable.Name)
		/*
			if ok {
				return pvr.jsonEscape(varValue.(string), doEscape)
			}
		*/
	}

	if err != nil {
		if !isJsonPathUnknownKey(err) {
			log.Error().Err(err).Msg(semLogContext)
			return "", variable.Deferred
		}
	}

	s, err = variable.ToString(varValue, doEscape)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
	}

	return s, variable.Deferred
}

func isJsonPathUnknownKey(err error) bool {
	if err != nil {
		return strings.HasPrefix(err.Error(), "unknown key")
	}

	return false
}

func (pvr *Context) jsonEscape(s string, doEscape bool) string {
	if doEscape {
		s = util.JSONEscape(s)
	}
	return s
}

func (pvr *Context) jsonPathValueToString(v interface{}) (string, error) {

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

func (pvr *Context) validatePrefix(pfix varResolver.VariablePrefix) (varResolver.VariablePrefix, error) {

	if pfix == varResolver.VariablePrefixNotSpecified {
		pfix = varResolver.VariablePrefixEnv
	}

	isValid := false
	switch pfix {
	case varResolver.VariablePrefixDollarSquareBracket:
		fallthrough

	case varResolver.VariablePrefixDollarDot:
		if pvr.input != nil {
			isValid = true
		}

	case varResolver.VariablePrefixVColon:
		if pvr.vars != nil {
			isValid = true
		}

	case varResolver.VariablePrefixHColon:
		if pvr.headers != nil {
			isValid = true
		}

	case varResolver.VariablePrefixEnv:
		isValid = true
	}

	if !isValid {
		return pfix, fmt.Errorf("found prefix but resover doesn't have data for resolving")
	}

	return pfix, nil
}
