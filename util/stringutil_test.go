package util_test

import (
	"fmt"
	"github.com/mario-imperato/tpm-common/util"
	"github.com/stretchr/testify/assert"
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
}

func TestStrings(t *testing.T) {

	assert := assert.New(t)

	var s []InputWanted
	var modS string

	// Decamelize
	s = []InputWanted{
		{"innerHTML", "inner_html"},
		{"action_name", "action_name"},
		{"css-class-name", "css-class-name"},
		{"my favorite items", "my favorite items"},
		{"CONTO BANCOPOSTA RETAIL", "conto bancoposta retail"},
	}

	for _, iw := range s {
		modS = util.Decamelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		assert.Equal(iw.wanted, modS, "decamelize: strings should match")
	}

	// Dasherize
	s = []InputWanted{
		{"innerHTML", "inner-html"},
		{"action_name", "action-name"},
		{"css-class-name", "css-class-name"},
		{"my favorite items", "my-favorite-items"},
		{"CONTO BANCOPOSTA RETAIL", "conto-bancoposta-retail"},
	}

	for _, iw := range s {
		modS = util.Dasherize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		assert.Equal(iw.wanted, modS, "dasherize: strings should match")
	}

	// Camelize
	s = []InputWanted{
		{"innerHTML", "innerHTML"},
		{"action_name", "actionName"},
		{"css-class-name", "cssClassName"},
		{"my favorite items", "myFavoriteItems"},
		{"My Favorite Items", "myFavoriteItems"},
		{"CONTO BANCOPOSTA RETAIL", "cONTOBANCOPOSTARETAIL"},
	}

	for _, iw := range s {
		modS = util.Camelize(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		assert.Equal(iw.wanted, modS, "camelize: strings should match")
	}

	// Classify
	s = []InputWanted{
		{"innerHTML", "InnerHTML"},
		{"action_name", "ActionName"},
		{"css-class-name", "CssClassName"},
		{"my favorite items", "MyFavoriteItems"},
		{"CONTO BANCOPOSTA RETAIL", "CONTOBANCOPOSTARETAIL"},
	}

	for _, iw := range s {
		modS = util.Classify(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		assert.Equal(iw.wanted, modS, "classify: strings should match")
	}

	// Underscore
	s = []InputWanted{
		{"innerHTML", "inner_html"},
		{"action_name", "action_name"},
		{"css-class-name", "css_class_name"},
		{"my favorite items", "my_favorite_items"},
		{"CONTO BANCOPOSTA RETAIL", "conto_bancoposta_retail"},
	}

	for _, iw := range s {
		modS = util.Underscore(iw.input)
		fmt.Printf("%s --> %s\n", iw.input, modS)
		assert.Equal(iw.wanted, modS, "underscore: strings should match")
	}
}
