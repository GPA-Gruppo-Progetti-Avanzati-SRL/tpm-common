package util_test

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsNumeric(t *testing.T) {

	sarr := []string{
		"T01009140904",
		"1+1",
		"+23.90",
	}

	for _, s := range sarr {
		t.Logf("Is numeric %s --> %t", s, util.IsNumeric(s))
	}
}

type InputWanted struct {
	input  string
	wanted string
	param  int
}

func TestStringJoin(t *testing.T) {

	s := []string{
		"0123456789", "ABCDEFGHIJ",
	}

	ns := util.StringJoin(s, "-", 15)
	require.Equal(t, 15, len(ns), "failed to join")
	t.Log(ns)

	ns = util.StringJoin(s, "-", -15)
	require.Equal(t, 15, len(ns), "failed to join")
	t.Log(ns)

	ns = util.StringJoin(s, "-", 0)
	require.Equal(t, 21, len(ns), "failed to join")
	t.Log(ns)
}

func TestMaxLengh(t *testing.T) {

	assert := assert.New(t)
	var s []InputWanted

	s = []InputWanted{
		{input: "0123456789", wanted: "01234", param: 5},
		{input: "0123456789", wanted: "0123456789", param: 10},
		{input: "0123456789", wanted: "0123456789", param: -10},
		{input: "0123456789", wanted: "789", param: -3},
		{input: "0123456789", wanted: "0123456789", param: 0},
	}

	for _, iw := range s {
		v, _ := util.ToMaxLength(iw.input, iw.param)
		fmt.Printf("%s (%d) --> %s\n", iw.input, iw.param, v)
		assert.Equal(iw.wanted, v, "to max length: strings should match")
	}
}

type InputWanted4PrefixWithWildCard struct {
	input        string
	prefix       string
	wildCardChar byte
	shouldMatch  bool
}

func TestHasPrefixWithWildCard(t *testing.T) {

	a := assert.New(t)
	var s []InputWanted4PrefixWithWildCard

	s = []InputWanted4PrefixWithWildCard{
		{input: "0123456789", prefix: "01234", wildCardChar: '*', shouldMatch: true},
		{input: "0123D56789", prefix: "01234", wildCardChar: '*', shouldMatch: false},
		{input: "0123456789", prefix: "01*34", wildCardChar: '*', shouldMatch: true},
	}

	for _, iw := range s {
		b := util.HasPrefixWithWildCard(iw.input, iw.prefix, iw.wildCardChar)
		a.Equal(iw.shouldMatch, b)
	}
}

func TestPadLengh(t *testing.T) {

	a := assert.New(t)
	var s []InputWanted

	s = []InputWanted{
		{input: "0123456789", wanted: "0123456789", param: 10},
		{input: "0123456789", wanted: "0123456789", param: 7},
		{input: "0123456789", wanted: "-----0123456789", param: -15},
		{input: "0123456789", wanted: "0123456789-----", param: 15},
		{input: "0123456789", wanted: "0123456789", param: 0},
	}

	for _, iw := range s {
		v, _ := util.Pad2Length(iw.input, iw.param, "-")
		fmt.Printf("%s (%d) --> %s\n", iw.input, iw.param, v)
		a.Equal(iw.wanted, v, "pad to length: strings should match")
	}
}

func TestStrings(t *testing.T) {

	a := assert.New(t)

	var s []InputWanted
	var modS string

	// Decamelize
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner_html"},
		{input: "action_name", wanted: "action_name"},
		{input: "css-class-name", wanted: "css-class-name"},
		{input: "my favorite items", wanted: "my favorite items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto bancoposta retail"},
		{input: "camt_029_001_09", wanted: "camt_029_001_09"},
	}

	for _, iw := range s {
		modS = util.Decamelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "decamelize: strings should match")
	}

	// Dasherize
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner-html"},
		{input: "action_name", wanted: "action-name"},
		{input: "css-class-name", wanted: "css-class-name"},
		{input: "my favorite items", wanted: "my-favorite-items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto-bancoposta-retail"},
		{input: "camt_029_001_09", wanted: "camt-029-001-09"},
	}

	for _, iw := range s {
		modS = util.Dasherize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "dasherize: strings should match")
	}

	// Camelize
	s = []InputWanted{
		{input: "innerHTML", wanted: "innerHTML"},
		{input: "action_name", wanted: "actionName"},
		{input: "css-class-name", wanted: "cssClassName"},
		{input: "my favorite items", wanted: "myFavoriteItems"},
		{input: "My Favorite Items", wanted: "myFavoriteItems"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "cONTOBANCOPOSTARETAIL"},
		{input: "camt_029_001_09", wanted: "camt02900109"},
		{input: "camt.029.001.09", wanted: "camt02900109"},
	}

	for _, iw := range s {
		modS = util.Camelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "camelize: strings should match")
	}

	// Classify
	s = []InputWanted{
		{input: "innerHTML", wanted: "InnerHTML"},
		{input: "action_name", wanted: "ActionName"},
		{input: "css-class-name", wanted: "CssClassName"},
		{input: "my favorite items", wanted: "MyFavoriteItems"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "CONTOBANCOPOSTARETAIL"},
		{input: "camt_029_001_09", wanted: "Camt02900109"},
		{input: "camt.029.001.09", wanted: "Camt.029.001.09"},
	}

	for _, iw := range s {
		modS = util.Classify(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "classify: strings should match")
	}

	// Underscore
	s = []InputWanted{
		{input: "innerHTML", wanted: "inner_html"},
		{input: "action_name", wanted: "action_name"},
		{input: "css-class-name", wanted: "css_class_name"},
		{input: "my favorite items", wanted: "my_favorite_items"},
		{input: "CONTO BANCOPOSTA RETAIL", wanted: "conto_bancoposta_retail"},
		{input: "camt_029_001_09", wanted: "camt_029_001_09"},
	}

	for _, iw := range s {
		modS = util.Underscore(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		a.Equal(iw.wanted, modS, "underscore: strings should match")
	}
}
