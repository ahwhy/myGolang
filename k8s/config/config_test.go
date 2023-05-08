package config_test

import (
	"context"

	"github.com/ahwhy/myGolang/k8s"
	"github.com/ahwhy/myGolang/k8s/config"
)

var (
	impl *config.Client
	ctx  = context.Background()
)

func init() {
	client, err := k8s.NewClientFromFile("../kube_config.yml")
	if err != nil {
		panic(err)
	}
	impl = client.Config()
}
