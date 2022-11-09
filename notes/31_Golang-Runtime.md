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

### 3. Constants && Variables
```go
	// Compiler 是编译工具链的名字，工具链会构建可执行的二进制文件
	const Compiler = "gc"
	// GOARCH是可执行程序的目标处理器架构(将要在该架构的机器上执行)：386、amd64或arm
	const GOARCH string = theGoarch
	// GOOS是可执行程序的目标操作系统(将要在该操作系统的机器上执行)：darwin、freebsd、linux等
	const GOOS string = theGoos

	// MemProfileRate 控制会在内存profile里记录和报告的内存分配采样频率
	var MemProfileRate int = 512 * 1024
```

### 4. Functions
```go
	type Error interface {
		error
		// RuntimeError是一个无操作的函数，仅用于区别运行时错误和普通错误。
		// 具有RuntimeError方法的错误类型就是运行时错误类型。
		RuntimeError()
	}

	// TypeAssertionError 表示一次失败的类型断言
	type TypeAssertionError struct { ... }
	func (e *TypeAssertionError) Error() string
	func (e *TypeAssertionError) RuntimeError()

	// GOROOT 返回Go的根目录
	func GOROOT() string
	// 返回Go的版本字符串
	func Version() string
	// NumCPU 返回本地机器的逻辑CPU个数
	func NumCPU() int
	// GOMAXPROCS 设置可同时执行的最大CPU数，并返回先前的设置
	func GOMAXPROCS(n int) int
	// SetCPUProfileRate 设置CPU profile记录的速率为平均每秒hz次
	func SetCPUProfileRate(hz int)
	// CPUProfile 返回二进制CPU profile堆栈跟踪数据的下一个chunk，函数会阻塞直到该数据可用
	func CPUProfile() []byte

	// GC 执行一次垃圾回收
	func GC()
	// SetFinalizer将x的终止器设置为f
	// 当垃圾收集器发现一个不能接触的(即引用计数为零，程序中不能再直接或间接访问该对象)具有终止器的块时，它会清理该关联(对象到终止器)并在独立go程调用f(x)，这使x再次可以接触，但没有了绑定的终止器
	// 如果SetFinalizer没有被再次调用，下一次垃圾收集器将视x为不可接触的，并释放x
	func SetFinalizer(x, f interface{})

	// MemProfile 返回当前内存profile中的记录数n
	// 大多数调用者应当使用runtime/pprof包或testing包的-test.memprofile标记，而非直接调用MemProfile
	func MemProfile(p []MemProfileRecord, inuseZero bool) (n int, ok bool)
	// Breakpoint 执行一个断点陷阱
	func Breakpoint()
	// Stack 将调用其的go程的调用栈踪迹格式化后写入到buf中并返回写入的字节数
	func Stack(buf []byte, all bool) int
	// Caller 报告当前go程调用栈所执行的函数的文件和行号信息
	func Caller(skip int) (pc uintptr, file string, line int, ok bool)
	// Callers 把当前go程调用栈上的调用栈标识符填入切片pc中，返回写入到pc中的项数
	func Callers(skip int, pc []uintptr) int

	// StackRecord 描述单条调用栈
	type StackRecord struct {
		Stack0 [32]uintptr // 该记录的调用栈踪迹，以第一个零值成员截止
	}
	// Stack 返回与记录相关联的调用栈踪迹，即r.Stack0的前缀
	func (r *StackRecord) Stack() []uintptr

	type Func struct { ... }
	// FuncForPC返回一个表示调用栈标识符pc对应的调用栈的*Func；如果该调用栈标识符没有对应的调用栈，函数会返回nil；每一个调用栈必然是对某个函数的调用
	func FuncForPC(pc uintptr) *Func

	// Name 返回该调用栈所调用的函数的名字
	func (f *Func) Name() string
	// FileLine 返回该调用栈所调用的函数的源代码文件名和行号
	// 如果pc不是f内的调用栈标识符，结果是不精确的
	func (f *Func) FileLine(pc uintptr) (file string, line int)
	// Entry 返回该调用栈的调用栈标识符
	func (f *Func) Entry() uintptr
	// NumCgoCall 返回当前进程执行的cgo调用次数
	func NumCgoCall() int64
	// NumGoroutine 返回当前存在的Go程数
	func NumGoroutine() int
	// Goexit 终止调用它的go程；其它go程不会受影响
	// Goexit 会在终止该go程前执行所有defer的函数
	func Goexit()
	// Gosched 使当前go程放弃处理器，以让其它go程运行
	func Gosched()
	// GoroutineProfile 返回活跃go程的堆栈profile中的记录个数
	func GoroutineProfile(p []StackRecord) (n int, ok bool)
	// LockOSThread 将调用的go程绑定到它当前所在的操作系统线程
	func LockOSThread()
	// UnlockOSThread 将调用的go程解除和它绑定的操作系统线程
	func UnlockOSThread()
	// ThreadCreateProfile 返回线程创建profile中的记录个数
	func ThreadCreateProfile(p []StackRecord) (n int, ok bool)

	// BlockProfileRecord 用于描述某个调用栈序列发生的阻塞事件的信息
	type BlockProfileRecord struct {
		Count  int64
		Cycles int64
		StackRecord
	}
	// SetBlockProfileRate 控制阻塞profile记录go程阻塞事件的采样频率
	func SetBlockProfileRate(rate int)
	// BlockProfile 返回当前阻塞profile中的记录个数
	func BlockProfile(p []BlockProfileRecord) (n int, ok bool)
```
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