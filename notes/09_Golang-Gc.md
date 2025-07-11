# Golang-Gc  Golang的Gc

## 一、经典的GC算法

### 1、引用计数(reference counting)
- 它是最简单的一种垃圾回收算法，和之前提到的智能指针异曲同工

- 实现方式
	- 对每个对象维护一个引用计数，当引用该对象的对象被销毁或更新时被引用对象的引用计数自动减一，当被引用对象被创建或被赋值给其他对象时引用计数自动加一
	- 当引用计数为0时则立即回收对象
	- 这种方法的优点是实现简单，并且内存的回收很及时

- 这种算法在内存比较紧张和实时性比较高的系统中使用的比较广泛，如ios cocoa框架，php，python等

- 简单引用计数算法也有明显的缺点
	- 频繁更新引用计数降低了性能
		- 一种简单的解决方法就是编译器将相邻的引用计数更新操作合并到一次更新
		- 还有一种方法是针对频繁发生的临时变量引用不进行计数，而是在引用达到0时通过扫描堆栈确认是否还有临时对象引用而决定是否释放
	- 循环引用问题
		- 当对象间发生循环引用时引用链中的对象都无法得到释放
		- 最明显的解决办法是避免产生循环引用，如cocoa引入了strong指针和weak指针两种指针类型
		- 或者系统检测循环引用并主动打破循环链，但这也增加了垃圾回收的复杂度

### 2、标记-清除(mark&sweep)
- 实现方式
	- 标记(Mark phase)，从根变量开始迭代得所有被引用的对象，对能够通过应用遍历访问到的对象都进行标记为"被引用" 
	- 清除(Sweep phase)，标记完成后进行清除操作，对没有标记过的内存进行回收(回收同时可能伴有碎片整理操作)

- 这种方法解决了引用计数的不足，但是也有比较明显的问题
	- STW(stop the world)，让程序暂停，程序出现卡顿
	- 标记需要扫描整个heap
	- 清除数据会产生heap碎片
	- 每次启动垃圾回收都会暂停当前所有的正常代码执行，回收使系统响应能力大大降低
	- 后续也出现了很多mark&sweep算法的变种(如三色标记法)优化了这个问题

### 3、分代收集(generation)
- 经过大量实际观察得知，在面向对象编程语言中，绝大多数对象的生命周期都非常短

- 实现方式
	- 分代收集的基本思想是，将堆划分为两个或多个称为 代(generation)的空间
	- 新创建的对象存放在称为 新生代(young generation)中(一般来说，新生代的大小会比 老年代小很多)
	- 随着垃圾回收的重复执行，生命周期较长的对象会被 提升(promotion)到老年代中
	- 因此，新生代垃圾回收和老年代垃圾回收两种不同的垃圾回收方式应运而生，分别用于对各自空间中的对象执行垃圾回收
	- 新生代垃圾回收的速度非常快，比老年代快几个数量级，即使新生代垃圾回收的频率更高，执行效率也仍然比老年代垃圾回收强，这是因为大多数对象的生命周期都很短，根本无需提升到老年代

### 4、复制收集(CopyandCollection)


## 二、Golang中的Gc算法

### 1、历史回顾
- go语言垃圾回收总体采用的是 标记-清除(mark&sweep) 算法

- v1.1 STW

- v1.3 Mark STW，Sweep 并行

- v1.5 三色标记法

- v1.8 hybrid write barrier 混合写屏障 + 三色标记法

- golang 的 GC 不断的提升并发性能并且**减少STW (Stop The World)**的时间

### 2、v1.1 标记清除算法
- 步骤
	- stw 暂停程序
	- mark 标记可达对象(分类出可对象和不可达对象)
	- sweep 清除不可达对象
	- stw 结束

### 3、v1.3 标记清除算法的优化
- 优化过程
	- 标记清除算法整个 GC 都是在 STW，STW 的时间过长
	- 为了减少 STW 的时间，golang 将 STW 的范围缩小
	- 只在 mark 时 STW，sweep 和程序并发执行

- 步骤
	- stw 暂停程序
	- mark 标记可达对象(分类出可对象和不可达对象)
	- stw 结束
	- 继续处理程序 同时并发 sweep 清除不可达对象

- 缺点
	- 整个的 STW 还是很长，特别是当需要 mark 标记的对象越多，需要 STW 的时间越长

### 4、v1.5 三色标记法
- 三色标记法
	- 动态的滚动式的做 gc，将一次长的 stw 分散到多次短暂的 stw 中去

- 步骤
	- 遍历根对象的第一层可达对象标记为灰色，不可达默认白色
	- 将灰色对象的下一层可达对象标记为灰色，自身标记为黑色
	- 多次重复步骤2，直到灰色对象为0，只剩下白色对象和黑色对象
	- sweep 白色对象

- 三色标记对象丢失
	- 如果 gc 期间不 stw 的话有可能对象丢失
		- 即 一个黑色对象在 gc 期间链接了白色对象，白色对象又没有任何的灰色对象可达就会导致对象的丢失
	- 解决思路
		- 对象丢失是在GC里面绝对不被允许的，可以暂时存在垃圾但是不能丢失
		- 所以还是要 STW，为了解决这个 STW 的问题，引入了强三色不变式和弱三色不变式，只要满足强三色不变式和弱三色不变式的任意一种就能解决 STW 的问题
			- 强三色不变式: 一个黑色对象在 gc 期间链接了白色对象
			- 弱三色不变式: 黑色对象链接的白色对象又没有任何的灰色对象可达

### 5、v1.8 三色标记法 + 混合写屏障
- "强-弱" 三色不变式
	- 强 三色不变式，不存在黑色对象引用到白色对象的指针
	- 弱 三色不变式，所以被黑色对象引用的白色对象都处于灰色保护状态

- 插入写屏障
	- 对象被引用时，触发的屏障，相当于拦截器
	- 具体操作: 在A对象引用B对象的时候，B对象被标记为灰色
	- 满足方式: 强三色不变式
	- 缺点，结束时需要STW来重新扫描栈，标记栈上引用的白色对象的存活

- 删除写屏障
	- 具体操作: 被删除的对象，如果自身为灰色或者白色，那么被标记为灰色
	- 满足方式: 弱三色不变式 (保护灰色对象到白色对象的路径不会断)
	- 缺点，回收精度低，GC开始时STW扫描堆栈来记录初始快照，这个过程会保护开始时刻的所有存活对象

- 混合写屏障
	- 具体操作
		- GC开始将栈上的对象全部引用扫描并标记为黑色 (之后不再进行第二次重复扫描，无需STW)
		- GC期间，任何在栈上创建的新对象，均为黑色，保证栈全为黑色对象
		- 被添加的对象标记为灰色，插入写屏障
		- 被删除的对象标记为灰色，删除写屏障
	- 满足: 变形的弱三色不变式

- 解决对象丢失问题
	- 一个黑色对象在 gc 期间链接了白色对象
	- 链接对象会经过混合写屏障，新插入白色的元素会标记为灰色
	- 接下来按照三色标记法继续遍历即可
	- [三色标记法和混合写屏障](https://blog.csdn.net/weixin_35655661/article/details/119509548)

### 6、Golang-GC && Java-GC 的区别
- golang 的 gc 的不断演化通过将 STW 分散化大大减少 STW 的时间

- 最新版本的 golang 的 gc 基本都是在 1ms 以内，但是这样故意的设计是牺牲吞吐量换来的

- gc 的频率会比其他的垃圾回收更高比如 jvm 和 .net高

- go 目前的混合写屏障 + 三色标记标记在 gc 期间的屏障也会给性能带来一定的损耗，但是总体的损耗目前的机器配置完全能 cover
	- go 的开发效率很高
	- go 的主要应用领域不是 cpu 敏感型，在高 IO 的领域有很好的性能
	- go 要避免跟 jvm 和 .net 这样给高吞吐的 gc 正面竞争，才能发挥出 go 超短 stw 的优势

### 7、参考资料
- [Golang GC 垃圾回收机制详解](https://blog.csdn.net/u010649766/article/details/80582153)
- [图解Golang的GC算法](https://studygolang.com/articles/18850?fr=sidebar)
- [golang 垃圾回收 gc 详解](https://blog.csdn.net/jarvan5/article/details/122970491)


## 三、Golang中的逃逸分析

### 1、逃逸分析
- 逃逸分析(escape analysis)就是在程序编译阶段根据程序代码中的数据流，对代码中哪些变量需要在栈上分配，哪些变量需要在堆上分配进行静态分析的方法
	- 帮助程序员将那些人们认为需要分配在栈上的变量尽可能保留在栈上，尽可能少的"逃逸"到堆上的算法

- 栈
	- 在程序中，每个函数块都会有自己的内存区域用来存自己的局部变量(内存占用少)、返回地址、返回值之类的数据，这一块内存区域有特定的结构和寻址方式，大小在编译时已经确定，寻址起来也十分迅速，开销很少
	- 这一块内存地址称为栈，栈是线程级别的，大小在创建的时候已经确定，所以当数据太大的时候，就会发生 "stack overflow"

- 堆
	- 在程序中，全局变量、内存占用大的局部变量、发生了逃逸的局部变量存在的地方就是堆，这一块内存没有特定的结构，也没有固定的大小，可以根据需要进行调整
	- 简单来说，有大量数据要存的时候，就存在堆里面
	- 堆是进程级别的，当一个变量需要分配在堆上的时候，开销会比较大，对于 go 这种带 GC 的语言来说，也会增加 gc 压力，同时也容易造成内存碎片

- 进程的内存结构
```
	0xc0000000  内核虚拟内存          <-- 内核使用
	0x40000000      栈区             <-- 程序运行时用于存放局部变量，可向下延伸空间
	            共享库的内存映像
	                堆区             <-- 程序运行时用于分配mallco和new申请的区域
	            可读写区(.data .bss) <-- 存放全局变量和静态变量
	0x08048000     只读区            <-- 存放程序和常量等
	        0      未使用
```

- golang中的 逃逸分析
	- 所谓逃逸分析是指由编译器决定内存分配的位置，不需要程序员指定
	- 函数中申请一个新的对象
		- 如果函数外部没有引用，则优先放到栈中
		- 如果函数外部存在引用，则必定放到堆中
	- go 在一定程度消除了堆和栈的区别，因为 go 在编译的时候进行逃逸分析，来决定一个对象放栈上还是放堆上，不逃逸的对象放栈上，可能逃逸的放堆上

- golang中的 内存逃逸
	- golang的变量会携带一组校验数据，用来证明它的整个生命周期是否是在运行；如果通过校验会在栈上分配内存，否则会在堆上分配内存，也就是逃逸
	- 关键点: "栈分配廉价, 堆分配昂贵"

- 一般而言，遇到以下情况会发生逃逸行为，Go编译器会将变量存储在heap上
	- 函数内局部变量在函数外部被引用
	- 接口(interface)类型的变量
	- 在interface上调用方法
	- size未知或者动态变化的变量，如slice，map，channel，[]byte等
	- 在slice中保存指针或带有指针的值
	- slice背后的数组被重新分配了
	- 发送指针或带有指针的值到channel中: 不能确定哪个goroutine会在channel上接收值，在编译时无法确认变量何时释放

- 内存逃逸分析工具
	- 内存逃逸分析是编译器在编译期就完成的，可以使用以编译下命令来做内存逃逸分析
		- `go build -gcflags="-m"` 可以展示逃逸分析、内联优化等各种优化结果
		- `go build -gcflags="-m -l"` -l 会禁用内联优化，可以过滤掉内联优化的结果展示，关注逃逸分析的结果
		- `go build -gcflags="-m -m"` -m 会展示更详细的分析结果

- 逃逸总结
	- [golang_逃逸分析](https://blog.csdn.net/qq_51537858/article/details/128239205)
	- [golang逃逸分析](https://www.cnblogs.com/xuweiqiang/p/16388143.html)
	- [Golang逃逸分析](https://cloud.tencent.com/developer/beta/article/2090795)
	- 栈上分配内存比在堆中分配内存有更高的效率
	- 栈上分配的内存不需要GC处理
	- 堆上分配的内存使用完毕会交给GC处理
	- 逃逸分析目的是决定内分配地址是栈还是堆
	- 逃逸分析在编译阶段完成


## 四、Golang中的内存泄漏

### 1、内存泄漏

- 内存泄漏，指程序运行过程中，内存因为某些原因无法释放或没有释放，让内存资源造成了浪费
	- 如果泄漏的内存越堆越多，就会占用程序正常运行的内存，比较轻的影响是程序开始运行越来越缓慢；
	- 严重的话，可能导致大量泄漏的内存堆积，最终导致程序没有内存可以运行，最终导致 OOM (Out Of Memory，即内存溢出)

- 造成内存泄露的场景
	- [浅谈Golang内存泄漏](https://cloud.tencent.com/developer/article/2134737)
	- [以kubelet为例使用go tool pprof分析服务性能](https://blog.csdn.net/buppt/article/details/127505818)
	- [strace的简单用法](https://blog.csdn.net/mijichui2153/article/details/85229307)
	- slice造成内存泄漏
		- 在slice1基础上，通过[1:3]切分出slice2；当 slice1 GC之后，底层数组其它没有被引用的位置 会产生泄露
		- 解决方法，使用 append 或者 copy
	- time.Ticker造成内存泄漏
		- 定时器未调用stop，导致发生内存泄漏
	- goroutine阻塞，造成内存泄漏
		- 向满的channel发送消息
		- 从空的channel接收消息
		- 当channel没有初始化的时候就会处于nil状态，向nil的channel发送或接收
		- 解决方法
			- 发生泄漏前，发送者和接收者的数量需要一致、channel需要初始化
			- 发生泄漏后，采用 `go tool pprof` 分析内存的占用和变化

- 案例：使用`go tool pprof`分析kubelet服务性能
	- CPU高，查看 profile、goroutine
	- Memory高，查看 heap、goroutine
	- 怀疑流程卡主，查看 profile、goroutine
		- 先确保没有 D 状态进程/线程
		- 随手执行下 `df > /dev/null`，查看 `echo $?`
		- 随手检查下 runtime
			- docker `docker ps -aq | xargs -n1 timeout 5 docker inspect > /dev/null`
			- containerd `crictl -r /var/run/containerd/containerd.sock ps -aq | xargs -n1 crictl -r /var/run/containerd/containerd.sock inspect > /dev/null`
```shell
# 通过 kubectl proxy 代理出 api，一般都是 8001 端口
$ kubectl proxy --address='0.0.0.0' --accept-hosts='^*$'

# 捞下 profile
$ curl 127.0.0.1:8001/api/v1/nodes/{nodename}/proxy/debug/pprof/profile -oprofile.bin

# 捞下 goroutine 信息
$ curl 127.0.0.1:8001/api/v1/nodes/{nodename}/proxy/debug/pprof/goroutine -ogoroutine.bin

# 捞下 heap 信息
$ curl 127.0.0.1:8001/api/v1/nodes/{nodename}/proxy/debug/pprof/heap -oheap.bin

# 分析pprof
# 上面命令直接使用： go tool pprof + api地址也可以查看
$ go tool pprof profile.bin
File: kubelet
Build ID: f822c0d1dfe86d6162b91cc75032946f93ace0a1
Type: cpu
Time: Apr 11, 2024 at 3:06pm (CST)
Duration: 30.12s, Total samples = 199.70s (663.04%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 156.52s, 78.38% of 199.70s total
Dropped 881 nodes (cum <= 1s)
Showing top 10 nodes out of 108
      flat  flat%   sum%        cum   cum%
   112.36s 56.26% 56.26%    112.97s 56.57%  syscall.Syscall
     8.17s  4.09% 60.36%     11.73s  5.87%  strings.Fields
     6.71s  3.36% 63.72%     14.20s  7.11%  runtime.scanobject
     6.11s  3.06% 66.78%      6.11s  3.06%  runtime.futex
     5.11s  2.56% 69.33%      5.11s  2.56%  memeqbody
     4.58s  2.29% 71.63%      4.61s  2.31%  path/filepath.(*lazybuf).append (inline)
     3.86s  1.93% 73.56%      9.93s  4.97%  path/filepath.Clean
     3.85s  1.93% 75.49%      3.85s  1.93%  runtime.memclrNoHeapPointers
     3.43s  1.72% 77.21%      4.16s  2.08%  runtime.findObject
     2.34s  1.17% 78.38%      2.34s  1.17%  runtime.memmove
(pprof) list syscall.Syscall
Total: 199.70s
ROUTINE ======================== syscall.Syscall in /usr/local/go/src/syscall/asm_linux_amd64.s
   112.36s    112.97s (flat, cum) 56.57% of Total
         .          .     13:
         .          .     14:// func rawVforkSyscall(trap, a1, a2 uintptr) (r1, err uintptr)
         .          .     15:TEXT ·rawVforkSyscall(SB),NOSPLIT|NOFRAME,$0-40
         .          .     16:	MOVQ	a1+8(FP), DI
         .          .     17:	MOVQ	a2+16(FP), SI
         .      350ms     18:	MOVQ	$0, DX
         .          .     19:	MOVQ	$0, R10
         .          .     20:	MOVQ	$0, R8
         .          .     21:	MOVQ	$0, R9
         .          .     22:	MOVQ	trap+0(FP), AX	// syscall entry
      40ms       40ms     23:	POPQ	R12 // preserve return address
   112.27s    112.27s     24:	SYSCALL
         .          .     25:	PUSHQ	R12
         .          .     26:	CMPQ	AX, $0xfffffffffffff001
         .          .     27:	JLS	ok2
         .          .     28:	MOVQ	$-1, r1+24(FP)
         .          .     29:	NEGQ	AX
         .          .     30:	MOVQ	AX, err+32(FP)
         .          .     31:	RET
         .          .     32:ok2:
      30ms       30ms     33:	MOVQ	AX, r1+24(FP)
      20ms       20ms     34:	MOVQ	$0, err+32(FP)
         .          .     35:	RET
         .      260ms     36:
         .          .     37:// func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
         .          .     38:TEXT ·rawSyscallNoError(SB),NOSPLIT,$0-48
         .          .     39:	MOVQ	a1+8(FP), DI
         .          .     40:	MOVQ	a2+16(FP), SI
         .          .     41:	MOVQ	a3+24(FP), DX
ROUTINE ======================== syscall.Syscall6 in /usr/local/go/src/syscall/asm_linux_amd64.s
     590ms      640ms (flat, cum)  0.32% of Total
         .          .     36:
         .          .     37:// func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
         .          .     38:TEXT ·rawSyscallNoError(SB),NOSPLIT,$0-48
         .          .     39:	MOVQ	a1+8(FP), DI
         .          .     40:	MOVQ	a2+16(FP), SI
         .       50ms     41:	MOVQ	a3+24(FP), DX
         .          .     42:	MOVQ	trap+0(FP), AX	// syscall entry
         .          .     43:	SYSCALL
         .          .     44:	MOVQ	AX, r1+32(FP)
         .          .     45:	MOVQ	DX, r2+40(FP)
         .          .     46:	RET
         .          .     47:
         .          .     48:// func gettimeofday(tv *Timeval) (err uintptr)
         .          .     49:TEXT ·gettimeofday(SB),NOSPLIT,$0-16
     590ms      590ms     50:	MOVQ	tv+0(FP), DI
         .          .     51:	MOVQ	$0, SI
         .          .     52:	MOVQ	runtime·vdsoGettimeofdaySym(SB), AX
         .          .     53:	TESTQ   AX, AX
         .          .     54:	JZ fallback
         .          .     55:	CALL	AX
(pprof)

$ go tool pprof goroutine.bin
File: kubelet
Build ID: f822c0d1dfe86d6162b91cc75032946f93ace0a1
Type: goroutine
Time: Apr 11, 2024 at 2:59pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 979, 99.39% of 985 total
Dropped 143 nodes (cum <= 4)
Showing top 10 nodes out of 136
      flat  flat%   sum%        cum   cum%
       954 96.85% 96.85%        954 96.85%  runtime.gopark
        25  2.54% 99.39%         25  2.54%  syscall.Syscall
         0     0% 99.39%          7  0.71%  bufio.(*Reader).Peek
         0     0% 99.39%         28  2.84%  bufio.(*Reader).Read
         0     0% 99.39%          8  0.81%  bufio.(*Reader).fill
         0     0% 99.39%         28  2.84%  bytes.(*Buffer).ReadFrom
         0     0% 99.39%          5  0.51%  crypto/tls.(*Conn).Read
         0     0% 99.39%          5  0.51%  crypto/tls.(*Conn).readFromUntil
         0     0% 99.39%          5  0.51%  crypto/tls.(*Conn).readRecord (inline)
         0     0% 99.39%          5  0.51%  crypto/tls.(*Conn).readRecordOrCCS
(pprof) list runtime.gopark
Total: 985
ROUTINE ======================== runtime.gopark in /usr/local/go/src/runtime/proc.go
       954        954 (flat, cum) 96.85% of Total
         .          .    301:		if forcegc.idle.Load() {
         .          .    302:			throw("forcegc: phase error")
         .          .    303:		}
         .          .    304:		forcegc.idle.Store(true)
         .          .    305:		goparkunlock(&forcegc.lock, waitReasonForceGCIdle, traceEvGoBlock, 1)
       954        954    306:		// this goroutine is explicitly resumed by sysmon
         .          .    307:		if debug.gctrace > 0 {
         .          .    308:			println("GC forced")
         .          .    309:		}
         .          .    310:		// Time-triggered, fully concurrent.
         .          .    311:		gcStart(gcTrigger{kind: gcTriggerTime, now: nanotime()})
ROUTINE ======================== runtime.goparkunlock in /usr/local/go/src/runtime/proc.go
         0        123 (flat, cum) 12.49% of Total
         .          .    307:		if debug.gctrace > 0 {
         .          .    308:			println("GC forced")
         .          .    309:		}
         .          .    310:		// Time-triggered, fully concurrent.
         .          .    311:		gcStart(gcTrigger{kind: gcTriggerTime, now: nanotime()})
         .        123    312:	}
         .          .    313:}
         .          .    314:
         .          .    315://go:nosplit
         .          .    316:
         .          .    317:// Gosched yields the processor, allowing other goroutines to run. It does not
(pprof)
```