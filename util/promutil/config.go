package promutil

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
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

const MetricsConfigReferenceLocalGroup = "local"

var DisabledMetricsConfigReference = MetricsConfigReference{
	GId:         "-",
	CounterId:   "-",
	HistogramId: "-",
	GaugeId:     "-",
}

type MetricsConfigReference struct {
	GId         string `yaml:"group-id,omitempty" mapstructure:"group-id,omitempty" json:"group-id,omitempty"`
	CounterId   string `yaml:"counter-id,omitempty" mapstructure:"counter-id,omitempty" json:"counter-id,omitempty"`
	HistogramId string `yaml:"histogram-id,omitempty" mapstructure:"histogram-id,omitempty" json:"histogram-id,omitempty"`
	GaugeId     string `yaml:"gauge-id,omitempty" mapstructure:"gauge-id,omitempty" json:"gauge-id,omitempty"`
}

func (mCfg MetricsConfigReference) IsLocal() bool {
	return mCfg.GId == MetricsConfigReferenceLocalGroup
}

func (mCfg MetricsConfigReference) IsZero() bool {
	return mCfg.GId == ""
}

func (mCfg MetricsConfigReference) IsEnabled() bool {
	return isMetricRefIdEnabled(mCfg.GId)
}

func (mCfg MetricsConfigReference) IsCounterEnabled() bool {
	return isMetricRefIdEnabled(mCfg.GId) && isMetricRefIdEnabled(mCfg.CounterId)
}

func (mCfg MetricsConfigReference) IsHistogramEnabled() bool {
	return isMetricRefIdEnabled(mCfg.GId) && isMetricRefIdEnabled(mCfg.HistogramId)
}

func (mCfg MetricsConfigReference) IsGaugeEnabled() bool {
	return isMetricRefIdEnabled(mCfg.GId) && isMetricRefIdEnabled(mCfg.GaugeId)
}

func isMetricRefIdEnabled(s string) bool {
	return s != "-" && s != ""
}

func (mCfg MetricsConfigReference) ResolveGroup(aGroup Group) (Group, bool, error) {

	var g Group
	var err error
	var ok bool
	if mCfg.IsEnabled() {
		if !mCfg.IsLocal() {
			g, err = GetGroup(mCfg.GId)
		} else {
			g = aGroup
		}

		if err == nil {
			if len(g) > 0 {
				ok = true
			}
		}
	}

	return g, ok, err
}

func CoalesceMetricsConfig(ref MetricsConfigReference, defaultVals MetricsConfigReference) MetricsConfigReference {
	gid := defaultVals
	if ref.GId != "" {
		gid.GId = ref.GId
	}

	if ref.CounterId != "" {
		gid.CounterId = ref.CounterId
	}

	if ref.HistogramId != "" {
		gid.HistogramId = ref.HistogramId
	}

	if ref.GaugeId != "" {
		gid.GaugeId = ref.GaugeId
	}

	return gid
}

type MetricLabelConfig struct {
	Id           string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Name         string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	DefaultValue string `yaml:"default-value,omitempty" mapstructure:"default-value,omitempty" json:"default-value,omitempty"`
}

type MetricLabelsConfig []MetricLabelConfig

func (li MetricLabelsConfig) Contains(n string) bool {
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
	Labels    MetricLabelsConfig
}

type MetricGroupConfig struct {
	GroupId    string         `yaml:"group-id,omitempty" mapstructure:"group-id,omitempty" json:"group-id,omitempty"`
	Namespace  string         `yaml:"namespace" mapstructure:"namespace" json:"namespace"`
	Subsystem  string         `yaml:"subsystem" mapstructure:"subsystem" json:"subsystem"`
	Collectors []MetricConfig `yaml:"collectors" mapstructure:"collectors" json:"collectors"`
}

func (c MetricGroupConfig) Merge(another MetricGroupConfig) MetricGroupConfig {
	const semLogContext = "metrics-group-config::merge"
	log.Info().Str("name", c.GroupId).Str("with", another.GroupId).Msg(semLogContext)

	c.Namespace = util.StringCoalesce(another.Namespace, c.Namespace)
	c.Subsystem = util.StringCoalesce(another.Subsystem, c.Subsystem)
	for _, def := range another.Collectors {
		targetCollectorNdx := -1
		for i, l := range c.Collectors {
			if l.Id == def.Id {
				targetCollectorNdx = i
				break
			}
		}

		if targetCollectorNdx == -1 {
			c.Collectors = append(c.Collectors, def)
		} else {
			c.Collectors[targetCollectorNdx] = def
		}
	}

	return c
}

type MetricConfig struct {
	Id      string                `yaml:"id" mapstructure:"id" json:"id"`
	Name    string                `yaml:"name" mapstructure:"name" json:"name"`
	Help    string                `yaml:"help" mapstructure:"help" json:"help"`
	Labels  []MetricLabelConfig   `yaml:"labels" mapstructure:"labels" json:"labels"`
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
