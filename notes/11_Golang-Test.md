# Golang-Test  Golang的测试

## 一、Go语言中的测试

### 1、test工具
- Go提供了test工具用于代码的测试
	- test工具会查找包下以 `_test.go`结尾的文件
	- 调用测试文件中以 `TestXxx` 或 `BenchmarkXxx` 开头的函数并给出运行结果
	- 这类文件将被排除在正常的程序包之外，但在运行 `go test` 命令时将被包含
	- 有关详细信息，请运行 `go help test` 和 `go help testflag` 了解

- testing
	- testing包 提供对 Go 包的自动化测试的支持
	- 通过 `go test` 命令，能够自动执行 `func TestXxx(*testing.T)` 形式的任何函数
		- 其中 Xxx 可以是任何字母数字字符串(但第一个字母不能是 [a-z])，用于识别测试例程
	- 在这些函数中，使用 `Error`, `Fail` 或相关方法来发出失败信号
	- 如果有需要，可以调用 `*T` 和 `*B` 的 Skip 方法，跳过该测试或基准测试
```go
	func TestTimeConsuming(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}
		...
	}
```

- testing/iotest
	- iotest包 实现用于测试的读取器和写入器

- testing/quick
	- quick包 实现一些实用函数来帮助黑盒测试

### 2、单元测试
- 定义
	- 单元是应用中最小的可测试部件，如函数和对象的方法
	- 单元测试是软件开发中对最小单位进行正确性检验的测试工作
	- Go语言原生支持了单元测试，使用上非常简单

- 意义
	- 保证变更/重构的正确性，特别是在一些频繁变动和多人合作开发的项目中
	- 简化调试过程: 可以轻松的测试出哪一部分代码出现问题
	- 单元测试是最好的文档 -> 在单测中直接给出具体接口的使用方法，是最好的实例代码

- 单元测试用例编写的原则
	- 单一原则: 一个测试用例只负责一个场景
	- 原子性: 结果只有两种情况: PASS/FAIL
	- 优先要核心组件和逻辑的测试用例
	- 对于高频使用库，如: util，需要重点覆盖

- 使用要求
	- 文件名命名要求: `{pkg_name}_test.go`
	- 需要导入testing包
	- 测试方法/函数命名要求: 定义以Test开头的函数
	- 函数参数: `(t *testing.T)`
	- 测试文件和被测试文件必须在一个包中
	- 在测试函数中调用函数进行返回值测试，当测试失败可通过`testing.T` 结构体的`Error*`函数抛出错误
	- 使用 `package packagename_test` 命名的测试包在构建二进制文件时不会被打包

- Go语言的测试框架
	- testing 测试框架
		- 执行命令
			- `go test -v .`
			- -cover测试覆盖率 `go test -v -cover .`
			- 只执行某个函数 `go test -run=TestAdd -v .`
			- 正则过滤函数名 `go test -run=TestM.* -v .`
			- 指定func/sub跑子测试 `go test -run=TestMul/正数 -v`
		- 子测试 t.Run
		- table-driven tests
			- 所有用例的数据组织在切片 cases 中，看起来就像一张表，借助循环创建子测试
			- 益处
				- 新增用例非常简单，只需给 cases 新增一条测试数据即可
				- 测试代码可读性好，直观地能够看到每个子测试的参数和期待的返回值
				- 用例失败时，报错信息的格式比较统一，测试报告易于阅读
				- 如果数据量较大，或是一些二进制数据，推荐使用相对路径从文件中读取
			- [prometheus的api_test](https://github.com/prometheus/prometheus/blob/main/web/api/v1/api_test.go)
```go
	// AllocsPerRun 返回在调用 f 期间内存平均分配次数。虽然返回值的类型为 float64，但它始终是一个整数值
	func AllocsPerRun(runs int, f func()) (avg float64)

	// PB 被 RunParallel 使用来运行并行基准测试
	type PB struct { ... }
	// Next 判断是否有更多的迭代要执行
	func (pb *PB) Next() bool

	// T 是传递给测试函数的一种类型，它用于管理测试状态并支持格式化测试日志；测试日志会在执行测试的过程中不断累积， 并在测试完成时转储至标准输出
	// 当一个测试的测试函数返回时， 又或者当一个测试函数调用 FailNow 、 Fatal 、 Fatalf 、 SkipNow 、 Skip 或者 Skipf 中的任意一个时， 该测试即宣告结束
	type T struct { ... }

	// 调用 Error 相当于在调用 Log 之后调用 Fail
	func (c *T) Error(args ...interface{})
	// 调用 Errorf 相当于在调用 Logf 之后调用 Fail
	func (c *T) Errorf(format string, args ...interface{})
	// 将当前测试标识为失败，但是仍继续执行该测试
	func (c *T) Fail()
	// 将当前测试标识为失败并停止执行该测试，在此之后，测试过程将在下一个测试或者下一个基准测试中继续
	func (c *T) FailNow()
	// Failed 用于报告测试函数是否已失败
	func (c *T) Failed() bool
	// 调用 Fatal 相当于在调用 Log 之后调用 FailNow
	func (c *T) Fatal(args ...interface{})
	// Log 使用与 Println 相同的格式化语法对它的参数进行格式化，然后将格式化后的文本记录到错误日志里面
	// 1）对于测试来说，格式化文本只会在测试失败或者设置了 -test.v 标志的情况下被打印出来
	// 2）对于基准测试来说，为了避免 -test.v 标志的值对测试的性能产生影响， 格式化文本总会被打印出来
	func (c *T) Log(args ...interface{})
	// Log 使用与 Printf 相同的格式化语法对它的参数进行格式化，然后将格式化后的文本记录到错误日志里面
	// 如果输入的格式化文本最末尾没有出现新行，那么将一个新行添加到格式化后的文本末尾
	func (c *T) Logf(format string, args ...interface{})
	// 返回正在运行的测试或基准测试的名字
	func (c *T) Name() string
	// Parallel 用于表示当前测试只会与其他带有 Parallel 方法的测试并行进行测试
	func (t *T) Parallel()
	// 执行名字为 name 的子测试 f，并报告 f 在执行过程中是否出现了任何失败；Run 将一直阻塞直到 f 的所有并行测试执行完毕
	func (t *T) Run(name string, f func(t *T)) bool
	// 调用 Skip 相当于在调用 Log 之后调用 SkipNow
	func (c *T) Skip(args ...interface{})
	// 将当前测试标识为“被跳过”并停止执行该测试
	func (c *T) SkipNow()
	// 调用 Skipf 相当于在调用 Logf 之后调用 SkipNow
	func (c *T) Skipf(format string, args ...interface{})
	// Skipped 用于报告测试函数是否已被跳过
	func (c *T) Skipped() bool
```

- 参考示例
	- GoConvey 测试框架 `go get github.com/smartystreets/goconvey`
	- testify 测试框架 `go get github.com/stretchr/testify/assert`
	- [单元测试覆盖率应用实例](https://github.com/m3db/m3/pull/3525)
```go
	// 子测试 t.Run
	func TestMul(t *testing.T) {
		t.Run("正数", func(t *testing.T) {
			if Mul(4, 5) != 20 {
				t.Fatal("muli.zhengshu.error")
			}
		})
	
		t.Run("负数", func(t *testing.T) {
			if Mul(2, -3) != -6 {
				t.Fatal("muli.fusshu.error")
			}
		})
	}
```

### 3、基准测试
- 用途
	- 基准测试常用于代码性能测试

- 用法
	- 函数需要导入testing包
	- 定义以 Benchmark开头的函数，参数为`(b *testing.B)`
	- 在测试函数中 循环调用函数多次

- bench的工作原理
	- 基准测试函数会被一直调用直到b.N无效，它是基准测试循环的次数
	- b.N 从 1 开始，如果基准测试函数在1秒内就完成(默认值)，则 b.N 增加，并再次运行基准测试函数
	- b.N 的值会按照序列 1,2,5,10,20,50,... 增加，同时再次运行基准测测试函数
	- 输出: `BenchmarkFib-12    183    6272054 ns/op    0 B/op    0 allocs/op`
		- 上述结果解读代表 1秒内运行了183次 每次 6272054 ns
		- -12 后缀和用于运行次测试的 GOMAXPROCS 值有关
		- 与GOMAXPROCS一样，此数字默认为启动时 Go
```go
	// B 是传递给基准测试函数的一种类型，它用于管理基准测试的计时行为，并指示应该迭代地运行测试多少次
	// 一个基准测试在它的基准测试函数返回时，又或者在它的基准测试函数调用 FailNow、Fatal、Fatalf、SkipNow、Skip 或者 Skipf 中的任意一个方法时，测试即宣告结束；至于其他报告方法，比如 Log 和 Error 的变种，则可以在其他 goroutine 中同时进行调用
	// 跟单元测试一样，基准测试会在执行的过程中积累日志，并在测试完毕时将日志转储到标准错误；但跟单元测试不一样的是，为了避免基准测试的结果受到日志打印操作的影响，基准测试总是会把日志打印出来
	type B struct {
	N int
	...
	}

	// 调用 Error 相当于在调用 Log 之后调用 Fail
	func (c *B) Error(args ...interface{})
	// 调用 Errorf 相当于在调用 Logf 之后调用 Fail
	func (c *B) Errorf(format string, args ...interface{})
	// 将当前的测试函数标识为“失败”，但仍然继续执行该函数
	func (c *B) Fail()
	// 将当前的测试函数标识为“失败”，并停止执行该函数；在此之后，测试过程将在下一个测试或者下一个基准测试中继续
	func (c *B) FailNow()
	// Failed 用于报告测试函数是否已失败
	func (c *B) Failed() bool
	// 调用 Fatal 相当于在调用 Log 之后调用 FailNow
	func (c *B) Fatal(args ...interface{})
	// Log 使用与 Println 相同的格式化语法对它的参数进行格式化，然后将格式化后的文本记录到错误日志里面
	func (c *B) Log(args ...interface{})
	// Log 使用与 Printf 相同的格式化语法对它的参数进行格式化， 然后将格式化后的文本记录到错误日志里面
	func (c *B) Logf(format string, args ...interface{})
	// 返回正在运行的测试或者基准测试的名字
	func (c *B) Name() string
	// 打开当前基准测试的内存统计功能，与使用 -test.benchmem 设置类似，但 ReportAllocs 只影响那些调用了该函数的基准测试
	func (b *B) ReportAllocs()
	// 对已经逝去的基准测试时间以及内存分配计数器进行清零
	func (b *B) ResetTimer()
	// 执行名字为 name 的子基准测试（subbenchmark）f ，并报告 f 在执行过程中是否出现了任何失败
	func (b *B) Run(name string, f func(b *B)) bool
	// 以并行的方式执行给定的基准测试
	func (b *B) RunParallel(body func(*PB))
	// 记录在单个操作中处理的字节数量
	func (b *B) SetBytes(n int64)
	// 将 RunParallel 使用的 goroutine 数量设置为 p*GOMAXPROCS ，如果 p 小于 1 ，那么调用将不产生任何效果
	func (b *B) SetParallelism(p int)
	// 调用 Skip 相当于在调用 Log 之后调用 SkipNow
	func (c *B) Skip(args ...interface{})
	// 将当前测试标识为“被跳过”并停止执行该测试
	func (c *B) SkipNow()
	// 调用 Skipf 相当于在调用 Logf 之后调用 SkipNow
	func (c *B) Skipf(format string, args ...interface{})
	// 报告测试是否已被跳过
	func (c *B) Skipped() bool
	// 开始对测试进行计时
	func (b *B) StartTimer()
	// 停止对测试进行计时
	func (b *B) StopTimer()

	// 基准测试运行的结果
	type BenchmarkResult struct {
	N         int           // The number of iterations.
	T         time.Duration // The total time taken.
	Bytes     int64         // Bytes processed in one iteration.
	MemAllocs uint64        // The total number of memory allocations.
	MemBytes  uint64        // The total number of bytes allocated.
	}

	// 测试单个函数
	func Benchmark(f func(b *B)) BenchmarkResult
	func (r BenchmarkResult) AllocedBytesPerOp() int64
	func (r BenchmarkResult) AllocsPerOp() int64
	func (r BenchmarkResult) MemString() string
	func (r BenchmarkResult) NsPerOp() int64
	func (r BenchmarkResult) String() string

	// TB 是一个接口，类型 T 和 B 实现了该接口
	type TB interface {
		Error(args ...interface{})
		Errorf(format string, args ...interface{})
		Fail()
		FailNow()
		Failed() bool
		Fatal(args ...interface{})
		Fatalf(format string, args ...interface{})
		Log(args ...interface{})
		Logf(format string, args ...interface{})
		Name() string
		Skip(args ...interface{})
		SkipNow()
		Skipf(format string, args ...interface{})
		Skipped() bool
		// contains filtered or unexported methods
	}
```

- 执行命令
	- go test 会在运行基准测试之前执行包里所有的单元测试
		- 通过 `go test` 命令，加上 `-bench flag` 来执行，多个基准测试按照顺序运行
		- 如果包里有很多单元测试，或者它们会运行很长时间，通过 go test -run 标识排除这些单元测试 `go test -bench=. -run=none`
	- 内存消耗情况
		- `go test -bench=. -benchmem -run=none`
	- CPU消耗情况
		- `go test -bench=. -cpu=1,2,4 -benchmem -run=none`
	- count多次运行基准测试
		- `go test -bench=. -count=10 -benchmem -run=none      // 热缩放、内存局部性、后台处理、gc活动等等会导致单次的误差`
	- benchtime指定运行秒数 
		- `go test -bench=. -benchtime=5s -benchmem -run=none  // 有的函数比较慢，为了更精确的结果，可以通过 -benchtime 标志指定运行时间，从而使它运行更多次`
```go
	// For example
	func BenchmarkHello(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fmt.Sprintf("hello")
		}
	}
	// 基准函数会运行目标代码 b.N 次。在基准执行期间，会调整 b.N 直到基准测试函数持续足够长的时间
	// Output: BenchmarkHello    10000000    282 ns/op
	// 意味着循环执行了 10000000 次，每次循环花费 282 纳秒(ns)

	// 如果在运行前基准测试需要一些耗时的配置，则可以通过 b.ResetTimer() 先重置定时器
	func BenchmarkBigLen(b *testing.B) {
		big := NewBig()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			big.Len()
		}
	}

	// 如果基准测试需要在并行设置中测试性能，则可以使用 RunParallel 辅助函数; 这样的基准测试一般与 go test -cpu 标志一起使用
	func BenchmarkTemplateParallel(b *testing.B) {
		templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			for pb.Next() {
				buf.Reset()
				templ.Execute(&buf, "World")
			}
		})
	}
```

- 斐波那契数列
```go
	func fib(n int) int {                   // fib.go
		if n == 0 || n == 1 {
			return n
		}
		return fib(n-2) + fib(n-1)
	}
	func BenchmarkFib(b *testing.B) {       // fib_test.go
		for n := 0; n < b.N; n++ {
			fib(30)
		}
	}
	// ResetTimer 如果基准测试在循环前需要一些耗时的配置，则可以先重置定时器
	func BenchmarkFib(b *testing.B) {
		time.Sleep(3 * time.Second)
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			fib(30)
		}
	}
```

- benchmem 展示内存消耗情况
	- G
	- 测试大cap的切片，直接用cap初始化，然后动态扩容
	- 结论: 用cap初始化好的性能可以高一个数据量级
```go
	// 制定大的cap的切片
	func generateWithCap(n int) []int {
		rand.Seed(time.Now().UnixNano())
		nums := make([]int, 0, n)
		for i := 0; i < n; i++ {
			nums = append(nums, rand.Int())
		}
		return nums
	}
	
	// 动态扩容的slice
	func generateDynamic(n int) []int {
		rand.Seed(time.Now().UnixNano())
		nums := make([]int, 0)
		for i := 0; i < n; i++ {
			nums = append(nums, rand.Int())
		}
		return nums
	}
	func BenchmarkGenerateWithCap(b *testing.B) {
		for n := 0; n < b.N; n++ {
			generateWithCap(100000)
		}
	}
	func BenchmarkGenerateDynamic(b *testing.B) {
		for n := 0; n < b.N; n++ {
			generateDynamic(100000)
		}
	}
```

- 测试函数复杂度 不带cap的slice 动态扩容
	- 结论: 输入变为原来的10倍，单次耗时也差不多是上一级的10倍，说明这个函数的复杂度是线性的
	- string拼接的 bench
		- `const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"`
		- `"+"`
		- `[]byte`
		- `strings.Builder`
		- `bytes.Buffer`
```go
	func benchmarkGenerate(i int, b *testing.B) {
		for n := 0; n < b.N; n++ {
			generateDynamic(i)
		}
	}
	func BenchmarkGenerateDynamic1000(b *testing.B)     { benchmarkGenerate(1000, b) }
	func BenchmarkGenerateDynamic10000(b *testing.B)    { benchmarkGenerate(10000, b) }
	func BenchmarkGenerateDynamic100000(b *testing.B)   { benchmarkGenerate(100000, b) }
	func BenchmarkGenerateDynamic1000000(b *testing.B)  { benchmarkGenerate(1000000, b) }
	func BenchmarkGenerateDynamic10000000(b *testing.B) { benchmarkGenerate(10000000, b) }
```