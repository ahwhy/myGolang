package watch_test

import (
	"testing"

	"github.com/ahwhy/myGolang/etcd/watch"
)

func TestCreate(t *testing.T) {
	watch.UpdateConfig("cmdb v3")
}

func TestDelete(t *testing.T) {
	watch.DeleteConfig()
}
