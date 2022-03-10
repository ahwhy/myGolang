# Golang-Etcd  Golang的Etcd

## 一、分布式应用

### 1. 分布式应用协助模型(分布式程序高可用方式)
- 基于协作协议，leader election
	- 通过程序某个端口，发送心跳，用协议通信
	- 如: etcd、elasticsearch
		- elasticsearch 基于9100端口，发送心跳，通过选举协议(leader election)选举主节点(leader)
		- 其他节点工作于备用模式或观察者模式或学习者模式
		- leader出现故障，剩下的人一而上进行下一轮选举，选举新的leader

- 基于分布式锁，leader election
	- 借助第三方组件(分布式锁)，来选举leader
	- 如：zookeeper、kube-controller-manager、kube-scheduler
		- zookeeper 通过分布式锁来选举leader，区别谁是主动，谁是被动
		- 原子发现协议，强一致性

- 事务的ACID特性
	- Atomicity: 原子性
	- Consistency: 一致性
	- Isolation: 隔离性
	- Durability: 持久性

- CAP理论，指在任何一个分布式系统中都无法同时满足CAP
	- C(Consistency): 表示一致性，所有的节点同一时间看到的是相同的数据
	- A(Avaliablity): 表示可用性，不管是否成功，确保一个请求都能接收到响应
	- P(Partion Tolerance): 分区容错性，系统任意分区后，在网络故障时，仍能操作

- BASE理论
	- Basically Available  基本可用
	- Soft State  软状态
	- Eventually Consistent 最终一致性 
	- 三个短语的缩写，核心思想: 既是无法做到强一致性(Strong consistency)，但每个应用都可以根据自身的业务特点，采用适当的方式来使系统达到最终一致性(Eventual consistency)
	- [BASE理论参考](https://juejin.cn/post/6844903621495095304)
	- [强一致性、弱一致性、最终一致性、读写一致性、单调读、因果一致性 的区别与联系](https://zhuanlan.zhihu.com/p/67949045)

## 二、Etcd

### 1. Etcd简介
- go语言开发的分布式应用
	- k/v存储，不需要其他依赖，所有操作都是对存在数据库里面的key-value进行操作
	- etcd运行多个实例时，直接基于raft协议，完成leader election，进而完成数据强一致性
	- 每一个节点均可以读写，但是写到任意一个节点的数据，都要同步到同一个集群的另外的节点，而且确保数据是强一致的
	- 特性: leader election 、 数据强一致

- raft协议是用于维护一组服务节点数据一致性的协议
	- raft协议是简装版的paxos协议，功能并不比其弱
	- 在今天各种分布式协作逻辑中，出现很多协议，都沿用paxos协议的思想，或简化、或另辟蹊径，都或多或少受到paxos协议的影响
	- 原作者穷10年之功，才设计出paxos
	- java开发分布式程序，就会使用paxos协议或 google研发的paxos变种zab协议(Zookeeper Atomic Broadcast)
	- [协议工作逻辑](https://www.bilibili.com/video/av77388641/)
	- 基于Redis，setnx
	- 基于ZooKeeper，paxios算法保证
	- 基于etcd，raft算法保证

### 2. Etcd搭建
```shell
	# 使用docker搭建单节点etcd
	# windows上注意不要使用绝对路径: /usr/local/bin/etcd
	# --listen-client-urls, --advertise-client-urls 必须带上, 后面使用api是的时候需要, 不然client 访问不到
	docker run \
	-itd \
	-p 2379:2379 \
	-p 2380:2380 \
	--name etcd quay.io/coreos/etcd:latest etcd \
	--listen-client-urls http://0.0.0.0:2379 \
	--advertise-client-urls http://0.0.0.0:2379
	
	# 通过命令查看当前etcd的版本
	$ docker exec -it etcd  etcdctl -version
	etcdctl version: 3.3.8
	API version: 2
	
	# 查看当前实例
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl member list -w table
	+------------------+---------+---------+-----------------------+-----------------------+
	|        ID        | STATUS  |  NAME   |      PEER ADDRS       |     CLIENT ADDRS      |
	+------------------+---------+---------+-----------------------+-----------------------+
	| 8e9e05c52164694d | started | default | http://localhost:2380 | http://localhost:2379 |
	+------------------+---------+---------+-----------------------+-----------------------+
```

### 3. Etcdctl集群管理
```shell
	# Etcd集群启动参数
	ExecStart=/usr/bin/etcd --election-timeout=3000 --auto-compaction-retention=1 --auto-compaction-mode=revision --heartbeat-interval=500 --cert-file=/var/lib/etcd/cert/etcd-server.pem --key-file=/var/lib/etcd/cert/etcd-server-key.pem --peer-client-cert-auth --peer-trusted-ca-file=/var/lib/etcd/cert/peer-ca.pem --peer-cert-file=/var/lib/etcd/cert/192.168.22.12-name-2.pem --peer-key-file=/var/lib/etcd/cert/192.168.22.12-name-2-key.pem --initial-advertise-peer-urls https://192.168.22.12:2380 --listen-peer-urls https://192.168.22.12:2380 --advertise-client-urls https://192.168.22.12:2379 --listen-client-urls https://192.168.22.12:2379 --initial-cluster 192.168.22.11-name-1=https://192.168.22.11:2380,192.168.22.12-name-2=https://192.168.22.12:2380,192.168.22.13-name-3=https://192.168.22.13:2380 --initial-cluster-state new --initial-cluster-token e21a2460-2096-49a9-b5d0-650362539766

	# Etcd集群健康检测
	$ export NODE_IPS="192.168.22.11 192.168.22.12 192.168.22.13" 
	$ for ip in ${NODE_IPS}; do ETCDCTL_API=3  etcdctl  --endpoints=https://${ip}:2379  --cacert=/etc/etcd/pki/ca.crt --cert=/etc/etcd/pki/client.crt --key=/etc/etcd/pki/client.key   endpoint health; done 
	https://192.168.22.11:2379 is healthy: successfully committed proposal: took = 2.487239ms
	https://192.168.22.12:2379 is healthy: successfully committed proposal: took = 1.77157ms
	https://192.168.22.13:2379 is healthy: successfully committed proposal: took = 3.064988ms

	# 集群密钥校验
	ETCDCTL_API=3 etcdctl --endpoints="https://192.168.22.11:2379,https://192.168.22.12:2379,https://192.168.22.13:2379" --cacert=/etc/kubernetes/pki/etcd/ca.pem --cert=/etc/kubernetes/pki/etcd/etcd-client.pem --key=/etc/kubernetes/pki/etcd/etcd-client-key.pem -h

	# 集群metrics数据获取
	curl -k --cert /etc/kubernetes/pki/etcd/etcd-client.pem --key /etc/kubernetes/pki/etcd/etcd-client-key.pem https://192.168.22.13:2379/metrics

	# 以下命令若不条件对应证书路径，需在配置文件中添加：http://127.0.0.1:2379
	# --endpoints=https://${ip}:2379  --cacert=/etc/etcd/pki/ca.crt --cert=/etc/etcd/pki/client.crt --key=/etc/etcd/pki/client.key
	# docker exec 容器ID sh -c "ETCDCTL_API=3 etcdctl --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key --endpoints="https://192.168.22.11:2379,https://192.168.22.12:2379,https://192.168.22.13:2379" endpoint health" 
	# docker exec 容器ID sh -c "ETCDCTL_API=3 etcdctl --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key --endpoints="https://192.168.22.11:2379,https://192.168.22.12:2379,https://192.168.22.13:2379" endpoint status" 
	# https://blog.csdn.net/weixin_30469895/article/details/99194344
	# docker cp $(docker ps |grep k8s_etcd_etcd |awk '{print $NF}'):/usr/local/bin/etcdctl  /usr/local/bin
	$ ETCDCTL_API=3  etcdctl  --help
	$ ETCDCTL_API=3  etcdctl  member list   --
	$ ETCDCTL_API=3  etcdctl  member remove id
	$ ETCDCTL_API=3  etcdctl  member add infra2 --peer-urls="https://192.168.22.13:2380"   #先加入再启动节点
	
	# 集群数据备份 https://www.cnblogs.com/chenqionghe/p/10622859.html
	$ ETCDCTL_API=3  etcdctl  snapshot save     snapshot.db                              #数据备份
	$ ETCDCTL_API=3  etcdctl  snapshot restore  snapshot.db  --data-dir=/opt/etcd-testdir #将数据恢复到一个新的不存在的目录中
```

### 4. Etcdctl基本命令
- 写入/修改数据
```shell
	# 写入数据
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object"
	OK
	
	# 读取数据
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get  /registry/configs/default/cmdb
	C:/Program Files/Git/registry/configs/default/cmdb
	cmdb config object
	## PS: 这里由于是windows的原因, key多出了一部分: C:/Program Files/Git, 同git bash有关
	
	# 再次往这个key写入数据，会覆盖之前的值, 也就实现了修改的效果
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v2"
	OK
	
	# 读取数据
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get  /registry/configs/default/cmdb -w json
	C:/Program Files/Git/registry/configs/default/cmdb
	{"header":{"cluster_id":14841639068965178418,"member_id":10276657743932975437,"revision":47,"raft_term":2},"kvs":[{"key":"L3JlZ2lzdHJ5L2NvbmZpZ3MvZGVmYXVsdC9jbWRi","create_revision":37,"mod_revision":47,"version":11,"value":"Y21kYiBjb25maWcgb2JqZWN0djk="}],"count":1}

```

- 读取数据
```shell
	# 读取数据采用get指令, 最基本的方法就是 get key_name, 这种方法只能读取一个key,
	# 如果想要读取多个key, 需要通过--prefix, 这样就能读取前缀为 key_prefix 的所有key, 从而实现list功能
	# 添加一个key: workflow的配置
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/workflow "workflow config object v1"
	OK
	
	# 查看当前注册的所有配置
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  /registry/configs
	C:/Program Files/Git/registry/configs/default/cmdb
	cmdb config object v2
	C:/Program Files/Git/registry/configs/default/workflow
	workflow config object v1
	
	# 一般获取数据的时候，都希望最近添加的数据方前面，也就是按照时间的倒叙, 可以使用 --order="DESCEND"
	# Order of results; ASCEND or DESCEND (ASCEND by default)
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --order="DESCEND"  /registry/configs
	C:/Program Files/Git/registry/configs/default/workflow
	workflow config object v1
	C:/Program Files/Git/registry/configs/default/cmdb
	cmdb config object v2
	
	# 只获取value, 不查询key
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --print-value-only  /registry/configs
	cmdb config object v2
	workflow config object v1
	
	# 只查询key, 不获取value, 相当于索引
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --keys-only  /registry/configs
	C:/Program Files/Git/registry/configs/default/cmdb
	C:/Program Files/Git/registry/configs/default/workflow
	
	# 如果数据过多, 可以使用limit 来限制返回的kv个数
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix --keys-only --limit 1  /registry/configs
	C:/Program Files/Git/registry/configs/default/cmdb
	
	# 相比于操作其他数据库，还有一个关键性需求就是offset, 当数据过多的时候，都需要分页读取, 着就是需要使用, 现在的etcd 并不能
	# 当前的keys
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  --keys-only   /registry/configs
	C:/Program Files/Git/registry/configs/default/cmdb
	C:/Program Files/Git/registry/configs/default/keyauth
	C:/Program Files/Git/registry/configs/default/workflow
	
	# 插入keyauth: docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/keyauth "keyauth config object"
	# 要查询从keyauth开始的后面的key, 也可以指定limit
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --from-key  --keys-only --limit=2  /registry/configs/default/keyauth
	C:/Program Files/Git/registry/configs/default/keyauth
	C:/Program Files/Git/registry/configs/default/workflow
```

- 删除数据
```shell
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

- watch
```shell
	# k8s里面的 watch/list 就是使用的 该功能, 可以通过watch 一个prefix key， 当这个key有变化的时候就可以收到变化的数据
	
	# 模拟node节点, watch 所有的config变化事件
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl watch --prefix  /registry/configs
	
	# 模拟API server 修改 cmdb的配置
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/cmdb "cmdb config object v5"
	
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl watch --prefix  /registry/configs
	PUT
	C:/Program Files/Git/registry/configs/default/cmdb
	cmdb config object v5
```

- lease和ttl
```shell
	# 为key 设置一个TTL, 比如 这个key 存活时间为60s
	
	# 首先需要生成一个租约, 使用lease命名进行操作
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lease grant 60
	lease 694d7d5f3050ef39 granted with TTL(60s)
	
	# 然后创建key的时候带上该租约
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl put /registry/configs/default/ttlkey --lease=694d7d5f3050ef39 "key with ttl"
	OK
	
	# 查看key
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get /registry/configs/default/ttlkey
	C:/Program Files/Git/registry/configs/default/ttlkey
	key with ttl
	
	# 60s秒过后, 该key就查不到了
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get /registry/configs/default/ttlkey
```

- 分布式锁
```shell
	# 一个服务启动了3个副本, 在修改数据A的时候需要 先获取锁才能修改
	# 通过lock 来创建一个锁, 该锁未释放之前, 其他想要获取该锁的实例都会阻塞
	
	# 模拟node1 获取了1个锁 /registry/configs/default/lockkey
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lock /registry/configs/default/lockkey
	C:/Program Files/Git/registry/configs/default/lockkey/694d7d5f3050ef43
	
	# 模拟node2 获取key
	# 直到node1释放后，node2才能获取到锁
	$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl lock /registry/configs/default/lockkey
	C:/Program Files/Git/registry/configs/default/lockkey/694d7d5f3050ef55
```

### 5. Eetcd client使用
- 初始化client
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
	// &{cluster_id:14841639068965178418 member_id:10276657743932975437 raft_term:2  [ID:10276657743932975437 name:"default" peerURLs:"http://localhost:2380" clientURLs:"http://0.0.0.0:2379" ] {} [] 0}
```

- put
```go
	// put
	key := "/registry/configs/default/cmdb"
	putResp, err := client.Put(ctx, key, "cmdb config v1")
	if err != nil {
		panic(err)
	}
	fmt.Println(putResp)
```

- get
```go
	// Get() 支持多个参数, 和命令行的含义一样, 这里使用编程的opt语法传入
	// get
	getResp, err := client.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(getResp)

	// get with prefix
	getResp, err = client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	fmt.Println(getResp.Kvs)
```

- delete
```go
	// delete
	// 参数WithPrevKV, 默认情况下 delete是无法获取到被删除的值得, 通过添加该参数获取当前被删除的值
	delResp, err := client.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		panic(err)
	}
	fmt.Println(delResp.PrevKvs)
```

- watch
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

	// api-server 负责修改对象
	watch.UpdateConfig("cmdb v3")
	
	// controller 负责watch对象变化
	watch.WatchConfig("/registry/configs")

	// 输出
	// revision 12 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:12 version:7 value:"cmdb v3" 
	// revision 13 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:13 version:8 value:"cmdb v3"
	// revision 14 PUT key:"/registry/configs/default/cmdb" create_revision:6 mod_revision:14 version:9 value:"cmdb v3"
	// revision 15 DELETE key:"/registry/configs/default/cmdb" mod_revision:15
```

- lock
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
	
	// 申请锁
	func (l *EtcdMutex) Lock() {
		// 不是标准的sync.Locker接口,需要传入Context对象,在获取锁时可以设置超时时间,或主动取消请求.
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
	
	// 模拟多进程 获取互斥锁
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
	
	// 如果不想获取锁的时候阻塞, 比如当前有人持有锁，直接放弃抢占锁, 可以使用Trylock
	// TryLock locks the mutex if not already locked by another session.
	// If lock is held by another session, return immediately after attempting necessary cleanup
	// The ctx argument is used for the sending/receiving Txn RPC.
	// func (m *Mutex) TryLock(ctx context.Context) error
```

### 6. Etcd集群选主
- 这一组服务节点构成一个Etcd集群，并且有一个Etcd主节点来对外提供服务
	- 当集群初始化，或者主节点挂掉后，面临一个选主问题
	- 集群中每个节点，任意时刻处于Leader（主）, Follower（从）, Candidate（候选）这三个角色之一

- 选举特点如下
	- 当集群初始化时候，每个节点都是Follower角色
	- 集群中存在至多1个有效的主节点，通过心跳与其他节点同步数据
	- 当Follower在一定时间内没有收到来自主节点的心跳，会将自己角色改变为Candidate，并发起一次选主投票；当收到包括自己在内超过半数节点赞成后，选举成功；当收到票数不足半数选举失败，或者选举超时。若本轮未选出主节点，将进行下一轮选举(出现这种情况，是由于多个节点同时选举，所有节点均为获得过半选票)
	- Candidate节点收到来自主节点的信息后，会立即终止选举过程，进入Follower角色
	- 为了避免陷入选主失败循环，每个节点未收到心跳发起选举的时间是一定范围内的随机值，这样能够避免2个节点同时发起选主

- 相关逻辑代码
```
	/*
	* 发起竞选
	* 未当选leader前，会一直阻塞在Campaign调用
	* 当选leader后，等待SIGINT、SIGTERM或session过期而退出
	* https://cloud.tencent.com/developer/article/1458456
	* https://github.com/etcd-io/etcd/blob/master/etcdctl/ctlv3/command/elect_command.go
	*/
	
	func campaign(c *clientv3.Client, election string, prop string) error {
		// NewSession函数中创建了一个lease，默认是60s TTL，并会调用KeepAlive，永久为这个lease自动续约(2/3生命周期的时候执行续约操作)
		s, err := concurrency.NewSession(c)
		if err != nil {
			return err
		}
		e := concurrency.NewElection(s, election)
		ctx, cancel := context.WithCancel(context.TODO())
		
		donec := make(chan struct{})
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigc
			cancel()
			close(donec)
		}()
		
		// 竞选逻辑，将展开分析
		if err = e.Campaign(ctx, prop); err != nil {
			return err
		}
		
		// print key since elected
		resp, err := c.Get(ctx, e.Key())
		if err != nil {
			return err
		}
		display.Get(*resp)
		
		select {
		case <-donec:
		case <-s.Done():
			return errors.New("elect: session expired")
		}
		
		return e.Resign(context.TODO())
	}

	/*
	* 类似于zookeeper的临时有序节点，etcd的选举也是在相应的prefix path下面创建key，该key绑定了lease并根据lease id进行命名
	* key创建后就有revision号，这样使得在prefix path下的key也都是按revision有序
	* https://github.com/etcd-io/etcd/blob/master/clientv3/concurrency/election.go
	*/
	func (e *Election) Campaign(ctx context.Context, val string) error {
		s := e.session
		client := e.session.Client()
		
		// 真正创建的key名为：prefix + lease id
		k := fmt.Sprintf("%s%x", e.keyPrefix, s.Lease())
		// Txn：transaction，依靠Txn进行创建key的CAS操作，当key不存在时才会成功创建
		txn := client.Txn(ctx).If(v3.Compare(v3.CreateRevision(k), "=", 0))
		txn = txn.Then(v3.OpPut(k, val, v3.WithLease(s.Lease())))
		txn = txn.Else(v3.OpGet(k))
		resp, err := txn.Commit()
		if err != nil {
			return err
		}
		e.leaderKey, e.leaderRev, e.leaderSession = k, resp.Header.Revision, s
		// 如果key已存在，则创建失败
		// 当key的value与当前value不等时，如果自己为leader，则不用重新执行选举直接设置value
		// 否则报错
		if !resp.Succeeded {
			kv := resp.Responses[0].GetResponseRange().Kvs[0]
			e.leaderRev = kv.CreateRevision
			if string(kv.Value) != val {
				if err = e.Proclaim(ctx, val); err != nil {
					e.Resign(ctx)
					return err
				}
			}
		}
		
		// 一直阻塞，直到确认自己的create revision为当前path中最小，从而确认自己当选为leader
		_, err = waitDeletes(ctx, client, e.keyPrefix, e.leaderRev-1)
		if err != nil {
			// clean up in case of context cancel
			select {
			case <-ctx.Done():
				e.Resign(client.Ctx())
			default:
				e.leaderSession = nil
			}
			return err
		}
		e.hdr = resp.Header
	
		return nil
	}
```

### 7. 其他
- 注意事项
	- 关于k/v设计
		- etcd 是kv数据库，没有where之类的操作
		- 如果想要设置索引过滤对象，需要提前设计key
	- 存储限制
		- 历史版本越多，存储空间越大，性能越差，直到etcd到达空间配额限制的时候，etcd的写入将会被禁止变为只读，影响线上服务，因此这些历史版本需要进行压缩
		- 数据压缩并不是清理现有数据，只是对给定版本之前的历史版本进行清理，清理后数据的历史版本将不能访问，但不会影响现有最新数据的访问
		- `etcdctl compaction 5`
		
- 参考文档
	- [etcd 问题、调优、监控](https://www.kubernetes.org.cn/7569.html)
	- [etcdv3.4 官网启动参数](https://etcd.io/docs/v3.4/op-guide/configuration/)
	- [etcdv3.4 CHANGELOG-3.4](https://github.com/etcd-io/etcd/blob/main/CHANGELOG-3.4.md)
	- [ETCD-内部原理](https://www.cnblogs.com/seattle-xyt/p/10366131.html)