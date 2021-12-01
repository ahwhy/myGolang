package etcd

import (
	"context"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// 初始化sync.Locker对象
func NewEtcdMutex(ctx context.Context, client *clientv3.Client, pfx string) (sync.Locker, error) {
	// 创建Session对象
	// 在调用concurrency.NewSession时,会设置ttl,默认为60秒
	// Session对象会持有对应的LeaseID,并会调用KeepAlive来续期
	// 使得锁在Unlock之前一直是有效的,其它想抢占分布式锁的程序只能是等待
	sess, err := concurrency.NewSession(client)
	if err != nil {
		return nil, err
	}

	// 创建Mutex对象. 需要指定锁的名称, 和命令行使用lock一样，就是key的prefix
	m := concurrency.NewMutex(sess, "/registry/locks")

	return &EtcdMutex{
		sess: sess,
		m:    m,
		pfx:  pfx,
		ctx:  ctx,
	}, nil
}

type EtcdMutex struct {
	sess *concurrency.Session
	m    *concurrency.Mutex
	pfx  string
	ctx  context.Context
}

// 申请锁
func (l *EtcdMutex) Lock() {
	// 不是标准的sync.Locker接口,需要传入Context对象,在获取锁时可以设置超时时间,或主动取消请求
	err := l.m.Lock(l.ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("获取锁")
}

// 释放锁
func (l *EtcdMutex) Unlock() {
	err := l.m.Unlock(l.ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("释放锁")
}
