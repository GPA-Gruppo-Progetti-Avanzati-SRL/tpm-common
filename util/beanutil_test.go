package util_test

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetProperty(t *testing.T) {
	m := map[string]interface{}{
		"propString": "propString-value",
		"propMap": map[string]interface{}{
			"propInt": 27,
			"propArr": []interface{}{
				map[string]interface{}{
					"propArrItem": "propArrItemValue",
					"propArrStringNested": []string{
						"string-array-nested-value-1",
						"string-array-nested-value-2",
					},
				},
			},
			"propArrString": []string{
				"string-array-value-1",
				"string-array-value-2",
			},
		},
	}

	testCases := []InputWanted{
		{input: "propString", wanted: "propString-value"},
		{input: "propMap.propArr[].propArrItem", wanted: "propArrItemValue"},
		{input: "propMap.propArr[].propArrStringNested[]", wanted: "string-array-nested-value-1"},
		{input: "propMap.propArrString[]", wanted: "string-array-value-1"},
		{input: "propMap.propArrString", wanted: "string-array-value-1"},
		{input: "propMap.propArrString[*]", wanted: "string-array-value-1,string-array-value-2"},
	}

	for i, c := range testCases {
		p := util.GetProperty(m, c.input)
		require.Equal(t, c.wanted, p, fmt.Sprintf("error on %d case", i))
	}

}
