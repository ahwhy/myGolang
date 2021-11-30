package basic

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.MemberList(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	// put
	key := "/registry/configs/default/cmdb"
	putResp, err := client.Put(ctx, key, "cmdb config v1")
	if err != nil {
		panic(err)
	}
	fmt.Println(putResp)

	// get
	getResp, err := client.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(getResp.Kvs)

	// get with prefix
	getResp, err = client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	fmt.Println(getResp.Kvs)

	// delete
	delResp, err := client.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		panic(err)
	}
	fmt.Println(delResp.PrevKvs)
}
