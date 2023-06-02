package collector_test

import (
	"github.com/ahwhy/myGolang/promethous/rocketMq_exporter/conf"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	retgistry = prometheus.NewRegistry()
)

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}
}
