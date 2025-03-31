package util

import "fmt"

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

func NewErrorRandomizer(s int, p int) ErrorRandomizerFunc {
	return func() error {
		const semLogContext = "kafka-sink-stage-queue::get-random-error"

		if rand.IntN(s) < p {
			log.Warn().Int("with-random-error", p).Int("scale", s).Msg(semLogContext)
			return errors.New("sink-stage queue random error")
		}

		return nil
	}
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
