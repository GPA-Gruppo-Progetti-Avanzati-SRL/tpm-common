package varResolver_test

import (
	vars "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindVariableReferences(t *testing.T) {

	sarr := []string{
		"${MY_VAR}",
		"${NOT_DEFINED_VAR}",
		"<%=PERC_VAR%>",
		"<#=DASH_VAR#>",
		"my var: ${MY_VAR}, second one: ${SECOND_ONE}, third var: ${THIRD_ONE}",
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
	}

	for _, s := range sarr {
		s1, err := vars.ResolveVariables(s, vars.DollarVariableReference, func(s string) string { return s }, false)
		require.NoError(t, err)
		t.Logf("string %s resolved to %s", s, s1)
	}

}
