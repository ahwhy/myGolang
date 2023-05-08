package admin_test

import (
	"context"

	"github.com/ahwhy/myGolang/k8s"
	"github.com/ahwhy/myGolang/k8s/admin"
)

var (
	impl *admin.Client
	ctx  = context.Background()
)

func init() {
	client, err := k8s.NewClientFromFile("./kubeConfig.yaml")
	if err != nil {
		panic(err)
	}

	impl = client.Admin()
}
