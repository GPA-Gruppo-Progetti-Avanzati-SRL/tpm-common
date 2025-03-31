package util_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrorRandomizer(t *testing.T) {

	var arr []struct {
		val       string
		shouldErr bool
	} = []struct {
		val       string
		shouldErr bool
	}{
		{val: "5/d", shouldErr: false},
		{val: "5/c", shouldErr: false},
		{val: "5/k", shouldErr: false},
		{val: "5/m", shouldErr: false},
	}

	for _, item := range arr {
		rnd, err := util.NewErrorRandomizer(item.val)
		if !item.shouldErr {
			require.NoError(t, err)

			for i := 0; i < 10000; i++ {
				if rnd.GenerateRandomError() != nil {
					t.Log("Error from GenerateRandomError at i:", i, " for: ", item.val)
					break
				}
			}
		} else {
			require.Error(t, err)
		}
	}

}
