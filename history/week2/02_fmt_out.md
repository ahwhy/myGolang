# 字符串格式化

Go语言用于控制文本输出常用的标准库是fmt

fmt中主要用于输出的函数有:

+ Print:   输出到控制台,不接受任何格式化操作
+ Println: 输出到控制台并换行
+ Printf : 只可以打印出格式化的字符串。只可以直接输出字符串类型的变量（不可以输出别的类型）
+ Sprintf：格式化并返回一个字符串而不带任何输出
+ Fprintf：来格式化并输出到 io.Writers 而不是 os.Stdout

我们通过Printf函数来了解下，Go语言里面的字符串格式化:
```
fmt.Sprintf(格式化样式, 参数列表…)
```

+ 格式样式: 字符串形式，格式化符号以 % 开头， %s 字符串格式，%d 十进制的整数格式
+ 参数列表: 多个参数以逗号分隔，个数必须与格式化样式中的个数一一对应，否则运行时会报错

比如:
```go
username := "boy"
fmt.Printf("welcome, %s", username)
```

## 整数类型

|格  式 | 描  述 |
|  ----  | --- |
|%b	 |整型以二进制方式显示|
|%o	 |整型以八进制方式显示|
|%d	 |整型以十进制方式显示|
|%x	 |整型以十六进制方式显示|
|%X	 |整型以十六进制、字母大写方式显示|
|%c	 |相应Unicode码点所表示的字符|
|%U	 |Unicode 字符, Unicode格式：123，等同于 "U+007B"|

```go
a := 255
fmt.Printf("二进制: %b\n", a)
fmt.Printf("八进制: %o\n", a)
fmt.Printf("十进制: %d\n", a)
fmt.Printf("十六进制: %x\n", a)
fmt.Printf("大写十六进制: %X\n", a)

fmt.Printf("十六进制: %d\n", Hex2Dec("4E2D"))
fmt.Printf("字符: %c\n", 20013)
fmt.Printf("Unicode格式: %U\n", '中') // U+4E2D
```

## 浮点数

|格  式 | 描  述 |
|  ----  | --- |
|%e |     科学计数法，例如 -1234.456e+78 |
|%E |     科学计数法，例如 -1234.456E+78 |
|%f |     有小数点而无指数，例如 123.456 |
|%g |     根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出 |
|%G |     根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出|

```go
fmt.Printf("%e", 12675757563.5345432567) //1.267576e+10
fmt.Printf("%E", 12675757563.5345432567) // 1.267576E+10
fmt.Printf("%f", 12675757563.5345432567) // 12675757563.534544
fmt.Printf("%g", 12675757563.5345432567) // 1.2675757563534544e+10
fmt.Printf("%G", 12675757563.5345432567) // 1.2675757563534544E+1
```

## 布尔

|格  式 | 描  述 |
|  ----  | --- |
|%t |     true 或 false |

```go
fmt.Printf("%t", true)
```

## 字符串

|格  式 | 描  述 |
|  ----  | --- |
|%s  |   字符串或切片的无解译字节 |
|%q  |   双引号围绕的字符串，由Go语法安全地转义 |
|%x  |   十六进制，小写字母，每字节两个字符 |
|%X  |   十六进制，大写字母，每字节两个字符|

```go
str := "I'm a boy"
fmt.Printf("%s", str)
fmt.Println()
fmt.Printf("%q", str)
fmt.Println()
fmt.Printf("%x", str)
fmt.Println()
fmt.Printf("%X", str)
fmt.Println()
```

## 指针

|格  式 | 描  述 |
|  ----  | --- |
|%p     |十六进制表示，前缀 0x |

```go
a := "I'm a boy"
b := &a
fmt.Printf("%p", b)
```

## 通用的占位符

|格  式 | 描  述 |
|  ----  | --- |
|%v  |    值的默认格式。|
|%+v |   类似%v，但输出结构体时会添加字段名|
|%#v |　 相应值的Go语法表示 |
|%T  |   相应值的类型的Go语法表示 |
|%%  |   百分号,字面上的%,非占位符含义|

默认格式%v下，对于不同的数据类型，底层会去调用默认的格式化方式：

``` go
bool:                    %t 
int, int8 etc.:          %d 
uint, uint8 etc.:        %d, %x if printed with %#v
float32, complex64, etc: %g
string:                  %s
chan:                    %p 
pointer:                 %p
```

```go
f := 1010.0101
s := "hey boy!"
fmt.Printf("%v", f)
fmt.Println()

fmt.Printf("%v", s)
fmt.Println()
```

如果是复杂对象的话,按照如下规则进行打印：

```go
struct:            {field0 field1 ...} 
array, slice:      [elem0 elem1 ...] 
maps:              map[key1:value1 key2:value2] 
pointer to above:  &{}, &[], &map[]
```

```go
user := User{"laowang", 33}
fmt.Printf("%v", user) // Go默认形式 {laowang 33}
fmt.Println()
fmt.Printf("%#v", user) //类型+值对象 day2.User{Name:"laowang", Age:33}
fmt.Println()
fmt.Printf("%+v", user) //输出字段名和字段值形式 {Name:laowang Age:33}
fmt.Println()
fmt.Printf("%T", user) //值类型的Go语法表示形式, day2.User
fmt.Println()
fmt.Printf("%%")
```

## 宽度

我们输出时 可能需要控制字符串的宽度和小数的精度

### 字符串宽度控制

宽度设置格式: 占位符中间加一个数字, 数字分正负, +: 右对齐, -: 左对齐

1.最小宽度, 不够部分可以选择补0

```go
fmt.Printf("|%s|", "aa") // 不设置宽度
fmt.Printf("|%5s|", "aa") // 5个宽度,  默认+， 右对齐
fmt.Printf("|%-5s|", "aa") // 5个宽度, 左对齐

fmt.Printf("|%05s|", "aa") // |000aa|
```

2.最大宽度, 超出的部分会被截断

```go
fmt.Printf("|%.2s|", "xxxx") // 最大宽度为5
```

### 浮点数精度控制

你也可以指定浮点型的输出宽度，同时也可以通过 宽度.精度 的语法来指定输出的精度

```go
a := 54.123456
fmt.Printf("|%f|", a)
fmt.Printf("|%5.1f|", a)
fmt.Printf("|%-5.1f|", a)
fmt.Printf("|%05.1f|", a)
```

思考: 不同语言的文字宽度并不一定相同, 比如

```go
fmt.Printf("|%2s|", "中国")
fmt.Printf("|%2s|", "ab")
```

可以参考: [获取字符的宽度](http://www.nbtuan.vip/2017/05/10/golang-char-width/)

## 格式化错误

1.类型错误或未知: %!verb(type=value)

```go
Printf("%d", "hi") // %!d(string=hi)
```

2.太多参数: %!(EXTRA type=value)

```go
Printf("hi", "guys") // hi%!(EXTRA string=guys)
```

3.太少参数: %!verb(MISSING)

```go
Printf("hi%d") // hi %!d(MISSING)
```

4.宽度/精度不是整数值: %!(BADWIDTH) or %!(BADPREC)

```go
Printf("%d", hi) // %!d(string=hi)
```

## 作业

请将这段二进制翻译成中文(unicode编码)

```go
// 1000 1000 0110 0011
// 0110 0111 0000 1101
// 0101 0101 1001 1100
// 0110 1011 0010 0010
// 0111 1010 0111 1111
// 0100 1110 0010 1101
// 0101 0110 1111 1101
// 0111 1110 1010 0010
```