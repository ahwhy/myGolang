package watch

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *clientv3.Client
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
)

func WatchConfig(key string) {
	wc := client.Watch(context.Background(), key, clientv3.WithPrefix())
	for v := range wc {
		e := v.Events[0]
		fmt.Println("revision", v.Header.Revision, e.Type, e.Kv)
	}
}

func UpdateConfig(data string) {
	// put
	key := "/registry/configs/default/cmdb"
	putResp, err := client.Put(ctx, key, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(putResp)
}

func DeleteConfig() {
	// del
	key := "/registry/configs/default/cmdb"
	putResp, err := client.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		panic(err)
	}
	fmt.Println(putResp)
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
