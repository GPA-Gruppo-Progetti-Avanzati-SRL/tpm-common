package funcs

import (
	"github.com/rs/zerolog/log"
	"time"
)

func Now(fmt string) (string, error) {
	return time.Now().Format(fmt), nil
}

func NowAfter(d string, fmt string) (string, error) {
	const semLogContext = "funcs::now-after"

	dur, err := time.ParseDuration(d)
	if err != nil {
		log.Error().Err(err).Msg(semLogContext)
		return "", err
	}
	return time.Now().Add(dur).Format(fmt), nil
}
