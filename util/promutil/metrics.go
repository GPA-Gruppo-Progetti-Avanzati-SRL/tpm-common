package promutil

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"strings"
)

type MetricRegistry []MetricInfo

func InitMetricsRegistry(metrics MetricsConfig) (MetricRegistry, error) {

	var metricsRegistry []MetricInfo
	for _, mCfg := range metrics.Collectors {
		if mc, err := NewCollector(metrics.Namespace, metrics.Subsystem, mCfg.Name, &mCfg); err != nil {
			log.Error().Err(err).Str("name", mCfg.Name).Msg("error creating metric")
			return nil, err
		} else {
			metricsRegistry = append(metricsRegistry, MetricInfo{Type: mCfg.Type, Id: mCfg.Id, Name: mCfg.Name, Collector: mc, Labels: mCfg.Labels})
		}
	}

	if len(metricsRegistry) == 0 {
		log.Warn().Msg("metrics registry is empty")
	}

	return metricsRegistry, nil
}

func (r MetricRegistry) FindCollectorByName(n string) MetricInfo {
	for _, c := range r {
		if c.Name == n {
			return c
		}
	}

	return MetricInfo{}
}

func (r MetricRegistry) FindCollectorById(id string) MetricInfo {
	for _, c := range r {
		if c.Id == id {
			return c
		}
	}

	return MetricInfo{}
}

func (r MetricRegistry) SetMetricValueByName(n string, v float64, labels prometheus.Labels) error {
	if c := r.FindCollectorByName(n); c.Type != "" {
		setMetricValue(c, v, labels)
	} else {
		err := errors.New("cannot find collector by name")
		log.Error().Err(err).Str("name", n).Send()
		return err
	}

	return nil
}

func (r MetricRegistry) SetMetricValueById(id string, v float64, labels prometheus.Labels) error {
	if c := r.FindCollectorById(id); c.Type != "" {
		setMetricValue(c, v, labels)
	} else {
		err := errors.New("cannot find collector by id")
		log.Error().Err(err).Str("id", id).Send()
		return err
	}

	return nil
}

func setMetricValue(c MetricInfo, v float64, labels prometheus.Labels) {

	labels = fixLabels(c.Labels, labels)

	switch c.Type {
	case MetricTypeCounter:
		cnter := c.Collector.(*prometheus.CounterVec)
		cnter.With(labels).Add(v)
	case MetricTypeGauge:
		gauger := c.Collector.(*prometheus.GaugeVec)
		gauger.With(labels).Set(v)
	case MetricTypeHistogram:
		hist := c.Collector.(*prometheus.HistogramVec)
		hist.With(labels).Observe(v)
	}
}

func fixLabels(cfgLabels LabelsInfo, providedLabels prometheus.Labels) prometheus.Labels {
	if len(cfgLabels) == 0 {
		return nil
	}

	actualLabels := make(prometheus.Labels)
	if len(providedLabels) == 0 {
		actualLabels = make(prometheus.Labels)
		for _, l := range cfgLabels {
			actualLabels[l.Name] = l.DefaultValue
		}

		return actualLabels
	}

	for _, l := range cfgLabels {
		if pl, ok := providedLabels[l.Name]; pl != "" && ok {
			actualLabels[l.Name] = pl
		} else {
			actualLabels[l.Name] = l.DefaultValue
		}
	}

	return actualLabels
}

func NewCollector(namespace string, subsystem string, opName string, metricConfig *MetricConfig) (prometheus.Collector, error) {

	var c prometheus.Collector
	switch metricConfig.Type {
	case MetricTypeCounter:
		c = NewCounter(namespace, subsystem, opName, metricConfig)
	case MetricTypeGauge:
		c = NewGauge(namespace, subsystem, opName, metricConfig)
	case MetricTypeHistogram:
		c = NewHistogram(namespace, subsystem, opName, metricConfig)
	default:
		return nil, errors.New("unknown metric type: " + metricConfig.Type)
	}

	if c == nil {
		return nil, errors.New("cannot instantiate metric: " + metricConfig.Name)
	}

	return c, nil
}

func NewCounter(namespace string, subsystem string, opName string, counterMetrics *MetricConfig) prometheus.Collector /* *prometheus.CounterVec */ {

	if counterMetrics.Type != MetricTypeCounter {
		log.Error().Str("type", counterMetrics.Type).Msg("type mismatch, not a counter")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("counter metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	var lbs []string
	if len(counterMetrics.Labels) != 0 {
		for _, l := range counterMetrics.Labels {
			lbs = append(lbs, l.Name)
		}
	}
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: metricSubsystem,
			Name:      counterMetrics.Name,
			Help:      counterMetrics.Help,
		},
		lbs)

	err := prometheus.Register(c)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", counterMetrics.Name).Msg("counter already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", counterMetrics.Name).Msg("counter error")
		}
	}

	return c
}

func NewGauge(namespace string, subsystem string, opName string, gaugeMetrics *MetricConfig) prometheus.Collector /* *prometheus.CounterVec */ {

	if gaugeMetrics.Type != MetricTypeGauge {
		log.Error().Str("type", gaugeMetrics.Type).Msg("type mismatch, not a gauge")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("gauge metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	var lbs []string
	if len(gaugeMetrics.Labels) != 0 {
		for _, l := range gaugeMetrics.Labels {
			lbs = append(lbs, l.Name)
		}
	}
	c := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: metricSubsystem,
			Name:      gaugeMetrics.Name,
			Help:      gaugeMetrics.Help,
		},
		lbs)

	err := prometheus.Register(c)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", gaugeMetrics.Name).Msg("counter already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", gaugeMetrics.Name).Msg("gauge error")
		}
	}

	return c
}

func NewHistogram(namespace string, subsystem string, opName string, histogramMetrics *MetricConfig) prometheus.Collector {

	if histogramMetrics.Type != MetricTypeHistogram {
		log.Error().Str("type", histogramMetrics.Type).Msg("type mismatch, not a histogram")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("histogram metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	var bck []float64
	switch t := histogramMetrics.Buckets.Type; t {
	case DefaultMetricsDurationBucketsTypeLinear:
		bck = prometheus.LinearBuckets(histogramMetrics.Buckets.Start, histogramMetrics.Buckets.WidthFactor, histogramMetrics.Buckets.Count)
	case DefaultMetricsDurationBucketsTypeExponential:
		bck = prometheus.ExponentialBuckets(histogramMetrics.Buckets.Start, histogramMetrics.Buckets.WidthFactor, histogramMetrics.Buckets.Count)
	case DefaultMetricsDurationBucketsTypeDefault:
		bck = prometheus.DefBuckets
	}

	var lbs []string
	if len(histogramMetrics.Labels) != 0 {
		for _, l := range histogramMetrics.Labels {
			lbs = append(lbs, l.Name)
		}
	}

	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: metricSubsystem,
		Name:      histogramMetrics.Name,
		Help:      histogramMetrics.Help,
		Buckets:   bck,
	}, lbs)

	err := prometheus.Register(h)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", histogramMetrics.Name).Msg("histogram already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", histogramMetrics.Name).Msg("histogram error")
		}
	}

	return h
}
