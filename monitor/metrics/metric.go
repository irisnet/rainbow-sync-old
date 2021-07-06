package metrics

import (
	"github.com/weichang-bianjie/metric-sdk/metrics"
	"github.com/weichang-bianjie/metric-sdk/metrics/counter"
	"github.com/weichang-bianjie/metric-sdk/metrics/gauge"
)

type (
	Metric  metrics.Metric
	Guage   gauge.Client
	Counter counter.Client
)

func NewGuage(nameSpace string, subSystem string, name string, help string, labels []string) Metric {
	return gauge.NewGauge(
		nameSpace,
		subSystem,
		name,
		help,
		labels,
	)
}

func NewCounter(nameSpace string, subSystem string, name string, help string, labels []string) Metric {
	return counter.NewCounter(
		nameSpace,
		subSystem,
		name,
		help,
		labels,
	)
}

func CovertGuage(metric Metric) (Guage, bool) {
	value, ok := metric.(Guage)
	return value, ok
}

func CovertCounter(metric Metric) (Counter, bool) {
	value, ok := metric.(Counter)
	return value, ok
}
