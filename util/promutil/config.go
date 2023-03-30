package promutil

import (
	"github.com/prometheus/client_golang/prometheus"
)

const DefaultMetricsDurationBucketsTypeLinear = "linear"
const DefaultMetricsDurationBucketsTypeExponential = "exponential"
const DefaultMetricsDurationBucketsTypeDefault = "default"

const DefaultMetricsDurationBucketsStart = 0.5
const DefaultMetricsDurationBucketsWidthFormat = 0.5
const DefaultMetricsDurationBucketsCount = 10

const MetricTypeCounter = "counter"
const MetricTypeGauge = "gauge"
const MetricTypeHistogram = "histogram"

//type MetricsCounterConfig struct {
//	Name   string
//	Help   string
//	Labels string
//}
//
//type MetricsGaugeConfig struct {
//	Name   string
//	Help   string
//	Labels string
//}

type LabelConfig struct {
	Id           string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Name         string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	DefaultValue string `yaml:"default-value,omitempty" mapstructure:"default-value,omitempty" json:"default-value,omitempty"`
}

type LabelsConfig []LabelConfig

func (li LabelsConfig) Contains(n string) bool {
	for _, ln := range li {
		if ln.Name == n {
			return true
		}
	}

	return false
}

type Metric struct {
	Id        string
	Type      string
	Name      string
	Collector prometheus.Collector
	Labels    LabelsConfig
}

type GroupConfig struct {
	GroupId    string   `yaml:"group-id,omitempty" mapstructure:"group-id,omitempty" json:"group-id,omitempty"`
	Namespace  string   `yaml:"namespace" mapstructure:"namespace" json:"namespace"`
	Subsystem  string   `yaml:"subsystem" mapstructure:"subsystem" json:"subsystem"`
	Collectors []Config `yaml:"collectors" mapstructure:"collectors" json:"collectors"`
}

type Config struct {
	Id      string                `yaml:"id" mapstructure:"id" json:"id"`
	Name    string                `yaml:"name" mapstructure:"name" json:"name"`
	Help    string                `yaml:"help" mapstructure:"help" json:"help"`
	Labels  []LabelConfig         `yaml:"labels" mapstructure:"labels" json:"labels"`
	Type    string                `yaml:"type" mapstructure:"type" json:"type"`
	Buckets HistogramBucketConfig `yaml:"buckets" mapstructure:"buckets" json:"buckets"`
}

/*type MetricsHistogramConfig struct {
	Name    string
	Help    string
	Labels  string
	Buckets HistogramBucketConfig
}
*/

type HistogramBucketConfig struct {
	Type        string  `yaml:"type" mapstructure:"type" json:"type"`
	Start       float64 `yaml:"start" mapstructure:"start" json:"start"`
	WidthFactor float64 `yaml:"width-factor" mapstructure:"width-factor" json:"width-factor"`
	Count       int     `yaml:"count" mapstructure:"count" json:"count"`
}
