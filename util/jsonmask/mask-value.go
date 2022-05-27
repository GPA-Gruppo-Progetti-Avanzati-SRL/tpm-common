package jsonmask

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func doMask(target interface{}, path string, val interface{}, fi FieldInfo) error {
	log.Info().Str("path", path).Interface("val", val).Msgf("%T data in parent type %T", val, target)

	ndx := strings.LastIndex(path, ".")
	if ndx < 0 {
		return fmt.Errorf("unsupported path %s", path)
	}

	propName := path[ndx+1:]
	propName = strings.TrimPrefix(propName, "[")
	propName = strings.TrimSuffix(propName, "]")

	var maskedVal interface{}
	var err error
	switch tVal := val.(type) {
	case string:
		maskedVal = randomMask(tVal, '*')
	case float64:
		sMaskedVal := randomMask(fmt.Sprintf("%f", tVal), '0')
		maskedVal, err = strconv.ParseFloat(sMaskedVal, 64)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type %T", val)
	}

	switch cnt := target.(type) {
	case []interface{}:
		propIndex, err := strconv.Atoi(propName)
		if err != nil {
			return err
		}

		cnt[propIndex] = maskedVal
	case map[string]interface{}:
		cnt[propName] = maskedVal
	default:
		return fmt.Errorf("unsupported type %T in doMask operation", val)
	}

	return nil
}

/*
 * Source: https://github.com/rkritchat/jsonmask/blob/master/jsonmask.go
 * Opened an issue because it fails to handle number sensitive data. replied nonsense...
 * https://github.com/rkritchat/jsonmask/issues/2
 */

func randomMask(c string, maskingChar rune) string {
	if len(c) == 0 {
		return c
	}
	var r = []rune(c)
	var cl = len(r)
	var size = initMaskSize(cl)
	var count int
	raffle := make(map[int]int, size)
	for i := 0; i < cl; i++ {
		count += 1 //avoid random forever
		if len(raffle) == size || count == 10 {
			//break if mask enough
			break
		}
		v := randPos(cl)
		if _, ok := raffle[v]; ok {
			i -= 1
			continue
		}
		//case not mask yet
		if len(r)-1 >= v {
			r[v] = maskingChar
			raffle[v] = v
		}
	}
	return string(r)
}

func randPos(max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	ra := rand.New(source)
	return ra.Intn(max)
}

func initMaskSize(l int) int {
	if l == 1 {
		return l
	}
	return l / 2
}
