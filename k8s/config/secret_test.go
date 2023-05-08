package config_test

import (
	"testing"

	"github.com/ahwhy/myGolang/k8s/meta"
)

func TestListSecret(t *testing.T) {
	req := meta.NewListRequest()
	v, err := impl.ListSecret(ctx, req)
	if err != nil {
		t.Log(err)
	}
	t.Log(v)
}

func TestGetSecret(t *testing.T) {
	req := meta.NewGetRequest("test-secret")
	v, err := impl.GetSecret(ctx, req)

	if err != nil {
		t.Log(err)
	}
	t.Log(v.Name)
}
