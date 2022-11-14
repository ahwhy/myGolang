# Golang-Strconv  Golang的类型转换

## 一、Golang的标准库 strconv包

### 1. strconv
- strconv
	- strconv包实现了基本数据类型和其字符串表示的相互转换
	- 提供了字符串与简单数据类型之间的类型转换功能，可以将简单类型转换为字符串，也可以将字符串转换为其它简单类型

### 2. string 与 bool
- string 转 bool
```go
	// ParseBool 返回字符串表示的bool值
	// 它接受1、0、t、f、T、F、true、false、True、False、TRUE、FALSE；否则返回错误
	// strconv.ParseBool("true")
	func ParseBool(str string) (value bool, err error)
```

- bool 转 string
```go
	// FormatBool 根据b的值返回"true"或"false"
	// strconv.FormatBool(true)
	func FormatBool(b bool) string
```

### 3. string 与 int
- string 转 int
```go
	// ParseInt 返回字符串表示的整数值，接受正负号
	// base指定进制(2到36)，如果base为0，则会从字符串前置判断，"0x"是16进制，"0"是8进制，否则是10进制
	// bitSize指定结果必须能无溢出赋值的整数类型，0、8、16、32、64 分别代表 int、int8、int16、int32、int64
	// 返回的err是*NumErr类型的，如果语法有误，err.Error = ErrSyntax；如果结果超出类型范围err.Error = ErrRang
	// strconv.ParseInt("11111111", 2, 16)
	func ParseInt(s string, base int, bitSize int) (i int64, err error)

	// Atoi 是 ParseInt(s, 10, 0)的简写
	func Atoi(s string) (i int, err error)
	// For Example
	i, err := strconv.Atoi("100x")
	fmt.Printf("type %v, value: %d, err: %v\n", reflect.TypeOf(i), i, err)
```

- int 转 string
```go
	// FormatInt 返回i的base进制的字符串表示。base 必须在2到36之间，结果中会使用小写字母'a'到'z'表示大于10的数字
	// strconv.FormatInt(255, 10)
	func FormatInt(i int64, base int) string

	// Itoa 是 FormatInt(i, 10) 的简写
	func Itoa(i int) string
	// For Example
	str := strconv.Itoa(100)
	fmt.Printf("type %v, value: %s\n", reflect.TypeOf(str), str)
```

### 4. string 与 uint
- string 转 uint
```go
	// ParseUint类似ParseInt但不接受正负号，用于无符号整型
	// strconv.ParseUint("4E2D", 16, 16)
	func ParseUint(s string, base int, bitSize int) (n uint64, err error)
```

- uint 转 string
```go
	// FormatUint 是FormatInt的无符号整数版本
	// strconv.FormatUint(255, 16)
	func FormatUint(i uint64, base int) string
```

### 5. string 与 int
- string 转 float
```go
	// ParseFloat 解析一个表示浮点数的字符串并返回其值，如果s合乎语法规则，函数会返回最为接近s表示值的一个浮点数(使用IEEE754规范舍入)
	// bitSize指定了期望的接收类型，32是float32(返回值可以不改变精确值的赋值给float32)，64是float64
	// 返回值err是*NumErr类型的，语法有误的，err.Error=ErrSyntax；结果超出表示范围的，返回值f为±Inf，err.Error= ErrRange
	// strconv.ParseFloat("3.1", 64)
	func ParseFloat(s string, bitSize int) (f float64, err error)
```

- float 转 string
```go
	// FormatFloat 将浮点数表示为字符串并返回
	// bitSize表示f的来源类型(32：float32、64：float64)，会据此进行舍入
	// fmt表示格式：'f'(-ddd.dddd)、'b'(-ddddp±ddd，指数为二进制)、'e'(-d.dddde±dd，十进制指数)、'E'(-d.ddddE±dd，十进制指数)、'g'(指数很大时用'e'格式，否则'f'格式)、'G'(指数很大时用'E'格式，否则'f'格式)
	// prec控制精度(排除指数部分)：对'f'、'e'、'E'，它表示小数点后的数字个数；对'g'、'G'，它控制总的数字个数。如果prec 为-1，则代表使用最少数量的、但又必需的数字来表示f
	// strconv.FormatFloat(3.1415, 'E', -1, 64)
	func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```