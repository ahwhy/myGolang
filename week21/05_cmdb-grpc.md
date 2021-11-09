# CMDB改造

+ gRPC改造
+ http改造

## gRPC改造

cmdb项目作为 devcloud的一个子服务, 需要提供内部调用接口

### 自定义Struct Tag

安装插件: 
```sh
go install github.com/favadi/protoc-go-inject-tag@latest
```

定义protobuf:
```protobuf
message IP {
  string Address = 1; // @gotags: valid:"ip" yaml:"ip" json:"overrided"
}
```

使用插件:
```
protoc-go-inject-tag -input=./test.pb.go
```

### 重构

使用protobuf 来定义接口和数据结构, 生成代码后 重构测试

## 框架进阶

现在我们没添加一个gRPC服务都很麻烦:
```go
// 初始化服务层 Ioc初始化
if err := impl.Service.Config(); err != nil {
  return err
}
pkg.Host = impl.Service

// Secret Service
if err := secretImpl.Service.Config(); err != nil {
  return err
}
pkg.Secret = secretImpl.Service

// Task Service
if err := taskImpl.Service.Config(); err != nil {
  return err
}
pkg.Task = taskImpl.Service

// resource Service
if err := searcher.Service.Config(); err != nil {
  return err
}
pkg.Resource = searcher.Service
```

暴露成HTTP服务也很麻烦:
```go
// 装置子服务路由
hostAPI.RegistAPI(s.r)
secretAPI.RegistAPI(s.r)
taskAPI.RegistAPI(s.r)
searchAPI.RegistAPI(s.r)
```

### Grpc服务重构

我们让我们更专注于app模块下的业务应用开发, 我们需要将这些逻辑做到框架层, 以注册的方式，
注册过后，其他逻辑有框架帮我们完成

我们的目标: 将实例化交给框架进行, 我们专注于 实现接口

#### Grpc Ioc层接口定义

```go
package app

import (
	"google.golang.org/grpc"
)

// GRPCService GRPC服务的实例
type GRPCApp interface {
	Registry(*grpc.Server)
	Config() error
	Name() string
}
```

#### grpc app注册管理

```go
var (
	grpcApps = map[string]GRPCApp{}
)

// RegistryService 服务实例注册
func RegistryGrpcApp(app GRPCApp) {
	// 已经注册的服务禁止再次注册
	_, ok := grpcApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("grpc app %s has registed", app.Name()))
	}

	grpcApps[app.Name()] = app
}

// LoadedGrpcApp 查询加载成功的服务
func LoadedGrpcApp() []string {
	return []string{}
}

func GetGrpcApp(name string) GRPCApp {
	app, ok := grpcApps[name]
	if !ok {
		panic(fmt.Sprintf("grpc app %s not registed", name))
	}

	return app
}

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

#### grpc模块 实现 GRPC APP接口定义

task服务实现了grpc接口, 并且通过init函数注册

```go
package grpc

import (
	"database/sql"

	"github.com/infraboard/cmdb/app"
	"github.com/infraboard/cmdb/app/host"
	"github.com/infraboard/cmdb/app/secret"
	"github.com/infraboard/cmdb/app/task"
	"github.com/infraboard/cmdb/conf"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db     *sql.DB
	log    logger.Logger
	host   host.ServiceServer
	secret secret.ServiceServer
	task.UnimplementedServiceServer
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)
	s.secret = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	return nil
}

func (s *service) Name() string {
	return task.AppName
}

func (s *service) Registry(server *grpc.Server) error {
	task.RegisterServiceServer(server, svr)
	return nil
}

func init() {
	app.RegistryGrpcApp(svr)
}
```


#### 统一导入所有Grpc APP

单独定义了all的包, 用于导入所有定义的app
```go
package all

import (
	_ "github.com/infraboard/cmdb/app/host/impl"
	_ "github.com/infraboard/cmdb/app/resource/impl"
	_ "github.com/infraboard/cmdb/app/secret/impl"
	_ "github.com/infraboard/cmdb/app/task/impl"
)
```

服务启动时 注册所有服务
```go
// 注册所有服务
_ "github.com/infraboard/cmdb/app/all"
```

#### 启动所有Grpc app

grpc服务启动时, 装置所有的grpc app
```go
// Start 启动GRPC服务
func (s *GRPCService) Start() {
	// 装载所有GRPC服务
	app.LoadGrpcApp(s.svr)

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPCAddr())
	if err != nil {
		s.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	s.l.Infof("GRPC 服务监听地址: %s", s.c.App.GRPCAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		s.l.Error("start grpc service error, %s", err.Error())
		return
	}
}
```

#### 如何获取依赖

通过GetGrpcApp获取依赖的服务, 为了确保服务名称能对应上, 我们为每个App 定义了一个 AppName的常量

```go
func (h *handler) Config() error {
	h.log = zap.L().Named(secret.AppName)
	h.service = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	return nil
}
```

#### 总结

到此 我们只需要编写我们对应的app服务的具体实现, 再也不用手动去写注册相关的代码了，这些都交给了外层框架

### http服务重构

重构的目的依然是: 业务和框架

#### http Ioc层接口定义

```go
// HTTPService Http服务的实例
type HTTPApp interface {
	Registry(router.SubRouter)
	Config() error
	Name() string
}
```

为了更好的对接后面的权限中心, 专注封装了httprouter, 为其定制权限注册和控制,
现在你可以把他当做httprouter使用

#### http app注册管理

```go
var (
	httpApps = map[string]HTTPApp{}
)

// RegistryHttpApp 服务实例注册
func RegistryHttpApp(app HTTPApp) {
	// 已经注册的服务禁止再次注册
	_, ok := httpApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", app.Name()))
	}

	httpApps[app.Name()] = app
}

// LoadedGrpcApp 查询加载成功的服务
func LoadedHttpApp() []string {
	return []string{}
}

func GetHttpApp(name string) HTTPApp {
	app, ok := httpApps[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApp 装载所有的http app
func LoadHttpApp(pathPrefix string, root router.Router) error {
	for _, api := range httpApps {
		if err := api.Config(); err != nil {
			return err
		}
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.SubRouter(pathPrefix + "/api/v1"))
	}
	return nil
}
```


#### 实现并注册HTTP app

```go
package http

import (
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/cmdb/app"
	"github.com/infraboard/cmdb/app/host"
)

var (
	h = &handler{}
)

type handler struct {
	service host.ServiceServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(host.AppName)
	h.service = app.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return host.AppName
}

func (h *handler) Registry(r router.SubRouter) {
	r.Handle("GET", "/hosts", h.QueryHost)
	r.Handle("POST", "/hosts", h.CreateHost)
	r.Handle("GET", "/hosts/:id", h.DescribeHost)
	r.Handle("DELETE", "/hosts/:id", h.DeleteHost)
	r.Handle("PUT", "/hosts/:id", h.PutHost)
	r.Handle("PATCH", "/hosts/:id", h.PatchHost)
}

func init() {
	app.RegistryHttpApp(h)
}
```

#### 统一导入所有Http app

```go
package all

import (
	_ "github.com/infraboard/cmdb/app/host/http"
	_ "github.com/infraboard/cmdb/app/resource/http"
	_ "github.com/infraboard/cmdb/app/secret/http"
	_ "github.com/infraboard/cmdb/app/task/http"
)
```

服务启动时 注册所有服务
```go
// 注册所有服务
_ "github.com/infraboard/cmdb/app/all"
```


#### 启动所有导入的 app

```go
// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	if err := app.LoadHttpApp(s.c.App.Name, s.r); err != nil {
		return err
	}

	// 启动 HTTP服务
	s.l.Infof("HTTP服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}
```

#### 总结

到此 我们就可以专注于写 http handler了, 无效关注其他框架层的事儿


## 总结

经过改造, 我们编写一个app的流程变更为:

+ 定义接口
+ 实现grpc服务并注册
+ 使用http服务并注册
+ all文件引入模块

