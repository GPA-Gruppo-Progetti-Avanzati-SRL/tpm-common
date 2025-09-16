package util_test

import (
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
)

func TestVersionUtil(t *testing.T) {

	iws := []InputWanted{
		{
			input: "6.5.4-alpha.1",
		},
		{
			input: "6.5.4-beta.1",
		},
		{
			input: "6.5.4-beta.21",
		},
		{
			input: "v1.2.3-RC5",
		},
	}

	t.Log("Testing New")
	for _, iw := range iws {
		v, err := util.NewVersionNumberFromString(iw.input)
		require.NoError(t, err)

		t.Log(iw.input, v.String())
	}

	iws = []InputWanted{
		{
			input: "6.5.4-alpha.1",
		},
		{
			input: "6.5.4-beta.1",
		},
		{
			input: "6.5.4-beta.12",
		},
		{
			input: "6.5.4-a",
		},
	}

	t.Log("Testing LessThan")
	var maxValue util.VersionNumber
	for i, iw := range iws {
		v, err := util.NewVersionNumberFromString(iw.input)
		require.NoError(t, err)

		t.Log("version-number: ", v.String())

		if i == 0 {
			maxValue = v
		} else {
			if maxValue.LessThan(v) {
				maxValue = v
			}
		}
	}

	t.Log("max-value: ", maxValue)

}
