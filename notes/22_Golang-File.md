# Golang-File  Golang的文件操作

## 一、基本概念

- 绝对路径
	- 文件路径字符串从根路径(盘符)开始
	
- 相对路径
	- 二进制文件运行的目录 `cd $directory`
	- 相对程序执行的路径，即当前shell 处于的路径

- 文件类型
	- cat/记事本 -> 文本内容(无乱码 .go .txt) -> 文本文件 -> `string`
	- cat/记事本 -> 有乱码(word, zip, excel) -> 二进制文件 -> `[]byte`

- I/O操作
	- 也叫输入/输出操作，其中I是指Input，O是指Output，主要用来读取或写入数据，很多语言中也叫做流操作
	- Go语言中 输入和输出操作是使用原语实现的
		- 这些原语将数据模拟成可以读或者可以写的字节流
		- DataSource -> `io.Reader` -> `Transfer buffer []byte` -> `io.Writer` -> Target


## 二、文件与目录

### 1. 基本操作(不带缓冲IO): 读、写
- 处理流程
	- 打开文件 -> 错误处理 -> 延迟关闭 -> 读/写/其他 -> 关闭文件

- Go语言中的os包
	- os包提供了操作系统函数的不依赖平台的接口
	- 设计为Unix风格的，虽然错误处理是go风格的；失败的调用会返回错误值而非错误码
	- 通常错误值里包含更多信息，例如，如果某个使用一个文件名的调用(如Open、Stat)失败了，打印错误时会包含该文件名，错误类型将为*PathError，其内部可以解包获得更多信息
	- os包的接口规定为在所有操作系统中都是一致的；非公用的属性可以从操作系统特定的syscall包获取

- os包中，标准输入、标准输出、标准错误
	- `os.Stdin` 标准输入 
	- `os.Stdout` 标准输出
	- `os.Stderr` 标准错误
```go
	// Stdin、Stdout和Stderr是指向标准输入、标准输出、标准错误输出的文件描述符
	var (
		Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
		Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
		Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
	)

	// os.Stdin 直接读取 
	content := make([]byte, 3)
	os.Stdin.Read(content)

	// os.StdOut.WriteString
	os.Stdout.WriteString("我是Stdout的输出")
	fmt.Fprintln(os.Stdout, "aaaaa")
	fmt.Fprintf(os.Stdout, "I am: %s", "aaaaa")

	// os.StdOut.Write 代替 fmt.Print
	os.Stdout.Write([]byte("aa和bb"))
	fmt.Println("aa和bb")
```

- os包中，操作系统信号
```go
	// Signal代表一个操作系统信号
	// 一般其底层实现是依赖于操作系统的：在Unix中，它是syscall.Signal类型
	type Signal interface {
		String() string
		Signal() // 用来区分其他实现了Stringer接口的类型
	}

	// 仅有的肯定会被所有操作系统提供的信号，Interrupt(中断信号) 和Kill (强制退出信号)
	var (
		Interrupt Signal = syscall.SIGINT
		Kill      Signal = syscall.SIGKILL
	)
```

- os包中，常用 对系统、进程的操作函数
	- `os.Hostname()` 获取主机名
	- `os.Getenv(key)` Getenv检索并返回名为key的环境变量的值
	- `os.Setenv(key)` Setenv设置名为key的环境变量
	- `os.Environ()` 获取所有环境变量
	- `os.Clearenv()` 删除所有环境变量
	- `os.Exit(code)` Exit让当前程序以给出的状态码code退出;一般来说，状态码0表示成功，非0表示出错;程序会立刻终止，defer的函数不会被执行
	- `os.Getuid()` 返回调用者的用户ID
	- `os.Getpid()` 返回调用者所在进程的进程ID
	- `os.Getppid()` 返回调用者所在进程的父进程的进程ID
	- 常用常量
```go
	// os.DevNull是操作系统空设备的名字
	// 在类似Unix的操作系统中，是"/dev/null"；在Windows中，为"NUL"
	const DevNull = "/dev/null"

	// os.Args 保管了命令行参数，第一个是程序名
	var Args []string

	// flag
	// 用于包装底层系统的参数用于Open函数，不是所有的flag都能在特定系统里使用的
	const (
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
	)
```

- os包中，常用 对目录和文件的操作函数
	- `os.Getwd()` 获取当前目录
	- `os.Chdir("dir")` 将当前工作目录修改为dir指定的目录
	- `os.Mkdir("a", os.ModePerm)` 创建文件夹
	- `os.MkdirAll("a/b/c", os.ModePerm)` 创建文件夹(父目录不存在逐层创建)
	- `os.IsExist(err)` 与 `os.Stat` 一起用于判断文件存在
	- `os.IsNotExist(err)` 与 `os.Stat` 一起用于判断文件不存在
	- `os.Link("oldname", "newname")` 创建一个名为newname指向oldname的硬链
	- `os.Create()` 创建文件
	- `os.Open("test.txt")` 读取文件
	- `os.Open().Stat`、`os.Stat` 获取文件属性 
	- `os.Chmod(name string, mode FileMode)` 修改文件权限
	- `os.Chown(name string, uid, gid int)` 修改文件所属用户，用户组
	- `os.Chtimes(name string, atime time.Time, mtime time.Time)` 修改文件访问时间和修改时间
	- `os.Truncate(name string, size int64)` 修改name指定的文件的大小
	- `os.Rename("a.txt", "b.txt")` 重命名文件
	- `os.Remove("b.txt")` 删除指定的文件或目录
	- `os.RemoveAll("path")` 删除path指定的文件，或目录及它包含的任何下级对象
	- `os.TempDir()` 返回一个用于保管临时文件的默认目录
	- `os.FileInfo` 用来描述一个文件对象
	- `os.FileMode` 代表文件的模式和权限位
```go
	type FileInfo interface {
		Name() string       // 文件的名字（不含扩展名）
		Size() int64        // 普通文件返回值表示其大小；其他文件的返回值含义各系统不同
		Mode() FileMode     // 文件的模式位
		ModTime() time.Time // 文件的修改时间
		IsDir() bool        // 等价于Mode().IsDir()
		Sys() interface{}   // 底层数据来源（可以返回nil）
	}

	// os.Stat返回一个描述name指定的文件对象的FileInfo
	// 如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接指向的文件的信息，本函数会尝试跳转该链接
	func Stat(name string) (fi FileInfo, err error)

	// os.Lstat返回一个描述name指定的文件对象的FileInfo
	// 如果指定的文件对象是一个符号链接，返回的FileInfo描述该符号链接的信息，本函数不会试图跳转该链接
	func Lstat(name string) (fi FileInfo, err error)

	// FileMode代表文件的模式和权限位
	// 这些字位在所有的操作系统都有相同的含义，因此文件的信息可以在不同的操作系统之间安全的移植
	// 不是所有的位都能用于所有的系统，唯一共有的是用于表示目录的ModeDir位
	type FileMode uint32

	const (
		// 单字符是被String方法用于格式化的属性缩写。
		ModeDir        FileMode = 1 << (32 - 1 - iota) // d: 目录
		ModeAppend                                     // a: 只能写入，且只能写入到末尾
		ModeExclusive                                  // l: 用于执行
		ModeTemporary                                  // T: 临时文件（非备份文件）
		ModeSymlink                                    // L: 符号链接（不是快捷方式文件）
		ModeDevice                                     // D: 设备
		ModeNamedPipe                                  // p: 命名管道（FIFO）
		ModeSocket                                     // S: Unix域socket
		ModeSetuid                                     // u: 表示文件具有其创建者用户id权限
		ModeSetgid                                     // g: 表示文件具有其创建者组id的权限
		ModeCharDevice                                 // c: 字符设备，需已设置ModeDevice
		ModeSticky                                     // t: 只有root/创建者能删除/移动文件
		// 覆盖所有类型位（用于通过&获取类型位），对普通文件，所有这些位都不应被设置
		ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
		ModePerm FileMode = 0777 // 覆盖所有Unix权限位（用于通过&获取类型位）
	)

	// IsDir报告m是否是一个目录
	func (m FileMode) IsDir() bool

	// IsRegular报告m是否是一个普通文件
	func (m FileMode) IsRegular() bool

	// Perm返回m的Unix权限位
	func (m FileMode) Perm() FileMode

	// String返回m的string格式
	func (m FileMode) String() string
```

- os包中，常用的结构体 File
	- `os.Create(name)` 创建文件并返回文件对象指针(文件不存在则创建，文件存在则清空)
	- `os.Open(name)` 打开文件并返回文件对象指针
	- `os.OpenFile(name, flag, perm)` 按指定权限打开文件，并返回文件指针对象
```go
	// os.File 代表一个打开的文件对象
	type File struct {
		// Has unexported fields.
	}
         	  
	// 常用函数
	// os.Create采用模式0666(任何人都可读写，不可执行)创建一个名为name的文件，如果文件已存在会截断它(为空文件)
	// 如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式；如果出错，错误底层类型是*PathError
	func Create(name string) (*os.File, error)
	// os.Open打开一个文件用于读取
	// 如果操作成功，返回的文件对象的方法可用于读取数据；对应的文件描述符具有O_RDONLY模式
	func Open(name string) (*os.File, error)
	// os.OpenFile 按指定权限打开文件，并返回文件指针对象
	func OpenFile(name string, flag int, perm FileMode) (*File, error)
	// os.NewFile 使用给出的Unix文件描述符和名称创建一个文件
	func NewFile(fd uintptr, name string) *File
	// os.Pipe 返回一对关联的文件对象，从r的读取将返回写入w的数据，本函数会返回两个文件对象和可能的错误
	func Pipe() (r *File, w *File, err error)
		
	// 常用方法
	// Name 返回文件名称
	func (f *File) Name() string
	// Stat 返回描述文件f的FileInfo类型值
	func (f *File) Stat() (fi FileInfo, err error)
	// Fd 返回与文件f对应的整数类型的Unix文件描述符
	func (f *File) Fd() uintptr
	// Chdir 将当前工作目录修改为f，f必须是一个目录
	func (f *File) Chdir() error
	// Chmod 修改文件的模式
	func (f *File) Chmod(mode FileMode) error
	// Chown 修改文件的用户ID和组ID
	func (f *File) Chown(uid, gid int) error

	// Readdir读取目录f的内容，返回一个有n个成员的[]FileInfo，这些FileInfo是被Lstat返回的，采用目录顺序
	func (f *File) Readdir(n int) (fi []FileInfo, err error)
	// Readdir读取目录f的内容，返回一个有n个成员的[]string，切片成员为目录中文件对象的名字，采用目录顺序
	func (f *File) Readdirnames(n int) (names []string, err error)

	// Truncate 改变文件的大小，它不会改变I/O的当前位置；如果截断文件，多出的部分就会被丢弃
	func (f *File) Truncate(size int64) error
	// Read 从f中读取最多len(b)字节数据并写入b，它返回读取的字节数和可能遇到的任何错误，文件终止标志是读取0个字节且返回值err为io.EOF
	func (f *File) Read(b []byte) (n int, err error)
	// ReadAt从指定的位置（相对于文件开始位置）读取len(b)字节数据并写入b
	func (f *File) ReadAt(b []byte, off int64) (n int, err error)
	// Write 向文件中写入len(b)字节数据，它返回写入的字节数和可能遇到的任何错误，如果返回值n!=len(b)，本方法会返回一个非nil的错误
	func (f *File) Write(b []byte) (n int, err error)
	// WriteString 接受一个字符串参数，写入字符串到文件
	func (f *File) WriteString(s string) (ret int, err error)
	// WriteAt 在指定的位置（相对于文件开始位置）写入len(b)字节数据
	func (f *File) WriteAt(b []byte, off int64) (n int, err error)
	// Seek 设置文件指针位置，即下一次读/写的位置
	// offset为相对偏移量，而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾，它返回新的偏移量(相对开头)和可能的错误
	// 在大部分编程语言中，不支持在文件的开始或中间插入数据(会直接从光标位置进行覆写)，只支持在文件末尾进行数据追加
	func (f *File) Seek(offset int64, whence int) (ret int64, err error)

	// Sync 递交文件的当前内容进行稳定的存储；一般来说，这表示将文件系统的最近写入的数据在内存中的拷贝刷新到硬盘中稳定保存  
	func (f *File) Sync() (err error)
	// Close 关闭文件f，使文件不能用于读写
	func (f *File) Close() error
```

- 参考示例
```go
	// 创建和写入文件
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

	// 读文件
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

	// os.Open -> 读文件，文件不存在则报错
	func Open(name string) (*File, error) {
		return os.OpenFile(name, os.O_RDONLY, 0777)
	}
	// os.Create -> 写文件，文件存在 截断，文件不存在 创建
	func Create(name string) (*File, error) {
		return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	}
	// os.OpenFile 
	file, err := os.OpenFile("test.txt", os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
```

- Go语言中的os/exec包
	- exec包执行外部命令
	- 它包装了os.StartProcess函数以便更容易的修正输入和输出，使用管道连接I/O，以及作其它的一些调整
```go
	// exec.LookPath 在环境变量PATH指定的目录中搜索可执行文件，如file中有斜杠，则只在当前目录搜索；返回完整路径或者相对于当前目录的一个相对路径。
	func LookPath(file string) (string, error)

	// Cmd 代表一个正在准备或者在执行中的外部命令
	type Cmd struct {
		// Path是将要执行的命令的路径
		Path string
		// Args保管命令的参数，包括命令名作为第一个参数；如果为空切片或者nil，相当于无参数命令
		Args []string
		// Env指定进程的环境，如为nil，则是在当前进程的环境下执行
		Env []string
		// Dir指定命令的工作目录。如为空字符串，会在调用者的进程当前目录下执行
		Dir string
		// Stdin指定进程的标准输入，如为nil，进程会从空设备读取（os.DevNull）

		Stdin io.Reader
		// Stdout和Stderr指定进程的标准输出和标准错误输出。
		// 如果任一个为nil，Run方法会将对应的文件描述符关联到空设备（os.DevNull）
		// 如果两个字段相同，同一时间最多有一个线程可以写入
		Stdout io.Writer
		Stderr io.Writer

		// ExtraFiles指定额外被新进程继承的已打开文件流，不包括标准输入、标准输出、标准错误输出
		// 如果本字段非nil，entry i会变成文件描述符3+i
		ExtraFiles []*os.File
		// SysProcAttr保管可选的、各操作系统特定的sys执行属性
		// Run方法会将它作为os.ProcAttr的Sys字段传递给os.StartProcess函数
		SysProcAttr *syscall.SysProcAttr
		// Process是底层的，只执行一次的进程
		Process *os.Process
		// ProcessState包含一个已经存在的进程的信息，只有在调用Wait或Run后才可用
	}

	// 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序；返回值只设定了Path和Args两个参数
	// 如果name不含路径分隔符，将使用LookPath获取完整路径；否则直接使用name；参数arg不应包含命令名
	func Command(name string, arg ...string) *Cmd

	// Run 执行c包含的命令，并阻塞直到完成
	func (c *Cmd) Run() error
	// Start 开始执行c包含的命令，但并不会等待该命令完成即返回
	func (c *Cmd) Start() error
	// Wait会阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	func (c *Cmd) Wait() error
	
	// Output 执行命令并返回标准输出的切片
	func (c *Cmd) Output() ([]byte, error)
	// CombinedOutput 执行命令并返回标准输出和错误输出合并的切片
	func (c *Cmd) CombinedOutput() ([]byte, error)

	// StdinPipe 返回一个在命令Start后与命令标准输入关联的管道
	func (c *Cmd) StdinPipe() (io.WriteCloser, error)
	// StdoutPipe 返回一个在命令Start后与命令标准输出关联的管道
	func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
	// StderrPipe 返回一个在命令Start后与命令标准错误输出关联的管道
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

### 2. 带缓冲IO 的 读、写
- bufio包
	- 实现了有缓冲的I/O，提供缓冲流的功能
	- 它包装一个 `io.Reader` 或 `io.Writer` 接口对象，创建另一个也实现了该接口，且同时还提供了缓冲和一些文本I/O的帮助函数的对象

- bufio包中，常用的结构体
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
```go
	const (
		// 用于缓冲一个token，实际需要的最大token尺寸可能小一些，例如缓冲中需要保存一整行内容
		MaxScanTokenSize = 64 * 1024
	)

	// bufio.NewReader创建一个具有默认大小缓冲、从r读取的*Reader
	func NewReader(rd io.Reader) *Reader

	// bufio.NewReaderSize创建一个具有最少有size尺寸的缓冲、从r读取的*Reader
	// 如果参数r已经是一个具有足够大缓冲的* Reader类型值，会返回r
	func NewReaderSize(rd io.Reader, size int) *Reader

	// bufio.Reader实现了给一个io.Reader接口对象附加缓冲
	type Reader struct { ... }

	// Reset丢弃缓冲中的数据，清除任何错误，将b重设为其下层从r读取数据
	func (b *Reader) Reset(r io.Reader)

	// Buffered返回缓冲中现有的可读取的字节数
	func (b *Reader) Buffered() int

	// Peek返回输入流的下n个字节，而不会移动读取位置
	func (b *Reader) Peek(n int) ([]byte, error)

	// Read读取数据写入p，本方法返回写入p的字节数
	// 本方法一次调用最多会调用下层Reader接口一次Read方法，因此返回值n可能小于len(p)；读取到达结尾时，返回值n将为0而err将为io.EOF
	func (b *Reader) Read(p []byte) (n int, err error)

	// ReadByte读取并返回一个字节。如果没有可用的数据，会返回错误
	func (b *Reader) ReadByte() (c byte, err error)

	// UnreadByte吐出最近一次读取操作读取的最后一个字节(只能吐出最后一个，多次调用会出问题)
	func (b *Reader) UnreadByte() error

	// ReadRune读取一个utf-8编码的unicode码值，返回该码值、其编码长度和可能的错误
	func (b *Reader) ReadRune() (r rune, size int, err error)

	// UnreadRune吐出最近一次ReadRune调用读取的unicode码值
	func (b *Reader) UnreadRune() error

	// ReadLine是一个低水平的行数据读取原语；大多数调用者应使用ReadBytes('\n')或ReadString('\n')代替，或者使用Scanner
	func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)

	// ReadSlice读取直到第一次遇到delim字节，返回缓冲里的包含已读取的数据和delim字节的切片
	func (b *Reader) ReadSlice(delim byte) (line []byte, err error)

	// ReadBytes读取直到第一次遇到delim字节，返回一个包含已读取的数据和delim字节的切片
	func (b *Reader) ReadBytes(delim byte) (line []byte, err error)

	// ReadString读取直到第一次遇到delim字节，返回一个包含已读取的数据和delim字节的字符串
	func (b *Reader) ReadString(delim byte) (line string, err error)

	// WriteTo方法实现了io.WriterTo接口
	func (b *Reader) WriteTo(w io.Writer) (n int64, err error)

	// bufio.NewWriter创建一个具有默认大小缓冲、写入w的*Writer
	func NewWriter(w io.Writer) *Writer

	// bufio.Writer实现了为io.Writer接口对象提供缓冲
	// 如果在向一个Writer类型值写入时遇到了错误，该对象将不再接受任何数据，且所有写操作都会返回该错误
	// 在说有数据都写入后，调用者有义务调用Flush方法以保证所有的数据都交给了下层的io.Writer
	// w := bufio.NewWriter(os.Stdout);w.Flush()
	type Writer struct { ... }

	// Reset丢弃缓冲中的数据，清除任何错误，将b重设为将其输出写入w
	func (b *Writer) Reset(w io.Writer)

	// Buffered返回缓冲中已使用的字节数
	func (b *Writer) Buffered() int

	// Available返回缓冲中还有多少字节未使用
	func (b *Writer) Available() int
	
	// Write将p的内容写入缓冲，返回写入的字节数
	func (b *Writer) Write(p []byte) (nn int, err error)

	// WriteString写入一个字符串，返回写入的字节数
	func (b *Writer) WriteString(s string) (int, error)

	// WriteByte写入单个字节
	func (b *Writer) WriteByte(c byte) error

	// WriteRune写入一个unicode码值（的utf-8编码），返回写入的字节数和可能的错误
	func (b *Writer) WriteRune(r rune) (size int, err error)

	// Flush方法将缓冲中的数据写入下层的io.Writer接口
	func (b *Writer) Flush() error

	// ReadFrom实现了io.ReaderFrom接口
	func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)

	// bufio.NewReadWriter申请创建一个新的、将读写操作分派给r和w 的ReadWriter
	func NewReadWriter(r *Reader, w *Writer) *ReadWriter

	// bufio.ReadWriter类型保管了指向Reader和Writer类型的指针，实现了io.ReadWriter接口
	type ReadWriter struct {
		*Reader
		*Writer
	}

	// bufio.SplitFunc类型代表用于对输出作词法分析的分割函数
	type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

	// bufio.ScanBytes是用于Scanner类型的分割函数(符合SplitFunc)，本函数会将每个字节作为一个token返回
	func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)

	// bufio.NewScanner创建并返回一个从r读取数据的Scanner，默认的分割函数是ScanLines
	func NewScanner(r io.Reader) *Scanner

	// bufio.Scanner类型提供了方便的读取数据的接口，如从换行符分隔的文本里读取每一行
	type Scanner struct { ... }

	// Split设置该Scanner的分割函数；本方法必须在Scan之前调用
	func (s *Scanner) Split(split SplitFunc)

	// Scan方法获取当前位置的token(该token可以通过Bytes或Text方法获得)，并让Scanner的扫描位置移动到下一个token
	// 当扫描因为抵达输入流结尾或者遇到错误而停止时，本方法会返回false；在Scan方法返回false后，Err方法将返回扫描时遇到的任何错误；除非是io.EOF，此时Err会返回nil
	func (s *Scanner) Scan() bool

	// Bytes方法返回最近一次Scan调用生成的token；底层数组指向的数据可能会被下一次Scan的调用重写
	func (s *Scanner) Bytes() []byte

	// Bytes方法返回最近一次Scan调用生成的token，会申请创建一个字符串保存token并返回该字符串
	func (s *Scanner) Text() string

	// Err返回Scanner遇到的第一个非EOF的错误
	func (s *Scanner) Err() error
```

- 参考示例
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

### 3. IO包
- io包
	- io包提供了对I/O原语的基本接口
	- 本包的基本任务是包装这些原语已有的实现(如os包里的原语)，使之成为共享的公共接口，这些公共接口抽象出了泛用的函数并附加了一些相关的原语的操作
```go
	// EOF当无法得到更多输入时，Read方法返回EO
	var EOF = errors.New("EOF")

	// 当从一个已关闭的Pipe读取或者写入时，会返回ErrClosedPipe
	var ErrClosedPipe = errors.New("io: read/write on closed pipe")

	// 某些使用io.Reader接口的客户端如果多次调用Read都不返回数据也不返回错误时，就会返回本错误，一般来说是io.Reader的实现有问题的标志
	var ErrNoProgress = errors.New("multiple Read calls return no data or error")

	// ErrShortBuffer表示读取操作需要大缓冲，但提供的缓冲不够大
	var ErrShortBuffer = errors.New("short buffer")

	// ErrShortWrite表示写入操作写入的数据比提供的少，却没有显式的返回错误
	var ErrShortWrite = errors.New("short write")

	// ErrUnexpectedEOF表示在读取一个固定尺寸的块或者数据结构时，在读取未完全时遇到了EOF
	var ErrUnexpectedEOF = errors.New("unexpected EOF")

	// io包中的常用接口，以流的方式高效处理数据，并不需要考虑数据是什么，数据来自哪里，数据要发送到哪里
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
	- ioutil包提供了一些I/O实用函数
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

### 4. Reader
- Reader
	- `io.Reader`表示一个读取器
		- 它从某个地方读取数据到传输的缓存区
		- 在缓存区里面，数据可以被流式的使用
```go
	// 接口签名
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

### 5. Writer
- Writer 
	- `io.Writer` 表示一个编写器，它从缓冲区读取数据，并将数据写入目标资源
```go
	// 接口签名
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

## 三、编码格式

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
