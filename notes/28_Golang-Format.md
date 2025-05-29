# Golang-Format  Golang的格式化

## 一、Golang的标准库 fmt包

- fmt
    - fmt包实现了类似C语言printf和scanf的格式化I/O
    - 格式化动作('verb')源自C语言但更简单

### 1. verb
- 类型
```
	%v : 变量的自然形式(natural format)，值的默认格式表示    
		%+v : 类似%v，但输出结构体时会添加字段名     // 类型+值对象
		%#v : 相应值的Go语法表示                  // 输出字段名和字段值形式
	%T : 相应值的类型的Go语法表示，变量的类型
	%% : 百分号，字面上的%，非占位符含义
```

- 字符串
```
	%s : 字符串类型，直接输出字符串或者[]byte        // %ns 打印字符前空n个宽度，默认+，右对齐；若-，左对齐
	%q : 该值对应的双引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
	%x : 每个字节用两字符十六进制数表示(使用a-f)
	%X : 每个字节用两字符十六进制数表示(使用A-F)
```

- 整型
```
	%t : bool类型
	%c : 相应Unicode码点所表示的字符                        // rune: Unico de co de point
	%q : 带单引号 的字符
	%b : 二进制
	%o : 八进制                                           // %#o 带 0 的前缀
	%d : 十进制
		%d : 十进制
			%+d  表示 对正整数 带 符号
			%nd  表示 最小 占位 n 个宽度且右对齐
			%-nd 表示 最小 占位 n 个宽度且左对齐
			%0nd 表示 最小 占位 n 个宽度且右对齐 空字符使用 0 填充
			"%d|%+d|%10d|%-10d|%010d|%+-10d|%+010d"
	%x : 十六进制，小写字母，每字节两个字符
	%X : 十六进制，大写字母，每字节两个字符                    // %#x(%# 带 0x(0X) 的前缀
	%U : Unicode 字符, Unicode格式: 123，等同于 "U+007B"    // %#U 带字符的 Unicode 码点
	%q : 该值对应的双引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
```
```golang
	var a int64 = 50
	// string 2, int32 50, float32 50.000000
	fmt.Printf("%T %[1]v, %T %[2]v, %T %[3]v", string(a), rune(a), float32(a))
```

- 浮点型
```
	%f 、%F : 十进制表示法                                 // %n.mf 表示最小占 n 个宽度并且保留 m 位小数
	%e 、%E : 科学记数法表示
	%g 、%G : 自动选择 最 紧凑的表示 方法 %e(E%) 或 %f(F%)
```

- 指针
```
	%v : 指针变量访问位置中存储的值                           // %#v接口的类型
	%q : 指针变量访问位置中存储的值(unicode 中文)
	%p : 十六进制表示，前缀 0x                               // 默认情况下，指针是已16进制存在的
```

- 特殊字符: 
```
	\ : 反斜线
	' : 单引号                                               // '' 只可以定义单一字符
	" : 双引号                                               // "" 可解析的字符串
	`` : 原始字符串/多行字符串
	\a : 响铃
	\b : 退格
	\f : 换页
	\n : 换行
	\r : 回车
	\t : 制表符
	\v : 垂直制表符
	\ooo : 3 个 8 位数字给定的八进制码点的 Unicode 字符(不能超过\377)
	\uhhhh : 4 个 16 位数字给定的十六进制码点的 Unicode 字符
	\Uhhhhhhhh : 8 个 32 位数字给定的十六进制码点的 Unicode 字符
	\xhh : 2 个 8 位数字给定的十六进制码点的 Unicode 字符
```

- Scan
    - 格式规则类似Printf，有如下区别
```
	%p : 未实现
	%T : 未实现
	%e 、%E 、%f 、%F 、%g 、%G : 效果相同，用于读取浮点数或复数类型
	%s 、%v : 用在字符串时会读取空白分隔的一个片段
	flag '#' 、'+' : 未实现   
```

### 2. Printing
- 函数种类
    - 查询语句 `go doc fmt | grep -Ei "func [FS]*Print`  
    - Print家族: 根据参数内容，串联所有输出生成并写入标准输出`os.Stdout`，包括 `Print()`、`Println()`、`Printf()`
    - Sprint家族: 根据参数内容，串联所有输出生成并返回一个字符串，包括 `Sprint()`、`Sprintln()`、`Sprintf()`
    - Fprint家族: 根据参数内容，串联所有输出生成并写入 `io.Writer`，包括 `Fscan()`、`Fscanf()`、`Fscanln()`
```go
	// State 代表一个传递给自定义Formatter接口的Format方法的打印环境
	type State interface {
		// Write方法用来写入格式化的文本
		Write(b []byte) (ret int, err error)
		// Width返回宽度值，及其是否被设置
		Width() (wid int, ok bool)
		// Precision返回精度值，及其是否被设置
		Precision() (prec int, ok bool)
		// Flag报告是否设置了flag c（一个字符，如+、-、#等）
		Flag(c int) bool
	}

	// 实现了Formatter接口的类型可以定制自己的格式化输出
	type Formatter interface {
		// c为verb，f提供verb的细节信息和Write方法用于写入生成的格式化文本
		Format(f State, c rune)
	}

	// Print 采用默认格式将其参数格式化并写入标准输出，不接受任何格式化操作
	func Print(a ...interface{}) (n int, err error)
	// Println 采用默认格式将其参数格式化并写入标准输出，即 输出到控制台并换行
	func Println(a ...interface{}) (n int, err error)
	// Printf 根据format参数生成格式化的字符串并写入标准输出，且 只可以直接输出字符串类型的变量(不可以输出别的类型)
	func Printf(format string, a ...interface{}) (n int, err error)

	// Sprint 采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串
	func Sprint(a ...interface{}) string
	// Sprintln 采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串
	func Sprintln(a ...interface{}) string
	// Sprintf 根据format参数生成格式化的字符串并返回该字符串
	func Sprintf(format string, a ...interface{}) string

	// Fprint 采用默认格式将其参数格式化并写入w
	func Fprint(w io.Writer, a ...interface{}) (n int, err error)
	// Fprintln 采用默认格式将其参数格式化并写入w
	func Fprintln(w io.Writer, a ...interface{}) (n int, err error)
	// Fprintf 根据format参数生成格式化的字符串并写入w
	func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
``` 

### 3. Scanning
- 函数种类
    - 查询语句 `go doc fmt | grep -Ei "func [FS]*Scan"`  
    - Scan家族: 从标准输入 `os.Stdin`中读取数据，包括 `Scan()`、`Scanln()`、`Scanf()` ，需要使用指针 `Scan(&name)`，scan会直接将输入的值存入指针所指的内存地址的值
    - SScan家族: 从字符串中读取数据，包括 `Sscan()`、`Sscanln()`、`Sscanf()`，即 从字符串扫描到变量
    - Fscan家族: 从 `io.Reader`中读取数据，包括 `Fscan()`、`Fscanln()`、`Fscanf()` ，即 从文件扫描到变量

- 注意事项
    - Scan、Sscan、Fscan将换行符当作空格处理
    - Scanln、Sscanln、Fscanln在遇到换行符的时候停止
    - Scanf、Sscanf、Fscanf根据给定的format格式读取，就像Printf一样
    - Scan家族函数从标准输入读取数据时，以空格做为分隔符分隔标准输入的内容，并将分隔后的各个记录保存到给定的变量中
```go
	// ScanState代表一个将传递给Scanner接口的Scan方法的扫描环境
	type ScanState interface {
		// 从输入读取下一个rune（Unicode码值），在读取超过指定宽度时会返回EOF
		// 如果在Scanln、Fscanln或Sscanln中被调用，本方法会在返回第一个'\n'后再次调用时返回EOF
		ReadRune() (r rune, size int, err error)
		// UnreadRune方法让下一次调用ReadRune时返回上一次返回的rune且不移动读取位置
		UnreadRune() error
		// SkipSpace方法跳过输入中的空白，换行被视为空白
		// 在Scanln、Fscanln或Sscanln中被调用时，换行被视为EOF
		SkipSpace()
		// 方法从输入中依次读取rune并用f测试，直到f返回假；将读取的rune组织为一个[]byte切片返回。
		// 如果skipSpace参数为真，本方法会先跳过输入中的空白。
		// 如果f为nil，会使用!unicode.IsSpace(c)；就是说返回值token将为一串非空字符。
		// 换行被视为空白，在Scanln、Fscanln或Sscanln中被调用时，换行被视为EOF。
		// 返回的切片指向一个共享内存，可能被下一次调用Token方法时重写；
		// 或被使用该Scanstate的另一个Scan函数重写；或者在本次调用的Scan方法返回时重写。
		Token(skipSpace bool, f func(rune) bool) (token []byte, err error)
		// Width返回返回宽度值，及其是否被设置。单位是unicode码值。
		Width() (wid int, ok bool)
		// 因为本接口实现了ReadRune方法，Read方法永远不应被在Scanner接口中使用。
		// 一个合法的ScanStat接口实现可能会选择让本方法总是返回错误。
		Read(buf []byte) (n int, err error)
	}

	// 当Scan、Scanf、Scanln或类似函数接受实现了Scanner接口的类型(其Scan方法的receiver必须是指针，该方法从输入读取该类型值的字符串表示并将结果写入receiver)作为参数时，会调用其Scan方法进行定制的扫描
	type Scanner interface {
		Scan(state ScanState, verb rune) error
	}

	// Scan 从标准输入扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数，换行视为空白
	func Scan(a ...interface{}) (n int, err error)
	// Scanln 类似Scan，但会在换行时才停止扫描
	func Scanln(a ...interface{}) (n int, err error)
	// Scanf 从标准输入扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数
	func Scanf(format string, a ...interface{}) (n int, err error)

	// Sscan 从字符串str扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数，换行视为空白
	func Sscan(str string, a ...interface{}) (n int, err error)
	// Sscanln 类似Sscan，但会在换行时才停止扫描
	func Sscanln(str string, a ...interface{}) (n int, err error)
	// Sscanf 从字符串str扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数
	func Sscanf(str string, format string, a ...interface{}) (n int, err error)

	// Fscan 从r扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数，换行视为空白
	func Fscan(r io.Reader, a ...interface{}) (n int, err error)
	// Fscanln 类似Fscan，但会在换行时才停止扫描
	func Fscanln(r io.Reader, a ...interface{}) (n int, err error)
	// Fscanf 从r扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数
	func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
```

### 4. Errorf
```go
	// Errorf 根据format参数生成格式化字符串并返回一个包含该字符串的错误
	func Errorf(format string, a ...interface{}) error
```

