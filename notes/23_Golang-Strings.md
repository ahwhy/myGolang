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
```
	s[0] = 'L' // compile error: cannot assign to s[0]
```
	- 不变性意味如果两个字符串共享相同的底层数据的话也是安全的，这使得复制任何长度的字符串代价是低廉的
	- 同样，一个字符串s和对应的子字符串切片s[7:]的操作也可以安全地共享相同的内存，因此字符串切片操作代价也是低廉的
	- 在这两种情况下都没有必要分配新的内存

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
	

## 二、ASCII && Unicode && UTF-8 
- ASCII编码
	- 英文和数字
	- [Go语言字符串的字节长度和字符个数](https://blog.csdn.net/qq_39397165/article/details/116178566)

- Unicode 
	- 称为Unicode字符集或者万国码, 就是将全球所有语言的字符通过编码
	- 所有语言都统一到一套编码，本质就是一张大的码表
	- 比如 `104 -> h`，`101 ->e` (数字 -> 字符 的映射机制，兼容assicc码)，即利用一个数字即可表示一个字符

- UTF-8
	- 目前互联网上使用最广泛的一种Unicode编码方式，最大特点就是可变长
	- 可以使用多个字节表示一个字符，根据字符的不同变换长度
	- UTF-8编码中，一个英文为一个字节，一个中文为三个字节
		- utf8使用变长字节编码，来表示这些unicode码
		- 编码规则如下
			- 如果只有一个字节则其最高二进制位为0
			- 如果是多字节，其第一个字节从最高位开始，连续的二进制位值为1的个数决定了其编码的位数，其余各字节均以10开头
			- UTF-8最多可用到6个字节
				- 如表
					1字节 0xxxxxxx
					2字节 110xxxxx 10xxxxxx 
					3字节 1110xxxx 10xxxxxx 10xxxxxx 
					4字节 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx 
					5字节 111110xx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx 
					6字节 1111110x 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx
				- ascii码 本来就是7个bit表示，所以完全兼容
	- Go语言
		- Go语言里的字符串的内部实现使用UTF8编码. 默认rune类型
		- GO中 uint8(byte类型 ascll) -> 0~127; int32(rune类型 utf-8) -> 128~0x10ffff;  
		- Ascll使用下标遍历，Unicode使用for range遍历
		- len只能表示字符串的 ASCII字符的个数或者字节长度
		- 使用 + 拼接多个字符串，支持换行

## 三、Golang的标准库strings包
- 用于实现字符串的一些常规操作

- 字符串比较
	- Compare 函数
		- 用于比较两个字符串的大小，如果两个字符串相等，返回为 0
		- 如果 a 小于 b ，返回 -1 ，反之返回 1
		- 不推荐使用这个函数，直接使用 == != > < >= <= 等一系列运算符更加直观
		- `func Compare(a, b string) int`
	- EqualFold 函数
		- 计算 s 与 t 忽略字母大小写后是否相等
		- `func EqualFold(s, t string) bool`

- 是否存在某个字符或子串
	- 子串 substr 在 s 中，返回 true
		- `func Contains(s, substr string) bool`
	- chars 中任何一个 Unicode 代码点在 s 中，返回 true 
		- `func ContainsAny(s, chars string) bool`
	- Unicode 代码点 r 在 s 中，返回 true 
		- `func ContainsRune(s string, r rune) bool`

- 子串出现次数
	- 在 Go 中，查找子串出现次数即字符串模式匹配，Count 函数的签名 `func Count(s, sep string) int`

- 字符切分
	- 通过分隔符来切割字符串
		- 带 N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数
			- 当 n < 0 时，返回所有的子字符串
			- 当 n == 0 时，返回的结果是 nil
			- 当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中，最后一个元素不会分割
		- 这4个函数都是通过genSplit内部函数来实现的, 通过 sep 进行分割，返回[]string
		- 如果 sep 为空，相当于分成一个个的 UTF-8 字符，如 `Split("abc","")`，得到的是`[a b c]`
		- `func genSplit(s, sep string, sepSave, n int) []string `
```go
	func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }               // Split 会将 s 中的 sep 去掉，而 SplitAfter 会保留 sep
	func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
	func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) } 
	func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
```

- 判断前缀和后缀
	- s 中是否以 prefix 开始
		- `func HasPrefix(s, prefix string) bool`
	- s 中是否以 suffix 结尾
		- `func HasSuffix(s, suffix string) bool`

- 字符串拼接
	- '+' 用加号连接
	- `func fmt.Sprintf(format string, a ...interface{}) string`
	- 将字符串数组(或 slice)连接起来可以通过 Join 实现
		- `func Join(a []string, sep string) string`
	- 拼接性能较高
		- strings.Builder 
		- bytes.Buffer

- 计算子串位置
```go
	// 查询子串的开始Index的函数有
	func Index(s, sep string) int                   // 在 s 中查找 sep 的第一次出现，返回第一次出现的索引
	func LastIndex(s, substr string) int
	func IndexByte(s string, c byte) int            // 在 s 中查找字节 c 的第一次出现，返回第一次出现的索引
	func IndexAny(s, chars string) int              // chars 中任何一个 Unicode 代码点在 s 中首次出现的位置
	func IndexRune(s string, r rune) int            // Unicode 代码点 r 在 s 中第一次出现的位置

	// 查找字串的结束Index的函数
	// 有三个对应的查找最后一次出现的位置
	func LastIndex(s, sep string) int
	func LastIndexByte(s string, c byte) int
	func LastIndexAny(s, chars string) int
	func LastIndexFunc(s string, f func(rune) bool) int
```

- 子串Count
	- `func Repeat(s string, count int) string`

- 字符和字符串替换
	- 字符替换: Map
		- 将 s 的每一个字符按照 mapping 的规则做映射替换，如果 mapping 返回值 <0 ，则舍弃该字符
		- 该方法只能对每一个字符做处理，但处理方式很灵活，可以方便的过滤，筛选汉字等
			- `func Map(mapping func(rune) rune, s string) string`
	- 字符串替换
		- 用 new 替换 s 中的 old，一共替换 n 个
		- 如果 n < 0，则不限制替换次数，即全部替换
			- `func Replace(s, old, new string, n int) string`
		- 该函数内部直接调用了函数 `Replace(s, old, new , -1)`
			- `func ReplaceAll(s, old, new string) string`

- 大小写转换
```go
	// ToLower，ToUpper 用于大小写转换
	func ToLower(s string) string
	func ToUpper(s string) string

	// ToLowerSpecial，ToUpperSpecial 可以转换特殊字符的大小写
	func ToLowerSpecial(c unicode.SpecialCase, s string) string 
	func ToUpperSpecial(c unicode.SpecialCase, s string) string
```

- 剔除子串
```go
	func Trim(s string, cutset string) string              // 将 s 左侧和右侧中匹配 cutset 中的任一字符的字符去掉
	func TrimLeft(s string, cutset string) string          // 将 s 左侧的匹配 cutset 中的任一字符的字符去掉
	func TrimRight(s string, cutset string) string         // 将 s 右侧的匹配 cutset 中的任一字符的字符去掉
	func TrimPrefix(s, prefix string) string               // 如果 s 的前缀为 prefix 则返回去掉前缀后的 string , 否则 s 没有变化。
	func TrimSuffix(s, suffix string) string               // 如果 s 的后缀为 suffix 则返回去掉后缀后的 string , 否则 s 没有变化。
	func TrimSpace(s string) string                        // 将 s 左侧和右侧的间隔符去掉。常见间隔符包括: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL)
	func TrimFunc(s string, f func(rune) bool) string      // 将 s 左侧和右侧的匹配 f 的字符去掉
	func TrimLeftFunc(s string, f func(rune) bool) string  // 将 s 左侧的匹配 f 的字符去掉
	func TrimRightFunc(s string, f func(rune) bool) string // 将 s 右侧的匹配 f 的字符去掉
```

## 五、其他
- Go语言源代码始终为UTF-8

- Go语言的字符串可以包含任意字节，字符底层是一个只读的byte数组

- Go语言中字符串可以进行循环，使用下表循环获取的acsii字符，使用range循环获取的unicode字符

- Go语言中提供了rune类型用来区分字符值和整数值，一个值代表的就是一个Unicode字符

- Go语言中获取字符串的字节长度使用len()函数，获取字符串的字符个数使用utf8.RuneCountInString函数或者转换为rune切片求其长度，这两种方法都可以达到预期结果。
