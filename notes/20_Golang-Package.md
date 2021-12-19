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
	$ go env -w GOPROXY=https://goproxy.io,direct`
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

## 十、常用包与函数  
- 标准库 https://studygolang.com/pkgdoc

### 1. time
```go
	time.Now()          // 获取当前时间
	time.Now().Unix()
	time.Now().Year()   // Month() Day()  Hour()  Minute()  Second()
	time.Now().Format("2006-01-02 15:04:05")
	time.Parse()        // 返回转换后的时间格式和一个判断信息（err)
	time.Sleep(1 * time.Second)
	time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
```

### 2. math
```go
	// 用于测试一个数是否是非数NaN
	// NaN非数，一般用于表示无效的除法操作结果0/0或Sqrt(-1)
	func IsNaN(f float64) (is bool)

	// 返回非数对应的值
	func NaN() float64
```

- math/rand
```go
	rand.Seed(time.Now().Unix())  // 使用当前时间设置随机数种子
	rand.Intn(100)    // 生产[0, 100)的随机数
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

### 17. gopkg.in/yaml.v2
```go
	yaml.Unmarshal()
```

### 18. net/http
```go
	http.StatusOK
```

### 19. github.com/gin-gonic/gin
- http客户端

### 20. github.com/go-playground/validator
- 校验参数

### 21. github.com/spf13/cobra
- cmd命令行客户端       

### 22. github.com/AlecAivazis/survey
- 隐藏密码

### 23. github.com/schollz/progressbar
- 进度条

### 24. github.com/BurntSushi/toml
- toml解析库，解析配置文件

### 25. github.com/caarlos0/env
- 环境变量解析库，解析环境变量

### 26. github.com/rs/xid
- ID生成器库

### 27. github.com/stretchr/testify/assert
- 测试

### 28. net/http/pprof
- http://127.0.0.1:8080/debug/pprof/goroutine?debug=1
