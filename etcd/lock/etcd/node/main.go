package main

import (
	"context"
	"time"

	"github.com/ahwhy/myGolang/etcd/lock/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *clientv3.Client
)

func main() {
	lock, err := etcd.NewEtcdMutex(context.Background(), client, "/registry/locks")
	if err != nil {
		panic(err)
	}

	lock.Lock()
	time.Sleep(10 * time.Second)
	lock.Unlock()
}

func init() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	
	client = c
}
