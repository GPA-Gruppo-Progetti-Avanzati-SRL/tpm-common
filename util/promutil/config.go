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

type MetricInfo struct {
	Id        string
	Type      string
	Name      string
	Collector prometheus.Collector
	Labels    string
}

type MetricsConfig struct {
	Namespace  string
	Subsystem  string
	Collectors []MetricConfig
}

type MetricConfig struct {
	Id      string
	Name    string
	Help    string
	Labels  string
	Type    string
	Buckets HistogramBucketConfig
}

/*type MetricsHistogramConfig struct {
	Name    string
	Help    string
	Labels  string
	Buckets HistogramBucketConfig
}
*/

type HistogramBucketConfig struct {
	Type        string
	Start       float64
	WidthFactor float64
	Count       int
}
