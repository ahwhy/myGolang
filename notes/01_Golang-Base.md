# Golang-Base

## 一、基础语法
### 1、底层语法
- 决定语法正确的是: 底层编译器
- 编译的基础单位: Lexical Token (词法标记)，如 `go/token` 包
- 标识符: 标识符是编程时所使用的名字，用于给变量、常量、函数、类型、接口、包名等进行命名，以建立名称和使用之间的关系
- 注释: 
	- `// 单行注释`
	- `/* 多行注释 */`
	- 特定场景注释: 构建注释，如: windows，linux  
	- 包注释，包文件开头注释等  

### 2、基础数据类型
- 内置常量: `true、false、nil、iota`
- 内置类型: `bool、byte、rune、int、int8、int16、int32、int64、uint、uint8、unit16、unit32、unit64、uintptr、float32、float64、complex64、complex128、string、error`
- 空白标识符: `_` 使用空白标识符，则无需引用
- 内置函数: `make、len、cap、new、append、copy、close、delete、complex、real、imag、panic、recover`
- 其他
	- 各数据类型所占字节大小
```
		in8            1byte
		int64/uint64   8bytes 2^64
		rune           4bytes
		string         16bytes(4 * 4)   // string = 2*int64 = 2*8bytes  []byte(string)
		slice          24bytes          // 切片的本质是一个slice结构体指针，指针为一个uint64内存地址，默认值为0，长度为24
```
	- 数据单位
```
		1Word(字) = 2Bytes(字节)
		1Byte = 8bit(位)  // 2^8
		1KB   = 1024B
		1MB   = 1024KB
		1GB   = 1024Mb
		1TB   = 1024GB
```
	- new和make对比
		- new 开辟一个类型对应的内存空间，返回一个内存空间的地址；且只能分配地址，一般用于基础类型的初始化；
		- `make{makeslice, makemap, makechannel}` make返回创建对象的内存地址  
		以slice为例，unsafe.Pointer --> slince struct --- {member: pointer ---> array}  表现为: []int  
		
		|函数类型 |适用范围|返回值|填充类型|
		|:------:|:------:|:------:|:------:|
		|new|new可以对所有类型进行分配|new返回指针|new填充零值|
		|make|make只能创建类型(slice、map、channel)|make返回引用|make填充非零值|  

### 3、25 关键字
- 声明: `import、package`
- 实体声明和定义: `chan、const、func、interface、map、struct、type、var`
- 流程控制: `break、case、continue、default、defer、else、fallthrough、for、go、goto、if、range、return、select、switch`  

### 4、操作符
- 算术运算符: `+、-、*、/、%、++、--` i++ 为表达式，非语句，无法赋值给变量
- 关系运算符: `>、>=、<、<=、==、!=` 判断 A 与 B 的关系，结构: 布尔值，注意函数不可以比较
- 逻辑运算符: `&&、||、!`
- 位运算符: `&、|、^、<<、>>、&^`
- 赋值运算符: `=、+=、-=、*=、/=、%=、&=、|=、^=、<<=、>>=` 值可能是数据，也可能是地址
- 其他运算符: `&(单目)、*(单目)、.(点)、-(单目)、…、<-` 单目运算符优先级最高
- 占位符: `_` /dev/null 1B<>,_ 就是丢弃值  

### 5、分割符
- `()`小括号  `[]`中括号  `{}`大括号  `;`分号  `,`逗号  

### 6、指针
- 指针与变量
	- 变量是一种使用方便的占位符，用于引用计算机内存地址
	- 可以通过`&`运算符获取变量在内存中的存储位置，指针是用来存储变量地址的变量
	- 本质 -> 内存地址，没有数据类型，占 8bytes
- 声明初始化与赋值
	- 指针声明需要指定存储地址中对应数据的类型，并使用`*`作为类型前缀
```
		var ip *int        // 指向整型
		var fp *float32    // 指向浮点型
```
	- 指针变量声明后会被初始化为 nil，表示空指针
	- 使用 new 函数初始化后，new 函数根据数据类型申请内存空间并使用零值填充，并返回申请空间地址
```	
		var a *int = new(int)
		fmt.Println(a)       // 0xc000014330
		fmt.Println(*a)      // 0
```
	- 指针赋值
```
		var a *int = new(int)
		*a = 10 
		fmt.Println(a, *a)   // 0xc000014330, 10
```
	- 指针运算
		- `&` 获取变量的指针 `b := &a 获取变量a在内存中的存储位置，即 取引用 - 取地址`
		- `*` 获取指针指向的值 `*b 指针变量存储地址中的值，即 解引用 - 返回内存地址中对应的对象 - 取值`
- 多维指针
	- 用来存储指针变量地址的变量叫做指针的指针
```
		var a ****int
		v := 10
		p1 := &v                        // *int
		p2 := &p1                       // **int    &p1，即指针变量本身的内存地址
		p3 := &p2                       // ***int
		a = &p3                         // ****int
		fmt.Println(v, p1, p2, p3, a)
```
  
### 7、基本注解
```
	// 当前程序的包名, main包表示入口包, 是编译构建的入口
	package main
	// 导入其他包
	import "fmt"
	// 常量定义
	const PI = 3.1415
	// 全局变量声明和赋值       变量本质: 内存地址; 值: 数据; 变量赋值: 修改值空间里存储的数据; 变量的声明: 强类型 -> 变量指向的值空间，存储的数据，受到类型的限制; 作用: 复用、配置、简洁易读
	var name = "fly"
	// 定义"别名"，counter类型实际还是int，比如 rune 为 int32； byte 为 uint8
	type counter = int
	// 一般类型声明
	type newType int
	// 函数声明
	type myFun func(x, y int) int
	// 结构体声明
	type student struct{}
	// 接口声明
	type reader interface{}
	// 程序入口
	func main() {
			fmt.Println("hello world, this is my first golang program!")
	}
```
  
  
## 二、输入输出
### 1、标准输入  
#### 1).fmt包
- Print: 输出到控制台，不接受任何格式化操作
- Println: 输出到控制台并换行
- Printf: 只可以打印出格式化的字符串；只可以直接输出字符串类型的变量(不可以输出别的类型)
- Sprintf: 格式化并返回一个字符串而不带任何输出
- Fprintf: 来格式化并输出到 io.Writers 而不是 os.Stdout `func fmt.Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)`

#### 2).占位符
- 类型
```
	%v : 值的默认格式。
		%+v : 类似%v，但输出结构体时会添加字段名  // 类型+值对象
		%#v : 相应值的Go语法表示                  // 输出字段名和字段值形式
	%T : 相应值的类型的Go语法表示
	%% : 百分号,字面上的%,非占位符含义
```
- 字符串
```
	%s : 字符串类型                               // %ns 打印字符前空n个宽度，默认+，右对齐；若-，左对齐
	%q : 双引号围绕的字符串
```
- 整型
```
	%t : bool类型
	%c : 相应Unicode码点所表示的字符                        // rune: Unico de co de point
	%q : 带单引号 的字符                                    
	%b : 二进制                                             
	%o : 八进制                                             // %#o 带 0 的前缀
	%d : 十进制 
		%d : 十进制 
			%+d  表示 对正整数 带 符号
			%nd  表示 最小 占位 n 个宽度且右对齐
			%-nd 表示 最小 占位 n 个宽度且左对齐
			%0nd 表示 最小 占位 n 个宽度且右对齐 空字符使用 0 填充
			"%d|%+d|%10d|%-10d|%010d|%+-10d|%+010d"	
	%x : 十六进制，小写字母，每字节两个字符
	%X : 十六进制，大写字母，每字节两个字符                                  // %#x(%# 带 0x(0X) 的前缀
	%U : Unicode 字符, Unicode格式: 123，等同于 "U+007B"    // %#U 带字符的 Unicode 码点
```
- 浮点型
```
	%f 、 %F : 十进制表示法                                 // %n.mf 表示最小占 n 个宽度并且保留 m 位小数
	%e 、 %E : 科学记数法表示
	%g 、 %G : 自动选择 最 紧凑的表示 方法 %e(E%) 或 %f(F%)
```
- 指针
```
	%t : 指针变量的具体内容
	%v : 指针变量访问位置中存储的值                         // %#v接口的类型
	%q : 指针变量访问位置中存储的值(unicode 中文)
	%p : 十六进制表示，前缀 0x                              // 默认情况下，指针是已16进制存在的
```
- 特殊字符: 
```
	\: 反斜线
	': 单引号                                               // '' 只可以定义单一字符
	": 双引号                                               // "" 可解析的字符串
	``: 原始字符串/多行字符串
	\a: 响铃
	\b: 退格
	\f: 换页
	\n: 换行
	\r: 回车
	\t: 制表符
	\v: 垂直制表符
	\ooo: 3 个 8 位数字给定的八进制码点的 Unicode 字符(不能超过\377)
	\uhhhh: 4 个 16 位数字给定的十六进制码点的 Unicode 字符
	\Uhhhhhhhh: 8 个 32 位数字给定的十六进制码点的 Unicode 字符
	\xhh: 2 个 8 位数字给定的十六进制码点的 Unicode 字符
```

### 2、标准输出  
#### fmt包 	`go doc fmt | grep -Ei "func [FS]*Scan"`  
- Scan家族: 从标准输入os.Stdin中读取数据，包括Scan()、Scanf()、Scanln() `需要使用指针 Scan(&name) scan会直接将输入的值存入指针所指的内存地址的值`
- SScan家族: 从字符串中读取数据，包括Sscan()、Sscanf()、Sscanln() `即从字符串扫描到变量 func fmt.Sscan(str string, a ...interface{}) (n int, err error)`
- Fscan家族: 从io.Reader中读取数据，包括Fscan()、Fscanf()、Fscanln() `即从文件扫描到变量 func fmt.Fscan(r io.Reader, a ...interface{}) (n int, err error)`  

#### 注意事项
- Scanln、Sscanln、Fscanln在遇到换行符的时候停止
- Scan、Sscan、Fscan将换行符当作空格处理
- Scanf、Sscanf、Fscanf根据给定的format格式读取，就像Printf一样
- Scan家族函数从标准输入读取数据时，以空格做为分隔符分隔标准输入的内容，并将分隔后的各个记录保存到给定的变量中。

### 3、流结束
- EOF: 用于标识 流(io stream)结束 
	- 如 parser/parser.go 会重复解析声明到文件的最后
```
		for p.tok != token.EOF {
			decls = append(decls, p.parseDecl(declStart))
		}
```  

### 4、格式转换
- strconv包
	- 提供了字符串与简单数据类型之间的类型转换功能，可以将简单类型转换为字符串，也可以将字符串转换为其它简单类型  `import "strconv"`
	- int转string   ->   strconv.Itoa()
```
		str := strconv.Itoa(100)
		fmt.Printf("type %v, value: %s\n", reflect.TypeOf(str), str)
```
	- string转int   ->   strconv.Atoi()
```
		i, err := strconv.Atoi("100x")
		fmt.Printf("type %v, value: %d, err: %v\n", reflect.TypeOf(i), i, err)
```
	- string转bool  ->  strconv.ParseBool()
```
		b, err := strconv.ParseBool("true")
```
	- string转float ->  strconv.ParseFloat()
```
		f, err := strconv.ParseFloat("3.1", 64)
```
	- string转int   ->  func ParseInt(s string, base int, bitSize int) (i int64, err error)
```
		i, err := strconv.ParseInt("11111111", 2, 16)
```
	- string转uint  ->  func ParseUint(s string, base int, bitSize int) (uint64, error)
```
		u, err = strconv.ParseUint("4E2D", 16, 16)
```
	- bool转string  ->  strconv.FormatBool(true)
	- float转string ->  strconv.FormatFloat(3.1415, 'E', -1, 64)  `func FormatFloat(f float64, fmt byte, prec, bitSize int) string`
	- int转string   ->  strconv.FormatInt(255, 10)
	- uint转string  ->  strconv.FormatUint(255, 16) 