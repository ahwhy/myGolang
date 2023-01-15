# Golang-Network  Golang的网络

## 一、网络通信

- 网络请求流程图
```
	                                     内核空间                               用户空间
	        1).网络请求 ->        2).copy(I/O模型、DMA)->        3).copy(MMAP) ->       4-1).处理请求
	client                   网卡                         内核缓冲区            web服务进程    |
	        7).返回数据 <-               6).copy <-                5).copy <-          4-2).构建Respense
```

- 网络请求过程
	- I/O模型
		- 阻塞式网络I/O
		- 非阻塞式网络I/O
		- 多路复用网络I/O
	- DMA
		- 网卡和磁盘数据拷贝到内存流程比较固定，不涉及到运算操作，且非常耗时
		- 在磁盘嵌入一个DMA芯片，完成上述拷贝工作，把CPU解脱出来，让CPU专注于运算
	- MMAP
		- 用户空间和内核空间映射同一块内存空间，从而达到省略将数据从内核缓冲区拷贝到用户空间的操作，用户空间通过映射直接操作内核缓冲区的数据

- socket
	- 用户进程(复数，即应用层) -> socket抽象层 -> TCP/UDP(传输层)
	- socket把复杂的传输层协议封装成简单的接口，使应用层可以像读写文件一样进行网络数据的传输
	- socket通信过程
```
	                           设置监听端口|设置监听队列|阻塞，循环等待客户端连接
	Server端        Socket() -> Bind() -> Listen() -> Accept()  -> Receive() -> Send()    -> Close()
	                                                  建立连接       发送数据      接收数据      关闭连接
	Client端        Socket()           ->             Connect() -> Send()    -> Receive() -> Close()
```

## 二、Socket编程

### 1. 网络进程标识
- 用三元组(ip地址，协议，端口号)唯一标示网络中的一个进程，如(172.122.121.111, tcp, 5656)
- IPv4的地址位数为32位，分为4段，每段最大取值为255
- IPv6的地址位数为128位，分为8段，各段用16进制表示，最大取值为ffff
- 端口: `0~1023` 被熟知的应用程序占用(普通应用程序不可以使用)，`49152~65535` 客户端程序运行时动态选择使用

### 2. TCP/CS架构
- TCP首部
	- 报文结构
```
	|      源端口(2Bytes)          |         目的端口(2Bytes)       |
	|                             序号                             |
	|                            确认号                            |
	| 数据偏移 | 保留位 |           Tcp Flags            |   窗口    | 
	|  (4bit)  | (6bit) | (URG、ACK、PSH、RST、SYN、FIN) | (2Bytes) |
	|      检验和(2Bytes)          |         紧急指针(2Bytes)       |
	|                           Tcp选项                            |
```

- TCP协议
	- 传输层协议
	- `MSS = MTU - ip首部 - tcp首部`，MTU视网络接口层的不同而不同
	- TCP在建立连接时通常需要协商双方的MSS值
	- 应用层传输的数据大于MSS时需要分段
	- TCP首部
		- 报文结构(见下)
		- 一位=1bit，8bit=1Byte(字节)
		- 前20个字节是固定的，后面还有4N个可选字节(TCP选项)
		- 端口在 TCP层指定，ip在 IP层指定
			- 端口占2个字节，则最大端口号为2^16-1=65535
		- 序号／确认号
			- 由于应用层的数据被进行分段，为了在接收端对数据按顺序重组，需要为每段数据编个"序号"
			- 32 位序号: 也称为顺序号(Sequence Number)，简写为SEQ
			- 32 位确认序号: 也称为应答号(Acknowledgment Number)，简写为ACK
				- 在握手阶段，确认序号将发送方的序号加1作为回答
		- 数据偏移
			- 4 位数据偏移／首部长度: TCP数据部分距TCP开头的偏移量(一个偏移量是4个字节)，亦即TCP首部的长度
			- 由于首部可能含有可选项内容，因此TCP报头的长度是不确定的
			- 报头不包含任何任选字段则长度为20字节，4位首部长度字段所能表示的最大值为1111，转化为10进制为15，`15*32/8 = 60`，故报头最大长度为60字节，TCP选项最多有40个字节
			- 首部长度也叫数据偏移，是因为首部长度实际上指示了数据区在报文段中的起始偏移值
		- 保留位: 为将来定义新的用途保留，现在一般置0
		- 6 位标志字段
			- URG: 紧急指针标志，为1时表示紧急指针有效，为0则忽略紧急指针
			- ACK: 确认序号标志，为1时表示确认号有效，为0表示报文中不含确认信息，忽略确认号字段；TCP协议规定在连接建立后所有传输的报文段都必须把ACK设置为1
			- PSH: push标志，为1表示是带有push标志的数据，指示接收方在接收到该报文段以后，应尽快将这个报文段交给应用程序，而不是在缓冲区排队
			- RST: 重置连接标志，用于重置由于主机崩溃或其他原因而出现错误的连接；或者用于拒绝非法的报文段和拒绝连接请求
			- SYN: 同步序号，用于建立连接过程，在连接请求中，SYN=1和ACK=0表示该数据段没有使用捎带的确认域；而连接应答捎带一个确认，即SYN=1和ACK=1
			- FIN: finish标志，用于释放连接，为1时表示发送方已经没有数据发送了，即关闭本方数据流
		- 窗口
			- 滑动窗口大小，用来告知发送端接受端的缓存大小，以此控制发送端发送数据的速率，从而达到流量控制
			- 窗口数据长度为16bit，因而窗口大小最大为65535
			- 理论上最优大小为: `网络带宽(b/s)*RTT(s)`
		- 校验和
			- 奇偶校验，此校验和是对整个的 TCP 报文段，包括 TCP 头部和 TCP 数据，以 16 位字进行计算所得
			- 由发送端计算和存储，并由接收端进行验证
		- 紧急指针
			- 只有当 URG 标志置 1 时紧急指针才有效
			- 紧急指针是一个正的偏移量，和顺序号字段中的值相加表示紧急数据最后一个字节的序号
			- TCP 的紧急方式是发送端向另一端发送紧急数据的一种方式
		- 选项和填充
			- 通常为空，最常见的可选字段是最长报文大小，用于发送方与接收方协商最大报文段长度(MSS)，或在高速网络环境下作窗口调节因子时使用
			- 首部字段还定义了一个时间戳选项
			- 最常见的可选字段是最长报文大小，又称为MSS(Maximum Segment Size)
				- 每个连接方通常都在握手的第一步中指明这个选项，它指明本端所能接收的最大长度的报文段，选项长度不一定是32位的整数倍，所以要加填充位，即在这个字段中加入额外的零，以保证TCP头是32的整数倍
				- 1460是以太网默认的大小
		- 数据部分
			- TCP 报文段中的数据部分是可选的
			- 在一个连接建立和一个连接终止时，双方交换的报文段仅有 TCP 首部
			- 如果一方没有数据要发送，也使用没有任何数据的首部来确认收到的数据
			- 在处理超时的许多情况中，也会发送不带任何数据的报文段
	- TCP连接
		- 建立连接 三次握手
			- 第一次握手: TCP首部SYN=1，初始化一个序号=J；SYN报文段不能携带数据
			- 第二次握手: TCP首部SYN=1，ACK=1，确认号=J+1，初始化一个序号=K；此报文同样不携带数据
			- 第三次握手: ACK=1，序号=J+1，确认号=K+1；此次一般会携带真正需要传输的数据
			- 确认号: 即希望下次对方发过来的序号值
			- SYN Flood 攻击始终不进行第三次握手，属于DDOS攻击的一种
		- 断开连接 四次挥手
			- TCP的连接是全双工(可以同时发送和接收)的连接，因此在关闭连接的时候，必须关闭传送和接收两个方向上的连接
			- 第一次挥手: FIN=1，序号=M
			- 第二次挥手: ACK=1，序号=M+1
			- 第三次挥手: FIN=1，序号=N
			- 第四次挥手: ACK=1，序号=N+1
			- 从TIME_WAIT进入CLOSED需要经过2个MSL(Maxinum Segment Lifetime)，RFC793建议MSL=2分钟
	- TCP连接状态
		- 客户端独有的: SYN_SENT、CLOSING
		- 服务器独有的: LISTEN、SYN_RCVD、CLOSE_WAIT
		- 共有的: CLOSED、ESTABLISHED、FIN_WAIT1、FIN_WAIT2、LAST_ACK、TIME_WAIT
		- 描述
			- CLOSED: 起始点，tcp连接 超时或关闭时进入此状态
			- LISTEN: 服务端等待连接时的状态，调用Socket、bind、listen函数就能进入此状态，称为被动打开
			- SYN_SENT: 客户端发起连接，发送SYN给服务端，若不能连接进入CLOSED
			- SYN_RCVD: 服务端接受客户端的SYN，由LISTEN进入SYN_RCVD；同时返回一个ACK，并发送一个SYN给服务端
				- 特殊情况: 客户端发起SYN的同时接收到服务端的SYN，客户端会由SYN-sent转为SYN-rcvd状态
			- ESTABLISHED: 可以传输数据的状态
			- FIN_WAIT1: 主动关闭连接，发送FIN，由ESTABLISHED转为此状态
			- CLOSE_WAIT: 收到FIN，发送ACK，被动关闭的一方关闭连接进入此状态
			- LAST_ACK: 发送FIN，同时在接受ACK时，由CLOSE_WAIT进入此状态；被关闭的一方发起关闭请求
			- FIN_WAIT2: 主动关闭连接，接收到对方的FIN+ACK，由FIN_WAIT1转为此状态
			- CLOSED: 终止点，tcp连接 超时或关闭时进入此状态
			- CLOSING: 两边同时发送关闭请求，由FIN_WAIT1进入此状态，收到FIN请求，同时响应一个ACK
			- TIME_WAIT
				- [记一次time_wait & close_wait的讨论总结](https://developer.aliyun.com/article/745776)
				- FIN_WAIT2到此状态: 双方不同时发送FIN的情况下，主动关闭的一方在完成自己的请求后，收到对方的FIN后的状态
				- CLOSING到此状态: 双方同时发起关闭，都发送了FIN，同时接受FIN并发送ACK后的状态
				- FIN_WAIT2到此状态: 对方发来的FIN的ACK同时到达后的状态，与上一条的区别是本身发送的FIN回应的ACK先于对方的FIN到达，而上一条是FIN先到达
				- 具体形成
					- 主动关闭端A: 发FIN，进入FIN-WAIT-1状态，并等待......
					- 被动关闭端P: 收到FIN后必须立即发ACK，进入CLOSE_WAIT状态，并等待......
					- 主动关闭端A: 收到ACK后进入FIN-WAIT-2状态，并等待......
					- 被动关闭端P: 发FIN，进入LAST_ACK状态，并等待......
					- 主动关闭端A: 收到FIN后必须立即发ACK，进入TIME_WAIT状态，等待2MSL后结束Socket
					- 被动关闭端P: 收到ACK后结束Socket
				- 因此，TIME_WAIT状态是出现在主动发起连接关闭的一点，和是谁发起的连接无关，可以是client端，也可以是server端
				- 而从TIME_WAIT状态到CLOSED状态，有一个超时设置，这个超时设置是 `2*MSL (Maximum Segment Lifetime  RFC793定义了MSL为2分钟，Linux设置成了30s)`
				- TIME_WAIT的作用
					- 为了确保两端能完全关闭连接
						- 假设A服务器是主动关闭连接方，B服务器是被动方
						- 如果没有TIME_WAIT状态，A服务器发出最后一个ACK就进入关闭状态，如果这个ACK对端没有收到，对端就不能完成关闭
						- 对端没有收到ACK，会重发FIN，此时连接关闭，这个FIN也得不到ACK，而有TIME_WAIT，则会重发这个ACK，确保对端能正常关闭连接
					- 为了确保后续的连接不会收到"脏数据"
						- 刚才提到主动端进入TIME_WAIT后，等待2MSL后CLOSE，这里的MSL是指(maximum segment lifetime，内核一般是30s，2MSL就是1分钟)，网络上数据包最大的生命周期
						- 这是为了使网络上由于重传出现的old duplicate segment都消失后，才能创建参数(四元组，源IP/PORT，目标IP/PORT)相同的连接
						- 如果等待时间不够长，又创建好了一样的连接，再收到old duplicate segment，数据就错乱了
				- TIME_WAIT 会导致的问题
					- 新建连接失败
						- TIME_WAIT到CLOSED，需要2MSL=60s的时间
						- 这个时间非常长，每个连接在业务结束之后，需要60s的时间才能完全释放
						- 如果业务上采用的是短连接的方式，会导致非常多的TIME_WAIT状态的连接，会占用一些资源，主要是本地端口资源
					- 一台服务器的本地可用端口是有限的，也就几万个端口，由这个参数控制
						- `$ sysctl net.ipv4.ip_local_port_range # net.ipv4.ip_local_port_range = 32768 61000`
						- 当服务器存在非常多的TIME_WAIT连接，将本地端口都占用了，就不能主动发起新的连接去连其他服务器
					- 这里需要注意，是主动发起连接，又是主动发起关闭的一方才会遇到这个问题
					- 如果是server端主动关闭client端建立的连接产生了大量的TIME_WAIT连接，这是不会出现这个问题的
					- 除非是其中涉及到的某个客户端的TIME_WAIT连接都有好几万个
				- TIME_WAIT条目超出限制
					- 这个限制，是由一个内核参数控制的 `$ sysctl net.ipv4.tcp_max_tw_buckets # net.ipv4.tcp_max_tw_buckets = 5000`
						- 超出了这个限制会报一条INFO级别的内核日志，然后继续关闭掉连接
						- 并没有什么特别大的影响，只是增加了刚才提到的收到脏数据的风险而已
					- 另外的风险就是，关闭掉TIME_WAIT连接后，刚刚发出的ACK如果对端没有收到，重发FIN包出来时，不能正确回复ACK，只是回复一个RST包，导致对端程序报错，说connection reset
					- 因此net.ipv4.tcp_max_tw_buckets这个参数是建议不要改小的，改小会带来风险，没有什么收益，只是表面上通过netstat看到的TIME_WAIT少了些而已	
					- 并且，建议是当遇到条目不够，增加这个值，仅仅是浪费一点点内存而已
				- 如何处理time_wait
					- 最佳方案是应用改造长连接，但是一般不太适用
					- 修改系统回收参数
						- 设置以下参数
							- `net.ipv4.tcp_timestamps = 1`
							- `net.ipv4.tcp_tw_recycle = 1`
						- 但如果这两个参数同时开启，会校验源ip过来的包携带的timestamp是否递增，如果不是递增的话，则会导致三次握手建联不成功，具体表现为抓包的时候看到syn发出，server端不响应syn ack
						- 通俗一些来讲就是，一个局域网有多个客户端访问服务端，如果有客户端的时间比别的客户端时间慢，就会建连不成功
					- 治标不治本的方式
						- 放大端口范围 `$ sysctl net.ipv4.ip_local_port_range # net.ipv4.ip_local_port_range = 32768 61000`
						- 放大time_wait的buckets `$ sysctl net.ipv4.tcp_max_tw_buckets # net.ipv4.tcp_max_tw_buckets = 180000`
				- 特殊场景
					- 本机会发起大量短链接
						- nginx结合php-fpm需要本地起端口
						- nginx反代，如：java，容器等
					- 解决方案
						- tcp_tw_reuse参数需要结合`net.ipv4.tcp_timestamps = 1`一起来用即 服务器即做客户端，也做server端的时候
						- tcp_tw_reuse参数用来设置是否可以在新的连接中重用TIME_WAIT状态的套接字
							- 注意，重用的是TIME_WAIT套接字占用的端口号，而不是TIME_WAIT套接字的内存等
							- 这个参数对客户端有意义，在主动发起连接的时候会在调用的inet_hash_connect()中会检查是否可以重用TIME_WAIT状态的套接字
							- 如果在服务器段设置这个参数的话，则没有什么作用，因为服务器端ESTABLISHED状态的套接字和监听套接字的本地IP、端口号是相同的，没有重用的概念
							- 但并不是说服务器端就没有TIME_WAIT状态套接字。
						- 该类场景最终建议
							- `net.ipv4.tcp_tw_recycle = 0` 关掉快速回收
							- `net.ipv4.tcp_tw_reuse = 1`   开启tw状态的端口复用(客户端角色)
							- `net.ipv4.tcp_timestamps = 1` 复用需要timestamp校验为1 
							- `net.ipv4.tcp_max_tw_buckets = 30000` 放大bucket
							- `net.ipv4.ip_local_port_range = 15000 65000` 放大本地端口范围
					- 内存开销测试
```shell
	# ss -s
	# 15000个socket消耗30多m内存
	Total: 15254 (kernel 15288)
	TCP:   15169 (estab 5, closed 15158, orphaned 0, synrecv 0, timewait 3/0), ports 0
	Transport Total     IP        IPv6
	*         15288     -         -        
	RAW       0         0         0        
	UDP       5         4         1        
	TCP       11        11        0        
	INET      16        15        1        
	FRAG      0         0         0        
```

- Go语言中的TCP编程接口
	- net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket，大部分使用者只需要Dial、Listen和Accept函数提供的基本接口，以及相关的Conn和Listener接口
	- net.ResolveTCPAddr
		- `func ResolveTCPAddr(network, address string) (*TCPAddr, error)`
		- net参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP4(IPv4-only)，TCP6(IPv6-only)或者TCP(IPv4,、IPv6的任意一个)
		- addr表示域名或者IP地址，例如"www.qq.com:80" 或者"127.0.0.1:22"
	- net.ListenTCP
		- `func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)`
		- 监听端口
	- net.Accept
		- `func (l *TCPListener) Accept() (Conn, error)`
		- 等待客户端连接
	- net.DialTCP
		- `func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)`
		- net参数是"tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only)
		- laddr表示本机地址，一般设置为nil
		- 向raddr端发起请求，建立tcp连接
	- net.DialTimeout
		- `func DialTimeout(network, address string, timeout time.Duration) (Conn, error)`
		- 创建连接时设置超时时间
		- netwok指定为tcp，建立连接时指定超时
	- net.SetReadDeadline
		- `func (c *TCPConn) SetReadDeadline(t time.Time) error`
		- 设置从一个tcp连接上读取的超时时间
	- net.SetWriteDeadline
		- `func (c *TCPConn) SetWriteDeadline(t time.Time) error`
		- 设置从一个tcp连接上写入的超时时间
	- net.SetKeepAlive
		- `func (c *TCPConn) SetKeepAlive(keepalive bool) error`
		- 当一个tcp连接上没有数据时，操作系统会间隔性地发送心跳包，如果长时间没有收到心跳包会认为连接已经断开
	- net.Write
		- `func (c *TCPConn) Write(b []byte) (int, error)`
		- 通过conn发送数据
	- net.Read
		- `func (c *TCPConn) Read(b []byte) (int, error)`
		- 从conn里读取数据
	- ioutil.ReadAll
		- `func ReadAll(r io.Reader) ([]byte, error)`
		- 从conn中读取所有内容，直到遇到error(比如连接关闭)或EOF
	- net.Close
		- `func (c *TCPConn) Close() error`
		- `func (l *TCPListener) Close() error`
		- 关闭连接
```golang
	// TCPAddr 代表一个TCP终端地址
	type TCPAddr struct {
		IP   IP
		Port int
		Zone string // IPv6范围寻址域
	}
	// ResolveTCPAddr 将addr作为TCP地址解析并返回
	// 参数addr格式为"host:port"或"[ipv6-host%zone]:port"，解析得到网络名和端口名；net必须是"tcp"、"tcp4"或"tcp6"
	// IPv6地址字面值/名称必须用方括号包起来，如"[::1]:80"、"[ipv6-host]:http"或"[ipv6-host%zone]:80"
	func ResolveTCPAddr(net, addr string) (*TCPAddr, error)
	// 返回地址的网络类型，"tcp"
	func (a *TCPAddr) Network() string
	func (a *TCPAddr) String() string

	// TCPListener 代表一个TCP网络的监听者，应尽量使用Listener接口而不是假设(网络连接为)TCP
	type TCPListener struct { ... }
	// ListenTCP 在本地TCP地址laddr上声明并返回一个*TCPListener，net参数必须是"tcp"、"tcp4"、"tcp6"，如果laddr的端口字段为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口
	func ListenTCP(net string, laddr *TCPAddr) (*TCPListener, error)
	// Addr返回l监听的的网络地址，一个*TCPAddr
	func (l *TCPListener) Addr() Addr
	// SetDeadline 设置监听器执行的期限，t为Time零值则会关闭期限限制
	func (l *TCPListener) SetDeadline(t time.Time) error
	// Accept 用于实现Listener接口的Accept方法；他会等待下一个呼叫，并返回一个该呼叫的Conn接口
	func (l *TCPListener) Accept() (Conn, error)
	// AcceptTCP 接收下一个呼叫，并返回一个新的*TCPConn
	func (l *TCPListener) AcceptTCP() (*TCPConn, error)
	// Close 停止监听TCP地址，已经接收的连接不受影响
	func (l *TCPListener) Close() error
	// File 方法返回下层的os.File的副本，并将该副本设置为阻塞模式
	// 使用者有责任在用完后关闭f，关闭c不影响f，关闭f也不影响c
	func (l *TCPListener) File() (f *os.File, err error)

	// TCPConn 代表一个TCP网络连接，实现了Conn接口
	type TCPConn struct { ... }
	// DialTCP 在网络协议net上连接本地地址laddr和远端地址raddr
	// net必须是"tcp"、"tcp4"、"tcp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址
	func DialTCP(net string, laddr, raddr *TCPAddr) (*TCPConn, error)
	// LocalAddr 返回本地网络地址
	func (c *TCPConn) LocalAddr() Addr
	// RemoteAddr 返回远端网络地址
	func (c *TCPConn) RemoteAddr() Addr
	// SetReadBuffer 设置该连接的系统接收缓冲
	func (c *TCPConn) SetReadBuffer(bytes int) error
	// SetWriteBuffer 设置该连接的系统发送缓冲
	func (c *TCPConn) SetWriteBuffer(bytes int) error
	// SetDeadline 设置读写操作期限，实现了Conn接口的SetDeadline方法
	func (c *TCPConn) SetDeadline(t time.Time) error
	// SetReadDeadline 设置读操作期限，实现了Conn接口的SetReadDeadline方法
	func (c *TCPConn) SetReadDeadline(t time.Time) error
	// SetWriteDeadline 设置写操作期限，实现了Conn接口的SetWriteDeadline方法
	func (c *TCPConn) SetWriteDeadline(t time.Time) error
	// SetKeepAlive 设置操作系统是否应该在该连接中发送keepalive信息
	func (c *TCPConn) SetKeepAlive(keepalive bool) error
	// SetKeepAlivePeriod 设置keepalive的周期，超出会断开
	func (c *TCPConn) SetKeepAlivePeriod(d time.Duration) error
	// SetLinger 设定当连接中仍有数据等待发送或接受时的Close方法的行为
	// 如果sec < 0	(默认)，Close方法立即返回，操作系统停止后台数据发送；如果 sec == 0，Close立刻返回，操作系统丢弃任何未发送或未接收的数据；
	// 如果sec > 0，Close方法阻塞最多sec秒，等待数据发送或者接收，在一些操作系统中，在超时后，任何未发送的数据会被丢弃
	func (c *TCPConn) SetLinger(sec int) error
	// SetNoDelay 设定操作系统是否应该延迟数据包传递，以便发送更少的数据包(Nagle's算法)
	// 默认为真，即数据应该在Write方法后立刻发送
	func (c *TCPConn) SetNoDelay(noDelay bool) error
	// Read 实现了Conn接口Read方法
	func (c *TCPConn) Read(b []byte) (int, error)
	// Write 实现了Conn接口Write方法
	func (c *TCPConn) Write(b []byte) (int, error)
	// ReadFrom 实现了io.ReaderFrom接口的ReadFrom方法
	func (c *TCPConn) ReadFrom(r io.Reader) (int64, error)
	// Close 关闭连接
	func (c *TCPConn) Close() error
	// CloseRead 关闭TCP连接的读取侧(以后不能读取)，应尽量使用Close方法
	func (c *TCPConn) CloseRead() error
	// CloseWrite 关闭TCP连接的写入侧(以后不能写入)，应尽量使用Close方法
	func (c *TCPConn) CloseWrite() error
	// File 方法设置下层的os.File为阻塞模式并返回其副本
	// 使用者有责任在用完后关闭f，关闭c不影响f，关闭f也不影响c
	func (c *TCPConn) File() (f *os.File, err error)
```

### 3. UDP/CS架构
- UDP协议
	- 不需要建立连接，直接收发数据，效率很高
	- 面向报文
		- 对应用层交下来的报文，既不合并也不拆分，直接加上边界交给IP层
		- TCP是面向字节流
	- 从机制上不保证顺序(在IP层要对数据分段)，可能会丢包(检验和如果出差错就会把这个报文丢弃掉)
		- 在内网环境下分片乱序和数据丢包极少发生
	- 支持一对一、一对多、多对一和多对多的交互通信
	- 首部开销小，只占8个字节

- UDP首部
	- 报文结构
```
	| 源端口(2Bytes) | 目的端口(2Bytes) |
	|  UDP报文长度   |       检验和     |
	|           用户数据部分            |
```

- Go语言中的UDP编程接口
	- net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket，大部分使用者只需要Dial、Listen和Accept函数提供的基本接口，以及相关的Conn和Listener接口
	- net.ResolveUDPAddr
		- `func ResolveUDPAddr(network, address string) (*UDPAddr, error)`
		- netwok指定为"udp", "udp4" (IPv4-only), "udp6" (IPv6-only)，解析成udp地址
	- net.ListenUDP
		- `func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)`
		- 直接调用Listen返回一个udp连接
	- net.DialUDP
		- `func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)`
		- netwok指定为"udp", "udp4" (IPv4-only), "udp6" (IPv6-only)
		- 建立udp连接(伪连接)
	- net.DialTimeout
		- `func DialTimeout(network, address string, timeout time.Duration) (Conn, error)`
		- 创建连接时设置超时时间
		- netwok指定为udp，建立连接时指定超时
	- net.SetReadDeadline
		- `func (c *UDPConn) SetReadDeadline(t time.Time) error`
	- net.SetWriteDeadline
		- `func (c *UDPConn) SetWriteDeadline(t time.Time) error`
	- net.Read
		- `func (c *UDPConn) Read(b []byte) (int, error)`
	- net.ReadFromUDP
		- `func (c *UDPConn) ReadFromUDP(b []byte) (int, *UDPAddr, error)`
		- 读数据，会返回remote的地址
	- net.Write
		- `func (c *UDPConn) Write(b []byte) (int, error)`
	- net.WriteToUDP
		- `func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)`
		- 写数据，需要指定remote的地址
	- net.Close
		- `func (c *UDPConn) Close() error`
```golang
	// UDPAddr 代表一个UDP终端地址
	type UDPAddr struct {
		IP   IP
		Port int
		Zone string // IPv6范围寻址域
	}
	// ResolveUDPAddr 将addr作为UDP地址解析并返回
	// 参数addr格式为"host:port"或"[ipv6-host%zone]:port"，解析得到网络名和端口名；net必须是"udp"、"udp4"或"udp6"
	func ResolveUDPAddr(net, addr string) (*UDPAddr, error)
	// 返回地址的网络类型，"UDP"
	func (a *UDPAddr) Network() string
	func (a *UDPAddr) String() string

	// UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口
	type UDPConn struct { ... }
	// DialUDP 在网络协议net上连接本地地址laddr和远端地址raddr
	// net必须是"UDP"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址
	func DialUDP(net string, laddr, raddr *UDPAddr) (*UDPConn, error)
	// ListenUDP 创建一个接收目的地是本地地址laddr的UDP数据包的网络连接
	// net必须是"udp"、"udp4"、"udp6"；如果laddr端口为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口
	// 返回的*UDPConn的ReadFrom和WriteTo方法可以用来发送和接收UDP数据包(每个包都可获得来源地址或设置目标地址)
	func ListenUDP(net string, laddr *UDPAddr) (*UDPConn, error)
	// ListenMulticastUDP 接收目的地是ifi接口上的组地址gaddr的UDP数据包；它指定了使用的接口，如果ifi是nil，将使用默认接口
	func ListenMulticastUDP(net string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
	// LocalAddr 返回本地网络地址
	func (c *UDPConn) LocalAddr() Addr
	// RemoteAddr 返回远端网络地址
	func (c *UDPConn) RemoteAddr() Addr
	// SetReadBuffer 设置该连接的系统接收缓冲
	func (c *UDPConn) SetReadBuffer(bytes int) error
	// SetWriteBuffer 设置该连接的系统发送缓冲
	func (c *UDPConn) SetWriteBuffer(bytes int) error
	// SetDeadline 设置读写操作期限，实现了Conn接口的SetDeadline方法
	func (c *UDPConn) SetDeadline(t time.Time) error
	// SetReadDeadline 设置读操作期限，实现了Conn接口的SetReadDeadline方法
	func (c *UDPConn) SetReadDeadline(t time.Time) error
	// SetWriteDeadline 设置写操作期限，实现了Conn接口的SetWriteDeadline方法
	func (c *UDPConn) SetWriteDeadline(t time.Time) error
	// Read 实现了Conn接口Read方法
	func (c *UDPConn) Read(b []byte) (int, error)
	// ReadFrom 实现PacketConn接口ReadFrom方法
	func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error)
	// ReadFromUDP 从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址；ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误
	func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
	// ReadMsgUDP 从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误
	func (c *UDPConn) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error)
	// Write 实现了Conn接口Write方法
	func (c *UDPConn) Write(b []byte) (int, error)
	// WriteTo 实现PacketConn接口WriteTo方法
	func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error)
	// WriteToUDP 通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节
	// WriteToUDP方法会在超过一个固定的时间点之后超时，并返回一个错误；在面向数据包的连接上，写入超时是十分罕见的
	func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
	// WriteMsgUDP通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数(包数据、带外数据)和可能的错误
	func (c *UDPConn) WriteMsgUDP(b, oob []byte, addr *UDPAddr) (n, oobn int, err error)
	// Close 关闭连接
	func (c *UDPConn) Close() error
	// File 方法设置下层的os.File为阻塞模式并返回其副本
	// 使用者有责任在用完后关闭f，关闭c不影响f，关闭f也不影响c
	func (c *UDPConn) File() (f *os.File, err error)
```

### 4. TLS协议
- TLS的特性
	- 很多应用层协议(http、ftp、smtp等)直接使用明文传输
	- TLS(Transport Layer Security，安全传输层)将应用层的报文进行加密后再交由TCP进行传输
	- TLS 在SSL v3.0 的基础上，提供了一些增强功能，两者差别很小

- TLS的作用
	- 保密性，信息加密传输(对称加密)
	- 完整性，MAC检验(散列函数)
	- 认证，双方都可以配备证书，防止身份被冒充
	
- TLS过程
	- C端获得S端的证书
	- C端生成一个随机的AES口令，然后用S端的公钥通过RSA加密这个口令，并发给S端
	- S端用自己的RSA私钥解密得到AES口令
	- 双方使用这个共享的AES口令用AES加密通信

- TLS证书
	- TLS 通过两个证书来实现服务端身份验证，以及对称密钥的安全生成
		- CA 证书: 浏览器/操作系统自带，用于验证服务端的 TLS 证书的签名，保证服务端证书可信
		- TLS 证书: 客户端和服务端使用 TLS 证书进行协商，以安全地生成一个对称密钥
	- 证书，是非对称加密中的公钥，加上一些别的信息组成的一个文件
	- 通过权威CA机构验证证书主人的真实身份
		- 验证签名
			- 私钥加密，公钥解密
	- 证书来源
		- 向权威CA机构申请证书需要收费(也有短期免费的)
		- 若所有通信全部在自家网站或APP内完成可以使用本地签名证书(即自己生成CA证书)
			- 编写证书签名请求的配置文件csr.conf，指定加密算法、授信域名、申请者信息等，规范参考(openssl.org)
			- 生成server的TLS证书私钥 `openssl genrsa -out server.key 2048`
			- 根据第1步的配置文件，生成证书签名请求(公钥加申请者信息) `openssl req -new -key server.key -out server.csr -config csr.conf`
			- 生成CA的私钥 `openssl genrsa -out ca.key 2048`
			- 生成CA证书，有效期1000天 `openssl req -x509 -new -nodes -key ca.key -subj "/CN=MaGeCA" -days 1000 -out ca.crt`
			- 签名，得到server的TLS证书，有效其365天
				- 包含四部分: 公钥+申请者信息 + 颁发者(CA)的信息 + 签名(使用 CA 私钥加密)
					`openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key  -CAcreateserial -out server.crt -days 365  -extensions v3_ext -extfile csr.conf`
	- 服务端安装TLS证书
		- 编辑Nginx配置文件 
		- 使Nginx重新加载配置文件 `nginx -s reload`
	- 客户端安装TLS证书
		- 将服务端的TLS证书添加到 OS 的默认证书列表中
			- linux
				- `sudo cp server.crt /usr/local/share/ca-certificates/server.crt`
				- `sudo update-ca-certificates`
			- mac
				- 打开钥匙串，把刚刚生成的server.crt 拖到证书一栏里
		- 使用 HTTPS 客户端的 api 指定使用的 TLS 证书
```shell
	// 编辑Nginx配置文件 
	$ vim /etc/nginx/nginx.conf 
	# Settings for a TLS enabled server.
	server {
		listen       443;
		server_name  july.go.edu;                     // 填写自己的域名
		ssl_certificate "/path/to/server.crt";        // 填写刚生成的TLS证书(包含公钥)
		ssl_certificate_key "/path/to/server.key";    // 填写刚生成的TLS证书私钥
	}
```

### 5. WebSocket协议
- 与http的异同
	- 通信方式
		- http 客户端和服务端每次通信，都需要进行客户端请求(如 POST/GET)，服务端返回的流程
		- websocket 客户端和服务端再进行握手确认后，直接使用全双工的方式进行通信，最后断开连接
	- 相似和关联
		- 都是应用层协议，基于tcp传输协议
		- 跟http有良好的兼容性，websocket和http的默认端口都是80，websockets和https的默认端口都是443
		- websocket在握手阶段采用http发送数据
	- 差异
		- http是半双工，而websocket通过多路复用实现了全双工
		- http只能由client主动发起数据请求，而websocket还可以由server主动向client推送数据
			- 在需要及时刷新的场景中，http只能靠client高频地轮询，浪费严重
		- http是短连接(也可以实现长连接，HTTP1.1的连接默认使用长连接)，每次数据请求都得经过三次握手重新建立连接，而websocket是长连接
		- http长连接中每次请求都要带上header，而websocket在传输数据阶段不需要带header

- websocket握手协议
	- Upgrade:websocket 和 Connection:Upgrade 指明使用WebSocket协议
	- Sec-WebSocket-Version 指定Websocket协议版本
	- Sec-WebSocket-Key是一个Base64 encode的值，是浏览器随机生成的
	- 服务端收到Sec-WebSocket-Key后拼接上一个固定的GUID，进行一次SHA-1摘要，再转成Base64编码，得到Sec-WebSocket-Accept返回给客户端
		- 客户端对本地的Sec-WebSocket-Key执行同样的操作跟服务端返回的结果进行对比，如果不一致会返回错误关闭连接
		- 该操作把websocket header跟http header区分开

```
	// Request Header
	Sec-Websocket-Version:13
	Upgrade:websocket
	Connection:Upgrade
	Sec-Websocket-Key:duR0pUQxNgBJsRQKj2Jxsw==

	// Response Header
	Upgrade:websocket
	Connection:Upgrade
	Sec-Websocket-Accept:a1y2oy1zvgHsVyHMx+hZ1AYrEHI=
```

- WebSocket/CS架构
	- 引用第三方包
		`go get github.com/gorilla/websocket`
	- 服务端需要将 http协议升级到 WebSocket协议
		`func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)`	
	- 客户端发起握手，请求连接
		`func (*websocket.Dialer) Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)`

- Go语言中的标准库 "net/http"
```go
	// http.Serve 启动http服务
	func Serve(l net.Listener, handler Handler) error

	// http.ListenAndServe net.Listen和http.Serve两步合成一步
	func ListenAndServe(addr string, handler Handler) error

	// http.Handler 需要实现Handler接口
	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
	func (ws *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request)

	// http.HandleFunc
	func HandleFunc(pattern string, handler func(ResponseWriter, *Request))

	// http.ServeFile 命名文件或者目录的内容，回复请求
	func ServeFile(w ResponseWriter, r *Request, name string)
```

- Go语言中的第三方库 `github.com/gorilla/websocket`
	- websocket发送的消息类型有5种: TextMessag, BinaryMessage, CloseMessag, PingMessage, PongMessage
	- TextMessag 和 BinaryMessage 分别表示发送文本消息和二进制消息
	- CloseMessage 关闭帧，接收方收到这个消息就关闭连接
	- PingMessage和PongMessage是保持心跳的帧
		- 发送方 -> 接收方是PingMessage
		- 接收方 -> 发送方是PongMessage
		- 目前浏览器没有相关api发送ping给服务器，只能由服务器发ping给浏览器，浏览器返回pong消息
```go
	// Upgrader结构体
	type Upgrader struct {
		HandshakeTimeout time.Duration                                                 // websocket握手超时时间
		ReadBufferSize, WriteBufferSize int                                            // io操作的缓存大小
		Error func(w http.ResponseWriter, r *http.Request, status int, reason error)   // http错误响应函数
		CheckOrigin func(r *http.Request) bool                                         // 用于统一的链接检查，以防止跨站点请求伪造
	}
```

- 示例 [聊天室实现](https://github.com/ahwhy/myGolang/tree/main/network/chatroom/doc)

## 三、Http编程

### 1. Http协议
- Http 超文本传输协议 Hyper Text Transfer Protocol
	- Http属于应用层协议，它在传输层用的是tcp协议
	- 无状态，对事务处理没有记忆能力(对比TCP协议里的确认号)
		- 如果要保存状态需要引用其他技术，如cookie
	- 无连接，每次连接只处理一个请求
		- 早期带宽和计算资源有限，这么做是为了追求传输速度快，后来通过`Connection: Keep-Alive`实现长连接
		- http1.1废弃了Keep-Alive，默认支持长连接

### 2. Http-Request
- Request报文结构
```
	| 请求方法 | 空格 | URL | 空格 | 协议版本 | \r | \n |      // 请求行
	|      字段名     |  :  |        值       | \r | \n |      // 请求头
	                         ...
	|      字段名     |  :  |        值       | \r | \n |
	| \r | \n |                                                // 空行
	|                        正文                       |
```

- Http Request
	- 请求方法
		- http 1.0
			- GET      请求获取Request-URI所标识的资源 
			- POST     向URI提交数据(例如提交表单或上传数据)
			- HEAD     类似于GET，返回的响应中没有具体的内容，用于获取报头
		- http 1.1
			- PUT      对服务器上已存在的资源进行更新
			- DELETE   请求服务器删除指定的页面
			- CONNECT  HTTP/1.1预留，能够将连接改为管道方式的代理服务器
			- OPTIONS  查看服务端性能
			- TRACE    回显服务器收到的请求，主要用于测试或诊断
			- PATCH    同PUT，可只对资源的一部分更新，资源不存在时会创建
		- 注意
			- 实际中server对各种request method的处理方式可能不是按协议标准来的
				- server收到PUT请求时执行DELETE操作
				- 仅用一个GET方法也能实现增删改查的全部功能
			- 大多数浏览器只支持GET和POST
	- URL
		- URI: uniform resource identifier，统一资源标识符，用来标识唯一的一个资源
		- URL: uniform resource locator，统一资源定位器，它是一种具体的URI，指明了如何locate这个资源
		- 示例见下
	- 协议版本
		- HTTP/1.0
		- HTTP/1.1
	- 请求头
		- Header字段解释及示例
			- Accept           指定客户端能够接收的内容类型，如 Accept: text/plain, text/html
			- Accept-Charset   浏览器可以接受的字符编码集，如 Accept-Charset: iso-8859-5
			- Accept-Encoding  指定浏览器可以支持的web服务器返回内容压缩编码类型，如 Accept-Encoding: compress, gzip
			- Accept-Language  浏览器可接受的语言，如 Accept-Language: en,zh
			- Authorization    HTTP授权的授权证书，如 Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
			- Cache-Control    指定请求和响应遵循的缓存机制，如 Cache-Control: no-cache
			- Connection       表示是否需要持久连接(HTTP 1.1默认进行持久连接)，如 Connection: close
			- Cookie           HTTP请求发送时，会把保存在该请求域名下的所有cookie值一起发送给web服务器，如 Cookie: $Version=1; Skin=new;
			- Content-Length   请求的内容长度，如 Content-Length: 348
			- Content-Type     指定正文(body)的数据格式，如 Content-Type: application/x-www-form-urlencoded
			- Date             请求发送日期和时间
			- User-Agent       浏览器信息，如 Mozilla/5.0 (Windows NT 6.1; Win64; x64)
		- Content-Type
			- application/x-www-form-urlencoded
				- 浏览器的原生form表单，如果不设置 Content-Type 属性，则默认以 application/x-www-form-urlencoded 方式传输数据
				- 正文示例: name=manu&message=this_is_great
			- multipart/form-data
				- 上传文件时使用multipart/form-data，支持多种文件格式
				- 正文示例: name="text"name="file"; filename="chrome.png"Content-Type: image/png... content of chrome.png
			- application/json
				- JSON数据格式
				- 正文示例: `{"title":"test","sub":[1,2,3]}`
			- application/xhtml+xml
				- XHTML格式
			- application/xml
				- XML数据格式
			- application/pdf
				- pdf格式
			- application/msword
				- word文档格式
			- application/octet-stream
				- 二进制流数据
			- text/xml
				- 示例见下
			- text/html
				- HTML格式
			- text/plain
				- 纯文本格式
			- image/gif
				- gif图片格式
			- image/jpeg
				- jpg图片格式
			- image/png
				- png图片格式
	- 请求正文
		- GET请求没有请求正文
		- POST可以包含GET，示例见下
		- GET和POST的区别
			- get的请求参数全部在url里，参数变时url就变
				- post可以把参数放到请求正文里，参数变时url不变
			- http协议并没有对url和请求正文做长度限制，但在实际中浏览器对url的长度限制到比请求正文要小很多
				- 所以post可以提交的数据比get要大得多
			- get比post更容易受到攻击(源于get的参数直接暴露在url里)
```xml
	// url
	https://baijiahao.baidu.com/s?id=1603848351636567407&wfr=spider&for=pc
	协议    域名                  参数
	http://www.qq.com:8080/news/tech/43253.html?id=432&name=f43s#pic
	协议   域名       端口 路径      文件名     参数             锚点
	
	// text/xml
	<?xml version="1.0"?>
	<methodCall>
		<methodName>
			examples.getStateName
		</methodName>
	</methodCall>`
	
	// POST请求
	POST /post?id=1234&page=1 HTTP/1.1
	Content-Type: application/x-www-form-urlencoded
	
	name=manu&message=this_is_great
```

### 3. Http-Response
- Response报文结构
```
	| 协议版本 | 状态码 | 原因话术 | \r | \n |      // 相应行
	|  字段名  |    :   |    值    | \r | \n |      // 相应头
					    ...
	|  字段名  |    :   |    值    | \r | \n |
	| \r | \n |                                     // 空行
	|                响应正文                |
```

- Http Response
	- 状态码及话术
		- Code\Phrase
			- 200 Ok                      请求成功
			- 400 Bad Request             客户端有语法错误，服务端不理解
			- 401 Unauthorized            请求未经授权
			- 403 Forbidden               服务端拒绝提供服务
			- 404 Not Found               请求资源不存在
			- 500 Internal Server Error   服务器发生不可预期的错误
			- 503 Server Unavailable      服务器当前有问题，过段时间可能恢复
	- 响应头
		- Header字段解释及示例
			- Allow              对某网络资源的有效的请求行为，如 Allow: GET, HEAD
			- Date               原始服务器消息发出的时间，如 Date: Tue, 15 Nov 2010 08:12:31 GMT
			- Content-Encoding   服务器支持的返回内容压缩编码类型，如 Content-Encoding: gzip
			- Content-Language   响应体的语言，如 Content-Language: en,zh
			- Content-Length     响应体的长度，如 Content-Length: 348
			- Cache-Control      指定请求和响应遵循的缓存机制，如 Cache-Control: no-cache
			- Content-Type       返回内容的MIME类型，如 Content-Type: text/html; charset=utf-8
	- 响应正文
		- Http Response示例
```html
	HTTP/1.1 200 OK 
	Date: Fri, 22 May 2009 06:07:21 GMT 
	Content-Type: text/html; charset=UTF-8 
			
	<html>                                  // 响应正文
		<head></head>
		<body>
			<!--body goes here--> 
		</body> 
	</html>
```

### 4. Https
- Https
	- http:  (应用层)HTTP -> TCP -> IP
	- https: (应用层)HTTP -> SSL/TCP -> TCP -> IP
		- HTTP + 加密 + 认证 + 完整性保护 = HTTPS(HTTP Secure)

### 5. Go语言中的HTTP编程接口
- net/http
	- http包提供了HTTP客户端和服务端的实现
```golang
	// Get、Head、Post和PostForm函数发出HTTP/HTTPS请求
	resp, err := http.Get("http://example.com/")
	// ...
	resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	// ...
	resp, err := http.PostForm("http://example.com/form",
		url.Values{"key": {"Value"}, "id": {"123"}})
	// ...

	// 程序在使用完回复后必须关闭回复的主体
	resp, err := http.Get("http://example.com/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...

	// 要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	resp, err := client.Get("http://example.com")
	// ...
	req, err := http.NewRequest("GET", "http://example.com", nil)
	// ...
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)
	// ...

	// 要管理代理、TLS配置、keep-alive、压缩和其他设置，创建一个Transport
	// Client和Transport类型都可以安全的被多个go程同时使用；出于效率考虑，应该一次建立、尽量重用
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://example.com")
	// ...

	// ListenAndServe使用指定的监听地址和处理器启动一个HTTP服务端
	// 处理器参数通常是nil，这表示采用包变量DefaultServeMux作为处理器
	// Handle和HandleFunc函数可以向DefaultServeMux添加处理器
	http.Handle("/foo", fooHandler)
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

	// 要管理服务端的行为，可以创建一个自定义的Serve
	s := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
```

- server/client
```golang
	// http-server
	// 把返回的内容写入http.ResponseWriter
	func HelloHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Boy") 
	}
	
	// 路由，请求要目录时去执行HelloHandler  
	http.HandleFunc("/", HelloHandler)  
	// ListenAndServe如果不发生error会一直阻塞
	// 为每一个请求创建一个协程去处理
	http.ListenAndServe(":5656", nil)

	// http-client
	if resp, err := http.Get("http://127.0.0.1:5656"); err != nil {
		panic(err)
	} else {
		// 注意一定要调用resp.Body.Close()，否则会协程泄漏(同时引发内存泄漏)
		defer resp.Body.Close()
		// 把resp.Body输出到标准输出流
		io.Copy(os.Stdout, resp.Body) 
	}
```

### 6. http-router
- Go语言中的第三方库 `github.com/julienschmidt/httprouter`
	- `go get -u github.com/julienschmidt/httprouter`
	- Router实现了`http.Handler`接口
	- 为各种request method提供了便捷的路由方式
	- 支持restful请求方式
	- 支持ServeFiles访问静态文件
	- 可以自定义捕获panic的方法

### 7. 关于请求校验的常见问题
- XSS
	- 跨站脚本攻击(Cross-site scripting，XSS)是一种安全漏洞，即通过注入脚本获取敏感信息
		- 攻击者可以利用这种漏洞在网站上注入恶意的客户端代码
		- 当被攻击者登陆网站时就会自动运行这些恶意代码，从而攻击者可以突破网站的访问权限，冒充受害者

- CSRF
	- 跨站请求伪造(Cross-site request forgery，CSRF)是一种冒充受信任用户，向服务器发送非预期请求的攻击方式
	- 例如，这些非预期请求可能是通过在跳转链接后的 URL 中加入恶意参数来完成

- jsonp 跨域
	- 主流浏览器不允许跨域访问数据(端口不同也属于跨域)
	- `<script>`标签的src属性不受同源策略限制
	- 通过script的src请求返回的数据，浏览器会当成js脚本去处理
		- 所以服务端可以返回一个在客户端存在的js函数
	- [跨域资源共享 CORS 详解](https://www.ruanyifeng.com/blog/2016/04/cors.html)

- validator
	- Go语言中的第三方库 `github.com/go-playground/validator`
		- `go get github.com/go-playground/validator`
	- 范围约束
		- 对于字符串、切片、数组和map，约束其长度: len=10, min=6, max=10, gt=10
		- 对于数值，约束其取值: min, max, eq, ne, gt, gte, lt, lte, oneof=6 8
	- 跨字段约束
		- 跨字段就在范围约束的基础上加field后缀
		- 如果还跨结构体(cross struct)就在跨字段的基础上在field前面加cs
			- 范围约束 cs, field
	- 字符串约束
		- contains包含子串
		- containsany包含任意unicode字符，containsany=abcd
		- containsrune包含rune字符，containsrune=☻
		- excludes不包含子串
		- excludesall不包含任意的unicode字符，excludesall=abcd
		- excludesrune不包含rune字符，excludesrune=☻
		- startswith以子串为前缀
		- endswith以子串为后缀
	- 唯一性uniq
		- 对于数组和切片，约束没有重复的元素
		- 对于map，约束没的重复的value
		- 对于元素类型为结构体的切片，unique约束结构体对象的某个字段不重复，通过unqiue=field指定这个字段名
		- ``` Friends []User `validate:"unique=Name"` ```
	- 自定义约束

```go
	// 范围约束
	type RegistRequest struct {
		UserName string `validate:"gt=0"`                // >0 长度大于0
		PassWord string `validate:"min=6,max=12"`        //密码长度[6, 12]
		PassRepeat string `validate:"eqfield=PassWord"`  //跨字段相等校验
		Email string `validate:"email"`                  //需要满足email的格式
	}

	// 自定义约束
	func validateEmail(fl validator.FieldLevel) bool {
		input := fl.Field().String()
		if pass, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,})\.([a-z]{2,4})$`, input); pass {
			return true
		}
		return false
	}
	
	//注册一个自定义的validator
	val.RegisterValidation("my_email", validateEmail)
	Email string `validate:"my_email"`
```

### 8. http中间件
- 中间件的作用
	- 将业务代码和非业务代码解耦
	- 非业务代码: 限流、超时控制、打日志等等

- 中间件的实现原理
	- 传入一个http.Handler，外面套上一些非业务功能代码，再返回一个http.Handler
	- 支持中间件层层嵌套
	- 通过HandlerFunc把一个`func(rw http.ResponseWriter, r *http.Request)`函数转为Handler
```go
	func timeMiddleWare(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			begin := time.Now()
			next.ServeHTTP(rw, r)
			timeElapsed := time.Since(begin)
			log.Printf("request %s use %d ms\n", r.URL.Path, timeElapsed.Milliseconds())
		})
	}
```

### 9. Go语言中的Http框架
- 自定义Web框架
	- 框架的作用
		- 节省封装的开发时间，统一各团队的编码风格，节省沟通和排查问题的时间
	- Web框架需要具备的功能
		- request参数获取
		- 参数校验 validator
		- 路由 httprouter
		- response生成和渲染
		- 中间件
		- 会话管理
	- Gorilla工具集
		- mux       一款强大的HTTP路由和URL匹配器
		- websocket 一个快速且被广泛应用的WebSocket实现
		- sessions  支持将会话跟踪信息保存到Cookie或文件系统
		- handler   为http服务提供很多有用的中间件
		- schema    表单数据和go struct互相转换
		- csrf      提供防止跨站点请求攻击的中间件

- Gin框架
	- Gin是一款高性能的、简单轻巧的Http Web框架
		`go get -u github.com/gin-gonic/gin`
	- 路由
		- Gin的路由是基于httprouter做的
		- 支持GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD
		- 支持路由分组，不用重复写上级路径
	- 参数获取
		- `c.Query()`         从GET请求的URL中获取参数
		- `c.Param()`         从Restful风格的url中获取参数
		- `c.PostForm()`      从post表单中获取参数
		- `c.FormFile()`      获取上传的文件，消息类型为form-data
		- `c.MultipartForm()` multipart/form-data可以上传多个form-data 并且用分隔符进行分割
	- 参数绑定
	- Gin生成response
		- `c.String()`   response Content-Type= text/plain
		- `c.JSON()`     response Content-Type= application/json
		- `c.XML()`      response Content-Type= application/xml
		- `c.HTML()`     前端写好模板，后端往里面填值
		- `c.Redirect()` 重定向
	- Gin参数校验
		- 基于`go-playground/validator`
		- Gin中间件
			- 丰富的第三方中间件 `gin-gonic/contrib (github.com)`
		- Gin会话
			- http是无状态的，即服务端不知道两次请求是否来自于同一个客户端
			- Cookie由服务端生成，发送给客户端，客户端保存在本地
			- 客户端每次发起请求时把Cookie带上，以证明自己的身份
			- HTTP请求中的Cookie头只会包含name和value信息(服务端只能取到name和value)，domain、path、expires等cookie属性是由浏览器使用的，对服务器来说没有意义
			- Cookie可以被浏览器禁用
			
```go
	// Gin参数绑定
	type Student struct {
		Name string `form:"username" json:"name" uri:"user" xml:"user" yaml:"user" binding:"required"`
		Addr string `form:"addr" json:"addr" uri:"addr" xml:"addr" yaml:"addr" binding:"required"`
	}
	var stu Student
	ctx.ShouldBindJSON(&stu)

	// Gin参数校验
	type Student struct {
		Name string `form:"name" binding:"required"` // required:必须上传name参数
		Score int `form:"score" binding:"gt=0"` // score必须为正数
		Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"` // 自定义验证before_today，日期格式东8区
		Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"` // 毕业时间要晚于入学时间
	}
	
	// 全局MiddleWare
	engine.Use(timeMiddleWare()) 
	// 局部MiddleWare
	engine.GET("/girl", limitMiddleWare(), func(c *gin.Context) {
		c.String(http.StatusOK, "hi girl")
	})
```

- Beego框架
	- beego简介
		- beego是一个大而全的http框架，用于快速开发go应用程序
		- bee工具提供诸多命令，帮助进行 beego 项目的创建、热编译、开发、测试、和部署
		- beego的八大模块互相独立，高度解耦，开发者可任意选取
		- 日志模块
		- ORM模块
		- Context模块  封装了request和response
		- Cache模块    封装了memcache、redis、ssdb
		- Config模块   解析.ini、.yaml、.xml、.json、.env等配置文件
		- httplib模块
		- Session模块  session保存在服务端，用于标识客户身份，跟踪会话
		- toolbox模块  健康检查、性能调试、访问统计、计划任务
	- 用bee工具创建web项目
		- MainController
		- MVC
			- View 前端页面
			- Controller 处理业务逻辑
			- Model 把Controller层重复的代码抽象出来
				- 在Model层可以使用beego提供的ORM功能

```golang
	// 使用方法
	$ go get github.com/astaxie/beego
	$ go get github.com/beego/bee
	$ cd $GOPATH/src
	$ bee new myweb
	$ cd myweb
	$ go build -mod=mod
	$ bee run
	
	// MainController
	type MainController struct { //继承自beego.Controller
		beego.Controller //beego.Controller里有Get、Post、Put等方法
	}
	func (c *MainController) Get() {
		c.Data["Website"] = "github.com/Orisun"
		c.Data["Email"] = zhchya@gmail.com
		//TplName是需要渲染的模板  .tpl经常被用来表示PHP模板
		c.TplName = "index.tpl"
		//Resquest和ResponseWriter都在beego.Controller.Ctx里
		fmt.Println("remote addr", c.Ctx.Request.RemoteAddr)
		//如果指定了response正文，就不会去渲染index.tpl了
		c.Ctx.WriteString("Hi boy") 
	}
```

## 四、Golang的标准库 net包

### 1. net
- net
	- net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket
	- 虽然本包提供了对网络原语的访问，大部分使用者只需要Dial、Listen和Accept函数提供的基本接口；以及相关的Conn和Listener接口
	- crypto/tls包提供了相同的接口和类似的Dial和Listen函数

- net.Const
```golang
	// Const
	const (
		IPv4len = 4
		IPv6len = 16
	)

	// Variables
	var (
		IPv4bcast     = IPv4(255, 255, 255, 255) // 广播地址
		IPv4allsys    = IPv4(224, 0, 0, 1)       // 所有主机和路由器
		IPv4allrouter = IPv4(224, 0, 0, 2)       // 所有路由器
		IPv4zero      = IPv4(0, 0, 0, 0)         // 本地地址，只能作为源地址(曾用作广播地址)
	)
	// 常用的IPv4地址
	var (
		IPv6zero                   = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		IPv6unspecified            = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		IPv6loopback               = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		IPv6interfacelocalallnodes = IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
		IPv6linklocalallnodes      = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
		IPv6linklocalallrouters    = IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}
	)
	// 常用的IPv6地址
	var (
		ErrWriteToConnected = errors.New("use of WriteTo with pre-connected connection")
	)
```

- net.Error
```golang
	// ParseError 代表一个格式错误的字符串，Type为期望的格式
	type ParseError struct {
		Type string
		Text string
	}
	func (e *ParseError) Error() string

	// Error 代表一个网络错误
	type Error interface {
		error
		Timeout() bool   // 错误是否为超时？
		Temporary() bool // 错误是否是临时的？
	}
	// UnknownNetworkError
	type UnknownNetworkError string
	func (e UnknownNetworkError) Error() string
	func (e UnknownNetworkError) Temporary() bool
	func (e UnknownNetworkError) Timeout() bool
	// InvalidAddrError
	type InvalidAddrError string
	func (e InvalidAddrError) Error() string
	func (e InvalidAddrError) Temporary() bool
	func (e InvalidAddrError) Timeout() bool
	// DNSConfigError 代表读取主机DNS配置时出现的错误
	type DNSConfigError struct {
		Err error
	}
	func (e *DNSConfigError) Error() string
	func (e *DNSConfigError) Temporary() bool
	func (e *DNSConfigError) Timeout() bool
	// DNSError 代表DNS查询的错误
	type DNSError struct {
		Err       string // 错误的描述
		Name      string // 查询的名称
		Server    string // 使用的服务器
		IsTimeout bool
	}
	func (e *DNSError) Error() string
	func (e *DNSError) Temporary() bool
	func (e *DNSError) Timeout() bool
	// AddrError
	type AddrError struct {
		Err  string
		Addr string
	}
	func (e *AddrError) Error() string
	func (e *AddrError) Temporary() bool
	func (e *AddrError) Timeout() bool
	// OpError是经常被net包的函数返回的错误类型；它描述了该错误的操作、网络类型和网络地址
	type OpError struct {
		// Op是出现错误的操作，如"read"或"write"
		Op  string
		// Net是错误所在的网络类型，如"tcp"或"udp6"
		Net string
		// Addr是出现错误的网络地址
		Addr Addr
		// Err是操作中出现的错误
		Err error
	}
	func (e *OpError) Error() string
	func (e *OpError) Temporary() bool
	func (e *OpError) Timeout() bool
```

- net.flag
```golang
	const (
		FlagUp           Flags = 1 << iota // 接口在活动状态
		FlagBroadcast                      // 接口支持广播
		FlagLoopback                       // 接口是环回的
		FlagPointToPoint                   // 接口是点对点的
		FlagMulticast                      // 接口支持组播
	)

	type Flags uint
	func (f Flags) String() string
```

- net.Interface
```golang
	// Interface类型代表一个网络接口(系统与网络的一个接点)，包含接口索引到名字的映射，也包含接口的设备信息
	type Interface struct {
		Index        int          // 索引，>=1的整数
		MTU          int          // 最大传输单元
		Name         string       // 接口名，例如"en0"、"lo0"、"eth0.100"
		HardwareAddr HardwareAddr // 硬件地址，IEEE MAC-48、EUI-48或EUI-64格式
		Flags        Flags        // 接口的属性，例如FlagUp、FlagLoopback、FlagMulticast
	}
	// InterfaceByIndex 返回指定索引的网络接口
	func InterfaceByIndex(index int) (*Interface, error)
	// InterfaceByName 返回指定名字的网络接口
	func InterfaceByName(name string) (*Interface, error)
	// Addrs 返回网络接口ifi的一或多个接口地址
	func (ifi *Interface) Addrs() ([]Addr, error)
	// MulticastAddrs 返回网络接口ifi加入的多播组地址
	func (ifi *Interface) MulticastAddrs() ([]Addr, error)

	// Interfaces 返回该系统的网络接口列表
	func Interfaces() ([]Interface, error)
	// InterfaceAddrs返回该系统的网络接口的地址列表
	func InterfaceAddrs() ([]Addr, error)
```

- net.IP
```golang
	// IP类型是代表单个IP地址的[]byte切片
	// 本包的函数都可以接受4字节(IPv4)和16字节(IPv6)的切片作为输入
	type IP []byte

	// IPv4 返回包含一个IPv4地址a.b.c.d的IP地址(16字节格式)
	func IPv4(a, b, c, d byte) IP

	// ParseIP 将s解析为IP地址，并返回该地址
	// 如果s不是合法的IP地址文本表示，ParseIP会返回nil
	// 字符串可以是小数点分隔的IPv4格式(如"74.125.19.99")或IPv6格式(如"2001:4860:0:2001::68")格式
	func ParseIP(s string) IP
	// 如果ip是全局单播地址，则返回真
	func (ip IP) IsGlobalUnicast() bool
	// 如果ip是链路本地单播地址，则返回真
	func (ip IP) IsLinkLocalUnicast() bool
	// 如果ip是接口本地组播地址，则返回真
	func (ip IP) IsInterfaceLocalMulticast() bool
	// 如果ip是链路本地组播地址，则返回真
	func (ip IP) IsLinkLocalMulticast() bool
	// 如果ip是组播地址，则返回真
	func (ip IP) IsMulticast() bool
	// 如果ip是环回地址，则返回真
	func (ip IP) IsLoopback() bool
	// 如果ip是未指定地址，则返回真
	func (ip IP) IsUnspecified() bool
	// 函数返回IP地址ip的默认子网掩码。只有IPv4有默认子网掩码；如果ip不是合法的IPv4地址，会返回nil
	func (ip IP) DefaultMask() IPMask
	// 如果ip和x代表同一个IP地址，Equal会返回真；代表同一地址的IPv4地址和IPv6地址也被认为是相等的
	func (ip IP) Equal(x IP) bool
	// To16将 一个IP地址转换为16字节表示。如果ip不是一个IP地址(长度错误)，To16会返回nil
	func (ip IP) To16() IP
	// To4 将一个IPv4地址转换为4字节表示。如果ip不是IPv4地址，To4会返回nil
	func (ip IP) To4() IP
	// Mask 方法认为mask为ip的子网掩码，返回ip的网络地址部分的ip(主机地址部分都置0)
	func (ip IP) Mask(mask IPMask) IP
	// String 返回IP地址ip的字符串表示
	// 如果ip是IPv4地址，返回值的格式为点分隔的，如"74.125.19.99"；否则表示为IPv6格式，如"2001:4860:0:2001::68"
	func (ip IP) String() string
	// MarshalText 实现了encoding.TextMarshaler接口，返回值和String方法一样
	func (ip IP) MarshalText() ([]byte, error)
	// UnmarshalText 实现了encoding.TextUnmarshaler接口；IP地址字符串应该是ParseIP函数可以接受的格式
	func (ip *IP) UnmarshalText(text []byte) error

	// IPMask 代表一个IP地址的掩码
	type IPMask []byte
	// IPv4Mask 返回一个4字节格式的IPv4掩码a.b.c.d
	func IPv4Mask(a, b, c, d byte) IPMask
	// CIDRMask 返回一个IPMask类型值，该返回值总共有bits个字位，其中前ones个字位都是1，其余字位都是0
	func CIDRMask(ones, bits int) IPMask
	// Size 返回m的前导的1字位数和总字位数；如果m不是规范的子网掩码(字位：/^1+0+$/)，将返会(0, 0)
	func (m IPMask) Size() (ones, bits int)
	// String 返回m的十六进制格式，没有标点
	func (m IPMask) String() string

	// IPNet 表示一个IP网络
	type IPNet struct {
		IP   IP     // 网络地址
		Mask IPMask // 子网掩码
	}
	// ParseCIDR将s作为一个CIDR(无类型域间路由)的IP地址和掩码字符串，如"192.168.100.1/24"或"2001:DB8::/48"，解析并返回IP地址和IP网络
	func ParseCIDR(s string) (IP, *IPNet, error)
	// Contains 报告该网络是否包含地址ip
	func (n *IPNet) Contains(ip IP) bool
	// Network 返回网络类型名："ip+net"，注意该类型名是不合法的
	func (n *IPNet) Network() string
	// String 返回n的CIDR表示，如"192.168.100.1/24"或"2001:DB8::/48"
	func (n *IPNet) String() string
```

- net.Addr
```golang
	// SplitHostPort 将格式为"host:port"、"[host]:port"或"[ipv6-host%zone]:port"的网络地址分割为host或ipv6-host%zone和port两个部分
	func SplitHostPort(hostport string) (host, port string, err error)
	// JoinHostPort 将host和port合并为一个网络地址。一般格式为"host:port"；如果host含有冒号或百分号，格式为"[host]:port"
	func JoinHostPort(host, port string) string

	// HardwareAddr 类型代表一个硬件地址(MAC地址)
	type HardwareAddr []byte
	// ParseMAC 解析一个IEEE 802 MAC-48、EUI-48或EUI-64硬件地址
	func ParseMAC(s string) (hw HardwareAddr, err error)
	func (a HardwareAddr) String() string

	// Addr 代表一个网络终端地址
	type Addr interface {
		Network() string // 网络名
		String() string  // 字符串格式的地址
	}

	// IPAddr 代表一个IP终端的地址
	type IPAddr struct {
		IP   IP
		Zone string // IPv6范围寻址域
	}
	// ResolveIPAddr 将addr作为一个格式为"host"或"ipv6-host%zone"的IP地址来解析
	// 函数会在参数net指定的网络类型上解析，net必须是"ip"、"ip4"或"ip6"
	func ResolveIPAddr(net, addr string) (*IPAddr, error)
	// Network 返回地址的网络类型："ip"
	func (a *IPAddr) Network() string
	func (a *IPAddr) String() string

	// TCPAddr 见上Tcp

	// UDPAddr 见上Udp

	// UnixAddr 代表一个Unix域socket终端地址
	type UnixAddr struct {
		Name string
		Net  string
	}	
	// ResolveUnixAddr 将addr作为Unix域socket地址解析，参数net指定网络类型："unix"、"unixgram"或"unixpacket"
	func ResolveUnixAddr(net, addr string) (*UnixAddr, error)
	// Network 返回地址的网络类型，"unix"，"unixgram"或"unixpacket"
	func (a *UnixAddr) Network() string
	func (a *UnixAddr) String() string
```

- net.Listener 
```golang
	// Listener 是一个用于面向流的网络协议的公用的网络监听器接口；多个线程可能会同时调用一个Listener的方法
	type Listener interface {
		// Addr返回该接口的网络地址
		Addr() Addr
		// Accept等待并返回下一个连接到该接口的连接
		Accept() (c Conn, err error)
		// Close关闭该接口，并使任何阻塞的Accept操作都会不再阻塞并返回错误
		Close() error
	}

	// 返回在一个本地网络地址laddr上监听的Listener
	// 网络类型参数net必须是面向流的网络："tcp"、"tcp4"、"tcp6"、"unix"或"unixpacket"；参见Dial函数获取laddr的语法
	func Listen(net, laddr string) (Listener, error)

	// TCPListener 见上Tcp

	// UnixListener代表一个Unix域scoket的监听者；使用者应尽量使用Listener接口而不是假设(网络连接为)Unix域scoket
	type UnixListener struct { ... }
	// ListenUnix 在Unix域scoket地址laddr上声明并返回一个*UnixListener，net参数必须是"unix"或"unixpacket"
	func ListenUnix(net string, laddr *UnixAddr) (*UnixListener, error)
	// Addr 返回l的监听的Unix域socket地址
	func (l *UnixListener) Addr() Addr
	// 设置监听器执行的期限，t为Time零值则会关闭期限限制
	func (l *UnixListener) SetDeadline(t time.Time) (err error)
	// Accept 用于实现Listener接口的Accept方法；他会等待下一个呼叫，并返回一个该呼叫的Conn接口
	func (l *UnixListener) Accept() (c Conn, err error)
	// AcceptUnix 接收下一个呼叫，并返回一个新的*UnixConn
	func (l *UnixListener) AcceptUnix() (*UnixConn, error)
	// Close 停止监听Unix域socket地址，已经接收的连接不受影响
	func (l *UnixListener) Close() error
	// File 方法返回下层的os.File的副本，并将该副本设置为阻塞模式
	// 使用者有责任在用完后关闭f。关闭c不影响f，关闭f也不影响c
	// 返回的os.File类型文件描述符和原本的网络连接是不同的
	func (l *UnixListener) File() (f *os.File, err error)
```

- net.Conn
```golang
	// Conn接口代表通用的面向流的网络连接；多个线程可能会同时调用同一个Conn的方法
	type Conn interface {
		// Read从连接中读取数据
		// Read方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
		Read(b []byte) (n int, err error)
		// Write从连接中写入数据
		// Write方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
		Write(b []byte) (n int, err error)
		// Close方法关闭该连接
		// 并会导致任何阻塞中的Read或Write方法不再阻塞并返回错误
		Close() error
		// 返回本地网络地址
		LocalAddr() Addr
		// 返回远端网络地址
		RemoteAddr() Addr
		// 设定该连接的读写deadline，等价于同时调用SetReadDeadline和SetWriteDeadline
		// deadline是一个绝对时间，超过该时间后I/O操作就会直接因超时失败返回而不会阻塞
		// deadline对之后的所有I/O操作都起效，而不仅仅是下一次的读或写操作
		// 参数t为零值表示不设置期限
		SetDeadline(t time.Time) error
		// 设定该连接的读操作deadline，参数t为零值表示不设置期限
		SetReadDeadline(t time.Time) error
		// 设定该连接的写操作deadline，参数t为零值表示不设置期限
		// 即使写入超时，返回值n也可能>0，说明成功写入了部分数据
		SetWriteDeadline(t time.Time) error
	}
	// 在网络network上连接地址address，并返回一个Conn接口
	// 可用的网络类型有："tcp"、"tcp4"、"tcp6"、"udp"、"udp4"、"udp6"、"ip"、"ip4"、"ip6"、"unix"、"unixgram"、"unixpacket" 
	// 对TCP和UDP网络，地址格式是host:port或[host]:port，参见函数JoinHostPort和SplitHostPort
	func Dial(network, address string) (Conn, error)
	// DialTimeout 类似Dial但采用了超时；timeout参数如果必要可包含名称解析
	func DialTimeout(network, address string, timeout time.Duration) (Conn, error)

	// Pipe 创建一个内存中的同步、全双工网络连接
	// 连接的两端都实现了Conn接口；一端的读取对应另一端的写入，直接将数据在两端之间作拷贝；没有内部缓冲
	func Pipe() (Conn, Conn)

	// Dialer 类型包含与某个地址建立连接时的参数
	// 每一个字段的零值都等价于没有该字段；因此调用Dialer零值的Dial方法等价于调用Dial函数
	type Dialer struct {
		// Timeout 是dial操作等待连接建立的最大时长，默认值代表没有超时
		// 如果Deadline字段也被设置了，dial操作也可能更早失败
		// 不管有没有设置超时，操作系统都可能强制执行它的超时设置
		// 例如，TCP(系统)超时一般在3分钟左右
		Timeout time.Duration
		// Deadline 是一个具体的时间点期限，超过该期限后，dial操作就会失败
		// 如果Timeout字段也被设置了，dial操作也可能更早失败
		// 零值表示没有期限，即遵守操作系统的超时设置
		Deadline time.Time
		// LocalAddr 是dial一个地址时使用的本地地址
		// 该地址必须是与dial的网络相容的类型
		// 如果为nil，将会自动选择一个本地地址
		LocalAddr Addr
		// DualStack 允许单次dial操作在网络类型为"tcp"，
		// 且目的地是一个主机名的DNS记录具有多个地址时，
		// 尝试建立多个IPv4和IPv6连接，并返回第一个建立的连接
		DualStack bool
		// KeepAlive 指定一个活动的网络连接的生命周期；如果为0，会禁止keep-alive
		// 不支持keep-alive的网络连接会忽略本字段
		KeepAlive time.Duration
	}
	// Dial在指定的网络上连接指定的地址，参见Dial函数获取网络和地址参数的描述
	func (d *Dialer) Dial(network, address string) (Conn, error)

	// PacketConn 接口代表通用的面向数据包的网络连接；多个线程可能会同时调用同一个Conn的方法
	type PacketConn interface {
		// ReadFrom 方法从连接读取一个数据包，并将有效信息写入b
		// ReadFrom 方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
		// 返回写入的字节数和该数据包的来源地址
		ReadFrom(b []byte) (n int, addr Addr, err error)
		// WriteTo 方法将有效数据b写入一个数据包发送给addr
		// WriteTo 方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
		// 在面向数据包的连接中，写入超时非常罕见
		WriteTo(b []byte, addr Addr) (n int, err error)
		// Close 方法关闭该连接
		// 会导致任何阻塞中的ReadFrom或WriteTo方法不再阻塞并返回错误
		Close() error
		// 返回本地网络地址
		LocalAddr() Addr
		// 设定该连接的读写deadline
		SetDeadline(t time.Time) error
		// 设定该连接的读操作deadline，参数t为零值表示不设置期限
		// 如果时间到达deadline，读操作就会直接因超时失败返回而不会阻塞
		SetReadDeadline(t time.Time) error
		// 设定该连接的写操作deadline，参数t为零值表示不设置期限
		// 如果时间到达deadline，写操作就会直接因超时失败返回而不会阻塞
		// 即使写入超时，返回值n也可能>0，说明成功写入了部分数据
		SetWriteDeadline(t time.Time) error
	}
	// ListenPacket函数监听本地网络地址laddr
	// 网络类型net必须是面向数据包的网络类型："ip"、"ip4"、"ip6"、"udp"、"udp4"、"udp6"、或"unixgram"；laddr的格式参见Dial函数
	func ListenPacket(net, laddr string) (PacketConn, error)

	// IPConn 类型代表IP网络连接，实现了Conn和PacketConn接口
	type IPConn struct { ... }
	// DialIP 在网络协议netProto上连接本地地址laddr和远端地址raddr，netProto必须是"ip"、"ip4"或"ip6"后跟冒号和协议名或协议号
	func DialIP(netProto string, laddr, raddr *IPAddr) (*IPConn, error)
	// ListenIP 创建一个接收目的地是本地地址laddr的IP数据包的网络连接，返回的*IPConn的ReadFrom和WriteTo方法可以用来发送和接收IP数据包(每个包都可获取来源址或者设置目标地址)
	func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error)
	// LocalAddr 返回本地网络地址
	func (c *IPConn) LocalAddr() Addr
	// RemoteAddr 返回远端网络地址
	func (c *IPConn) RemoteAddr() Addr
	// SetReadBuffer 设置该连接的系统接收缓冲
	func (c *IPConn) SetReadBuffer(bytes int) error
	// SetWriteBuffer 设置该连接的系统发送缓冲
	func (c *IPConn) SetWriteBuffer(bytes int) error
	// SetDeadline 设置读写操作绝对期限，实现了Conn接口的SetDeadline方法
	func (c *IPConn) SetDeadline(t time.Time) error
	// SetReadDeadline 设置读操作绝对期限，实现了Conn接口的SetReadDeadline方法
	func (c *IPConn) SetReadDeadline(t time.Time) error
	// SetWriteDeadline 设置写操作绝对期限，实现了Conn接口的SetWriteDeadline方法
	func (c *IPConn) SetWriteDeadline(t time.Time) error
	// Read实现Conn接口Read方法
	func (c *IPConn) Read(b []byte) (int, error)
	// ReadFrom实现PacketConn接口ReadFrom方法
	// 注意本方法有bug，应避免使用
	func (c *IPConn) ReadFrom(b []byte) (int, Addr, error)
	// ReadFromIP 从c读取一个IP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址
	// ReadFromIP 方法会在超过一个固定的时间点之后超时，并返回一个错误；注意本方法有bug，应避免使用
	func (c *IPConn) ReadFromIP(b []byte) (int, *IPAddr, error)
	// ReadMsgIP 从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误
	func (c *IPConn) ReadMsgIP(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error)
	// Write 实现Conn接口Write方法
	func (c *IPConn) Write(b []byte) (int, error)
	// WriteTo 实现PacketConn接口WriteTo方法
	func (c *IPConn) WriteTo(b []byte, addr Addr) (int, error)
	// WriteToIP 通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节
	// WriteToIP 方法会在超过一个固定的时间点之后超时，并返回一个错误；在面向数据包的连接上，写入超时是十分罕见的
	func (c *IPConn) WriteToIP(b []byte, addr *IPAddr) (int, error)
	// WriteMsgIP 通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数(包数据、带外数据)和可能的错误
	func (c *IPConn) WriteMsgIP(b, oob []byte, addr *IPAddr) (n, oobn int, err error)
	// Close 关闭连接
	func (c *IPConn) Close() error
	// File 方法设置下层的os.File为阻塞模式并返回其副本
	// 使用者有责任在用完后关闭f，关闭c不影响f，关闭f也不影响c；返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会(也可能不会)得到期望的效果
	func (c *IPConn) File() (f *os.File, err error)

	// TCPConn 见上Tcp

	// UDPConn 见上Udp

	// UnixConn 代表一个Unix域socket终端地址UnixConn代表Unix域socket连接，实现了Conn和PacketConn接口
	type UnixConn struct { ... }
	// DialUnix 在网络协议net上连接本地地址laddr和远端地址raddr
	// net必须是"unix"、"unixgram"、"unixpacket"，如果laddr不是nil将使用它作为本地地址，否则自动选择一个本地地址
	func DialUnix(net string, laddr, raddr *UnixAddr) (*UnixConn, error)
	// ListenUnixgram接收目的地是本地地址laddr的Unix datagram网络连接
	// net必须是"unixgram"，返回的*UnixConn的ReadFrom和WriteTo方法可以用来发送和接收数据包(每个包都可获取来源址或者设置目标地址)
	func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error)
	// LocalAddr 返回本地网络地址
	func (c *UnixConn) LocalAddr() Addr
	// RemoteAddr 返回远端网络地址
	func (c *UnixConn) RemoteAddr() Addr
	// SetReadBuffer 设置该连接的系统接收缓冲
	func (c *UnixConn) SetReadBuffer(bytes int) error
	// SetWriteBuffer 设置该连接的系统发送缓冲
	func (c *UnixConn) SetWriteBuffer(bytes int) error
	// SetDeadline 设置读写操作绝对期限，实现了Conn接口的SetDeadline方法
	func (c *UnixConn) SetDeadline(t time.Time) error
	// SetReadDeadline 设置读操作绝对期限，实现了Conn接口的SetReadDeadline方法
	func (c *UnixConn) SetReadDeadline(t time.Time) error
	// SetWriteDeadline 设置写操作绝对期限，实现了Conn接口的SetWriteDeadline方法
	func (c *UnixConn) SetWriteDeadline(t time.Time) error
	// Read实现Conn接口Read方法
	func (c *UnixConn) Read(b []byte) (int, error)
	// ReadFrom实现PacketConn接口ReadFrom方法
	func (c *UnixConn) ReadFrom(b []byte) (int, Addr, error)
	// ReadFromUnix 从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址
	// ReadFromUnix 方法会在超过一个固定的时间点之后超时，并返回一个错误；
	func (c *UnixConn) ReadFromUnix(b []byte) (n int, addr *UnixAddr, err error)
	// ReadMsgUnix 从c读取一个数据包，将有效负载拷贝进b，相关的带外数据拷贝进oob，返回拷贝进b的字节数，拷贝进oob的字节数，数据包的flag，数据包来源地址和可能的错误
	func (c *UnixConn) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error)
	// Write 实现Conn接口Write方法
	func (c *UnixConn) Write(b []byte) (int, error)
	// WriteTo 实现PacketConn接口WriteTo方法
	func (c *UnixConn) WriteTo(b []byte, addr Addr) (int, error)
	// WriteToUnix 通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节
	// WriteToUnix 方法会在超过一个固定的时间点之后超时，并返回一个错误；在面向数据包的连接上，写入超时是十分罕见的
	func (c *UnixConn) WriteToUnix(b []byte, addr *UnixAddr) (n int, err error)
	// WriteMsgUnix 通过c向地址addr发送一个数据包，b和oob分别为包有效负载和对应的带外数据，返回写入的字节数(包数据、带外数据)和可能的错误
	func (c *UnixConn) WriteMsgUnix(b, oob []byte, addr *UnixAddr) (n, oobn int, err error)
	// Close 关闭连接
	func (c *UnixConn) Close() error
	// CloseRead关闭TCP连接的读取侧(以后不能读取)，应尽量使用Close方法
	func (c *UnixConn) CloseRead() errorv
	// CloseWrite关闭TCP连接的写入侧(以后不能写入)，应尽量使用Close方法
	func (c *UnixConn) CloseWrite() error
	// File 方法设置下层的os.File为阻塞模式并返回其副本
	// 使用者有责任在用完后关闭f，关闭c不影响f，关闭f也不影响c；返回的os.File类型文件描述符和原本的网络连接是不同的。试图使用该副本修改本体的属性可能会(也可能不会)得到期望的效果
	func (c *UnixConn) File() (f *os.File, err error)
```

- net.Conn
```golang
	// FileListener 返回一个下层为文件f的网络监听器的拷贝
	// 调用者有责任在使用结束后改变l。关闭l不会影响f；关闭f也不会影响l；本函数与各种实现了Listener接口的类型的File方法是对应的
	func FileListener(f *os.File) (l Listener, err error)

	// FileConn 返回一个下层为文件f的网络连接的拷贝
	// 调用者有责任在结束程序前关闭f；关闭c不会影响f，关闭f也不会影响c；本函数与各种实现了Conn接口的类型的File方法是对应的
	func FileConn(f *os.File) (c Conn, err error)

	// FilePacketConn函数返回一个下层为文件f的数据包网络连接的拷贝
	// 调用者有责任在结束程序前关闭f；关闭c不会影响f，关闭f也不会影响c；本函数与各种实现了PacketConn接口的类型的File方法是对应的
	func FilePacketConn(f *os.File) (c PacketConn, err error)
```

- net.DNS
```golang
	// MX 代表一条DNS MX记录(邮件交换记录)，根据收信人的地址后缀来定位邮件服务器
	type MX struct {
		Host string
		Pref uint16
	}

	// NS 代表一条DNS NS记录(域名服务器记录)，指定该域名由哪个DNS服务器来进行解析
	type NS struct {
		Host string
	}

	// SRV 代表一条DNS SRV记录(资源记录)，记录某个服务由哪台计算机提供
	type SRV struct {
		Target   string
		Port     uint16
		Priority uint16
		Weight   uint16
	}

	// LookupPort 函数查询指定网络和服务的(默认)端口
	func LookupPort(network, service string) (port int, err error)
	// LookupCNAME 函数查询name的规范DNS名(但该域名未必可以访问)
	// 如果调用者不关心规范名可以直接调用LookupHost或者LookupIP；这两个函数都会在查询时考虑到规范名
	func LookupCNAME(name string) (cname string, err error)
	// LookupHost 函数查询主机的网络地址序列
	func LookupHost(host string) (addrs []string, err error)
	// LookupIP 函数查询主机的ipv4和ipv6地址序列
	func LookupIP(host string) (addrs []IP, err error)
	// LookupAddr 查询某个地址，返回映射到该地址的主机名序列，本函数和LookupHost不互为反函数
	func LookupAddr(addr string) (name []string, err error)
	// LookupMX 函数返回指定主机的按Pref字段排好序的DNS MX记录
	func LookupMX(name string) (mx []*MX, err error)
	// LookupNS 函数返回指定主机的DNS NS记录
	func LookupNS(name string) (ns []*NS, err error)
	// LookupSRV 函数尝试执行指定服务、协议、主机的SRV查询
	// 协议proto为"tcp" 或"udp"。返回的记录按Priority字段排序，同一优先度按Weight字段随机排序
	func LookupSRV(service, proto, name string) (cname string, addrs []*SRV, err error)
	// LookupTXT函数返回指定主机的DNS TXT记录
	func LookupTXT(name string) (txt []string, err error)
```

### 2. net/http
- net/http
	- http包提供了HTTP客户端和服务端的实现

- http.const
```golang
	// HTTP状态码
	const (
		StatusContinue           = 100
		StatusSwitchingProtocols = 101
		StatusOK                   = 200
		StatusCreated              = 201
		StatusAccepted             = 202
		StatusNonAuthoritativeInfo = 203
		StatusNoContent            = 204
		StatusResetContent         = 205
		StatusPartialContent       = 206
		StatusMultipleChoices   = 300
		StatusMovedPermanently  = 301
		StatusFound             = 302
		StatusSeeOther          = 303
		StatusNotModified       = 304
		StatusUseProxy          = 305
		StatusTemporaryRedirect = 307
		StatusBadRequest                   = 400
		StatusUnauthorized                 = 401
		StatusPaymentRequired              = 402
		StatusForbidden                    = 403
		StatusNotFound                     = 404
		StatusMethodNotAllowed             = 405
		StatusNotAcceptable                = 406
		StatusProxyAuthRequired            = 407
		StatusRequestTimeout               = 408
		StatusConflict                     = 409
		StatusGone                         = 410
		StatusLengthRequired               = 411
		StatusPreconditionFailed           = 412
		StatusRequestEntityTooLarge        = 413
		StatusRequestURITooLong            = 414
		StatusUnsupportedMediaType         = 415
		StatusRequestedRangeNotSatisfiable = 416
		StatusExpectationFailed            = 417
		StatusTeapot                       = 418
		StatusInternalServerError     = 500
		StatusNotImplemented          = 501
		StatusBadGateway              = 502
		StatusServiceUnavailable      = 503
		StatusGatewayTimeout          = 504
		StatusHTTPVersionNotSupported = 505
	)
	// DefaultMaxHeaderBytes 是HTTP请求的头域最大允许长度；可以通过设置Server.MaxHeaderBytes字段来覆盖
	const DefaultMaxHeaderBytes = 1 << 20 // 1 MB
	// DefaultMaxIdleConnsPerHost 是Transport的MaxIdleConnsPerHost的默认值
	const DefaultMaxIdleConnsPerHost = 2
	// TimeFormat 是当解析或生产HTTP头域中的时间时，用与time.Parse或time.Format函数的时间格式
	const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
```

- http.variables&&error
```golang
	// DefaultClient 是用于包函数Get、Head和Post的默认Client
	var DefaultClient = &Client{}
	// DefaultServeMux 是用于Serve的默认ServeMux
	var DefaultServeMux = NewServeMux()

	// HTTP请求的解析错误
	var (
		ErrHeaderTooLong        = &ProtocolError{"header too long"}
		ErrShortBody            = &ProtocolError{"entity body too short"}
		ErrNotSupported         = &ProtocolError{"feature not supported"}
		ErrUnexpectedTrailer    = &ProtocolError{"trailer header without chunked transfer encoding"}
		ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
		ErrNotMultipart         = &ProtocolError{"request Content-Type isn't multipart/form-data"}
		ErrMissingBoundary      = &ProtocolError{"no multipart boundary param in Content-Type"}
	)
	// 会被HTTP服务端返回的错误
	var (
		ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
		ErrBodyNotAllowed  = errors.New("http: request method or response status code does not allow body")
		ErrHijacked        = errors.New("Conn has been hijacked")
		ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
	)
	// 在Resquest或Response的Body字段已经关闭后，试图从中读取时，就会返回ErrBodyReadAfterClose
	// 这个错误一般发生在：HTTP处理器中调用完ResponseWriter 接口的WriteHeader或Write后从请求中读取数据的时候
	var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")
	// 在处理器超时以后调用ResponseWriter接口的Write方法，就会返回ErrHandlerTimeout
	var ErrHandlerTimeout = errors.New("http: Handler timeout")
	var ErrLineTooLong = errors.New("header line too long")
	// 当请求中没有提供给FormFile函数的文件字段名，或者该字段名不是文件字段时，该函数就会返回ErrMissingFile
	var ErrMissingFile = errors.New("http: no such file")
	var ErrNoCookie = errors.New("http: named cookie not present")
	var ErrNoLocation = errors.New("http: no Location header in response")

	// HTTP请求解析错误
	type ProtocolError struct {
		ErrorString string
	}
	func (err *ProtocolError) Error() string
```

- http.ConnState
```golang
	const (
		// StateNew代表一个新的连接，将要立刻发送请求。
		// 连接从这个状态开始，然后转变为StateAlive或StateClosed。
		StateNew ConnState = iota
		// StateActive代表一个已经读取了请求数据1到多个字节的连接。
		// 用于StateAlive的Server.ConnState回调函数在将连接交付给处理器之前被触发，
		// 等到请求被处理完后，Server.ConnState回调函数再次被触发。
		// 在请求被处理后，连接状态改变为StateClosed、StateHijacked或StateIdle。
		StateActive
		// StateIdle代表一个已经处理完了请求、处在闲置状态、等待新请求的连接。
		// 连接状态可以从StateIdle改变为StateActive或StateClosed。
		StateIdle
		// 代表一个被劫持的连接。这是一个终止状态，不会转变为StateClosed。
		StateHijacked
		// StateClosed代表一个关闭的连接。
		// 这是一个终止状态。被劫持的连接不会转变为StateClosed。
		StateClosed
	)

	// ConnState 代表一个客户端到服务端的连接的状态；本类型用于可选的Server.ConnState回调函数
	type ConnState int
	func (c ConnState) String() string
```

- http.Header
```golang
	// Header 代表HTTP头域的键值对
	type Header map[string][]string

	// Get 返回键对应的第一个值，如果键不存在会返回""
	// 如要获取该键对应的值切片，请直接用规范格式的键访问map
	func (h Header) Get(key string) string
	// Set 添加键值对到h，如键已存在则会用只有新值一个元素的切片取代旧值切片
	func (h Header) Set(key, value string)
	// Add 添加键值对到h，如键已存在则会将新的值附加到旧值切片后面
	func (h Header) Add(key, value string)
	// Del删除键值对
	func (h Header) Del(key string)
	// Write以有线格式将头域写入w
	func (h Header) Write(w io.Writer) error
	// WriteSubset 以有线格式将头域写入w
	// 当exclude不为nil时，如果h的键值对的键在exclude中存在且其对应值为真，该键值对就不会被写入w
	func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error
```

- http.Cookie
```golang
	// Cookie 代表一个出现在HTTP回复的头域中Set-Cookie头的值里或者HTTP请求的头域中Cookie头的值里的HTTP cookie
	type Cookie struct {
		Name       string
		Value      string
		Path       string
		Domain     string
		Expires    time.Time
		RawExpires string
		// MaxAge=0表示未设置Max-Age属性
		// MaxAge<0表示立刻删除该cookie，等价于"Max-Age: 0"
		// MaxAge>0表示存在Max-Age属性，单位是秒
		MaxAge   int
		Secure   bool
		HttpOnly bool
		Raw      string
		Unparsed []string // 未解析的“属性-值”对的原始文本
	}
	// String 返回该cookie的序列化结果
	// 如果只设置了Name和Value字段，序列化结果可用于HTTP请求的Cookie头或者HTTP回复的Set-Cookie头；如果设置了其他字段，序列化结果只能用于HTTP回复的Set-Cookie头
	func (c *Cookie) String() string

	// CookieJar管理cookie的存储和在HTTP请求中的使用；CookieJar的实现必须能安全的被多个go程同时使用
	// net/http/cookiejar包提供了一个CookieJar的实现
	type CookieJar interface {
		// SetCookies管理从u的回复中收到的cookie
		// 根据其策略和实现，它可以选择是否存储cookie
		SetCookies(u *url.URL, cookies []*Cookie)
		// Cookies返回发送请求到u时应使用的cookie
		// 本方法有责任遵守RFC 6265规定的标准cookie限制
		Cookies(u *url.URL) []*Cookie
	}
```

- http.Request
```golang
	// Request 类型代表一个服务端接受到的或者客户端发送出去的HTTP请求
	// Request各字段的意义和用途在服务端和客户端是不同的；除了字段本身上方文档，还可参见Request.Write方法和RoundTripper接口的文档
	type Request struct {
		// Method指定HTTP方法(GET、POST、PUT等)。对客户端，""代表GET。
		Method string
		// URL在服务端表示被请求的URI，在客户端表示要访问的URL。
		//
		// 在服务端，URL字段是解析请求行的URI(保存在RequestURI字段)得到的，
		// 对大多数请求来说，除了Path和RawQuery之外的字段都是空字符串。
		// (参见RFC 2616, Section 5.1.2)
		//
		// 在客户端，URL的Host字段指定了要连接的服务器，
		// 而Request的Host字段(可选地)指定要发送的HTTP请求的Host头的值。
		URL *url.URL
		// 接收到的请求的协议版本。本包生产的Request总是使用HTTP/1.1
		Proto      string // "HTTP/1.0"
		ProtoMajor int    // 1
		ProtoMinor int    // 0
		// Header字段用来表示HTTP请求的头域。如果头域(多行键值对格式)为：
		//	accept-encoding: gzip, deflate
		//	Accept-Language: en-us
		//	Connection: keep-alive
		// 则：
		//	Header = map[string][]string{
		//		"Accept-Encoding": {"gzip, deflate"},
		//		"Accept-Language": {"en-us"},
		//		"Connection": {"keep-alive"},
		//	}
		// HTTP规定头域的键名(头名)是大小写敏感的，请求的解析器通过规范化头域的键名来实现这点。
		// 在客户端的请求，可能会被自动添加或重写Header中的特定的头，参见Request.Write方法。
		Header Header
		// Body是请求的主体。
		//
		// 在客户端，如果Body是nil表示该请求没有主体买入GET请求。
		// Client的Transport字段会负责调用Body的Close方法。
		//
		// 在服务端，Body字段总是非nil的；但在没有主体时，读取Body会立刻返回EOF。
		// Server会关闭请求的主体，ServeHTTP处理器不需要关闭Body字段。
		Body io.ReadCloser
		// ContentLength记录相关内容的长度。
		// 如果为-1，表示长度未知，如果>=0，表示可以从Body字段读取ContentLength字节数据。
		// 在客户端，如果Body非nil而该字段为0，表示不知道Body的长度。
		ContentLength int64
		// TransferEncoding按从最外到最里的顺序列出传输编码，空切片表示"identity"编码。
		// 本字段一般会被忽略。当发送或接受请求时，会自动添加或移除"chunked"传输编码。
		TransferEncoding []string
		// Close在服务端指定是否在回复请求后关闭连接，在客户端指定是否在发送请求后关闭连接。
		Close bool
		// 在服务端，Host指定URL会在其上寻找资源的主机。
		// 根据RFC 2616，该值可以是Host头的值，或者URL自身提供的主机名。
		// Host的格式可以是"host:port"。
		//
		// 在客户端，请求的Host字段(可选地)用来重写请求的Host头。
		// 如过该字段为""，Request.Write方法会使用URL字段的Host。
		Host string
		// Form是解析好的表单数据，包括URL字段的query参数和POST或PUT的表单数据。
		// 本字段只有在调用ParseForm后才有效。在客户端，会忽略请求中的本字段而使用Body替代。
		Form url.Values
		// PostForm是解析好的POST或PUT的表单数据。
		// 本字段只有在调用ParseForm后才有效。在客户端，会忽略请求中的本字段而使用Body替代。
		PostForm url.Values
		// MultipartForm是解析好的多部件表单，包括上传的文件。
		// 本字段只有在调用ParseMultipartForm后才有效。
		// 在客户端，会忽略请求中的本字段而使用Body替代。
		MultipartForm *multipart.Form
		// Trailer指定了会在请求主体之后发送的额外的头域。
		//
		// 在服务端，Trailer字段必须初始化为只有trailer键，所有键都对应nil值。
		// (客户端会声明哪些trailer会发送)
		// 在处理器从Body读取时，不能使用本字段。
		// 在从Body的读取返回EOF后，Trailer字段会被更新完毕并包含非nil的值。
		// (如果客户端发送了这些键值对)，此时才可以访问本字段。
		//
		// 在客户端，Trail必须初始化为一个包含将要发送的键值对的映射。(值可以是nil或其终值)
		// ContentLength字段必须是0或-1，以启用"chunked"传输编码发送请求。
		// 在开始发送请求后，Trailer可以在读取请求主体期间被修改，
		// 一旦请求主体返回EOF，调用者就不可再修改Trailer。
		//
		// 很少有HTTP客户端、服务端或代理支持HTTP trailer。
		Trailer Header
		// RemoteAddr允许HTTP服务器和其他软件记录该请求的来源地址，一般用于日志。
		// 本字段不是ReadRequest函数填写的，也没有定义格式。
		// 本包的HTTP服务器会在调用处理器之前设置RemoteAddr为"IP:port"格式的地址。
		// 客户端会忽略请求中的RemoteAddr字段。
		RemoteAddr string
		// RequestURI是被客户端发送到服务端的请求的请求行中未修改的请求URI
		// (参见RFC 2616, Section 5.1)
		// 一般应使用URI字段，在客户端设置请求的本字段会导致错误。
		RequestURI string
		// TLS字段允许HTTP服务器和其他软件记录接收到该请求的TLS连接的信息
		// 本字段不是ReadRequest函数填写的。
		// 对启用了TLS的连接，本包的HTTP服务器会在调用处理器之前设置TLS字段，否则将设TLS为nil。
		// 客户端会忽略请求中的TLS字段。
		TLS *tls.ConnectionState
	}

	// NewRequest 使用指定的方法、网址和可选的主题创建并返回一个新的*Request
	// 如果body参数实现了io.Closer接口，Request返回值的Body 字段会被设置为body，并会被Client类型的Do、Post和PostFOrm方法以及Transport.RoundTrip方法关闭
	func NewRequest(method, urlStr string, body io.Reader) (*Request, error)
	// ReadRequest 从b读取并解析出一个HTTP请求(本函数主要用在服务端从下层获取请求)
	func ReadRequest(b *bufio.Reader) (req *Request, err error)
	// ProtoAtLeast 报告该请求使用的HTTP协议版本至少是major.minor
	func (r *Request) ProtoAtLeast(major, minor int) bool
	// UserAgent 返回请求中的客户端用户代理信息(请求的User-Agent头)
	func (r *Request) UserAgent() string
	// Referer 返回请求中的访问来路信息(请求的Referer头)
	func (r *Request) Referer() string
	// AddCookie 向请求中添加一个cookie，按照RFC 6265 section 5.4的规定，AddCookie不会添加超过一个Cookie头字段，这表示所有的cookie都写在同一行，用分号分隔(cookie内部用逗号分隔属性)
	func (r *Request) AddCookie(c *Cookie)
	// SetBasicAuth 使用提供的用户名和密码，采用HTTP基本认证，设置请求的Authorization头，HTTP基本认证会明码传送用户名和密码
	func (r *Request) SetBasicAuth(username, password string)
	// Write 方法以有线格式将HTTP/1.1请求写入w(用于将请求写入下层TCPConn等)
	// // 本方法会考虑请求的如下字段：
	// // Host
	// // URL
	// // Method (defaults to "GET")
	// // Header
	// // ContentLength
	// // TransferEncoding
	// // Body
	// // 如果存在Body，ContentLength字段<= 0且TransferEncoding字段未显式设置为["identity"]，Write方法会显式添加"Transfer-Encoding: chunked"到请求的头域；Body字段会在发送完请求后关闭
	func (r *Request) Write(w io.Writer) error
	// WriteProxy 类似Write但会将请求以HTTP代理期望的格式发送
	func (r *Request) WriteProxy(w io.Writer) error
	// Cookies 解析并返回该请求的Cookie头设置的cookie
	func (r *Request) Cookies() []*Cookie
	// Cookie 返回请求中名为name的cookie，如果未找到该cookie会返回nil, ErrNoCookie
	func (r *Request) Cookie(name string) (*Cookie, error)
	// ParseForm 解析URL中的查询字符串，并将解析结果更新到r.Form字段
	// 对于POST或PUT请求，ParseForm还会将body当作表单解析，并将结果既更新到r.PostForm也更新到r.Form；解析结果中，POST或PUT请求主体要优先于URL查询字符串(同名变量，主体的值在查询字符串的值前面)
	// 如果请求的主体的大小没有被MaxBytesReader函数设定限制，其大小默认限制为开头10MB
	// ParseMultipartForm会自动调用ParseForm；重复调用本方法是无意义的
	func (r *Request) ParseForm() error
	// ParseMultipartForm 将请求的主体作为multipart/form-data解析
	// 请求的整个主体都会被解析，得到的文件记录最多maxMemery字节保存在内存，其余部分保存在硬盘的temp文件里；如果必要，ParseMultipartForm会自行调用ParseForm；重复调用本方法是无意义的。
	func (r *Request) ParseMultipartForm(maxMemory int64) error
	// FormValue 返回key为键查询r.Form字段得到结果[]string切片的第一个值
	// POST和PUT主体中的同名参数优先于URL查询字符串；如果必要，本函数会隐式调用ParseMultipartForm和ParseForm
	func (r *Request) FormValue(key string) string
	// PostFormValue返回key为键查询r.PostForm字段得到结果[]string切片的第一个值；如果必要，本函数会隐式调用ParseMultipartForm和ParseForm
	func (r *Request) PostFormValue(key string) string
	// FormFile返回以key为键查询r.MultipartForm字段得到结果中的第一个文件和它的信息；如果必要，本函数会隐式调用ParseMultipartForm和ParseForm。查询失败会返回ErrMissingFile错误
	func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
	// 如果请求是multipart/form-data POST请求，MultipartReader返回一个multipart.Reader接口，否则返回nil和一个错误。使用本函数代替ParseMultipartForm，可以将r.Body作为流处理
	func (r *Request) MultipartReader() (*multipart.Reader, error)
```

- http.Response
```golang
	// Response 代表一个HTTP请求的回复
	type Response struct {
		Status     string // 例如"200 OK"
		StatusCode int    // 例如200
		Proto      string // 例如"HTTP/1.0"
		ProtoMajor int    // 例如1
		ProtoMinor int    // 例如0
		// Header保管头域的键值对。
		// 如果回复中有多个头的键相同，Header中保存为该键对应用逗号分隔串联起来的这些头的值
		// (参见RFC 2616 Section 4.2)
		// 被本结构体中的其他字段复制保管的头(如ContentLength)会从Header中删掉。
		//
		// Header中的键都是规范化的，参见CanonicalHeaderKey函数
		Header Header
		// Body代表回复的主体。
		// Client类型和Transport类型会保证Body字段总是非nil的，即使回复没有主体或主体长度为0。
		// 关闭主体是调用者的责任。
		// 如果服务端采用"chunked"传输编码发送的回复，Body字段会自动进行解码。
		Body io.ReadCloser
		// ContentLength记录相关内容的长度。
		// 其值为-1表示长度未知(采用chunked传输编码)
		// 除非对应的Request.Method是"HEAD"，其值>=0表示可以从Body读取的字节数
		ContentLength int64
		// TransferEncoding按从最外到最里的顺序列出传输编码，空切片表示"identity"编码。
		TransferEncoding []string
		// Close记录头域是否指定应在读取完主体后关闭连接。(即Connection头)
		// 该值是给客户端的建议，Response.Write方法的ReadResponse函数都不会关闭连接。
		Close bool
		// Trailer字段保存和头域相同格式的trailer键值对，和Header字段相同类型
		Trailer Header
		// Request是用来获取此回复的请求
		// Request的Body字段是nil(因为已经被用掉了)
		// 这个字段是被Client类型发出请求并获得回复后填充的
		Request *Request
		// TLS包含接收到该回复的TLS连接的信息。 对未加密的回复，本字段为nil。
		// 返回的指针是被(同一TLS连接接收到的)回复共享的，不应被修改。
		TLS *tls.ConnectionState
	}

	// ReadResponse从r读取并返回一个HTTP 回复
	// req参数是可选的，指定该回复对应的请求(即是对该请求的回复)；如果是nil，将假设请求是GET请求
	// 客户端必须在结束resp.Body的读取后关闭它；读取完毕并关闭后，客户端可以检查resp.Trailer字段获取回复的trailer的键值对(本函数主要用在客户端从下层获取回复)
	func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)
	// ProtoAtLeast 报告该回复使用的HTTP协议版本至少是major.minor
	func (r *Response) ProtoAtLeast(major, minor int) bool
	// Cookies 解析并返回该回复中的Set-Cookie头设置的cookie
	func (r *Response) Cookies() []*Cookie
	// Location 返回该回复的Location头设置的URL。相对地址的重定向会相对于该回复对应的请求来确定绝对地址。如果回复中没有Location头，会返回nil, ErrNoLocation
	func (r *Response) Location() (*url.URL, error)
	// Write 以有线格式将回复写入w(用于将回复写入下层TCPConn等)
	// // 本方法会考虑如下字段：
	// // StatusCode
	// // ProtoMajor
	// // ProtoMinor
	// // Request.Method
	// // TransferEncoding
	// // Trailer
	// // Body
	// // ContentLength
	// // Header(不规范的键名和它对应的值会导致不可预知的行为)
	// // Body字段在发送完回复后会被关闭
	func (r *Response) Write(w io.Writer) error

	// ResponseWriter 接口被HTTP处理器用于构造HTTP回复
	type ResponseWriter interface {
		// Header返回一个Header类型值，该值会被WriteHeader方法发送。
		// 在调用WriteHeader或Write方法后再改变该对象是没有意义的。
		Header() Header
		// WriteHeader该方法发送HTTP回复的头域和状态码。
		// 如果没有被显式调用，第一次调用Write时会触发隐式调用WriteHeader(http.StatusOK)
		// WriterHeader的显式调用主要用于发送错误码。
		WriteHeader(int)
		// Write向连接中写入作为HTTP的一部分回复的数据。
		// 如果被调用时还未调用WriteHeader，本方法会先调用WriteHeader(http.StatusOK)
		// 如果Header中没有"Content-Type"键，
		// 本方法会使用包函数DetectContentType检查数据的前512字节，将返回值作为该键的值。
		Write([]byte) (int, error)
	}

	// HTTP处理器ResponseWriter接口参数的下层如果实现了Flusher接口，可以让HTTP处理器将缓冲中的数据发送到客户端
	// 注意：即使ResponseWriter接口的下层支持Flush方法，如果客户端是通过HTTP代理连接的，缓冲中的数据也可能直到回复完毕才被传输到客户端
	type Flusher interface {
		// Flush将缓冲中的所有数据发送到客户端
		Flush()
	}

	// HTTP处理器ResponseWriter接口参数的下层如果实现了CloseNotifier接口，可以让用户检测下层的连接是否停止；如果客户端在回复准备好之前关闭了连接，该机制可以用于取消服务端耗时较长的操作
	type CloseNotifier interface {
		// CloseNotify返回一个通道，该通道会在客户端连接丢失时接收到唯一的值
		CloseNotify() <-chan bool
	}

	// HTTP处理器ResponseWriter接口参数的下层如果实现了Hijacker接口，可以让HTTP处理器接管该连接
	type Hijacker interface {
		// Hijack让调用者接管连接，返回连接和关联到该连接的一个缓冲读写器
		// 调用本方法后，HTTP服务端将不再对连接进行任何操作，调用者有责任管理、关闭返回的连接
		Hijack() (net.Conn, *bufio.ReadWriter, error)
	}

	// SetCookie 在w的头域中添加Set-Cookie头，该HTTP头的值为cookie
	func SetCookie(w ResponseWriter, cookie *Cookie)
	// Redirect 回复请求一个重定向地址urlStr和状态码code，该重定向地址可以是相对于请求r的相对地址
	func Redirect(w ResponseWriter, r *Request, urlStr string, code int)
	// NotFound回复请求404状态码(not found：目标未发现)
	func NotFound(w ResponseWriter, r *Request)
	// Error 使用指定的错误信息和状态码回复请求，将数据写入w；错误信息必须是明文
	func Error(w ResponseWriter, error string, code int)
	// ServeContent 使用提供的ReadSeeker的内容回复请求；
	// ServeContent比起io.Copy函数的主要优点，是可以处理范围类请求(只要一部分内容)、设置MIME类型，处理If-Modified-Since请求
	// 如果未设定回复的Content-Type头，本函数首先会尝试从name的文件扩展名推断数据类型；如果失败，会用读取content的第1块数据并提供给DetectContentType推断类型；之后会设置Content-Type头。参数name不会用于别的地方，甚至于它可以是空字符串，也永远不会发送到回复里
	// 如果modtime不是Time零值，函数会在回复的头域里设置Last-Modified头；如果请求的头域包含If-Modified-Since头，本函数会使用modtime参数来确定是否应该发送内容；如果调用者设置了w的ETag头，ServeContent会使用它处理包含If-Range头和If-None-Match头的请求。
	// 参数content的Seek方法必须有效：函数使用Seek来确定它的大小
	// 注意：本包File接口和*os.File类型都实现了io.ReadSeeker接口
	func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)
	// ServeFile 回复请求name指定的文件或者目录的内容
	func ServeFile(w ResponseWriter, r *Request, name string)
	// MaxBytesReader 类似io.LimitReader，但它是用来限制接收到的请求的Body的大小的
	// 不同于io.LimitReader，本函数返回一个ReadCloser，返回值的Read方法在读取的数据超过大小限制时会返回非EOF错误，其Close方法会关闭下层的io.ReadCloser接口r
	// MaxBytesReader 预防客户端因为意外或者蓄意发送的“大”请求，以避免尺寸过大的请求浪费服务端资源
	func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser
```

- http.Transport
```golang
	// Transport类型实现了RoundTripper接口，支持http、https和http/https代理
	// Transport类型可以缓存连接以在未来重用
	type Transport struct {
		// Proxy指定一个对给定请求返回代理的函数。
		// 如果该函数返回了非nil的错误值，请求的执行就会中断并返回该错误。
		// 如果Proxy为nil或返回nil的*URL置，将不使用代理。
		Proxy func(*Request) (*url.URL, error)
		// Dial指定创建TCP连接的拨号函数；如果Dial为nil，会使用net.Dial。
		Dial func(network, addr string) (net.Conn, error)
		// TLSClientConfig指定用于tls.Client的TLS配置信息。
		// 如果该字段为nil，会使用默认的配置信息。
		TLSClientConfig *tls.Config
		// TLSHandshakeTimeout指定等待TLS握手完成的最长时间。零值表示不设置超时。
		TLSHandshakeTimeout time.Duration
		// 如果DisableKeepAlives为真，会禁止不同HTTP请求之间TCP连接的重用。
		DisableKeepAlives bool
		// 如果DisableCompression为真，会禁止Transport在请求中没有Accept-Encoding头时，
		// 主动添加"Accept-Encoding: gzip"头，以获取压缩数据。
		// 如果Transport自己请求gzip并得到了压缩后的回复，它会主动解压缩回复的主体。
		// 但如果用户显式的请求gzip压缩数据，Transport是不会主动解压缩的。
		DisableCompression bool
		// 如果MaxIdleConnsPerHost!=0，会控制每个主机下的最大闲置连接。
		// 如果MaxIdleConnsPerHost==0，会使用DefaultMaxIdleConnsPerHost。
		MaxIdleConnsPerHost int
		// ResponseHeaderTimeout指定在发送完请求(包括其可能的主体)之后，
		// 等待接收服务端的回复的头域的最大时间。零值表示不设置超时。
		// 该时间不包括获取回复主体的时间。
		ResponseHeaderTimeout time.Duration
		// ...
	}

	// DefaultTransport 是被包变量DefaultClient使用的默认RoundTripper接口
	// 它会根据需要创建网络连接，并缓存以便在之后的请求中重用这些连接；它使用环境变量$HTTP_PROXY和$NO_PROXY(或$http_proxy和$no_proxy)指定的HTTP代理
	var DefaultTransport RoundTripper = &Transport{
		Proxy: ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// RegisterProtocol 注册一个新的名为scheme的协议
	// t会将使用scheme协议的请求转交给rt；rt有责任模拟HTTP请求的语义
	// RegisterProtocol可以被其他包用于提供"ftp"或"file"等协议的实现
	func (t *Transport) RegisterProtocol(scheme string, rt RoundTripper)
	// RoundTrip 方法实现了RoundTripper接口
	// 高层次的HTTP客户端支持(如管理cookie和重定向)请参见Get、Post等函数和Client类型
	func (t *Transport) RoundTrip(req *Request) (resp *Response, err error)
	// CloseIdleConnections 关闭所有之前的请求建立但目前处于闲置状态的连接，本方法不会中断正在使用的连接
	func (t *Transport) CloseIdleConnections()
	// CancelRequest 通过关闭请求所在的连接取消一个执行中的请求
	func (t *Transport) CancelRequest(req *Request)

	// ProxyURL 返回一个代理函数(用于Transport类型)，该函数总是返回同一个URL
	func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)
	// ProxyFromEnvironment f使用环境变量$HTTP_PROXY和$NO_PROXY(或$http_proxy和$no_proxy)的配置返回用于req的代理
	// 如果代理环境不合法将返回错误；如果环境未设定代理或者给定的request不应使用代理时，将返回(nil, nil)；如果req.URL.Host字段是"localhost"(可以有端口号，也可以没有)，也会返回(nil, nil)
	func ProxyFromEnvironment(req *Request) (*url.URL, error)
```

- http.Client
```golang
	// Client类型 代表HTTP客户端；它的零值(DefaultClient)是一个可用的使用DefaultTransport的客户端
	// Client的Transport字段一般会含有内部状态(缓存TCP连接)，因此Client类型值应尽量被重用而不是每次需要都创建新的；Client类型值可以安全的被多个go程同时使用。
	// Client类型的层次比RoundTripper接口(如Transport)高，还会管理HTTP的cookie和重定向等细节
	type Client struct {
		// Transport指定执行独立、单次HTTP请求的机制。
		// 如果Transport为nil，则使用DefaultTransport。
		Transport RoundTripper
		// CheckRedirect指定处理重定向的策略。
		// 如果CheckRedirect不为nil，客户端会在执行重定向之前调用本函数字段。
		// 参数req和via是将要执行的请求和已经执行的请求(切片，越新的请求越靠后)。
		// 如果CheckRedirect返回一个错误，本类型的Get方法不会发送请求req，
		// 而是返回之前得到的最后一个回复和该错误。(包装进url.Error类型里)
		//
		// 如果CheckRedirect为nil，会采用默认策略：连续10此请求后停止。
		CheckRedirect func(req *Request, via []*Request) error
		// Jar指定cookie管理器。
		// 如果Jar为nil，请求中不会发送cookie，回复中的cookie会被忽略。
		Jar CookieJar
		// Timeout指定本类型的值执行请求的时间限制。
		// 该超时限制包括连接时间、重定向和读取回复主体的时间。
		// 计时器会在Head、Get、Post或Do方法返回后继续运作并在超时后中断回复主体的读取。
		//
		// Timeout为零值表示不设置超时。
		//
		// Client实例的Transport字段必须支持CancelRequest方法，
		// 否则Client会在试图用Head、Get、Post或Do方法执行请求时返回错误。
		// 本类型的Transport字段默认值(DefaultTransport)支持CancelRequest方法。
		Timeout time.Duration
	}

	// Do 方法发送请求，返回HTTP回复；它会遵守客户端c设置的策略(如重定向、cookie、认证)
	// 如果客户端的策略(如重定向)返回错误或存在HTTP协议错误时，本方法将返回该错误；如果回应的状态码不是2xx，本方法并不会返回错误
	// 如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它；如果返回值resp的主体未关闭，c下层的RoundTripper接口(一般为Transport类型)可能无法重用resp主体下层保持的TCP连接去执行之后的请求
	// 请求的主体，如果非nil，会在执行后被c.Transport关闭，即使出现错误
	// 一般应使用Get、Post或PostForm方法代替Do方法
	func (c *Client) Do(req *Request) (resp *Response, err error)

	// Head 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Head会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	func (c *Client) Head(url string) (resp *Response, err error)
	// Get 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Get会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	// 如果c.CheckRedirect执行失败或存在HTTP协议错误时，本方法将返回该错误；如果回应的状态码不是2xx，本方法并不会返回错误；如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	func (c *Client) Get(url string) (resp *Response, err error)
	// Post向指定的URL发出一个POST请求
	// bodyType为POST数据的类型，body为POST数据，作为请求的主体；如果参数body实现了io.Closer接口，它会在发送请求后被关闭；调用者有责任在读取完返回值resp的主体后关闭它
	func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)
	// PostForm 向指定的URL发出一个POST请求，url.Values类型的data会被编码为请求的主体
	// POST数据的类型一般会设为"application/x-www-form-urlencoded"。如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)

	// Head 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Head会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	// Head是对包变量DefaultClient的Head方法的包装
	func Head(url string) (resp *Response, err error)
	// Get 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Get会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	// 如果c.CheckRedirect执行失败或存在HTTP协议错误时，本方法将返回该错误；如果回应的状态码不是2xx，本方法并不会返回错误；如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	// Get是对包变量DefaultClient的Get方法的包装
	// For example
	// // res, err := http.Get("http://www.google.com/robots.txt")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // robots, err := ioutil.ReadAll(res.Body)
	// // res.Body.Close()
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // fmt.Printf("%s", robots)
	func Get(url string) (resp *Response, err error)
	// Post 向指定的URL发出一个POST请求
	// bodyType为POST数据的类型，body为POST数据，作为请求的主体；如果参数body实现了io.Closer接口，它会在发送请求后被关闭；调用者有责任在读取完返回值resp的主体后关闭它
	// Post是对包变量DefaultClient的Post方法的包装
	func Post(url string, bodyType string, body io.Reader) (resp *Response, err error)
	// PostForm 向指定的URL发出一个POST请求，url.Values类型的data会被编码为请求的主体
	// 如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	// PostForm是对包变量DefaultClient的PostForm方法的包装
	func PostForm(url string, data url.Values) (resp *Response, err error)

```

- http.Server
```golang
	// Server类型定义了运行HTTP服务端的参数；Server的零值是合法的配置
	type Server struct {
		Addr           string        // 监听的TCP地址，如果为空字符串会使用":http"
		Handler        Handler       // 调用的处理器，如为nil会调用http.DefaultServeMux
		ReadTimeout    time.Duration // 请求的读取操作在超时前的最大持续时间
		WriteTimeout   time.Duration // 回复的写入操作在超时前的最大持续时间
		MaxHeaderBytes int           // 请求的头域最大长度，如为0则用DefaultMaxHeaderBytes
		TLSConfig      *tls.Config   // 可选的TLS配置，用于ListenAndServeTLS方法
		// TLSNextProto(可选地)指定一个函数来在一个NPN型协议升级出现时接管TLS连接的所有权。
		// 映射的键为商谈的协议名；映射的值为函数，该函数的Handler参数应处理HTTP请求，
		// 并且初始化Handler.ServeHTTP的*Request参数的TLS和RemoteAddr字段(如果未设置)。
		// 连接在函数返回时会自动关闭。
		TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
		// ConnState字段指定一个可选的回调函数，该函数会在一个与客户端的连接改变状态时被调用。
		// 参见ConnState类型和相关常数获取细节。
		ConnState func(net.Conn, ConnState)
		// ErrorLog指定一个可选的日志记录器，用于记录接收连接时的错误和处理器不正常的行为。
		// 如果本字段为nil，日志会通过log包的标准日志记录器写入os.Stderr。
		ErrorLog *log.Logger
		// ...
	}

	// SetKeepAlivesEnabled 控制是否允许HTTP闲置连接重用(keep-alive)功能；
	// 默认该功能总是被启用的，只有资源非常紧张的环境或者服务端在关闭进程中时，才应该关闭该功能
	func (s *Server) SetKeepAlivesEnabled(v bool)
	// Serve 会接手监听器l收到的每一个连接，并为每一个连接创建一个新的服务go程；该go程会读取请求，然后调用srv.Handler回复请求
	func (srv *Server) Serve(l net.Listener) error
	// ListenAndServe 监听srv.Addr指定的TCP地址，并且会调用Serve方法接收到的连接；如果srv.Addr为空字符串，会使用":http"
	func (srv *Server) ListenAndServe() error
	// ListenAndServeTLS监听srv.Addr确定的TCP地址，并且会调用Serve方法处理接收到的连接
	// 必须提供证书文件和对应的私钥文件；如果证书是由权威机构签发的，certFile参数必须是顺序串联的服务端证书和CA证书；如果srv.Addr为空字符串，会使用":https"
	func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error

	// Serve 会接手监听器l收到的每一个连接，并为每一个连接创建一个新的服务go程
	// 该go程会读取请求，然后调用handler回复请求；handler参数一般会设为nil，此时会使用DefaultServeMux
	func Serve(l net.Listener, handler Handler) error
	// ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接
	// handler参数一般会设为nil，此时会使用DefaultServeMux
	// For example
	// // // hello world, the web server
	// // func HelloServer(w http.ResponseWriter, req *http.Request) {
	// // 	io.WriteString(w, "hello, world!\n")
	// // }
	// // http.HandleFunc("/hello", HelloServer)
	// // err := http.ListenAndServe(":12345", nil)
	// // if err != nil {
	// // 	log.Fatal("ListenAndServe: ", err)
	// // }
	// // 
	func ListenAndServe(addr string, handler Handler) error
	// ListenAndServeTLS 和ListenAndServe 的行为基本一致，除了它期望HTTPS连接之外
	// 此外，必须提供证书文件和对应的私钥文件；如果证书是由权威机构签发的，certFile参数必须是顺序串联的服务端证书和CA证书；如果srv.Addr为空字符串，会使用":https"
	// 使用crypto/tls包的generate_cert.go文件来生成cert.pem和key.pem两个文件
	// For example
	// // func handler(w http.ResponseWriter, req *http.Request) {
	// // 	w.Header().Set("Content-Type", "text/plain")
	// // 	w.Write([]byte("This is an example server.\n"))
	// // }
	// // http.HandleFunc("/", handler)
	// // log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	// // err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error
```

- http.Handler
```golang
	// 实现了Handler接口的对象可以注册到HTTP服务端，为特定的路径及其子树提供服务
	// ServeHTTP应该将回复的头域和数据写入ResponseWriter接口然后返回；返回标志着该请求已经结束，HTTP服务端可以转移向该连接上的下一个请求
	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
	// NotFoundHandler 返回一个简单的请求处理器，该处理器会对每个请求都回复"404 page not found"
	func NotFoundHandler() Handler
	// RedirectHandler 返回一个请求处理器，该处理器会对每个请求都使用状态码code重定向到网址url
	func RedirectHandler(url string, code int) Handler
	// TimeoutHandler 返回一个采用指定时间限制的请求处理器
	// 返回的Handler会调用h.ServeHTTP去处理每个请求，但如果某一次调用耗时超过了时间限制，该处理器会回复请求状态码503 Service Unavailable，并将msg作为回复的主体(如果msg为空字符串，将发送一个合理的默认信息)；在超时后，h对它的ResponseWriter接口参数的写入操作会返回ErrHandlerTimeout
	func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler
	// StripPrefix 返回一个处理器，该处理器会将请求的URL.Path字段中给定前缀prefix去除后再交由h处理。StripPrefix会向URL.Path字段中没有给定前缀的请求回复404 page not found
	// For example
	// // To serve a directory on disk (/tmp) under an alternate URL
	// // path (/tmpfiles/), use StripPrefix to modify the request
	// // URL's path before the FileServer seest:
	// http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
	func StripPrefix(prefix string, h Handler) Handler

	// HandlerFunc type是一个适配器，通过类型转换让我们可以将普通的函数作为HTTP处理器使用，如果f是一个具有适当签名的函数，HandlerFunc(f)通过调用f实现了Handler接口
	type HandlerFunc func(ResponseWriter, *Request)
	// ServeHTTP 方法会调用f(w, r)
	func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)

	// Handle 注册HTTP处理器handler和对应的模式pattern(注册到DefaultServeMux)
	// 如果该模式已经注册有一个处理器，Handle会panic；ServeMux的文档解释了模式的匹配机制、
	func Handle(pattern string, handler Handler)
	// HandleFunc注册一个处理器函数handler和对应的模式pattern(注册到DefaultServeMux)；ServeMux的文档解释了模式的匹配机制
	func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

- http.ServeMux
```golang
	// ServeMux类型 是HTTP请求的多路转接器；它会将每一个接收的请求的URL与一个注册模式的列表进行匹配，并调用和URL最匹配的模式的处理器
	// 模式是固定的、由根开始的路径，如"/favicon.ico"，或由根开始的子树，如"/images/"(注意结尾的斜杠)；较长的模式优先于较短的模式，因此如果模式"/images/"和"/images/thumbnails/"都注册了处理器，后一个处理器会用于路径以"/images/thumbnails/"开始的请求，前一个处理器会接收到其余的路径在"/images/"子树下的请求
	// 注意，因为以斜杠结尾的模式代表一个由根开始的子树，模式"/"会匹配所有的未被其他注册的模式匹配的路径，而不仅仅是路径"/"
	// 模式也能(可选地)以主机名开始，表示只匹配该主机上的路径；指定主机的模式优先于一般的模式，因此一个注册了两个模式"/codesearch"和"codesearch.google.com/"的处理器不会接管目标为"http://www.google.com/"的请求
	// ServeMux还会注意到请求的URL路径的无害化，将任何路径中包含"."或".."元素的请求重定向到等价的没有这两种元素的URL。(参见path.Clean函数)
	type ServeMux struct { ... }

	// NewServeMux 创建并返回一个新的*ServeMux
	// For example
	// // mux := http.NewServeMux()
	// // mux.Handle("/api/", apiHandler{})
	// // mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// //     // The "/" pattern matches everything, so we need to check
	// //     // that we're at the root here.
	// //     if req.URL.Path != "/" {
	// //         http.NotFound(w, req)
	// //         return
	// //     }
	// //     fmt.Fprintf(w, "Welcome to the home page!")
	// // })
	func NewServeMux() *ServeMux
	// Handle 注册HTTP处理器handler和对应的模式pattern；如果该模式已经注册有一个处理器，Handle会panic
	func (mux *ServeMux) Handle(pattern string, handler Handler)
	// HandleFunc 注册一个处理器函数handler和对应的模式pattern
	func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	// Handler 根据r.Method、r.Host和r.URL.Path等数据，返回将用于处理该请求的HTTP处理器；它总是返回一个非nil的处理器；如果路径不是它的规范格式，将返回内建的用于重定向到等价的规范路径的处理器
	// Handler也会返回匹配该请求的的已注册模式；在内建重定向处理器的情况下，pattern会在重定向后进行匹配。如果没有已注册模式可以应用于该请求，本方法将返回一个内建的"404 page not found"处理器和一个空字符串模式
	func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
	// ServeHTTP 将请求派遣到与请求的URL最匹配的模式对应的处理器
	func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

- http.File
```golang
	// File 是被FileSystem接口的Open方法返回的接口类型，可以被FileServer等函数用于文件访问服务；该接口的方法的行为应该和*os.File类型的同名方法相同
	type File interface {
		io.Closer
		io.Reader
		Readdir(count int) ([]os.FileInfo, error)
		Seek(offset int64, whence int) (int64, error)
		Stat() (os.FileInfo, error)
	}

	// FileSystem 接口实现了对一系列命名文件的访问，文件路径的分隔符为'/'，不管主机操作系统的惯例如何
	type FileSystem interface {
		Open(name string) (File, error)
	}

	// Dir 使用限制到指定目录树的本地文件系统实现了http.FileSystem接口
	// 空Dir被视为"."，即代表当前目录
	type Dir string
	func (d Dir) Open(name string) (File, error)

	// NewFileTransport 返回一个RoundTripper接口，使用FileSystem接口fs提供文件访问服务；返回的RoundTripper接口会忽略接收的请求的URL主机及其他绝大多数属性
	// NewFileTransport函数的典型使用情况是给Transport类型的值注册"file"协议，如下所示：
	// // t := &http.Transport{}
	// // t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	// // c := &http.Client{Transport: t}
	// // res, err := c.Get("file:///etc/passwd")
	// // ...
	func NewFileTransport(fs FileSystem) RoundTripper

	// FileServer 返回一个使用FileSystem接口root提供文件访问服务的HTTP处理器
	// 要使用操作系统的FileSystem接口实现，可使用http.Dir：
	// // http.Handle("/", http.FileServer(http.Dir("/tmp")))
	// // // Simple static webserver:
	// // log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))
	// // // To serve a directory on disk (/tmp) under an alternate URL
	// // // path (/tmpfiles/), use StripPrefix to modify the request
	// // // URL's path before the FileServer sees it:
	// // http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
	func FileServer(root FileSystem) Handler
```

- http.Func
```golang
	// CanonicalHeaderKey 返回头域(表示为Header类型)的键s的规范化格式
	// 范化过程中让单词首字母和'-'后的第一个字母大写，其余字母小写；例如，"accept-encoding"规范化为"Accept-Encoding"
	func CanonicalHeaderKey(s string) string

	// DetectContentType 实现了http://mimesniff.spec.whatwg.org/描述的算法，用于确定数据的Content-Type
	// 函数总是返回一个合法的MIME类型；如果它不能确定数据的类型，将返回"application/octet-stream"；它最多检查数据的前512字节
	func DetectContentType(data []byte) string

	// ParseHTTPVersion 解析HTTP版本字符串。如"HTTP/1.0"返回(1, 0, true)
	func ParseHTTPVersion(vers string) (major, minor int, ok bool)
	// ParseTime 用3种格式TimeFormat, time.RFC850和time.ANSIC尝试解析一个时间头的值(如Date: header)
	func ParseTime(text string) (t time.Time, err error)

	// StatusText 返回HTTP状态码code对应的文本，如220对应"OK"。如果code是未知的状态码，会返回""
	func StatusText(code int) string

	// ProxyURL 返回一个代理函数(用于Transport类型)，该函数总是返回同一个URL
	func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)
	// ProxyFromEnvironment使用环境变量$HTTP_PROXY和$NO_PROXY(或$http_proxy和$no_proxy)的配置返回用于req的代理
	// 如果代理环境不合法将返回错误；如果环境未设定代理或者给定的request不应使用代理时，将返回(nil, nil)；如果req.URL.Host字段是"localhost"(可以有端口号，也可以没有)，也会返回(nil, nil)
	func ProxyFromEnvironment(req *Request) (*url.URL, error)

	// SetCookie 在w的头域中添加Set-Cookie头，该HTTP头的值为cookie
	func SetCookie(w ResponseWriter, cookie *Cookie)
	// Redirect 回复请求一个重定向地址urlStr和状态码code，该重定向地址可以是相对于请求r的相对地址
	func Redirect(w ResponseWriter, r *Request, urlStr string, code int)
	// NotFound回复请求404状态码(not found：目标未发现)
	func NotFound(w ResponseWriter, r *Request)
	// Error 使用指定的错误信息和状态码回复请求，将数据写入w；错误信息必须是明文
	func Error(w ResponseWriter, error string, code int)
	// ServeContent 使用提供的ReadSeeker的内容回复请求；
	// ServeContent比起io.Copy函数的主要优点，是可以处理范围类请求(只要一部分内容)、设置MIME类型，处理If-Modified-Since请求
	// 如果未设定回复的Content-Type头，本函数首先会尝试从name的文件扩展名推断数据类型；如果失败，会用读取content的第1块数据并提供给DetectContentType推断类型；之后会设置Content-Type头。参数name不会用于别的地方，甚至于它可以是空字符串，也永远不会发送到回复里
	// 如果modtime不是Time零值，函数会在回复的头域里设置Last-Modified头；如果请求的头域包含If-Modified-Since头，本函数会使用modtime参数来确定是否应该发送内容；如果调用者设置了w的ETag头，ServeContent会使用它处理包含If-Range头和If-None-Match头的请求。
	// 参数content的Seek方法必须有效：函数使用Seek来确定它的大小
	// 注意：本包File接口和*os.File类型都实现了io.ReadSeeker接口
	func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)
	// ServeFile 回复请求name指定的文件或者目录的内容
	func ServeFile(w ResponseWriter, r *Request, name string)
	// MaxBytesReader 类似io.LimitReader，但它是用来限制接收到的请求的Body的大小的
	// 不同于io.LimitReader，本函数返回一个ReadCloser，返回值的Read方法在读取的数据超过大小限制时会返回非EOF错误，其Close方法会关闭下层的io.ReadCloser接口r
	// MaxBytesReader 预防客户端因为意外或者蓄意发送的“大”请求，以避免尺寸过大的请求浪费服务端资源
	func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser

	// Head 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Head会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	// Head是对包变量DefaultClient的Head方法的包装
	func Head(url string) (resp *Response, err error)
	// Get 向指定的URL发出一个HEAD请求，如果回应的状态码如下，Get会在调用c.CheckRedirect后执行重定向：
	// // 301 (Moved Permanently)
	// // 302 (Found)
	// // 303 (See Other)
	// // 307 (Temporary Redirect)
	// 如果c.CheckRedirect执行失败或存在HTTP协议错误时，本方法将返回该错误；如果回应的状态码不是2xx，本方法并不会返回错误；如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	// Get是对包变量DefaultClient的Get方法的包装
	// For example
	// // res, err := http.Get("http://www.google.com/robots.txt")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // robots, err := ioutil.ReadAll(res.Body)
	// // res.Body.Close()
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // fmt.Printf("%s", robots)
	func Get(url string) (resp *Response, err error)
	// Post 向指定的URL发出一个POST请求
	// bodyType为POST数据的类型，body为POST数据，作为请求的主体；如果参数body实现了io.Closer接口，它会在发送请求后被关闭；调用者有责任在读取完返回值resp的主体后关闭它
	// Post是对包变量DefaultClient的Post方法的包装
	func Post(url string, bodyType string, body io.Reader) (resp *Response, err error)
	// PostForm 向指定的URL发出一个POST请求，url.Values类型的data会被编码为请求的主体
	// 如果返回值err为nil，resp.Body总是非nil的，调用者应该在读取完resp.Body后关闭它
	// PostForm是对包变量DefaultClient的PostForm方法的包装
	func PostForm(url string, data url.Values) (resp *Response, err error)

	// Handle 注册HTTP处理器handler和对应的模式pattern(注册到DefaultServeMux)
	// 如果该模式已经注册有一个处理器，Handle会panic；ServeMux的文档解释了模式的匹配机制、
	func Handle(pattern string, handler Handler)
	// HandleFunc注册一个处理器函数handler和对应的模式pattern(注册到DefaultServeMux)；ServeMux的文档解释了模式的匹配机制
	func HandleFunc(pattern string, handler func(ResponseWriter, *Request))

	// Serve 会接手监听器l收到的每一个连接，并为每一个连接创建一个新的服务go程
	// 该go程会读取请求，然后调用handler回复请求；handler参数一般会设为nil，此时会使用DefaultServeMux
	func Serve(l net.Listener, handler Handler) error
	// ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接
	// handler参数一般会设为nil，此时会使用DefaultServeMux
	// For example
	// // // hello world, the web server
	// // func HelloServer(w http.ResponseWriter, req *http.Request) {
	// // 	io.WriteString(w, "hello, world!\n")
	// // }
	// // http.HandleFunc("/hello", HelloServer)
	// // err := http.ListenAndServe(":12345", nil)
	// // if err != nil {
	// // 	log.Fatal("ListenAndServe: ", err)
	// // }
	// // 
	func ListenAndServe(addr string, handler Handler) error
	// ListenAndServeTLS 和ListenAndServe 的行为基本一致，除了它期望HTTPS连接之外
	// 此外，必须提供证书文件和对应的私钥文件；如果证书是由权威机构签发的，certFile参数必须是顺序串联的服务端证书和CA证书；如果srv.Addr为空字符串，会使用":https"
	// 使用crypto/tls包的generate_cert.go文件来生成cert.pem和key.pem两个文件
	// For example
	// // func handler(w http.ResponseWriter, req *http.Request) {
	// // 	w.Header().Set("Content-Type", "text/plain")
	// // 	w.Write([]byte("This is an example server.\n"))
	// // }
	// // http.HandleFunc("/", handler)
	// // log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	// // err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error
```

### 3. net/http包的其他类型
- net/http/cookiejar
	- cookiejar包实现了保管在内存中的符合RFC 6265标准的http.CookieJar接口
```golang
	// PublicSuffixList 提供域名的公共后缀，例如：
	// - "example.com"的公共后缀是"com"
	// - "foo1.foo2.foo3.co.uk"的公共后缀是"co.uk"
	// - "bar.pvt.k12.ma.us"的公共后缀是"pvt.k12.ma.us"
	// PublicSuffixList接口的实现必须是并发安全的；一个总是返回""的实现是合法的，也可以通过测试；但却是不安全的：它允许HTTP服务端跨域名设置cookie
	// 推荐实现：code.google.com/p/go.net/publicsuffix
	type PublicSuffixList interface {
		// 返回域名的公共后缀。
		// TODO：域名的格式化应该由调用者还是接口方法负责还没有确定。
		PublicSuffix(domain string) string
		// 返回公共后缀列表的来源的说明，该说明一般应该包含时间戳和版本号。
		String() string
	}

	// Options是创建新Jar是的选项
	type Options struct {
		// PublicSuffixList是公共后缀列表，用于决定HTTP服务端是否能给某域名设置cookie
		// nil值合法的，也可以通过测试；但却是不安全的：它允许HTTP服务端跨域名设置cookie
		PublicSuffixList PublicSuffixList
	}
	// Jar类型实现了net/http包的http.CookieJar接口
	type Jar struct { ... }
	// 返回一个新的Jar，nil指针等价于Options零值的指针
	func New(o *Options) (*Jar, error)
	// 实现CookieJar接口的Cookies方法，如果URL协议不是HTTP/HTTPS会返回空切片
	func (j *Jar) Cookies(u *url.URL) (cookies []*http.Cookie)
	// 实现CookieJar接口的SetCookies方法，如果URL协议不是HTTP/HTTPS则不会有实际操作
	func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie)
```

- net/http/cgi
	- cgi包实现了CGI(Common Gateway Interface，公共网关协议)，参见RFC 3875
	- 注意使用CGI意味着对每一个请求开始一个新的进程，这显然要比使用长期运行的服务程序要低效
	- 本包主要是为了兼容现有的系统

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

### 4. net/mail
- net/smtp
	- smtp包实现了简单邮件传输协议(SMTP)

- net/mail
	- mail包实现了邮件的解析，大部分都遵守RFC 5322规定的语法
		- 旧格式地址和嵌入远端信息的地址不会被解析
		- 组地址不会被解析
		- 不支持全部的间隔符(CFWS语法元素)，如分属两行的地址

### 5. net/rpc
- net/rpc/jsonrpc
	- rpc包提供了通过网络或其他I/O连接对一个对象的导出方法的访问

- net/rpc/jsonrpc
	- jsonrpc包实现了JSON-RPC的ClientCodec和ServerCodec接口，可用于rpc包

### 6. net/other
- net/textproto
	- textproto实现了对基于文本的请求/回复协议的一般性支持，包括HTTP、NNTP和SMTP
	- 本包提供
		- 错误，代表服务端回复的错误码
		- Pipeline，以管理客户端中的管道化的请求/回复
		- Reader，读取数值回复码行，键值对形式的头域，一个作为后续行先导的空行，以及以只有一个"."的一行为结尾的整个文本块
		- Writer，写入点编码的文本
		- Conn，对Reader、Writer和Pipline的易用的包装，用于单个网络连接

- net/url
	- url包解析URL并实现了查询的逸码
```golang
	// URL类型代表一个解析后的URL(或者说，一个URL参照)
	// URL基本格式：scheme://[userinfo@]host/path[?query][#fragment]
	// scheme后不是冒号加双斜线的URL被解释的格式：scheme:opaque[?query][#fragment]
	// For example
	// // u, err := url.Parse("http://bing.com/search?q=dotnet")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // u.Scheme = "https"
	// // u.Host = "google.com"
	// // q := u.Query()
	// // q.Set("q", "golang")
	// // u.RawQuery = q.Encode()
	// // fmt.Println(u)
	// OutPut：https://google.com/search?q=golang
	type URL struct {
		Scheme   string
		Opaque   string    // 编码后的不透明数据
		User     *Userinfo // 用户名和密码信息
		Host     string    // host或host:port
		Path     string
		RawQuery string // 编码后的查询字符串，没有'?'
		Fragment string // 引用的片段(文档位置)，没有'#'
	}

	// Parse 解析rawurl为一个URL结构体，rawurl可以是绝对地址，也可以是相对地址
	func Parse(rawurl string) (url *URL, err error)
	// ParseRequestURI 解析rawurl为一个URL结构体，本函数会假设rawurl是在一个HTTP请求里，因此会假设该参数是一个绝对URL或者绝对路径，并会假设该URL没有#fragment后缀(网页浏览器会在去掉该后缀后才将网址发送到网页服务器)
	func ParseRequestURI(rawurl string) (url *URL, err error)

	// IsAbs 在URL是绝对URL时才返回真
	func (u *URL) IsAbs() bool
	// Query 解析RawQuery字段并返回其表示的Values类型键值对
	func (u *URL) Query() Values
	// RequestURI 返回编码好的path?query或opaque?query字符串，用在HTTP请求里
	func (u *URL) RequestURI() string
	// String 将URL重构为一个合法URL字符串
	func (u *URL) String() string
	// Parse 以u为上下文来解析一个URL，ref可以是绝对或相对URL；本方法解析失败会返回nil, err；否则返回结果和ResolveReference一致
	func (u *URL) Parse(ref string) (*URL, error)
	// ResolveReference 根据一个绝对URI将一个URI补全为一个绝对URI，参见RFC 3986 节 5.2
	// 参数ref可以是绝对URI或者相对URI。ResolveReference总是返回一个新的URL实例，即使该实例和u或者ref完全一样。如果ref是绝对URI，本方法会忽略参照URI并返回ref的一个拷贝
	func (u *URL) ResolveReference(ref *URL) *URL
	// Userinfo 类型是一个URL的用户名和密码细节的一个不可修改的封装。一个真实存在的Userinfo值必须保证有用户名(但根据 RFC 2396可以是空字符串)以及一个可选的密码
	type Userinfo struct { ... }
	// User 返回一个用户名设置为username的不设置密码的*Userinfo
	func User(username string) *Userinfo
	// UserPassword 返回一个用户名设置为username、密码设置为password的*Userinfo；这个函数应该只用于老式的站点，因为风险很大，不建议使用，参见RFC 2396
	func UserPassword(username, password string) *Userinfo
	// Username 方法返回用户名
	func (u *Userinfo) Username() string
	// 如果设置了密码返回密码和真，否则会返回假
	func (u *Userinfo) Password() (string, bool)
	// String 返回编码后的用户信息，格式为"username[:password]"
	func (u *Userinfo) String() string

	// Values将建映射到值的列表；它一般用于查询的参数和表单的属性。不同于http.Header这个字典类型，Values的键是大小写敏感的
	type Values map[string][]string

	// ParseQuery 解析一个URL编码的查询字符串，并返回可以表示该查询的Values类型的字典
	// 本函数总是返回一个包含了所有合法查询参数的非nil字典，err用来描述解码时遇到的(如果有)第一个错误
	func ParseQuery(query string) (m Values, err error)

	// Get 会获取key对应的值集的第一个值
	// 如果没有对应key的值集会返回空字符串获取值集请直接用map
	func (v Values) Get(key string) string
	// Set 将key对应的值集设为只有value，它会替换掉已有的值集
	func (v Values) Set(key, value string)
	// Add 将value添加到key关联的值集里原有的值的后面
	func (v Values) Add(key, value string)
	// Del 删除key关联的值集
	func (v Values) Del(key string)
	// Encode方法将v编码为url编码格式("bar=baz&foo=quux")，编码时会以键进行排序
	func (v Values) Encode() string
```
