package workload

import (
	"context"

	"github.com/ahwhy/myGolang/k8s/meta"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ListDaemonSet(ctx context.Context, req *meta.ListRequest) (*appsv1.DaemonSetList, error) {
	return c.appsv1.DaemonSets(req.Namespace).List(ctx, req.Opts)
}

func (c *Client) CreateDaemonSet(ctx context.Context, obj *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return c.appsv1.DaemonSets(obj.Namespace).Create(ctx, obj, v1.CreateOptions{})
}

func (c *Client) DeleteDaemonSet(ctx context.Context, req *meta.DeleteRequest) error {
	return c.appsv1.DaemonSets(req.Namespace).Delete(ctx, req.Name, req.Opts)
}

func GetDaemonSetStatus(*appsv1.DaemonSet) *WorkloadStatus {
	return NewWorklaodStatus()
}
