# Golang-Time  Golang的时间

## 一、Golang的time标准库
- Go 语言通过标准库 time 包处理日期和时间相关的问题
 
### 1. 打印当前时间
- 方法: `time.Now()`
	- 时间戳返回 int64
	- 10位数时间戳 是秒单位
	- 13位数时间戳 是毫秒单位，毫秒=纳秒/1e6 且prometheus默认查询就是毫秒
	- 19位数时间戳 是纳秒单位

```go
	func numLen(n int64) int {
		return len(strconv.Itoa(int(n)))
	}
	now := time.Now()
	log.Printf("[当前时间对象为: %v]", now)
	log.Printf("[当前时间戳 秒级: %v][位数: %v]", now.Unix(), numLen(now.Unix()))
	log.Printf("[当前时间戳 毫秒级: %v][位数:%v]", now.UnixNano()/1e6, numLen(now.UnixNano()/1e6))
	log.Printf("[当前时间戳 纳秒级: %v][位数:%v]", now.UnixNano(), numLen(now.UnixNano()))
	log.Printf("[当前时间戳 纳秒小数部分: %v]", now.Nanosecond())
	log.Printf("[当前时间 %v %v %v %v %v %v",now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	log.Printf("[今天是 %d年中的第 %d天 星期 %d]", now.Year(), now.YearDay(), now.Weekday())
	zone, offset := now.Zone() && log.Printf("[直接获取时区  %v，和东utc时区差 %d 个小时]", zone, offset/3600)
	
	/* 输出
	2021/07/17 16:34:10 [当前时间对象为: 2021-07-17 16:34:10.2999431 +0800 CST m=+0.007295301]
	2021/07/17 16:34:10 [当前时间戳 秒级: 1626510850][位数:10]
	2021/07/17 16:34:10 [当前时间戳 毫秒级: 1626510850299][位数:13]
	2021/07/17 16:34:10 [当前时间戳 纳秒级: 1626510850299943100][位数:19]
	2021/07/17 16:34:10 [当前时间戳 纳秒小数部分: 299943100]
	2021/07/17 16:34:10 [当前时间 2021 June 17 16 34 10]
	2021/07/17 16:44:57 [今天是 2021年中的第 198天 星期]
	2021/07/17 16:44:57 [直接获取时区  CST，和东utc时区差 8 个小时]
	*/
```


### 2. Time结构体
- Now()返回的是个Time结构体，这也是Go内部表示时间的数据结构
	- Time 代表一个纳秒精度的时间点
	- 程序中应使用 Time 类型值来保存和传递时间，而不是指针；就是说，表示时间的变量和字段，应为 time.Time 类型，而不是 *time.Time. 类型
	- 时间点可以使用 Before、After 和 Equal 方法进行比较
		- `Sub` 方法让两个时间点相减，生成一个 Duration 类型值(代表时间段)
		- `Add` 方法给一个时间点加上一个时间段，生成一个新的 Time 类型时间点
	- Time 零值代表时间点 January 1, year 1, 00:00:00.000000000 UTC
		- 因为本时间点一般不会出现在使用中，IsZero 方法提供了检验时间是否是显式初始化的一个简单途径
	- Time是有时区的s通过 == 比较 Time 时，Location 信息也会参与比较，因此 Time 不应该作为 map 的 key
```
	type Time struct {
		// Has unexported fields.
	}
```

- now() 的具体实现在 runtime 包中，由汇编实现的，和平台有关，一般在`sys_{os_platform}_amd64.s` 中

### 3. 时间的格式化
- 方法: time.Now().Format()

- 格式
```go
	/*
	模板   占位
	 年  ->  2006
	 月  ->  01
	 日  ->  02
	 时  ->  03(12h) / 15(24h)
	 分  ->  04
	 秒  ->  05
	*/
	format := "2006-01-02 15:04:05"                                // string
	fmt.Printf("%T %#v\n", now.Format(format), now.Format(format)) // string "2021-06-20 20:08:45"
```

### 4. 时间戳转换
- Unix时间戳(Unix timestamp)定义为从1970年01月01日00时00分00秒(UTC)起至现在经过的总秒数

- 不论东西南北、在地球的每一个角落都是相同

- Go语言中的转换
```
	              -> time.Unix(sec int64, nsec int64)              -> time.Format()
	时间戳(Timestamp)                             time.Time                             日期格式
	              <- time.Unix()                                   <- time.Parse()				
```

- 时间戳(Timestamp) 转成 time.Time 类型 再格式化成日期格式
```go
	// Time -> Timestamp  方法: time.Now().Unix()
	ts := time.Now().Unix()             // 时间戳
	layout := "2006-01-02 15:04:05"
	
	// Timestamp -> Time  方法: time.Unix()	
	t := time.Unix(ts, 0)               // 构造时间对象
	log.Printf(t.Format(layout))
```

- 日期格式字符串 转成 time.Time 类型 再转成时间戳
```go
	// time.Parse()                 func Parse(layout, value string) (Time, error)  返回转换后的时间格式和判断信息(err) 
	d1, err := time.Parse("2006-01-02 15:04:05", "2021-06-18 12:12:12")
	log.Println(d1.Unix())
	
	// time.ParseInLocation()       func ParseInLocation(layout, value string, loc *Location) (Time, error) 可以指定时区
	tStr := "2021-07-17 16:52:59"
	layout := "2006-01-02 15:04:05"
	t1, _ := time.ParseInLocation(layout, tStr, time.Local)
	t2, _ := time.ParseInLocation(layout, tStr, time.UTC)
	log.Printf("[ %s的 CST时区的时间戳为 : %d]", tStr, t1.Unix())
	log.Printf("[ %s的 UTC时区的时间戳为 : %d]", tStr, t2.Unix())
	log.Printf("[UTC - CST =%d 小时]", (t2.Unix()-t1.Unix())/3600)
```

### 5. 时间的比较
- Before、After 和 Equal
```go
	now := time.Now()
	t1, _ := time.ParseDuration("1h")
	m1 := now.Add(t1)
	log.Printf("[a.after(b) a在b之后: %v]", m1.After(now))      // func After(d Duration) <-chan Time
	log.Printf("[a.Before(b) a在b之前: %v]", now.Before(m1))    // func (t Time) Before(u Time) bool
	log.Printf("[a.Equal(b) a=b: %v]", m1.Equal(now))           // func (t Time) Equal(u Time) bool
```
	
### 6. 时间长度 Duration
- `time.Duration`表示时间长度
	- 以纳秒为基数
	- 底层数据类型为int64
	- int64 类型的变量不能直接和time.Duration类型相乘，需要显示转换，常量除外
		- 不行:  `num * time.Second`
		- 可以:  `time.Duration(num) * time.Second`
		- 可以:  `5 * time.Second`

### 7. 时长计算
- Add
	- 让一个时间点加上一个时间段，生成一个新的 Time 类型时间点
```go
	// func (t Time) Add(d Duration) Time
	now := time.Now()
	after := now.Add(time.Hour * 24)
	fmt.Println(after)             // 2021/07/19 17:33:21 2021-07-20 17:33:21.12884928 +0800 CST m=+86400.000024700
```

- Sub
	- 让两个时间点相减，生成一个 Duration 类型值(代表时间段)
```go
	// func (t Time) Sub(u Time) Duration
	fmt.Println(now.Sub(after))    // 2021/07/19 17:33:21 -24h0m0s
```

- ParseDuration 时间差
```go
	// func ParseDuration(s string) (Duration, error)
	now := time.Now()
	var layout = "2006-01-02 15:04:05"
	func tTostr(t time.Time) string {
		return time.Unix(t.Unix(), 0).Format(layout)
	}

	t1, _ := time.ParseDuration("1h1m1s")  // 1小时1分1秒后
	m1 := now.Add(t1)
	log.Printf("[ 1小时1分1秒后时间为: %v]", tTostr(m1))
	
	t2, _ := time.ParseDuration("-1h1m1s") // 1小时1分1秒前
	m2 := now.Add(t2)
	log.Printf("[ 1小时1分1秒前时间为: %v]", tTostr(m2))

	sub1 := now.Sub(m2)                    // sub计算两个时间差
	log.Printf("[ 时间差: %s 、相差小时数: %v 、相差分钟数: %v ]", sub1.String(), sub1.Hours(), sub1.Minutes())

	t3, _ := time.ParseDuration("-3h3m3s")
	m3 := now.Add(t3)
	log.Printf("[time.since 当前时间与t的时间差: %v]", time.Since(m3))    // func Since(t Time) Duration
	log.Printf("[time.until t与当前时间的时间差: %v]", time.Until(m3))    // func Until(t Time) Duration
	m4 := now.AddDate(0, 0, 5)                                            // func (t Time) AddDate(years int, months int, days int) Time
	log.Printf("[5天后的时间: %v]", m4)
```

### 8. Sleep
- 方法: `time.Sleep()`
```go
	time.Sleep(time.Hour * 24)
```

### 9. 定时器
- 定时器是进程规划自己在未来某一时刻接获通知的一种机制，共有2种

- 单次触发: `Timer`
	- 通过 `time.After` 实现同步等待
	- 通过 `time.AfterFunc` 中断循环，触发自定义函数
		- Timer的stop
			- 如果定时器还未触发，Stop 会将其移除，并返回 true；否则返回 false
			- 后续再对该 Timer 调用 Stop，直接返回 false
		- Timer的Reset
			- Reset 先调用 stopTimer 再调用 startTimer
			- 类似于废弃之前的定时器，重新启动一个定时器
			- 返回值和 Stop 一样
```go
	// Timer数据结构
	type Timer struct {
		C <-chan Time   // C: 一个存放Time对象的Channel
		r runtimeTimer  // runtimeTimer: 它定义在 sleep.go 文件中，必须和 runtime 包中 time.go 文件中的 timer 必须保持一致
	}

	// 通过 time.After 实现同步等待
	m := time.NewTimer(5 * time.Second)
	fmt.Println(<-m.C)
	fmt.Println("exit")

	// 通过 time.AfterFunc 中断循环，触发自定义函数
	stop := false
	time.AfterFunc(5*time.Second, func() {          // func AfterFunc(d Duration, f func()) *Timer
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


- 周期性触发: `Ticker`
```go
	// Ticker数据结构	
	type Ticker struct {
		C <-chan Time // The channel on which the ticks are delivered.
		r runtimeTimer
	}

	// time.NewTicker 实现同步等待
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