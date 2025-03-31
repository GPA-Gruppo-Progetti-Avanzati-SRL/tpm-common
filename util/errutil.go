package util

import (
	"fmt"
	"strconv"
	"strings"
)

import (
	"errors"
	"github.com/rs/zerolog/log"
	"math/rand/v2"
)

type ErrorRandomizer interface {
	GenerateRandomError() error
}

type ErrorRandomizerFunc func() error

func (f ErrorRandomizerFunc) GenerateRandomError() error {
	return f()
}

func NewErrorRandomizer(p string) (ErrorRandomizerFunc, error) {

	if len(p) == 0 || p == "0" {
		return nil, nil
	}

	if len(p) < 3 {
		return nil, errors.New("invalid randomizer format: form is d+/(c|d|k|m)")
	}

	unit := strings.ToLower(p[len(p)-2:])
	value, err := strconv.Atoi(p[0 : len(p)-2])
	if err != nil {
		return nil, err
	}

	if value <= 0 {
		return nil, nil
	}

	var scale int
	switch unit {
	case "/d":
		scale = 10
	case "/c":
		scale = 100
	case "/k":
		scale = 1000
	case "/m":
		scale = 10000
	default:
		return nil, errors.New("invalid randomizer unit: " + unit + ", supported units: c, d, k, m")
	}

	if value >= scale {
		return nil, errors.New("invalid randomizer value " + p)
	}

	return func() error {
		const semLogContext = "kafka-sink-stage-queue::get-random-error"

		if rand.IntN(scale) < value {
			log.Warn().Str("with-random-error", p).Msg(semLogContext)
			return errors.New("sink-stage queue random error")
		}

		return nil
	}, nil
}

type ErrorWithCode struct {
	Rc  string
	Err error
}

func NewError(c string, err error) *ErrorWithCode {
	if c == "" {
		c = "500"
	}
	return &ErrorWithCode{Rc: c, Err: err}
}

func (r *ErrorWithCode) Code() string {
	return r.Rc
}

func (r *ErrorWithCode) Error() string {
	return fmt.Sprintf("[%s] %s", r.Rc, r.Err.Error())
}

func (r *ErrorWithCode) Unwrap() error {
	return r.Err
}
