# Golang-Test  Golang的测试

## 一、Go语言中的测试

### 1、单元测试
- Go提供了test工具用于代码的测试
	- test工具会查找包下以 `_test.go`结尾的文件
	- 调用测试文件中以Test或Benchmark开头的函数并给出运行结果

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

- 参考示例
	- GoConvey 测试框架 `go get github.com/smartystreets/goconvey`
	- testify 测试框架 `go get github.com/stretchr/testify/assert`
	- [单元测试覆盖率应用实例](https://github.com/m3db/m3/pull/3525)

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

- 执行命令
	- go test 会在运行基准测试之前执行包里所有的单元测试
		- 如果包里有很多单元测试，或者它们会运行很长时间
		- 通过 go test 的-run 标识排除这些单元测试
		- `go test -bench=. -run=none`
	- 内存消耗情况
		- `go test -bench=. -benchmem -run=none`
	- CPU消耗情况
		- `go test -bench=. -cpu=1,2,4 -benchmem -run=none`
	- count多次运行基准测试
		- `go test -bench=. -count=10 -benchmem -run=none      // 热缩放、内存局部性、后台处理、gc活动等等会导致单次的误差`
	- benchtime指定运行秒数 
		- `go test -bench=. -benchtime=5s -benchmem -run=none  // 有的函数比较慢，为了更精确的结果，可以通过 -benchtime 标志指定运行时间，从而使它运行更多次`


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