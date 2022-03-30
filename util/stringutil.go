package util

import (
	"fmt"
	"regexp"
	"strings"
	"text/scanner"
)

func Coalesce(args ...string) string {
	for _, str := range args {
		if str != "" {
			return str
		}
	}
	return ""
}

var NumericStringRegexp = regexp.MustCompile("^[-+]?\\d+(\\.\\d+)?$")

func IsNumeric(inputData string) bool {
	return NumericStringRegexp.Match([]byte(inputData))
}

func ParseSetClause(src string, clauseDelimiter rune) (map[string]interface{}, error) {

	var s scanner.Scanner
	var err error
	s.Init(strings.NewReader(src))
	s.Filename = "set-clause"

	st := 0
	propertyName := ""
	propertyValue := ""
	var resultMap map[string]interface{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		switch st {
		case 0:
			st, propertyName, err = parseSetClauseHandleState0(tok, &s)

		case 1:
			st, err = parseSetClauseHandleState1(tok, &s)

		case 2:
			st, propertyValue, err = parseSetClauseHandleState2(tok, &s)
			if err == nil {
				if resultMap == nil {
					resultMap = make(map[string]interface{}, 0)
				}
				resultMap[propertyName] = propertyValue
			}
		case 3:
			st, err = parseSetClauseHandleState3(tok, &s, clauseDelimiter)
		}

		if err != nil {
			return resultMap, err
		}

	}

	return resultMap, nil

}

func parseSetClauseHandleState0(tok rune, s *scanner.Scanner) (int, string, error) {
	if tok != scanner.Ident {
		return 0, "", fmt.Errorf("property name expected and found %s", s.TokenText())
	}

	return 1, s.TokenText(), nil
}

func parseSetClauseHandleState1(tok rune, s *scanner.Scanner) (int, error) {
	if tok != '=' {
		return 1, fmt.Errorf("equal sign expected and found %s", s.TokenText())
	}

	return 2, nil
}

func parseSetClauseHandleState2(tok rune, s *scanner.Scanner) (int, string, error) {

	if tok == scanner.String {
		pval := s.TokenText()
		return 3, pval[1 : len(pval)-1], nil
	} else {
		stb := strings.Builder{}
		for tok != scanner.EOF && tok != ',' {
			stb.WriteString(s.TokenText())
			tok = s.Scan()
		}
		return 0, stb.String(), nil
	}
}

func parseSetClauseHandleState3(tok rune, s *scanner.Scanner, clauseDelimiter rune) (int, error) {
	if tok != clauseDelimiter {
		return 3, fmt.Errorf("delimiter sign expected and found %s", s.TokenText())
	}

	return 0, nil
}

// Code from https://github.com/angular/angular-cli/blob/master/packages/angular_devkit/core/src/utils/strings.ts

const (
	STRING_DASHERIZE_REGEXP    = `[ _]`
	STRING_DECAMELIZE_REGEXP   = `([a-z\d])([A-Z])`
	STRING_CAMELIZE_REGEXP     = `(-|_|\.|\s)+(.)?`
	STRING_CAMELIZE_REGEXP_2   = `^([A-Z])`
	STRING_UNDERSCORE_REGEXP_1 = `([a-z\d])([A-Z]+)`
	STRING_UNDERSCORE_REGEXP_2 = `-|\s+`
)

func Decamelize(s string) string {
	m := regexp.MustCompile(STRING_DECAMELIZE_REGEXP)
	// fmt.Printf("%q\n", m.FindAllString(s, -1))
	return strings.ToLower(m.ReplaceAllString(s, "${1}_${2}"))
}

func Dasherize(s string) string {
	m := regexp.MustCompile(STRING_DASHERIZE_REGEXP)
	return m.ReplaceAllString(Decamelize(s), "-")
}

func Camelize(s string) string {
	m := regexp.MustCompile(STRING_CAMELIZE_REGEXP)
	s1 := m.ReplaceAllStringFunc(s, func(r string) string { return strings.ToUpper(r[len(r)-1:]) })

	m1 := regexp.MustCompile(STRING_CAMELIZE_REGEXP_2)
	return m1.ReplaceAllStringFunc(s1, func(r string) string { return strings.ToLower(r) })
}

func Classify(s string) string {
	sarr := strings.Split(s, ".")
	for i := 0; i < len(sarr); i++ {
		sarr[i] = Capitalize(Camelize(sarr[i]))
	}

	return strings.Join(sarr, ".")
}

func Underscore(s string) string {
	m := regexp.MustCompile(STRING_UNDERSCORE_REGEXP_1)
	s1 := m.ReplaceAllString(s, "${1}_${2}")

	m2 := regexp.MustCompile(STRING_UNDERSCORE_REGEXP_2)
	s2 := m2.ReplaceAllString(s1, "_")
	return strings.ToLower(s2)
}

func Capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}