package funcs

import (
	"github.com/PaesslerAG/gval"
	"strings"
	"time"
)

func Builtins() map[string]interface{} {
	builtins := make(map[string]interface{})
	builtins["isDef"] = IsDefined
	builtins["_now"] = time.Now
	builtins["_nowAfterDuration"] = func(dur, fmt string) string {
		v, _ := NowAfterDuration(dur, fmt)
		return v
	}
	return builtins
}

func GValFunctions() []gval.Language {
	gvalFuncs := []gval.Language{
		gval.Function("_nowAfterDuration", func(dur string, fmt string) (string, error) {
			return NowAfterDuration(dur, fmt)
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

func IsDefined(variable interface{}) bool {
	return variable != nil
}
