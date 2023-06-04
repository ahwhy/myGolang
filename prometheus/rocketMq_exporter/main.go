package main

import (
	"fmt"
	"os"

	"github.com/ahwhy/myGolang/prometheus/rocketMq_exporter/collector"
	"github.com/ahwhy/myGolang/prometheus/rocketMq_exporter/conf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func main() {
	// 读取配置
	if err := conf.LoadConfigFromEnv(); err != nil {
		fmt.Println("load config error, ", err)
		os.Exit(1)
	}

	// 设置参数
	if len(os.Args) > 1 {
		conf.C().CmdConfig.Target = os.Args[1]
	}

	// 注册采集器
	c := collector.NewCollector()
	retgistry := prometheus.NewRegistry()
	err := retgistry.Register(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 采集数据
	mf, err := retgistry.Gather()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 编码输出
	enc := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)
	for i := range mf {
		enc.Encode(mf[i])
	}
}
