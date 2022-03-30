package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	varResolver "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTemplate(t *testing.T) {

	tmplInfo := []util.TemplateInfo{
		{Name: "t1", Content: "Subst: {{.prop_1}}, {{.prop2}}"},
		{Name: "t2", Content: "Subst: {{.prop_1}}, {{.prop2}}"},
	}

	data := make(map[string]interface{})
	data["prop_1"] = "value_1"
	data["prop_2"] = "value_2"

	for _, ti := range tmplInfo {
		tmpl, err := util.ParseTemplates([]util.TemplateInfo{
			{Name: ti.Name, Content: ti.Content},
		}, nil)
		require.NoError(t, err)

		b, err := util.ProcessTemplate(tmpl, data, false)
		require.NoError(t, err)

		t.Log(string(b))
	}

}

func TestPreprocessTemplate(t *testing.T) {

	tmplInfo := []util.TemplateInfo{
		{Name: "t1", Content: "Subst: {{.prop_1}}, <%=prop_2>"},
		{Name: "t2", Content: "Subst: <#=prop_1>, {{.prop2}}"},
	}

	data := make(map[string]interface{})
	data["prop_1"] = "value_1"
	data["prop_2"] = "value_2"

	for _, ti := range tmplInfo {

		tmplContent := util.TemplatePreprocessVariableReferences(ti.Content, varResolver.AnyVariableReference)
		t.Log("pre-processed template content", tmplContent)
		tmpl, err := util.ParseTemplates([]util.TemplateInfo{
			{Name: ti.Name, Content: tmplContent},
		}, nil)
		require.NoError(t, err)

		b, err := util.ProcessTemplate(tmpl, data, false)
		require.NoError(t, err)

		t.Log(string(b))
	}

}
