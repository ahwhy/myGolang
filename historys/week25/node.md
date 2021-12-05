# Node 开发

Node的核心是 监听Step对象, 当发现是调度给自己的时候，就运行它


## Informer

不同于调度器的Step Informer， Node的Informer只关心 调度给自己的 Step, 因此我们创建了一个Filter Informer

```go
func NewFilterInformer(client *clientv3.Client, filter informer.StepFilterHandler) informer.Informer {
	return &Informer{
		log:     zap.L().Named("Step"),
		client:  client,
		filter:  filter,
		indexer: cache.NewIndexer(informer.MetaNamespaceKeyFunc, informer.DefaultStoreIndexers()),
	}
}
```

在k8s的代码里也经常看到这样的Filter, 比如 Namespace Filter, 底层的逻辑比较简单, 就是对象过滤
```go
type StepFilterHandler func(obj *pipeline.Step) error

func NewNodeFilter(node *node.Node) StepFilterHandler {
	return func(obj *pipeline.Step) error {
		if !node.IsMatch(obj.ScheduledNodeName()) {
			return fmt.Errorf("step %s not match this node [%s], expect [%s]", obj.Key, node.Name(), obj.ScheduledNodeName())
		}
		return nil
	}
}
```

具体的比对逻辑：和当前Node的name进行比较, 来决定是否过滤掉
```go
func (n *Node) IsMatch(nodeName string) bool {
	return n.InstanceName == nodeName
}
```

我们看Node 启动是逻辑, 就携带了Node Filter
```go
func newService(cfg *conf.Config) (*service, error) {
	// Controller 实例
	rn := MakeRegistryNode(cfg)

	// 实例化Informer
	info := si_impl.NewFilterInformer(cfg.Etcd.GetClient(), informer.NewNodeFilter(rn))

	ctl := controller.NewController(rn.Name(), info, client.C())
	ctl.Debug(zap.L().Named("Node"))

	svr := &service{
		info: info,
		log:  zap.L().Named("CLI"),
		node: rn,
		ctl:  ctl,
	}
	return svr, nil
}
```

其他逻辑和scheduler的 Step Informer一样, 毕竟Informer是 共享包, Controller按需加载

## Controller

Controller的主题框架基本是一样的, 只有Sync和 Handler不同, 因此k8s有专门的工具来生成着部分代码

### 加载Step

```go
func (c *Controller) sync(ctx context.Context) error {
	// 调用Lister 获得所有的cronjob 并添加cron
	c.log.Info("starting sync(List) all steps")
	steps, err := c.informer.Lister().List(ctx)
	if err != nil {
		return err
	}

	// 新增所有的job
	for i := range steps {
		c.enqueueForAdd(steps[i])
	}
	c.log.Infof("sync all(%d) steps success", len(steps))

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}

	return nil
}
```

那这里在List的时候 有没有过滤出，只属于本Node节点的Step喃?

下面是Informer Lister的逻辑, 可以看到 如果有filter是要过滤的, 因此我们可以放心加载
```go
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
```

### 处理变更

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
		c.log.Debugf("remove step: %s, skip", key)
	}

	st, isOK := obj.(*pipeline.Step)
	if !isOK {
		return errors.New("invalidate *pipeline.Step obj")
	}

	// 添加step
	if err := c.addStep(st); err != nil {
		return err
	}
	return nil
}
```

### 添加Step任务

如果一个Step需要Run我们专门一个做个Engine模块, 交给他来整体负责

```go
func (c *Controller) addStep(s *pipeline.Step) error {
	status := s.Status.Status
	switch status {
	case pipeline.STEP_STATUS_PENDDING:
		engine.RunStep(s)
		return nil
	case pipeline.STEP_STATUS_RUNNING:
		// TODO: 判断引擎中该step状态是否一致
		// 如果不一致则同步状态, 但是不作再次运行
		c.log.Debugf("step is running, no thing todo")
	case pipeline.STEP_STATUS_CANCELING:
		return c.cancelStep(s)
	case pipeline.STEP_STATUS_SUCCEEDED,
		pipeline.STEP_STATUS_FAILED,
		pipeline.STEP_STATUS_CANCELED,
		pipeline.STEP_STATUS_SKIP,
		pipeline.STEP_STATUS_REFUSE:
		return fmt.Errorf("step %s status is %s has complete", s.Key, status)
	case pipeline.STEP_STATUS_AUDITING:
		return fmt.Errorf("step %s is %s, is auditing", s.Key, status)
	}

	return nil
}

func (c *Controller) cancelStep(s *pipeline.Step) error {
	c.log.Infof("receive cancel object: %s", s)
	if err := s.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return nil
	}

	// 已经完成的step不作处理
	if s.IsComplete() {
		c.log.Debugf("step [%s] is complete, skip cancel", s.Key)
	}

	engine.CancelStep(s)
	return nil
}
```

### 执行引擎

引擎负责Step的具体执行, 不同类型的Step 由引擎交给不同的Runner进行运行

```go
var (
	engine = &Engine{}
)

func RunStep(s *pipeline.Step) {
	// 开始执行, 更新状态
	s.Run()
	engine.updateStep(s)

	// 执行step
	go engine.Run(s)
}

func CancelStep(s *pipeline.Step) {
	engine.CancelStep(s)
}

func Init(wc *client.Client, recorder step.Recorder) (err error) {
	if wc == nil {
		return fmt.Errorf("init runner error, workflow client is nil")
	}

	engine.log = zap.L().Named("Runner.Engine")
	engine.recorder = recorder
	engine.wc = wc
	engine.docker, err = docker.NewRunner()
	engine.k8s = k8s.NewRunner()
	engine.local = local.NewRunner()

	if err != nil {
		return err
	}

	engine.init = true
	return nil
}

type Engine struct {
	recorder step.Recorder
	wc       *client.Client
	docker   runner.Runner
	k8s      runner.Runner
	local    runner.Runner
	init     bool
	log      logger.Logger
}
```

引擎在运行时, 会把需要合并的参数合并后，一并传入runner, runner就个干活的，他本身并没有啥状态

运行Step的核心逻辑:
+ 通过step的定义的Action, 查询出Action对象, 比如 go_build@v1
+ 然后合并参数:
    + action 默认参数
    + pipeline 全局参数
    + pipeline 运行时的参数, 动态生成的(还未实现)
    + step本身参数
+ 最后根据Action类型, 调用不同的runner来运行

```go
// Run 运行Step
// step的参数加载优先级:
//   1. step 本身传人的
//   2. pipeline 运行中产生的
//   3. pipeline 全局传人
//   4. action 默认默认值
func (e *Engine) run(req *runner.RunRequest, resp *runner.RunResponse) {
	if !e.init {
		resp.Failed("engine not init")
		return
	}

	s := req.Step

	e.log.Debugf("start run step: %s status %s", s.Key, s.Status)

	// 1.查询step对应的action定义
	descA := action.NewDescribeActionRequest(s.ActionName(), s.ActionVersion())
	ctx := gcontext.NewGrpcOutCtx()
	actionIns, err := e.wc.Action().DescribeAction(ctx.Context(), descA)
	if err != nil {
		resp.Failed("describe step action error, %s", err)
		return
	}

	// 2.加载Action默认参数
	req.LoadRunParams(actionIns.DefaultRunParam())

	// 3.查询Pipeline, 加载全局参数
	if s.IsCreateByPipeline() {
		descP := pipeline.NewDescribePipelineRequestWithID(s.GetPipelineId())
		descP.Namespace = s.GetNamespace()
		pl, err := e.wc.Pipeline().DescribePipeline(ctx.Context(), descP)
		if err != nil {
			resp.Failed("describe step pipeline error, %s", err)
			return
		}
		req.LoadRunParams(pl.With)
		req.LoadMount(pl.Mount)
	}

	// 4. 加载step传递的参数
	req.LoadRunParams(s.With)

	// 校验run参数合法性
	if err := actionIns.ValidateRunParam(req.RunParams); err != nil {
		resp.Failed(err.Error())
		return
	}

	// 加载Runner运行需要的参数
	req.LoadRunnerParams(actionIns.RunnerParam())

	e.log.Debugf("choice %s runner to run step", actionIns.RunnerType)
	// 3.根据action定义的runner_type, 调用具体的runner
	switch actionIns.RunnerType {
	case action.RUNNER_TYPE_DOCKER:
		e.docker.Run(context.Background(), req, resp)
	case action.RUNNER_TYPE_K8s:
		e.k8s.Run(context.Background(), req, resp)
	case action.RUNNER_TYPE_LOCAL:
		e.local.Run(context.Background(), req, resp)
	default:
		resp.Failed("unknown runner type: %s", actionIns.RunnerType)
		return
	}
}
```

### Docker Runner

这里主要实现Docker Runner, 负责调用Docker 来运行 run step

如何使用SDK操作Docker, 主要参考官方文档

```go
// Docker官方SDK使用说明: https://docs.docker.com/engine/api/sdk/examples/
// Docker官方API使用说明: https://docs.docker.com/engine/api/v1.41/
type Runner struct {
	cli           *client.Client
	log           logger.Logger
	store         store.StoreFactory
	cancelTimeout *time.Duration
}

// ContainerCreate参数说明:  https://docs.docker.com/engine/api/v1.41/#operation/ContainerCreate
// Runner Params:
//   IMAGE_URL: 镜像URL, 比如: docker-build
//   IMAGE_PULL_SECRET: 拉起镜像的凭证
//   IMAGE_PUSH_SECRET: 推送镜像的凭证
// Run Params:
//   IMAGE_VERSION: 镜像版本 比如: v1
//   GIT_SSH_URL: 代码仓库地址, 比如: git@gitee.com:infraboard/keyauth.git
//   IMAGE_PUSH_URL: 代码推送地址
func (r *Runner) Run(ctx context.Context, in *runner.RunRequest, out *runner.RunResponse) {
	r.log.Debugf("docker start run step: %s", in.Step.Key)

	req := newDockerRunRequest(in)
	if err := req.Validate(); err != nil {
		out.Failed("validate docker run request error, %s", err)
		return
	}

	// 创建容器
	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image: req.Image(),
		Env:   req.ContainerEnv(),
		Cmd:   req.ContainerCMD(),
	}, nil, nil, nil, req.ContainerName())
	if err != nil {
		out.Failed("create container error, %s", err)
		return
	}
	// 退出时销毁容器
	defer r.removeContainer(resp.ID)

	// 更新状态
	up := r.store.NewFileUploader(req.Step.Key)
	out.UpdateReponseMap("log_driver", up.DriverName())
	out.UpdateReponseMap("log_path", up.ObjectID())
	out.UpdateReponseMap(CONTAINER_ID_KEY, resp.ID)
	out.UpdateReponseMap(CONTAINER_WARN_KEY, strings.Join(resp.Warnings, ","))
	out.UpdateResponse(in.Step)

    // TODO: 如果镜像不存在, 要提前拉去

	// 启动容器
	err = r.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		// 启动失败则删除容器
		r.removeContainer(resp.ID)
		out.Failed("run container error, %s", err)
		return
	}

	// 记录容器的日志
	logStream, err := r.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		out.Failed("get container log error, %s", err)
		return
	}

	// 上传容器日志
	if err := up.Upload(ctx, logStream); err != nil {
		out.Failed(err.Error())
		return
	}

	// 等待容器退出
	statusCh, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			out.Failed(err.Error())
			return
		}
	case <-statusCh:
	}

	// 容器退出
	if err := r.containerExit(resp.ID); err != nil {
		out.Failed(err.Error())
		return
	}
}
```

### Runner日志

Step在run的时候 会产生任务日志, 所有我们需要定义个 存储模块来记录这些日志, 接口是按照OSS来设计的, 主要是支持stream

```go
// 保存Runner运行中的日志
type StoreFactory interface {
	NewFileUploader(key string) Uploader
}

// 用于上传日志
type Uploader interface {
	DriverName() string
	ObjectID() string
	Upload(ctx context.Context, steam io.ReadCloser) error
}

func NewStore() *Store {
	return &Store{}
}

type Store struct{}

func (s *Store) NewFileUploader(key string) Uploader {
	return file.NewUploader(key)
}
```

默认采用文件模式实现: 就是把日志保存在 runner_log下，按照时间做文件夹，以step命名, 比如

![](./images/runner_log.png)

```go
func NewUploader(id string) *Uploader {
	return &Uploader{
		id:     id,
		root:   "runner_log",
		parent: dateDir(),
	}
}

type Uploader struct {
	id     string
	root   string
	parent string
}

func (u *Uploader) DriverName() string {
	return "local_file"
}
func (u *Uploader) ObjectID() string {
	return path.Join(u.root, u.parent, u.id)
}
func (u *Uploader) Upload(ctx context.Context, stream io.ReadCloser) error {
	defer stream.Close()

	f, err := u.createFile(ctx)
	if err != nil {
		return fmt.Errorf("create file error, %s", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.ReadFrom(stream)
	if err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush file error, %s", err)
	}
	return nil
}

func (u *Uploader) createFile(ctx context.Context) (*os.File, error) {
	fp := u.ObjectID()
	if checkFileIsExist(fp) {
		return os.OpenFile(fp, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	}

	if err := os.MkdirAll(path.Dir(fp), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(fp)
}

func dateDir() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%d/%d/%d", year, int(month), day)
}

// 判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filepath string) bool {
	var exist = true
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
```

## 流水线与CI CD

到此我们基于Pipeline 的 执行引擎就完成了，但是貌似和CI/CD 还没挂上, 接下来我们开发API Server, 让流水线更加SCM的事件来触发
