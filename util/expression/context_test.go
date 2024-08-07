package expression_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var j = []byte(`
{
  "can-ale": "APPP",
  "beneficiario": {
    "natura": "PP",
    "tipologia": "ALIAS",
    "numero": "8188602",
    "intestazione": "MARIO ROSSI"
  },
  "ordinante": {
    "natura": "DT",
    "tipologia": "ALIAS",
    "numero": "7750602",
    "codiceFiscale": "LPRSPM46H85U177S"
  },
  "operazione": {
    "divisa": "EUR",
    "importo": 0,
    "descrizione": "string",
    "tipo": "RPAU"
  },
  "additionalProperties": {
    "additionalProp1": {},
    "additionalProp2": {},
    "additionalProp3": {}
  },
  "operazioni": [{
      "errori-ope": [{
          "dsc-errore": "mio errore"
      }],
      "pippo": "pluto"
  }]
}
`)

func TestContextEvaluation(t *testing.T) {
	arr := []struct {
		expr     string
		expected interface{}
	}{
		{
			expr:     `{$}`,
			expected: "8188602",
		},
		{
			expr:     `{$.beneficiario.numero}`,
			expected: "8188602",
		},
		{
			expr:     `{$.beneficiario.numero,len=10,pad=0}`,
			expected: "8188602000",
		},
		{
			expr:     `{$.beneficiario.numero,len=-10,pad=0}`,
			expected: "0008188602",
		},
		{
			expr:     `{v:var01,len=10,pad=.}`,
			expected: "OK........",
		},
		{
			expr:     `!e:{v:var01,len=10,pad=.}`,
			expected: "OK........",
		},
		{
			expr:     `!e:{$.missing.props,len=-10,pad=0,defer}`,
			expected: "!e:{$.missing.props,len=-10,pad=0}", // the onf=keep-ref tag gets deleted when evaluated first time
		},
		{
			expr:     `e:_nowAfter("1m", "2006-01-02")`,
			expected: time.Now().Format("2006-01-02"),
		},
	}

	exprCtx, err := expression.NewContext(expression.WithJsonInput(j), expression.WithVars(map[string]interface{}{"var01": "OK"}))
	require.NoError(t, err)

	for i, input := range arr {
		v, err := exprCtx.EvalOne(input.expr)
		require.NoError(t, err, "[%d] error", i)
		require.EqualValues(t, input.expected, v, "[%d] Expected doesn't match actual", i)
	}

}

func TestContextBoolEvaluation(t *testing.T) {

	arr := []struct {
		rules    []string
		expected bool
		mode     expression.EvaluationMode
	}{
		{
			rules:    []string{`"{$.propNotPresent}" == "OK"`},
			expected: false,
			mode:     expression.AllMustMatch,
		},
		{
			rules:    []string{`"{$.beneficiario.natura}" == "DT"`, `"{$.beneficiario.numero}" == "8188602"`},
			expected: true,
			mode:     expression.ExactlyOne,
		},
		{
			rules:    []string{`"{v:var01}" == "OK"`, `"{$.beneficiario.numero}" == "8188602"`},
			expected: true,
			mode:     expression.AllMustMatch,
		},
		{
			rules:    []string{`"{v:var01}" == "OK"`, `"{$.beneficiario.numero}" == "8188602-NO"`},
			expected: false,
			mode:     expression.AllMustMatch,
		},
		{
			rules:    []string{`"{v:varNotPresent}" == "OK"`},
			expected: false,
			mode:     expression.AllMustMatch,
		},
	}

	exprCtx, err := expression.NewContext(expression.WithJsonInput(j), expression.WithVars(map[string]interface{}{"var01": "OK"}))
	require.NoError(t, err)

	for i, input := range arr {
		found, firstFailed, err := exprCtx.BoolEvalMany(input.rules, input.mode)
		require.NoError(t, err)
		require.EqualValues(t, input.expected, found, "Expected doesn't match actual")
		if found {
			t.Logf("[expr:%d] evaluated to true", i)
		} else {
			t.Logf("[expr:%d] evaluated to false, first failed was %d", i, firstFailed)
		}
	}

}
