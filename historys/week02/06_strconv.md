# 数据类型转换

Go语言是强类型, 如果需要对数据进行类型转换，需要手动进行, 转换数据类型的简单方法是 直接通过`类型()`方式:
```go
valueOfTypeB = typeB(valueOfTypeA)
```

例如:
```go
a := 3.14
b := int(a)
```

我们可以通过： 标准库 reflect包中的TypeOf方法查看一个变量的类型, 比如:
```go
a := 10
b := 0.1314
c := "hello"

fmt.Printf("a type: %v\n", reflect.TypeOf(a))
fmt.Printf("b type: %v\n", reflect.TypeOf(b))
fmt.Printf("c type: %v\n", reflect.TypeOf(c))
```

## 自定义类型转换

Go允许我们通过type定义自己的类型,自己定义的类型和该类型不是同一类型了，比如:
```go
type Age int
var a Age = 10
var b int = 20
fmt.Println(reflect.TypeOf(a))  // day2.Age
fmt.Println(reflect.TypeOf(b))  // int
```

此时我们定义的Age类型已经不再是int类型了, 只是该类型底层的值为int

```go
// Age 底层数据结构为 int
type Age int
// a 类型是Age 底层为 int 10
var a Age = 10

// 将a转化成int类型,
// 由于a是Age, 转化成int后, 他们不是同一种类型，不能再次赋值回去: a = int(a) 是不行的
b := int(a)
// 现在b是int类型
fmt.Println(reflect.TypeOf(b))

// 反过来我们也可以将int类型转换为Age类型
c := Age(10)
// 现在c就是Age类型，而不是int类型了
fmt.Println(reflect.TypeOf(c))
```

## 直接转换

对于数值类 的类型 比如 int类 uint类 float类 他们直接是可以相互转换的

```go
var a float64 = 5.3232223232323
fmt.Println(float32(a))
fmt.Println(int(a))
fmt.Println(int8(a))
fmt.Println(int16(a))
fmt.Println(int32(a))
fmt.Println(int64(a))
fmt.Println(uint(a))
```

1.精度损失问题
低精度转换为高精度时是安全的，高精度的值转换为低精度时会丢失精度。例如int32转换为int16，float32转换为int

2.跨大类型转换无法转换
不是所有数据类型都能转换的，例如字母格式的string类型"abcd"转换为int肯定会失败

```go
var a float64 = 5.3232223232323
fmt.Println(string(a)) // 这种会报错
```

如果我们想要跨大类型转换可以使用strconv包提供的函数

## strconv

strconv包提供了字符串与简单数据类型之间的类型转换功能，可以将简单类型转换为字符串，也可以将字符串转换为其它简单类型

## string和int转换

1. int转string的方法是: Itoa()

```go
str := strconv.Itoa(100)
fmt.Printf("type %v, value: %s\n", reflect.TypeOf(str), str)
```

2.string转int的方法是:

```go
i, err := strconv.Atoi("100")
fmt.Printf("type %v, value: %d, err: %v\n", reflect.TypeOf(i), i, err)
```

并不是所有string都能转化为int, 所以可能会报错:

```go
i, err := strconv.Atoi("100x")
fmt.Printf("type %v, value: %d, err: %v\n", reflect.TypeOf(i), i, err)
```

## string转其他类型

strconv包提供的Parse类函数用于将字符串转化为给定类型的值：ParseBool()、ParseFloat()、ParseInt()、ParseUint()
由于字符串转换为其它类型可能会失败，所以这些函数都有两个返回值，第一个返回值保存转换后的值，第二个返回值判断是否转换成功。

1.转bool

```go
b, err := strconv.ParseBool("true")
fmt.Println(b, err)
```

2.转float

```go
f1, err := strconv.ParseFloat("3.1", 32)
fmt.Println(f1, err)
f2, err := strconv.ParseFloat("3.1", 64)
fmt.Println(f2, err)
```

由于浮点数的小数部分 并不是所有小数都能在计算机中精确的表示, 这就造成了浮点数精度问题, 比如下面

```go
var n float64 = 0
for i := 0; i < 1000; i++ {
    n += .01
}
fmt.Println(n)
```

关于浮点数精度问题: c[计算机不都是0101吗，你有想过计算机是怎么表示的小数吗](https://www.21ic.com/article/881429.html), 简单理解就是:

将其整数部分与小树部分分开, 比如5.25

+ 对于整数部分 5 ，我们使用"不断除以2取余数"的方法，得到 101
+ 对于小数部分 .25 ，我们使用"不断乘以2取整数"的方法，得到 .01

听说有一个包可以解决这个问题: github.com/shopspring/decimal

3.转int

```go
func ParseInt(s string, base int, bitSize int) (i int64, err error)
```

+ base: 进制，有效值为0、2-36。当base=0的时候，表示根据string的前缀来判断以什么进制去解析：0x开头的以16进制的方式去解析，0开头的以8进制方式去解析，其它的以10进制方式解析
+ bitSize: 多少位，有效值为0、8、16、32、64。当bitSize=0的时候，表示转换为int或uint类型。例如bitSize=8表示转换后的值的类型为int8或uint8

```go
fmt.Println(bInt8(-1)) // 0000 0001(原码) -> 1111 1110(反码) -> 1111 1111

// Parse 二进制字符串
i, err := strconv.ParseInt("11111111", 2, 16)
fmt.Println(i, err)
// Parse 十进制字符串
i, err = strconv.ParseInt("255", 10, 16)
fmt.Println(i, err)
// Parse 十六进制字符串
i, err = strconv.ParseInt("4E2D", 16, 16)
fmt.Println(i, err)
```

4.转uint

```go
func ParseUint(s string, base int, bitSize int) (uint64, error)
```

用法和转int一样, 只是转换后的数据类型是uint64

```go
u, err := strconv.ParseUint("11111111", 2, 16)
fmt.Println(u, err)
u, err = strconv.ParseUint("255", 10, 16)
fmt.Println(u, err)
u, err = strconv.ParseUint("4E2D", 16, 16)
fmt.Println(u, err)
```

## 其他类型转string

将给定类型格式化为string类型：FormatBool()、FormatFloat()、FormatInt()、FormatUint()。

```go
fmt.Println(strconv.FormatBool(true))

// 问题又来了
fmt.Println(strconv.FormatInt(255, 2))
fmt.Println(strconv.FormatInt(255, 10))
fmt.Println(strconv.FormatInt(255, 16))

fmt.Println(strconv.FormatUint(255, 2))
fmt.Println(strconv.FormatUint(255, 10))
fmt.Println(strconv.FormatUint(255, 16))

fmt.Println(strconv.FormatFloat(3.1415, 'E', -1, 64))
```

```go
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```

+ bitSize表示f的来源类型（32：float32、64：float64），会据此进行舍入。
+ fmt表示格式：'f'（-ddd.dddd）、'b'（-ddddp±ddd，指数为二进制）、'e'（-d.dddde±dd，十进制指数）、'E'（-d.ddddE±dd，十进制指数）、'g'（指数很大时用'e'格式，否则'f'格式）、'G'（指数很大时用'E'格式，否则'f'格式）。
+ prec控制精度（排除指数部分）：对'f'、'e'、'E'，它表示小数点后的数字个数；对'g'、'G'，它控制总的数字个数。如果prec 为-1，则代表使用最少数量的、但又必需的数字来表示f。