package k8s_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ahwhy/myGolang/k8s"
)

var (
	kubeConfig = ""
	system     = "kube-system"
)

func TestGetter(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	v, err := client.ServerVersion()
	should.NoError(err)
	fmt.Printf("%v\n", v)
	fmt.Printf("%v\n", client.CurrentContext())
	fmt.Printf("%v\n", client.CurrentCluster())

	ctx := context.Background()
	node, err := client.ListNode(ctx)
	should.NoError(err)
	ns, err := client.ListNamespace(ctx)
	should.NoError(err)
	fmt.Printf("%v\n", node)
	fmt.Printf("%v\n", ns)
}

func TestListConfigMap(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	ctx := context.Background()
	oper := k8s.NewListConfigMapRequest(system)
	v, err := client.ListConfigMap(ctx, oper)
	should.NoError(err)
	fmt.Printf("%v\n", v.Items)
	fmt.Printf("%v\n", v.ListMeta)
}

func TestListCronJob(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	ctx := context.Background()
	oper := k8s.NewListCronJobRequest(system)
	v, err := client.ListCronJob(ctx, oper)
	should.NoError(err)
	fmt.Printf("%v\n", v)
	fmt.Printf("%v\n", v.ListMeta)
}

func TestListDaemonSet(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	ctx := context.Background()
	oper := k8s.NewListDaemonSetRequest(system)
	v, err := client.ListDaemonSet(ctx, oper)
	should.NoError(err)
	fmt.Printf("%v\n", v.Items)
	// fmt.Printf("%v\n", v.ListMeta)
}

func TestNewListDeployment(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	ctx := context.Background()
	oper := k8s.NewListDeploymentRequest(system)
	v, err := client.ListDeployment(ctx, oper)
	should.NoError(err)
	fmt.Printf("%v\n", v.Items)
	// fmt.Printf("%v\n", v.ListMeta)
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	kc, err := ioutil.ReadFile(filepath.Join(wd, "kube_config.yml"))
	if err != nil {
		panic(err)
	}
	kubeConfig = string(kc)
}
