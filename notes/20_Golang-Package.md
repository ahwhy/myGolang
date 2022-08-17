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
	// 推送到远程
	// or push an existing repository from the command line
	git remote add origin https://github.com/ahwhy/myGolang.git
	git branch -M main
	git push -u origin main

	// 上传tag
	git log --oneline
	// 创建lightweight类型的tag
	git tag v0.0.2-lw
	// 创建annotated类型的tag
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

### 3. context
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
	// CancelFunc 通知操作放弃其工作，且不会等待工作停止，在第一次调用之后，对CancelFunc的后续调用不起作用
	type CancelFunc func()

	type Context interface {
		Deadline() (deadline time.Time, ok bool)
		Done() <-chan struct{}
		Err() error
		Value(key interface{}) interface{}
	}

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

- 示例
```go
	// context.WithDeadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50 * time.Millisecond)) // context.WithTimeout(context.Background(), 50 * time.Millisecond)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
	// Output: 
	// context deadline exceeded

	// context.WithValue
	type favContextKey string
	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")
	f(ctx, k)
	f(ctx, favContextKey("color"))
	// Output:
	// found value: Go
	// key not found: color
```

### 4. encoding
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

### 3. reflect
```go
	//  获取数据类型，同Printf("%T")
	reflect.TypeOf()
	reflect.ValueOf()
```

### 4. os
```go
	// os文件处理
	*os.PathError           // PathError records an error and the operation and file path that caused it.
	*os.LinkError           // LinkError records an error during a link or symlink or rename system call and the paths that caused it.
	*os.SyscallError        // SyscallError records an error from a specific system call.

	//系统退出
	os.Esxit(1) 
	os.Args                 // 接收命令行参数生成切片，从程序本身的路径开始 var os.Args []string
	os.Stat("test.txt")     // func os.Stat(name string) (fs.FileInfo, error)  Stat returns a FileInfo describing the named file. If there is an error, it will be of type *PathError.
```

### 5. strings
```go
	strings.FieldsFunc()  // 将字符串进行分段，返回切片 func strings.FieldsFunc(s string, f func(rune) bool) []string
	strings.Contains()    // func Contains(s, substr string) bool  判断字符串s中是否存在对应字符substr
	strings.ToLower()     // func ToLower(s string) string         将字符串统一转成小写
	strings.NewReader("")     // 从字符串创建一个reader对象
	strings.Reader.Reader()   // func (*strings.Reader).Read(b []byte) (n int, err error)
```

### 6. errors
```go
	errors.New() // 创建错误 或使用 fmt.Errorf()
```

### 7. sort
```go
	type StringSlice []string
	func sort.Strings(x []string)
	func sort.Sort(data sort.Interface)
```

### 8. flag
```go
	flag.Parse()           // 解析命令行参数
	flag.IntVar()          // 设置int类型参数
	flag.BoolVar()		   // 设置bool类型参数
	flag.StringVar()       // 设置string类型参数    func flag.StringVar(p *string, name string, value string, usage string)
	flag.PrintDefaults()   // 获取自动生成的参数信息
```

### 9. crypto
- crypto/md5
```go
	md5.Sum([]byte(""))     // 计算byte切片中字符的MD5
	md5.New()               // 解码
```
- crypto/sha1
```go
	sha1.Sum([]byte(""))    // 计算byte切片中字符消息摘要(Hash)
	sha256.Sum([]byte(""))
	sha512.Sum([]byte(""))
	sha256.Sum256([]byte(""))
```

### 10. encoding/base64
```go
	base64.stdEncoding.EncodeToString([]byte(""))     // 计算byte切片中字符的base64加密
	base64.StdEncoding.DecodeString()                 // 计算base64解码
	base64.RawStdEncoding.EncodeToString([]byte(""))  // 计算byte切片中字符的base64加密且不使用=填充
	base64.URLEncoding.EncodeToString([]byte(""))     // 计算byte切片中字符的url加密
	encoding/gob
	encoding/csv
	encoding/json
```

### 11. log
```go
	log.Printf("aa")  // 2021/07/04 15:32:10 aa
```

### 12. sync
```go
	sync.Mutex 互斥锁
	sync.RWMutex 读写锁
	sync.Map 线程安全的Map
	sync.WaitGroup(计数信号量)
	sync.Once
```

### 13. runtime
```go
	func runtime.Gosched()
```

### 14. bufio
```go
	// 提供缓冲流的功能
	bufio.NewScanner(os.Stdin)
```

### 15. io 
```go
	func io.Copy(dst io.Writer, src io.Reader) (written int64, err error)
```

### 16. io/ioutil

### 17. net/http
```go
	http.StatusOK
```

### 18. net/http/pprof
- http://127.0.0.1:8080/debug/pprof/goroutine?debug=1

### 19. bytes
- `func bytes.TrimSpace(s []byte) []byte` 去除首尾空格
- `func bytes.Replace(s []byte, old []byte, new []byte, n int) []byte` 替换字符
- `func bytes.Join(s [][]byte, sep []byte) []byte`


## 十、常用公共库 

### 1. gopkg.in/yaml.v2
```go
	yaml.Unmarshal()
```

### 2. github.com/gin-gonic/gin
- http客户端

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

### 11. github.com/gorilla/websocket
- WebSocket

### 12. github.com/julienschmidt/httprouter
- Httprouter

### 13.github.com/distribution/distribution
- Registry

### 14.k8s.io/apimachinery/pkg/util/uuid
- Uuid

### 15.github.com/pkg/errors
- Errors

### 16.k8s.io/client-go
- k8s.io/client-go/rest
- k8s.io/client-go/tools/clientcmd