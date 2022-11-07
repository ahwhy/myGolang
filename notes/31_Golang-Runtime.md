# Golang-Runtime  Golang的运行时

## 一、Golang的标准库 runtime包

### 1. runtime
- runtime
	- runtime包提供和go运行时环境的互操作，如控制go程的函数
	- 它也包括用于reflect包的低层次类型信息
	- 也包括用于reflect包的低层次类型信息

### 2. Environment Variables
- 环境变量
	- $name或%name%，这依赖于主机的操作系统，其控制go程序的运行时行为
	- 它们的含义和用法可能在各发行版之间改变

- GOGC
	- 设置最初的垃圾收集目标百分比
	- 当新申请的数据和前次垃圾收集剩下的存活数据的比率达到该百分比时，就会触发垃圾收集
	- 默认 `GOGC=100`
	- 设置 `GOGC=off` 会完全关闭垃圾收集
	- runtime/debug包 的 `SetGCPercent`函数允许在运行时修改该百分比

- GODEBUG
	- 控制运行时的debug输出
	- GODEBUG的值是逗号分隔的name=val对
```
	allocfreetrace: 设置其为1，会导致每次分配都会被记录剖面，会记录每一个对象的分配、释放及其堆栈踪迹。
	efence: 设置其为1，会导致分配器运行模式为：每个对象申请在独立的页和地址，且永不循环利用。
	gctrace: 设置其为1，会导致垃圾收集器每次收集都向标准错误输出写入单行的数据，概述收集的总内存的大小和暂停的总时间长度。设置其为2，会写入同样的概述，但也会写入每次收集的两个数据。
	gcdead: 设置其为1，会导致垃圾收集器摧毁任何它认为已经死掉的执行堆栈。
	schedtrace: 设置其为X，会导致调度程序每隔X毫秒输出单行信息到标准错误输出，概述调度状态。
	scheddetail: 设置schedtrace为X并设置其为1，会导致调度程序每隔X毫秒输出详细的多行信息，描述调度、进程、线程和go程的状态。
```

- GOMAXPROCS
	- 限制可以同时运行用户层次的go代码的操作系统进程数
	- 没有对代表go代码的、可以在系统调用中阻塞的go程数的限制
	- 那些阻塞的go程不与GOMAXPROCS限制冲突
	- 本包的`GOMAXPROCS`函数可以查询和修改该限制

- GOTRACEBACK
	- 控制当go程序因为不能恢复的panic或不期望的运行时情况失败时的输出
	- 失败的程序默认会打印所有现存go程的堆栈踪迹(省略运行时系统中的函数)，然后以状态码2退出
	- 如果GOTRACEBACK为0，会完全忽略所有go程的堆栈踪迹
	- 如果GOTRACEBACK为1，会采用默认行为。如果GOTRACEBACK为2，会打印所有现存go程包括运行时函数的堆栈踪迹
	- 如果GOTRACEBACK为crash，会打印所有现存go程包括运行时函数的堆栈踪迹，并且如果可能会采用操作系统特定的方式崩溃，而不是退出
	- 例如，在Unix系统里，程序会释放SIGABRT信号以触发核心信息转储

- 其他
	- 环境变量GOARCH、GOOS、GOPATH和GOROOT构成完整的go环境变量集合，它们影响go程序的构建
	- 环境变量GOARCH、GOOS和GOROOT在编译时被记录并可用本包的常量和函数获取，但它们不会影响运行时环境


## 二、Golang的标准库 cgo包

- runtime/cgo
	- cgo 包含有 cgo 工具生成的代码的运行时支持

## 三、Golang的标准库 debug包

- runtime/debug
	- debug包 包含程序在运行时自我调试的工具

## 四、Golang的标准库 pprof包

- runtime/pprof
	- pprof包以pprof可视化工具期望的格式书写运行时剖面数据

## 五、Golang的标准库 race包

- runtime/race
	- race包实现了数据竞争检测逻辑
	- 没有提供公共接口

## 六、Golang的标准库 trace包

- runtime/trace
	- 执行追踪器；tracer捕获各种执行事件，如goroutine创建/阻塞/解除阻塞、syscall进入/退出/阻塞、GC相关事件、堆大小的变化、处理器启动/停止等，并将它们写入io
	- 对于大多数事件，都会捕获精确到纳秒级的时间戳和堆栈跟踪
	- 使用 `go tool trace` 命令来分析跟踪