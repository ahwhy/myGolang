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