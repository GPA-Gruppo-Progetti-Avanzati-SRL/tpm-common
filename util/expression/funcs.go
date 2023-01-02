package expression

import (
	"strings"
)

func BuiltinFuncMap() map[string]interface{} {
	builtins := make(map[string]interface{})
	builtins["isDef"] = IsDefined
	return builtins
}

var expressionSmell = []string{
	"now",
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
func IsExpression(e string) bool {
	if e == "" {
		return false
	}

	for _, s := range expressionSmell {
		if strings.Contains(e, s) {
			return true
		}
	}

	return false
}

func IsDefined(variable interface{}) bool {
	return variable != nil
}
