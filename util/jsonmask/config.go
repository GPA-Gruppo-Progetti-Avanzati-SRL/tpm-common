package jsonmask

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type FieldInfo struct {
	Path     string `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path,omitempty"`
	MaskType string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
	uxPath   string
	indexes  []int
}

type Domain struct {
	Name   string      `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Fields []FieldInfo `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`
}

type Config map[string]Domain

type builder struct {
	fn       string
	yamlData []byte
}

type Option func(cfg *builder)

func FromFileName(fn string) Option {
	return func(cfg *builder) {
		cfg.fn = fn
	}
}

func FromData(b []byte) Option {
	return func(cfg *builder) {
		cfg.yamlData = b
	}
}

func readConfigFromFile(fn string) (Config, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	return readConfig(b)
}

func readConfig(b []byte) (Config, error) {
	cfg := make(Config)
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
