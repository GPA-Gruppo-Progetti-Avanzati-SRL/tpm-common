package util_test

import (
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog"
)

func TestEnvVar(t *testing.T) {
	util.LogEnvVars(zerolog.InfoLevel, "^GO")
}
