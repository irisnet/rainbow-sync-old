package metrics

import (
	"fmt"
	"github.com/weichang-bianjie/metric-sdk"
	"github.com/weichang-bianjie/metric-sdk/metrics"
	"github.com/weichang-bianjie/metric-sdk/types"
)

type Monitor interface {
	RegisterMetrics(metric ...metrics.Metric)
	Report(reports ...func())
}

type client struct {
	metric_sdk.MetricClient
}

func NewMonitor(port int) Monitor {
	metricClient := metric_sdk.NewClient(types.Config{
		Address: fmt.Sprintf(":%v", port),
	})

	return client{metricClient}
}

func (c client) RegisterMetrics(metric ...metrics.Metric) {
	c.MetricClient.RegisterMetric(metric...)
}

func (c client) Report(reports ...func()) {
	c.MetricClient.Start(func() {
		for _, report := range reports {
			go report()
		}
	})
}
