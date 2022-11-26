# Golang-Text  Golang的文本处理

## 一、Golang的标准库 scanner包

### 1. text/scanner
- text/scanner
	- scanner包提供对utf-8文本的token扫描服务
		- 它会从一个io.Reader获取utf-8文本，通过对Scan方法的重复调用获取一个个token
		- 为了兼容已有的工具，NUL字符不被接受
		- 如果第一个字符是表示utf-8编码格式的BOM标记，会自动忽略该标记
	- 一般Scanner会跳过空白和Go注释，并会识别所有go语言规格的字面量
		- 它可以定制为只识别这些字面量的一个子集，也可以识别不同的空白字符
```go
	// TokenString 返回一个token或unicode码值的可打印的字符串表示
	func TokenString(tok rune) string

	// 代表资源里的一个位置
	type Position struct {
		Filename string // 文件名（如果存在）
		Offset   int    // 偏移量，从0开始
		Line     int    // 行号，从1开始
		Column   int    // 列号，从1开始（每行第几个字符）
	}
	// IsValid 返回所处的位置是否合法
	func (pos *Position) IsValid() bool
	func (pos Position) String() string

	// Scanner类型实现了token和unicode字符(从io.Reader中)的读取
	type Scanner struct {
		// 每一次出现错误时都会调用该函数；如果Error为nil，则会将错误报告到os.Stderr。
		Error func(s *Scanner, msg string)
		// 每一次出现错误时，ErrorCount++
		ErrorCount int
		// 控制那些token被识别。如要识别整数，就将Mode的ScanInts位设为1。随时都可以修改Mode。
		Mode uint
		// 控制那些字符识别为空白。如果要将一个码值小于32的字符视为空白，只需将码值对应的位设为1；
		// 空格码值是32，大于32的位设为1的行为未定义。随时都可以修改Whitespace。
		Whitespace uint64
		// 最近一次扫描到的token的开始位置，由Scan方法设定
		// 调用Init或Next方法会使位置无效（Line==0），Scanner不会操作Position.Filename字段
		// 如果发生错误且Position不合法，此时扫描位置不在token内，应调用Pos获取错误发生的位置
		Position
		...
	}
	// Init 使用src创建一个Scanner，并将Error设为nil，ErrorCount设为0，Mode设为GoTokens，Whitespace 设为GoWhitespace
	func (s *Scanner) Init(src io.Reader) *Scanner
	// Pos 返回上一次调用Next或Scan方法后读取结束时的位置
	func (s *Scanner) Pos() (pos Position)
	// Peek 返回资源的下一个unicode字符而不移动扫描位置；如果扫描位置在资源的结尾会返回EOF
	func (s *Scanner) Peek() rune
	// Next 读取并返回下一个unicode字符
	func (s *Scanner) Next() rune
	// Scan 从资源读取下一个token或者unicode字符并返回它
	func (s *Scanner) Scan() rune
	// TokenText 返回最近一次扫描的token对应的字符串；应该在Scan方法后调用
	func (s *Scanner) TokenText() string

	// For Example
	var s scanner.Scanner
	s.Init(src)
	tok := s.Scan()
	for tok != scanner.EOF {
		// do something with tok
		tok = s.Scan()
	}
```

## 二、Golang的标准库 tabwriter包

### 1. text/tabwrite
- text/tabwrite
	- tabwriter包实现了写入过滤器 `tabwriter.Writer`，可以将输入的缩进修正为正确的对齐文本
```go
	// Writer是一个过滤器，会在输入的tab划分的列进行填充，在输出中对齐它们
	type Writer struct { ... }

	// 创建并初始化一个tabwriter.Writer，参数用法和Init函数类似
	func NewWriter(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer

	// 初始化一个Writer，第一个参数指定格式化后的输出目标，其余的参数控制格式化
	// minwidth 最小单元长度
	// tabwidth tab字符的宽度
	// padding  计算单元宽度时会额外加上它
	// padchar  用于填充的ASCII字符，如果是'\t'，则Writer会假设tabwidth作为输出中tab的宽度，且单元必然左对齐
	// flags    格式化控制
	func (b *Writer) Init(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer
	// 将buf写入b，实现io.Writer接口，只有在写入底层输出流是才可能发生并返回错误
	func (b *Writer) Write(buf []byte) (n int, err error)
	// 在最后一次调用Write后，必须调用Flush方法以清空缓存，并将格式化对齐后的文本写入生成时提供的output中
	func (b *Writer) Flush() (err error)

	// For Example
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "a\tb\tc\td\t.")
	fmt.Fprintln(w, "123\t12345\t1234567\t123456789\t.")
	fmt.Fprintln(w)
	w.Flush()
	// Format right-aligned in space-separated columns of minimal width 5
	// and at least one blank of padding (so wider column entries do not
	// touch each other).
	w.Init(os.Stdout, 5, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "a\tb\tc\td\t.")
	fmt.Fprintln(w, "123\t12345\t1234567\t123456789\t.")
	fmt.Fprintln(w)
	w.Flush()
	// Output:
	// a	b	c	d		.
	// 123	12345	1234567	123456789	.
	//     a     b       c         d.
	//   123 12345 1234567 123456789.
```

## 三、Golang的标准库 template包

### 1. text/template
- text/template
	- template包实现了数据驱动的用于生成文本输出的模板
	- 如果要生成HTML格式的输出，参见html/template包，该包提供了和本包相同的接口，但会自动将输出转化为安全的HTML格式输出，可以抵抗一些网络攻击
	- 通过将模板应用于一个数据结构(即该数据结构作为模板的参数)来执行，来获得输出
		- 模板中的注释引用数据接口的元素(一般如结构体的字段或者字典的键)来控制执行过程和获取需要呈现的值
		- 模板执行时会遍历结构并将指针表示为'.'(称之为"dot")指向运行过程中数据结构的当前位置的值
	- 用作模板的输入文本必须是utf-8编码的文本
		- "Action" — 数据运算 和 控制单位 — 由"{{"和"}}"界定
		- 在Action之外的所有文本都不做修改的拷贝到输出中
		- Action内部不能有换行，但注释可以有换行
	- 经解析生成模板后，一个模板可以安全的并发执行
```go
	// For Example
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} of {{.Material}}")
	if err != nil { panic(err) }
	err = tmpl.Execute(os.Stdout, sweaters)
	if err != nil { panic(err) }
```

- Actions
	 - 下面是一个action(动作)的列表
	 - "Arguments"和"pipelines"代表数据的执行结果，细节定义如下
```
	{{/* a comment */}}
		注释，执行时会忽略。可以多行。注释不能嵌套，并且必须紧贴分界符始止，就像这里表示的一样。
	{{pipeline}}
		pipeline的值的默认文本表示会被拷贝到输出里。
	{{if pipeline}} T1 {{end}}
		如果pipeline的值为empty，不产生输出，否则输出T1执行结果。不改变dot的值。
		Empty值包括false、0、任意nil指针或者nil接口，任意长度为0的数组、切片、字典。
	{{if pipeline}} T1 {{else}} T0 {{end}}
		如果pipeline的值为empty，输出T0执行结果，否则输出T1执行结果。不改变dot的值。
	{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
		用于简化if-else链条，else action可以直接包含另一个if；等价于：
			{{if pipeline}} T1 {{else}}{{if pipeline}} T0 {{end}}{{end}}
	{{range pipeline}} T1 {{end}}
		pipeline的值必须是数组、切片、字典或者通道。
		如果pipeline的值其长度为0，不会有任何输出；
		否则dot依次设为数组、切片、字典或者通道的每一个成员元素并执行T1；
		如果pipeline的值为字典，且键可排序的基本类型，元素也会按键的顺序排序。
	{{range pipeline}} T1 {{else}} T0 {{end}}
		pipeline的值必须是数组、切片、字典或者通道。
		如果pipeline的值其长度为0，不改变dot的值并执行T0；否则会修改dot并执行T1。
	{{template "name"}}
		执行名为name的模板，提供给模板的参数为nil，如模板不存在输出为""。
	{{template "name" pipeline}}
		执行名为name的模板，提供给模板的参数为pipeline的值。
	{{with pipeline}} T1 {{end}}
		如果pipeline为empty不产生输出，否则将dot设为pipeline的值并执行T1。不修改外面的dot。
	{{with pipeline}} T1 {{else}} T0 {{end}}
		如果pipeline为empty，不改变dot并执行T0，否则dot设为pipeline的值并执行T1。
```

- Arguments
	- 参数代表一个简单的，由下面的某一条表示的值
		- go语法的布尔值、字符串、字符、整数、浮点数、虚数、复数，视为无类型字面常数，字符串不能跨行
		- 关键字nil，代表一个go的无类型的nil值
		- 字符'.'（句点，用时不加单引号），代表dot的值
		- 变量名，以美元符号起始加上（可为空的）字母和数字构成的字符串，如：$piOver2和$；执行结果为变量的值，变量参见下面的介绍
		- 结构体数据的字段名，以句点起始，如：.Field；
			- 执行结果为字段的值，支持链式调用：.Field1.Field2；
			- 字段也可以在变量上使用（包括链式调用）：$x.Field1.Field2；
		- 字典类型数据的键名；以句点起始，如：.Key；
			- 执行结果是该键在字典中对应的成员元素的值；
			- 键也可以和字段配合做链式调用，深度不限：.Field1.Key1.Field2.Key2；
			- 虽然键也必须是字母和数字构成的标识字符串，但不需要以大写字母起始；
			- 键也可以用于变量（包括链式调用）：$x.key1.key2；
		- 数据的无参数方法名，以句点为起始，如：.Method；
			- 执行结果为dot调用该方法的返回值，dot.Method()；
			- 该方法必须有1到2个返回值，如果有2个则后一个必须是error接口类型；
			- 如果有2个返回值的方法返回的error非nil，模板执行会中断并返回给调用模板执行者该错误；
			- 方法可和字段、键配合做链式调用，深度不限：.Field1.Key1.Method1.Field2.Key2.Method2；
			- 方法也可以在变量上使用（包括链式调用）：$x.Method1.Field；
		- 无参数的函数名，如：fun；
			- 执行结果是调用该函数的返回值fun()；对返回值的要求和方法一样；函数和函数名细节参见后面
		- 上面某一条的实例加上括弧（用于分组）,执行结果可以访问其字段或者键对应的值：
			- print (.F1 arg1) (.F2 arg2)
			- (.StructValuedMethod "arg").Field

- Pipelines
	- pipeline通常是将一个command序列分割开，再使用管道符'|'连接起来(但不使用管道符的command序列也可以视为一个管道)
	- 在一个链式的pipeline里，每个command的结果都作为下一个command的最后一个参数
	- pipeline最后一个command的输出作为整个管道执行的结果

- Examples
	- 面是一些单行模板，展示了pipeline和变量
	- 所有都生成加引号的单词"output"
```
	{{"\"output\""}}
		字符串常量
	{{`"output"`}}
		原始字符串常量
	{{printf "%q" "output"}}
		函数调用
	{{"output" | printf "%q"}}
		函数调用，最后一个参数来自前一个command的返回值
	{{printf "%q" (print "out" "put")}}
		加括号的参数
	{{"put" | printf "%s%s" "out" | printf "%q"}}
		玩出花的管道的链式调用
	{{"output" | printf "%s" | printf "%q"}}
		管道的链式调用
	{{with "output"}}{{printf "%q" .}}{{end}}
		使用dot的with action
	{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
		创建并使用变量的with action
	{{with $x := "output"}}{{printf "%q" $x}}{{end}}
		将变量使用在另一个action的with action
	{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
		以管道形式将变量使用在另一个action的with action  
```

- Functions
	- 执行模板时，函数从两个函数字典中查找：首先是模板函数字典，然后是全局函数字典
	- 一般不在模板内定义函数，而是使用Funcs方法添加函数到模板里
	- 预定义的全局函数
```
	and
		函数返回它的第一个empty参数或者最后一个参数；
		就是说"and x y"等价于"if x then y else x"；所有参数都会执行；
	or
		返回第一个非empty参数或者最后一个参数；
		亦即"or x y"等价于"if x then x else y"；所有参数都会执行；
	not
		返回它的单个参数的布尔值的否定
	len
		返回它的参数的整数类型长度
	index
		执行结果为第一个参数以剩下的参数为索引/键指向的值；
		如"index x 1 2 3"返回x[1][2][3]的值；每个被索引的主体必须是数组、切片或者字典。
	print
		即fmt.Sprint
	printf
		即fmt.Sprintf
	println
		即fmt.Sprintln
	html
		返回其参数文本表示的HTML逸码等价表示。
	urlquery
		返回其参数文本表示的可嵌入URL查询的逸码等价表示。
	js
		返回其参数文本表示的JavaScript逸码等价表示。
	call
		执行结果是调用第一个参数的返回值，该参数必须是函数类型，其余参数作为调用该函数的参数；
		如"call .X.Y 1 2"等价于go语言里的dot.X.Y(1, 2)；
		其中Y是函数类型的字段或者字典的值，或者其他类似情况；
		call的第一个参数的执行结果必须是函数类型的值（和预定义函数如print明显不同）；
		该函数类型值必须有1到2个返回值，如果有2个则后一个必须是error接口类型；
		如果有2个返回值的方法返回的error非nil，模板执行会中断并返回给调用模板执行者该错误；
```

## 四、Golang的标准库 parse包

### 1. text/template/parse
- text/template/parse
	- parse包 由 text/template包 和 html/template包定义的模版，构建解析树
	- 客户端应用这两个包来构造模版，parse包提供了不用于一般用途的共享内部数据结构
