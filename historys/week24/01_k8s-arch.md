# K8s架构与Etcd使用

![](./images/ci_cd.png)

workflow的架构思想借鉴了k8s, 因此在开发workflow之前, 我们先介绍下开k8s的架构和etcd的一些基础用法

## watch-list

![](./images/k8s-watch-list.png)

+ 对象存储在Etcd 以Kev Value的方式
```
prefix + "/" + 资源类型 + "/" + namespace + "/" + 具体资源名
例如：/registry/deployments/default/helloworld <object>
```

+ 监听对象的变化，面向事件编程
```
比如你修改了该Object
/registry/deployments/default/helloworld <object>
  |
  |
  V
node: 收到一个PUT事件, 并且把对象的key和Value的数据给到node

为了提升性能, 事件会在node本地缓存, 然后把事件通知真正需要处理业务的控制器
```

+ 面向期望的编程模式
```
控制器 获取变更的操作和对象

目标对象: helloworld deploy对象, 期望你把hello world部署3份
  |
  |
  V
node: 当前没有-----各种操作----> 完成期望
```

## etcd 搭建

由于我们的测试使用，因此使用docker搭建单节点etcd:
```sh
# windows上注意不要使用绝对路径: /usr/local/bin/etcd
# --listen-client-urls, --advertise-client-urls 必须带上, 后面使用api是的时候需要, 不然client 访问不到
docker run \
  -itd \
  -p 2379:2379 \
  -p 2380:2380 \
  --name etcd quay.io/coreos/etcd:latest etcd \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379
```

通过命令查看当前etcd的版本
```sh
$ docker exec -it etcd  etcdctl -version
etcdctl version: 3.3.8
API version: 2
```

查看当前实例
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl member list -w table
+------------------+---------+---------+-----------------------+-----------------------+
|        ID        | STATUS  |  NAME   |      PEER ADDRS       |     CLIENT ADDRS      |
+------------------+---------+---------+-----------------------+-----------------------+
| 8e9e05c52164694d | started | default | http://localhost:2380 | http://localhost:2379 |
+------------------+---------+---------+-----------------------+-----------------------+
```

注意: etcd容器没有shell, 你可以把他当做一个 二进制包来使用, 只是名字有点长而已

```sh
$ docker exec -it etcd  etcdctl 
NAME:
   etcdctl - A simple command line client for etcd.

WARNING:
   Environment variable ETCDCTL_API is not set; defaults to etcdctl v2.
   Set environment variable ETCDCTL_API=3 to use v3 API or ETCDCTL_API=2 to use v2 API.

USAGE:
   etcdctl [global options] command [command options] [arguments...]

VERSION:
   3.3.8
   
COMMANDS:
     backup          backup an etcd directory
     cluster-health  check the health of the etcd cluster
     mk              make a new key with a given value
     mkdir           make a new directory
     rm              remove a key or a directory
     rmdir           removes the key if it is an empty directory or a key-value pair
     get             retrieve the value of a key
     ls              retrieve a directory
     set             set the value of a key
     setdir          create a new directory or update an existing directory TTL
     update          update an existing key with a given value
     updatedir       update an existing directory
     watch           watch a key for changes
     ...
```

etcd有2个版本, 默认是v2, ETCDCTL_API=3 来使用v3版本的
```
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl
NAME:
        etcdctl - A simple command line client for etcd3.

USAGE:
        etcdctl

VERSION:
        3.3.8

API VERSION:
        3.3


COMMANDS:
        get                     Gets the key or a range of keys
        put                     Puts the given key into the store
        del                     Removes the specified key or range of keys [key, range_end)
        txn                     Txn processes all the requests in one transaction
        compaction              Compacts the event history in etcd
        ...
```

## etcd cli基本操作

etcd 是kv数据库, 我们的所有操作都是对存在数据库里面的key-value进行操作

下面我们演示一个简单的基于etcd的配置中心的流程:
+ key设计: /registry/configs/default/app_name
+ value设计: json object

### 写入数据(包含修改)

通过put 可以往数据库里面添加一条数据
```sh
# 写入数据
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object"
OK

# 读取数据
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get  /registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object

## 注意, 这里由于是windows的原因, key多出了一部分: C:/Program Files/Git
## 估计和windows下的 git bash有关系, 直接使用etcd client是没有这个问题的
```

我们我们再次往这个key写入数据，会覆盖之前的值, 也就实现了修改的效果
```sh
# 写入数据
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v2"
OK

# 读取数据
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get  /registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v2
```

### 读取数据

读取数据采用get指令, 最基本的方法就是 get key_name, 这种方法只能读取一个key, 如果我们想要读取多个key, 需要通过--prefix, 这样就能读取前缀为 key_prefix 的所有key， 也就实现了list功能

因此我们在添加一个key: workflow的配置
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/workflow "workflow config object v1"
OK
```

然后我们查看 当前注册的所有配置
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  /registry/configs
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v2
C:/Program Files/Git/registry/configs/default/workflow
workflow config object v1
```

+ 一般我们获取数据的时候，都希望最近添加的数据方前面，也就是按照时间的倒叙, 可以使用 --order="DESCEND" 
```
# Order of results; ASCEND or DESCEND (ASCEND by default)
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --order="DESCEND"  /registry/configs
C:/Program Files/Git/registry/configs/default/workflow
workflow config object v1
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v2
```

+ 你也可以只获取value, 不查询key
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --print-value-only  /registry/configs
cmdb config object v2
workflow config object v1
```

+ 你也可以只查询key, 不获取value, 这个就相当于索引了
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --keys-only  /registry/configs
C:/Program Files/Git/registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/workflow
```

+ 如果可以过多, 我们可以使用limit 来限制返回的kv个数
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --keys-only --limit 1  /registry/configs
C:/Program Files/Git/registry/configs/default/cmdb
```

+ 相比于我们操作其他数据库，还有一个关键性需求就是offset, 当数据过多的时候，我们都需要分页读取, 着就是需要使用, 现在的etcd 并不能
```sh
# 当前的keys
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  --keys-only   /registry/configs
C:/Program Files/Git/registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/keyauth
C:/Program Files/Git/registry/configs/default/workflow

## 我们要查询从keyauth开始的后面的key, 当然你也可以指定limit
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --from-key  --keys-only --limit=2  /registry/configs/default/keyauth
C:/Program Files/Git/registry/configs/default/keyauth
C:/Program Files/Git/registry/configs/default/workflow

# 可以看出，如果你要把他当场景的mysql或者mongo使用 是比较蓝瘦的, 你并不能指望数据来给你处理很多业务(MySQL 存储过程， Mongo MapReduce和函数)
```

### 版本

etcd的value是有版本概念的, 我们我们之前修改过一次cmdb的配置, 当前是v2的配置
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get  /registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v2
```

其实没次key有修改的时候 都会返回一个该key的版本号的, 但是需要-w json 才能查看到, 下面我么修改2次cmdb的配置
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v3" -w json
{"header":{"cluster_id":14841639068965178418,"member_id":10276657743932975437,"revision":6,"raft_term":2}}

$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v4" -w json
{"header":{"cluster_id":14841639068965178418,"member_id":10276657743932975437,"revision":7,"raft_term":2}}
```

如果我们没指定获取那个版本, 默认获取最新版本, 如果要获取历史版本需要--rev指定版本号
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get   /registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v4
# 获取版本6
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --rev=6   /registry/configs/default/cmdb
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v3
```

因此要使用etcd作为配置中心， 你需要保持好key的版本号

### 删除

现在我们workflow服务下线了, 我需要删除他的配置怎么办? etcd 可以通过del 删除指定key, 也可以指定--prefix 删除一批key
```
# 删除后返回删除的个数
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl del   /registry/configs/default/workflow
1

# 再次查看 workflow的配置已经删除了
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  /registry/configs
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v4
C:/Program Files/Git/registry/configs/default/keyauth
keyauth config object v1
```

### watch

像k8s里面的 watch list 就是使用的 该功能, 我们可以通过watch 一个prefix key， 当这个key有变化的时候 我们可以收到变化的数据

模拟node节点, watch 所有的config变化事件
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl watch --prefix  /registry/configs
```

模拟API server 修改 cmdb的配置
```sh
docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v5"
```

此时node就会收到该事件
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl watch --prefix  /registry/configs
PUT
C:/Program Files/Git/registry/configs/default/cmdb
cmdb config object v5
```

### lease和ttl

如果我们要为key 设置一个TTL喃?, 比如 这个key 存活时间为60s

+ 首先我们需要生成一个租约, 使用lease命名进行操作
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lease grant 60
lease 694d7d5f3050ef39 granted with TTL(60s)
```

+ 然后我们创建key的时候带上该租约
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/ttlkey --lease=694d7d5f3050ef39 "key with ttl"
OK
```

+ 查看key
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get /registry/configs/default/ttlkey
C:/Program Files/Git/registry/configs/default/ttlkey
key with ttl

# 60s秒过后, 该key就查不到了
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get /registry/configs/default/ttlkey
```

### 分布式锁

比如我们一个服务启动了3个副本, 在修改数据A的时候需要 先获取锁才能修改, 不然修改顺序就乱了, 怎么办?

通过lock 来创建一个锁, 该锁未释放之前, 其他想要获取该锁的实例都会阻塞

+ 模拟node1 获取了1个锁 /registry/configs/default/lockkey
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lock /registry/configs/default/lockkey
C:/Program Files/Git/registry/configs/default/lockkey/694d7d5f3050ef43
```

+ 模拟node2 获取key
```
# 直到node1释放后，node2才能获取到锁
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lock /registry/configs/default/lockkey
C:/Program Files/Git/registry/configs/default/lockkey/694d7d5f3050ef55
```

## etcd client使用

etcd client的使用逻辑和cli基本一致

+ client: go.etcd.io/etcd/client/v3 
+ 版本要求: v3.5.1(最好大于3.5)

### 初始化client

```go
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
}
```

正常的情况下, 我们可以打印出当前etcd的member节点信息
```
&{cluster_id:14841639068965178418 member_id:10276657743932975437 raft_term:2  [ID:10276657743932975437 name:"default" peerURLs:"http://localhost:2380" clientURLs:"http://0.0.0.0:2379" ] {} [] 0}
```

### put

```go
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
fmt.Println(getResp)
```

### get

```go
// get with prefix
getResp, err = client.Get(ctx, key, clientv3.WithPrefix())
if err != nil {
    panic(err)
}
fmt.Println(getResp.Kvs)
```

get 支持多个参数, 和命令行的含义一样, 只是这里使用编程的opt语法传入

### del

```go
// delete
delResp, err := client.Delete(ctx, key, clientv3.WithPrevKV())
if err != nil {
    panic(err)
}
fmt.Println(delResp.PrevKvs)
```
注意 这里添加一个参数WithPrevKV, 默认情况下 delete是无法获取到被删除的值得, 通过添加该参数可以 获取当前被删除的值

你也可以通过命令行确认结果
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  /registry/configs
```

### watch

先编写一个用于测试的包
```go
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
		fmt.Println(v)
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
```

api server 负责修改对象
```go
package main

import "gitee.com/infraboard/go-course/day24/etcd/watch"

func main() {
	watch.UpdateConfig("cmdb v3")
}
```

controller 负责watch 对象变化
```go
package main

import "gitee.com/infraboard/go-course/day24/etcd/watch"

func main() {
	watch.WatchConfig("/registry/configs")
}
```

然后我们测试修改
+ 添加
+ 修改
+ 删除

```sh
$ go run controler/main.go
revision 12 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:12 version:7 value:"cmdb v3" 
revision 13 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:13 version:8 value:"cmdb v3"
revision 14 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:14 version:9 value:"cmdb v3"
revision 15 DELETE key:"/registry/configs/default/cmdb" mod_revision:15
```

### lock

之前 1000 goroutine 修改全局变量累加的例子:
```go
package main

import (
	"sync"
)

// 全局变量
var counter int

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	println(counter)
}
```

当在同一个进程的时候我们可以使用: 互斥锁，可以解决并行抢占的问题
```go
func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			lock.Lock()
			counter++
			lock.Unlock()
		}()
	}

	wg.Wait()
	println(counter)
}
```

如果我们不是在一个进程内部, 比如 你起了3个进程, 这3个进程 可能分布在不同的主机上, 这个时候进程锁 就无法满足需求了, 需要使用分布式锁

可以实现分布式锁能力的服务主要有(主要看数据一致性模型)
+ 基于Redis的setnx, 不太推荐
+ 基于ZooKeeper, paxios算法保证
+ 基于etcd， raft算法保证

当然我们选择使用etcd实现, etcd 提供的 concurrency 就是解决分布式并发问题的: 

![](./images/etcd-lock.png)

下面我们利用etcd 实现一个分布式互斥锁
```go
package etcd

import (
	"context"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// 初始化sync.Locker对象.
func NewEtcdMutex(ctx context.Context, client *clientv3.Client, pfx string) (sync.Locker, error) {
	// 创建Session对象.
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

// 申请锁.
func (l *EtcdMutex) Lock() {
	// 不是标准的sync.Locker接口,需要传入Context对象,在获取锁时可以设置超时时间,或主动取消请求.
	err := l.m.Lock(l.ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("获取锁")
}

// 释放锁.
func (l *EtcdMutex) Unlock() {
	err := l.m.Unlock(l.ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("释放锁")
}
```

写个node模拟多进程 获取互斥锁:
```go
package main

import (
	"context"
	"time"

	"gitee.com/infraboard/go-course/day24/etcd/lock/etcd"
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
```

测试2个锁 是否能循环

如果你不想获取锁的时候阻塞, 比如当前有人持有锁，直接放弃抢占锁, 可以使用Trylock
```go
// TryLock locks the mutex if not already locked by another session.
// If lock is held by another session, return immediately after attempting necessary cleanup
// The ctx argument is used for the sending/receiving Txn RPC.
func (m *Mutex) TryLock(ctx context.Context) error
```


## 关于kv设计

因为etcd 是kv数据库, 所以没有where之类的操作, 如果你想要设置索引过滤对象, 只能设计好你的key

比如:
```sh
/registry/configs/namesapce/resource_name
```

## 注意事项

历史版本越多，存储空间越大，性能越差，直到etcd到达空间配额限制的时候，etcd的写入将会被禁止变为只读，影响线上服务，因此这些历史版本需要进行压缩

数据压缩并不是清理现有数据，只是对给定版本之前的历史版本进行清理，清理后数据的历史版本将不能访问，但不会影响现有最新数据的访问

```sh
etcdctl compact 5
```

## 总结

+ watch list 设计理念
+ etcd client基本操作
+ watch and lock

基础准备好了后，我们接下来 开始workflow的 API Server的开发

## 参考

+ [etcd 问题、调优、监控](https://www.kubernetes.org.cn/7569.html)







