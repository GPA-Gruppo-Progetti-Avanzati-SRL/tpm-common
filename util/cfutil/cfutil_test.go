package cfutil_test

import (
	"math/rand"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/cfutil"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestCF(t *testing.T) {

	cf, err := cfutil.CalculateCF("imperato", "Mario Alessandro", "M", "19621121", false, "G337")
	require.NoError(t, err)
	t.Log(cf)

	// With omocodie.... Doesn't support check digit
	err = cfutil.CheckCFAgainstCFInfo("MPRMLS62S21G337J", "imperato", "Mario Alessandro", "M", "19621121", false, "G337")
	require.NoError(t, err)
}

func TestExtractInfo(t *testing.T) {
	fiscalCodeInfo, err := cfutil.ExtractInfo("MPRMLS62S21G337J")
	require.NoError(t, err)

	t.Log("fiscalCodeInfo", fiscalCodeInfo)
}

func TestRandomCF(t *testing.T) {
	nameSurnameCode := RandStringRunes([]rune("BCDFGLMNPQRSTVZ"), 6)
	year := rand.Intn(40)
	month := RandStringRunes(monthRune, 1)
	day := rand.Intn(28)

	log.Info().
		Str("nameSurnameCode", nameSurnameCode).
		Int("year", year).
		Str("month", month).
		Int("day", day).
		Msg("FiscalCode")
}

var monthRune = []rune("ABCDEHLMPRST")

func RandStringRunes(letterRunes []rune, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
