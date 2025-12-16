package util

import (
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogEnvVars(level zerolog.Level, withPattern string) {
	const semLogContext = "util::log-env-vars"
	var err error

	var patternRegexp *regexp.Regexp
	if withPattern != "" {
		patternRegexp, err = regexp.Compile(withPattern)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext)
			return
		}
	}

	for _, v := range os.Environ() {
		ndx := strings.Index(v, "=")
		if ndx == -1 {
			log.Warn().Str("var", v).Msg(semLogContext + " - no equal sign")
			continue
		}

		vn := v[:ndx]
		vv := v[ndx:]
		if patternRegexp != nil && patternRegexp.MatchString(vn) || patternRegexp == nil {
			log.WithLevel(level).Str("var-name", vn).Str("var-value", vv).Msg(semLogContext)
		}
	}

}
