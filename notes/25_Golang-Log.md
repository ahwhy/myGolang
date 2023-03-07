# Golang-Log  Golang的日志

## 一、Golang的标准库 log包

- log包
	- 实现了简单的日志服务
	- 内置日志库，提供基本的日志功能，但是没有提供日志级别，比如: debug、warning、error
	- 定义了Logger类型，该类型提供了一些格式化输出的方法
	- Fatal系列函数会在写入日志信息后调用 `os.Exit(1)`
	- Panic系列函数会在写入日志信息后 `panic`

- log/syslog
	- syslog 包提供一个简单的系统日志服务的接口

### 1. Constants
```go
	// itoa枚举依次是 1，2，4，8，16，32
	const (
		// 字位共同控制输出日志信息的细节。不能控制输出的顺序和格式。
		// 在所有项目后会有一个冒号：2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
		Ldate         = 1 << iota     // 日期：2009/01/23
		Ltime                         // 时间：01:23:23
		Lmicroseconds                 // 微秒分辨率：01:23:23.123123（用于增强Ltime位）
		Llongfile                     // 文件全路径名+行号： /a/b/c/d.go:23
		Lshortfile                    // 文件无路径名+行号：d.go:23（会覆盖掉Llongfile）
		LstdFlags     = Ldate | Ltime // 标准logger的初始值
	)
	
	// log.flag && iota
	// 因为可以组由组合标志位，后端进行 &判断
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {}  // 代表原来的flag中有Ldate|Ltime|Lmicroseconds
	// logger不能决定字段输出顺序  a|b = b|a
```

### 2. 基本使用
- 常用函数
```go
	// SetFlags 设置标准logger的输出选项
	func SetFlags(flag int)
	// Flags 返回标准logger的输出选项
	func Flags() int

	// SetPrefix 设置标准logger的输出前缀
	func SetPrefix(prefix string)
	// Prefix 返回标准logger的输出前缀
	func Prefix() string

	// SetOutput 设置标准logger的输出目的地，默认是标准错误输出
	func SetOutput(w io.Writer)

	// 调用Output将生成的格式化字符串输出到标准logger
	func Printf(format string, v ...interface{})
	func Print(v ...interface{})
	func Println(v ...interface{})

	// 等价于{Printf(v...); os.Exit(1)}
	func Fatalf(format string, v ...interface{})
	func Fatal(v ...interface{})
	func Fatalln(v ...interface{})

	// 等价于{Printf(v...); panic(...)}
	func Panicf(format string, v ...interface{})
	func Panic(v ...interface{})
	func Panicln(v ...interface{})
```

- 参考示例
	- 输出文件，只要实现接口 `io.Writer` 的类型都可以作为文件的输出
```go
		func logPrint(baseStr string) {
			for i := 0; i < 10; i++ {
				msg := fmt.Sprintf("%s_%d", baseStr, i)
				log.Println(msg)
			}
		}
		
		// 创建文件对象
		file, err := os.OpenFile("my.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)       // func Fatal(v ...interface{}) - Fatal is equivalent to Print() followed by a call to os.Exit(1).
		}
		
		// 设置log输出到文件
		log.SetOutput(file)      // func SetOutput(w io.Writer)
		go logPrint("aaaaa")
		time.Sleep(1 * time.Second)
```

### 2. Logger的使用
- Logger
	- log包提供了一个预定义的“标准”Logger，可以通过辅助函数`Print[f|ln]`、`Fatal[f|ln]`和`Panic[f|ln]`访问，比手工创建一个Logger对象更容易使用
	- Logger会打印每条日志信息的日期、时间，默认输出到标准错误
```go
	// Logger类型表示一个活动状态的记录日志的对象，它会生成一行行的输出写入一个io.Writer接口
	// 每一条日志操作会调用一次io.Writer接口的Write方法
	// Logger类型的对象可以被多个线程安全的同时使用，它会保证对io.Writer接口的顺序访问
	type Logger struct { ... }

	// New创建一个Logger
	// 参数out设置日志信息写入的目的地
	// 参数prefix会添加到生成的每一条日志前面
	// 参数flag定义日志的属性(时间、文件等等)
	func New(out io.Writer, prefix string, flag int) *Logger

	// SetFlags 设置logger的输出选项
	func (l *Logger) SetFlags(flag int)
	// Flags 返回logger的输出选项
	func (l *Logger) Flags() int

	// SetPrefix 设置logger的输出前缀
	func (l *Logger) SetPrefix(prefix string)
	// Prefix 返回logger的输出前缀
	func (l *Logger) Prefix() string

	// Output写入输出一次日志事件
	// 参数s包含在Logger根据选项生成的前缀之后要打印的文本
	func (l *Logger) Output(calldepth int, s string) error

	// 调用l.Output将生成的格式化字符串输出到logger
	func (l *Logger) Printf(format string, v ...interface{})
	func (l *Logger) Print(v ...interface{})
	func (l *Logger) Println(v ...interface{})

	// 等价于{l.Printf(v...); os.Exit(1)}
	func (l *Logger) Fatalf(format string, v ...interface{})
	func (l *Logger) Fatal(v ...interface{})
	func (l *Logger) Fatalln(v ...interface{})

	// 等价于{l.Printf(v...); panic(...)}
	func (l *Logger) Panicf(format string, v ...interface{})
	func (l *Logger) Panic(v ...interface{})
	func (l *Logger) Panicln(v ...interface{})
```

- 参考示例
```go
	var (
		WarningLogger *log.Logger
		InfoLogger    *log.Logger
		ErrorLogger   *log.Logger
	)
	
	func init() {
		file, err := os.OpenFile("c.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
				log.Fatal(err)
		}
		InfoLogger = log.New(file, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
		WarningLogger = log.New(file, "[WARNING]", log.LstdFlags|log.Lshortfile)
		ErrorLogger = log.New(file, "[ERROR]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
	
    // 创建文件对象
    InfoLogger.Println("[常见写法]启动服务....")
    InfoLogger.Println("[日期简写]正常上报....")
    WarningLogger.Println("[文件长路径]不严重的错误，报个warining....")
    ErrorLogger.Println("[微秒时间戳]严重的错误，报个error....")
```


## 二、常用的公共log包

### 1. github.com/sirupsen/logrus
- 实现利用logrus包，通过钉钉机器人发送日志
	- 首先定义相关结构体，然后实现`Levels`和`Fire`方法 --> 实现 Hook接口
		- `Levels`中定义日志等级
		- `Fire`中处理日志发送逻辑
			- 比如发送到redis、es、钉钉、logstash
	- 调用 `AddHook()`，直接打印日志并发送
```go
	// ogrus源码
	type Hook interface {
		Levels() []Level
		Fire(*Entry) error
	}
	func AddHook(hook Hook) {
		std.AddHook(hook)
	}
```

### 2. github.com/rifflock/lfshook
- 结合`logrotate github.com/lestrrat-go/file-rotatelogs`
	- 保留4个文件 `rotatelogs.WithRotationCount(4)`
	- 切割时间 `rotatelogs.WithRotationTime(1*time.Second)`
	- 删除时间 `rotatelogs.WithMaxAge(2*time.Minute)`
