package collector_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/prometheus/rocketMq_exporter/collector"
	"github.com/ahwhy/myGolang/prometheus/rocketMq_exporter/conf"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	retgistry = prometheus.NewRegistry()
)

func TestCollect(t *testing.T) {
	c := collector.NewCollector()
	c.Conf.FileConfig.Path = "../data/sample.txt"

	if err := retgistry.Register(c); err != nil {
		t.Fatal(err)
	}

	mf, err := retgistry.Gather()
	if err != nil {
		t.Fatal(err)
	}

	// 编码输出
	b := bytes.NewBuffer([]byte{})
	enc := expfmt.NewEncoder(b, expfmt.FmtText)
	for i := range mf {
		enc.Encode(mf[i])
	}

	fmt.Println(b.String())
}

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}
}
