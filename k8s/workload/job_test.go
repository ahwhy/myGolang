package workload_test

import (
	"testing"

	"github.com/ahwhy/myGolang/k8s/meta"
	"github.com/ahwhy/myGolang/utils/tools"

	v1 "k8s.io/api/batch/v1"
)

func TestListJob(t *testing.T) {
	req := meta.NewListRequest()
	req.Namespace = "default"
	list, err := impl.ListJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	// 序列化
	for _, v := range list.Items {
		t.Log(tools.MustToYaml(v))
	}
}

func TestGetJob(t *testing.T) {
	req := meta.NewGetRequest("job-test")
	req.Namespace = "default"
	ins, err := impl.GetJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	// 序列化
	t.Log(tools.MustToYaml(ins))
}

func TestCreateJob(t *testing.T) {
	job := &v1.Job{}
	tools.MustReadYamlFile("job.yml", job)
	job, err := impl.CreateJob(ctx, job)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(job)
}

func TestDeleteJob(t *testing.T) {
	req := meta.NewDeleteRequest("job-test")
	err := impl.DeleteJob(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
}
