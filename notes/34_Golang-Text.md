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
```


## 二、Golang的标准库 tabwriter包

### 1. text/tabwrite
- text/tabwrite
	- tabwriter包实现了写入过滤器 `tabwriter.Writer`，可以将输入的缩进修正为正确的对齐文本


## 三、Golang的标准库 template包

### 1. text/template
- text/template
	- template包实现了数据驱动的用于生成文本输出的模板
	- 如果要生成HTML格式的输出，参见html/template包，该包提供了和本包相同的接口，但会自动将输出转化为安全的HTML格式输出，可以抵抗一些网络攻击


## 四、Golang的标准库 parse包

### 1. text/template/parse
- text/template/parse
	- parse包 由 text/template包 和 html/template包定义的模版，构建解析树
	- 客户端应用这两个包来构造模版，parse包提供了不用于一般用途的共享内部数据结构