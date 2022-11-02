# Golang-Regexp  Golang的正则表达式

## 一、Golang的标准库 regexp包

### 1. regexp
- regexp
	- regexp包实现了正则表达式搜索
	- 正则表达式采用RE2语法（除了\c、\C），和Perl、Python等语言的正则基本一致
	- 参见 [Syntax](http://code.google.com/p/re2/wiki/Syntax)

### 2. regexp_Syntax
- 本包采用的正则表达式语法，默认采用perl标志；某些语法可以通过切换解析时的标志来关闭

- 单字符
```
	.              任意字符（标志s==true时还包括换行符）
	[xyz]          字符族
	[^xyz]         反向字符族
	\d             Perl预定义字符族
	\D             反向Perl预定义字符族
	[:alpha:]      ASCII字符族
	[:^alpha:]     反向ASCII字符族
	\pN            Unicode字符族（单字符名），参见unicode包
	\PN            反向Unicode字符族（单字符名）
	\p{Greek}      Unicode字符族（完整字符名）
	\P{Greek}      反向Unicode字符族（完整字符名）
```

- 结合
```
	xy             匹配x后接着匹配y
	x|y            匹配x或y（优先匹配x）
```

- 重复
	- 实现的限制：计数格式x{n}等（不包括x*等格式）中n最大值1000
	- 负数或者显式出现的过大的值会导致解析错误，返回ErrInvalidRepeatSize
```
	x*             重复>=0次匹配x，越多越好（优先重复匹配x）
	x+             重复>=1次匹配x，越多越好（优先重复匹配x）
	x?             0或1次匹配x，优先1次
	x{n,m}         n到m次匹配x，越多越好（优先重复匹配x）
	x{n,}          重复>=n次匹配x，越多越好（优先重复匹配x）
	x{n}           重复n次匹配x
	x*?            重复>=0次匹配x，越少越好（优先跳出重复）
	x+?            重复>=1次匹配x，越少越好（优先跳出重复）
	x??            0或1次匹配x，优先0次
	x{n,m}?        n到m次匹配x，越少越好（优先跳出重复）
	x{n,}?         重复>=n次匹配x，越少越好（优先跳出重复）
	x{n}?          重复n次匹配x
```

- 分组
```
	(re)           编号的捕获分组
	(?P<name>re)   命名并编号的捕获分组
	(?:re)         不捕获的分组
	(?flags)       设置当前所在分组的标志，不捕获也不匹配
	(?flags:re)    设置re段的标志，不捕获的分组
```

- 标志
	- 标志的语法为xyz（设置）、-xyz（清楚）、xy-z（设置xy，清楚z）
```
	I              大小写敏感（默认关闭）
	m              ^和$在匹配文本开始和结尾之外，还可以匹配行首和行尾（默认开启）
	s              让.可以匹配\n（默认关闭）
	U              非贪婪的：交换x*和x*?、x+和x+?……的含义（默认关闭）
```

- 边界匹配
```
	^              匹配文本开始，标志m为真时，还匹配行首
	$              匹配文本结尾，标志m为真时，还匹配行尾
	\A             匹配文本开始
	\b             单词边界（一边字符属于\w，另一边为文首、文尾、行首、行尾或属于\W）
	\B             非单词边界
	\z             匹配文本结尾
```

- 转义序列
```
	\a             响铃符（\007）
	\f             换纸符（\014）
	\t             水平制表符（\011）
	\n             换行符（\012）
	\r             回车符（\015）
	\v             垂直制表符（\013）
	\123           八进制表示的字符码（最多三个数字）
	\x7F           十六进制表示的字符码（必须两个数字）
	\x{10FFFF}     十六进制表示的字符码
	\*             字面值'*'
	\Q...\E        反斜线后面的字符的字面值
```

- 字符族
	- 预定义字符族之外，方括号内部
```
	x              单个字符
	A-Z            字符范围（方括号内部才可以用）
	\d             Perl字符族
	[:foo:]        ASCII字符族
	\pF            单字符名的Unicode字符族
	\p{Foo}        完整字符名的Unicode字符族
```

- 预定义字符族作为字符族的元素
```
	[\d]           == \d
	[^\d]          == \D
	[\D]           == \D
	[^\D]          == \d
	[[:name:]]     == [:name:]
	[^[:name:]]    == [:^name:]
	[\p{Name}]     == \p{Name}
	[^\p{Name}]    == \P{Name}
```

- Perl字符族
```
	\d             == [0-9]
	\D             == [^0-9]
	\s             == [\t\n\f\r ]
	\S             == [^\t\n\f\r ]
	\w             == [0-9A-Za-z_]
	\W             == [^0-9A-Za-z_]
```

- ASCII字符族
```
	[:alnum:]      == [0-9A-Za-z]
	[:alpha:]      == [A-Za-z]
	[:ascii:]      == [\x00-\x7F]
	[:blank:]      == [\t ]
	[:cntrl:]      == [\x00-\x1F\x7F]
	[:digit:]      == [0-9]
	[:graph:]      == [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~]
	[:lower:]      == [a-z]
	[:print:]      == [ -~] == [ [:graph:]]
	[:punct:]      == [!-/:-@[-`{-~]
	[:space:]      == [\t\n\v\f\r ]
	[:upper:]      == [A-Z]
	[:word:]       == [0-9A-Za-z_]
	[:xdigit:]     == [0-9A-Fa-f]
```

### 3. regexp_Method
```go
	// QuoteMeta 返回将s中所有正则表达式元字符都进行转义后字符串
	// 该字符串可以用在正则表达式中匹配字面值s；例如，QuoteMeta(`[foo]`)会返回`\[foo\]`
	func QuoteMeta(s string) string
	// Match 检查b中是否存在匹配pattern的子序列；更复杂的用法需使用Compile函数和Regexp对象
	func Match(pattern string, b []byte) (matched bool, err error)
	// MatchString 类似Match，但匹配对象是字符串
	func MatchString(pattern string, s string) (matched bool, err error)
	// MatchReader 类似Match，但匹配对象是io.RuneReader
	func MatchReader(pattern string, r io.RuneReader) (matched bool, err error)

	// Compile 解析并返回一个正则表达式；如果成功返回，该Regexp就可用于匹配文本
	func Compile(expr string) (*Regexp, error)
	// CompilePOSIX 类似Compile但会将语法约束到POSIX ERE（egrep）语法，并将匹配模式设置为leftmost-longest
	func CompilePOSIX(expr string) (*Regexp, error)
	// MustCompile 类似Compile但会在解析失败时panic，主要用于全局正则表达式变量的安全初始化
	func MustCompile(str string) *Regexp
	// MustCompilePOSIX 类似CompilePOSIX但会在解析失败时panic，主要用于全局正则表达式变量的安全初始化
	func MustCompilePOSIX(str string) *Regexp
	// Regexp 代表一个编译好的正则表达式；Regexp可以被多线程安全地同时使用
	type Regexp struct { ... }

	// String 返回用于编译成正则表达式的字符串
	func (re *Regexp) String() string
	// LiteralPrefix 返回一个字符串字面值prefix，任何匹配本正则表达式的字符串都会以prefix起始
	func (re *Regexp) LiteralPrefix() (prefix string, complete bool)
	// NumSubexp 返回该正则表达式中捕获分组的数量
	func (re *Regexp) NumSubexp() int
	// SubexpNames返回该正则表达式中捕获分组的名字
	func (re *Regexp) SubexpNames() []string
	// Longest 让正则表达式在之后的搜索中都采用"leftmost-longest"模式
	func (re *Regexp) Longest()

	// Match 检查b中是否存在匹配pattern的子序列
	func (re *Regexp) Match(b []byte) bool
	// MatchString 类似Match，但匹配对象是字符串
	func (re *Regexp) MatchString(s string) bool
	// MatchReader 类似Match，但匹配对象是io.RuneReader
	func (re *Regexp) MatchReader(r io.RuneReader) bool

	// Find 返回保管正则表达式re在b中的最左侧的一个匹配结果的[]byte切片
	func (re *Regexp) Find(b []byte) []byte
	// FindString 返回保管正则表达式re在b中的最左侧的一个匹配结果的字符串
	func (re *Regexp) FindString(s string) string
	// FindIndex 返回保管正则表达式re在b中的最左侧的一个匹配结果的起止位置的切片（显然len(loc)==2）
	func (re *Regexp) FindIndex(b []byte) (loc []int)
	func (re *Regexp) FindStringIndex(s string) (loc []int)
	func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int)
	// FindSubmatch 返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的[][]byte切片
	func (re *Regexp) FindSubmatch(b []byte) [][]byte
	func (re *Regexp) FindStringSubmatch(s string) []string
	func (re *Regexp) FindSubmatchIndex(b []byte) []int
	func (re *Regexp) FindStringSubmatchIndex(s string) []int
	func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int
	// FindAll 返回保管正则表达式re在b中的所有不重叠的匹配结果的[][]byte切片
	func (re *Regexp) FindAll(b []byte, n int) [][]byte
	// FindAllString 返回保管正则表达式re在b中的所有不重叠的匹配结果的[]string切片
	func (re *Regexp) FindAllString(s string, n int) []string
	// FindAllIndex 返回保管正则表达式re在b中的所有不重叠的匹配结果的起止位置的切片
	func (re *Regexp) FindAllIndex(b []byte, n int) [][]int
	func (re *Regexp) FindAllStringIndex(s string, n int) [][]int
	// FindAllSubmatch 返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的[][][]byte切片
	func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte
	func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string
	// FindAllSubmatchIndex 返回一个保管正则表达式re在b中的所有不重叠的匹配结果及其对应的（可能有的）分组匹配的结果的起止位置的切片（第一层表示第几个匹配结果，完整匹配和分组匹配的起止位置对在第二层）
	func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int
	func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int

	// Split 将re在s中匹配到的结果作为分隔符将s分割成多个字符串，并返回这些正则匹配结果之间的字符串的切片
	func (re *Regexp) Split(s string, n int) []string
	// Expand 返回新生成的将template添加到dst后面的切片
	func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte
	func (re *Regexp) ExpandString(dst []byte, template string, src string, match []int) []byte

	// ReplaceAllLiteral 返回src的一个拷贝，将src中所有re的匹配结果都替换为repl
	func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte
	func (re *Regexp) ReplaceAllLiteralString(src, repl string) string
	func (re *Regexp) ReplaceAll(src, repl []byte) []byte
	func (re *Regexp) ReplaceAllString(src, repl string) string
	// ReplaceAllFunc 返回src的一个拷贝，将src中所有re的匹配结果（设为matched）都替换为repl(matched)
	func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte
	func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string
```

## 二、Golang的标准库 regexp/syntax包

### 1. regexp/syntax
- regexp/syntax
	- syntax包将正则表达式解析成解析树，并将解析树编译成程序
	- 一般使用regexp包的功能

### 2. syntax_Syntax
- 同上