package restclient

import (
	"github.com/opentracing/opentracing-go"
	"time"
)

const (
	RequestTraceNameOpNamePlaceHolder    = "{op-name}"
	RequestTraceNameRequestIdPlaceHolder = "{req-id}"
	RequestIdTraceTag                    = "req-id"
	OpNameTraceTag                       = "op-name"
	LraHttpContextTraceTag               = "long-running-action"
)

type Header struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type Config struct {
	RestTimeout      time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	SkipVerify       bool          `mapstructure:"skip-verify" json:"skip-verify" yaml:"skip-verify"`
	Headers          []Header      `mapstructure:"headers" json:"headers" yaml:"headers"`
	TraceGroupName   string        `mapstructure:"trace-group-name" json:"trace-group-name" yaml:"trace-group-name"`
	TraceRequestName string        `mapstructure:"trace-req-name" json:"trace-req-name" yaml:"trace-req-name"`
	RetryCount       int           `mapstructure:"retry-count" json:"retry-count" yaml:"retry-count"`
	RetryWaitTime    time.Duration `mapstructure:"retry-wait-time" json:"retry-wait-time" yaml:"retry-wait-time"`
	RetryMaxWaitTime time.Duration `mapstructure:"retry-max-wait-time" json:"retry-max-wait-time" yaml:"retry-max-wait-time"`
	RetryOnHttpError []int         `mapstructure:"retry-on-errors" json:"retry-on-errors" yaml:"retry-on-errors"`

	Span opentracing.Span
}

type Option func(o *Config)

func WithSpan(span opentracing.Span) Option {
	return func(o *Config) {
		o.Span = span
	}
}

func WithTraceGroupName(opn string) Option {
	return func(o *Config) {
		o.TraceGroupName = opn
	}
}

func WithTraceRequestName(opn string) Option {
	return func(o *Config) {
		o.TraceRequestName = opn
	}
}

func WithSkipVerify(b bool) Option {
	return func(o *Config) {
		o.SkipVerify = b
	}
}

func WithTimeout(to time.Duration) Option {
	return func(o *Config) {
		o.RestTimeout = to
	}
}
