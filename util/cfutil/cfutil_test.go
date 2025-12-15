package cfutil_test

import (
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/cfutil"
	"github.com/stretchr/testify/require"
)

func TestCF(t *testing.T) {

	cf := cfutil.CalculateCF("imperato", "Mario Alessandro", "M", "19621121", false, "G337")
	t.Log(cf)

	// With omocodie.... Doesn't support check digit
	err := cfutil.CheckCF("MPRMLSSNSNMGPPTJ", "imperato", "Mario Alessandro", "M", "19621121", false, "G337")
	require.NoError(t, err)
}
