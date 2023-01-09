package restclient

import (
	"github.com/opentracing/opentracing-go"
	"time"
)

type Header struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type Config struct {
	RestTimeout      time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	SkipVerify       bool          `mapstructure:"skip-verify" json:"skip-verify" yaml:"skip-verify"`
	Headers          []Header      `mapstructure:"headers" json:"headers" yaml:"headers"`
	TraceOpName      string        `mapstructure:"trace-op-name" json:"trace-op-name" yaml:"trace-op-name"`
	RetryCount       int           `mapstructure:"retry-count" json:"retry-count" yaml:"retry-count"`
	RetryWaitTime    time.Duration `mapstructure:"retry-wait-time" json:"retry-wait-time" yaml:"retry-wait-time"`
	RetryMaxWaitTime time.Duration `mapstructure:"retry-max-wait-time" json:"retry-max-wait-time" yaml:"retry-max-wait-time"`
	RetryOnHttpError []int         `mapstructure:"retry-on-errors" json:"retry-on-errors" yaml:"retry-on-errors"`

	NestTraceSpans bool `mapstructure:"nest-trace-ops" json:"nest-trace-ops" yaml:"nest-trace-ops"`
	Span           opentracing.Span
}

type Option func(o *Config)

func WithNestedTraceSpan(b bool) Option {
	return func(o *Config) {
		if b || o.Span == nil {
			o.NestTraceSpans = b
		}
	}
}

func WithSpan(span opentracing.Span) Option {
	return func(o *Config) {
		o.Span = span
		if span != nil {
			o.NestTraceSpans = true
		}
	}
}

func WithTraceOperationName(opn string) Option {
	return func(o *Config) {
		o.TraceOpName = opn
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
