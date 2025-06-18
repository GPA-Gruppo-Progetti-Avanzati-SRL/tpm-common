package promutil_test

import (
	"embed"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

var testMetric = []byte(`
namespace: mdb_symphony
subsystem: activity
collectors:
  - help: numero richieste
    id: activity-counter
    labels:
      - default-value: N/A
        id: type
        name: type
      - default-value: N/A
        id: name
        name: name
      - default-value: N/A
        id: endpoint
        name: endpoint
      - default-value: N/A
        id: status-code
        name: status_code
    name: counter
    type: counter
  - buckets:
      count: 10
      start: 0.5
      type: linear
      width-factor: 0.5
    help: durata lavorazione richiesta
    id: activity-duration-new
    labels:
      - default-value: N/A
        id: type
        name: type
      - default-value: N/A
        id: name
        name: name
      - default-value: N/A
        id: endpoint
        name: endpoint
      - default-value: N/A
        id: status-code
        name: status_code
    name: duration
    type: histogram
`)

//go:embed test-metrics/*
var embeddedMetricsDefs embed.FS

func TestRegistry(t *testing.T) {

	embedded, err := promutil.ReadEmbeddedMetricGroupConfig("test-metrics", embeddedMetricsDefs)
	require.NoError(t, err)

	var gc promutil.MetricGroupConfig
	err = yaml.Unmarshal(testMetric, &gc)
	require.NoError(t, err)

	_, err = promutil.InitRegistry(embedded, map[string]promutil.MetricGroupConfig{"activity": gc})
	require.NoError(t, err)
}
