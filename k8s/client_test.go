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
	system     = "default"
	pod        = "demo-87d85fccd-nt9n9"
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

func TestPod(t *testing.T) {
	should := assert.New(t)
	client, err := k8s.NewClient(kubeConfig)
	should.NoError(err)

	ctx := context.Background()
	list := k8s.NewListPodtRequest(system)
	v, err := client.ListPod(ctx, list)
	should.NoError(err)
	fmt.Printf("%v\n", v.Items)
	fmt.Println("------------------------")

	get := k8s.NewGetPodRequest(system, pod)
	v2, err := client.GetPod(ctx, get)
	should.NoError(err)
	fmt.Printf("%v\n", v2)
	fmt.Println("------------------------")

	del := k8s.NewDeletePodRequest(system, pod)
	err = client.DeletePod(ctx, del)
	should.NoError(err)
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
