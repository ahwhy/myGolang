# 基本语法

如下是一个程序的基本组成结构:
```go
// 当前程序的包名, main包表示入口包, 是编译构建的入口
package main

// 导入其他包
import "fmt"

// 常量定义
const PI = 3.1415

// 全局变量声明和赋值
var name = "fly"

// 一般类型声明
type newType int

// 结构体声明
type student struct{}

// 接口声明
type reader interface{}

// 程序入口
func main() {
        fmt.Println("hello world, this is my first golang program!")
}
```

在这个结构里面必须要符合Go程序的语法, 编写出来的程序才是合法的，才能被编译器识别, 编译器识别代码的基础单位是Lexical Token(词法标记)，比如 如下一段代码:
```go
func main() {
        fmt.Println("hello world, this is my first golang program!")
}
```
他包含的Lexical Token(词法标记)，有如下12个:
```
func     // 关键字 func
main     // 标识符 函数名称
(        // LPAREN 左小括号
{        // LBRACE 左花括号
fmt      // 标识符  包名称  
.        // PERIOD 调用符
Println  // 标识符 函数名称
(        // LPAREN 左小括号
"hello world, this is my first golang program!" // 标识符 字符串常量
)        // RPAREN 右小括号
}        // RBRACE 右花括号
)        // LPAREN 右小括号
```

## 非法的词法

ILLEGAL： 标识非法的词法
比如如下一段代码:
```go
func main() {
   中文
}
```
这个`中文`就是一个非法的词法, 编译的时候会直接报错
```sh
src\day1\hello.go:8:2: undefined: 中文
```

## 流结束

EOF： 用于标识 流(io stream)结束
比如parser/parser.go会重复解析声明到文件的最后:
```go
for p.tok != token.EOF {
    decls = append(decls, p.parseDecl(declStart))
}
```

## 注释

COMMENT：注释
Go 支持两种注释方式，行注释和块注释：
+ 行注释：以//开头，例如： //我是行注释
+ 块注释：以/*开头，以*/结尾，例如：/*我是块注释*/

```go
// 这是一个行注释

/*
这是一个块注释
*/
```

## 行分隔符
行分隔符为: ;
但是我们写程序的时候往往都不需要收到写；, 比如Hello world里面的打印语句
```go
fmt.Println("Hello World!")
```
这是因为Go语言的编译器会自动为我们添加

如果你打算将多个语句写在同一行，它们则必须使用 ; 人为区分，但在实际开发中我们并不鼓励这种做法。
而且go fmt 会自动帮你分成2行
```go
fmt.Println("Hello World!");fmt.Println("Hello World!")
```


## 标识符

标识符用来命名变量、类型等程序实体，一个标识符实际上就是一个或是多个字母(A~Z和a~z)数字(0~9)、下划线_组成的序列，但是第一个字符必须是字母或下划线而不能是数字

Go 语言标识符的命名规则如下：
+ 只能由非空字母(Unicode)、数字、下划线(_)组成
+ 只能以字母或下划线开头
+ 不能 Go 语言关键字
+ 避免使用 Go 语言预定义标识符
+ 建议使用驼峰式
+ 标识符区分大小写

下面这些就是一些合法的标识符
```sh
username   xxx   M   user_name   user1
_temp   temp_   heelo1  MMXXX  中文
```

比如这段代码是合法的
```go
package main

import (
	"fmt"
)

func main() {
	中文 := "你好，中文"
	fmt.Println(中文)
}
```

而下面这些就是一个非法的标识符
```
1user  // 数字打头
for    // 关键字不能作为标识符
m*m   // 运算符是不允许的
中 午 // 有空格
```

下面这段代码就会报错
```go
package main

import (
	"fmt"
)

func main() {
	m*2 := "stirng"
	fmt.Println(m*2)
}
```

## 扩展

go 语言支持的所有词法标记如下:
```go
var tokens = [...]string{
        // 特殊词法
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

        // 标识符
	IDENT:  "IDENT",

        // 基本类型
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

        // 操作符
	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

        // 25个关键字
	BREAK:    "break",
	CASE:     "case",
	CHAN:     "chan",
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",

	FUNC:   "func",
	GO:     "go",
	GOTO:   "goto",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	SELECT: "select",
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",
}
```