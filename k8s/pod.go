package k8s

import (
	"context"

	"github.com/go-playground/validator/v10"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	validate = validator.New()
)

type ListPodtRequest struct {
	Namespace string
}

func (c *Client) ListPod(ctx context.Context, req *ListPodtRequest) (*v1.PodList, error) {
	if req.Namespace == "" {
		req.Namespace = v1.NamespaceDefault
	}

	return c.client.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{})
}

type GetPodRequest struct {
	Namespace string
	Name      string
}

func (c *Client) GetPod(ctx context.Context, req *GetPodRequest) (*v1.Pod, error) {
	if req.Namespace == "" {
		req.Namespace = v1.NamespaceDefault
	}

	return c.client.CoreV1().Pods(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
}

func (c *Client) DeletePod(ctx context.Context) error {
	return c.client.CoreV1().Pods("").Delete(ctx, "", metav1.DeleteOptions{})
}
