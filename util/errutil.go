package util

import "fmt"

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
