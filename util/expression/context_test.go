package expression_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression"
	"github.com/stretchr/testify/require"
	"testing"
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
			expr:     `{$.beneficiario.numero}`,
			expected: "8188602",
		},
	}

	exprCtx, err := expression.NewContext(expression.WithJsonInput(j), expression.WithVars(map[string]interface{}{"var01": "OK"}))
	require.NoError(t, err)

	for i, input := range arr {
		v, err := exprCtx.EvalOne(input.expr)
		require.NoError(t, err)
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
