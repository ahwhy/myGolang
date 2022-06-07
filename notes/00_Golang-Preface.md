# Golang-Preface  Golang的序言

## 一、初识Go语言

### Golang的历程

Go 语言是由谷歌的开发工程师(罗伯特·格瑞史莫、罗勃·派克、肯·汤普逊等)于 2007 年开始设计，利用 20%的自由时间开发的实验项目，并于 2009 年以 BSD-style 授权(完全开源)首次公开发布，于 2012 年正式发布

- [历史上的Go版本](https://golang.google.cn/doc/devel/release.html)

- go1.16 (released 2021/02/16)
	- runtime/metrics
	- GO111MODULE的默认值为on, 可以把GO111MODULE设置回auto
	- io/ioutil为废弃，函数被转移到了os和io这两个包里
	- 1.16优化了链接器，现在它在linux/amd64上比1.15快了20-25%，内存占用减少了5-15%

- go1.15 (released 2020/08/11)
	- Go 1.15版本不再对darwin/386和darwin/arm两个32位平台提供支持
	- 新版linker的起点，新链接器的性能要提高20%，内存占用减少30% 后续持续优化
	- GOPROXY支持设置为多个proxy的列表，多个proxy之间采用逗号分隔,默认值依旧为https://proxy.golang.org,direct

- go1.14 (released 2020/02/25)
	- Go Module 已经可以用于生产环境
	- defer zero overhead

- go1.13 (released 2019/09/03)
	- 重写了逃逸分析，减少了 Go 程序中堆上的内存申请的空间，1.3版本被诟病后改了
	- defer 性能提升30%
	- GO111MODULE auto

- go1.12 (released 2019/02/25)
	- 内存分配优化
	- GO111MODULE off
 
- go1.11 (released 2018/08/24)
	- Go modules，社区积怨已久
	- 实验性的 WebAssembly

- go1.10 (released 2018/02/16)
	- gc快优化到头了
	- 为了加快构建速度，go build 引入了构建缓存

- go1.9 (released 2017/08/24)
	- type alias
	- Concurrent Map  sync.Map结构

- go1.8 (released 2017/02/16)
	- gc 的停顿时间减少到了 1 毫秒以下(100 microseconds and often as low as 10 microseconds)
	- defer 开销降低到原来的一半

- go1.7 (released 2016/08/15)
	- 对编译工具链也作了优化，编译速度更快
	- gc进一步提升
	- Context包推出 对goroutine的生命周期管理提出的标准方案
 
- go1.6 (released 2016/02/17)
	- gc进一步提升

- go1.5 (released 2015/08/19) 发布时间推迟了两个月2个大招
	- [Go 1.5 concurrent garbage collector pacing](https://docs.google.com/document/d/1wmjrocXIWTr1JxU-3EQBI6BK6KgtiFArkG47XK73xIQ/edit#)
	- 完成了自举，.c文件被全部重写，用Go语言编写Go语言的编译器，来回测X，直到全部通过

- go1.4 (released 2014/12/10) 
	- 继续优化GC，  go generate诞生

- go1.3 (released 2014/06/18) 
	- 尝试优化GC，同时推出了 sync.Pool

- go1.2 (released 2013/12/01) 语言层面的优化
	- 基础库的性能提升 
	- SetMaxThreads defaults (10000)

- go1.1 (released 2013/05/13) 
	- [Scalable Go Scheduler Design Doc](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#heading=h.mmq8lm48qfcw)

- go1 (released 2012/03/28)  
	- [Go 1 and the Future of Go Programs](https://golang.org/doc/go1compat)


### Golang的优势

- 简单的部署方式
	- 可直接编译成机械码
	- 不依赖其他库
	- 直接运行即可部署

- 静态类型语言
	- 编译的时候可以检测出隐藏的大多数问题
	- 强类型方便阅读与重构

- 语言层面的并发
	- 天生的基于支持
	- 充分利用多核

- 工程化比较优秀
	- GoDoc 可以直接从代码和注释生成漂亮的文档
	- GoFmt 统一的代码格式
	- GoLint 代码语法提示
	- 测试框架内置

- 强大的标准库
	- Runtime系统调度机制
	- 高效的GC垃圾回收
	- 丰富的标准库

- 简单易学
	- 25个关键字
	- C语言简洁基因，内嵌C语言语法支持
	- 面向对象特征(继承, 多态，封装)
	- 跨平台


### Golang的缺陷

- 包管理，大部分都在 github上
	- 作者修改项目名称
	- 作者删库跑路
	- vendor 到 mod 迁移麻烦，很多遗留依赖问题

- 无泛型， 2.0 有计划加上 (传言)
	- interface{}可以解决该问题, 但是不易于代码阅读

- 没有Exception，使用Error来处理异常
	- error处理不太优雅, 很多重复代码


### Golang的应用

- Go 语言主要用于服务端开发，其定位是开发大型软件，常用于
	- 服务器编程: 日志处理、数据打包、虚拟机处理、文件系统、分布式系统、数据库代理等 
	- 网络编程: Web 应用、API 应用、下载应用
	- 内存数据库
	- 云平台
	- 机器学习
	- 区块链
	- ……

- [使用 Go 开发的项目列表](https://github.com/golang/go/wiki/Projects)
	- Go 
	- docker
	- kubernetes
	- lantern
	- etcd
	- Prometheus
	- Influxdb
	- Consul
	- nsq
	- nats
	- beego
	- ……

- [使用 Go 开发的组织](http://go-lang.cat-v.org/organizations-using-go)
  - 国外: Google、CloudFlare……
  - 国内: 阿里、腾讯、百度、京东、爱奇艺、小米、今日头条、滴滴、美团、饿了么、360、七牛、B 站、盛大、搜狗……


## 二、学习Golang的一些方法

- [Golang官网](https://golang.google.cn/doc/)

- 不错的gitbook入门: 
    - [Go语言圣经](https://yar999.gitbook.io/gopl-zh/) 
    - [Go Web 编程](https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/) 
    - 不错的微信公众号: 现在还在更新的就只有GoCN了
	
- 阅读一些优秀的项目的源码
    - [kubernetes](https://github.com/kubernetes/kubernetes)
    - [nats](https://github.com/nats-io/nats-server)
    - [telegraf](https://github.com/influxdata/telegraf)
    - [beats](https://github.com/elastic/beats)
    - [etcd raft](https://github.com/etcd-io/etcd/tree/main/raft)