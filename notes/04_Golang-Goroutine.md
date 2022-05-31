# Golang-Goroutine  Golang的并发编程

## 一、Golang的并发编程定义
- 并发编程开发，将一个过程按照 并行算法 拆分为多个可以独立执行的代码块，从而充分利用多核和多处理器提高系统吞吐率
	- Go采用的并发编程思想是CSP(Communicating Sequential Process，通讯顺序进程)，CSP有着精确的数学模型，其思想的核心是同步通讯
	- Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，Goroutine 之间会通过 Channel 传递数据
	- Channel遵循了先进先出的设计(FIFO)，具体规则如下
		- 先从 Channel 读取数据的 Goroutine 会先接收到数据
		- 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利
		- 因此可以把channel当作队列来使用，这也是Goroutine的通信方式

- 顺序、并发与并行
	- CPU调度
		- CPU是支持分时调度的，以时间片的形式来跑指令集
		- OS层面操作系统会轮换调度进程(运行中的程序)的一部分指令集给CPU运行
		- 硬件(CPU)是串行进行处理的，但是进程在OS层面是并发的
	- 顺序 指所有的指令都是以串行的方式执行，在相同的时刻有且仅有一个CPU在顺序执行程序的指令
	- 并发 指同一时间应对(dealing with)多件事情的能力，并发更关注的是程序的设计层面
		- 并发程序通信方式
			- 共享数据 -> 同步
				- 多个并发程序需要对同一个资源进行访问，则需要先申请资源的访问权限，同时再使用完成后释放资源的访问权
				- 当资源被其他程序已申请访问权后，程序应该等待访问权被释放并被申请到时进行访问操作
				- 同一时间资源只能被一个程序访问和操作
			- 管道 -> 异步
				- 数据处理者处理完数据后将数据放入缓冲区中，数据接收者从缓冲区中获取数据，处理者不用等待接收者是否准备好处理数据
	- 并行 指同一时间动手(doing)做多件事情的能力，并行更关注的是程序的运行层面
		- 任何语言的并行，到操作系统层面，都是内核线程的并行
		- 并行架构
			- 位级(bit-level)并行，比如寄存器从32位升级到64位
			- 指令级(instruction-level)并行，比如CPU指令流水线，乱序执行和猜测执行
			- 数据级(data-level)并行，比如CPU的SIMD指令以及GPU向量计算指令
			- 任务级(task-level)并行，比如多处理器架构以及分布式系统架构
		- 并发与并行
			- 并发是问题域中的概念，程序需要被设计成能够处理多个同时(或者几乎同时)发生的事情
			- 而并行则是方法域中的概念，通过将问题中的多个部分并行执行来加速解决问题

## 二、Golang的协程(goroutine)

### 1. 进程&&线程
- 进程: 资源分配的基本单位
	- 独占资源: 地址空间、网络、内存地址、文件句柄
	- 进程是操作系统资源分配的最小单元
	- 所有的进程都是由操作系统的内核管理的
	- 每个进程之间是独立的，每一个进程都会有自己单独的内存空间以及上下文信息
	- 一个进程挂了不会影响其他进程的运行
	- 进程通过内核的ipc进行通信
	- 操作系统进程并发的模型
		- Master/Worker模型
		- 缺点 -> 每次创建进程的结构都需要复制一遍(Fork), 开销很大

- 线程: CPU调度的基本单位
	- 独占资源: 栈、局部变量、程序计数器
	- 线程更加轻量级，创建和销毁的成本都很低
	- 线程之间通信以及共享内存非常方便
	- 多进程相比开销要小得多，但频繁的创建和销毁线程，仍然是不小的负担
	- 线程池
		- 创建一大批线程放入线程池当中，需要用的时候拿出来使用，用完了再放回
		- 优点: 复用线程，回收和领用代替了创建和销毁两个操作，大大提升了性能
		- 缺点: 资源的共享，由于线程之间资源共享更加频繁，所以在一些场景当中需要加上锁等设计，避免并发带来的数据紊乱
	- 每个系统级线程都会有一个固定大小的栈(一般默认可能是2MB)，这个栈主要用来保存函数递归调用时参数和局部变量

### 2. 协程
- Goroutine 轻量级线程
	- 协程最大优势在于"轻量级"，可以轻松创建上百万个而不会导致系统资源衰竭
		- 线程和进程通常最多也不能超过1万
	- 创建Goroutine时为其分配4k堆栈，随着程序的执行自动增长删除
		- 当遇到深度递归导致当前栈空间不足时，Goroutine会根据需要动态地伸缩栈的大小(主流实现中栈的最大值可达到1GB)
		- 创建线程时必须指定堆栈且是固定的，通常以M为单位
	- goroutine是由Go runtime负责管理的，创建和销毁的消耗非常小，是用户级
		- Thread创建和销毀会有巨大的消耗，要和操作系统打交道，是内核级的
		- 通常解决的办法就是线程池
	- goroutines 切换只需保存三个寄存器，约200纳秒
		- 线程切换时需要保存各种寄存器状态，以便恢复，约1000-1500纳秒
	- go的协程是非抢占式的，由Goruntime主动交出控制权(对于开发者而言是抢占式的)
		- 线程在时间片用完后，由 CPU 中断任务强行将其调度走，这时就必须多保存很多信息
	- 从进程到线程再到协程，其实是一个不断共享，不断减少切换成本的过程
	
- 协程的调度不是基于操作系统的而是基于用户空间的程序的
	- 一般由库或者程序的运行时提供调度
	- 根据具体函数只保存必要的寄存器，切换的代价要比系统线程低得多
	
- Go语言中的MPG模型
	- 协程调度器，其将协程调度给操作系统的线程运行，由Go的Runtime实现
	- 调度
		- 名词定义
			- G(Goroutine)
				- Goroutine协程
				- 本质上是一种轻量级的线程
			- P(Processor)
				- 虚拟处理器
				- 代表M所需的上下文环境，是处理用户级代码逻辑的处理器
				- P的数量由环境变量中的GOMAXPROCS决定，默认情况下就是核数
				- 所有的P都在程序启动时创建，并保存在数组中，最多有GOMAXPROCS(可配置)个
				- `func runtime.GOMAXPROCS(n int) int`
			- M(Machine)
				- M对应一个内核线程(Thread)，并且这个对应关系是确定的
				- 线程想运行任务就得获取P，从P的本地队列获取G
				- P队列为空时，M也会尝试从全局队列拿一批G放到P的本地队列或从其他P的本地队列偷一半放到自己P的本地队列
				- M运行G，G执行之后，M会从P获取下一个G，不断重复下去
				- M进入系统调用时，会抛弃P，P被挂到其他M上，然后继续执行G队列
				- 系统调用返回后，相应的G进入全局的可运行队列(runqueue)中，P会周期性扫描全局的runqueue，使上面的G得到执行
			- 全局队列(Global Queue)
				- 存放等待运行的G
			- P的本地队列
				- 同全局队列类似，存放的也是等待运行的G，但存的数量有限，不超过 256 个
				- 新建G’时，G’优先加入到P的本地队列，如果队列满了，则会把本地队列中一半的G移动到全局队列
		- 流程
			- 创建Goroutine
			- 调度器将Goroutine调度入队(本地队列或者全局队列)
			- M获取G执行(绑定的P的队列，全局队列，其他P的队列)
			- M 循环的调度 G 执行
			- 如果M阻塞，创建一个M或者从休眠的M队列中挑选一个M来服务于该P
			- G执行完成，销毁G，返回，此时M空闲，放入休眠中的M队列
	- 创建
		- go + 语句
		- main函数也是由一个协程来启动执行，这个协程称为主协程，其他协程叫工作协程
			- 主协程结束后工作协程也会随之销毁，即主协程不等待工作协程的结束
			- 无从属关系，只是两种协程类型
			- 通过 sync.WaitGroup(计数信号量)来维护执行协程执行状态
			- 通过 runtime包中的GoSched让协程主动让出CPU，或通过 time.Sleep让协程休眠从而让出CPU
	- 过程分析
		- Go语言中提供trace工具，用于分析程序的运行过程
			- 执行程序后，会生成trace.out文件，再运行go tool trace trace.out  
		- go run或go build时添加-race参数检查资源竞争
	- 闭包陷阱
		- 闭包使用函数外变量，当协程执行时，外部变量已经发生变化，导致打印内容不正确，可使用在创建协程时通过函数传递参数(值拷贝)方式避免

```go
	// 使用 sync.WaitGroup(计数信号量)来维护执行协程执行状态
	func printChars(prefix string, group *sync.WaitGroup) {
		defer group.Done()                      // 通过信号量通知执行结束  等待信号量 -1
		for ch := 'A'; ch <= 'Z'; ch++ {
			fmt.Printf("%s:%c\n", prefix, ch)
			runtime.Gosched()                   // 让出CPU的调度 同time.Sleep(1 * time.Microsecond)
		}
	}
	group := &sync.WaitGroup{}                  // 定义信号量  启动协程之前 +1，协程结束时 -1
	n := 10
	group.Add(n)                                // 定义信号量 10
	
	for i := 0; i < n; i++ {                    // 创建协程
		go printChars(fmt.Sprintf("gochars%0d\n", i), group)
	}
	group.Wait()                                // 等待所有协程执行结束
	fmt.Println("over")

	// trace工具
	// 创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	// 启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	
	// 启动goroutine
	asyncRun()
	wg.Wait()

	// 闭包陷阱
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)        //  全部打印10
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("-----------------------")
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)        //  打印0 ~ 9
			wg.Done()
		}(i)
	}
	wg.Wait()
```

## 三、Golang的管道(channel) 
- 定义
	- go语言中可以通过 chan 来定义管道
	- 通过操作符 `<-` 对管道进行读取和写入操作，通过管道维护协程状态

- 本质
	- 一个 mutex 锁加上一个环状缓存、一个发送方队列和一个接收方队列
		- 管道底层是一个环形队列(先进先出)，send(插入) 和 recv(取走) 从同一个位置沿同一个方向顺序执行
		- sendx 表示最后一次插入的元素
		- recvx 表示最后一次取走元素的位置
	- 发送过程包含三个步骤
		- 持有锁
		- 入队，拷贝要发送的数据
			- 找到是否有正在阻塞的接收方，是则直接发送
			- 找到是否有空余的缓存，是则存入
			- 阻塞直到被唤醒
		- 释放锁
	- 接收过程包含三个步骤
		- 上锁
		- 从缓存中出队，拷贝要接收的数据
			- 如果 Channel 已被关闭，且 Channel 没有数据，立刻返回
			- 如果存在正在阻塞的发送方，说明缓存已满，从缓存队头取一个数据，再复始一个阻塞的发送方
			- 否则，检查缓存，如果缓存中仍有数据，则从缓存中读取，读取过程会将队列中的数据拷贝一份到接收方的执行栈中
			- 没有能接收的数据，阻塞当前的接收方 Goroutine
		- 解锁
```go
	// src/runtime/chan.go
	type hchan struct {
		qcount   uint           // 队列中的所有数据数
		dataqsiz uint           // 环形队列的大小
		buf      unsafe.Pointer // 指向大小为 dataqsiz 的数组
		elemsize uint16         // 元素大小
		closed   uint32         // 是否关闭
		elemtype *_type         // 元素类型
		sendx    uint           // 发送索引
		recvx    uint           // 接收索引
		recvq    waitq          // recv 等待列表，即( <-ch )
		sendq    waitq          // send 等待列表，即( ch<- )
		lock mutex
	}
	type waitq struct {         // 等待队列 sudog 双向队列
		first *sudog
		last  *sudog
	}
```

- 声明
	- channel 管道，声明时需要指定管道存放数据的类型
	- 管道可以存放任何类型，但只建议用于存放值类型或者只包含值类型的结构体
	- 在管道声明后，会被初始化为 nil
```go
	var channel chan int
	fmt.Printf("%T %v", channel, channel)  // chan int <nil>
```

- 初始化
	- 使用 make函数初始化chan并分配内存
		- `make(chan type)` 不带 len参数创建无缓存区的管道
			- 无缓冲管道(unbuffered channel)特点
				- sender端向channel中send一个数据，然后阻塞，直到receiver端将数据receive
				- receiver端一直阻塞，直到sender端向channel发送了一个数据
				- 无缓冲管道是阻塞的，常用于同步通信模式
				- channel的死锁问题
					- Channel满了，就阻塞写；Channel空了，就阻塞读
					- 阻塞之后会交出cpu，去执行其他协程，希望其他协程能帮自己解除阻塞
					- 如果阻塞发生在main协程里，并且没有其他子协程可以执行，那就可以确定"(希望)永远等不来"，就会自已把自己杀掉，再报一个fatal error: deadlock出来
					- 如果阻塞发生在子协程里，就不会发生死锁，因为至少main协程是一个值得等待的"希望"，会一直等下去
		- `make(chan type, len)` 使用 len参数创建指定缓冲区长度的管道
			- 带缓冲管道(buffered channel )
				- 会创建一个环形缓冲队列，队列满时send操作会阻塞或fatal error
				- buffered channel有两个属性: 容量(capacity)和长度(length)
					- capacity: 表示bufffered channel最多可以缓冲多少个数据
					- length: 表示buffered channel当前已缓冲多少个数据
				- 特点
					- 未满之前是非阻塞，异步模式
					- 填满之后是阻塞的，同步模式
```go
	channel = make(chan int)           // %T chan int;    %v 0xc00001e0c0; len 0
	channel2 := make(chan string, 10)  // %T chan string; %v 0xc00005c180; len 0
```

- 读取和写入
	- 通过操作符 `<-` 对管道进行读取和写入操作
		- send: 当ch出现在<-的左边 `ch<-`
		- recv: 当ch出现在<-的右边 `<-ch`
	- 当写入 无缓冲区管道 或 缓冲区管道 已满时，写入则会阻塞，直到管道中元素被其他协程读取
	- 当管道中无元素时，读取也会阻塞，直到管道被其他协程写入元素
		- 只有在协程中读取才会阻塞
		- 在main会直接报错 fatal error，非panic 不能通过recover捕获
	- `channel <- struct{}{}` 起通知作用
		- 空结构体变量的内存占用为0，因此struct{}类型的管道比bool类型的管道还要省内存
	- 只读&&只写
		- 在函数参数时声明管道
			- `chan<-` 表示管道只写
			- `<-chan` 表示管道只读

```go
	// 协程中读取才会阻塞，在main会直接报错 fatal error
	channel2 <- "1"
	fmt.Println(<-channel2)

	// 管道只写
	var wchannel chan<- int
	func Write(cl chan<- rune) { }
	
	/// 管道只读
	var rchannel <-chan int
	func Read(cl <-chan rune) { }
```

- 关闭管道
	- 通过close函数关闭管道
		- 关闭channel后，send操作将导致painc
		- channel不能重复close，否则会panic
		- 关闭channel后，recv操作将返回对应类型的0值以及一个状态码false
		- 使用`close()`时，建议加上defer，只在sender端上显式使用close()关闭channel
	- 当读取到最后一个元素后可通过读取的第二个参数用于判断是否结束
		- `e, ok := <-channel2`
			- `ok == true` 代表管道还没有关闭
			- 读取到最后一个元素返回false
```go
	for {
		if ele, ok := <-channel2; ok {                                  // ok==true代表管道还没有close
			fmt.Printf("receive %d\n", ele)
		} else {                                                        // 管道关闭后，读操作会立即返回"0值"
			fmt.Printf("channel have been closed, receive %d\n", ele)
			break
		}
	}
```


- `for-range` 遍历管道
	- 只有当管道关闭时，才能通过range遍历管道里的数据，否则会发生fatal error
	- 如果使用 channel 或者 sync.WaitGroup 等待协程退出，range未关闭的管道会报错
	- 如果使用 time.Sleep，则到达时间后会直接退出，range未关闭的管道无影响
```go
	channel03 := make(chan int)
	go func() {
		for e := range channel03 {
			fmt.Println(e)
		}
		channel <- 0                     // 利用chan的特性进行阻塞
	}()
	go func() {
		for i := 0; i < 100; i++ {
			channel03 <- i
		}
		close(channel03)
	}()
	<-channel
```

- Go语言time包实现了Tick函数，可以用于实现定时机制，Tick函数返回一个只读管道
	- `func Tick(d Duration) <-chan Time`
```go
	for now := range time.Tick(3 * time.Second) {
		fmt.Println(time.Now())      // 每隔3s打印一次时间
	}
```

## 四、Golang的多路复用(select-case)

### 1. 定义
- I/O模型
	- 操作系统级的I/O模型有
		- 阻塞I/O
		- 非阻塞I/O
		- 信号驱动I/O
		- 异步I/O
		- 多路复用I/O
			- select系统调用可同时监听1024个文件描述符的可读或可写状态
			- poll用链表存储文件描述符，摆脱了1024的上限
			- 各操作系统实现了自己的I/O多路复用函数，如epoll、evport和kqueue 等
			- go语言中多路复用函数以netpoll为前缀，针对不同的操作系统做了不同的封装，以达到最优的性能
				- 在编译go语言时会根据目标平台选择特定的分支进行编译

- 文件描述符
	- Linux下，一切皆文件，包括普通文件、目录文件、字符设备文件(键盘、鼠标)、块设备文件(硬盘、光驱)、套接字socket等等
	- 文件描述符(File descriptor，FD)是访问文件资源的抽象句柄，读写文件都要通过它
	- 文件描述符就是个非负整数，每个进程默认都会打开3个文件描述符: 0标准输入、1标准输出、2标准错误
	- 由于内存限制，文件描述符是有上限的，可通过ulimit –n查看，文件描述符用完后应及时关闭

### 2. 使用
- 若需要同时对多个管道进行监听(写入或读取)，则可以使用`select-case`语句
	- select的行为模式主要是对channel是否可读进行轮询
		- 如果所有的case语句块评估时都被阻塞，则阻塞直到某个语句块可以被处理
		- 如果多个case同时满足条件，则随机选择一个进行处理
			- 对于这一次的选择，其它的case都不会被阻塞，而是处理完被选中的case后进入下一轮select(如果select在循环中)或者结束select(如果select不在循环中或循环次数结束)
		- 如果存在default且其它case都不满足条件，则执行default
			- default必须要可执行而不能阻塞
	- select会被return、break关键字中断: return是退出整个函数，break是退出当前select
	- 当所有case都失败，则执行default语句
	- defalut语句是可选的，不允许fall through行为，但允许case语句块为空块
```go
	select {
	case v, ok := <-channel:
		fmt.Println("channel", v, ok)
	case v, ok := <-channel02:
		fmt.Println("channel02", v, ok)
	default:
		fmt.Println("default")
	}
```

### 3. 超时机制
- 通过select-case 实现对执行操作超时的控制
	- select-case语句监听每个case语句中管道的读取，当某个case语句中管道读取成功则执行对应子语句
```go
	var timeout chan int           
	go func() {
		time.Sleep(3 * time.Second)
		close(timeout)                         // 设定3s后关闭管道
	}()
	select {
	case v, ok := <-channel:
		fmt.Println("success:", r)
	case <-timeout:                            // 3s后关闭管道，读取timeout成功
		fmt.Println("timeout")
	}
```

- Go语言中的标准库 "time" 实现了After函数，可以用于实现超时机制，After函数返回一个只读管道
	- `func After(d Duration) <-chan Time`
```go
	select {
	case v, ok := <-channel:
		fmt.Println("success:", r)
	case <-time.After(3 * time.Second):    // 3s后关闭管道，读取管道成功
		fmt.Println("timeout")
	}
```

- Go语言中的标准库 "context" 实现了timeout
	- 调用cancel()将关闭ctx.Done()对应的管道
	- 调用cancel()或到达超时时间都将关闭ctx.Done()对应的管道

```go
	// context包常用函数
	func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
	func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
	type CancelFunc func()    // CancelFunc会告诉操作放弃它的工作
	type Context interface{ 
		Deadline() (deadline time.Time, ok bool)
		Done() <-chan struct{}
		Err() error
		Value(key interface{}) interface{}
	}
	
	// 关闭管道
	ctx, cancel := context.WithCancel(context.Background())
	
	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background (),time.Microsecond*100)
	ctx.Done()   // 管道关闭后读操作将立即返回
```

## 五、Golang的共享数据

### 1. 定义
- 多个协程对同一个内存资源进行修改，未对资源进行同步限制，导致修改数据混乱

- 临界区
	- 如果一部分程序会被并发访问或修改，为了避免并发访问导致的意向不到的结果，则这部分程序需要被保护起来
	- 如果多个线程同时访问或操作临界区，会造成访问错误
		- 此时需要限定临界区，在同一时间只能有1个线程持有，保证操作的顺序
			- 当临界区由一个线程持有的时候，其它线程如果想进入这个临界区，就会返回失败，或者是等待
			- 直到持有的线程退出临界区，这些等待线程中的某一个才有机会接着持有这个临界区

- 锁
	- 用于解决隔离性的一种机制
	- 某个协程(线程)在访问某个资源时先锁住，防止其它协程的访问，等访问完毕解锁后其他协程再来加锁进行访问
	- Go语言中的锁
		- 互斥锁
			- 互斥锁加锁失败后，线程会释放 CPU，给其他线程
		- 自旋锁
			- 自旋锁加锁失败后，线程会忙等待，直到它拿到锁
		- 读写锁

### 2. sync包
- `sync.Mutex` 互斥锁，用于对资源加锁和释放锁提供对资源同步方式访问
	- 获取到锁的任务，阻塞其他任务; 意味着同一时间只有一个任务可以获取锁
```go
	var HcMutex sync.Mutex
	HcMutex.Lock()    // 获取锁
	HcMutex.UnLock()  // 释放锁
```

- `sync.RWMutex` 读写锁 
	- 写锁阻塞所有锁(所有读锁和写锁)，目的是修改时其他人不要读取，也不要修改
	- 读锁阻塞写锁，读锁可以同时施加多个，目的是不要让修改数据影响读取结果 
		- 同时多个读任务，可以施加多个读锁，阻塞写锁
		- 同时多个写任务，只可以施加一个写锁，阻塞其他所有锁，并且退化成互斥锁
		- 读写混合: 若有写锁，等待释放后能施加 读或写
		- 读写混合: 若有读锁，只能再施加读锁，阻塞写锁
```go
	var rwMutex sync.RWMutex
	rwMutex.Lock      // 获取写入锁
	rwMutex.Unlock    // 释放写入锁 
	rwMutex.RLock     // 获取读取锁
	rwMutex.RUnlock   // 释放读取锁
```

- sync.Map
	- 并发修改map会发生panic，因为map的value是不可寻址的
	- go 1.9引入的内置方法，并发线程安全的map
	- sync.Map 将key和value 按照interface{}存储
	- 查询出来后要类型断言 x.(int) x.(string)
	- 遍历使用Range() 方法，需要传入一个匿名函数作为参数，匿名函数的参数为k,v interface{}，每次调用匿名函数将结果返回
```go
	m := sync.Map{}
	m.Store(k,v)  // 读
	m.Load(k)     // 写
	m.Delete(k)   // 删除
	m.Range(func(k,v interface{}bool{   // 遍历
		k := k.(string)
		v := v.(string)
		log.Printf("[找到了][%s=%d]", key, value)
		return true
	})
	m.LoadOrstore(k,v)   // 若没有key，则添加
	m.LoadAndDelete("key")  // 加载并删除
```

### 3. 原子操作
 - 原子操作是指过程不能中断的操作
 - CAS函数(Compare And Swap)
	- `sync/atomic`包中提供了五类原子操作函数，其操作对象为整数型或整数指针
		- `Add*`: 增加/减少
		- `Load*`: 载入
		- `Store*`: 存储
		- `Swap*`: 更新
		- `CompareAndSwap*`: 比较第一个参数引用值是否与第二个参数值相同，若相同则将第一个参数值更新为第三个参数

## 六、Golang的CSP并发设计模式

### 1. CSP
- Go语言是采用CSP编程思想的典范，它将CSP发挥到了极致，而Goroutine和Channel 就是这种思想的体现

- Go语言的设计者 Rob Pike: *Do not communicate by sharing memory; instead, share memory by communicating.
	- 即 *不要使用共享内存通信，而是应该使用通信去共享内存
		- Thread1 -> Memory -> Thread2
		- Goroutine1 -> Channel -> Goroutine2
		- 使用发送消息来同步信息相比于直接使用共享内存和互斥锁是一种更高级的抽象，使用更高级的抽象能够为程序设计上提供更好的封装，让程序的逻辑更加清晰
		- 消息发送在解耦方面与共享内存相比也有一定优势，将线程的职责分成生产者和消费者，并通过消息传递的方式将它们解耦，不需要再依赖共享内存

### 2. 基于CSP的常见设计模式
- Barrier 模式
	- barrier 屏障模式是一种屏障，用来阻塞直到聚合所有 goroutine 返回结果
	- 使用 channel 来实现
	- 使用场景
		- 多个网络请求并发，聚合结果
		- 粗粒度任务拆分并发执行，聚合结果
	- 网页爬虫

```go
	var (
		client = http.Client{
			Timeout: time.Duration(1 * time.Second),
		}
	)
	type SiteResp struct {
		Err    error
		Resp   string
		Status int
		Cost   int64
	}
	
	// 构造请求
	func doSiteRequest(out chan<- SiteResp, url string) {
		res := SiteResp{}
		startAt := time.Now()
		defer func() {
			res.Cost = time.Since(startAt).Milliseconds()
			out <- res
		}()
		resp, err := client.Get(url)
		if resp != nil {
			res.Status = resp.StatusCode
		}
		if err != nil {
			res.Err = err
			return
		}
		// 暂不处理结果
		_, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			res.Err = err
			return
		}
		// res.Resp = string(byt)
	}
	
	// 聚合返回
	func mergeResponse(resp <-chan SiteResp, ret *[]SiteResp, down chan struct{}) {
		defer func() {
			down <- struct{}{}
		}()
		count := 0
		for v := range resp {
			*ret = append(*ret, v)
			count++
			// 填充完成,  返回
			if count == cap(*ret) {
				return
			}
		}
	}
	
	// 执行请求
	func BatchSiteReqeust() {
		endpoints := []string{
			"https://www.baidu.com",
			"https://segmentfault.com/",
			"https://blog.csdn.net/",
			"https://www.jd.com/",
		}
		// 每个endpoints返回一个结果，缓冲数根据len确定
		respChan := make(chan SiteResp, len(endpoints))
		defer close(respChan)
		// 并行爬取
		for _, endpoints := range endpoints {
			go doSiteRequest(respChan, endpoints)
		}
		// 聚合结果, 返回结束事件, 避免轮询
		down := make(chan struct{})
		ret := make([]SiteResp, 0, len(endpoints))
		go mergeResponse(respChan, &ret, down)
		// 等待结束
		<-down
		// 打印爬取信息
		for _, v := range ret {
			fmt.Println(v)
		}
	}
```

- Pipeline 模式
	- 利用多核的优势把一段粗粒度逻辑分解成多个 goroutine 执行

```go
	// getRandNum | addRandNum | printRes
	var wg sync.WaitGroup
	
	// getRandNum()用于生成随机整数，并将生成的整数放进第一个channel ch1中
	func getRandNum(out chan<- int) {
		defer wg.Done()
		var random int
		// 总共生成10个随机数
		for i := 0; i < 10; i++ {
			// 生成[0,30)之间的随机整数并放进channel out
			random = rand.Intn(30)
			out <- random
		}
		close(out)
	}
	
	// addRandNum()用于接收ch1中的数据(来自第一个函数)，将其输出，然后对接收的值加1后放进第二个channel ch2中
	func addRandNum(in <-chan int, out chan<- int) {
		defer wg.Done()
		for v := range in {
			// 输出从第一个channel中读取到的数据
			// 并将值+1后放进第二个channel中
			fmt.Println("before +1:", v)
			out <- (v + 1)
		}
		close(out)
	}
	
	// printRes() 接收ch2中的数据并将其输出
	func printRes(in <-chan int) {
		defer wg.Done()
		for v := range in {
			fmt.Println("after +1:", v)
		}
	}
	
	// 如果将函数认为是Linux的命令，则类似于下面的命令行: ch1相当于第一个管道，ch2相当于第二个管道
	func PipelineMode() {
		wg.Add(3)
		// 创建两个channel
		ch1 := make(chan int)
		ch2 := make(chan int)
		// 3个goroutine并行
		go getRandNum(ch1)
		go addRandNum(ch1, ch2)
		go printRes(ch2)
		wg.Wait()
	}
```

- Producer/Consumer 模式
	- 生产者消费者模型，该模式主要通过平衡生产线程和消费线程的工作能力来提高程序的整体处理数据的速度
		- 即生产者生产一些数据，然后放到成果队列中，同时消费者从成果队列中来取这些数据
		- 让生产、消费变成了异步的两个过程

```go
	// Producer 生产者: 生成 factor 整数倍的序列
	func Producer(factor int, out chan<- int) {
		maxCount := 0
		for i := 0; ; i++ {
			out <- i * factor
			// 最多生成10个
			maxCount++
			if maxCount > 10 {
				break
			}
		}
	}
	
	// Consumer 消费者
	func Consumer(in <-chan int) {
		for v := range in {
			fmt.Println(v)
		}
	}
	
	func ProducerConsumerMode() {
		ch := make(chan int, 64) // 成果队列
		go Producer(3, ch) // 生成 3 的倍数的序列
		go Producer(5, ch) // 生成 5 的倍数的序列
		go Consumer(ch)    // 消费 生成的队列
		// 运行一定时间后退出
		time.Sleep(5 * time.Second)
	}
```

- Pub/Sub 模式
	- pub/sub 也就是发布订阅模型
		- 在这个模型中，消息生产者成为发布者(publisher)，而消息消费者则成为订阅者(subscriber)，生产者和消费者是M:N的关系
		- 在传统生产者和消费者模型中，是将消息发送到一个队列中，而发布订阅模型则是将消息发布给一个主题

```go
	type (
		subscriber chan interface{}         // 订阅者为一个管道
		topicFunc  func(v interface{}) bool // 订阅者处理消息的函数, bool是方便判断是否处理成功, 这里不作retry实现
	)
	
	// 构建一个发布者对象, 可以设置发布超时时间和缓存队列的长度
	func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
		return &Publisher{
			buffer:      buffer,
			timeout:     publishTimeout,
			subscribers: make(map[subscriber]topicFunc),
		}
	}
	
	// 发布者对象
	type Publisher struct {
		m           sync.RWMutex             // 读写锁
		buffer      int                      // 订阅队列的缓存大小
		timeout     time.Duration            // 发布超时时间
		subscribers map[subscriber]topicFunc // 订阅者信息
	}
	
	// 添加一个新的订阅者，订阅全部主题
	func (p *Publisher) Subscribe() chan interface{} {
		return p.SubscribeTopic(nil)
	}
	
	// 添加一个新的订阅者，订阅过滤器筛选后的主题
	func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
		ch := make(chan interface{}, p.buffer)
	
		p.m.Lock()
		defer p.m.Unlock()
	
		p.subscribers[ch] = topic // channel引用类型，这里的ch是一个内存地址，将这个地址作为key注册到map
	
		return ch
	}
	
	// 退出订阅
	func (p *Publisher) Evict(sub chan interface{}) {
		p.m.Lock()
		defer p.m.Unlock()
	
		delete(p.subscribers, sub)
		close(sub)
	}
	
	// 发布一个主题
	func (p *Publisher) Publish(v interface{}) {
		p.m.RLock()
		defer p.m.RUnlock()
	
		var wg sync.WaitGroup
		for sub, topic := range p.subscribers {
			wg.Add(1)
			go p.sendTopic(sub, topic, v, &wg)
		}
		wg.Wait()
	}
	
	// 发送主题，可以容忍一定的超时
	func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		if topic != nil && !topic(v) {
			return
		}
	
		select {
		case sub <- v:
			fmt.Printf("【消息发送中】%v->%v \n", sub, v)
		case <-time.After(p.timeout):
			fmt.Println("【消息发送超时】")
		}
	}
	
	// 关闭发布者对象，同时关闭所有的订阅者管道
	func (p *Publisher) Close() {
		p.m.Lock()
		defer p.m.Unlock()
	
		for sub := range p.subscribers {
			delete(p.subscribers, sub)
			close(sub)
		}
	}
	
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()
	// 订阅所有
	all := p.Subscribe()
	// 通过过滤订阅一部分信息
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})
	// 发布者 发布信息
	p.Publish("hello,   python!")
	p.Publish("godbybe, python!")
	p.Publish("hello,   golang!")
	// 订阅者查看消息
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()
	// 订阅者查看消息
	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()
	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
```

- Workers Pool 模式
	- Go语言中 goroutine 已经足够轻量，甚至 net/http server 的处理方式也是 goroutine-per-connection 的，所以比起其他语言来说可能场景稍微少一些
		- 每个 goroutine 的初始内存消耗在 2~8kb，当有大批量任务的时候，需要起很多goroutine来处理，这会给系统代理很大的内存开销和GC压力，这个时候就可以考虑一下协程池
	- 参考Go的老版调度实现一个任务工作队列
		- 有(最多)4个worker，每个worker是一个goroutine，它们有worker ID
		- 每个worker都从一个队列中取出待执行的任务(Task)，并发执行
		- 队列容量为10，即最多只允许10个任务进行排队
		- 任务的执行方式: 随机睡眠0-1秒钟，并将任务标记为完成

```go
	// worker的数量，即使用多少goroutine执行任务
	const workerNum = 4
	var wg sync.WaitGroup
	type Task struct {
		ID         int
		JobID      int
		Status     string
		CreateTime time.Time
	}
	
	// Run 执行任务
	func (t *Task) Run() {
		sleep := rand.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		t.Status = "Completed"
	}
	
	// 从buffered channel中读取任务，并执行任务
	func worker(in <-chan *Task, workID int) {
		defer wg.Done()
		for v := range in {
			fmt.Printf("Worker%d: recv a request: TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
			v.Run()
			fmt.Printf("Worker%d: Completed for TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
		}
	}
	
	// 将待执行任务放进buffered channel，共15个任务
	func produceTask(out chan<- *Task) {
		for i := 1; i <= 15; i++ {
			// fmt.Println(i)
			out <- &Task{
				ID:         i,
				JobID:      100 + i,
				CreateTime: time.Now(),
			}
		}
	}
	
	func RunTaskWithPool() {
		wg.Add(workerNum)
		// 创建容量为10的bufferd channel
		taskQueue := make(chan *Task, 10)
		// 激活goroutine，执行任务
		for workID := 1; workID <= workerNum; workID++ {
			go worker(taskQueue, workID)
		}
		// 生成消息
		produceTask(taskQueue)
		// 5秒后 关闭管道，通知所有worker退出
		time.Sleep(5 * time.Second)
		close(taskQueue)
		wg.Wait()
	}
```

## 七、Golang的并发注意事项
- 协程泄漏
	- 原因
		- 协程阻塞，未能如期结束
		- 协程阻塞最常见的原因都跟channel有关
		- 由于每个协程都要占用内存，所以协程泄漏也会导致内存泄漏
	- 排查
		- go run对应程序
		- 在浏览器访问 `http://127.0.0.1:8080/debug/pprof/goroutine?debug=1 `
		- 在终端执行 `go tool pprof http://0.0.0.0:8080/debug/pprof/goroutine`
			- 通过list查看函数每行代码产生了多少协程 `list + 对应函数`
		- 通过traces打印调用堆栈 `traces`
			- 在pprof中输入web命令 web
		- 终端执行 `go tool pprof --http=:8081 /Users/zhangchaoyang/pprof/pprof.goroutine.001.pb.gz`
			- 在source view下可看到哪行代码生成的协程最多
```go
	// 示例
	import (
		"net/http"
		_ "net/http/pprof"
	)
	func main() {
		go func() {
			if err := http.ListenAndServe("localhost:8080", nil); err != nil {
				panic(err)
			}
		}()
	}
```

- 控制并发数
	- Goroutine会以一个很小的栈启动(可能是2KB或4KB)，但是系统和调度器的能力总是有上限的，在面对大规模的并发请求时(千万或者亿)是要考虑goroutine的销毁成本的
	- 方法
		- 使用goroutine pool控制gotourine数量
		- 做好系统的限流与上限控制
		- 管理好goroutine的退出，不让goroutine泄露

- 并发的安全退出
	- 有时候需要通知goroutine停止它正在干的事情，特别是当它工作在错误的方向上的时候
```go
	func worker(wg *sync.WaitGroup, cancel chan bool) {
		defer wg.Done()
		for {
			select {
			default:
				fmt.Println("hello")
				time.Sleep(100 * time.Millisecond)
			case <-cancel:
				return
			}
		}
	}
	func CancelWithDown() {
		cancel := make(chan bool)
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go workerv2(&wg, cancel)
		}
		time.Sleep(time.Second)
		// 发送退出信号
		close(cancel)
		// 等待goroutine 安全退出
		wg.Wait()
	}
```

- 使用回调函数替代channel
	- 回调函数就是一个被作为参数传递的函数
	- 使用回调函数进行网页爬虫
```go
	var (
		client = http.Client{
			Timeout: time.Duration(1 * time.Second),
		}
	)
	type SiteResp struct {
		Err    error
		Resp   string
		Status int
		Cost   int64
	}
	type SiteRespCallBack func(SiteResp)
	
	// 构造请求
	func doSiteRequest(cb SiteRespCallBack, url string) {
		res := SiteResp{}
		startAt := time.Now()
		defer func() {
			res.Cost = time.Since(startAt).Milliseconds()
			cb(res)
			wg.Done()
		}()
	
		resp, err := client.Get(url)
		if resp != nil {
			res.Status = resp.StatusCode
		}
		if err != nil {
			res.Err = err
			return
		}
	
		// 暂不处理结果
		_, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			res.Err = err
			return
		}
	
		// res.Resp = string(byt)
	}
	
	// 主函数完成回调处理逻辑
	func CallBackMode() {
		endpoints := []string{
			"https://www.baidu.com",
			"https://segmentfault.com/",
			"https://blog.csdn.net/",
			"https://www.jd.com/",
	}
	
	// 一个endpoints返回一个结果, 缓冲可以确定
	respChan := make(chan SiteResp, len(endpoints))
	defer close(respChan)
	
	// 回调处理逻辑
	ret := make([]SiteResp, 0, len(endpoints))
	cb := func(resp SiteResp) {
		ret = append(ret, resp)
	}
	
	// 并行爬取
	for _, endpoints := range endpoints {
		wg.Add(1)
		go doSiteRequest(cb, endpoints)
	}
	
	// 等待结束
	wg.Wait()
	
	for _, v := range ret {
		fmt.Println(v)
	}
```

- 使用守护进程优雅退出
	- 守护协程
		- 独立于控制终端和用户请求的协程，它一直存在，周期性执行某种任务或等待处理某些发生的事件
		- 伴随着main协程的退出，守护协程也退出
	- kill命令不是杀死进程，它只是向进程发送信号`kill –s pid`，s的默认值是15
	- 常见的终止信号
		- `SIGINT   2`   Ctrl + c触发
		- `SIGKILL  9`   无条件结束程序，不能捕获、阻塞或忽略
		- `SIGTERM  15`  结束程序，可以捕获、阻塞或忽略
			- 当Context的deadline到期或调用了CancelFunc后，Context的Done()管道会关闭，该管道上关联的读操作会解除阻塞，然后执行协程退出前的清理工作
```go
	type Context interface {
		Deadline() (deadline time.Time, ok bool)
		Done() <-chan struct{}
	}l
	func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

- 协程管理组件 `github.com/x-mod/routine`
	- 封装了常规的业务逻辑: 初始化、收尾清理、工作协程、守护协程、监听term信号
	- 封装了常见的协程组织形式: 并行、串行、定时任务、超时控制、重试、profiling

## 八、Golang的sync包
- sync包提供了同步原语

- 常用结构体
	- `sync.Mutex`:  互斥锁
	- `sync.RWMutex`: 读写锁
	- `sync.Cond`: 条件等待
	- `sync.Once`: 单次执行，且多个Once也只执行一个，可以确保在高并发的场景下有些事情只执行一次，比如加载配置文件、关闭管道等
	- `sync.Map`: 协程安全映射
	- `sync.Pool`: 对象池
	- `sync.WaitGroup`: 组等待
	
## 九、Golang的runtime包
- runtime包提供了与 Go运行时系统交互的操作

- 常用函数
	- `runtime.Gosched()`: 当前 goroutine 让出时间片
	- `runtime.GOROOT()`: 获取 Go 安装路径
	- `runtime.NumCPU()`: 获取可使用的逻辑 CPU 数量
	- `runtime.GOMAXPROCS(1)`: 设置当前进程可使用的逻辑 CPU 数量
	- `runtime.NumGoroutine()`: 获取当前进程中 goroutine 的数量