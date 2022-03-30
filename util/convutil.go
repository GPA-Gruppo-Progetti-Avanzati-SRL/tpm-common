package util

import (
	"github.com/rs/zerolog/log"
	"strconv"
)

func ConvToInt64(v interface{}) int64 {

	switch val := v.(type) {
	case int32:
		return int64(val)
	case int64:
		return val
	case string:
		ival, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("conversion error")
		}

		return int64(ival)
	}

	return 0
}

func ConvToInt32(v interface{}) int32 {

	switch val := v.(type) {
	case int32:
		return val
	case int64:
		return int32(val)
	case string:
		ival, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.Error().Err(err).Msg("conversion error")
		}

		return int32(ival)
	}

	return 0
}
