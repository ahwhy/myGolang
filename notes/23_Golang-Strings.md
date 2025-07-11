# Golang-Strings  Golang的字符串

## 一、Golang的字符串

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
```
	s1 := "My name is 小明😀"                         // 字符串里可以包含任意Unicode字条
	s2 := "He say:\"i'm fine.\" \n\\Thank\tyou.\\"    // 包含转义字符
	s3 := `There is first line.
	
	There is third line.`                             // 反引号里的转义字符无效，反引号里的内容原封不动地进行输出，包括空白符和换行符
```

- 字符串是不可修改的，修改字符串内部数据的操作也是被禁止的
	- 不变性意味如果两个字符串共享相同的底层数据的话也是安全的，这使得复制任何长度的字符串代价是低廉的
	- 同样，一个字符串s和对应的子字符串切片s[7:]的操作也可以安全地共享相同的内存，因此字符串切片操作代价也是低廉的
	- 在这两种情况下都没有必要分配新的内存
```
	s[0] = 'L' // compile error: cannot assign to s[0]
```

- `+` 操作符将两个字符串链接构造一个新字符串

- `byte` 和 `rune`
	- string 中每个元素叫 "字符"，字符有两种
		- byte: 兼容 ASCLL 码的字符，是 byte 类型，即 uint8 的别名，占用 1 个字节
		- rune: 汉字等字符，unicode，是 rune 类型，即 int32 的别名，占用 4 个字节
	- string 底层是 byte 数组，string 的长度就是该 byte 数组的长度，UTF-8 编码下一个汉字占 3 个 byte，即一个汉字占 3 个长度
		- UTF-8 为目前互联网上使用最广泛的一种 Unicode 的编码方式，最大特点就是可变长
		- 用`len()`可以查询 string 长度
		- golang 中 len 只看字节数，其他语言 len 是看字符数
	- string 可以转换为 `[]byte` 或 `[]rune` 类型
```golang
	// rune 类型字面量 rune int32
	// 单引号定义，使用的是 Unicode 码点(整数)
	// 只能放一个字符，多了报错：more than one character in rune literal
	// 没有长度，不是容器，所以不能用 len
	// 报错：invalid argument: s1 (variable of type rune) for len
	s1 := 'a'         // rune 4bytes 
	var s2 byte = 'a' // 1bytes
	s3 := '0'         // 表示 0 的字节
	s4 := '测'        // rune 4bytes unicode 

	// 字符串字面量 string
	// 线性数据结构，可以索引
	// UTF-8 编码的字节序列，不可变
	// 双引号定义
	s5 := "abc"          // 3 len，3 bytes，616263(ASCLL)
	s6 := "测试"          // 6 len，6 bytes，UTF-8
	fmt.Println(s5, s6)  // 3 6
	
	// string 有序字节序列
	t1 := []byte(s5)                                   // []byte(x) 强制类型转换; []byte{x} byte切片字面量的定义
	fmt.Println(t1, len(t1), cap(t1), &t1[0], &t1[1])  // [97 98 99] 3 8 0xc0000a6018 0xc0000a6019 偏移量offset 1byte
	t2 := []rune(s5)
	fmt.Println(t2, len(t2), cap(t2), &t2[0], &t2[1])  // [97 98 99] 3 4 0xc0000a6040 0xc0000a6044 偏移量offset 4byte

	t3 := []byte(s6)                                   // UTF-8 编码的字节序列
	fmt.Println(t3, len(t3))                           // [230 181 139 232 175 149] 6
	t4 := []rune(s6)                                   // 强制转换为 rune，即 [rune rune] Unicode; UTF-8序列 --> Unicode序列 [27979 35797]
	fmt.Println(t4, len(t4))                           // [27979 35797] 2 

	fmt.Println(s5[0], s5[1], s5[2])                   // []byte，其中一个byte uint8  输出结果:  97 98 99
	fmt.Println(string([]byte{0x61, '\x62', 0x63}))    // byte sequence --> UTF-8 string sequence  输出结果: abc

	fmt.Println(s6[0], s6[1], s6[2])                   // []byte，其中一个byte uint8  输出结果:  230 181 139
	fmt.Println(string(27979), string(35797))          // string(Unicode码)  输出结果: 测 试
	fmt.Println(string([]rune{27979, 35797}))          // Unicode序列 --> UTF-8序列 string  输出结果: 测试
	fmt.Println([]byte{230 181 139 232 175 149})       // 输出结果: 测试
``` 

- 强制类型转换
	- `byte`  和 `int` 可以相互转换
	- `float` 和 `int` 可以相互转换，小数位会丢失(float到int会丢失精度)
	- `boot`  和 `int` 不可以相互转换
	- 不同长度的 `int` 和 `float` 之间可以相互转换
	- `string` 可以转换为 `[]byte` 或 `[]rune` 类型，`byte` 或 `rune` 类型可以转换为 `string`
	- `boot` 和 `int` 不能相互转换
	- 低精度向高精度转换没有问题，高精度向低精度转换会丢失位数
	- 无符号向有符号转换，最高位是无符号

- 单引号、双引号、反引号
	- 单引号，表示byte类型或rune类型，对应 uint8和int32类型，默认是 rune 类型
		- byte用来强调数据是raw data，而不是数字
		- 而rune用来表示Unicode的code point
	- 双引号，才是字符串，实际上是字符数组
		- 可以用索引号访问某字节
		- 也可以用len()函数来获取字符串所占的字节长度
	- 反引号，表示字符串字面量，但不支持任何转义序列
		- 字面量 raw literal string 的意思是，定义时写的啥样，它就啥样
		- 有换行，它就换行
		- 如果写转义字符，它也就展示转义字符。


## 二、ASCII && Unicode && UTF-8

- ASCII编码
	- ASCII (American Standard Code for Information Interchange)
	- 美国信息交换标准代码是基于拉丁字母的一套电脑编码系统，主要用于显示现代英语和其他西欧语言，即英文和数字
	- [Go语言字符串的字节长度和字符个数](https://blog.csdn.net/qq_39397165/article/details/116178566)
	- 常见特殊字符

	|序号|转义字符|十进制数|说明|
	|:------:|:------:|:------:|:------:|
	|1|\x00|0|null字符，表中第一项，C语言中的字符串结束符|
	|2|\x09 \t|9|tab字符|
	|3|\x0d\x0a \r\n|13 10|回车和换行|
	|4|\x30~\x39|48～57|字符0～9|
	|5|\x31|49|字符1|
	|6|\x41|65|字符A|
	|7|\x61|97|字符a|

- Unicode 
	- 一个字符集标准，它为全球所有语言的字符分配唯一的数字(称为 码点，Code Point)
		- 'A' 对应码点 U+0041
		- '汉' 对应码点 U+6C49
		- `104 -> h`，`101 ->e` (数字 -> 字符 的映射机制，兼容ASCII编码)，即利用一个数字即可表示一个字符
	- 所有语言都统一到一套编码，本质就是一张大的码表
	- Unicode 本身不涉及编码方式，它只是定义了字符与码点的映射关系

- UTF-8
	- UTF-8 是 Unicode 的一种 编码方式(Encoding)，它定义了如何将 Unicode 码点转换为二进制字节序列
	- 目前互联网上使用最广泛的一种 Unicode 编码方式，特点就是兼容 ASCII 以及可变长
	- 向后兼容 ASCII(ASCII 字符在 UTF-8 中仍用单字节表示)
	- UTF-8编码中，一个英文为一个字节，一个中文为三个字节
		- UTF-8使用变长字节编码，来表示这些Unicode码
		- 可以使用多个字节表示一个字符，根据字符的不同变换长度，一个 Unicode 码点在 UTF-8 中可能占用 1~4 个字节。
		- 编码规则如下
			- 如果只有一个字节则其最高二进制位为0
			- 如果是多字节，其第一个字节从最高位开始，连续的二进制位值为1的个数决定了其编码的位数，其余各字节均以10开头
			- UTF-8最多可用到6个字节
				- 如表
					|1字节|0xxxxxxx|
					|2字节|110xxxxx|10xxxxxx|
					|3字节|1110xxxx|10xxxxxx|10xxxxxx|
					|4字节|11110xxx|10xxxxxx|10xxxxxx|10xxxxxx|
					|5字节|111110xx|10xxxxxx|10xxxxxx|10xxxxxx|10xxxxxx|
					|6字节|1111110x|10xxxxxx|10xxxxxx|10xxxxxx|10xxxxxx|10xxxxxx|
				- ascii码 本来就是7个bit表示，所以完全兼容
			- 例如：'A'(U+0041)在 UTF-8 中占 1 字节；'汉'(U+6C49)占 3 字节
	- Go语言
		- 在 Golang 中，字符串以 UTF-8 编码存储，但处理 Unicode 字符时需使用 rune 类型和 utf8 包
		- `uint8`(byte类型 ASCII) -> 0~127
		- `int32`(rune类型 UTF-8) -> 128~0x10ffff;  
		- Ascll使用下标遍历，Unicode使用`for range`遍历
		- len只能表示字符串的 ASCII字符 的个数或者字节长度
		- 使用 `+` 拼接多个字符串，支持换行


## 三、Golang的标准库strings包

- strings包
	- 实现了用于操作字符的简单函数
### 1. 字符串查询
- 字符串比较
	- Compare 函数
		- 用于比较两个字符串的大小，如果两个字符串相等，返回为 0
		- 如果 a 小于 b ，返回 -1 ，反之返回 1
		- 不推荐使用这个函数，直接使用 == != > < >= <= 等一系列运算符更加直观
	- EqualFold 函数
		- 判断两个utf-8编码字符串(将unicode大写、小写、标题三种格式字符视为相同)是否相同
```go
	func Compare(a, b string) int       // Compare
	func EqualFold(s, t string) bool    // EqualFold
```

- 判断前缀和后缀
```go
	func HasPrefix(s, prefix string) bool   // 判断 s 是否有前缀字符串 prefix
	func HasSuffix(s, suffix string) bool   // 判断 s 是否有后缀字符串 suffix
```

- 是否存在某个字符或子串
```go
	func Contains(s, substr string) bool       // 子串 substr 在 s 中，返回 true
	func ContainsAny(s, chars string) bool     // chars 中任何一个 Unicode 代码点在 s 中，返回 true 
	func ContainsRune(s string, r rune) bool   // Unicode 代码点 r 在 s 中，返回 true 
```

- 子串出现次数
```go
	// Count 查找子串出现次数即字符串模式匹配
	func Count(s, sep string) int
```

- 计算子串位置
```go
	// 查询子串的开始Index的函数有
	func Index(s, sep string) int                     // 在 s 中查找 sep 的第一次出现，返回第一次出现的索引，不存在则返回-1
	func IndexByte(s string, c byte) int              // 在 s 中查找字节 c 的第一次出现，返回第一次出现的索引
	func IndexAny(s, chars string) int                // chars 中任何一个 Unicode 代码点在 s 中首次出现的位置
	func IndexRune(s string, r rune) int              // Unicode 代码点 r 在 s 中第一次出现的位置
	func IndexFunc(s string, f func(rune) bool) int   // s 中第一个满足函数 f 的位置 i (该处的utf-8码值r满足f(r)==true)

	// 查找字串的结束Index的函数
	// 有三个对应的查找最后一次出现的位置
	func LastIndex(s, sep string) int
	func LastIndexByte(s string, c byte) int
	func LastIndexAny(s, chars string) int
	func LastIndexFunc(s string, f func(rune) bool) int
```

### 2. 字符串替换
- 字符串大小写转换
```go
	// 返回s中每个单词的首字母都改为标题格式的字符串拷贝
	// strings.Title("her royal highness")
	// Output: Her Royal Highness
	func Title(s string) string
	// 返回s中所有字母都转为对应的标题版本的拷贝
	func ToTitle(s string) string
	// 使用_case规定的字符映射，返回s中所有字母都转为对应的标题版本的拷贝
	func ToTitleSpecial(_case unicode.SpecialCase, s string) string

	// ToLower，ToUpper 用于大小写转换
	func ToLower(s string) string
	func ToUpper(s string) string

	// ToLowerSpecial，ToUpperSpecial 可以转换特殊字符的大小写
	func ToLowerSpecial(c unicode.SpecialCase, s string) string 
	func ToUpperSpecial(c unicode.SpecialCase, s string) string
```

- `strings.Replace`
```go
	// 返回将s中前n个不重叠old子串都替换为new的新字符串，如果 n<0 会替换所有old子串
	func Replace(s, old, new string, n int) string
	// 替换所有old子串，该函数内部直接调用了函数 Replace
	func ReplaceAll(s, old, new string) string

	// 示例 Example
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
	// Output:
	// oinky oinky oink
	// moo moo moo
```

- `strings.Map`
	- 将 s 的每一个字符按照 mapping 的规则做映射替换，如果 mapping 返回值 <0 ，则舍弃该字符
	- 该方法只能对每一个字符做处理，但处理方式很灵活，可以方便的过滤，筛选汉字等
```go
	// 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝
	// 如果mapping返回一个负值，将会丢弃该码值而不会被替换 (如果mapping返回一个负值，将会丢弃该码值而不会被替换)
	func Map(mapping func(rune) rune, s string) string

	// 示例 Example
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))
	// Output: 'Gjnf oevyyvt naq gur fyvgul tbcure...
```

### 3. 字符串剔除
- 剔除子串
```go
	func Trim(s string, cutset string) string              // 将 s 左侧和右侧中匹配 cutset 中的任一字符(包含的utf-8码值)的字符去掉
	func TrimLeft(s string, cutset string) string          // 将 s 左侧的匹配 cutset 中的任一字符的字符去掉
	func TrimRight(s string, cutset string) string         // 将 s 右侧的匹配 cutset 中的任一字符的字符去掉
	func TrimPrefix(s, prefix string) string               // 如果 s 的前缀为 prefix 则返回去掉前缀后的 string , 否则 s 没有变化。
	func TrimSuffix(s, suffix string) string               // 如果 s 的后缀为 suffix 则返回去掉后缀后的 string , 否则 s 没有变化。
	func TrimSpace(s string) string                        // 将 s 左侧和右侧的间隔符去掉。常见间隔符包括: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL)
	func TrimFunc(s string, f func(rune) bool) string      // 将 s 左侧和右侧的匹配 f 的字符去掉
	func TrimLeftFunc(s string, f func(rune) bool) string  // 将 s 左侧的匹配 f 的字符去掉
	func TrimRightFunc(s string, f func(rune) bool) string // 将 s 右侧的匹配 f 的字符去掉
```

### 4. 字符串分割
- 通过空白字符来分割字符串
```go
	// 返回将字符串按照空白(unicode.IsSpace确定，可以是一到多个连续的空白字符)分割的多个字符串
	// 如果字符串全部是空白或者是空字符串的话，会返回空切片
	func Fields(s string) []string
	// 类似Fields，但使用函数f来确定分割符(满足f的unicode码值)
	func FieldsFunc(s string, f func(rune) bool) []string
```

- 通过分隔符来分割字符串
	- 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片
	- 每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割
	- 带 N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数
		- 当 n < 0 时，返回所有的子字符串
		- 当 n == 0 时，返回的结果是 nil
		- 当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中 最后一个元素不会分割
	- 这4个函数都是通过genSplit内部函数来实现的, 通过 sep 进行分割，返回 `[]string`
	- 如果 sep 为空，相当于分成一个个的 UTF-8 字符，如 `Split("abc","")`，得到的是`[a b c]`
	- `func genSplit(s, sep string, sepSave, n int) []string`
```go
	func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }               // Split 会将 s 中的 sep 去掉，而 SplitAfter 会保留 sep
	func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
	func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) } 
	func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
```

### 5. 字符串拼接
- 字符串拼接
	- '+' 用加号连接
	- `func fmt.Sprintf(format string, a ...interface{}) string`
	- 将字符串数组(或 slice)连接起来可以通过 Join 实现
		- `func Join(a []string, sep string) string` 将一系列字符串连接为一个字符串，之间用sep来分隔
	- 拼接性能较高
		- `strings.Builder`
		- `bytes.Buffer`
```golang
	var builder strings.Builder // 构建器
	builder.Write([]byte(s1))
	builder.WriteByte('-')
	builder.WriteString(s2)
	s3 := builder.String()
```

- 子串Count
	- 返回count个s串联的字符串
```go
	// "ba" + strings.Repeat("na", 2)
	// Output: banana
	func Repeat(s string, count int) string
```

### 6. Strings包中的Reader
- `Strings.Reader`
	- Reader 类型通过从一个字符串读取数据，实现了io.Reader、io.Seeker、io.ReaderAt、io.WriterTo、io.ByteScanner、io.RuneScanner接口
```go
	type Reader struct { ... }

	// NewReader创建一个从s读取数据的Reader
	func NewReader(s string) *Reader

	// Len返回r包含的字符串还没有被读取的部分
	func (r *Reader) Len() int

	func (r *Reader) Read(b []byte) (n int, err error)
	func (r *Reader) ReadByte() (b byte, err error)
	func (r *Reader) UnreadByte() error
	func (r *Reader) ReadRune() (ch rune, size int, err error)
	func (r *Reader) UnreadRune() error
	func (r *Reader) Seek(offset int64, whence int) (int64, error)  // Seek实现了io.Seeker接口
	func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)
	func (r *Reader) WriteTo(w io.Writer) (n int64, err error)      // WriteTo实现了io.WriterTo接口
```


## 五、其他
- Go语言源代码始终为UTF-8

- Go语言的字符串可以包含任意字节，字符底层是一个只读的`byte`数组

- Go语言中字符串可以进行循环，使用下表循环获取的 ASCII字符，使用`for range`循环获取的 Unicode字符 或使用 `utf8.DecodeRuneInString`

- Go语言中提供了`rune`类型用来区分字符值和整数值，一个值代表的就是一个 Unicode字符

- Go语言中获取字符串的字节长度使用`len()`函数，获取字符串的字符个数使用`utf8.RuneCountInString`函数或者转换为 `[]rune` 切片求其长度，这两种方法都可以达到预期结果。