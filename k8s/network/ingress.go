package network

import (
	"context"

	"github.com/ahwhy/myGolang/k8s/meta"
	v1 "k8s.io/api/networking/v1"
)

func (c *Client) ListIngress(ctx context.Context, req *meta.ListRequest) (*v1.IngressList, error) {
	return c.networkingv1.Ingresses(req.Namespace).List(ctx, req.Opts)
}
