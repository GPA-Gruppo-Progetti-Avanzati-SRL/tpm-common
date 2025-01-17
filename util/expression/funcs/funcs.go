package funcs

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/PaesslerAG/gval"
	"strings"
)

func Builtins() map[string]interface{} {
	builtins := make(map[string]interface{})
	// builtins["isDef"] = IsDefined
	//builtins["_now"] = func(fmt string) string {
	//	v, _ := Now(fmt)
	//	return v
	//}

	builtins["_nowAfter"] = func(dur, fmt string) string {
		v, _ := NowAfter(dur, fmt)
		return v
	}

	builtins["_pad"] = func(s string, len int) string {
		v, _ := util.Pad2Length(s, len, "0")
		return v
	}

	builtins["now"] = Now
	builtins["age"] = Age
	builtins["isDate"] = IsDate
	builtins["parseDate"] = ParseDate
	builtins["parseAndFormatDate"] = ParseAndFmtDate
	builtins["dateDiff"] = DateDiff
	builtins["printf"] = Printf
	builtins["amtConv"] = AmtConv
	builtins["amtCmp"] = AmtCmp
	builtins["amtAdd"] = AmtAdd
	builtins["amtDiff"] = AmtDiff
	builtins["padLeft"] = PadLeft
	builtins["left"] = Left
	builtins["right"] = Right
	builtins["len"] = Len
	builtins["substr"] = Substr
	builtins["isDef"] = IsDefined
	builtins["b64"] = Base64
	builtins["uuid"] = Uuid
	builtins["regexMatch"] = RegexMatch
	builtins["regexExtractFirst"] = RegexExtractFirst
	builtins["lenJsonArray"] = LenJsonArray
	builtins["isJsonArray"] = IsJsonArray
	builtins["stringIn"] = StringIn
	return builtins
}

// GValFunctions Not used actually. Not needed for expression. The funcMap should be enough
func GValFunctions() []gval.Language {
	gvalFuncs := []gval.Language{
		gval.Function("_nowAfter", func(dur string, fmt string) (string, error) {
			return NowAfter(dur, fmt)
		}),
		gval.Function("_now", func(fmt string) (string, error) {
			return Now(fmt), nil
		}),
		gval.Function("_pad", func(s string, len int) (string, error) {
			v, _ := util.Pad2Length(s, len, "0")
			return v, nil
		}),
	}

	return gvalFuncs
}

var expressionSmell = []string{
	"_now",
	"isDef",
	">",
	"<",
	"(",
	")",
	"=",
	// "\"",
}

// IsExpression In order not to clutter the process vars assignments in simple cases.... try to detect if this is an expression or not.
// didn't parse the thing but try to find if there is any 'reserved' word in there.
// example: 'hello' is not an expression, '"hello"' is an expression which evaluates to 'hello'. This trick is to avoid something like
// value: '"{$.operazione.commissione}"' in the yamls. Someday I'll get to there.... sure...
func IsExpression(e string) (string, bool) {
	if e == "" {
		return e, false
	}

	if strings.HasPrefix(e, ":") {
		return strings.TrimPrefix(e, ":"), true
	}

	if strings.HasPrefix(e, "e:") {
		return strings.TrimPrefix(e, "e:"), true
	}

	if strings.HasPrefix(e, "!e:") {
		return strings.TrimPrefix(e, "!e:"), false
	}

	for _, s := range expressionSmell {
		if strings.Contains(e, s) {
			return e, true
		}
	}

	return e, false
}
