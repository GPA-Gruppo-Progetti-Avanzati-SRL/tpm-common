package varResolver_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/expression/funcs"
	vars "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/stretchr/testify/require"
)

func TestFindVariableReferences(t *testing.T) {

	sarr := []string{
		"${MY_VAR}",
		"${NOT_DEFINED_VAR}",
		"<%=PERC_VAR%>",
		"<#=DASH_VAR#>",
		"my var: ${MY_VAR}, second one: ${SECOND_ONE}, third var: ${THIRD_ONE}",
		"{$.myvar}",
		"{$[\"my-var\"]}",
	}

	for _, s := range sarr {
		varr, err := vars.FindVariableReferences(s, vars.AnyVariableReference)
		require.NoError(t, err)

		t.Logf("find variables of type %v in %s --> %d", vars.AnyVariableReference, s, len(varr))
		for i, v := range varr {
			t.Logf("[%d] found var of type %v = %s", i, v.RefType, v.VarName)
		}
	}

}

func TestResolveVariableReferences(t *testing.T) {

	sarr := []string{
		"${MY_VAR}",
		"${NOT_DEFINED_VAR}",
		"<%=PERC_VAR%>",
		"<#=DASH_VAR#>",
		"my var: ${MY_VAR}, second one: ${SECOND_ONE}, third var: ${THIRD_ONE}",
		"${MY_VAR,quoted,08d}",
		"${NOT_DEFINED_VAR,onf=false,ont=true}",
	}

	for _, s := range sarr {
		s1, deferred, err := vars.ResolveVariables(s, vars.DollarVariableReference, func(a, s string) (string, bool) { return s, true }, false)
		require.NoError(t, err)
		t.Logf("string %s deferred(%t) resolved to %s", s, deferred, s1)
	}

	sarr = []string{
		"{v:numero,len=-12,pad=0}",
		"${ctx-id}:${today,20060102}${seq-id,03d}:${check-digit,len=-10,pad=.}",
		"${not-present,onf=now,20060102}",
		"${now,2006-01-02}",
		"${NOT_DEFINED_VAR,onf=false,ont=true}",
		"${MY_VAR,onf=false,ont=true}",
		fmt.Sprintf("${not-present,%s} - ${ctx-id}", vars.DeferOption),
		"${MY_VAR,onf=false,ont=true,quoted-ont}",
		"${MY_VAR_NOT_EXISTENT,onf=false,ont=true,quoted-ont}",
		"${num-rapporto,atoi}",
		"${num-rapporto,onf=null,atoi}",
		"%now,2006-01-02%",
	}

	m := map[string]interface{}{
		"ctx-id":       "BPMIFI",
		"seq-id":       22,
		"MY_VAR":       "MY_VAR_VALUE",
		"today":        time.Now(),
		"now":          time.Now,
		"add-duration": funcs.NowAfter,
		"num-rapporto": "012345",
		"numero":       "1008250111,1008250222,1008250333,1008250444",
		"check-digit": func(a, s string) string {
			a = strings.Replace(a, fmt.Sprintf("${%s}", s), "", -1)
			return fmt.Sprint(len(s))
		},
	}

	for i, s := range sarr {
		s1, deferred, err := vars.ResolveVariables(s, vars.DollarVariableReference, vars.SimpleMapResolver(m), false)
		require.NoError(t, err, "test: %d - %s", i, s)
		t.Logf("[%d] string %s deferred(%t) resolved to %s", i, s, deferred, s1)

		s1, deferred, err = vars.ResolveVariables(s, vars.SimpleVariableReference, vars.SimpleMapResolver(m), false)
		require.NoError(t, err, "test: %d - %s", i, s)
		t.Logf("[%d] string %s deferred(%t) resolved to %s", i, s, deferred, s1)

		s1, deferred, err = vars.ResolveVariables(s, vars.WritersideVariableReference, vars.SimpleMapResolver(m), false)
		require.NoError(t, err, "test: %d - %s", i, s)
		t.Logf("[%d] string %s deferred(%t) resolved to %s", i, s, deferred, s1)
	}
}

func TestVariableValueFormat(t *testing.T) {

	var tests = []struct {
		input  string
		value  string
		wanted string
	}{
		{
			input:  "{$.whatever}",
			value:  "ALM3067 ",
			wanted: "ALM3067",
		},
	}

	for _, s := range tests {
		v, err := vars.ParseVariable(s.input)
		require.NoError(t, err)

		res, err := v.ToString(s.value, true, false)
		require.NoError(t, err)

		require.Equal(t, s.wanted, res)
	}

}
