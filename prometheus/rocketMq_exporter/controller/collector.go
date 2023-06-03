package collector

import (
	"bufio"
	"io"
	"os"

	"github.com/ahwhy/myGolang/prometheus/rocketMq_exporter/conf"
	"github.com/prometheus/client_golang/prometheus"
)

func NewCollector() *Collector {
	return &Collector{
		Conf: conf.C(),
		count: prometheus.NewDesc(
			"rocketmq_count",
			"rocketmq_count",
			// 动态标签的key列表
			[]string{"group", "version", "type", "model"},
			// 静态标签
			prometheus.Labels{"module": "rocketmq"},
		),
		tps: prometheus.NewDesc(
			"rocketmq_tps",
			"rocketmq_tps",
			// 动态标签的key列表
			[]string{"group", "version", "type", "model"},
			// 静态标签
			prometheus.Labels{"module": "rocketmq"},
		),
		diff: prometheus.NewDesc(
			"rocketmq_diff_total",
			"rocketmq_diff_total",
			// 动态标签的key列表
			[]string{"group", "version", "type", "model"},
			// 静态标签
			prometheus.Labels{"module": "rocketmq"},
		),
	}
}

type Collector struct {
	Conf  *conf.Config
	count *prometheus.Desc
	tps   *prometheus.Desc
	diff  *prometheus.Desc
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.count
	ch <- c.tps
	ch <- c.diff
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	f, err := os.Open(c.Conf.FileConfig.Path)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		a, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		m := ParseLine(string(a))
		if m.Group != "#Group" {
			ch <- prometheus.MustNewConstMetric(c.count, prometheus.GaugeValue, m.IntCount(), m.Group, m.Version, m.Type, m.Model)
			ch <- prometheus.MustNewConstMetric(c.tps, prometheus.GaugeValue, m.IntTPS(), m.Group, m.Version, m.Type, m.Model)
			ch <- prometheus.MustNewConstMetric(c.diff, prometheus.GaugeValue, m.IntDiffTotal(), m.Group, m.Version, m.Type, m.Model)
		}
	}
}
