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
- 项目目录下 go mod init github.com/ahwhy/myGolang
- git init 
- 添加 .gitignore 文件去掉一些和代码无关的文件/文件夹
- git add . && git commit -m "Record me learning golang"
- github上新建一个仓库
- 推送到远程
- 上传tag
	- 通过tag可以返回到项目的特定状态下，可以将tag看作是在大量commit中设定的书签
	- https://www.jianshu.com/p/79ecf4fe5079
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
- 标准库 https://studygolang.com/pkgdoc

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