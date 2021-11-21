# 框架介绍

经过了cmdb的一系列进步，我们的框架也基本形成, 这里再次介绍下 整个框架的思路, 用户中心keyauth也是基于此架构的


## 架构思想

常见的代码组织方式有类:

### 功能式架构

按照功能分层, MVC就是其中典型, cmdb项目就是采用这种功能架构
+ Dao层: 数据对象转换, 与数据库交互
+ Model层: 模型定义, 主要是结构体定义
+ Controller层: 业务控制层, 一些数据处理的Handler都在该层
+ View: 视图层，前后端分离后, 这层一般有前端MVVM替换

![](./pic/mvc.jpeg)


### 分区式架构


按照业务领域进行分区, 需要划分出业务模块的边界, 本项目就是采用这种架构
+ 业务定义层: 定义该业务需要的数据结构和方法
+ 业务实现层: 更加接口定义, 实现具体的业务功能

![](./pic/keyauth-app.png)

## App开发

比如我们按照分区架构来开发一个App, 注意 这里说的App就是 app目录下一个具体的包

### 接口定义

下面是Token业务领域提供的接口定义:

```protobuf
syntax = "proto3";

package infraboard.keyauth.token;
option go_package = "github.com/infraboard/keyauth/app/token";

import "app/token/pb/request.proto";
import "app/token/pb/token.proto";

service TokenService {
    rpc IssueToken(IssueTokenRequest) returns (Token) {};
    rpc ValidateToken (ValidateTokenRequest) returns (Token) {};
	rpc DescribeToken(DescribeTokenRequest) returns (Token) {};
    rpc RevolkToken(RevolkTokenRequest) returns (Token) {};
    rpc BlockToken(BlockTokenRequest) returns (Token) {};
    rpc ChangeNamespace(ChangeNamespaceRequest) returns (Token) {};
	rpc QueryToken(QueryTokenRequest) returns (Set) {};
	rpc DeleteToken(DeleteTokenRequest) returns (DeleteTokenResponse) {};
}
```

生成的接口
```go
// TokenServiceServer is the server API for TokenService service.
// All implementations must embed UnimplementedTokenServiceServer
// for forward compatibility
type TokenServiceServer interface {
	IssueToken(context.Context, *IssueTokenRequest) (*Token, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*Token, error)
	DescribeToken(context.Context, *DescribeTokenRequest) (*Token, error)
	RevolkToken(context.Context, *RevolkTokenRequest) (*Token, error)
	BlockToken(context.Context, *BlockTokenRequest) (*Token, error)
	ChangeNamespace(context.Context, *ChangeNamespaceRequest) (*Token, error)
	QueryToken(context.Context, *QueryTokenRequest) (*Set, error)
	DeleteToken(context.Context, *DeleteTokenRequest) (*DeleteTokenResponse, error)
	mustEmbedUnimplementedTokenServiceServer()
}
```

+ context: context必须是接口的第一个参数
+ 请求参数: 必须把请求封装成一个请求结构体,  比如UpdatePolicy, 就是*UpdatePolicyRequest
+ 响应返回: 返回必须封装为一个响应结构体，具体就是你的业务数据结构


#### 构造函数

尽量不要直接使用结构体, 每个结构体，为其实现一个构造函数, 构造函数命名规范: New+结构体名称

```go
// NewIssueTokenRequest 默认请求
func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{}
}
```

#### 面向接口

接口 就是一个标准规范, 我们的业务是按照标准规范建立起来的. 这就很好的把 标准和实现解耦合了

在接口层定义都是抽象的业务能力, 不能对实体有任何依赖, 简而言之 我们就是在 面向接口编程, 基于接口搭建上层建筑

> 这一层的依赖都是接口，不是实体

下面是一个Token办法服务的定义(接口定义)
```go
// Issuer todo
type Issuer interface {
	CheckClient(ctx context.Context, clientID, clientSecret string) (*application.Application, error)
	IssueToken(context.Context, *token.IssueTokenRequest) (*token.Token, error)
}
```

下面是一个token的颁发对象, 他依赖各种接口, 在接口之上定义上次业务, 上层业务也是接口
```go
// NewTokenIssuer todo
func NewTokenIssuer() (Issuer, error) {
	issuer := &issuer{
		user:    app.GetGrpcApp(user.AppName).(user.UserServiceServer),
		domain:  app.GetGrpcApp(domain.AppName).(domain.DomainServiceServer),
		token:   app.GetGrpcApp(token.AppName).(token.TokenServiceServer),
		ldap:    app.GetInternalApp(provider.AppName).(provider.LDAP),
		app:     app.GetGrpcApp(application.AppName).(application.ApplicationServiceServer),
		emailRE: regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9\.]+)\.([a-zA-Z0-9]+)`),
		log:     zap.L().Named("Token Issuer"),
	}
	return issuer, nil
}

// TokenIssuer 基于该数据进行扩展
type issuer struct {
	app     application.ApplicationServiceServer
	token   token.TokenServiceServer
	user    user.UserServiceServer
	domain  domain.DomainServiceServer
	ldap    provider.LDAP
	emailRE *regexp.Regexp
	log     logger.Logger
}
```

## App实现

在接口定义的当面目录下, 就是接口的具体实现, 实现的包 按照具体实现依赖进行了命名:

+ mysql: 基于MySQL实现
+ etcd: 基于etcd实现
+ impl: 表示这是其中实现


#### App实例

比如下面是一个基于mongo实现的domain服务, 一个domain就是一个组织(Org) 他依赖Mongo的Collection对象和很多其他内部服务

```go
type service struct {
	col           *mongo.Collection
	enableCache   bool
	notifyCachPre string
	domain.UnimplementedDomainServiceServer
}
```


最后依赖mongodb我们实现接口定义的方式，这个服务实力就开发完成了


#### App依赖

上面的domain service是很简单的一个服务，他本身没有其他依赖, 而真实的服务可能需要依赖的很多其他的服务, 比如第三方服务, 内部模块，或者内部系统的其他服务

面向对象: 在面向对象设计的软件系统中，它的底层都是由N个对象构成的，各个对象之间通过相互合作，最终实现系统地业务逻辑

下面直接依赖的结果

![](./pic/object-dep.png)

为了解决对象之间的耦合度过高的问题，软件专家Michael Mattson 1996年提出了IOC理论，用来实现对象之间的“解耦”，目前这个理论已经被成功地应用到实践当中

IOC理论提出的观点大体是这样的：借助于“第三方”实现具有依赖关系的对象之间的解耦:

![](./pic/ioc.png)


IOC是Inversion of Control的缩写，多数书籍翻译成“控制反转”

我们的PKG包就是这个IOC容器, 所有的服务都注册到PKG包下面, 如果需要依赖他，只需要从IOC层取出即可

下面Token依赖的服务

```go
type service struct {
	token.UnimplementedTokenServiceServer
	col           *mongo.Collection
	log           logger.Logger
	enableCache   bool
	notifyCachPre string

	app      application.ApplicationServiceServer
	user     user.UserServiceServer
	domain   domain.DomainServiceServer
	policy   policy.PolicyServiceServer
	issuer   issuer.Issuer
	endpoint endpoint.EndpointServiceServer
	session  session.ServiceServer
	checker  security.Checker
	code     verifycode.VerifyCodeServiceServer
	ns       namespace.NamespaceServiceServer
}
```


#### App注册

下面是token服务的注册

```go
func init() {
	app.RegistryGrpcApp(svr)
}
```

这里的Ioc层已经是一个第三方库了: mcube/app
```go
var (
	grpcApps = map[string]GRPCApp{}
)

// GRPCService GRPC服务的实例
type GRPCApp interface {
	Registry(*grpc.Server)
	Config() error
	Name() string
}

// RegistryService 服务实例注册
func RegistryGrpcApp(app GRPCApp) {
	// 已经注册的服务禁止再次注册
	_, ok := grpcApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("grpc app %s has registed", app.Name()))
	}

	grpcApps[app.Name()] = app
}
```


注册过后这个服务的实例会交给IOC层管理了


### 注册所有App到Ioc

我们通过一个app下的all包 定义需要加载的app, 下面是导入需要加载的grpc app清单

```go
package all

import (
	_ "github.com/infraboard/keyauth/app/application/impl"
	_ "github.com/infraboard/keyauth/app/department/impl"
	_ "github.com/infraboard/keyauth/app/domain/impl"
	_ "github.com/infraboard/keyauth/app/endpoint/impl"
	_ "github.com/infraboard/keyauth/app/mconf/impl"
	_ "github.com/infraboard/keyauth/app/micro/impl"
	_ "github.com/infraboard/keyauth/app/namespace/impl"
	_ "github.com/infraboard/keyauth/app/permission/impl"
	_ "github.com/infraboard/keyauth/app/policy/impl"
	_ "github.com/infraboard/keyauth/app/role/impl"
	_ "github.com/infraboard/keyauth/app/session/impl"
	_ "github.com/infraboard/keyauth/app/tag/impl"
	_ "github.com/infraboard/keyauth/app/token/impl"
	_ "github.com/infraboard/keyauth/app/user/impl"
	_ "github.com/infraboard/keyauth/app/verifycode/impl"
)
```

这样我们导入该包后, 这个服务就通过init导入，完成注册

我们可以在main包，导入后，打印下 当前已经加载的包

## 框架加载Ioc App

服务开发完成后，需要框架加载后才能正常运行

### App初始化

在服务启动的时候，需要初始化所有的App, 这样我们App才能正常工作， IOC提供了服务初始化的方法

```go
// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcApp(server *grpc.Server) error {
	for name, app := range grpcApps {
		err := app.Config()
		if err != nil {
			return fmt.Errorf("config grpc app %s error %s", name, err)
		}

		app.Registry(server)
	}
	return nil
}
```

+ 初始化好服务
+ 加载给Grpc Server暴露


### 运行App

由整个框架负Protocol层责初始化App并运行他们, 在 protocol中 可以知道相关逻辑


#### Grpc 导出

下面是grpc服务在导出的时候，需要提前加载好内部服务和grpc服务
```go
// GRPCService grpc服务
type GRPCService struct {
	svr *grpc.Server
	l   logger.Logger
	c   *conf.Config

	micro micro.MicroServiceServer
}

// Start 启动GRPC服务
func (s *GRPCService) Start() error {
	// 加载内部服务
	if err := app.LoadInternalApp(); err != nil {
		return err
	}

	// 装载所有GRPC服务
	if err := app.LoadGrpcApp(s.svr); err != nil {
		return err
	}

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPCAddr())
	if err != nil {
		return err
	}

	s.l.Infof("GRPC 服务监听地址: %s", s.c.App.GRPCAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		return fmt.Errorf("start service error, %s", err.Error())
	}

	return nil
}
```


#### HTTP导出

如果你的业务模块 需要通过HTTP API 对外暴露, 需要单独编写一个HTTP模块, 用于该服务实例的暴露

```go
var (
	api = &handler{}
)

type handler struct {
	service token.TokenServiceServer
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("token")

	r.BasePath("/oauth2/tokens")
	r.Handle("POST", "/", h.IssueToken).DisableAuth()
	r.Handle("GET", "/", h.ValidateToken)
	r.Handle("DELETE", "/", h.RevolkToken)

	r.BasePath("/self/tokens")
	r.Handle("GET", "/", h.QueryToken)
	r.Handle("POST", "/", h.ChangeNamespace)
	r.Handle("DELETE", "/", h.DeleteToken)
}

func (h *handler) Config() error {
	h.service = app.GetGrpcApp(token.AppName).(token.TokenServiceServer)
	return nil
}

func (h *handler) Name() string {
	return token.AppName
}

func init() {
	app.RegistryHttpApp(api)
}
```

下面是 ioc层里的 http app 注册逻辑
```go
var (
	httpApps = map[string]HTTPApp{}
)

// HTTPService Http服务的实例
type HTTPApp interface {
	Registry(router.SubRouter)
	Config() error
	Name() string
}

// RegistryHttpApp 服务实例注册
func RegistryHttpApp(app HTTPApp) {
	// 已经注册的服务禁止再次注册
	_, ok := httpApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", app.Name()))
	}

	httpApps[app.Name()] = app
}
```

##### http app 开发

我们要暴露那个ioc里面的服务，我就依赖那个, 然后按需暴露, 比如暴露 token service 的issuetoken方法

```go
// IssueToken 颁发资源访问令牌
func (h *handler) IssueToken(w http.ResponseWriter, r *http.Request) {
	req := token.NewIssueTokenRequest()
	req.WithUserAgent(r.UserAgent())
	req.WithRemoteIPFromHTTP(r)

	// 从Header中获取client凭证, 如果有
	req.ClientId, req.ClientSecret, _ = r.BasicAuth()
	req.VerifyCode = r.Header.Get(CodeHeaderKeyName)
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	d, err := h.service.IssueToken(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, d)
	return
}
```

#### Ioc http app注册

在 all下的 http.go 中加载了所有的 http app 到ioc
```go
package all

import (
	// 加载服务模块
	_ "github.com/infraboard/keyauth/app/application/http"
	_ "github.com/infraboard/keyauth/app/department/http"
	_ "github.com/infraboard/keyauth/app/domain/http"
	_ "github.com/infraboard/keyauth/app/endpoint/http"
	_ "github.com/infraboard/keyauth/app/ip2region/http"
	_ "github.com/infraboard/keyauth/app/mconf/http"
	_ "github.com/infraboard/keyauth/app/micro/http"
```


#### 框架加载Ioc http app

在我们http app加载的时候，会注册自己的路由给 root httprouter, 这样我们的子http app就开发完成了

这层逻辑在 protocol的 http.go中

```go
// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	if err := app.LoadHttpApp(s.PathPrefix(), s.r); err != nil {
		return err
	}

	// 启动HTTP服务
	s.l.Infof("HTTP 服务开始启动, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}
```


### 服务启动

最终我们启动我们的http server 和 grpc server 内部的 app 就可以对我提供访问了

```go
func (s *service) start() error {
	if s.bm != nil {
		if err := s.bm.Connect(); err != nil {
			s.log.Errorf("connect bus error, %s", err)
		}
	}

	s.log.Infof("loaded grpc app: %s", app.LoadedGrpcApp())
	s.log.Infof("loaded http app: %s", app.LoadedHttpApp())
	s.log.Infof("loaded internal app: %s", app.LoadedInternalApp())

	go s.grpc.Start()
	return s.http.Start()
}
```

最后总结下我们3中app的开发流程:
+ grpc app
+ internal app
+ http app


## Http Router封装

主要是基于Httprouter定制 认证功能和注册功能, httprouter 本质上就是一个 中间件, 要搞清楚 httprouter如何封装，需要了解中间件的原理


### Web中间件技术

[gin 中间件](https://github.com/gin-gonic/contrib)

```go
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

以及较流行的开源Go语言框架chi：

```
compress.go
  => 对http的响应体进行压缩处理
heartbeat.go
  => 设置一个特殊的路由，例如/ping，/healthcheck，用来给负载均衡一类的前置服务进行探活
logger.go
  => 打印请求处理处理日志，例如请求处理时间，请求路由
profiler.go
  => 挂载pprof需要的路由，如`/pprof`、`/pprof/trace`到系统中
realip.go
  => 从请求头中读取X-Forwarded-For和X-Real-IP，将http.Request中的RemoteAddr修改为得到的RealIP
requestid.go
  => 为本次请求生成单独的requestid，可一路透传，用来生成分布式调用链路，也可用于在日志中串连单次请求的所有逻辑
timeout.go
  => 用context.Timeout设置超时时间，并将其通过http.Request一路透传下去
throttler.go
  => 通过定长大小的channel存储token，并通过这些token对接口进行限流
```

我们基于httprouter封装的路由也是这样的一套框架

#### Hello World

```go
package main

func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe(":8080", nil)
    ...
}
```

现在突然来了一个新的需求，我们想要统计之前写的hello服务的处理耗时，需求很简单，我们对上面的程序进行少量修改
```go
func hello(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    log.Println(timeElapsed)
}
```

这样便可以在每次接收到http请求时，打印出当前请求所消耗的时间

#### 重复的代码

完成了这个需求之后，我们继续进行业务开发，提供的API逐渐增加，现在我们的路由看起来是这个样子

```go
// 省略了一些相同的代码
package main

func helloHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showInfoHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showEmailHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showFriendsHandler(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("your friends is tom and alex"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/info/show", showInfoHandler)
    http.HandleFunc("/email/show", showEmailHandler)
    http.HandleFunc("/friends/show", showFriendsHandler)
    // ...
}
```

每一个handler里都有之前提到的记录运行时间的代码，每次增加新的路由我们也同样需要把这些看起来长得差不多的代码拷贝到我们需要的地方去。因为代码不太多，所以实施起来也没有遇到什么大问题。

渐渐的我们的系统增加到了30个路由和handler函数，每次增加新的handler，我们的第一件工作就是把之前写的所有和业务逻辑无关的周边代码先拷贝过来


#### 代码泥潭

接下来系统安稳地运行了一段时间，突然有一天，老板找到你，我们最近找人新开发了监控系统，为了系统运行可以更加可控，需要把每个接口运行的耗时数据主动上报到我们的监控系统里。给监控系统起个名字吧，叫metrics。现在你需要修改代码并把耗时通过HTTP Post的方式发给metrics系统了

```go
func helloHandler(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
    // 新增耗时上报
    metrics.Upload("timeHandler", timeElapsed)
}
```

修改到这里，本能地发现我们的开发工作开始陷入了泥潭


#### 使用中间件剥离非业务逻辑

我们来分析一下，一开始在哪里做错了呢？我们只是一步一步地满足需求，把我们需要的逻辑按照流程写下去呀？

我们犯的最大的错误，是把业务代码和非业务代码揉在了一起。对于大多数的场景来讲，非业务的需求都是在http请求处理前做一些事情，并且在响应完成之后做一些事情。

我们有没有办法使用一些重构思路把这些公共的非业务功能代码剥离出去呢？

写一个Wrapper函数, 返回一个符合之前函数签名的函数, 也就是我们常说的 装饰器(修饰器/面向切面)编程模式

```go
func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}

func timeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
        timeStart := time.Now()

        // next handler
        next.ServeHTTP(wr, r)

        timeElapsed := time.Since(timeStart)
        logger.Println(timeElapsed)
    })
}

func main() {
	// http.HandleFunc("/", hello)
	// HandlerFunc 是一个类型, 我们把hello这个函数，转换为了HandlerFunc类型, 这我们使用int(a)是一个语法 为啥要这样喃?
	// 因为 HandlerFunc 实现了ServeHTTP方法, 这样我们的hello函数对象也就有了ServeHTTP方法
	// 是不是很骚, HandlerFunc 是个函数， 我们把他定义为一个Type, 然后给这个Type 绑定了一个函数
	// 这样凡事 被转换为HandlerFunc的类型，都是一个http.Handler
	// 当然 要完成 Type()的转换，函数签名必须一致
	http.Handle("/", timeMiddleware(http.HandlerFunc(hello)))
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
```

这样就非常轻松地实现了业务与非业务之间的剥离，魔法就在于这个timeMiddleware。可以从代码中看到，我们的timeMiddleware()也是一个函数，其参数为http.Handler，http.Handler的定义在net/http包中

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

任何方法实现了ServeHTTP，即是一个合法的http.Handler，读到这里你可能会有一些混乱，我们先来梳理一下http库的Handler，HandlerFunc和ServeHTTP的关系
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

只要你的handler函数签名是
```go
func (ResponseWriter, *Request)
```


#### 更多中间件

如果我们有更多的中间件喃?

比如再加一个限速的中间件:
```go
customizedHandler = logger(timeout(ratelimit(helloHandler)))
```

![](./pic/middleware_flow.png)

再直白一些，这个流程在进行请求处理的时候就是不断地进行函数压栈再出栈，有一些类似于递归的执行流
```
[exec of logger logic]           函数栈: []

[exec of timeout logic]          函数栈: [logger]

[exec of ratelimit logic]        函数栈: [timeout/logger]

[exec of helloHandler logic]     函数栈: [ratelimit/timeout/logger]

[exec of ratelimit logic part2]  函数栈: [timeout/logger]

[exec of timeout logic part2]    函数栈: [logger]

[exec of logger logic part2]     函数栈: []
```

#### 更优雅的中间件写法

上一节中解决了业务功能代码和非业务功能代码的解耦，但也提到了，看起来并不美观，如果需要修改这些函数的顺序，或者增删中间件还是有点费劲，本节我们来进行一些“写法”上的优化。

思路很简单: 递归变迭代, 我们把这些函数放到一个[]middleware里面, 进行迭代调用

```go
package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []Middleware
    // mux map[string] http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(m Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Merge(h http.Handler) http.Handler {
	var mergedHandler = h

	// customizedHandler = logger(timeout(ratelimit(helloHandler)))
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

    // r.mux[route] = mergedHandler
	return mergedHandler
}

func TimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		// next handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		log.Println(timeElapsed)
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		log.Println("start")

		// next handler
		next.ServeHTTP(wr, r)

		log.Println("end")
	})
}
```

然后我们就可以进行调用了: 
```go
func main() {
	r := middleware.NewRouter()
	r.Use(middleware.LogMiddleware)
	r.Use(middleware.TimeMiddleware)
	http.Handle("/", r.Merge(http.HandlerFunc(hello)))
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
```

### 路由创建

通过Use使用中间件: 中间件时链式调用，是栈结构, 有先后顺序

+ Recovery: Hold住所有的Panic
+ AccessLog: 记录Access Log
+ cors: 允许跨域
+ EnableAPIRoot: / 可以访问注册的 Endpoint，及API ROOT
+ SetAuther: 添加认证拦截器
+ Auth: 全局变量，后面每个endpoinit可单独覆盖

```go
r := httprouter.New()
r.Use(recovery.NewWithLogger(zap.L().Named("Recovery")))
r.Use(accesslog.NewWithLogger(zap.L().Named("AccessLog")))
r.Use(cors.AllowAll())
r.EnableAPIRoot()
r.SetAuther(pkg.NewAuther())
r.Auth(true)
```

### 路由配置

r.Handle 提供路由处理, path语法和httprouter语法一样， 只是把路径参数封装到了ctx里面了

```go
// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	self := router.ResourceRouter("application")
	self.BasePath("applications")
	self.Handle("POST", "/", h.CreateApplication)
	self.Handle("GET", "/", h.QueryApplication)
	self.Handle("GET", "/:id", h.GetApplication)
	self.Handle("DELETE", "/:id", h.DestroyApplication)
}
```

### 请求上下文

context包封装了请求上下午, auther中间件认证完后的结果也存储在当前ctx中, 通过GetContext获取
```go
func (h *handler) QueryApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	page := request.NewPageRequestFromHTTP(r)
	req := application.NewQueryApplicationRequest(page)
	req.Account = tk.Account

	apps, err := h.service.QueryApplication(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, apps)
}
```

req 请求对象
```go
// ReqContext context
type ReqContext struct {
        Entry    *router.Entry
        PS       httprouter.Params
        AuthInfo interface{}
}
```

+ entry 当前服务条目
+ ps: path参数
+ AuthInfo: 认证完成后的数据

### 认证中间件

那这个 上下文是如何设置到每个请求当中的喃? 答案是 认证中间件

我们在router上 留了 auther设置的口子: r.SetAuther(pkg.NewAuther())
```go
// Auther 设置受保护路由使用的认证器
// Header 用于鉴定身份
// Entry 用于鉴定权限
type Auther interface {
	Auth(req *http.Request, entry httppb.Entry) (authInfo interface{}, err error)
	ResponseHook(http.ResponseWriter, *http.Request, httppb.Entry)
}
```

该中间件会获取token，并根据entry做权限判定, 在后面的keyauth部分会讲 认证中间件, 这里只需要知道留有口子

## 程序配置

全局实例: conf.C()

```go
func (s *service) Config() error {
	db := conf.C().Mongo.GetDB()
	ac := db.Collection("application")
    ...
}
```

## 程序日志

全局实例: zap.L()

```go
func (s *service) Config() error {
    ...
	s.log = zap.L().Named("Department")
	return nil
}
```