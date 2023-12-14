package util_test

import (
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
						"string-array-nested-value",
					},
				},
			},
			"propArrString": []string{
				"string-array-value",
			},
		},
	}

	p := util.GetProperty(m, "propString")
	require.NotNil(t, p)
	t.Log("propString", p)

	p = util.GetProperty(m, "propMap.propArr[].propArrItem")
	require.NotNil(t, p)
	t.Log("propMap.propArr[].propArrItem", p)

	p = util.GetProperty(m, "propMap.propArr[].propArrStringNested[]")
	require.NotNil(t, p)
	t.Log("propMap.propArr[].propArrStringNested[]", p)

	p = util.GetProperty(m, "propMap.propArrString[]")
	require.NotNil(t, p)
	t.Log("propMap.propArrString[]", p)
}
