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
- Print: 输出到控制台，不接受任何格式化操作
- Println: 输出到控制台并换行
- Printf: 只可以打印出格式化的字符串；只可以直接输出字符串类型的变量(不可以输出别的类型)
- Sprintf: 格式化并返回一个字符串而不带任何输出
- Fprintf: 来格式化并输出到 io.Writers 而不是 os.Stdout `func fmt.Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)`

### 3. Scanning
- 函数种类
    - 查询语句 `go doc fmt | grep -Ei "func [FS]*Scan"`  
    - Scan家族: 从标准输入os.Stdin中读取数据，包括Scan()、Scanf()、Scanln() `需要使用指针 Scan(&name) scan会直接将输入的值存入指针所指的内存地址的值`
    - SScan家族: 从字符串中读取数据，包括Sscan()、Sscanf()、Sscanln() `即从字符串扫描到变量 func fmt.Sscan(str string, a ...interface{}) (n int, err error)`
    - Fscan家族: 从io.Reader中读取数据，包括Fscan()、Fscanf()、Fscanln() `即从文件扫描到变量 func fmt.Fscan(r io.Reader, a ...interface{}) (n int, err error)`  

- 注意事项
    - Scanln、Sscanln、Fscanln在遇到换行符的时候停止
    - Scan、Sscan、Fscan将换行符当作空格处理
    - Scanf、Sscanf、Fscanf根据给定的format格式读取，就像Printf一样
    - Scan家族函数从标准输入读取数据时，以空格做为分隔符分隔标准输入的内容，并将分隔后的各个记录保存到给定的变量中。