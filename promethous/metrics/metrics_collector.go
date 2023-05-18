package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewDemoCollector() *DemoCollector {
	return &DemoCollector{
		queueLengthDesc: prometheus.NewDesc(
			"myGolang_prome_metrics_queue_length",
			"The number of items in the queue.",
			// 动态标签的key列表
			[]string{"instnace_id", "instnace_name"},
			// 静态标签
			prometheus.Labels{"module": "http-server"},
		),
		// 动态标的value列表, 必须与声明的动态标签的key一一对应
		labelValues: []string{"mq_001", "kafka01"},
	}
}

type DemoCollector struct {
	queueLengthDesc *prometheus.Desc
	labelValues     []string
}

func (c *DemoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.queueLengthDesc
}

func (c *DemoCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.queueLengthDesc, prometheus.GaugeValue, 100, c.labelValues...)
}