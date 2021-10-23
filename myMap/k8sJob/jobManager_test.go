package k8sJob_test

import (
	"testing"
	"time"

	"github.com/ahwhy/myGolang/myMap/k8sJob"
)

func TestJobManager(t *testing.T) {
	jm := k8sJob.NewJobManager()
	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewK8sD("a"), k8sJob.NewK8sD("b"), k8sJob.NewK8sD("c")})
	time.Sleep(10 * time.Second)

	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewHost("A"), k8sJob.NewHost("B"), k8sJob.NewHost("C")})
	time.Sleep(10 * time.Second)

	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewK8sD("b"), k8sJob.NewK8sD("c"), k8sJob.NewK8sD("d")})
	time.Sleep(10 * time.Second)

	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewK8sD("c"), k8sJob.NewK8sD("d"), k8sJob.NewK8sD("e")})
	time.Sleep(10 * time.Second)

	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewHost("A"), k8sJob.NewHost("B"), k8sJob.NewHost("C"), k8sJob.NewHost("D")})
	time.Sleep(10 * time.Second)

	go jm.Sync([]k8sJob.DeployJob{k8sJob.NewK8sD("c"), k8sJob.NewK8sD("d"), k8sJob.NewK8sD("e"), k8sJob.NewK8sD("f")})
	time.Sleep(10 * time.Second)
}
