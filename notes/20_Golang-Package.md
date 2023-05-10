# Golang-Package  Golong的包

## 一、定义
- 包是函数和数据的集合，将有相关特性的函数和数据放在统一的文件目录进行管理，每个包都可以作为独立的单元维护并提供给其他项目进行使用
- 声明所在包，包名告知编译器哪些是包的源代码用于编译库文件，其次包名用于限制包内成员对外的可见性，最后包名用于在包外对公开成员的访问
- 在源文件中加上`package xxx`就可以声明xxx的包
 
## 二、成员可见性
- Go 语言使用名称首字母大小写来判断对象(常量、变量、函数、类型、结构体、方法等)的访问权限，首字母大写标识包外可见(公开的)，否者仅包内可访问(内部的);

## 三、main 包与 main 函数
- main包用于声明告知编译器 将包编译为二进制 可执行文件
- main包中的 main 函数是程序的入口，无返回值，无参数

## 四、init 函数
- init函数是初始化包使用，无返回值，无参数。建议每个包只定义一个； 
- init函数在import包时自动被调用(const -->var -->init)。

## 五、标准包
- Go提供了大量标准包，可查看 https://golang.google.cn/pkg/ &&  https://godoc.org 
	- `go list std` 查看所有标准包
	- `go doc packagename` 查看包的帮助信息
	- `go doc packagename.element` 查看包内成员 帮助信息

## 六、包的维护
- 包的提供者 -> 打tag  -> git tag
- 包的使用者 -> 改版本 -> go mod

## 七、关系说明
- import 导入的是路径，而非包名
- 包名和目录名不强制一致，但推荐一致
- 在代码中引用包的成员变量或者函数时，使用的包名不是目录名
- 在同一目录下，所有的源文件必须使用相同的包名
	- Multiple packages in directory: pk2, pk3 
- 文件名不限制，但不能有中文

## 八、设置 go mod 和 go proxy
- 设置两个环境变量
```shell
	$ go env -w GO111MODULE=on
	$ go env -w GOPROXY=https://goproxy.io,direct
	$ go env -w GOPROXY=https://goproxy.cn
```

- [安装 Go语言扩展需要的工具集](https://gitee.com/infraboard/go-course/blob/master/zh-cn/base/install.md#%E5%AE%89%E8%A3%85go-%E8%AF%AD%E8%A8%80%E6%89%A9%E5%B1%95%E9%9C%80%E8%A6%81%E7%9A%84%E5%B7%A5%E5%85%B7%E9%9B%86)
```shell
	$ go install -v golang.org/x/tools/gopls@latest
	$ go install -v honnef.co/go/tools/cmd/staticcheck@latest
	$ go install -v github.com/go-delve/delve/cmd/dlv@latest
	$ go install -v github.com/haya14busa/goplay/cmd/goplay@latest
	$ go install -v github.com/josharian/impl@latest
	$ go install -v github.com/fatih/gomodifytags@latest
	$ go install -v github.com/cweill/gotests/gotests@latest
	$ go install -v github.com/ramya-rao-a/go-outline@latest
	$ go install -v github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
```

## 九、创建git，发布到github
- 项目目录下 `go mod init github.com/ahwhy/myGolang`
- `git init`
- 添加 .gitignore 文件去掉一些和代码无关的文件/文件夹
- `git add . && git commit -m "Record me learning golang" --author "ahwhya <ahwhya@outlook.com>"`
- github上新建一个仓库
- 推送到远程
- 上传tag
	- 通过tag可以返回到项目的特定状态下，可以将tag看作是在大量commit中设定的书签
	- [git中tag与release的创建以及两者的区别](https://www.jianshu.com/p/79ecf4fe5079)
```shell
	# 推送到远程
	# or push an existing repository from the command line
	git remote add origin https://github.com/ahwhy/myGolang.git
	git branch -M main
	git push -u origin main

	# 上传tag
	git log --oneline
	# 创建lightweight类型的tag
	git tag v0.0.2-lw
	# 创建annotated类型的tag
	git tag -a v0.0.2 -m "updata"
	git push origin --tags
	git push -u origin --tags v0.0.2
```

## 十、常用内建库与函数  
- [golang的标准库](https://studygolang.com/pkgdoc)

- [Go中文开发手册](https://www.php.cn/manual/view/35126.html)

### 1. time
- time包提供了时间的显示和测量用的函数
	- 日历的计算采用的是公历
	- 详见 [Golang-Time](./21_Golang-Time.md)

- time包中的格式转换
	- `layout := "2006-01-02 15:04:05"`
```
	              -> time.Unix(sec int64, nsec int64)             -> time.Format(layout)
	时间戳(Timestamp)                             time.Time                                       日期格式
	                       <- time.Unix()                    <- time.Parse(layout, value string)			
```

- 示例
```go
	time.Now()          // 获取当前时间
	time.Now().Unix()
	time.Now().Year()   // Month() Day()  Hour()  Minute()  Second()
	time.Now().Format("2006-01-02 15:04:05")
	time.Parse("2006-01-02 15:04:05", "2022-08-08 09:36:58")        // 返回转换后的时间格式和一个判断信息（err)
	time.Sleep(1 * time.Second)
	time.Now().Add(30 * time.Second)
	time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
```

### 2. math
- math
	- math包提供了基本的数学常数和数学函数
```go
	// 用于测试一个数是否是非数NaN
	// NaN非数，一般用于表示无效的除法操作结果0/0或Sqrt(-1)
	func IsNaN(f float64) (is bool)

	// 返回非数对应的值，一个IEEE 754 “这不是一个数字” 值
	func NaN() float64

	// 返回x的绝对值
	func Abs(x float64) float64

	// 返回x和y中最大值
	func Max(x, y float64) float64

	// 返回x和y中最小值
	func Min(x, y float64) float64
```

- math/big
	- big包实现了大数字的多精度计算
	- 支持如下数字类型
		- Int	有符号整数
		- Rat	有理数

- math/cmplx
	- cmplx包提供了复数的常用常数和常用函数。

- math/rand
	- rand包实现了伪随机数生成器
	- 伪随机数，用确定性的算法计算出来自[0,1]均匀分布的随机数序列。并不真正的随机，但具有类似于随机数的统计特征，如均匀性、独立性等
	- 在计算伪随机数时，若使用的初值（种子）不变，那么伪随机数的数序也不变
	- 伪随机数可以用计算机大量生成，在模拟研究中为了提高模拟效率，一般采用伪随机数代替真正的随机数
	- 模拟中使用的一般是循环周期极长并能通过随机数检验的伪随机数，以保证计算结果的随机性
```go
	// 使用给定的seed将默认资源初始化到一个确定的状态
	// 如未调用Seed，默认资源的行为就好像调用了Seed(1)
	// 使用当前时间设置随机数种子  rand.Seed(time.Now().Unix())
	func Seed(seed int64)

	// 返回一个非负的伪随机int值
	func Int() int

	// 返回一个取值范围在[0,n)的伪随机int值，如果 n<=0 会 panic
	// 生产[0, 100)的随机数 rand.Intn(100)
	func Intn(n int) int

	// 返回一个取值范围在[0.0, 1.0)的伪随机float64值
	func Float64() float64

	// 返回一个有n个元素的，[0,n)范围内整数的伪随机排列的切片
	func Perm(n int) []int
```
### 3. os
- os
	- os包提供了操作系统函数的不依赖平台的接口
	- os包的接口规定为在所有操作系统中都是一致的
	- 非公用的属性可以从操作系统特定的syscall包获取
	- `os.Hostname()` 获取主机名
	- `os.Getenv(key)` Getenv检索并返回名为key的环境变量的值
	- `os.Getpid()` 返回调用者所在进程的进程ID
	- `os.Getppid()` 返回调用者所在进程的父进程的进程ID
	- `os.Exit(code)` Exit让当前程序以给出的状态码code退出;一般来说，状态码0表示成功，非0表示出错;程序会立刻终止，defer的函数不会被执行
	- `os.Getwd()` 获取当前目录
	- `os.Mkdir("a", os.ModePerm)` 创建文件夹
	- `os.MkdirAll("a/b/c", os.ModePerm)` 创建文件夹(父目录不存在逐层创建)
	- `os.IsExist(err)` 与 `os.Stat` 一起用于判断文件存在
	- `os.IsNotExist(err)` 与 `os.Stat` 一起用于判断文件不存在
	- `os.Link("oldname", "newname")` 创建一个名为newname指向oldname的硬链
	- `os.Create()` 创建文件
	- `os.Open("test.txt")` 读取文件
	- `os.Open().Stat` 获取文件属性 
	- `os.Chmod(name string, mode FileMode)` 修改文件权限
	- `os.Chown(name string, uid, gid int)` 修改文件所属用户，用户组
	- `os.Truncate(name string, size int64)` 修改name指定的文件的大小
	- `os.Rename("a.txt", "b.txt")` 重命名文件
	- `os.Remove("b.txt")` 删除指定的文件或目录
	- `os.RemoveAll("path")` 删除path指定的文件，或目录及它包含的任何下级对象
	- `os.TempDir()` 返回一个用于保管临时文件的默认目录
	- `os.Create(name)` 创建文件并返回文件对象指针(文件不存在则创建，文件存在则清空)
	- `os.Open(name)` 打开文件并返回文件对象指针
	- `os.OpenFile(name, flag, perm)` 按指定权限打开文件，并返回文件指针对象

- os/exec
	- exec包提供了启动一个外部进程并使用标准输入和输出进行通信
	- `exec.Command("date").Output()` 执行命令并返回标准输出的切片

- os/signal
	- signal包实现了对输入信号的访问
```go
	// Notify函数让signal包将输入信号转发到c
	// 如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
	func Notify(c chan<- os.Signal, sig ...os.Signal)

	// Stop 让signal包停止向c转发信号
	// 它会取消之前使用c调用的所有Notify的效果；当Stop返回后，会保证c不再接收到任何信号
	func Stop(c chan<- os.Signal)

	// Example
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Got signal:", s)
```

- os/user
	- user包允许通过名称或ID查询用户帐户
```go
	// User代表一个用户帐户
	type User struct {
		Uid      string // 用户ID
		Gid      string // 初级组ID
		Username string
		Name     string
		HomeDir  string
	}

	// user.Current 返回当前的用户帐户
	func Current() (*User, error)

	// user.Lookup 根据用户名查询用户
	func Lookup(username string) (*User, error)

	// user.LookupId 根据用户ID查询用户
	func LookupId(uid string) (*User, error)
```

### 4. bufio
- bufio包
	- 实现了有缓冲的I/O，提供缓冲流的功能
	- 它包装一个 `io.Reader` 或 `io.Writer` 接口对象，创建另一个也实现了该接口，且同时还提供了缓冲和一些文本I/O的帮助函数的对象
	- Reader
		- 常用函数 
			- `bufio.NewReader` 创建缓冲输入流
		- 常用方法
			- `Reset` 重设缓冲流
			- `Buffered` 回缓冲中现有的可读取的字节数
			- `Peek` 返回输入流的下n个字节，而不会移动读取位置
			- `Read` 读取数据到切片中
			- `ReadByte` 读取并返回一个字节
			- `ReadLine` 读取一行内容到字节切片中
			- `ReadSlice` 根据分隔符读取数据到字节切片
			- `ReadString` 根据分隔符读取数据到字符串
			- `WriteTo` 将数据写入到输出流
	- Writer
		- 常用函数	
			- `bufio.NewWriter` 创建缓冲输出流
		- 常用方法
			- `Reset` 重置输出流
			- `Buffered` 返回缓冲中已使用的字节数
			- `Available` 返回缓冲中还有多少字节未使用
			- `Write` 将字节切片内容写入
			- `WriteString` 将字符串写入
			- `WriteByte` 写入单个字节
			- `WriteRune` 写入一个unicode码值(的utf-8编码)，返回写入的字节数和可能的错误
			- `Flush` 刷新数据到输出流
			- `ReadFrom` ReadFrom实现了io.ReaderFrom接口
	- ReadWriter
		- 常用函数	
			- `bufio.NewReadWriter` 申请创建一个新的、将读写操作分派给r和w 的ReadWriter
	- Scanner
		- 常用函数
			- `bufio.NewScanner` 创建扫描对象
		- 常用方法
			- `Split` 定义流分割函数，默认空格，必须在Scan前
			- `Scan` 扫描数据
			- `Bytes` 读取数据，返回byte数组
			- `Text` 读取数据，返回字符串
			- `Err` 获取错误

### 5. io
- io
	- io包提供了对I/O原语的基本接口
	- 本包的基本任务是包装这些原语已有的实现(如os包里的原语)，使之成为共享的公共接口，这些公共接口抽象出了泛用的函数并附加了一些相关的原语的操作
```go
	// io.Reader 接口用于包装基本的读取方法
	type Reader interface {
		// Read方法读取len(p)字节数据写入p
		Read(p []byte) (n int, err error)
	}

	// io.Writer接口用于包装基本的写入方法
	type Writer interface {
		// Write方法len(p) 字节数据从p写入底层的数据流
		Write(p []byte) (n int, err error)
	}

	// io.Closer接口用于包装基本的关闭方法
	type Closer interface {
		Close() error
	}

	// io.Seeker接口用于包装基本的移位方法
	type Seeker interface {
		// Seek方法设定下一次读写的位置：偏移量为offset，校准点由whence确定：0表示相对于文件起始；1表示相对于当前位置；2表示相对于文件结尾
		// Seek方法返回新的位置以及可能遇到的错误
		// 移动到一个绝对偏移量为负数的位置会导致错误;移动到任何偏移量为正数的位置都是合法的，但其下一次I/O操作的具体行为则要看底层的实现
		Seek(offset int64, whence int) (int64, error)
	}

	// 将src的数据拷贝到dst，直到在src上到达EOF或发生错误，返回拷贝的字节数和遇到的第一个错误
	func Copy(dst Writer, src Reader) (written int64, err error)

	// 从src拷贝n个字节数据到dst，直到在src上到达EOF或发生错误，返回复制的字节数和遇到的第一个错误
	func CopyN(dst Writer, src Reader, n int64) (written int64, err error)

	// ReadAtLeast从r至少读取min字节数据填充进buf
	func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)

	// ReadFull从r精确地读取len(buf)字节数据填充进buf
	func ReadFull(r Reader, buf []byte) (n int, err error)

	// WriteString函数将字符串s的内容写入w中
	func WriteString(w Writer, s string) (n int, err error)
```

- io/ioutil
	- 包ioutil提供了一些I/O实用函数
```go
	// ReadAll从r读取数据直到EOF或遇到error，返回读取的数据和遇到的错误，且成功的调用返回的err为nil而非EOF
	func ReadAll(r io.Reader) ([]byte, error)

	// ReadFile 从filename指定的文件中读取数据并返回文件的内容，成功的调用返回的err为nil而非EOF
	func ReadFile(filename string) ([]byte, error)

	// 函数向filename指定的文件中写入数据，如果文件不存在将按给出的权限创建文件，否则在写入数据之前清空文件
	func WriteFile(filename string, data []byte, perm os.FileMode) error

	// 返回dirname指定的目录的目录信息的有序列表
	func ReadDir(dirname string) ([]os.FileInfo, error)

	// 在dir目录里创建一个新的、使用prfix作为前缀的临时文件夹，并返回文件夹的路径
	func TempDir(dir, prefix string) (name string, err error)

	// 在dir目录下创建一个新的、使用prefix为前缀的临时文件，以读写模式打开该文件并返回os.File指针
	func TempFile(dir, prefix string) (f *os.File, err error)
```

### 6. encoding
- encoding
	- encoding包定义了供其它包使用的可以将数据在字节水平和文本表示之间转换的接口
	- `encoding/gob`, `encoding/json`, `encoding/xml`, 三个包都会检查使用这些接口

- encoding/base64
	- base64实现了RFC 4648规定的base64编码
```go
	// RFC 4648定义的标准base64编码字符集
	var StdEncoding = NewEncoding(encodeStd)

	// RFC 4648定义的另一base64编码字符集，用于URL和文件名
	var URLEncoding = NewEncoding(encodeURL)

	// 双向的编码/解码协议
	type Encoding struct { ... }

	// 使用给出的字符集生成一个*Encoding，字符集必须是64字节的字符串
	func NewEncoding(encoder string) *Encoding

	// 将src的数据解码后存入dst，最多写DecodedLen(len(src))字节数据到dst，并返回写入的字节数
	// 如果src包含非法字符，将返回成功写入的字符数和CorruptInputError
	// 换行符（\r、\n）会被忽略。
	func (enc *Encoding) Decode(dst, src []byte) (n int, err error)

	// 返回base64编码的字符串s代表的数据
	// base64.StdEncoding.DecodeString(str)
	func (enc *Encoding) DecodeString(s string) ([]byte, error)

	// 将src的数据编码后存入dst，最多写EncodedLen(len(src))字节数据到dst，并返回写入的字节数
	// 函数会把输出设置为4的倍数，因此不建议对大数据流的独立数据块执行此方法，使用NewEncoder()代替
	func (enc *Encoding) Encode(dst, src []byte)

	// 返回将src编码后的字符串
	// base64.StdEncoding.EncodeToString(data)
	func (enc *Encoding) EncodeToString(src []byte) string

	// 创建一个新的base64流解码器
	func NewDecoder(enc *Encoding, r io.Reader) io.Reader

	// 创建一个新的base64流编码器
	// 写入的数据会在编码后再写入w，base32编码每3字节执行一次编码操作
	// 写入完毕后，使用者必须调用Close方法以便将未写入的数据从缓存中刷新到w中
	func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser
```

- encoding/csv
	- `encoding/csv` 包提供对 csv 文件读写的操作

- encoding/json
	- json包实现了json对象的编解码
	- [JSON and Go](http://golang.org/doc/articles/json_and_go.html)
```go
	// json.Marshal
	// Marshal函数返回v的json编码
	// 结构体标签值里的"json"键为键名，后跟可选的逗号和选项，具体如下:
	// // 字段被本包忽略
	// Field int `json:"-"`
	// // 字段在json里的键为"myName"
	// Field int `json:"myName"`
	// // 字段在json里的键为"myName"且如果字段为空值将在对象中省略掉
	// Field int `json:"myName,omitempty"`
	// // 字段在json里的键为"Field"（默认值），但如果字段为空值会跳过；注意前导的逗号
	// Field int `json:",omitempty"
	// // "string"选项标记一个字段在编码json时应编码为字符串；它只适用于字符串、浮点数、整数类型的字段
	// Int64String int64 `json:",string"`
	func Marshal(v interface{}) ([]byte, error)

	// json.Unmarshal
	// Unmarshal 函数解析json编码的数据并将结果存入v指向的值
	// Unmarshal 和Marshal 做相反的操作，必要时申请映射、切片或指针，有如下的附加规则
	// JSON 的 null 值解码为go的接口、指针、切片时会将它们设为nil，因为null在json里一般表示“不存在”；解码json的null值到其他go类型时，不会造成任何改变，也不会产生错误
	// 要将json数据解码写入一个接口类型值，函数会将数据解码为如下类型写入接口:
	// Bool                   对应JSON布尔类型
	// float64                对应JSON数字类型
	// string                 对应JSON字符串类型
	// []interface{}          对应JSON数组
	// map[string]interface{} 对应JSON对象
	// nil                    对应JSON的null
	func Unmarshal(data []byte, v interface{}) error

	// Decoder从输入流解码json对象
	type Decoder struct { ... }
	
	// NewDecoder创建一个从r读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据
	func NewDecoder(r io.Reader) *Decoder

	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	func (dec *Decoder) Decode(v interface{}) error

	// Encoder将json对象写入输出流
	type Encoder struct { ... }

	// NewEncoder创建一个将数据写入w的 *Encoder
	func NewEncoder(w io.Writer) *Encoder

	// Encode将v的json编码写入输出流，并会写入一个换行符
	func (enc *Encoder) Encode(v interface{}) error
```

- encoding/gob
	- gob包管理gob流，在编码器(发送器)和解码器(接受器)之间交换的binary值
	- go特有的编码格式，不能跨语言，提供了对数据结构进行二进制序列化的功能
	- 一般用于传递远端程序调用(RPC)的参数和结果，如net/rpc包就有提供

### 7. strings
- strings包
	- 实现了用于操作字符的简单函数
	- `strings.EqualFold` 判断两个字符串是否相同
	- `strings.HasPrefix` 判断 s 是否有前缀字符串 prefix
	- `strings.HasSuffix` 判断 s 是否有后缀字符串 suffix
	- `strings.Contains` 子串 substr 在 s 中，返回 true
	- `strings.Count` 查找子串出现次数即字符串模式匹配
	- `strings.Repeat` 返回count个s串联的字符串
	- `strings.Index` 在 s 中查找 sep 的第一次出现，返回第一次出现的索引，不存在则返回-1
	- `strings.LastIndex`
	- `strings.Title` 返回s中每个单词的首字母都改为标题格式的字符串拷贝
	- `strings.ToLower` 返回将所有字母都转为对应的小写版本的拷贝
	- `strings.ToUpper` 返回将所有字母都转为对应的大写版本的拷贝
	- `strings.Replace` 返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串
	- `strings.Map` 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝
	- `strings.Trim` 返回将s前后端所有cutset包含的utf-8码值都去掉的字符串
	- `strings.Fields` 返回将字符串按照空白（unicode.IsSpace确定，可以是一到多个连续的空白字符）分割的多个字符串
	- `strings.SplitN` 用去掉s中出现的sep的方式进行分割，会分割到结尾，并返回生成的所有片段组成的切片
	- `strings.Join` 将一系列字符串连接为一个字符串，之间用sep来分隔
	- `strings.NewReader` NewReader创建一个从s读取数据的Reader

### 8. context
- Context 包定义了上下文类型，该上下文类型跨越 API 边界和进程之间传递截止期限，取消信号和其他请求范围值
	- 对服务器的传入请求应创建一个 Context，对服务器的传出调用应接受 Context
		- 它们之间的函数调用链必须传播 Context，可以用使用 `WithCancel`，`WithDeadline`，`WithTimeout` 或 `WithValue` 创建的派生上下文替换
		- 当 Context 被取消时，从它派生的所有 Context 也被取消
	- 即使函数允许，也不要传递 nil Context
		- 如果不确定要使用哪个Context，请传递 `context.TODO`
	- 使用上下文值仅适用于传输进程和 API 的请求范围数据，而不用于将可选参数传递给函数
	- 相同的上下文可以传递给在不同 goroutine 中运行的函数; 
		- 上下文对于多个 goroutine 同时使用是安全的。
```go
	// Background 返回non-nil(非零)，空的 Context; 其从未被取消，没有值，也没有最后期限
	// 通常由主函数，初始化和测试使用，并作为传入请求的top-level Context (顶级上下文)
	func Background() Context

	// TODO 返回非零空的上下文，一般在不清楚使用哪个Context 或 它尚不可用时(因为周围的函数尚未扩展为接受Context参数)
	func TODO() Context

	// WithCancel 返回父上下文的副本，同时返回一个 CancelFunc
	// 取消这个上下文可以释放与它相关的资源，因此只要在这个Context 中运行的操作完成，代码就应该立即调用 cancel，下同
	func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

	// WithDeadline 返回父上下文的副本，并将截止日期调整为不晚于deadline，同时返回一个 CancelFunc
	func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

	// WithTimeout 返回 WithDeadline(parent, time.Now().Add(timeout))
	func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

	// WithValue 返回父键的副本，其中与键关联的值是val
	// 提供的密钥必须具有可比性，不应该是字符串类型或任何其他内置类型，以避免使用上下文的包之间发生冲突
	func WithValue(parent Context, key, val interface{}) Context
```

### 9. log
- log包
	- 实现了简单的日志服务
	- 定义了Logger类型，该类型提供了一些格式化输出的方法
		- `log.New` 创建一个Logger
		- `logger.Println` 调用Output将生成的格式化字符串输出到标准logger
	- Fatal系列函数会在写入日志信息后调用 `os.Exit(1)`
		- `log.Fatalf` 等价于 `{Printf(v...); os.Exit(1)}`
	- Panic系列函数会在写入日志信息后 `panic`
		- `log.Panicf` 等价于 `{Printf(v...); panic(...)}`

- log/syslog
	- syslog 包提供一个简单的系统日志服务的接口

### 10. archive
- archive/tar包
	- tar包实现了tar格式压缩文件的存取
	- 本包目标是覆盖大多数tar的变种，包括GNU和BSD生成的tar文件

- archive/zip包
	- zip包提供了zip档案文件的读写服务
	- 本包不支持跨硬盘的压缩

### 11. compress
- compress/bzip2
	- bzip2包实现bzip2的解压

- compress/flate
	- flate包实现了deflate压缩数据格式
	- gzip包和zlib包实现了对基于deflate的文件格式的访问

- compress/gzip
	- gzip包实现了gzip格式压缩文件的读写

- compress/gzip
	- gzip包实现了gzip格式压缩文件的读写

- compress/lzw 
	- lzw包实现了Lempel-Ziv-Welch数据压缩格式

- compress/zlib
	- zlib包实现了对zlib格式压缩数据的读写
	- 本包的实现提供了在读取时解压和写入时压缩的滤镜

### 12. container
- container/heap
	- heap包提供了对任意类型(实现了heap.Interface接口)的堆操作
	- (最小)堆是具有"每个节点都是以其为根的子树中最小值"属性的树
	- 树的最小元素为其根元素，索引0的位置

- container/list
	- list包实现了双向链表

- container/ring
	- ring包实现了环形链表的操作

### 13. crypto
- crypto
	- crypto包搜集了常用的密码(算法)常量

- crypto/aes
	- aes包实现了AES加密算法

- crypto/cipher
	- cipher包实现了多个标准的用于包装底层块加密算法的加密算法实现

- crypto/des
	- des包实现了DES标准和TDEA算法

- crypto/dsa
	- dsa包实现FIPS 186-3定义的数字签名算法(Digital Signature Algorithm)，即DSA算法

- crypto/ecdsa
	- ecdsa包实现了椭圆曲线数字签名算法

- crypto/elliptic
	- elliptic包实现了几条覆盖素数有限域的标准椭圆曲线

- crypto/hmac
	- hmac包实现了 U.S. Federal Information Processing Standards Publication 198规定的HMAC(加密哈希信息认证码)。
	- HMAC是使用key标记信息的加密hash
	- 接收者使用相同的key逆运算来认证hash

- crypto/md5
	- md5包实现了MD5哈希算法
```go
	md5.Sum([]byte(""))     // 计算byte切片中字符的MD5
	md5.New()               // 解码
```

- crypto/rand
	- rand包实现了用于加解密的更安全的随机数生成器

- crypto/rc4
	- rc4包实现了RC4加密算法

- crypto/rsa
	- rsa包实现了PKCS#1规定的RSA加密算法

- crypto/sha1
	- sha1包实现了SHA1哈希算法
```go
	sha1.Sum([]byte(""))    // 计算byte切片中字符消息摘要(Hash)
	sha256.Sum([]byte(""))
	sha512.Sum([]byte(""))
	sha256.Sum256([]byte(""))
```

- crypto/sha256
	- sha256包实现了SHA224和SHA256哈希算法

- crypto/sha512
	- sha512包实现了SHA384和SHA512哈希算法

- crypto/tls
	- tls包实现了TLS 1.2

- crypto/x509
	- x509包解析X.509编码的证书和密钥

- crypto/x509/pkix
	- pkix包提供了共享的、低层次的结构体，用于ASN.1解析和X.509证书、CRL、OCSP的序列化

### 14. debug
- debug/dwarf
	- dwarf包提供对从可执行文件加载的dwarf调试信息的访问，如 http://dwarfstd.org/doc/dwarf-2.0.0.pdf 的DWARF 2.0标准中所定义的

- debug/elf
	- elf包提供对 ELF 文件的访问

- debug/gosym
	- gosym包提供对gc编译器生成的Go二进制文件中嵌入的Go符号和行号表的访问

- debug/macho
	- macho包提供对 Mach-O 件的访问

- debug/pe
	- pe包提供对对 PE(Microsoft Windows Portable Executable) 文件的访问

- debug/plan9obj
	- plan9obj包提供对 Plan 9 a.out 文件的访问

### 15. runtime
- runtime
	- runtime包提供和go运行时环境的互操作，如控制go程的函数
	- 它也包括用于reflect包的低层次类型信息
	- 也包括用于reflect包的低层次类型信息

- runtime/cgo
	- cgo 包含有 cgo 工具生成的代码的运行时支持

- runtime/debug
	- debug包 包含程序在运行时自我调试的工具

- runtime/pprof
	- pprof包以pprof可视化工具期望的格式书写运行时剖面数据

- runtime/race
	- race包实现了数据竞争检测逻辑
	- 没有提供公共接口

- runtime/trace
	- 执行追踪器；tracer捕获各种执行事件，如goroutine创建/阻塞/解除阻塞、syscall进入/退出/阻塞、GC相关事件、堆大小的变化、处理器启动/停止等，并将它们写入io
	- 对于大多数事件，都会捕获精确到纳秒级的时间戳和堆栈跟踪
	- 使用 `go tool trace` 命令来分析跟踪

### 16. errors
- errors
	- errors包实现了创建错误值的函数
```go
	// 使用字符串创建一个错误，同fmt.Errorf，类似于 New(fmt.Sprintf(...))
	func New(text string) error
```

### 17. expvar
- expvar
	- expvar包提供了公共变量的标准接口，如服务的操作计数器
	- 本包通过HTTP在/debug/vars位置以JSON格式导出了这些变量

### 18. flag
- flag
	- flag包实现了命令行参数的解析
	- [flag 示例](../cloudstation/simple_tool/simple_call.go)
```go
	flag.Parse()           // 解析命令行参数
	flag.IntVar()          // 设置int类型参数
	flag.BoolVar()		   // 设置bool类型参数
	flag.StringVar()       // 设置string类型参数    func flag.StringVar(p *string, name string, value string, usage string)
	flag.PrintDefaults()   // 获取自动生成的参数信息
```

### 19. hash
- hash
	- hash包提供hash函数的接口

- hash/adler32
	- adler32包实现了Adler-32校验和算法

- hash/crc32
	- crc32包实现了32位循环冗余校验(CRC-32)的校验和算法

- hash/crc64
	- crc64包实现了64位循环冗余校验(CRC-64)的校验和算法

- hash/fnv
	- fnv包实现了FNV-1和FNV-1a(非加密hash函数)

### 20. html
- html
	- html包提供了用于转义和解转义HTML文本的函数
```go
	// EscapeString函数将特定的一些字符转为逸码后的字符实体，如"<"变成"&lt;"
	// 它只会修改五个字符：<、>、&、'、"
	// UnescapeString(EscapeString(s)) == s总是成立，但是两个函数顺序反过来则不一定成立
	func EscapeString(s string) string
	// UnescapeString函数将逸码的字符实体如"&lt;"修改为原字符"<"
	func UnescapeString(s string) string
```

- html/template
	- template包实现了数据驱动的模板，用于生成可对抗代码注入的安全HTML输出
	- 本包提供了和text/template包相同的接口，无论何时当输出是HTML的时候都应使用本包
	- [html/template 包](https://studygolang.com/static/pkgdoc/pkg/html_template.htm)

### 21. image
- image
	- image实现了基本的2D图片库
	- 基本接口叫作Image，图片的色彩定义在image/color包
	- Image接口可以通过调用如`NewRGBA`和`NewPaletted`函数等获得；也可以通过调用`Decode`函数解码包含GIF、JPEG或PNG格式图像数据的输入流获得
	- 解码任何具体图像类型之前都必须注册对应类型的解码函数；注册过程一般是作为包初始化的副作用，放在包的init函数里；要解码PNG图像，只需在程序的main包里嵌入如下代码 `import _ "image/png"` _表示导入包但不使用包中的变量/函数/类型，只是为了包初始化函数的副作用

- image/color
	- color包实现了基本色彩库

- image/color/palette
	- palette包提供了标准的调色板

- image/draw
	- draw包提供了图像合成函数

- image/gif
	- gif包实现了gif文件的编码器和解码器

- image/jpeg
	- jpeg包实现了jpeg格式图像的编解码

- image/png
	- png包实现了PNG图像的编解码

### 22. reflect
- index/suffixarray
	- suffixarrayb包通过使用内存中的后缀树实现了对数级时间消耗的子字符串搜索

### 23. mime
- mime
	- mime包实现了MIME的部分规定

- mime/multipart
	- multipart包实现了MIME的multipart解析，参见 [RFC 2046](http://tools.ietf.org/html/rfc2046)
	- 该实现适用于HTTP([RFC 2046](https://www.rfc-editor.org/rfc/rfc2388))和常见浏览器生成的multipart主体

- mime/quotedprintable
	- quotedprintable包实现了quoted-printable encoding，参见 [RFC 2045](http://tools.ietf.org/html/rfc2045)

### 24. path
- path
	- path包实现了对斜杠分隔的路径的实用操作函数

- path/filepath
	- filepath包实现了兼容各操作系统的文件路径的实用操作函数

### 25. plugin
- plugin
	- plugin包实现Go plugins的加载和符号解析
	- plugin是一个带有导出函数和变量的Go main包，目前plugin只在Linux上工作 `go build -buildmode=plugin`
	- 当一个plugin第一次打开时，所有不属于程序的包的init函数都会被调用；主函数不会运行，插件只初始化一次，并且无法关闭。

### 26. reflect
- reflect
	- reflect包实现了运行时反射，允许程序操作任意类型的对象
	- 典型用法是用静态类型 `interface{}`保存一个值，通过调用 `TypeOf`获取其动态类型信息，该函数返回一个 `Type`类型值
	- 调用 `ValueOf`函数返回一个 `Value`类型值，该值代表运行时的数据
	- `Zero`接受一个 `Type`类型参数并返回一个代表该类型零值的 `Value`类型值。
```go
	// 获取数据类型，同Printf("%T")
	reflect.TypeOf()
	reflect.ValueOf()
```

### 27. fmt
- fmt
    - fmt包实现了类似C语言printf和scanf的格式化I/O
    - 格式化动作('verb')源自C语言但更简单

### 28. regexp
- regexp
	- regexp包实现了正则表达式搜索
	- 正则表达式采用RE2语法（除了\c、\C），和Perl、Python等语言的正则基本一致
	- 参见 [Syntax](http://code.google.com/p/re2/wiki/Syntax)z

- regexp/syntax
	- syntax包将正则表达式解析成解析树，并将解析树编译成程序
	- 一般使用regexp包的功能

### 29. sort
- sort
	- sort包提供了排序切片和用户自定义数据集的函数
```go
	// Search函数采用二分法搜索找到[0, n)区间内最小的满足f(i)==true的值i
	// Search函数希望f在输入位于区间[0, n)的前面某部分(可以为空)时返回假，而在输入位于剩余至结尾的部分(可以为空)时返回真；Search函数会返回满足f(i)==true的最小值i；如果没有该值，函数会返回n
	func Search(n int, f func(int) bool) int
	
	// Sort 排序data，它调用1次data.Len确定长度，调用O(n*log(n))次data.Less和data.Swap
	func Sort(data Interface)

	// Reverse包装一个Interface接口并返回一个新的Interface接口，对该接口排序可生成递减序列
	func Reverse(data Interface) Interface
	// Ints 函数将a排序为递增顺序
	func Ints(a []int)
	// Float64s函数将a排序为递增顺序
	func Float64s(a []float64)
	// Strings 函数将a排序为递增顺序
	func Strings(a []string)
```

### 30. strconv
- strconv
	- strconv包实现了基本数据类型和其字符串表示的相互转换
	- 提供了字符串与简单数据类型之间的类型转换功能，可以将简单类型转换为字符串，也可以将字符串转换为其它简单类型
	- string 转 bool `strconv.ParseBool("true")`
	- bool 转 string `strconv.FormatBool(true)`
	- string 转 int `strconv.ParseInt("11111111", 2, 16)` 、 `strconv.Atoi("100x")`
	- int 转 string `strconv.FormatInt(255, 10)` 、 `strconv.Itoa(100)`
	- string 转 uint `strconv.ParseUint("4E2D", 16, 16)`
	- uint 转 string `strconv.FormatUint(255, 16)`
	- string 转 float `strconv.ParseFloat("3.1", 64)`
	- float 转 string `strconv.FormatFloat(3.1415, 'E', -1, 64)`

### 31. sync
- sync
	- sync包提供了基本的同步基元，如互斥锁
	- 除了Once和WaitGroup类型，大部分都是适用于低水平程序线程，高水平的同步使用channel通信更好一些
	- 本包的类型的值不应被拷贝
	- 互斥锁 `sync.Mutex`
	- 读写锁 `sync.RWMutex`
	- 线程安全的Map `sync.Map`
	- 等待一组线程 `sync.WaitGroup(计数信号量)`
	- 执行一次动作 `sync.Once`

- sync/atomic
	- atomic包提供了底层的原子级内存操作，对于同步算法的实现很有用
	- 应通过通信来共享内存，而不通过共享内存实现通信

### 32. syscall
- syscall
	- syscall包 包含一个到低级操作系统原语的接口
	- 具体细节因底层系统而异，默认情况下，godoc将显示当前系统的syscall文档
	- 如果希望godoc显示另一个系统的syscall文档，可以将$GOOS和$GOARCH设置为所需的系统

### 33. text
- text/scanner
	- scanner包提供对utf-8文本的token扫描服务
		- 它会从一个io.Reader获取utf-8文本，通过对Scan方法的重复调用获取一个个token
		- 为了兼容已有的工具，NUL字符不被接受
		- 如果第一个字符是表示utf-8编码格式的BOM标记，会自动忽略该标记

- text/tabwrite
	- tabwriter包实现了写入过滤器 `tabwriter.Writer`，可以将输入的缩进修正为正确的对齐文本

- text/template
	- template包实现了数据驱动的用于生成文本输出的模板

- text/template/parse
	- parse包 由 text/template包 和 html/template包定义的模版，构建解析树

### 34. unicode
- unicode
	- unicode包提供数据和函数来测试Unicode码位的一些属性

- unicode/utf16
	- utf16包实现了UTF-16序列的编解码

- unicode/utf8
	- utf8包实现了对utf-8文本的常用函数和常数的支持，包括rune和utf-8编码byte序列之间互相翻译的函数

### 35. net
- net
	- net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket
	- 虽然本包提供了对网络原语的访问，大部分使用者只需要Dial、Listen和Accept函数提供的基本接口；以及相关的Conn和Listener接口
	- crypto/tls包提供了相同的接口和类似的Dial和Listen函数

- net/http
	- http包提供了HTTP客户端和服务端的实现
	- Get、Head、Post和PostForm函数发出HTTP/ HTTPS请求
		- `http.Get("http://example.com/")`
		- `http.Post("http://example.com/upload", "image/jpeg", &buf)`
		- `http.PostForm("http://example.com/form", url.Values{"key": {"Value"}, "id": {"123"}})`
		- `http.StatusOK`

- net/http/cgi
	- cgi包实现了CGI(Common Gateway Interface，公共网关协议)，参见RFC 3875
	- 注意使用CGI意味着对每一个请求开始一个新的进程，这显然要比使用长期运行的服务程序要低效
	- 本包主要是为了兼容现有的系统

- net/http/cookiejar
	- cookiejar包实现了保管在内存中的符合RFC 6265标准的http.CookieJar接口

- net/http/fcgi
	- fcgi包实现了FastCGI协议
	- 目前只支持响应器的角色

- net/http/httptest
	- httptest包提供了HTTP测试的常用函数

- net/http/httptrace
	- httptrace包提供了跟踪HTTP客户端请求中的事件的机制

- net/http/httputil
	- httputil包提供了HTTP公用函数，是对net/http包的更常见函数的补充

- net/http/pprof
	- pprof包通过它的HTTP服务端提供pprof可视化工具期望格式的运行时剖面文件数据服务
	- 本包一般只需导入获取其注册HTTP处理器的副作用。处理器的路径以/debug/pprof/开始
		- http://127.0.0.1:8080/debug/pprof/goroutine?debug=1
		- `func bytes.TrimSpace(s []byte) []byte` 去除首尾空格
		- `func bytes.Replace(s []byte, old []byte, new []byte, n int) []byte` 替换字符
		- `func bytes.Join(s [][]byte, sep []byte) []byte`

- net/smtp
	- smtp包实现了简单邮件传输协议(SMTP)

- net/mail
	- mail包实现了邮件的解析

- net/rpc/jsonrpc
	- rpc包提供了通过网络或其他I/O连接对一个对象的导出方法的访问

- net/rpc/jsonrpc
	- jsonrpc包实现了JSON-RPC的ClientCodec和ServerCodec接口，可用于rpc包。

- net/smtp
	- smtp包实现了简单邮件传输协议(SMTP)

- net/textproto
	- textproto实现了对基于文本的请求/回复协议的一般性支持，包括HTTP、NNTP和SMTP

- net/url
	- url包解析URL并实现了查询的逸码

## 十、常用公共库 

### 1. Log
- `github.com/sirupsen/logrus` logrus库 完全兼容标准log库，被广泛使用
- `github.com/uber-go/zap` 快速、结构化、分级的日志库
- `github.com/rs/zerolog` zerolog包提供了一个快速简单的日志记录器，专门用于JSON输出
- `github.com/rifflock/lfshook` 为logrus库设计的钩子

### 2. Db
- SQLBuilder
	- `github.com/parkingwang/go-sqlbuilder` 创建SQL语句
	- `github.com/didi/gendr` 辅助操作数据库

- ORM
	- GORM
		- `gorm.io/gorm`
		- `gorm.io/driver/mysql`

- MongoDB
	- `go.mongodb.org/mongo-driver`

- Redis
	- `github.com/go-redis/redis`

- Cache
	- `github.com/bluele/gcache` Supports expirable Cache, LFU, LRU and ARC.

### 3. 爬虫
- `github.com/andeya/pholcus`
	- 一款纯 Go 语言编写的支持分布式的高并发爬虫软件

### 4. Web
- 框架
	- `github.com/gin-gonic/gin`

- Restful
	- `github.com/emicklei/go-restful/v3`
	- `github.com/go-openapi/spec` Marshal and unmarshal Swagger API

- WebSocket
	- `github.com/gorilla/websocket`

- Httprouter
	- `github.com/julienschmidt/httprouter`

### 5. kubernetes
- `k8s.io/client-go`
	- `k8s.io/client-go/rest`
	- `k8s.io/client-go/kubernetes`
	- `k8s.io/client-go/kubernetes/typed/apps/v1`
	- `k8s.io/client-go/kubernetes/typed/core/v1`
	- `k8s.io/client-go/kubernetes/typed/batch/v1`
	- `k8s.io/client-go/kubernetes/typed/networking/v1`
	- `k8s.io/client-go/tools/clientcmd`
	- `k8s.io/client-go/tools/clientcmd/api`
	- [kubernetes client-go功能介绍](https://www.cnblogs.com/haiyux/p/17162339.html)

- Api
	- `k8s.io/apimachinery/pkg/apis/meta/v1`
	- `k8s.io/api/autoscaling/v1`
	- `k8s.io/api/networking/v1`

- Etcd
	- `go.etcd.io/etcd/client/v3`
	- `go.etcd.io/etcd/client/v3/concurrency` 选举

### 6. Marshal && UNMarshal
- `gopkg.in/yaml.v2`



### 3. github.com/go-playground/validator
- 校验参数

### 4. github.com/spf13/cobra
- cmd命令行客户端       

### 5. github.com/AlecAivazis/survey
- 隐藏密码

### 6. github.com/schollz/progressbar
- 进度条

### 7. github.com/BurntSushi/toml
- toml解析库，解析配置文件

### 8. github.com/caarlos0/env
- 环境变量解析库，解析环境变量

### 9. github.com/rs/xid
- ID生成器库

### 10. github.com/stretchr/testify/assert
- 测试

### 13.github.com/distribution/distribution
- Registry

### 14.k8s.io/apimachinery/pkg/util/uuid
- Uuid

### 15.github.com/pkg/errors
- Errors

### 18. github.com/jedib0t/go-pretty
- 美化表格、列表、进度条、文本等控制台输出
