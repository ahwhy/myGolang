# 项目部署

到此我们开发完成的服务有:
+ devcloud: 项目前端
+ cmdb: 资产中心
+ keyauth: 权限中心

接下来我们就看看如何把这套服务部署上

## 部署架构


## 依赖部署

+ mysql: cmdb依赖
+ mongodb: keyauth依赖
+ nginx: devcloud依赖

## keyauth部署

我们本地开发的启动方式:
```sh
go run main.go start -f etc/keyauth.toml
```

### 编译

我们需要部署到的服务器是Linux, 因此我们需要交叉编译出一个Linux版本的 keyauth二进制文件:
```sh
GOOS=linux GOARCH=amd64 go build -o keyauth main.go
```

为了编译瘦身, 降低二进制包的大小, 添加一些编译参数: 通过ldflags可以传染一些参数，控制编译的过程
+ -s 的作用是去掉符号信息。去掉符号表，golang panic 时 stack trace 就看不到文件名及出错行号信息了。
+ -w 的作用是去掉 DWARF tables 调试信息。结果就是得到的程序就不能用 gdb 调试了

```sh
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o keyauth main.go
```

为了给编译的产物 打上一些 当前编译信息, 需要定义个version包: 该包定义了4个导出变量:
+ GIT_BRANCH: 当前构建代码的commit号
+ GIT_COMMIT: 当前分支
+ BUILD_TIME: 构建时间
+ GO_VERSION: go版本信息

```go
package version

import (
	"fmt"
)

const (
	// ServiceName 服务名称
	ServiceName = "keyauth"

	// Description 服务描述
	Description = "微服务权限中心"
)

var (
	GIT_COMMIT string
	GIT_BRANCH string
	BUILD_TIME string
	GO_VERSION string
)

// FullVersion show the version info
func FullVersion() string {
	version := fmt.Sprintf("Build Time: %s\nGit Branch: %s\nGit Commit: %s\nGo Version: %s\n", BUILD_TIME, GIT_BRANCH, GIT_COMMIT, GO_VERSION)
	return version
}

// Short 版本缩写
func Short() string {
	commit := ""
	if len(GIT_COMMIT) > 8 {
		commit = GIT_COMMIT[:8]
	}
	return fmt.Sprintf("%s[%s]", GIT_BRANCH, commit)
}
```

那如何在编译的时候注入这些变量喃?
```sh
# 使用 -X 参数 pkg.VAR=value 这种方式往包中注入的变量注入值: 
# -ldflags "-X gitee.com/infraboard/go-course/day14/demo/api/version.GIT_TAG='v0.0.1'"

-ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}'
```

配合上我们的Makefile， 这个指令就是这样的
```sh
PROJECT_NAME=keyauth
MAIN_FILE=main.go
PKG := "github.com/infraboard/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "${PKG}/version"

build: dep ## Build the binary file
	@go build -a -o dist/${PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${MAIN_FILE}

linux: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${PROJECT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${MAIN_FILE}
```

改文件已经在keyauth的Makefile中编写好了, 我们直接make linux就可以了
```sh
# 先构建个本地版本的试试
make build
# 打印注入的构建信息
./dist/keyauth -v
Build Time: 2021-11-19 21:47:54
Git Branch: master
Git Commit: 1b40552b29d253f385da2e19f3b7d298088f60c9
Go Version: go1.16.5 darwin/amd64
# 最后我们构建一个linux版本
make linux
```

### 部署

我们将dist/keyauth 的linux版本 cp到服务器上启动即可

```sh
# 上传到服务器
scp -P 4774  dist/keyauth user@ip:~
## copy到 我们的PATH里面去
mv /home/etb/keyauth /usr/local/bin/
```

添加配置文件: /etc/infraboard/keyauth/keyauth.toml
```
mkdir -pv /etc/infraboard/keyauth
vim /etc/infraboard/keyauth/keyauth.toml
```
下面是配置文件的样例:
```toml
[app]
name = "keyauth"
host = "0.0.0.0"
port = "8050"
key  = "this is your app key"

[mongodb]
endpoints = ["0.0.0.0:17232"]
username = "keyauth"
password = "xxx"
database = "keyauth"

[log]
level = "debug"
path = "logs"
format = "text"
to = "stdout"
```

先尝试启动服务:
```sh
keyauth start -f /etc/infraboard/keyauth/keyauth.toml
```

能正常启动服务后，我们开始初始化我们的keyauth服务:
```sh
[etb@VM_0_7_centos ~]$ keyauth init -f /etc/infraboard/keyauth/keyauth.toml 
2021-11-19T22:27:19.062+0800    INFO    [INIT]  cmd/start.go:200        log level: debug
2021-11-19T22:27:19.062+0800    INFO    [INIT]  cmd/start.go:212        use cache in local memory
2021-11-19T22:27:19.740+0800    DEBUG   [User]  impl/dao.go:27  find filter: map[domain: type:SUPPER]
? 请输入公司(组织)名称: 基础设施服务中心
? 请输入管理员用户名称: admin
? 请输入管理员密码: ******
? 再次输入管理员密码: ******

开始初始化...
初始化用户: admin [成功]
初始化域: 基础设施服务中心   [成功]
初始化应用: admin-web [成功]
应用客户端ID: wIJvf2gOcsz1ym0XZ2vA96kl
应用客户端凭证: Y3sqeVcYNUqf8shJt77qsmj0AipVgtj2
初始化应用: admin-micro [成功]
应用客户端ID: oB5VEpOTQNN7tn3PhvEe0p2D
应用客户端凭证: ajr9FkrXh14OuPm0xneERE676Gzgm106
2021-11-19T22:27:27.314+0800    DEBUG   [Role]  impl/permission.go:27   query permission filter: map[role_id:c6br8ju1l0cvabpa7fcg]
2021-11-19T22:27:27.316+0800    DEBUG   [Role]  impl/permission.go:27   query permission filter: map[role_id:c6br8ju1l0cvabpa7fd0]
初始化角色: admin [成功]
初始化角色: visitor [成功]
初始化根部门: 基础设施服务中心   [成功]
初始化默认空间: 系统初始化创建   [成功]
初始化系统配置: v1   [成功]
```

这里要注意, 保存好客户端凭证, 后面前端访问需要:
+ 应用客户端ID
+ 应用客户端凭证

直接通过nohup启动服务:
```sh
nohup keyauth start -f /etc/infraboard/keyauth/keyauth.toml &> keyauth.log &
```

可以看到服务已经开始监听端口:
```sh
[etb@VM_0_7_centos ~]$ netstat -tlnup | grep keyatuh
tcp6       0      0 :::18050                :::*                    LISTEN      14447/keyauth       
tcp6       0      0 :::8050                 :::*                    LISTEN      14447/keyauth 
```

当然你也可以使用systemd 来管理你的服务

### 配置代理

由于我们之前的keyauth的后台页面 还没有迁移过devcloud来, 所以需要

我们等下需要使用API来创建一个cmdb服务, 这样cmdb服务才可以接入
```nginx
server {
    listen       80;
    server_name  devcloud.nbtuan.vip;
    root         /usr/share/nginx/devcloud;

    # 说明下vue-router的默认hash模式——使用URL的hash来模拟一个完整的URL，当URL改变时，页面不会重新加载
    # history 模式 这种模式充分利用 history.pushState API 来完成 URL 跳转而无须重新加载页面
    try_files $uri $uri/ /index.html;

    # Load configuration files for the default server block.
    include /etc/nginx/default.d/*.conf;

    location /keyauth/api/v1/ {
    proxy_set_header X-Real-IP  $remote_addr;
    proxy_set_header Host-Real-IP  $http_host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_pass http://127.0.0.1:8050;
    }

    location /cmdb/api/v1/ {
    proxy_set_header X-Real-IP  $remote_addr;
    proxy_set_header Host-Real-IP  $http_host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_pass http://127.0.0.1:8060;
    }

    error_page 404 /404.html;
        location = /40x.html {
    }

    error_page 500 502 503 504 /50x.html;
        location = /50x.html {
    }
}
```

重启下我们的代理服务
```sh
# 如果你只是修改server 可以使用reload
systemctl restart nginx
```

此时我们的 keyauth API server 就暴露出去了,  可以通过API访问测试下:
```
curl http://ip:port/keyauth/api/v1/namespaces/
```

## devcloud部署

作为前后端分离的项目, 我们首先需要配置 后端API网关的地址: http://devcloud.nbtuan.vip

然后构建前端时 使用生产的网关地址

### 开发环境

vue 支持通过环境变量的方式加载配置
+ .env.development: 开发环境变量文件
+ .env.production: 测试环境变量文件

比如我们.env.production文件:
```
# just a flag
NODE_ENV = 'production'

# base api
VUE_APP_BASE_API = 'http://devcloud.nbtuan.vip'

// client
VUE_APP_CLIENT_ID = 'xxx'
VUE_APP_CLIENT_SECRET = 'xxx'
```

我们通过`process.env.VUE_APP_CLIENT_ID`的方式读取改变量
```js
export function LOGIN(data, params) {
  return request({
    url: `${keyauth.baseURL}/oauth2/tokens`,
    method: "post",
    auth: {
      username: process.env.VUE_APP_CLIENT_ID,
      password: process.env.VUE_APP_CLIENT_SECRET,
    },
    data,
    params,
  });
}
```

但是如果才知道读取的是生产环境的变量还是开发环境变量喃?

我们通过--mode 在运行时传递, 默认是development, 这里我们添加了production参数, 对应的就是.env.production文件
```json
  "scripts": {
    "dev": "vue-cli-service serve",
    "build": "vue-cli-service build --mode production",
    "lint": "vue-cli-service lint",
    "svgo": "svgo -f src/icons/svg --config=src/icons/svgo.yml"
  },
```

### 构建

配置生产环境变量 .env.production:
```
# just a flag
NODE_ENV = 'production'

# base api
VUE_APP_BASE_API = 'http://devcloud.nbtuan.vip'

// client
VUE_APP_CLIENT_ID = 'oB5VEpOTQNN7tn3PhvEe0p2D'
VUE_APP_CLIENT_SECRET = 'ajr9FkrXh14OuPm0xneERE676Gzgm106'
```

然后我们开始构建:
```sh
npm run build
```

> 注意: 如果 你node版本过高, 比如我 v16.13.0, 会报如下错误:
```sh
npm ERR! 2 warnings generated.
npm ERR! In file included from ../src/binding.cpp:1:
npm ERR! In file included from ../../nan/nan.h:58:
npm ERR! In file included from /Users/g7/.node-gyp/16.13.0/include/node/node.h:63:
npm ERR! In file included from /Users/g7/.node-gyp/16.13.0/include/node/v8.h:30:
npm ERR! /Users/g7/.node-gyp/16.13.0/include/node/v8-internal.h:492:38: error: no template named 'remove_cv_t' in namespace 'std'; did you mean 'remove_cv'?
npm ERR!             !std::is_same<Data, std::remove_cv_t<T>>::value>::Perform(data);
npm ERR!                                 ~~~~~^~~~~~~~~~~
npm ERR!                                      remove_cv
npm ERR! /Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/c++/v1/type_traits:776:50: note: 'remove_cv' declared here
npm ERR! template <class _Tp> struct _LIBCPP_TEMPLATE_VIS remove_cv
```

这是我们项目依赖中 node-sass 的版本不支持 Node.js16 ，只能在 15 及以下的环境下运行, 可以通过nvm来管理多个版本, 或者
```
根据提示: no template named 'remove_cv_t' in namespace 'std'; did you mean 'remove_cv'?

于是前往/Users/g7/.node-gyp/16.13.0/include/node/v8-internal.h，将492行的remove_cv_t改为remove_cv
```

这我们前端的产物就是dist目录下放好了

### 部署

上传前端到我们服务器:
```
scp -r -P 4774  dist/* user@ip:~/devcloud
```

将dist下的产物配置到 服务器的一个静态站点上, 比如我们上面的配置
```
    listen       80;
    server_name  devcloud.nbtuan.vip;
    root         /usr/share/nginx/devcloud;

    # 说明下vue-router的默认hash模式——使用URL的hash来模拟一个完整的URL，当URL改变时，页面不会重新加载
    # history 模式 这种模式充分利用 history.pushState API 来完成 URL 跳转而无须重新加载页面
    try_files $uri $uri/ /index.html;
```

这样我们的前端就部署完成了, 然后登录验证下

## cmdb部署

我们本地开发的启动方式:
```sh
go run main.go start -f etc/cmdb-api.toml
```

### 构建

cmdb服务之前没有打印版本信息, 因此在root cmd里面, 开启-v参数
```go
var vers bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cmdb-api",
	Short: "cmdb-api 管理系统",
	Long:  `cmdb-api ...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return errors.New("no flags find")
	},
}

func init() {
    ...
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the cmdb version")
}
```

与keyauth一样, 提前写好makefile
```sh
make linux
```

### 部署

```toml
[app]
name = "cmdb"
http_host = "0.0.0.0"
http_port = "8060"
grpc_host = "0.0.0.0"
grpc_port = "18060"
key  = "this is your app key"

[mysql]
host = "xxx"
port = "xxxx"
username = "cmdb"
password = "xxxx"
database = "cmdb"

[keyauth]
host = "127.0.0.1"
port = "18050"
client_id = "pz3HiVQA3indzSHzFKtLHaJW"
client_secret = "vDvlAtqN3rS9CZcHugXp6QBuk28zRjud"

[log]
level = "debug"
path = "logs"
format = "text"
to = "stdout"
```



## 参考

+ [Systemctl和Unit file](https://blog.51cto.com/u_15131458/3119541)
+ [Systemd的Unit文件; systemctl增加服务详细介绍](https://blog.csdn.net/shuaixingi/article/details/49641721)

