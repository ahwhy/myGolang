# Golang-Base  Golang的基础语法

## 一、Golang特性简介

- Go语言的编译
	- 决定语法正确的是: 底层编译器
	- 编译的基础单位: Lexical Token (词法标记)，如 `go/token` 包

- Go语言的注释 
	- `// 单行注释`
	- `/* 多行注释 */`
	- 特定场景注释: 构建注释，如: windows，linux  
	- 包注释，包文件开头注释等  

- Go语言中将数据类型分为四类: 基础类型、复合类型、引用类型和接口类型
	- 基础类型，包括: 数字、字符串和布尔型
	- 复合数据类型，包括: 数组和结构体，即通过组合简单类型，来表达更加复杂的数据结构
	- 引用类型，包括指针、切片、字典、函数、通道，虽然数据种类很多，但它们都是对程序中一个变量或状态的间接引用，这意味着对任一引用类型数据的修改都会影响所有该引用的拷贝
	- 接口类型


## 二、程序结构

### 1. 命名&&声明
- 标识符: 标识符是编程时所使用的名字，用于给变量、常量、函数、类型、接口、包名等进行命名，以建立名称和使用之间的关系

- 内置常量: `true、false、nil、iota`

- 内置类型: `bool、byte、rune、int、int8、int16、int32、int64、uint、uint8、unit16、unit32、unit64、uintptr、float32、float64、complex64、complex128、string、error`

- 空白标识符: `_` 使用空白标识符，则无需引用

- 内置函数: `make、len、cap、new、append、copy、close、delete、complex、real、imag、panic、recover`

- 25 关键字，Go语言中类似if和switch的关键字有25个，且关键字不能用于自定义名字，只能在特定语法结构中使用
	- 引用包: `import、package`
	- 实体声明和定义: `const、var、type、func、interface、map、struct、chan`
	- 流程控制: `break、case、continue、default、defer、else、fallthrough、for、go、goto、if、range、return、select、switch`

### 2. 变量
- 在声明变量时，如果初始化表达式被省略，那么将用零值初始化该变量
	- 数值类型变量对应的零值是0
	- 布尔类型变量对应的零值是false
	- 字符串类型对应的零值是空字符串
	- 接口或引用类型(包括slice、map、chan和函数)变量对应的零值是nil
	- 数组或结构体等聚合类型对应的零值是每个元素或字段都是对应该类型的零值

- 初始化表达式可以是字面量或任意的表达式
	- 在包级别声明的变量会在main入口函数执行前完成初始化
	- 局部变量将在声明语句被执行到的时候完成初始化

### 3. 基本语句
```go
	// 当前程序的包名, main包表示入口包, 是编译构建的入口
	package main
	// 导入其他包
	import "fmt"
	// 常量定义
	const PI = 3.1415
	// 全局变量声明和赋值     变量本质: 内存地址; 值: 数据; 变量赋值: 修改值空间里存储的数据; 变量的声明: 强类型 -> 变量指向的值空间，存储的数据，受到类型的限制; 作用: 复用、配置、简洁易读
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


## 三、基础数据类型

- 从底层而言，所有的数据都是由比特组成，但计算机一般操作的是固定大小的数，如整数、浮点数、比特数组、内存地址等

- 进一步将这些数组织在一起，就可表达更多的对象，例如数据包、像素点、诗歌，甚至其他任何对象

- Go语言提供了丰富的数据组织形式，这依赖于Go语言内置的数据类型

- 这些内置的数据类型，兼顾了硬件的特性和表达复杂数据结构的便捷性

### 1. 整型
- Go语言同时提供了有符号和无符号类型的整数运算
	- `int8、int16、int32、int64` 四种截然不同大小的有符号整形数类型，分别对应 `8、16、32、64bit` 大小的有符号整形数，与此对应的是 `uint8、uint16、uint32、uint64` 四种无符号整形数类型
	- `int 、uint` 两种一般对应特定 CPU平台机器字大小的有符号和无符号整数
		- 其中 int是应用最广泛的数值类型
		- 这两种类型都有同样的大小，32或64bit，但不同的编译器即使在相同的硬件平台上可能产生不同的大小
	- `rune` 类型的 Unicode字符是和 int32等价的类型，通常用于表示一个 Unicode码点
	- `byte` 类型的 Unicode字符是和 uint8等价的类型，byte类型一般用于强调数值是一个原始的数据而不是一个小的整数
	- `uintptr` 一种无符号的整数类型，没有指定具体的 bit大小但是足以容纳指针
		- uintptr类型只有在底层编程是才需要，特别是 Go语言和 C语言函数库或操作系统接口相交互的地方

- 不管它们的具体大小，int、uint和uintptr是不同类型的兄弟类型
	- 其中 int和 int32也是不同的类型，即使 int的大小也是 32bit，在需要将 int当作 int32类型的地方需要一个显式的类型转换操作，反之亦然
	- 其中有符号整数采用 2的补码形式表示，也就是最高bit位用作表示符号位，一个 n-bit的有符号数的值域是从 `-2^(n-1)` 到 `2^(n-1) - 1`
	- 无符号整数的所有 bit位都用于表示非负数，值域是 `0` 到 `2^n - 1`
	- 例如，int8类型整数的值域是从 -128到 127，而 uint8类型整数的值域是从 0到 255

### 2. 浮点型
- Go语言提供了两种精度的浮点数，float32和float64

- 它们的算术规范由IEEE754浮点数国际标准定义，该浮点数规范被所有现代的CPU支持

### 3. 复数
- Go语言提供了两种精度的复数类型：complex64和complex128，分别对应float32和float64两种浮点数精度

- 内置的complex函数用于构建复数，内建的real和imag函数分别返回复数的实部和虚部

### 4. 布尔型
- 一个布尔类型的值只有两种：true和false

- if和for语句的条件部分都是布尔类型的值，并且==和<等比较操作也会产生布尔型的值

### 5. 字符串
- 字符串是 Go 语言中的基础数据类型
	- 虽然字符串往往被看做一个整体，但是它实际上是一片连续的内存空间，也可以将它理解成一个由字符组成的数组  
	- 字符串中的每一个元素叫做"字符"

- 字符串的本质
	- 字符串是由字符组成的数组[]byte
	- 数组会占用一片连续的内存空间，而内存空间存储的字节共同组成了字符串
	- Go语言中的字符串只是一个只读的字节数组
```go
	// runtime/string.go
	type stringStruct struct {
		str unsafe.Pointer
		len int
	}
```

- 一个字符串是一个不可改变的字节序列
	- 字符串可以包含任意的数据，包括byte值0，但是通常是用来包含人类可读的文本
	- 文本字符串通常被解释为采用UTF8编码的Unicode码点(rune)序列

- 内置的len函数可以返回一个字符串中的字节数目(不是rune字符数目)
	- 索引操作s[i]返回第i个字节的字节值，i必须满足0 ≤ i< len(s)条件约束
	- 试图访问超出字符串索引范围的字节将会导致panic异常 `panic: index out of range`
	- 子字符串操作s[i:j]基于原始的s字符串的第i个字节开始到第j个字节(并不包含j本身)生成一个新字符串，成的新字符串将包含j-i个字节

- 字符串赋值
```go
	s1 := "My name is 小明😀"                         // 字符串里可以包含任意Unicode字条
	s2 := "He say:\"i'm fine.\" \n\\Thank\tyou.\\"    // 包含转义字符
	s3 := `There is first line.
	
	There is third line.`                             // 反引号里的转义字符无效，反引号里的内容原封不动地进行输出，包括空白符和换行符
```

- 字符串是不可修改的，修改字符串内部数据的操作也是被禁止的
	- 不变性意味如果两个字符串共享相同的底层数据的话也是安全的，这使得复制任何长度的字符串代价是低廉的
	- 同样，一个字符串s和对应的子字符串切片s[7:]的操作也可以安全地共享相同的内存，因此字符串切片操作代价也是低廉的
	- 在这两种情况下都没有必要分配新的内存
```go
	s[0] = 'L' // compile error: cannot assign to s[0]
```

- `+` 操作符将两个字符串链接构造一个新字符串

- byte 和 rune
	- string 中每个元素叫"字符"，字符有两种
		- byte: 1个字节，代表 ASCLL码 的一个字符
		- rune: 4个字节，代表一个 UTF-8字符，一个汉字可用一个 rune 表示
	- string 底层是byte数组，string的长度就是该byte数组的长度，UTF-8 编码下一个汉字占 3个byte，即一个汉字占3个长度
	- string 可以转换为 []byte 或 []rune 类型

- 强制类型转换
	- byte  和 int 可以相互转换
	- float 和 int 可以相互转换，小数位会丢失
	- boot  和 int 不可以相互转换
	- 不同长度的 int 和 float 之间可以相互转换
	- string 可以转换为 []byte 或 []rune 类型，byte 或 rune 类型可以转换为string
	- 低精度向高精度转换没有问题，高精度向低精度转换会丢失位数
	- 无符号向有符号转换，最高位是无符号

### 6. 常量
- 常量表达式的值在编译期计算，而不是在运行期，每种常量的潜在类型都是基础类型: bool、string或数字

- 一个常量的声明语句定义了常量的名字，和变量的声明语法类似，常量的值不可修改，这样可以防止在运行期被意外或恶意的修改

### 7. 其他
- 操作符
	- 算术运算符: `+、-、*、/、%、++、--` 如 i++ 为表达式，非语句，无法赋值给变量
	- 关系运算符: `>、>=、<、<=、==、!=` 判断 A 与 B 的关系，结构: 布尔值，注意函数不可以比较
	- 逻辑运算符: `&&、||、!` &&对应逻辑乘法，||对应逻辑加法，乘法比加法优先级要高
	- 位运算符: `&、|、^、<<、>>、&^` x<<n 左移运算等价于乘以2^n，x>>n 右移运算等价于除以 2^n
	- 赋值运算符: `=、+=、-=、*=、/=、%=、&=、|=、^=、<<=、>>=` 值可能是数据，也可能是地址
	- 其他运算符: `&(单目)、*(单目)、.(点)、-(单目)、…、<-` 单目运算符优先级最高
	- 占位符: `_` /dev/null 1B<>,_ 就是丢弃值  

- 分割符
	- `()`小括号  `[]`中括号  `{}`大括号  `;`分号  `,`逗号  

- 基础数据类型
	|类型|长度(字节byte)|默认值|说明|
	|:------:|:------:|:------:|:------:|
	|bool|1|false|2^8|
	|byte|1|0|uint8, 取值范围[0, 255], 2^8|
	|rune|4|0|Unicode Code Point, int32 , 2^32|
	|int, uint|4或8|0|32 或 64 位, 取决于操作系统, 2^16 或 2^32|
	|int8, uint8|1|0|取值范围[-128, 127], [0, 255]|
	|int16, uint16|2|0|取值范围[-32768, 32767], [0, 65535]|
	|int32, uint32|4|0|取值范围[-21亿, 21亿], [0, 42亿], rune 是 int32 的别名|
	|int64, uint64|8|0|取值范围[-42亿, 42亿], [0, 2^64]|
	|float32|4|0.0|uint8，取值范围[0, 255]，2^8|
	|float64|8|0.0|uint8，取值范围[0, 255]，2^8|
	|complex64|8|||
	|complex128|16|||
	|uintptr|4或8||以存储指针的 uint32 或 uint64 整数|

- 复合数据类型
	|类型|长度(字节byte)|默认值|说明|
	|:------:|:------:|:------:|:------:|
	|array|||值类型|
	|struct|||值类型|
	|string|16bytes(4 * 4)|""|UTF-8 字符串, string = `2*int64` = `2*8bytes`, []byte(string)|
	|slice|24bytes|nil|引用类型, 切片的本质是一个slice结构体指针，指针为一个uint64内存地址，默认值为0，长度为24|
	|map||nil|引用类型|
	|channel||nil|引用类型|
	|interface||nil|接口|
	|function||nil|函数|

- 数据单位
```
	1Word(字) = 2Bytes(字节)
	1Byte = 8bit(位)  // 2^8
	1KB   = 1024B
	1MB   = 1024KB
	1GB   = 1024Mb
	1TB   = 1024GB
```


## 四、流程控制

### 1. 条件语句 - if语句
- 对于条件语句必须有 if 语句，可以有 0 个或多个 else if 语句 ，最多有 1 个 else 语句
 
- if嵌套
```go
	if bool1 {
		/* bool1 = true */
	} else {
		if bool2 {
			/* bool2 = true */
		} else {
			if bool3 {
				/* bool3 = true */
			} else {
				/* bool3 = false */
			}
		}
	}
```

- 多重判断
```go
	if bool1 {
		/* bool1 = true */
	} else if bool表达式2 {
		/* bool2 = true */
	} else if bool表达式3 {
		/* bool3 = true */
	} else {
		/* bool = false */
	}
```

### 2. 选择语句 - switch
- 对于选择语句可以有0个或多个case语句，最多有1个default语句选择条件为true的case语句块开始执行并退出
 
- 若所有条件为false，则执行default语句块并退出

- 可以通过fallthrough修改执行退出行为，继续执行下一条的case或default语句块
```go
	switch var1 {
	case var2 :
		...
		fallthrough  // 只要执行成功，就无视case2的条件，强制执行下一个语句
	case var3 :
		...
	case var4,var5 :
		...
		if(...){
			break
		}
		fallthrough // 此时switch会执行case3和case4，但是如果满足if条件，则只执行case3
	case var6 :
		...
	default:
		...
	}
```

### 3. 循环语句 - for
- 基本用法
	- 基本表达式
	- `break`     用于跳出循环，当条件满足则结束循环
	- `continue`  用于跳过循环，当条件满足这跳过本次循环进行后置或条件子语句执行
```go
	// init:  一般为赋值表达式，给控制变量赋初值
	// condition:  关系表达式或逻辑表达式，循环控制条件
	// post:  一般为赋值表达式，给控制变量增量或减量
	// 执行顺序为
	//     a) 初始化子语句 init
	//     b) 条件子语句   condition
	//     c) 语句块
	//     d) 后置子语句   post
	//     e) b -->c -->d
	//     f) 直到条件子语句为 false 结束循环
	for init; condition; post { 
		...
	}
```

	
- 类while
	- for子语句可以只保留条件子语句，此时类似于其他语言中的 while 循环
```go
	for condition { 
		...
	}
```

- 无限循环
	- for子语句全部省略，则为无限循环(死循环)，常与 break 结合使用
	- `for i := 0;;i++`
	- `for ;; { }`
	- `for true { }`
	- `for { }`

- for-range 
	- 用于遍历 可迭代对象中的每个元素，例如字符串，数组，切片，映射，通道 等
	- 针对包含Unicode 字符的字符串遍历是需要使用 for range
	- range 返回两个元素分别为字节索引index 和 rune 字符，可通过空白标识符_ 忽略需要接收的变量
```go
	for index,value := range iterable {
		...
	}
```

- label 与 goto
	- 通过 goto 语句任意跳转到当前函数指定的 label 位置

- 嵌套循环:
```go
	for [condition | ( init; condition; increment ) | Range] {
		for [condition |  ( init; condition; increment ) | Range]{
			statement(s);
		}
		statement(s);
	}
```

### 4. 示例
- 乘法口诀表
```go
	// 正三角
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%-2d * %-2d = %-2d\t", j, i, i*j)
		}
		fmt.Println()
	}

	// 倒三角
	for i := 1; i < 10; i++ {
		for j := 1; j < i; j++ {
			var n string
			fmt.Printf("%-2s   %-2s   %-2s\t", n, n, n)
		}
		for j := i; j < 10; j++ {
			fmt.Printf("%-2d * %-2d = %-2d\t", j, i, i*j)
		}
		fmt.Println()
	}
```

- 求100以内素数的和
```go
	var sum int
	i := 2
	var isP bool
	for i < 101 {
		isP = true
		j := 2
		for j <= (i / j) {
			if i%j == 0 {
				// fmt.Printf("%d不是素数\n",i)
				isP = false
				break
			}
			j++
		}
		if isP {
			fmt.Printf("%d是素数\n", i)
			sum += i
		}
		i++
	}
	fmt.Println(sum)
	}
```


## 五、复合数据类型

- 基本数据类型，用以构建程序中数据结构，是Go语言的世界的原子

- 复合数据类型，以不同的方式组合基本类型，构造出来的复合数据类型
	- 数组和结构体是聚合类型，它们的值由许多元素或成员字段的值组成
		- 数组是由同构的元素组成——每个数组元素都是完全相同的类型
		- 结构体则是由异构的元素组成的
	- 数组和结构体都是有固定内存大小的数据结构
	- 而 slice 和 map 则是动态的数据结构，它们将根据需要动态增长

### 1. 数组 array
- 数组是具有相同数据类型的数据项组成的一组长度固定的序列，数据项叫做数组的元素，数组的长度必须是非负整数的常量，长度也是类型的一部分
	- 占用内存空间 = length * 数据类型的字节大小
	- 数组的长度是数组类型的一个组成部分，因此[3]int和[4]int是两种不同的数组类型
	- 数组的长度必须是常量表达式，因为数组的长度需要在编译阶段确定

- 当在Go中声明一个数组之后，会在内存中开辟一段固定长度的、连续的空间存放数组中的各个元素，这些元素的数据类型完全相同，可以是内置的简单数据类型(int、string等)，也可以是自定义的struct类型
	- 固定长度: 这意味着数组不可增长、不可缩减；想要扩展数组，只能创建新数组，将原数组的元素复制到新数组
	- 连续空间: 这意味可以在缓存中保留的时间更长，搜索速度更快，是一种非常高效的数据结构，同时还意味着可以通过数值index的方式访问数组中的某个元素
	- 数据类型: 意味着限制了每个block中可以存放什么样的数据，以及每个block可以存放多少字节的数据

- 声明 && 初始化
	- 指定数组长度 `var name [length]type = [length]type{v1, v2, …,vlength}`
	- 使用初始化元素数量推到数组长度 `name := [...]type{v1, v2, …,vlength}`
	- 对指定位置元素进行初始化  `var name [length]type = [length]type{im:vm, …, sin:in}`
```go
	type Currency int
	
	const (
		USD Currency = iota // 美元
		EUR                 // 欧元
		GBP                 // 英镑
		RMB                 // 人民币
	)
	
	symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}
	
	fmt.Println(RMB, symbol[RMB]) // "3 ￥"
```

- 数组的比较
	- 如果一个数组的元素类型是可以相互比较的，那么数组类型也是可以相互比较的
	- 可以直接通过 == 比较运算符 来比较两个数组，只有当两个数组的所有元素都是相等的时候，数组才是相等的
	- 不相等比较运算符 != 遵循同样的规则
	- 切片为引用类型，不能比较
```go
	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}
	fmt.Println(a == b, a == c, b == c) // "true false false"
	d := [3]int{1, 2}
	fmt.Println(a == d) // compile error: cannot compare [2]int == [3]int
```

- 指针数组
	- 声明一个指针类型的数组，这样数组中就可以存放指针
	- 指针的默认初始化值为nil

```go
	a := [4]*int{0: new(int), 3: new(int)}   // [0xc00011a300 <nil> <nil> 0xc00011a308]
	a[1] = new(int)                          // 空指针直接赋值会报错
	*a[1] = 10                               // [0xc00011a300 0xc00011a310 <nil> 0xc00011a308]
	b := a                                   // [0xc00011a300 0xc00011a310 <nil> 0xc00011a308]

	func zero(ptr *[32]byte) {
		*ptr = [32]byte{}
	}
```

- 遍历数组
```go
	for i := 0; i < len(name); i++ {
		fmt.Println(i, name[i])
	}
	
	for i, j := range name {
		fmt.Printf("%d %q\n", i, j)
	}
```

- 多维数组
```go
	// 声明 && 初始化
	var name [vlength][vvlength]type = [vlength][vvlength]type{{v1,v2, …,vvlength}, {v1,v2, …,vvlength}, …,{vlength,vvlength}}
	name := [...][vvlength]type{{v1,v2, …,vvlength}, {v1,v2, …,vvlength}, …,{vlength,vvlength}}  // 多维数组只有第一维长度可使用变量数量推测
	name := [vlength][vvlength]type{0:{0:v1,3:v2},5:{2:v1,5:v2, …,m:v3}, …,n:{6:v1,m:vvlength}}

	// 遍历
	for i := 0; i < len(name); i++ {
		for j := 0; j < len(name[i]); j++ {
			fmt.Printf("[%d ,%d]: %q\n", i, j, name[i][j]) 
		}
	}
	
	for i, line := range name {
		for n, m := range line {
			fmt.Printf("[%d ,%d]: %q\n", i, n, m)
		}
	}
```

### 2. 切片 slice
- 切片是长度可变的'数组'，即 具有相同数据类型的数据项组成的一组长度可变的序列 ，切片由三部分组成
	- 指针(array): 指向 第一个slice元素对应的底层数组元素的地址，要注意的是slice的第一个元素并不一定就是数组的第一个元素
	- 长度(length): 切片元素的数量
	- 容量(capacity): 切片开始到结束位置(可容纳)元素的数量

- Go语言中的slice依赖于数组，它的底层就是数组，所以数组具有的优点，slice都有
```go
	// runtime/slice.go
	type slice struct {
		array unsafe.Pointer // 数组指针
		len   int // 长度 
		cap   int // 容量
	}
```

- Go语言中的slice支持通过append向slice中追加元素，长度不够时会动态扩展，通过再次slice切片，可以得到得到更小的slice结构，可以迭代、遍历等
	- 切片共享底层数组，若某个切片元素发生变化，则数组和其他有共享元素的切片也会发生变化
	- 切片底层是一个长度和数据类型固定的数组，只有在切片的长度大于底层数组的长度后，该切片的底层才会在内存中更换新的数组

- 声明 && 初始化
	- 使用字面量初始化空切片 `[]type{}`，初始化为零值 nil
	- 使用字面量初始化 `var name []type = []type{v1, v2, …,vn}`
	- 指定长度和容量字面量初始化 `[]type{im:vm, in:vn, ilength:vlength]`
	- 使用make函数初始化 `make([]type, len)/make([]type, len, cap)`，即 通过 make 函数创建长度为 len，容量为 cap 的切片，且 len 必须小于等于 cap
	- 使用数组切片操作初始化 
```go
	array[start:end]            // end <= src_cap ; 新创建切片长度和容量计算: new_len: end-start, new_cap: src_cap-start ;
	array[start:end:end_cap]    // 用于限制新切片的容量值(end<=cap<=src_cap)；新创建切片长度和容量计算 new_len: end-start, new_cap: end_cap-start
```

- 遍历切片 
	- 同数组
	- 关于Golang遍历中存在的问题
		- 因为for range在遍历值类型时，其中的v变量是一个值的拷贝
		- 当使用&获取指针时，实际上是获取到v这个临时变量的指针，而v变量在for range中只会创建一次，之后循环中会被一直重复使用
		- 所以在arr2赋值的时候其实都是v变量的指针，而&v最终会指向arr1最后一个元素的值拷贝
		- [Go语言中for range的"坑"](https://www.jianshu.com/p/3bef2c245102)
		- [Go官方 CommonMistakes](https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable)

```go
	// 结果 3, 3, 3
	arr1 := []int{1, 2, 3}
	arr2 := make([]*int, len(arr1))
	for i, v := range arr1 {
		arr2[i] = &v
	}
	for _, v := range arr2 {
		fmt.Println(*v)
	}
	
	// 改进 - 传递原始指针
	for i := range arr1 {
		arr2[i] = &arr1[i]
	}
	
	// 改进 - 使用临时变量
	for i, v := range arr1 {
		t := v
		arr2[i] = &t
	}
	
	// 改进 - 使用闭包
	for i, v := range arr1 {
		func(v int){
			arr2[i] = &v
		}(v)
	}
```

- 增加元素 
	- 使用append函数对切片增加一个或多个元素并返回修改后切片，当长度在容量范围内时只增加长度，容量和底层数组不变。
	- 当长度超过容量范围则会创建一个新的底层数组并对容量进行智能运算(元素数量<1024时，约按原容量1倍增加，>1024时约按原容量0.25倍增加)
```go
	append(slice, 1, 2, ...,n)  
	// 移除元素 
	append(slince[:n-1], slince[n+1]...)
```

- 复制切片
```go
	copy(drc_slice, src_slice)  
	// 移除元素 
	copy(slice[3:], slice[4:])
```

- 切片的比较
	- 数组可以比较，但切片不能进行比较
	- 标准库提供了高度优化的bytes.Equal函数来判断两个字节型slice是否相等`[]byte`，
	- 对于其他类型的slice，必须展开每个元素进行比较
	- slice唯一合法的比较操作是和nil比较
```go
	// 展开每个元素进行比较
	func equal(x, y []string) bool {
		if len(x) != len(y) {
			return false
		}
		for i := range x {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	}

	// 和nil比较
	if summer == nil { /* ... */ }
```

- nil silce
	- 一个零值的slice等于nil
	- 一个nil值的slice并没有底层数组
	- 一个nil值的slice的长度和容量都是0，但是也有非nil值的slice的长度和容量也是0的，例如`[]int{}`或`make([]int, 3)[3:]`
	- 与任意类型的nil值一样，可以用`[]int(nil)`类型转换表达式来生成一个对应类型slice的nil值
```go
	var s []int    // len(s) == 0, s == nil
	s = nil        // len(s) == 0, s == nil
	s = []int(nil) // len(s) == 0, s == nil
	s = []int{}    // len(s) == 0, s != nil
```

- 应用
```go
	// 用切片实现队列
	queue := []int{}
	queue = append(queue, 1)
	queue = append(queue, 2)
	queue = queue[1:]

	// 用切片实现堆栈
	stack := []int{}
	stack = append(stack, 1)	
	stack = append(stack, 2)
	stack = stack[:len(stack)-1]
```

- 多维切片
```go
	// 声明&&初始化
	var name [][]type = [][]type{{v1,v2, …,vvlength}, {v1,v2, …,vvlength}, …,{vlength,vvlength}}
	name := [][]type{0:{0:v1,3:v2},5:{2:v1,5:v2, …,m:v3}, …,n:{6:v1,m:vvlength}}

	// append
	slice = append(slice, []int{1, 2, 3})
	slice[0] = append(point[0], 1)

	// copy
	slice2 := [][]int{{}, {}}
	copy(slice2, slice)
```

- 其他
```go
	// 字节切片 bytes包
	[]byte{string}
	string([]byte{})

	// rune切片
	[]rune{string}
	string([]rune{})
```

### 3. 映射 map
- 哈希表是一种巧妙并且实用的数据结构，它是一个无序的key/value对的集合，其中所有的key都是不同的，然后通过给定的key可以在常数时间复杂度O(1)内检索、更新或删除对应的value

- 在Go语言中，一个 map 就是一个哈希表的引用，map类型可以写为 `map[K]V` ，其中 K 和 V 分别对应key和value
	- key 经过哈希后，对哈希结果使用map槽位总数取模，得到的结果，就插入对应槽位
		- 比如 key1 哈希后得到 5437，使用槽位总数取模，得到 5 ，即插入第 5 个槽位
		- 若插入时冲突太多，则需要增加槽位数，给每个 key 重新分配槽位
	- map中所有的key都有相同的类型，所有的value也有着相同的类型，但是key和value之间可以是不同的数据类型
	- 开销
		- 稀疏型数据结构，牺牲空间换取时间
		- 相对而言的数组，完整型数据结构
	- map删除数据时存在延迟，所以最好不作为内存存储

- Go语言中只要是可比较的类型都可以作为 key，除开 slice，map，functions 这几种类型，其他类型都是 OK 的 
	- 具体包括: 布尔值、数字、字符串、指针、通道、接口类型、结构体、只包含上述类型的数组
	- 虽然 浮点数类型也是支持相等运算符比较的，但是将浮点数用做key类型则是一个坏的想法，最坏的情况是可能出现的NaN和任何浮点数都不相等
	- 对于V对应的value数据类型则没有任何的限制。
	- 这些类型的共同特征是支持 == 和 != 操作符，k1 == k2 时，可认为 k1 和 k2 是同一个 key 
	- 如果是结构体，则需要它们的字段值都相等，才被认为是相同的 key
	- key和map本身都可以被哈希(类似漏斗，可以把很大一部分数据，筛选成一小部分特征，又叫散列)

- 声明&&初始化
	- map声明需要指定组成 元素key 和 value 的类型
		- map类型的零值是nil，也就是没有引用任何哈希表
		- 在声明后，若不注意会被初始化为 nil，表示暂不存在的映射 
		- nil map，它将不会做任何初始化，不会指向任何数据结构，添加元素会报错，而直接赋值会报空指针，map类型实际上就是一个指针，具体为 `*hmap`
		- 如在结构体中新增map字段，在后面的引用代码中需要添加make初始化map
	- 初始化
		- 使用字面量初始化 `map[ktype]vtype{k1:v1, k2:v2, …, kn:vn}`，key 包括 string、int、bool、array
		- 使用字面量初始化空映射 `map[ktype]vtype{}`，若不加{}，则初始化为nil，即无法添加key
		- 使用make函数初始化 `make(map[ktype]vtype)`，通过make函数创建映射，它会先创建好底层数据结构，然后再创建map，并让map指向底层数据结构
	- 判断是否存在
		- 通过key访问元素时可接收两个值，第一个值为value，第二个值为bool类型表示元素是否存在，若存在为true，否则为false
	- 修改&增加
		- 使用key对映射赋值时当key存在则修改key对应的value，若key不存在则增加key和value
	- 删除
		- 使用delete函数删除映射中已经存在的key
	- 比较
	- map中的元素并不是一个变量，禁止对map元素取址，原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而可能导致之前的地址无效
```go
	// 判断元素是否存在
	map_01, ok := map[1]
	fmt.Printf("%t, %v\n", ok, map_01)

	// 删除元素
	delete(map, 3)
	delete(map[2], "XX")
	
	// 比较元素
	func equal(x, y map[string]int) bool {
		if len(x) != len(y) {
			return false
		}
		for k, xv := range x {
			if yv, ok := y[k]; !ok || yv != xv {
				return false
			}
		}
		return true
	}	
	
	// 禁止对map元素取址
	_ = &ages["bob"] // compile error: cannot take address of map element
```
	
- 多维映射
```go
	map := map[int]map[string]string{1: map[string]string{"name": "aa", "tel": "123"}, 2: map[string]string{"name": "bb", "tel": "456"}}
```

- 映射的遍历
	- 基本语句
	- 顺序遍历
		- Map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序
		- 每次都使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现
		- 如果要按顺序遍历key/value对，必须显式地对key进行排序，可以使用sort包的Strings函数对字符串slice进行排序

```go
	// 映射的遍历
	for k, v := range map {
		fmt.Printf("%v:%v\n", k, v)
	}

	// 顺序遍历
	import "sort"
	
	names := make([]string, 0, len(ages))
	for name := range ages {
		names = append(names, name)
	}
	
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
```


- 函数中的应用
```go
	// map作为函数参数
	// map是一种指针，所以将map传递给函数，仅仅只是复制这个指针，所以函数内部对map的操作会直接修改外部的map
	a := map[int]string{1: "a", 2: "b", 3: "c"}
	func(map[int]string) {
		delete(a, 1)
	}(a)

	// map值为函数
	op := map[string]func(x, y int) int{
		"+": func(x, y int) int {
			return x + y
		},
		"-": func(x, y int) int {
			return x - y
		},
		"*": func(x, y int) int {
			return x * y
		},
		"/": func(x, y int) int {
			return x / y
		},
	}
	fmt.Println(op["+"](1, 2))
	fmt.Println(op["-"](1, 2))

	// slice 作为 map 的key
	// 有时候需要一个 map 或 set 的 key 是 slice类型，但是 map 的 key 必须是可比较的类型，slice并不满足这个条件
	// 可以通过两个步骤绕过这个限制
	//   1、定义一个辅助函数k，将slice转为map对应的string类型的key，确保只有x和y相等时k(x) == k(y)才成立
	//   2、创建一个key为string类型的map，在每次对map操作时先用k辅助函数将slice转化为string类型
	var m = make(map[string]int)
	func k(list []string) string { return fmt.Sprintf("%q", list) }
	func Add(list []string)       { m[k(list)]++ }
	func Count(list []string) int { return m[k(list)] }
```

- 存在问题 Go语言中原生的map线程不安全
	- 解决方案一: 加锁
	- 解决方案二: 使用 `sync.map`
		- go 1.9引入的内置方法，并发线程安全的map
		- sync.Map 将 key和 value 按照 interface{}存储
		- 查询出来后要类型断言 x.(int) x.(string)
		- 遍历使用Range() 方法，需要传入一个匿名函数作为参数，匿名函数的参数为k,v interface{}，每次调用匿名函数将结果返回
		- sync.map 性能对比
			- https://studygolang.com/articles/27515
			- 性能对比结论
				- 只读场景: sync.map > rwmutex >> mutex
				- 读写场景(边读边写): rwmutex > mutex >> sync.map
				- 读写场景(读80% 写20%): sync.map > rwmutex > mutex
				- 读写场景(读98% 写2%): sync.map > rwmutex >> mutex
				- 只写场景: sync.map >> mutex > rwmutex
			- sync.Map使用场景的建议
				- 读多: 给定的key-v只写一次，但是读了很多次，只增长的缓存场景
				- key不相交: 覆盖更新的场景比少
				- 结构体复杂的case多不用sync.Map
	- 解决方案三: 分片锁 并发map github.com/orcaman/concurrent-map
```go
	// 方案一
	type concurrentMap struct {
		sync.RWMutex
		mp map[int]int
	}
	func (c *concurrentMap) Set(key int, value int) {
		// 获取写锁
		c.Lock()
		c.mp[key] = value
		c.Unlock()
	}
	func (c *concurrentMap) Get(key int) int {
		// 获取读锁
		c.RLock()
		res := c.mp[key]
		c.RUnlock()
		return res
	}
	c := concurrentMap{
		mp: (map[int]int{}),
	}
	// 写map的goroutine
	go func() {
		for i := 0; i < 10000; i++ {
			c.Set(i, i)
		}
	}()
	// 读map的goroutine
	go func() {
		for i := 0; i < 10000; i++ {
			res := c.Get(i)
			fmt.Printf("[cmap.get][%d=%d]\n", i, res)
		}
	}()
	time.Sleep(20 * time.Second)
```

- 使用map统计 "I have a dream" 中，各字母出现次数
```go
	article := `
		I have a dream
		...
		`
	stats := map[rune]int{}
	
	for _, ch := range article {
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
			stats[ch]++
		}
	}
	for ch, cnt := range stats {
		fmt.Printf("%c, %v\n", ch, cnt)
	}
```

- Set的实现
```go
	type Set struct {
		sync.RWMutex
		Mset map[interface{}]interface{}
	}

	func NewSet(item ...interface{}) *Set {
		s := &Set{}
		s.Mset = make(map[interface{}]interface{})
		s.Add(item...)

		return s
	}

	// Add 添加元素
	func (s *Set) Add(items ...interface{}) error {
		s.Lock()
		defer s.Unlock()

		for _, item := range items {
			s.Mset[item] = item
		}

		return nil
	}

	// Delete 删除元素
	func (s *Set) Delete(items ...interface{}) error {
		s.Lock()
		defer s.Unlock()

		for _, item := range items {
			delete(s.Mset, item)
		}

		return nil
	}

	// Modify 修改元素
	func (s *Set) Modify(items ...interface{}) error {
		s.Lock()
		defer s.Unlock()

		for _, item := range items {
			_, ok := s.Mset[item]
			if !ok {
				continue
			}

			switch value := item.(type) {
			case string:
				s.Mset[item] = value + "modify"
			case int:
				s.Mset[item] = value + 10
			case bool:
				s.Mset[item] = true
			case float64:
				s.Mset[item] = value + 3.1415926535
			}
		}

		return nil
	}

	// Contains 包含&&查询
	func (s *Set) Contains(item interface{}) bool {
		s.RLock()
		defer s.RUnlock()

		_, ok := s.Mset[item]

		return ok
	}

	// Size 容量
	func (s *Set) Size() int {
		s.RLock()
		defer s.RUnlock()

		return len(s.Mset)
	}

	// Clear 清除
	func (s *Set) Clear() {
		s.Lock()
		defer s.Unlock()

		s.Mset = make(map[interface{}]interface{})
	}

	// Equel 比较
	func (s *Set) Equel(other *Set) bool {
		// 如果两者Size不相等，返回false
		if s.Size() != other.Size() {
			return false
		}

		// 迭代查询遍历
		for key := range s.Mset {
			// 若不存在返回false
			if !other.Contains(key) {
				return false
			}
		}

		return true
	}

	// IsSubset 子集
	func (s *Set) IsSubset(other *Set) bool {
		// s的size长于other，返回false
		if s.Size() > other.Size() {
			return false
		}

		// 迭代遍历
		for key := range s.Mset {
			if !other.Contains(key) {
				return false
			}
		}

		return true
	}
```

- 类k8s-deployment的实现
```go
	/*
	- interface
	- map
	- 发布系统，能发布k8s node(jtype字段) 不同的类型的任务
	- 增量更新：开启新的，停掉旧的
	- 原有的是a,b,c ，b,c,d --> 开启d, 停掉a
	*/

	type jobManager struct {
		jobMutex sync.RWMutex
		// 增量更新
		activeJobs map[string]deployJob
	}

	func NewJobManager() *jobManager {
		return &jobManager{
			jobMutex:   sync.RWMutex{},
			activeJobs: map[string]deployJob{},
		}
	}

	func (jm *jobManager) sync(jobs []deployJob) {
		// ***增量更新
		jb := make(map[string]deployJob)

		for _, v := range jobs {
			jb[v.GetName()] = v
		}

		jm.jobMutex.Lock()
		defer jm.jobMutex.Unlock()
		for k, jmv := range jm.activeJobs {
			if jbv, ok := jb[k]; ok && jmv.GetKind() == jbv.GetKind() {
				delete(jb, k)
			} else {
				go jm.activeJobs[k].stop()
			}
		}
		for k, v := range jb {
			jm.activeJobs[k] = v
			go jm.activeJobs[k].start()
		}
	}

	type deployJob interface {
		start()
		stop()
		SetFlag(b bool)
		GetName() string
		GetFlag() bool
		GetKind() string
	}

	type k8sD struct {
		Name string
		Flag bool
		Kind string
	}

	func NewK8sD(name string) *k8sD {
		return &k8sD{
			Name: name,
			Kind: "K8sD",
		}
	}

	func (k *k8sD) start() {
		k.SetFlag(true)
		log.Printf("Pod %s 开始运行 \n", k.Name)
		count := 0
		for {
			if k.Flag {
				count++
				log.Printf("Pod %s 正在运行 %d 秒\n", k.Name, count)
				time.Sleep(1 * time.Second)
			} else {
				log.Printf("Pod %s 退出运行 \n", k.Name)
				break
			}
		}
	}

	func (k *k8sD) stop() {
		k.SetFlag(false)
		log.Printf("Pod %s 退出运行 \n", k.Name)
	}

	func (k *k8sD) SetFlag(b bool) {
		k.Flag = b
	}

	func (k *k8sD) GetName() string {
		return k.Name
	}

	func (k *k8sD) GetFlag() bool {
		return k.Flag
	}

	func (k *k8sD) GetKind() string {
		return k.Kind
	}

	type host struct {
		Name string
		Flag bool
		Kind string
		// Pod  []k8sD
	}

	func NewHost(name string) *host {
		return &host{
			Name: name,
			Kind: "Host",
		}
	}

	func (h *host) start() {
		h.SetFlag(true)
		log.Printf("Host %s 开始运行 \n", h.Name)
		count := 0
		for {
			count++
			if h.Flag {
				log.Printf("Host %s 正在运行 %d 秒\n", h.Name, count)
				time.Sleep(1 * time.Second)
			} else {
				log.Printf("Host %s 退出运行 \n", h.Name)
				break
			}
		}
	}

	func (h *host) stop() {
		h.SetFlag(false)
		log.Printf("Host %s 退出运行 \n", h.Name)
	}

	func (h *host) SetFlag(b bool) {
		h.Flag = b
	}

	func (h *host) GetName() string {
		return h.Name
	}

	func (h *host) GetFlag() bool {
		return h.Flag
	}

	func (h *host) GetKind() string {
		return h.Kind
	}
```

- 带过期时间的map
	- map做缓存用的 垃圾堆积k1、k2 
	- 希望缓存存活时间 5分钟
	- 将加锁的时间控制在最低
	- 耗时的操作在加锁外侧做
```go
	type item struct {
		value interface{} // 值
		ts    int64       // 时间戳，item被创建出来的时间，或者被更新的时间
	}

	func NewItem(i interface{}) *item {
		return &item{
			value: i,
			ts:    time.Now().Unix(),
		}
	}

	// 带过期缓存的map
	type Cache struct {
		sync.RWMutex
		mp map[string]*item
	}

	func NewCache() *Cache {
		return &Cache{
			mp: make(map[string]*item),
		}
	}

	func (c *Cache) Get(key string) *item {
		c.RLock()
		defer c.RUnlock()

		return c.mp[key]
	}

	func (c *Cache) Set(key string, value *item) {
		c.Lock()
		defer c.Unlock()

		c.mp[key] = value
	}

	// GC 先加读锁 -> 检查确实有需要回收的数据 -> 合并写锁回收
	func (c *Cache) Gc(timeDelta int64) {
		for {
			toDelKeys := make([]string, 0)
			now := time.Now().Unix()
			c.RLock()

			// 遍历缓存中的项目，对比时间戳，超过 timeDelta 就删除该项目
			for k, v := range c.mp {
				if now-v.ts > timeDelta {
					log.Printf("[项目已过期][key %s]", k)
					toDelKeys = append(toDelKeys, k)
				}
			}
			c.RUnlock()

			c.Lock()
			for _, k := range toDelKeys {
				delete(c.mp, k)
			}
			c.Unlock()

			time.Sleep(5 * time.Second)
		}
	}

	c := NewCache()
	// 异步执行，删除过期项目的任务
	go c.Gc(10)

	// 写入数据
	for i := 0; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		item := NewItem(i)

		// 设置缓存
		log.Printf("[设置项目缓存[key:%s][v:%v]", key, item)
		c.Set(key, item)
	}

	time.Sleep(15 * time.Second)

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		item := NewItem(i + 2)

		// 更新缓存
		log.Printf("[更新项目缓存][key:%s][v:%v]", key, item)
		c.Set(key, item)
	}

	time.Sleep(10 * time.Second)
```

- 带过期时间的缓存 github.com/patrickmn/go-cache
```go
	var (
		defaultInterval = time.Second * 30
		UserCache       = cache.New(defaultInterval, defaultInterval)
	)

	type user struct {
		Name  string
		Email string		
		Phone int64
	}

	// HttpGetUser 初始化user对象
	func HttpGetUser(name string) user {
		u := user{
			Name:  name,
			Email: "qq.com",
			Phone: time.Now().Unix(),
		}

		return u
	}

	func GetUser(name string) user {
		res, found := UserCache.Get(name)
		if found {
			u := res.(user)
			log.Printf("[found_user_in_cache][name:%s][value:%v]", name, u)

			return u
		} else {
			res := HttpGetUser(name)
			// 给每个key 单独设置超时时间
			UserCache.Set(name, res, defaultInterval)
			log.Printf("[not_found_user_in_cache][query_by_http][name:%s][value:%v]", name, res)

			return res
		}
	}
```


### 4. 结构体 struct
- 定义
	- 结构体是由一些列属性组成的复合数据类型，每个属性都具有名称、类型和值，结构体将属性组合在一起由程序进行处理
	- 结构体定义使用 struct 标识，需要指定其包含的属性(名和类型)，在定义结构体时可以为 结构体指定结构体名(命名结构体)，用于后续声明结构体变量使用
```go
	type struct_variable_type struct {
		member definition
		...
	}
```

- 声明
	- 声明结构体变量只需要定义变量类型为结构体名，变量中的每个属性被初始化为对应类型的零值 `var struct_name struct_variable_type`
	- 可以声明结构体指针变量，此时变量被初始化为 nil `var struct_name *struct_variable_type`
	
- 初始化
```go
	// 使用结构体零值初始化结构体值对象
	var struct_name struct_variable_type = struct_variable_type{}

	// 使用结构体字面量初始化结构体值对象
	var struct_name struct_variable_type = struct_variable_type{
		...
	}
	struct_name := struct_variable_type{a: 123, c: 789...}    //初始化部分

	// 使用 new 函数进行初始化结构体指针对象
	var struct_name *struct_variable_type = new(struct_variable_type)

	// 使用结构体字面量初始化结构体指针对象
	var struct_name *struct_variable_type = &struct_variable_type{
		...
	}
	struct_name := &struct_variable_type{a: 123, c: 789...}
```
	
- 属性的访问和修改 
	- 通过结构体对象名/结构体指针对象.属性名的方式来访问和修改对象的属性值 `struct_name.a`
	
- 匿名结构体
	- 在定义变量时将类型指定为结构体的结构，此时叫匿名结构体
	- 匿名结构体常用于初始化一次结构体变量的场景，例如项目配置
```go
	a := struct{
		a1 int
		a2 string
		...
	}{1, "aa"}
```

- 命名嵌入
```go
	// 结构体命名嵌入是指结构体中的属性对应的类型也是结构体
	type Address struct{
		region string
		street string
		num    string
	}
	type User struct{
		name string
		tel  string
		addr Address
	}
	
	// 声明&&初始化 
	var u1 User
	var a1 Address = Address{
		"浙江省杭州市",
		"转塘街道",
		"001",
	}
	var u2 User = User{"atlantis", "18100000000", a1}
	u3 := User{
		name: "lywq",
		tel: "18600000000",
		addr: Address{
			"安徽省合肥市",
			"望江西路",
			"001",
		},
	}
	
	// 属性的访问和修改
	u3.addr.num = "002"
```

- 匿名嵌入
	- 结构体匿名嵌入是指将已定义的结构体名直接声明在新的结构体中，从而实现对以后已有类型的扩展和修改
	- 声明和初始化
		- 在初始化匿名嵌入的结构体对象时需要遵循树状声明的结构，对于匿名嵌入的结构体可以使用结构体名来指定初始化参数
	- 属性的访问和修改	
		- 在访问和修改嵌入结构体的属性值时，可以通过对象名.结构体名称.属性名的方式进行访问和修改，结构体名称可以省略（匿名成员有一个隐式的名称），因此不能嵌套两个相同名称的结构体
		- 当被嵌入结构体和嵌入结构体有相同的属性名时，在访问和修改嵌入结构体成员的属性值时不能省略结构体名称
```go
	type Employee struct {
		User
		salary float64
		title  string
	}	
```

- 指针类型嵌入
```go
	// 结构体嵌入(命名&匿名)类型也可以为结构体指针
	type PUser *struct{              // 命名嵌入结构体指针
		name string
		tel  string
		addr *Address
	}
	type PEmployee struct {          // 匿名嵌入结构体指针
		*PUser
		salary float64
		title  string
	}
	// 声明&&初始化
	puser := PUser{"atlantis", "18100000000", &a1}
	pemployee := PEmployee{
		PUser:  &puser,
		salary: 25000,
		title:  "aa",
	}
```

- 属性的访问和修改	
	- 使用属性为指针类型底层共享数据结构，当底层数据发生变化，所有引用都会发生影响
 
- 空结构体
	- `struct{}`
		- struct{} 是一个无元素的结构体类型，通常在没有信息存储时使用
		- 优点是大小为0，不需要内存来存储 struct{}类型的值
	- `struct{}{}`
		- struct{}{} 是一个复合字面量，它构造了一个 struct{}类型的值，该值也是空 `var empty = struct{}{}`

- 可见性
	- 结构体名首字母大写则包外可见(公开的)，小写则仅包内可访问(内部的)
	- 结构体属性名首字母大写包外可见(公开的)，小写则仅包内可访问(内部的)
		- 结构体名首字母大写，属性名大写: 结构体可在包外使用，且访问其大写的属性名
		- 结构体名首字母大写，属性名小写: 结构体可在包外使用，且不能访问其小写的属性名
		- 结构体名首字母小写，属性名大写: 结构体只能在包内使用，属性访问在结构体嵌入时由被嵌入结构体(外层)决定，被嵌入结构体名首字母大写时属性名包外可见，否者只能在包内使用
		- 结构体名首字母小写，属性名小写: 结构体只能在包内使用(可以使用函数返回首字母小写的结构体，以达到控制访问的作用)


- 结构体方法
	- 结构体的行为，需要用方法才能完成
	- Go语言中的方法是作用在指定的数据类型上的(和指定的数据类型绑定)
	- 自定义类型都可以有方法，不仅仅是结构体 
	- 方法的声明与调用
```go
	type A struct {
		Name string  
	}  
	func (a A) say(){           // 表示结构体A有一个方法，名称为say
		fmt.Println(a.Name)     // (a A) 体现say 方法与A 类型绑定  
	}
```	


## 六、引用类型

### 1. 指针
- 指针与变量
	- 变量是一种使用方便的占位符，用于引用计算机内存地址
	- 可以通过`&`运算符获取变量在内存中的存储位置，指针是用来存储变量地址的变量
	- 本质 -> 内存地址，没有数据类型，占 8bytes

- 声明初始化与赋值
	- 指针声明需要指定存储地址中对应数据的类型，并使用`*`作为类型前缀
	- 指针变量声明后会被初始化为 nil，表示空指针
	- 使用 new 函数初始化后，new 函数根据数据类型申请内存空间并使用零值填充，并返回申请空间地址
	- 指针赋值
	- 指针运算
		- `&` 获取变量的地址，赋值成为指针 `b := &a`，获取变量a在内存中的存储位置，即 取引用 - 取地址
		- `*` 获取指针指向的值 `*b`，指针变量存储地址中的值，即 解引用 - 返回内存地址中对应的对象 - 取值
```go
	// 指定存储地址中对应数据的类型
	// 指向整型
	var ip *int
	// 指向浮点型	
	var fp *float32
	
	// 使用 new 函数
	var a *int = new(int)
	fmt.Println(a)       // 0xc000014330
	fmt.Println(*a)      // 0

	// 指针赋值
	var a *int = new(int)
	*a = 10 
	fmt.Println(a, *a)   // 0xc000014330, 10
```

- 多维指针
	- 用来存储指针变量地址的变量叫做指针的指针
```go
	var a ****int
	v := 10
	p1 := &v                        // *int
	p2 := &p1                       // **int    &p1，即指针变量本身的内存地址
	p3 := &p2                       // ***int
	a = &p3                         // ****int
	fmt.Println(v, p1, p2, p3, a)
```

### 2. 其他
- slice、map、channel 是Go语言里的3种引用类型，都可以通过make函数进行初始化(申请分配内存)，因为其包含一个指向底层数据结构的指针，所以称为"引用"类型

- 引用类型未初始化时都是nil，可以对它们指向len()函数，返回0

- 深浅拷贝
	- 定义&&区别
		- 深拷贝 copy的是数据本身
			- 复制的时候会创建一个新的对象
			- 指向完全不同的内存地址
			- 修改互不影响
		- 浅拷贝 复制的是对象的指针，
			- 新老对象指向同一块内存区域
			- a修改字段，b中的字段也一同被修改
			- 内存消耗是一致的
	- 类型
		- 值类型的都是深拷贝 `int、float、bool、array、struct`
		- 引用类型都是浅拷贝 `slice、map、function`

- new和make对比
	- new 开辟一个类型对应的内存空间，返回一个内存空间的地址；且只能分配地址，一般用于基础类型的初始化；
	- `make{makeslice, makemap, makechannel}` make返回创建对象的内存地址  
	以slice为例，unsafe.Pointer --> slince struct --- {member: pointer ---> array}  表现为: []int  
	
	|函数类型|适用范围|返回值|填充类型|
	|:------:|:------:|:------:|:------:|
	|new|new可以对所有类型进行分配|new返回指针|new填充零值|
	|make|make只能创建类型(slice、map、channel)|make返回引用|make填充非零值|