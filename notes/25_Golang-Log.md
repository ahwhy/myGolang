# Golang-Log  Golang的日志收集

## 一、标准库 log 包
- Go 语言可以通过标准库 log 包处理日志
	- 内置日志库
	- 提供基本的日志功能，但是没有提供日志级别，比如: debug、warning、error

- 简单使用
	- 输出文件，只要实现接口io.Writer的类型都可以作为文件的输出
```go
		func logPrinta(baseStr string) {
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
		time.Sleep(1 * time.Hour)
```

- 自定义的logger
```go
	var (
		WarningLogger *log.Logger          // type Logger struct { // Has unexported fields. }
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

- log.flag && iota
```go
	// itoa枚举依次是 1，2，4，8，16，32
	const (
		Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
		Ltime                         // the time in the local time zone: 01:23:23
		Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
		Llongfile                     // full file name and line number: /a/b/c/d.go:23
		Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
		LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
		Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
		LstdFlags     = Ldate | Ltime // initial values for the standard logger
	)
	
	// 因为可以组由组合标志位，后端进行 &判断
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {}  // 代表原来的flag中有Ldate|Ltime|Lmicroseconds
	// logger不能决定字段输出顺序  a|b = b|a
```

## 二、其他的log包

### 1. github.com/sirupsen/logrus
- 实现利用logrus包，通过钉钉机器人发送日志
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
	- 首先定义相关结构体，然后实现Levels和Fire方法 --> 实现 Hook接口
		- Levels中定义日志等级
		- Fire中处理日志发送逻辑
			- 比如发送到redis、es、钉钉、logstash
	- 调用AddHook()，直接打印日志并发送

### 2. github.com/rifflock/lfshook
- 结合`logrotate github.com/lestrrat-go/file-rotatelogs`
	- 保留4个文件 `rotatelogs.WithRotationCount(4)`
	- 切割时间 `rotatelogs.WithRotationTime(1*time.Second)`
	- 删除时间 `rotatelogs.WithMaxAge(2*time.Minute)`

