package templateutil

import (
	"bytes"
	"errors"
	varResolver "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/vars"
	"github.com/rs/zerolog/log"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Info struct {
	Name    string
	Content string
}

func PreprocessVariableReferences(tmpl string, refType varResolver.VariableReferenceType) string {

	varr, err := varResolver.FindVariableReferences(tmpl, refType)
	if err != nil {
		log.Warn().Err(err).Msg("error in replacing variable references in temp,late")
		return tmpl
	}

	for _, v := range varr {
		var sb strings.Builder
		sb.WriteString("{{.")
		sb.WriteString(v.VarName)
		sb.WriteString("}}")
		tmpl = strings.ReplaceAll(tmpl, v.Match, sb.String())
	}

	return tmpl
}

func EvaluateSimpleTemplateWithVars(t string, meta map[string]interface{}, variablesFilter varResolver.VariableReferenceType) (string, error) {

	realTmpl := PreprocessVariableReferences(t, variablesFilter)

	tmpl := MustParse([]Info{{
		Name:    "evaluate-template",
		Content: realTmpl,
	}}, nil)

	n, err := Process(tmpl, meta, false)
	if err != nil {
		return "", err
	}

	return string(n), nil
}

func Parse(templates []Info, fMaps template.FuncMap) (*template.Template, error) {
	if len(templates) == 0 {
		return nil, errors.New("no template provided")
	}

	mainTemplate := template.New(templates[0].Name)
	if len(fMaps) > 0 {
		mainTemplate = mainTemplate.Funcs(fMaps)
	}

	var err error
	if mainTemplate, err = mainTemplate.Parse(templates[0].Content); err != nil {
		return nil, err
	} else {
		for i := 1; i < len(templates); i++ {
			if _, err = mainTemplate.New(templates[i].Name).Parse(templates[i].Content); err != nil {
				return nil, err
			}
		}
	}

	return mainTemplate, nil
}

func MustParse(templates []Info, fMaps template.FuncMap) *template.Template {
	t, err := Parse(templates, fMaps)
	if err != nil {
		panic(err)
	}

	return t
}

/*
func ParseTemplateProcessWrite2File(templateContent string, templateData interface{}, outputFile string, formatSource bool) error {

	if pkgTemplate, err := template.New("css").Parse(templateContent); err != nil {
		return err
	} else {
		if err := ProcessWrite2File(pkgTemplate, templateData, outputFile, formatSource); err != nil {
			return err
		}
	}

	return nil
}

func ParseTemplateWithFuncMapsProcessWrite2File(templateContent string, fMaps template.FuncMap, templateData interface{}, outputFile string, formatSource bool) error {

	if pkgTemplate, err := template.New("css").Funcs(fMaps).Parse(templateContent); err != nil {
		return err
	} else {
		if err := ProcessWrite2File(pkgTemplate, templateData, outputFile, formatSource); err != nil {
			return err
		}
	}

	return nil
}
*/

func Process(pkgTemplate *template.Template, templateData interface{}, formatSource bool) ([]byte, error) {

	builder := &bytes.Buffer{}

	if err := pkgTemplate.Execute(builder, templateData); err != nil {
		return nil, err
	}

	var data []byte
	if formatSource {
		var err error
		if data, err = format.Source(builder.Bytes()); err != nil {
			return nil, err
		}
	} else {
		data = builder.Bytes()
	}

	return data, nil
}

func MustProcess(pkgTemplate *template.Template, templateData interface{}, formatSource bool) []byte {

	builder := &bytes.Buffer{}

	if err := pkgTemplate.Execute(builder, templateData); err != nil {
		panic(err)
	}

	var data []byte
	if formatSource {
		var err error
		if data, err = format.Source(builder.Bytes()); err != nil {
			panic(err)
		}
	} else {
		data = builder.Bytes()
	}

	return data
}

func ProcessWrite2File(pkgTemplate *template.Template, templateData interface{}, outputFile string, formatSource bool) error {

	data, err := Process(pkgTemplate, templateData, formatSource)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputFile, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func LoadProcessWrite2File(templateFileName string, templateData interface{}, outputFile string, formatSource bool) error {

	if f, err := ioutil.ReadFile(templateFileName); err != nil {
		return err
	} else {
		var pkgTemplate *template.Template
		if pkgTemplate, err = template.New("css").Parse(string(f)); err != nil {
			return err
		} else {
			if err = ProcessWrite2File(pkgTemplate, templateData, outputFile, formatSource); err != nil {
				return err
			}
		}
	}

	return nil

}
