package funcs

import (
	"github.com/rs/zerolog/log"
	"time"
)

func NowAfterDuration(d string, fmt string) (string, error) {
	const semLogContext = "funcs::now-after-duration"

	dur, err := time.ParseDuration(d)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}
	return time.Now().Add(dur).Format(fmt), nil
}
