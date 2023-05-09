package k8s_test

import (
	"testing"

	"github.com/ahwhy/myGolang/k8s"
	"github.com/ahwhy/myGolang/utils/tools"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	client *k8s.Client
)

func TestServerVersion(t *testing.T) {
	v, err := client.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestServerResources(t *testing.T) {
	rs, err := client.ServerResources()
	if err != nil {
		t.Log(err)
	}
	for i := range rs {
		t.Log(rs[i].GroupVersion, rs[i].APIVersion)
		for _, r := range rs[i].APIResources {
			t.Log(r)
		}
	}
}

func init() {
	zap.DevelopmentSetup()

	kubeConf := tools.MustReadContentFile("kube_config.yml")
	c, err := k8s.NewClient(kubeConf)
	if err != nil {
		panic(err)
	}
	client = c
}
