# 日期与时间

Go 语言通过标准库 time 包处理日期和时间相关的问题

## 简单示例

打印当前时间
```go
now := time.Now()
fmt.Println(now) // 2021-06-17 13:29:40.801445 +0800 CST m=+0.001636524
fmt.Println(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) // 2021 June 17 13 29 40
```

## Time结构体

我们看到Now()返回的是个Time结构体, 这也是Go内部表示时间的数据结构
```go
type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
```

+ Time 代表一个纳秒精度的时间点
+ 程序中应使用 Time 类型值来保存和传递时间，而不是指针。就是说，表示时间的变量和字段，应为 time.Time 类型，而不是 *time.Time. 类型
+ 时间点可以使用 Before、After 和 Equal 方法进行比较
+ Sub 方法让两个时间点相减，生成一个 Duration 类型值（代表时间段）
+ Add 方法给一个时间点加上一个时间段，生成一个新的 Time 类型时间点
+ Time 零值代表时间点 January 1, year 1, 00:00:00.000000000 UTC。因为本时间点一般不会出现在使用中，IsZero 方法提供了检验时间是否是显式初始化的一个简单途径
+ Time是有时区的
+ 通过 == 比较 Time 时，Location 信息也会参与比较，因此 Time 不应该作为 map 的 key

now() 的具体实现在 runtime 包中, 由汇编实现的, 和平台有关, 一般在sys_{os_platform}_amd64.s 中

## 时间的格式化

方法：time.Now().Format()


## 时间的解析

方法：time.Parse()，返回转换后的时间格式和一个判断信息（err)

```go
d1, err := time.Parse("2006-01-02 15:04:05", "2021-06-18 12:12:12")
fmt.Println(d1, err)
```

为什么是 2006-01-02 15:04:05
这是固定写法，类似于其他语言中 Y-m-d H:i:s 等。为什么采用这种形式, 最直接的说法是很好记：2006 年 1 月 2 日 3 点 4 分 5 秒

## 时间戳转换

Time -> Timestamp 方法：time.Now().Unix()

```go
ts := time.Now().Unix()
fmt.Println(ts)
```

Timestamp -> Time 方法: 

```go
ts := time.Now().Unix()
fmt.Println(ts)

fmt.Println(time.Unix(ts, 0))
```


## 时间的比较

使用 Before、After 和 Equal 方法进行比较

```go
d1, err := time.Parse("2006-01-02 15:04:05", "2021-06-18 12:12:12")
fmt.Println(d1, err)

now := time.Now()
fmt.Println(now.After(d1))
fmt.Println(now.Before(d1))
fmt.Println(now.Equal(d1))
```

## 时间长度: Duration

```go
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64

const (
	minDuration Duration = -1 << 63
	maxDuration Duration = 1<<63 - 1
)

// Common durations. There is no definition for units of Day or larger
// to avoid confusion across daylight savings time zone transitions.
//
// To count the number of units in a Duration, divide:
//	second := time.Second
//	fmt.Print(int64(second/time.Millisecond)) // prints 1000
//
// To convert an integer number of units to a Duration, multiply:
//	seconds := 10
//	fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
//
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)
```

time.Duration表示时间长度
+ 以纳秒为基数
+ 底层数据类型为int64

int64 类型的变量不能直接和time.Duration类型相乘，需要显示转换，常量除外
+ 不行：num * time.Second
+ 可以： time.Duration(num) * time.Second
+ 可以： 5 * time.Second


## 时长计算

1.Add: Add 方法给一个时间点加上一个时间段，生成一个新的 Time 类型时间点

```go
now := time.Now()
after1 := now.Add(time.Hour * 24)
fmt.Println(after1)
```



2.Sub: 方法让两个时间点相减，生成一个 Duration 类型值（代表时间段）
```go
now := time.Now()
after1 := now.Add(time.Hour * 24)
fmt.Println(now.Sub(after1))
```

## Sleep

```go
```

## 定时器

定时器是进程规划自己在未来某一时刻接获通知的一种机制。定时器有2种:
+ 单次触发: Timer
+ 周期性触发: Ticker

### Timer

注意：Timer 的实例必须通过 NewTimer 或 AfterFunc 获得, 我们先看2个简单的用法:

1.通过 time.AfterFunc中断循环， 到时触发自定义函数

```go
stop := false
time.AfterFunc(5*time.Second, func() {
	stop = true
})

for {
	if stop {
		fmt.Println("exit")
		break
	}

	time.Sleep(1 * time.Second)
}
```


2.通过time.After实现同步等待
```go
m := time.NewTimer(5 * time.Second)
fmt.Println(<-m.C)
fmt.Println("exit")
```

3.timer的stop

如果定时器还未触发，Stop 会将其移除，并返回 true；否则返回 false；后续再对该 Timer 调用 Stop，直接返回 false。

```go
```


4.timer的Reset

Reset 会先调用 stopTimer 再调用 startTimer，类似于废弃之前的定时器，重新启动一个定时器。返回值和 Stop 一样

```go
```

5.timer数据结构

```go
// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	C <-chan Time
	r runtimeTimer
}
```
+ C: 一个存放Time对象的Channel
+ runtimeTimer: 它定义在 sleep.go 文件中，必须和 runtime 包中 time.go 文件中的 timer 必须保持一致

### Ticker

```go
tk := time.NewTicker(2 * time.Second)

count := 1
for {
	if count > 2 {
		tk.Stop()
		break
	}
	fmt.Println(<-tk.C)
	count++
}
```

我们看看Ticker结构体
```go
// A Ticker holds a channel that delivers ``ticks'' of a clock
// at intervals.
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	r runtimeTimer
}
```