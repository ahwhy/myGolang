# Step Scheduler 开发

Step调度器的核心逻辑:
+ Watch Step对象, 监听变化事件
+ 当有新的Step对象被创建时, 修改Step对象的 Scheduler属性, 为其挑选一个可以Node Node来处理Step调度


## Informer

### 定义接口

```go
// Informer 负责事件通知
type Informer interface {
	Watcher() Watcher
	Lister() Lister
	Recorder() Recorder
	GetStore() cache.Store
}
```

#### watcher: 监听对象
```go
// Watcher 负责事件通知
type Watcher interface {
	// Run starts and runs the shared informer, returning after it stops.
	// The informer will be stopped when stopCh is closed.
	Run(ctx context.Context) error
	// AddEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddStepEventHandler(handler StepEventHandler)
}

// StepEventHandler can handle notifications for events that happen to a
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
type StepEventHandler interface {
	OnAdd(obj *pipeline.Step)
	OnUpdate(old, new *pipeline.Step)
	OnDelete(obj *pipeline.Step)
}

// StepEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.
type StepEventHandlerFuncs struct {
	AddFunc    func(obj *pipeline.Step)
	UpdateFunc func(oldObj, newObj *pipeline.Step)
	DeleteFunc func(obj *pipeline.Step)
}
```

#### lister: 查询对象
```go
type Lister interface {
	Get(ctx context.Context, key string) (*pipeline.Step, error)
	List(ctx context.Context) ([]*pipeline.Step, error)
}
```

#### recorder: 更新对象
```go
type Recorder interface {
	Update(*pipeline.Step) error
}
```

#### store: 缓存存储
```go
type Store interface {
	Reader
	Writer
	Manager
}
```

### Informer实现

同理基于etcd实现

#### watcher实现

第一步是watch key: 这里watch的key 为: inforboard/workflow/steps
```go
func (i *shared) watch(ctx context.Context) {
	// 监听事件
	stepWatchKey := pipeline.EtcdStepPrefix()
	i.watchChan = i.client.Watch(ctx, stepWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd step resource key: %s", stepWatchKey)
}

```

然后处理事件:
```go
func (i *shared) dealEvent() {
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
```

处理PUT事件
```go
func (i *shared) handlePut(event *clientv3.Event, eventVersion int64) error {
	i.log.Debugf("receive step put event, %s", event.Kv.Key)

	// 解析对象
	new, err := pipeline.LoadStepFromBytes(event.Kv.Value)
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
		i.log.Debugf("update step store key: %s, status %s", new.Key, new.Status)
		if err := i.indexer.Update(new); err != nil {
			i.log.Errorf("update indexer cache error, %s", err)
		}
		i.handler.OnUpdate(old.(*pipeline.Step), new)
	} else {
		// 添加缓存
		i.log.Debugf("add step store key: %s, status %s", new.Key, new.Status)
		if err := i.indexer.Add(new); err != nil {
			i.log.Errorf("add indexer cache error, %s", err)
		}
		i.handler.OnAdd(new)
	}

	return nil
}
```

处理删除事件
```go
func (i *shared) handleDelete(event *clientv3.Event) error {
	key := event.Kv.Key
	i.log.Debugf("receive step delete event, %s", key)

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

	i.handler.OnDelete(obj.(*pipeline.Step))
	return nil
}
```

#### Lister实现

使用etcd 完成 GET操作: inforboard/workflow/steps

```go
type lister struct {
	log    logger.Logger
	client clientv3.KV
	filter step.StepFilterHandler
}

func (l *lister) List(ctx context.Context) (ret []*pipeline.Step, err error) {
	listKey := pipeline.EtcdStepPrefix()

	l.log.Infof("list etcd step resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	set := pipeline.NewStepSet()
	for i := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadStepFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}

		if l.filter != nil {
			if err := l.filter(ins); err != nil {
				l.log.Error(err)
				continue
			}
		}

		ins.ResourceVersion = resp.Header.Revision
		set.Add(ins)
	}

	return set.Items, nil
}

func (l *lister) Get(ctx context.Context, key string) (*pipeline.Step, error) {
	descKey := pipeline.StepObjectKey(key)
	l.log.Infof("describe etcd step resource key: %s", descKey)
	resp, err := l.client.Get(ctx, descKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	if resp.Count > 1 {
		return nil, exception.NewInternalServerError("step find more than one: %d", resp.Count)
	}

	ins := pipeline.NewDefaultStep()
	for index := range resp.Kvs {
		// 解析对象
		ins, err = pipeline.LoadStepFromBytes(resp.Kvs[index].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}
	}
	return ins, nil
}
```

#### recorder实现

使用etcd 完成PUT操作: inforboard/workflow/steps/{namespace}/{step_id}

```go
type recorder struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *recorder) Update(step *pipeline.Step) error {
	step.UpdateAt = ftime.Now().Timestamp()
	objKey := pipeline.StepObjectKey(step.Key)
	objValue, err := json.Marshal(step)
	if err != nil {
		return err
	}

	l.log.Debugf("update step %s status %s %s ...", objKey, step.Status, string(objValue))
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline step '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}
```

#### Informer 

最后Informer就是整合他们
```go
func NewFilterInformer(client *clientv3.Client, filter informer.StepFilterHandler) informer.Informer {
	return &Informer{
		log:     zap.L().Named("Step"),
		client:  client,
		filter:  filter,
		indexer: cache.NewIndexer(informer.MetaNamespaceKeyFunc, informer.DefaultStoreIndexers()),
	}
}

// NewSInformer todo
func NewInformer(client *clientv3.Client) informer.Informer {
	return NewFilterInformer(client, nil)
}

// Informer todo
type Informer struct {
	log      logger.Logger
	client   *clientv3.Client
	shared   *shared
	lister   *lister
	recorder *recorder
	indexer  cache.Indexer
	filter   informer.StepFilterHandler
}
```

## 调度算法

Step调度算法：
+ 调度池 是 cache.Store，里面存储的是当前注册上的所有Node(需要筛选出Node，应该还有Schuster)
+ 轮询算法: 采用取模方式, 轮流调度(需要记录当前调度的值)

```go
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

	ns := []*node.Node{}
	for i := range nodes {
		n := nodes[i].(*node.Node)
		if n.Type == node.NodeType {
			ns = append(ns, n)
		}
	}

	if len(ns) == 0 {
		return nil, fmt.Errorf("has no available node nodes")
	}

	n := ns[p.next]
	// 修改状态
	p.next = (p.next + 1) % len(ns)

	return n, nil
}
```

## Controller

其他逻辑和Pipeline调度器一样, 不同之处就在于现状处理的对象换成了Step

### 加载Step

List 所有Step, 加载到indexer中
```go
func (c *Controller) sync(ctx context.Context) error {
	// 获取所有的pipeline
	listCount := 0
	steps, err := c.informer.Lister().List(ctx)
	if err != nil {
		return err
	}

	// 看看是否有需要调度的
	for i := range steps {
		s := steps[i]
		if s.IsComplete() {
			c.log.Debugf("step %s is complete, skip schedule", s.Key)
			continue
		}

		if s.IsScheduled() {
			c.log.Debugf("step %s is scheduler to %s, skip schedule", s.Key, s.ScheduledNodeName())
			continue
		}
		c.informer.GetStore().Add(s)
		c.enqueueForAdd(s)
		listCount++
	}
	c.log.Infof("%d step need schedule", listCount)
	return nil
}
```

### 处理事件

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
		c.log.Debugf("removed step: %s, skip", key)
		return nil
	}

	st, isOK := obj.(*pipeline.Step)
	if !isOK {
		return fmt.Errorf("object %T invalidate, is not *pipeline.Step obj, ", obj)
	}

	// 添加
	if err := c.addStep(st); err != nil {
		return err
	}

	return nil
}
```

### 调度Step

+ 已经调度的任务不处理
+ 调度失败的任务不处理 (后面可以考虑重新调度)
+ 如果开启审核，需要通过后，才能调度执行
+ 最后执行调度(设置Step的调度节点属性)

```go
func (c *Controller) addStep(s *pipeline.Step) error {
	c.log.Infof("receive add step: %s", s)
	if err := s.Validate(); err != nil {
		return fmt.Errorf("invalidate node error, %s", err)
	}

	// 已经调度的任务不处理
	if s.IsScheduled() {
		return fmt.Errorf("step %s has schedule to node %s, skip add", s.Key, s.ScheduledNodeName())
	}

	// 调度失败的任务不处理
	if s.IsScheduledFailed() {
		return fmt.Errorf("step %s schedule failed, skip add", s.Key)
	}

	// 如果开启审核，需要通过后，才能调度执行
	if !c.isAllow(s) {
		return fmt.Errorf("step not allow")
	}

	if err := c.scheduleStep(s); err != nil {
		return err
	}

	return nil
}
```

标准Step的调度节点属性: ScheduleNode

```go
// Step任务调度
func (c *Controller) scheduleStep(step *pipeline.Step) error {
	node, err := c.picker.Pick(step)
	if err != nil || node == nil {
		c.log.Warnf("step %s pick node error, %s", step.Name, err)
		step.ScheduleFailed(err.Error())
		// 清除一下其他数据
		if err := c.informer.Recorder().Update(step.Clone()); err != nil {
			c.log.Errorf("update scheduled step error, %s", err)
		}
		return err
	}

	c.log.Debugf("choice [%s] %s for step %s", node.Type, node.InstanceName, step.Key)
	step.SetScheduleNode(node.InstanceName)
	// 清除一下其他数据
	if err := c.informer.Recorder().Update(step.Clone()); err != nil {
		c.log.Errorf("update scheduled step error, %s", err)
	}
	return nil
}
```

## 测试

启动如下服务:
+ API Server
+ Scheduler

```sh
make run-sch
```

启动服务后, 我们注册的Node就只有当前调度器实例, 并没有处理Step的Node被注册上来, 因此任务依然无法调度

### Node 节点注册

```sh
make run-node
```

我们通过etcd可以看到已经注册的服务:
```sh
$ docker exec -it -e "ETCDCTL_API=3" etcd  etcdctl get --prefix  inforboard/workflow/service

inforboard/workflow/service/node/DESKTOP-HOVMR7V
{"instance_name":"DESKTOP-HOVMR7V","service_name":"workflow","type":"node","address":"127.0.0.1","online":1638598176543}
inforboard/workflow/service/scheduler/DESKTOP-HOVMR7V
{"instance_name":"DESKTOP-HOVMR7V","service_name":"workflow","type":"scheduler","address":"127.0.0.1","online":16385980107
```

现在我们暂时不处理 调度失败的情况, 调度失败有 用户通过API 进行对Step进行手动重试

### 新增Pipeline调试

我们把之前的Pipeline删除后, 重新创建一个新的pipeline, 然后查看结果:
+ 我们看到调度逻辑是否正常
```json
{
    "code": 0,
    "data": {
        "total": 0,
        "items": [
            {
                "id": "c6lhgn13n7pjoq14b8hg",
                "resource_version": 124,
                "domain": "admin-domain",
                "namespace": "c6br8ju1l0cvabpa7fdg",
                "create_at": 1638602844653,
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
                    "start_at": 1638602844696,
                    "end_at": 1638602845868,
                    "status": "COMPLETE",
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
                                "key": "c6br8ju1l0cvabpa7fdg.c6lhgn13n7pjoq14b8hg.1.1",
                                "create_type": "PIPELINE",
                                "namespace": "c6br8ju1l0cvabpa7fdg",
                                "pipeline_id": "c6lhgn13n7pjoq14b8hg",
                                "create_at": 1638602844711,
                                "update_at": 1638602845841,
                                "deploy_id": "",
                                "resource_version": 120,
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
                                    "start_at": 1638602844756,
                                    "end_at": 1638602845840,
                                    "status": "FAILED",
                                    "scheduled_node": "DESKTOP-HOVMR7V",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "create container error, Error response from daemon: No such image: busybox:latest",
                                    "response": null
                                }
                            },
                            {
                                "key": "c6br8ju1l0cvabpa7fdg.c6lhgn13n7pjoq14b8hg.1.2",
                                "create_type": "PIPELINE",
                                "namespace": "c6br8ju1l0cvabpa7fdg",
                                "pipeline_id": "c6lhgn13n7pjoq14b8hg",
                                "create_at": 1638602844711,
                                "update_at": 1638602846805,
                                "deploy_id": "",
                                "resource_version": 123,
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
                                    "start_at": 1638602844761,
                                    "end_at": 1638602846805,
                                    "status": "FAILED",
                                    "scheduled_node": "DESKTOP-HOVMR7V",
                                    "audit_at": 0,
                                    "audit_response": "UOD",
                                    "audit_message": "",
                                    "notify_at": 0,
                                    "notify_error": "",
                                    "message": "create container error, Error response from daemon: No such image: busybox:latest",
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

这里我们Node节点在执行Docker run的时候并没有主动去获取镜像, 我们可以先手动拉去镜像到本地, 后面修改Node添加自动拉去镜像的逻辑
```sh
$ docker pull busybox:latest
latest: Pulling from library/busybox
3aab638df1a9: Pull complete
Digest: sha256:52817dece4cfe26f581c834d27a8e1bcc82194f914afe6d50afad5a101234ef1
Status: Downloaded newer image for busybox:latest
docker.io/library/busybox:latest
```


### 调试审核逻辑

再次删除重建Pipeline, 然后流程会卡在 审核处等待, 最后我们调用 审核API 来通过审核: http://{{HOST}}/workflow/api/v1/steps/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.3/audit

```json
{
    "audit_reponse": "ALLOW",
    "audit_message": "good job"
}
```

最后 等待我们Step运行完成
```json
{
    "code": 0,
    "data": {
        "id": "c6lhm1h3n7pjoq14b8l0",
        "resource_version": 169,
        "domain": "admin-domain",
        "namespace": "c6br8ju1l0cvabpa7fdg",
        "create_at": 1638603526474,
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
            "current_flow": 4,
            "start_at": 1638603526514,
            "end_at": 1638603662242,
            "status": "COMPLETE",
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
                        "key": "c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.1",
                        "create_type": "PIPELINE",
                        "namespace": "c6br8ju1l0cvabpa7fdg",
                        "pipeline_id": "c6lhm1h3n7pjoq14b8l0",
                        "create_at": 1638603526529,
                        "update_at": 1638603537677,
                        "deploy_id": "",
                        "resource_version": 140,
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
                                "status": {
                                    "start_at": 1638603537684,
                                    "cost": 457,
                                    "success": true,
                                    "message": "{\"Extra\":null,\"StatusCode\":0,\"StatusMessage\":\"success\"}"
                                }
                            }
                        ],
                        "node_selector": null,
                        "status": {
                            "flow_number": 1,
                            "start_at": 1638603526575,
                            "end_at": 1638603537677,
                            "status": "SUCCEEDED",
                            "scheduled_node": "DESKTOP-HOVMR7V",
                            "audit_at": 0,
                            "audit_response": "UOD",
                            "audit_message": "",
                            "notify_at": 0,
                            "notify_error": "",
                            "message": "",
                            "response": {
                                "container_id": "a37f0bb3b5111f1ccdee6aaacbfe2a6ee1b2dfcb091cf60187bc44e8c8b8c496",
                                "container_warn": "",
                                "log_driver": "local_file",
                                "log_path": "runner_log/2021/12/4/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.1"
                            }
                        }
                    },
                    {
                        "key": "c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.2",
                        "create_type": "PIPELINE",
                        "namespace": "c6br8ju1l0cvabpa7fdg",
                        "pipeline_id": "c6lhm1h3n7pjoq14b8l0",
                        "create_at": 1638603526529,
                        "update_at": 1638603537850,
                        "deploy_id": "",
                        "resource_version": 141,
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
                                "status": {
                                    "start_at": 1638603538146,
                                    "cost": 244,
                                    "success": true,
                                    "message": "{\"Extra\":null,\"StatusCode\":0,\"StatusMessage\":\"success\"}"
                                }
                            }
                        ],
                        "node_selector": null,
                        "status": {
                            "flow_number": 1,
                            "start_at": 1638603526576,
                            "end_at": 1638603537849,
                            "status": "SUCCEEDED",
                            "scheduled_node": "DESKTOP-HOVMR7V",
                            "audit_at": 0,
                            "audit_response": "UOD",
                            "audit_message": "",
                            "notify_at": 0,
                            "notify_error": "",
                            "message": "",
                            "response": {
                                "container_id": "19a300308786206706c5cffb19f5c5b280680f1dbd1707baae9259ec8c110e37",
                                "container_warn": "",
                                "log_driver": "local_file",
                                "log_path": "runner_log/2021/12/4/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.2"
                            }
                        }
                    },
                    {
                        "key": "c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.3",
                        "create_type": "PIPELINE",
                        "namespace": "c6br8ju1l0cvabpa7fdg",
                        "pipeline_id": "c6lhm1h3n7pjoq14b8l0",
                        "create_at": 1638603538406,
                        "update_at": 1638603640141,
                        "deploy_id": "",
                        "resource_version": 153,
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
                                "status": {
                                    "start_at": 1638603640150,
                                    "cost": 277,
                                    "success": true,
                                    "message": "{\"Extra\":null,\"StatusCode\":0,\"StatusMessage\":\"success\"}"
                                }
                            }
                        ],
                        "node_selector": null,
                        "status": {
                            "flow_number": 2,
                            "start_at": 1638603627342,
                            "end_at": 1638603640140,
                            "status": "SUCCEEDED",
                            "scheduled_node": "DESKTOP-HOVMR7V",
                            "audit_at": 1638603627289,
                            "audit_response": "ALLOW",
                            "audit_message": "good job",
                            "notify_at": 0,
                            "notify_error": "",
                            "message": "",
                            "response": {
                                "container_id": "8ba414b587808d82c736b1c8e3234368108ba911980343c520fe2f46fb34d118",
                                "container_warn": "",
                                "log_driver": "local_file",
                                "log_path": "runner_log/2021/12/4/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.1.3"
                            },
                            "context_map": {
                                "AUDIT_NOTIFY_HAS_SEND": "true"
                            }
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
                        "key": "c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.2.1",
                        "create_type": "PIPELINE",
                        "namespace": "c6br8ju1l0cvabpa7fdg",
                        "pipeline_id": "c6lhm1h3n7pjoq14b8l0",
                        "create_at": 1638603640450,
                        "update_at": 1638603651137,
                        "deploy_id": "",
                        "resource_version": 160,
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
                                "status": {
                                    "start_at": 1638603651143,
                                    "cost": 197,
                                    "success": true,
                                    "message": "{\"Extra\":null,\"StatusCode\":0,\"StatusMessage\":\"success\"}"
                                }
                            }
                        ],
                        "node_selector": null,
                        "status": {
                            "flow_number": 3,
                            "start_at": 1638603640480,
                            "end_at": 1638603651137,
                            "status": "SUCCEEDED",
                            "scheduled_node": "DESKTOP-HOVMR7V",
                            "audit_at": 0,
                            "audit_response": "UOD",
                            "audit_message": "",
                            "notify_at": 0,
                            "notify_error": "",
                            "message": "",
                            "response": {
                                "container_id": "32a66ef58f366b8da68816d4a1ada3c5e4f0aecfd61dbc88e3f6072f166cd088",
                                "container_warn": "",
                                "log_driver": "local_file",
                                "log_path": "runner_log/2021/12/4/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.2.1"
                            }
                        }
                    },
                    {
                        "key": "c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.2.2",
                        "create_type": "PIPELINE",
                        "namespace": "c6br8ju1l0cvabpa7fdg",
                        "pipeline_id": "c6lhm1h3n7pjoq14b8l0",
                        "create_at": 1638603651362,
                        "update_at": 1638603662065,
                        "deploy_id": "",
                        "resource_version": 167,
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
                                "status": {
                                    "start_at": 1638603662072,
                                    "cost": 151,
                                    "success": true,
                                    "message": "{\"Extra\":null,\"StatusCode\":0,\"StatusMessage\":\"success\"}"
                                }
                            }
                        ],
                        "node_selector": null,
                        "status": {
                            "flow_number": 4,
                            "start_at": 1638603651394,
                            "end_at": 1638603662064,
                            "status": "SUCCEEDED",
                            "scheduled_node": "DESKTOP-HOVMR7V",
                            "audit_at": 0,
                            "audit_response": "UOD",
                            "audit_message": "",
                            "notify_at": 0,
                            "notify_error": "",
                            "message": "",
                            "response": {
                                "container_id": "52e894b9754ba3a252f1fd63882a5995992c206d0d8f73dd263a82765ec3f389",
                                "container_warn": "",
                                "log_driver": "local_file",
                                "log_path": "runner_log/2021/12/4/c6br8ju1l0cvabpa7fdg.c6lhm1h3n7pjoq14b8l0.2.2"
                            }
                        }
                    }
                ]
            }
        ]
    }
}
```

最后查看飞书群，看看看看是否通知(信息还没格式化):

![](./images/feishu_notify.jpg)

我们在下个小节讲 Pipeline Hook机制

## 总结

不要纠结 Node是如何执行的, 我们先关注调度逻辑, Node如何执行，在Node开发部分会介绍