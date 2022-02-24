package k8s

import (
	"context"

	v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewListCronJobRequest(namespace string) *ListCronJobRequest {
	return &ListCronJobRequest{
		Namespace: namespace,
	}
}

type ListCronJobRequest struct {
	Namespace string
}

func (c *Client) ListCronJob(ctx context.Context, req *ListCronJobRequest) (*v1.CronJobList, error) {
	if req.Namespace == "" {
		req.Namespace = apiv1.NamespaceDefault
	}

	return c.client.BatchV1().CronJobs(req.Namespace).List(ctx, metav1.ListOptions{})
}
