# Golang-File  Golang的文件操作

## 一、基本概念

- 绝对路径
	- 文件路径字符串从根路径(盘符)开始
	
- 相对路径
	- 二进制文件运行的目录 `cd $directory`
	- 相对程序执行的路径，即当前shell 处于的路径

- 文件类型
	cat/记事本 -> 文本内容(无乱码 .go .txt) -> 文本文件 -> `string`
	cat/记事本 -> 有乱码(word, zip, excel) -> 二进制文件 -> `[]byte`

- I/O操作
	- 也叫输入/输出操作，其中I是指Input，O是指Output，主要用来读取或写入数据，很多语言中也叫做流操作
	- Go语言中 输入和输出操作是使用原语实现的
		- 这些原语将数据模拟成可以读或者可以写的字节流
		- DataSource -> `io.Reader` -> `Transfer buffer []byte` -> `io.Writer` -> Target

## 二、文件

### 1. 基本操作(不带缓冲IO): 读、写
- 处理流程
	- 打开文件 -> 错误处理 -> 延迟关闭 -> 读/写/其他 -> 关闭文件

- Go语言中的os包，提供了对文件、系统和进程的操作函数
	- 创建 `os.Create()`
	- 读取 `os.Open()`
	- 获取属性 `os.Open().Stat`，`os.Stat` 
	- 修改属性 -> 权限，所属人 
		`os.Chmod()`，
		`os.Chown()`
	- 重命名 `os.Rename("a.txt", "b.txt")`
	- 删除文件 `os.Remove("b.txt")`
	- 常用常量
		- flag
```go
			// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
			O_RDONLY int = syscall.O_RDONLY // open the file read-only.
			O_WRONLY int = syscall.O_WRONLY // open the file write-only.
			O_RDWR   int = syscall.O_RDWR   // open the file read-write.
			// The remaining values may be or'ed in to control behavior.
			O_APPEND int = syscall.O_APPEND // append data to the file when writing.
			O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
			O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
			O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O. 使用同步 I/O
			O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened. 截断(清空)文件
```
		- ModePerm
```go
			// The single letters are the abbreviations
			// used by the String method's formatting.
			ModeDir        = fs.ModeDir        // d: is a directory
			ModeAppend     = fs.ModeAppend     // a: append-only
			ModeExclusive  = fs.ModeExclusive  // l: exclusive use
			ModeTemporary  = fs.ModeTemporary  // T: temporary file; Plan 9 only
			ModeSymlink    = fs.ModeSymlink    // L: symbolic link
			ModeDevice     = fs.ModeDevice     // D: device file
			ModeNamedPipe  = fs.ModeNamedPipe  // p: named pipe (FIFO)
			ModeSocket     = fs.ModeSocket     // S: Unix domain socket
			ModeSetuid     = fs.ModeSetuid     // u: setuid
			ModeSetgid     = fs.ModeSetgid     // g: setgid
			ModeCharDevice = fs.ModeCharDevice // c: Unix character device, when ModeDevice is set
			ModeSticky     = fs.ModeSticky     // t: sticky
			ModeIrregular  = fs.ModeIrregular  // ?: non-regular file; nothing else is known about this file
			// Mask for the type bits. For regular files, none will be set.
			ModeType = fs.ModeType
			ModePerm = fs.ModePerm // Unix permission bits, 0o777
```
		- I/O
```go
			Stdin    // Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
			Stdout   // Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
			Stderr   // Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
```
	- 常用函数
		- `Args`: 获取命令行参数
		- `Hostname`: 获取主机名
		- `Getpid`: 获取当前进程名
		- `Getenv`: 获取一条环境变量
		- `Environ`: 获取所有环境变量
		- `Getwd`: 获取当前目录
		- `Chmod`: 修改文件权限
		- `Chown`: 修改文件所属用户，用户组
		- `Chtimes`: 修改文件访问时间和修改时间
		- `IsExist`: 与 `os.Stat` 一起用于判断文件存在
		- `IsNotExist`: 与 os.Stat 一起用于判断文件不存在
		- `Link`: 创建软链接
		- `Mkdir`: 创建文件夹
		- `MkdirAll`: 创建文件夹(父目录不存在逐层创建)
		- `Remove`: 移除文件或空文件夹
		- `RemoveAll`: 移除所有文件
		- `Rename`: 重命名
	- 常用结构体 
```go
		// File: 对文件操作 
		type File struct {
				// Has unexported fields.
		}
         	  
		// 常用函数
		// Create: 创建文件并返回文件对象指针(文件不存在则创建，文件存在则清空)
		func os.Create(name string) (*os.File, error)
		// Open: 打开文件并返回文件对象指针
		func os.Open(name string) (*os.File, error)
		// OpenFile: 按指定权限打开文件，并返回文件指针对象
		func OpenFile(name string, flag int, perm FileMode) (*File, error)
		
		// 常用方法
		// Read: 读取文件到字节切片
		func (*os.File).Read(b []byte) (n int, err error)
		// Write: 写入字节切片到文件
		func (*os.File).Write(b []byte) (n int, err error)
		// WriteString: 写入字符串到文件
		func (*os.File).WriteString(s string) (n int, err error)
		// Readdir: 获取目录下所有文件信息
		func (*os.File).ReadDir(n int) ([]fs.DirEntry, error)
		// Readdirnames: 获取目录下所有文件名
		func (*os.File).Readdirnames(n int) (names []string, err error)
		// Seek: 设置文件指针位置
		func (f *File) Seek(offset int64, whence int) (ret int64, err error)  whence: 0文件开始, 1当前位置, 2文件末尾; 在大部分编程语言中，不支持在文件的开始或中间插入数据(会直接从光标位置进行覆写)，只支持在文件末尾进行数据追加
		// Stat: 获取文件状态信息
		func (file *File) Stat() (FileInfo, error)
		// Sync: 同步文件到硬盘      
		func (f *File) Sync() error
		// Close: 关闭文件
		func (*os.File).Close() error
```

- 参考示例
	- 创建和写入文件
```go
		path := "test.txt"
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		
		name := "aa"
		fmt.Fprintf(file, "I am %s\n", name)
		file.WriteString("ccc联动")
		file.Write([]byte("123456789\n"))
		fmt.Println(file.Write([]byte("bb")))    // 两个返回值 一个是 []byte 的长度，一个是 error
```
	- 读文件
```go
		path := "test.txt"
		file, err := os.Open(path)
		fmt.Println(file, err)
		if err != nil {
			return
		}
		defer file.Close()
		
		content := make([]byte, 3)
		
		for {
			n, err := file.Read(content)
			if err != nil {
				if err != io.EOF {               // EOF(End Of File) -> 标识文件读取结束
					fmt.Print(err)
				} else {
					fmt.Print(err)
				}
				break
			}
			fmt.Println(string(content[:n]))
		}
```
	- `os.OpenFile`使用
```go
		// os.Open -> 读文件，文件不存在则报错
		func Open(name string) (*File, error) {
			return os.OpenFile(name, os.O_RDONLY, 0777)
		}
		// os.Create -> 写文件，文件存在 截断，文件不存在 创建
		func Create(name string) (*File, error) {
			return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		}
		// 示例
		file, err := os.OpenFile("test.txt", os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
```



### 2. 标准输入、标准输出、标准错误
- 标准输入 -> `os.Stdin`
	- `os.Stdin` 直接读取
```go
		content := make([]byte, 3)
		
		fmt.Print("请输入内容: ")
		fmt.Println(os.Stdin.Read(content))
		fmt.Printf("%q\n", string(content))
```
	- `os.Stdin` 作为脚本的输入内容
```go
		// os/exec exec包提供了启动一个外部进程并使用标准输入和输出进行通信
		// 常用结构体
		// Cmd: 执行命令 
		// 常用函数
		// Command
		func Command(name string, arg ...string) *Cmd
		// 常用方法
		// Output: 执行并获取标准输出结果
		// Run: 自行命令
		func (c *Cmd) Run() error
		// Start: 启动命令
		func (c *Cmd) Start() error
		// Wait: 与 Start一起使用等待命令结束
		// StdoutPipe: 输出管道
		func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
		// StdinPipe: 输入管道
		func (c *Cmd) StdinPipe() (io.WriteCloser, error)
		// StderrPipe: 错误管道
		func (c *Cmd) StderrPipe() (io.ReadCloser, error)

		// 写一个脚本，和命令一个输入的文件
		// 文件作为脚本的stdin，执行
		// echo "ss -ntlp " > a.txt
		// go run a.go < a.txt
		cmd := exec.Command("sh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("run.err", err)
			return
		}
```

- 标准输出 -> `os.Stdout`
```go
		os.Stdout.WriteString("我是Stdout的输出")
		fmt.Fprintln(os.Stdout, "aaaaa")
		fmt.Fprintf(os.Stdout, "I am: %s", "aaaaa")

		// os.StdOut.Write 代替 fmt.Print
		os.Stdout.Write([]byte("aa和bb"))
		fmt.Println("aa和bb")
```

- 标准错误  -> `os.Stderr`


### 3. 带缓冲的IO
- bufio包提供缓冲流的功能
	- 常用结构体
		- Reader
			- 常用函数 
				`NewReader`: 创建缓冲 输入 流
			- 常用方法
				- `Read`: 读取数据到切片中
				- `ReadLine`: 读取一行内容到字节切片中
				- `ReadSlice`: 根据分隔符读取数据到字节切片
				- `ReadString`: 根据分隔符读取数据到字符串
				- `Reset`: 重设缓冲流
				- `WriteTo`: 将数据写入到输出流
		- Scanner
			- 常用函数 
				`NewScanner`: 创建扫描对象
			- 常用方法
				- `Scan`: 扫描数据
				- `Split`: 定义流分割函数，默认 空格
				- `Text`: 读取数据
				- `Err`: 获取错误
		- Writer
			- 常用函数 
				`NewWriter`: 创建缓冲输出流
			- 常用方法
				- `Write`: 将字节切片内容写入
				- `WriteString`: 将字符串写入
				- `Reset`: 重置输出流
				- `Flush`: 刷新数据到输出流
	- 示例
```go
		scanner := bufio.NewScanner(os.Stdin)        // os.Stdin
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			break
		}
		
		func ScanInt() (int, error) {
			// 读取一行 进行转换
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				return strconv.Atoi(scanner.Text())
			}
			return 0, scanner.Err()
		}
		num, err := ScanInt()
		fmt.Println(num, err)
```

### 4. IO库
- io库属于底层接口定义库，作用是定义一些基本接口和基本常量，如io.EOF

- io.Copy `func io.Copy(dst io.Writer, src io.Reader) (written int64, err error)`

- io库常用接口有:`Reader`、`Writer`、`Close` 以流的方式高效处理数据，并不需要考虑数据是什么，数据来自哪里，数据要发送到哪里
	- Reader
		- `io.Reader`表示一个读取器
			- 它从某个地方读取数据到传输的缓存区
			- 在缓存区里面，数据可以被流式的使用
			- 接口签名
```go
				type Reader interface {
					Read(p []byte) (n int, err error)
				}
```
		- `strings.NewReader`
			- `io.Reader` 接口只有一个方法: Read方法
			- 即只要有个对象实现了Read方法，那么这个对象就是一个读取器
			- `Read()` 首先要有一个读缓冲区的参数
			- `Read()` 返回两个值，第一个是读取到的字节数，第二个是读取时发生的错误 `func (*strings.Reader).Read(b []byte) (n int, err error)`
			- 注意: 返回到的读取字节个数n可能小于缓冲区的大小
			- io.EOF 表示输入的流已经读到头了
```go
				// 实现一个 reader 每次读取4个字节
				// 从字符串创建一个reader对象
				reader := strings.NewReader("马哥教育 2021 第005期 golang")
				// new一个3字节的读取缓冲区
				p := make([]byte, 3)
				for {
					// reader对象读取数据
					n, err := reader.Read(p)
					if err != nil {
						if err == io.EOF {
							log.Printf("[数据已读完 EOF:%d]", n)
							break
						}
						log.Printf("[未知错误:%v]", err)
						return
					}
					log.Printf("[打印读取的字节数:%d 内容:%s]", n, string(p[:n]))
				}
```
		- 自定义Reader
			- 要求: 过滤输入字符串中的非字母字符
			- 输入 "mage jiaoyue 2021 go !!!!"
			- 输出 "magejiaoyuego"
```go
				type zimuguolv struct {
					src string
					cur int
				}
				func alpha(r byte) byte {
					// r在 A-Z 或者 a-z
					if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
						return r
					}
					return 0
				}
				func (z *zimuguolv) Read(p []byte) (int, error) {
					// 当前位置 >= 字符串长度，说明已经读取到结尾了，返回 EOF
					if z.cur >= len(z.src) {
						return 0, io.EOF
					}
					// 定义一个剩余还没读到的长度
					x := len(z.src) - z.cur
					// bound叫做本次读取长度
					// n代表本次遍历 bound的索引
					n, bound := 0, 0
					if x >= len(p) {
						// 剩余长度超过缓冲区大小，说明本次可以完全填满换冲区
						bound = len(p)
					} else {
						// 剩余长度小于缓冲区大小，使用剩余长度输出，缓冲区填不满
						bound = x
					}
				
					buf := make([]byte, bound)
				
					for n < bound {
						if char := alpha(z.src[z.cur]); char != 0 {
							buf[n] = char
						}
						// 索引++
						n++
						z.cur++
					}
					copy(p, buf)
					return n, nil
				}
				zmreader := zimuguolv{
					src: "mage jiaoyu 2021 go !!!!",
				}
				p := make([]byte, 4)
				for {
					n, err := zmreader.Read(p)
					if err == io.EOF {
						log.Printf("[EOF错误]")
						break
					}
					log.Printf("[读取到的长度%d 内容%s]", n, string(p[:n]))
				}
```
		- 组合多个Reader
			- 标准库里面已经有了很多Reader
			- 使用一个Reader A作为一个Reader B的一部分
			- 目的是重用和屏蔽下层实现的复杂度；即复用逻辑，流式处理
			- 复用的`io.Reader`
```go
				type alphaReader struct {
					ioReader io.Reader
				}
				func (a *alphaReader) Read(p []byte) (int, error) {
					// 复用io.reader的read方法
					n, err := a.ioReader.Read(p)
					if err != nil {
						return n, err
					}
				
					buf := make([]byte, n)
					for i := 0; i < n; i++ {
						if char := alpha(p[i]); char != 0 {
							buf[i] = char
						}
					}
					copy(p, buf)
					return n, nil
				}
				func alpha(r byte) byte {
					// r在 A-Z 或者 a-z
					if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
						return r
					}
					return 0
				}
				myReader := alphaReader{
					strings.NewReader("mage jiaoyu 2021 go !!!"),
				}
				p := make([]byte, 4)
				for {
					n, err := myReader.Read(p)
					if err == io.EOF {
						log.Printf("[EOF错误]")
						break
					}
					log.Printf("[读取到的长度%d 内容%s]", n, string(p[:n]))
				}
```
		- os.File 结合
			- os.Open得到一个file对象 ，实现了io.Reader的Read方法
			- 以下代码展示了 alphaReader 如何与 os.File 结合以过滤掉文件中的非字母字符
```go
				file, err := os.Open("test.txt")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()
				myReader := alphaReader{
					file,
				}
				p := make([]byte, 4)
				for {
					n, err := myReader.Read(p)
					if err == io.EOF {
						log.Printf("[EOF错误]")
						break
					}
					log.Printf("[读取到的长度%d 内容%s]", n, string(p[:n]))
				}
```
		- 读取文件 `ioutil.ReadFile` vs `bufio`
			- 都提供了文件读写的能力
			- bufio多了一层缓存的能力，优势体现在读取大文件的时候
			- `ioutil.ReadFile`是一次性将内容加载到内存，大文件容易爆掉
```go
				fileName := "test.txt"
				
				// ioutil.ReadFile
				bytes, err := ioutil.ReadFile(fileName)
				if err != nil {
					return
				}
				
				// os.Open + ioutil.ReadAll  
				// func ReadAll(r io.Reader) ([]byte, error)
				file, err := os.Open(fileName)
				if err != nil {
					return
				}
				bytes, err = ioutil.ReadAll(file)
				if err != nil {
					return
				}
				file.Close()
				
				// os.Open + file.Read
				file, _ = os.Open(fileName)
				buf := make([]byte, 50)
				_, err = file.Read(buf)
				if err != nil {
					return
				}
				file.Close()
				
				// os.Open + bufio.Read
				file, _ = os.Open(fileName)
				rd := bufio.NewReader(file)  	// bufio.NewReader
				buf1 := make([]byte, 50)
				_, err = rd.Read(buf1)
				if err != nil {
					return
				}
				file.Close()
```
	- Writer 
		- `io.Writer` 表示一个编写器，它从缓冲区读取数据，并将数据写入目标资源
			- 接口签名
```go
				type Writer interface {
					Write(p []byte) (n int, err error)
				}
```
		- Write() 方法有两个返回值，一个是写入到目标资源的字节数，一个是发生错误时的错误。
			- closer
			- `bytes.Buffer`库 
				- bytes.Buffer 的针对的是内存到内存的缓存
			- `io/ioutil` ioutil库 工具包
				- 在io目录下，它是一个工具包，实现一些实用的工具
```go
				// ReadFile 读取文件                           
				// func ReadFile(filename string) ([]byte, error)
				fileName := "golang.txt"
				bytes, err := ioutil.ReadFile(fileName)
				if err != nil {
					fmt.Println(err)
					return
				}

				// WriteFile 写入文件
				// func WriteFile(filename string, data []byte, perm fs.FileMode) error
				fileName := "test.txt"
				err := ioutil.WriteFile(fileName, []byte("123\n456"), 0644)
				fmt.Println(err)

				// ReadDir 读取目录下的文件元信息
				// func ReadDir(dirname string) ([]fs.FileInfo, error)
				fs, err := ioutil.ReadDir("./")
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, f := range fs {
					fmt.Printf("[name:%v][size:%v][isDir:%v][mode:%v][ModTime:%v]\n",
						f.Name(),
						f.Size(),
						f.IsDir(),
						f.Mode(),
						f.ModTime(),
					)
				}
```

## 三、目录

- Go语言中的os包，提供了对目录的操作
	- 创建 
		`os.Mkdir("a", os.ModePerm)`，
		`os.MkdirAll("a/b/c", os.ModePerm)`
	- 读取 `os.Open("test.txt")`
	- 获取属性 `os.Open().Stat`/`os.Stat `
	- 修改属性 -> 权限，所属人 
		`os.Chmod()`，
		`os.Chown()`
	- 重命名 `fmt.Println(os.Rename("b", "d:\\d"))`
	- 删除文件夹 
		`os.Remove("a")`，
		`os.RemoveAll("a")`
	- FileInfo: 文件状态信息
		- 常用函数
			- `Lstat`: 获取文件路径文件信息（对于链接返回连接文件信息）
			- `Stat`: 获取文件路径文件信息（对于链接返回连接到的文件的信息）
		- 常用方法
			- `Name`: 获取文件名
			- `Size`: 获取文件大小
			- `Mode`: 获取文件模式 `func (fs.FileInfo).Mode() fs.FileMode`
			- `ModTime`: 获取修改时间
			- `IsDir`: 判断是否为文件夹  
	- FileMode: 文件模式
		- 常用方法
			`IsDir`: 判断是否为文件夹 `func (fs.FileMode).IsDir() bool`

## 四、编码格式

- 处理流程
	- 注册，打开文件，创建对象，编码/解码

- gob
	- go特有的编码格式，不能跨语言
	- `encoding/gob` 包提供了对数据结构进行二进制序列化的功能
		- 常用函数
			- `Register`: 注册 gob 编解码记录值 `func Register(value interface{})`
			- `RegisterName`: 注册 gob 编解码记录值，并指定名称 `func RegisterName(name string, value interface{})`
		- 常用结构体
			- Encoder `type Encoder struct{ ... }`
				- 常用函数 
					- `NewEncoder`: 创建编码器 `func NewEncoder(w io.Writer) *Encoder`
				- 常用方法
					- `Encode`: 将对象进行编码到流对象中 `func (enc *Encoder) Encode(e interface{}) error`
			- Decoder `type Decoder struct{ ... }`
				- 常用函数
					- `NewDecoder`: 创建解码器 `func NewDecoder(r io.Reader) *Decoder`
				- 常用方法
					- `Decode`: 将流对象中的数据编码到对象中 `func (dec *Decoder) Decode(e interface{}) error`
```go
		type User struct {
			Id   int
			Name string
		}
		enusers := []User{
			{1, "aa"},
			{2, "bb"},
		}
		// 注册
		gob.Register(User{})
		// 编码
		file, err := os.Create("users.gob")
		if err != nil {
			fmt.Println(err)
			return
		}
		encoder := gob.NewEncoder(file)
		fmt.Println(encoder.Encode(enusers))
		file.Close()
		// 解码
		file, err := os.Open("users.gob")
		if err != nil {
			return
		}
		decoder := gob.NewDecoder(file)
		var deusers []User
		fmt.Println(decoder.Decode(&deusers))
		fmt.Println(deusers)
		file.Close()
```
		
- csv
	- `encoding/csv` 包提供对 csv 文件读写的操作
	- 常用结构体
		- Reader
			- 常用函数
				- `func NewReader(r io.Reader) *Reader`
			- 常用方法
				- `func (r *Reader) Read() (record []string, err error)`
				- `func (r *Reader) ReadAll() (records [][]string, err error)`
		- Writer
			- 常用函数
				- `func NewWriter(w io.Writer) *Writer`
			- 常用方法
				- `func (w *Writer) WriteAll(records [][]string) error`
				- `func (w *Writer) Write(record []string) error`
				- `func (w *Writer) Flush()`
				- `func (w *Writer) Error() error`
```go
		wusers := []User{
			{1, "aa"},
			{2, "bb"},
		}
		// 写入
		file, err := os.Create("users.csv")
		if err != nil {
			return
		}
		writer := csv.NewWriter(file)
		for _, user := range wusers {
			writer.Write([]string{strconv.Itoa(user.Id), user.Name})
		}
		writer.Flush()
		file.Close()
		// 读取
		file, err = os.Open("users.csv")
		if err != nil {
			return
		}
		var rusers []User
		reader := csv.NewReader(file)
		for {
			line, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			id, _ := strconv.Atoi(line[0])
			rusers = append(rusers, User{id, line[1]})
		}
		fmt.Println(rusers)
```

## 五、参考范例

- 真实生产应用 
	- 夜莺监控发送告警，调用python的send.py脚本，将发送的内容作为stdin传过去
    - [go代码](https://github.com/didi/nightingale/blob/master/alert/consume.go#L183)
    - [python 脚本](https://github.com/didi/nightingale/blob/master/etc/script/notify.py)
