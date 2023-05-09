package workload_test

import (
	"context"

	"github.com/ahwhy/myGolang/k8s"
	"github.com/ahwhy/myGolang/k8s/workload"
	"github.com/ahwhy/myGolang/utils/conf"
)

var (
	impl *workload.Client
	ctx  = context.Background()
)

func init() {
	client, err := k8s.NewClientFromFile("./kubeConfig.yaml")
	if err != nil {
		panic(err)
	}

	// 加载单元测试的变量
	conf.LoadConfigFromEnv()
	impl = client.WorkLoad()
}
