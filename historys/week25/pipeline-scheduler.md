# Pipeline调度器开发

Pipeline调度器的核心逻辑:
+ Watch Pipeline对象, 监听变化事件
+ 当有新的Pipeline对象被创建时, 修改Pipeline对象的 Scheduler属性, 为其挑选一个可以Scheduler Node来处理Pipline

![](./images/k8s-watch-list.png)

因此第一步是编写Informer, Watch Pipeline对象的变化

## PIpeline Informer

### 定义接口

我们先定义Informer的接口
```go
package pipeline

import (
	"context"

	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/common/cache"
)

// Informer 负责事件通知
type Informer interface {
	Watcher() Watcher
	Lister() Lister
	Recorder() Recorder
	GetStore() cache.Store
}

type Lister interface {
	List(ctx context.Context, opts *pipeline.QueryPipelineOptions) (*pipeline.PipelineSet, error)
}

type Recorder interface {
	Update(*pipeline.Pipeline) error
}

// Watcher 负责事件通知
type Watcher interface {
	// Run starts and runs the shared informer, returning after it stops.
	// The informer will be stopped when stopCh is closed.
	Run(ctx context.Context) error
	// AddEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddPipelineTaskEventHandler(handler PipelineEventHandler)
}

// PipelineEventHandler can handle notifications for events that happen to a
// resource. The events are informational only, so you can't return an
// error.
//  * OnAdd is called when an object is added.
//  * OnUpdate is called when an object is modified. Note that oldObj is the
//      last known state of the object-- it is possible that several changes
//      were combined together, so you can't use this to see every single
//      change. OnUpdate is also called when a re-list happens, and it will
//      get called even if nothing changed. This is useful for periodically
//      evaluating or syncing something.
//  * OnDelete will get the final state of the item if it is known, otherwise
//      it will get an object of type DeletedFinalStateUnknown. This can
//      happen if the watch is closed and misses the delete event and we don't
//      notice the deletion until the subsequent re-list.
type PipelineEventHandler interface {
	OnAdd(obj *pipeline.Pipeline)
	OnUpdate(old, new *pipeline.Pipeline)
	OnDelete(obj *pipeline.Pipeline)
}

// PipelineEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.
type PipelineTaskEventHandlerFuncs struct {
	AddFunc    func(obj *pipeline.Pipeline)
	UpdateFunc func(oldObj, newObj *pipeline.Pipeline)
	DeleteFunc func(obj *pipeline.Pipeline)
}

// OnAdd calls AddFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnAdd(obj *pipeline.Pipeline) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnUpdate(oldObj, newObj *pipeline.Pipeline) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnDelete(obj *pipeline.Pipeline) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}

type PipelineFilterHandler func(obj *pipeline.Pipeline) error
```

### Etcd实现

结下来我们使用etcd 来实现定义


#### 实现Lister

lister其实就是一个带 prefix的 Get操作

```go
type lister struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *lister) List(ctx context.Context, opts *pipeline.QueryPipelineOptions) (*pipeline.PipelineSet, error) {
	listKey := pipeline.EtcdPipelinePrefix()
	l.log.Infof("list etcd pipeline resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	ps := pipeline.NewPipelineSet()
	for i := range resp.Kvs {
		// 解析对象
		pt, err := pipeline.LoadPipelineFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Errorf("load pipeline [key: %s, value: %s] error, %s", resp.Kvs[i].Key, string(resp.Kvs[i].Value), err)
			continue
		}

		pt.ResourceVersion = resp.Header.Revision
		ps.Add(pt)
	}
	return ps, nil
}
```

#### 实现Recoder

recorder 其实就是一个 put操作

```go
type recorder struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *recorder) Update(t *pipeline.Pipeline) error {
	if t == nil {
		return fmt.Errorf("update nil pipeline")
	}

	if l.client == nil {
		return fmt.Errorf("etcd client is nil")
	}

	objKey := t.EtcdObjectKey()
	objValue, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline task '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}
```


#### 实现Watcher

watcher的逻辑稍微复杂些:
+ 需要Watch 所有Pipeline对象
+ 处理事件, 把对象缓存到本地的indexer中, 这里Indexer就是一个内存缓存
+ 调用Event处理函数，把事件传播出去


注意Watch是一个阻塞操作, 所以Watcher 是几个里面为一个需要Run的

```go
type shared struct {
	log       logger.Logger
	client    clientv3.Watcher
	indexer   cache.Indexer
	handler   informer.PipelineEventHandler
	filter    informer.PipelineFilterHandler
	watchChan clientv3.WatchChan
}

// AddPipelineEventHandler 添加事件处理回调
func (i *shared) AddPipelineTaskEventHandler(h informer.PipelineEventHandler) {
	i.handler = h
}

// Run 启动 Watch
func (i *shared) Run(ctx context.Context) error {
	// 是否准备完成
	if err := i.isReady(); err != nil {
		return err
	}

	// 监听事件
	i.watch(ctx)

	// 后台处理事件
	go i.dealEvents()
	return nil
}

func (i *shared) dealEvents() {
	// 处理所有事件
	for {
		select {
		case nodeResp := <-i.watchChan:
			for _, event := range nodeResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					if err := i.handlePut(event, nodeResp.Header.GetRevision()); err != nil {
						i.log.Error(err)
					}
				case mvccpb.DELETE:
					if err := i.handleDelete(event); err != nil {
						i.log.Error(err)
					}
				default:
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.handler == nil {
		return errors.New("PipelineEventHandler not add")
	}
	return nil
}

func (i *shared) watch(ctx context.Context) {
	ppWatchKey := pipeline.EtcdPipelinePrefix()
	i.watchChan = i.client.Watch(ctx, ppWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd pipeline resource key: %s", ppWatchKey)
}

func (i *shared) handlePut(event *clientv3.Event, eventVersion int64) error {
	i.log.Debugf("receive pipeline put event, %s", event.Kv.Key)

	// 解析对象
	new, err := pipeline.LoadPipelineFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	new.ResourceVersion = eventVersion

	old, hasOld, err := i.indexer.GetByKey(new.MakeObjectKey())
	if err != nil {
		return err
	}

	if i.filter != nil {
		if err := i.filter(new); err != nil {
			return err
		}
	}

	// 区分Update
	if hasOld {
		// 更新缓存
		i.log.Debugf("update pipeline: %s", new.ShortDescribe())
		if err := i.indexer.Update(new); err != nil {
			i.log.Errorf("update indexer cache error, %s", err)
		}
		i.handler.OnUpdate(old.(*pipeline.Pipeline), new)
	} else {
		// 添加缓存
		i.log.Debugf("add pipeline: %s", new.ShortDescribe())
		if err := i.indexer.Add(new); err != nil {
			i.log.Errorf("add indexer cache error, %s", err)
		}
		i.handler.OnAdd(new)
	}

	return nil
}

func (i *shared) handleDelete(event *clientv3.Event) error {
	key := event.Kv.Key
	i.log.Debugf("receive pipeline delete event, %s", key)

	obj, ok, err := i.indexer.GetByKey(string(key))
	if err != nil {
		i.log.Errorf("get key %s from store error, %s", key)
	}
	if !ok {
		i.log.Warnf("key %s found in store", key)
	}

	// 清除缓存
	if err := i.indexer.Delete(obj); err != nil {
		i.log.Errorf("delete indexer cache error, %s", err)
	}

	i.handler.OnDelete(obj.(*pipeline.Pipeline))
	return nil
}
```


#### 实现Informer

最后我们把这些组合起来我们的Informer就完成了

```go
// NewScheduleInformer todo
func NewInformerr(client *clientv3.Client, filter informer.PipelineFilterHandler) informer.Informer {
	return &Informer{
		log:     zap.L().Named("Pipeline"),
		client:  client,
		filter:  filter,
		indexer: cache.NewIndexer(informer.MetaNamespaceKeyFunc, informer.DefaultStoreIndexers()),
	}
}

// Informer todo
type Informer struct {
	log      logger.Logger
	client   *clientv3.Client
	shared   *shared
	lister   *lister
	recorder *recorder
	indexer  cache.Indexer
	filter   informer.PipelineFilterHandler
}

func (i *Informer) GetStore() cache.Store {
	return i.indexer
}

func (i *Informer) Debug(l logger.Logger) {
	i.log = l
	i.shared.log = l
	i.lister.log = l
}

func (i *Informer) Watcher() informer.Watcher {
	if i.shared != nil {
		return i.shared
	}
	i.shared = &shared{
		log:     i.log.Named("Watcher"),
		client:  clientv3.NewWatcher(i.client),
		indexer: i.indexer,
		filter:  i.filter,
	}
	return i.shared
}

func (i *Informer) Lister() informer.Lister {
	if i.lister != nil {
		return i.lister
	}
	i.lister = &lister{
		log:    i.log.Named("Lister"),
		client: clientv3.NewKV(i.client),
	}
	return i.lister
}

func (i *Informer) Recorder() informer.Recorder {
	if i.recorder != nil {
		return i.recorder
	}
	i.recorder = &recorder{
		log:    i.log.Named("Recorder"),
		client: clientv3.NewKV(i.client),
	}
	return i.recorder
}
```


#### 测试

基于etcdcli可以测试我们Informer功能是否正常:
+ 测试Lister的逻辑
+ 测试Watcher的逻辑
+ 测试Recorder的逻辑

## 服务注册

既然需要调度, 因此我们的Node节点需要注册到中央来，我们才能知道如何调度

### 服务接口

我们先关注:
+ 服务的接口规范
+ 服务的数据结构

我们把workflow的服务抽象成了接口:
```go
type Register interface {
	Debug(logger.Logger)
	Registe() error
	UnRegiste() error
}
```

定义workflow服务的类型:
```go
const (
	// API 提供API访问的服务
	APIType = Type("api")
	// Worker 后台作业服务
	NodeType = Type("node")
	// Scheduler 调度器
	SchedulerType = Type("scheduler")
)
```

Node数据结构 用于泛指一个服务
```go
// Node todo
type Node struct {
	Region          string            `json:"region,omitempty"`
	ResourceVersion int64             `json:"resourceVersion,omitempty"`
	InstanceName    string            `json:"instance_name,omitempty"`
	ServiceName     string            `json:"service_name,omitempty"`
	Type            Type              `json:"type,omitempty"`
	Address         string            `json:"address,omitempty"`
	Version         string            `json:"version,omitempty"`
	GitBranch       string            `json:"git_branch,omitempty"`
	GitCommit       string            `json:"git_commit,omitempty"`
	BuildEnv        string            `json:"build_env,omitempty"`
	BuildAt         string            `json:"build_at,omitempty"`
	Online          int64             `json:"online,omitempty"`
	Tag             map[string]string `json:"tag,omitempty"`

	Prefix   string        `json:"-"`
	Interval time.Duration `json:"-"`
	TTL      int64         `json:"-"`
}
```

### 基于Etcd的注册中心

我们通过etcd来实现注册器

etcd有个租约的概念, 我们可以通过租约来控制一个key 的TTL, 我们基于此来实现注册的心跳功能

+ 往etcd里面写一个服务的key/value，并通过租约设置TTL
+ 每隔一个心跳周期，就KeepAliveOnce 把改租约 续约一次, 也就是心跳机制
+ 最好服务停止时，主动注销服务

1. 初次注册

+ key结构: inforboard/workflow/service/scheduler/{name_name}
+ value: node结构体json数据

```go
func (e *etcd) addOnce() error {
	// 获取leaseID
	resp, err := e.client.Lease.Grant(context.TODO(), e.node.TTL)
	if err != nil {
		return fmt.Errorf("get etcd lease id error, %s", err)
	}
	e.leaseID = resp.ID

	// 写入key
	if _, err := e.client.Put(context.Background(), e.instanceKey, e.instanceValue, clientv3.WithLease(e.leaseID)); err != nil {
		return fmt.Errorf("registe service '%s' with ttl to etcd3 failed: %s", e.instanceKey, err.Error())
	}
	e.instanceKey = e.instanceKey
	return nil
}
```

2. 续约

```go
func (e *etcd) keepAlive(ctx context.Context) {
	// 不停的续约
	interval := e.node.TTL / 5
	e.Infof("keep alive lease interval is %d second", interval)
	tk := time.NewTicker(time.Duration(interval) * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			e.Infof("keepalive goroutine exit")
			return
		case <-tk.C:
			Opctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := e.client.Lease.KeepAliveOnce(Opctx, e.leaseID)
			if err != nil {
				if strings.Contains(err.Error(), "requested lease not found") {
					// 避免程序卡顿造成leaseID失效(比如mac 电脑休眠))
					if err := e.addOnce(); err != nil {
						e.Errorf("refresh registry error, %s", err)
					} else {
						e.Warn("refresh registry success")
					}
				}
				e.Errorf("lease keep alive error, %s", err)
			} else {
				e.Debugf("lease keep alive key: %s", e.instanceKey)
			}
		}
	}
}
```

3. 注销

+ 删除注册上去的服务实例信息
+ 删除租约
+ 停止KeepAlive续约Goroutine

```go
// UnRegiste delete nodeed service from etcd, if etcd is down
// unnode while timeout.
func (e *etcd) UnRegiste() error {
	if e.isStopped {
		return errors.New("the instance has unregisted")
	}
	// delete instance key
	e.stopInstance <- struct{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if resp, err := e.client.Delete(ctx, e.instanceKey); err != nil {
		e.Warnf("unregiste '%s' failed: connect to etcd server timeout, %s", e.instanceKey, err.Error())
	} else {
		if resp.Deleted == 0 {
			e.Infof("unregiste '%s' failed, the key not exist", e.instanceKey)
		} else {
			e.Infof("服务实例(%s)注销成功", e.instanceKey)
		}
	}
	// revoke lease
	_, err := e.client.Lease.Revoke(context.TODO(), e.leaseID)
	if err != nil {
		e.Warnf("revoke lease error, %s", err)
		return err
	}
	e.isStopped = true
	// 停止续约心态
	e.keepAliveStop()
	return nil
}
```

### 服务启动时注册

在服务启动时，调用etcd 注册器来将服务实例信息注册到 注册到etcd

cmd/start.go
```go
// 注册服务
r, err := etcd_register.NewEtcdRegister(svr.node)
if err != nil {
	svr.log.Warn(err)
}
r.Debug(zap.L().Named("Register"))
defer r.UnRegiste()
if err := r.Registe(); err != nil {
	return err
}
```


### 测试

+ 测试正常流程，使用etcdctl 查看 etcd里面改服务的实例是否存在
+ 测试TTL超时，不能完成续约的情况

## Scheduler节点注册

当node节点注册后, 有专门的Node Controller复杂关注Node的变化

### Node加载

这里只需要关注 注册后Node的的加载情况, 因为后面都需要 要访问当前集群有哪些Node都是通过访问Indexer获取的
```go
// controller/node/controller.go

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) run(ctx context.Context, async bool) error {
	// Start the informer factories to begin populating the informer caches
	c.log.Infof("starting node control loop")

	// 调用Lister 获得所有的cronjob 并添加cron
	c.log.Info("starting sync(List) all nodes")
	nodes, err := c.informer.Lister().ListAll(ctx)
	if err != nil {
		return err
	}

	// 更新node存储
	for i := range nodes {
		c.informer.GetStore().Add(nodes[i])
		c.enqueueForAdd(nodes[i])
	}
	c.log.Infof("sync all(%d) nodes success", len(nodes))

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}

	if async {
		go c.waitDown(ctx)
	} else {
		c.waitDown(ctx)
	}
	return nil
}
```

### 测试

我们启动调度器
```sh
make run-sch

...
2021-12-04T12:00:38.628+0800    INFO    [Node]  node/controller.go:82   starting node control loop
2021-12-04T12:00:38.629+0800    INFO    [Node]  node/controller.go:85   starting sync(List) all nodes
2021-12-04T12:00:38.629+0800    INFO    [Node.Lister]   etcd/lister.go:30       list etcd node resource key: inforboard/workflow/service
2021-12-04T12:00:38.631+0800    INFO    [Node.Lister]   etcd/lister.go:47       total nodes: 1
...
```

确保 当前调度器节点已经注册成功, 并且已经加载到node store中
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  inforboard

inforboard/workflow/service/scheduler/DESKTOP-HOVMR7V
{"instance_name":"DESKTOP-HOVMR7V","service_name":"workflow","type":"scheduler","address":"127.0.0.1","online":1638590438613}
```


## 调度算法

调度算法的核心逻辑:
+ 挑选一个Node

我们把如何调度抽象成一个Picker

### Picker接口

这里我们有这种资源需要被调度:
+ Pipeline: 多个调度器中挑选一个来处理Pipeline
+ Step: 多个Node中挑选一个执行Step

下面就是Node的分类:
```go
const (
	// API 提供API访问的服务
	APIType = Type("api")
	// Worker 后台作业服务
	NodeType = Type("node")
	// Scheduler 调度器
	SchedulerType = Type("scheduler")
)

type Type string

type Node struct {
	Region          string            `json:"region,omitempty"`
	ResourceVersion int64             `json:"resourceVersion,omitempty"`
	InstanceName    string            `json:"instance_name,omitempty"`
	ServiceName     string            `json:"service_name,omitempty"`
	Type            Type              `json:"type,omitempty"`
	Address         string            `json:"address,omitempty"`
	Version         string            `json:"version,omitempty"`
	GitBranch       string            `json:"git_branch,omitempty"`
	GitCommit       string            `json:"git_commit,omitempty"`
	BuildEnv        string            `json:"build_env,omitempty"`
	BuildAt         string            `json:"build_at,omitempty"`
	Online          int64             `json:"online,omitempty"`
	Tag             map[string]string `json:"tag,omitempty"`

	Prefix   string        `json:"-"`
	Interval time.Duration `json:"-"`
	TTL      int64         `json:"-"`
}
```

接口定义:
```go
// Picker 挑选一个合适的node 运行Step
type StepPicker interface {
	Pick(*pipeline.Step) (*node.Node, error)
}

type PipelinePicker interface {
	Pick(*pipeline.Pipeline) (*node.Node, error)
}
```

### Roundrobin Picker

Picker就是我们挑选Node的算法, 因此我们先实现一个最简单的算法: RR

step picker:
```go
type roundrobinPicker struct {
	mu    *sync.Mutex
	next  int
	store cache.Store
}

// NewStepPicker 实现分调度
func NewStepPicker(nodestore cache.Store) (algorithm.StepPicker, error) {
	base := &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}
	return &stepPicker{base}, nil
}

type stepPicker struct {
	*roundrobinPicker
}

func (p *stepPicker) Pick(step *pipeline.Step) (*node.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nodes := p.store.List()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	n := nodes[p.next]

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return n.(*node.Node), nil
}
```

pipeline picker:
```go
// NewPipelinePicker 实现分调度
func NewPipelinePicker(nodestore cache.Store) (algorithm.PipelinePicker, error) {
	base := &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}
	return &pipelinePicker{base}, nil
}

type pipelinePicker struct {
	*roundrobinPicker
}

func (p *pipelinePicker) Pick(step *pipeline.Pipeline) (*node.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nodes := p.store.List()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	schs := []*node.Node{}
	for i := range nodes {
		n := nodes[i].(*node.Node)
		if n.Type == node.SchedulerType {
			schs = append(schs, n)
		}
	}

	n := schs[p.next]
	// 修改状态
	p.next = (p.next + 1) % len(schs)

	return n, nil
}
```


## Pipeline Controller

已经准备好了:
+ Informer
+ Node Register
+ Picker

现在Controller负责将他们串起来


### 启动流程

controller 启动逻辑:
+ Watch: 在Controller启动之前， watcher已经启动, Controller通过Watch订阅Pipeline事件，把Watch到的事件添加到工作队列(work queue)
+ List: Controller 首先Sync一下，拉去当前所有Pipeline，判定是否有需要处理的，添加到工作队列(work queue)
+ Worker: 然后启动worker，处理工作队列里面的事件(worker queue)
+ Sgin: 等待Controller 结束

```go
// PipelineTaskScheduler 调度器控制器
type Controller struct {
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue      workqueue.RateLimitingInterface
	informer       informer.Informer
	step           step.Informer
	log            logger.Logger
    // 启动多个worker来处理事件
	workerNums     int
    // 当前running中的worker
	runningWorkers map[string]bool
	wLock          sync.Mutex
    // 调度挑选算法
	picker         algorithm.PipelinePicker
    // 调度器的名称
	schedulerName  string
}
```

下面是运行逻辑:
```go
// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) run(ctx context.Context, async bool) error {
	// Start the informer factories to begin populating the informer caches
	c.log.Infof("starting pipeline control loop, schedule name: %s", c.schedulerName)

	// 获取所有的pipeline
	if err := c.sync(ctx); err != nil {
		return err
	}

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}

	if async {
		go c.waitDown(ctx)
	} else {
		c.waitDown(ctx)
	}

	return nil
}
```

下面是sync的逻辑：
```go

func (c *Controller) sync(ctx context.Context) error {
	// 获取所有的pipeline
	listCount := 0
	ps, err := c.informer.Lister().List(ctx, nil)
	if err != nil {
		return err
	}

	// 看看是否有需要调度的
	for i := range ps.Items {
		p := ps.Items[i]

        // 判定Pipeline是否已经执行完成, 已经完成的Pipeline无效处理
        // 由此可见，我们Etcd里面是不适合存储大量历史数据的
		if p.IsComplete() {
			c.log.Debugf("pipline %s is complete, skip schedule",
				p.ShortDescribe())
			continue
		}

        // 判定改Pipeline是否需要当前调度器处理
        // 这里是多个Controller竞争一个Pipeline调度
		if !p.MatchScheduler(c.schedulerName) {
			c.log.Debugf("pipeline %s scheduler %s is not match this scheduler %s",
				p.ShortDescribe(), p.ScheduledNodeName(), c.schedulerName)
			continue
		}

        // 如果都不是，这添加到工作队列，等待调度
		c.enqueueForAdd(p)
		listCount++
	}
	c.log.Infof("%d pipeline need schedule", listCount)
	return nil
}
```

我们继续看Run worker的逻辑:
```go
// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker(name string) {
	isRunning, ok := c.runningWorkers[name]
	if ok && isRunning {
		c.log.Warnf("worker %s has running", name)
		return
	}
	c.wLock.Lock()
	c.runningWorkers[name] = true
	c.log.Infof("start worker %s", name)
	c.wLock.Unlock()
	for c.processNextWorkItem() {
	}
	if isRunning, ok = c.runningWorkers[name]; ok {
		c.wLock.Lock()
		delete(c.runningWorkers, name)
		c.wLock.Unlock()
		c.log.Infof("worker %s has stopped", name)
	}
}
```

### Worker处理队列数据

worker 启动后会从队列中消费数据, k8s的这个workqueue 在事件处理完成后需要调用Forget才能正在代表消费完了这条数据(和kafka 的commit机制一样)

我们看下worker处理事件的流程:
+ 从队列里面取出一条数据(key)
+ 如果队列里面的数据不是一个string(key是string), 直接丢弃掉, 不处理
+ 然后我们把key传递给handler，后面有handler处理具体的业务
+ handler处理成功，则Forget该事件，该事件成功处理完成
+ handler处理失败, 事件标记为Down, 但是并不会Forget， 等待下次再次处理

```go
// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	c.log.Debugf("get obj from queue: %s", obj)
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			c.log.Errorf("expected string in workqueue but got %#v", obj)
			return nil
		}
		c.log.Debugf("wait sync: %s", key)

		// Run the syncHandler, passing it the namespace/name string of the
		// Network resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		c.log.Infof("successfully synced '%s'", key)
		return nil
	}(obj)
	if err != nil {
		c.log.Error(err)
		return true
	}
	return true
}
```

### Handler面向期望

从worker queue接收到key后, 和 判断缓冲池(informer indexer):
+ 如果有 就 是期望 新增一个对象
	+ 当新增一个对象的时候, 我们就需要调度器来运行这个Pipeline, 这里的运行指的也是调度，因为Pipeline定义的是task的编排, 具体的task不是由 Node节点负责运行的
	+ 运行一个pipeline，流程如下
		+ 判断是否已经完成, 已完不做处理
		+ 判断是否需要调度, 未调度先调度
		+ 判断是否已经运行, 未运行先标记为运行状态
		+ 如果是运行状态, 判断pipline是否需要中断
		+ 如果pipeline正常, 则开始运行定义的 Next Step
+ 如果没有 就是期望 删除一个对象
	+ 删除一个对象, 暂时没啥好处理的

```go
// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Network resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	obj, ok, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return err
	}

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("remove pipeline: %s, skip", key)
		return nil
	}

	ins, isOK := obj.(*pipeline.Pipeline)
	if !isOK {
		return errors.New("invalidate *pipeline.Pipeline obj")
	}

	// 运行pipeline
	if err := c.runPipeline(ins); err != nil {
		return err
	}

	return nil
}
```

运行pipeline流程
```go
// 运行一个pipeline，流程如下
// 1. 判断是否已经完成, 已完不做处理
// 2. 判断是否需要调度, 未调度先调度
// 3. 判断是否已经运行, 未运行先标记为运行状态
// 4. 如果是运行状态, 判断pipline是否需要中断
// 5. 如果pipeline正常, 则开始运行定义的 Next Step
func (c *Controller) runPipeline(p *pipeline.Pipeline) error {
	c.log.Debugf("receive add pipeline: %s status: %s", p.ShortDescribe(), p.Status.Status)
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalidate pipeline error, %s", err)
	}

	// 已经处理完成的无需处理
	if p.IsComplete() {
		return fmt.Errorf("skip run complete pipeline %s, status: %s", p.ShortDescribe(), p.Status.Status)
	}

	// TODO: 使用分布式锁trylock处理 多个实例竞争调度问题

	// 未调度的选进行调度后, 再处理
	if !p.IsScheduled() {
		if err := c.schedulePipeline(p); err != nil {
			return err
		}
		return nil
	}

	// 标记开始执行, 并更新保存
	if !p.IsRunning() {
		p.Run()
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s start status to store error, %s", p.ShortDescribe(), err)
		} else {
			c.log.Debugf("update pipeline %s start status to store success", p.ShortDescribe())
		}
		return nil
	}

	// 判断pipeline没有要执行的下一步, 则结束整个Pipeline
	steps := c.nextStep(p)
	c.log.Debugf("pipeline %s start run next steps: %s", p.ShortDescribe(), steps)
	return c.runPipelineNextStep(steps)
}
```

### 创建Step任务

我们首先要挑选出 第一批需要运行的Step 任务来进行创建, 如何挑选喃?
![](./images/pipeline-flow.png)

+ 如果是并行的任务需要一批同时运行, 我们把这个并行运行的概念定义为一个flow, 一个flow运行完后，才会运行下一个flow
+ 如何判断哪些step是一个flow喃?
	+ flow 必须处于stage内部, 不能夸state出现flow
	+ flow 连续的并行任务别标记为一个flow, 因此最大的并行情况就是 一个stage的任务都是并行的, 整个stage就是一个flow
	+ 由此可见 stage一定是串联的, 这样才是流水线,整体逻辑是串行的, 但是stage内的任务 支持并行



#### 挑选Flow steps

我们看下挑选逻辑
	+ pipeline 会一直
```go
func (c *Controller) nextStep(p *pipeline.Pipeline) []*pipeline.Step {
	//pipeline 下次执行需要的step
	steps, isComplete := p.NextStep()
	if isComplete {
		p.Complete()
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s end status to store error, %s", p.ShortDescribe(), err)
		} else {
			c.log.Debugf("pipeline is complete, update pipeline status to db success")
		}
		return nil
	}
	...
}
```

下面是挑选下一批需要运行的steps步骤(flow概念在内部):
+ 判断当前Flow是否运行完成
	+ 当前flow是否中断, flow中的某个任务执行失败
	+ 判断flow中的任务是否全部完成(是否处于running中), 如果没有完成，需要等待都完成后，才进行下一个flow
+ 获取下个flow
	+ 判断下一个flow 是否为nil, 从而判断pipeline是否完成，无其他step

```go
// 只有上一个flow执行完成后, 才会有下个fow
// 注意: 多个并行的任务是不能跨stage同时执行的
//      也就是说stage一定是串行的
func (p *Pipeline) NextStep() (steps []*Step, isComplete bool) {
	// 判断当前Flow是否运行完成
	if f := p.GetCurrentFlow(); f != nil {
		// 如果有flow中断, pipeline 提前结束
		if f.IsBreak() {
			isComplete = true
			return
		}

		// 如果flow没有pass 说明还是在运行中, 不需要调度下以组step
		if !f.IsPassed() {
			return
		}
	}

	f := p.GetNextFlow()

	// 判断是不是最后一个Flow了
	if f == nil {
		isComplete = true
		return
	}

	// 如果不是则获取flow中的step
	steps = f.items
	return
}
```

#### 同步step状态到pipeline

当 一个flow里面的step 只有一部分完成时，需要把step当前的状态同步到pipeline上 对象上做持久化

```go
func (c *Controller) nextStep(p *pipeline.Pipeline) []*pipeline.Step {
	...

	// 找出需要同步的step
	needSync := []*pipeline.Step{}
	for i := range steps {
		ins := steps[i]

		// 判断step是否已经运行, 如果已经运行则更新Pipeline状态
		old, err := c.step.Lister().Get(context.Background(), ins.Key)
		if err != nil {
			c.log.Errorf("get step %s by key error, %s", ins.Key, err)
			return nil
		}

		if old == nil {
			c.log.Debugf("step %s not found in db", ins.Key)
			continue
		}

		// 状态相等 则无需同步
		if ins.Status.Status.Equal(old.Status.Status) {
			c.log.Debugf("pipeline step status: %s, etcd step status: %s, has sync",
				ins.Status.Status, old.Status.Status)
			continue
		}

		needSync = append(needSync, old)
	}

	// 同步step到pipeline上
	if len(needSync) > 0 {
		for i := range needSync {
			c.log.Debugf("sync step %s to pipeline ...", needSync[i].Key)
			p.UpdateStep(needSync[i])
		}
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline status error, %s", err)
			return nil
		}
		c.log.Debugf("sync %d steps ok", len(needSync))
		return nil
	}

	return steps
}
```

#### 创建Flow的step task

我们需要把挑选出来的step对象 写入etcd, 交给后面的 step调度器进行调度

```go
func (c *Controller) runPipelineNextStep(steps []*pipeline.Step) error {
	// 将需要调度的任务, 交给step调度器调度
	if c.step == nil {
		return fmt.Errorf("step recorder is nil")
	}

	// 有step则进行执行
	for i := range steps {
		ins := steps[i]

		c.log.Debugf("create pipeline step: %s", ins.Key)
		if err := c.step.Recorder().Update(ins.Clone()); err != nil {
			c.log.Errorf(err.Error())
		}
	}

	return nil
}
```

### 测试

到此我们 Pipeline的 控制器就完成了, 可以测试下流程

启动服务:
+ API Server
+ Schduler

```sh
make run-api
make run-sch
```



我们使用API server创建一个pipeline: POST http://{{HOST}}/workflow/api/v1/pipelines/
+ stage1: 2个并行，一个串行
+ stage2: 2个串行

```json
{
    "name": "pipeline01",
    "stages": [
        {
            "name": "stage1",
            "steps": [
                {
                    "name": "step1.1",
                    "action": "go_build@v1",
                    "with": {
                        "ENV1": "env1",
                        "ENV2": "env2"
                    },
                    "with_audit": false,
                    "is_parallel": true,
                    "webhooks": [
                        {
                            "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                            "events": [
                                "SUCCEEDED"
                            ]
                        }
                    ]
                },
                {
                    "name": "step1.2",
                    "action": "go_build@v1",
                    "with_audit": false,
                    "is_parallel": true,
                    "with": {
                        "ENV1": "env1",
                        "ENV2": "env2"
                    },
                    "webhooks": [
                        {
                            "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                            "events": [
                                "SUCCEEDED"
                            ]
                        }
                    ]
                },
                {
                    "name": "step1.3",
                    "action": "go_build@v1",
                    "with_audit": true,
                    "with": {
                        "ENV1": "env1",
                        "ENV2": "env2"
                    },
                    "webhooks": [
                        {
                            "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                            "events": [
                                "SUCCEEDED"
                            ]
                        }
                    ]
                }
            ]
        },
        {
            "name": "stage2",
            "steps": [
                {
                    "name": "step2.1",
                    "action": "go_build@v1",
                    "with": {
                        "ENV1": "env3",
                        "ENV2": "env4"
                    },
                    "webhooks": [
                        {
                            "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                            "events": [
                                "SUCCEEDED"
                            ]
                        }
                    ]
                },
                {
                    "name": "step2.2",
                    "action": "go_build@v1",
                    "with": {
                        "ENV1": "env1",
                        "ENV2": "env2"
                    },
                    "webhooks": [
                        {
                            "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                            "events": [
                                "SUCCEEDED"
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}
```

我们通过API查看当前结果:
+ Pipeline是否调度
+ Piptline step 是否创建(flow1 2个step)
+ Piptline step 是否调度(由于我们还没写node, 也没有node节点注册上来, 所以应该是调度报错)
```json
{
    "code": 0,
    "data": {
        "total": 0,
        "items": [
            {
                "id": "c6lepk93n7pjoq14b8c0",
                "resource_version": 85,
                "domain": "admin-domain",
                "namespace": "c6br8ju1l0cvabpa7fdg",
                "create_at": 1638591697979,
                "create_by": "admin",
                "template_id": "",
                "name": "pipeline01",
                "with": null,
                "mount": null,
                "tags": null,
                "description": "",
                "on": null,
                "hook_event": null,
                "status": {
                    "current_flow": 1,
                    "start_at": 1638591698020,
                    "end_at": 0,
                    "status": "EXECUTING",
                    "scheduler_node": "DESKTOP-HOVMR7V",
                    "message": ""
                },
                "stages": [
                    {
                        "id": 1,
                        "name": "stage1",
                        "needs": null,
                        "steps": [
                            {
                                "key": "c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.1",
                                "create_type": "PIPELINE",
                                "namespace": "c6br8ju1l0cvabpa7fdg",
                                "pipeline_id": "c6lepk93n7pjoq14b8c0",
                                "create_at": 1638591698036,
                                "update_at": 1638591698066,
                                "deploy_id": "",
                                "resource_version": 82,
                                "id": 1,
                                "name": "step1.1",
                                "action": "go_build@v1",
                                "with": {
                                    "ENV1": "env1",
                                    "ENV2": "env2"
                                },
                                "is_parallel": true,
                                "ignore_failed": false,
                                "with_audit": false,
                                "audit_params": null,
                                "with_notify": false,
                                "notify_params": null,
                                "webhooks": [
                                    {
                                        "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                                        "header": null,
                                        "events": [
                                            "SUCCEEDED"
                                        ],
                                        "description": "",
                                        "status": null
                                    }
                                ],
                                "node_selector": null,
                                "status": {
                                    "flow_number": 1,
                                    "start_at": 0,
                                    "end_at": 1638591698066,
                                    "status": "SCHEDULE_FAILED",
                                    "scheduled_node": "",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "has no available node nodes",
                                    "response": null
                                }
                            },
                            {
                                "key": "c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.2",
                                "create_type": "PIPELINE",
                                "namespace": "c6br8ju1l0cvabpa7fdg",
                                "pipeline_id": "c6lepk93n7pjoq14b8c0",
                                "create_at": 1638591698036,
                                "update_at": 1638591698069,
                                "deploy_id": "",
                                "resource_version": 83,
                                "id": 2,
                                "name": "step1.2",
                                "action": "go_build@v1",
                                "with": {
                                    "ENV1": "env1",
                                    "ENV2": "env2"
                                },
                                "is_parallel": true,
                                "ignore_failed": false,
                                "with_audit": false,
                                "audit_params": null,
                                "with_notify": false,
                                "notify_params": null,
                                "webhooks": [
                                    {
                                        "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                                        "header": null,
                                        "events": [
                                            "SUCCEEDED"
                                        ],
                                        "description": "",
                                        "status": null
                                    }
                                ],
                                "node_selector": null,
                                "status": {
                                    "flow_number": 1,
                                    "start_at": 0,
                                    "end_at": 1638591698068,
                                    "status": "SCHEDULE_FAILED",
                                    "scheduled_node": "",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "has no available node nodes",
                                    "response": null
                                }
                            },
                            {
                                "create_type": "PIPELINE",
                                "namespace": "",
                                "create_at": 0,
                                "update_at": 0,
                                "deploy_id": "",
                                "id": 3,
                                "name": "step1.3",
                                "action": "go_build@v1",
                                "with": {
                                    "ENV1": "env1",
                                    "ENV2": "env2"
                                },
                                "is_parallel": false,
                                "ignore_failed": false,
                                "with_audit": true,
                                "audit_params": null,
                                "with_notify": false,
                                "notify_params": null,
                                "webhooks": [
                                    {
                                        "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                                        "header": null,
                                        "events": [
                                            "SUCCEEDED"
                                        ],
                                        "description": "",
                                        "status": null
                                    }
                                ],
                                "node_selector": null,
                                "status": {
                                    "flow_number": 0,
                                    "start_at": 0,
                                    "end_at": 0,
                                    "status": "PENDDING",
                                    "scheduled_node": "",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "",
                                    "response": {}
                                }
                            }
                        ]
                    },
                    {
                        "id": 2,
                        "name": "stage2",
                        "needs": null,
                        "steps": [
                            {
                                "create_type": "PIPELINE",
                                "namespace": "",
                                "create_at": 0,
                                "update_at": 0,
                                "deploy_id": "",
                                "id": 1,
                                "name": "step2.1",
                                "action": "go_build@v1",
                                "with": {
                                    "ENV1": "env3",
                                    "ENV2": "env4"
                                },
                                "is_parallel": false,
                                "ignore_failed": false,
                                "with_audit": false,
                                "audit_params": null,
                                "with_notify": false,
                                "notify_params": null,
                                "webhooks": [
                                    {
                                        "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                                        "header": null,
                                        "events": [
                                            "SUCCEEDED"
                                        ],
                                        "description": "",
                                        "status": null
                                    }
                                ],
                                "node_selector": null,
                                "status": {
                                    "flow_number": 0,
                                    "start_at": 0,
                                    "end_at": 0,
                                    "status": "PENDDING",
                                    "scheduled_node": "",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "",
                                    "response": {}
                                }
                            },
                            {
                                "create_type": "PIPELINE",
                                "namespace": "",
                                "create_at": 0,
                                "update_at": 0,
                                "deploy_id": "",
                                "id": 2,
                                "name": "step2.2",
                                "action": "go_build@v1",
                                "with": {
                                    "ENV1": "env1",
                                    "ENV2": "env2"
                                },
                                "is_parallel": false,
                                "ignore_failed": false,
                                "with_audit": false,
                                "audit_params": null,
                                "with_notify": false,
                                "notify_params": null,
                                "webhooks": [
                                    {
                                        "url": "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479",
                                        "header": null,
                                        "events": [
                                            "SUCCEEDED"
                                        ],
                                        "description": "",
                                        "status": null
                                    }
                                ],
                                "node_selector": null,
                                "status": {
                                    "flow_number": 0,
                                    "start_at": 0,
                                    "end_at": 0,
                                    "status": "PENDDING",
                                    "scheduled_node": "",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "",
                                    "response": {}
                                }
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

我们也可以通过etcd看到当前的2个step task
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  inforboard

inforboard/workflow/steps/c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.1
{"key":"c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.1","create_type":"PIPELINE","namespace":"c6br8ju1l0cvabpa7fdg","pipeline_id":"c6lepk93n7pjoq14b8c0","create_at":1638591698036,"update_at":1638591698066,"deploy_id":"","resource_version":80,"id":1,"name":"step1.1","action":"go_build@v1","with":{"ENV1":"env1","ENV2":"env2"},"is_parallel":true,"ignore_failed":false,"with_audit":false,"audit_params":null,"with_notify":false,"notify_params":null,"webhooks":[{"url":"https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479","header":null,"events":["SUCCEEDED"],"description":"","status":null}],"node_selector":null,"status":{"flow_number":1,"start_at":0,"end_at":1638591698066,"status":"SCHEDULE_FAILED","scheduled_node":"","audit_at":0,"audit_response":"UOD","audit_message":"","notify_at":0,"notify_error":"","message":"has no available node nodes","response":null}}
inforboard/workflow/steps/c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.2
{"key":"c6br8ju1l0cvabpa7fdg.c6lepk93n7pjoq14b8c0.1.2","create_type":"PIPELINE","namespace":"c6br8ju1l0cvabpa7fdg","pipeline_id":"c6lepk93n7pjoq14b8c0","create_at":1638591698036,"update_at":1638591698069,"deploy_id":"","resource_version":81,"id":2,"name":"step1.2","action":"go_build@v1","with":{"ENV1":"env1","ENV2":"env2"},"is_parallel":true,"ignore_failed":false,"with_audit":false,"audit_params":null,"with_notify":false,"notify_params":null,"webhooks":[{"url":"https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479","header":null,"events":["SUCCEEDED"],"description":"","status":null}],"node_selector":null,"status":{"flow_number":1,"start_at":0,"end_at":1638591698068,"status":"SCHEDULE_FAILED","scheduled_node":"","audit_at":0,"audit_response":"UOD","audit_message":"","notify_at":0,"notify_error":"","message":"has no available node nodes","response":null}}
```

到此Pipeline调度器可以正常运行了