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

type LabelInfo struct {
	Name         string `yaml:"name" mapstructure:"name" json:"name"`
	DefaultValue string `yaml:"default-value" mapstructure:"default-value" json:"default-value"`
}

type LabelsInfo []LabelInfo

func (li LabelsInfo) Contains(n string) bool {
	for _, ln := range li {
		if ln.Name == n {
			return true
		}
	}

	return false
}

type MetricInfo struct {
	Id        string
	Type      string
	Name      string
	Collector prometheus.Collector
	Labels    LabelsInfo
}

type MetricsConfig struct {
	Namespace  string         `yaml:"namespace" mapstructure:"namespace" json:"namespace"`
	Subsystem  string         `yaml:"subsystem" mapstructure:"subsystem" json:"subsystem"`
	Collectors []MetricConfig `yaml:"collectors" mapstructure:"collectors" json:"collectors"`
}

type MetricConfig struct {
	Id      string                `yaml:"id" mapstructure:"id" json:"id"`
	Name    string                `yaml:"name" mapstructure:"name" json:"name"`
	Help    string                `yaml:"help" mapstructure:"help" json:"help"`
	Labels  []LabelInfo           `yaml:"labels" mapstructure:"labels" json:"labels"`
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
