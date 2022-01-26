# Golang-Network  Golang的网络

## 一、网络通信

- 网络请求流程图
```
	                                     内核空间                               用户空间
			1).网络请求 ->        2).copy(I/O模型、DMA)->        3).copy(MMAP) ->      4-1).处理请求
	client                   网卡                         内核缓冲区            web服务进程    |
			7).返回数据 <-               6).copy <-                  6).copy <-        4-2).构建Respense
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
	Server端        Socket() -> Bind() -> Listen() -> Accept() -> Receive() -> Send()    -> Close()
	                                                 建立连接    发送数据     接收数据     关闭连接
	Client端        Socket()           ->            Connect() -> Send()    -> Receive() -> Close()
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
	|                             序号                              |
	|                            确认号                             |
	| 数据偏移 | 保留位 |           tcp flags            |   窗口   | 
	|  (4bit)  | (6bit) | (URG、ACK、PSH、RST、SYN、FIN) | (2Bytes) |
	|      检验和(2Bytes)          |         紧急指针(2Bytes)       |
	|                           TCP选项                             |
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
			- 32 位序号: 也称为顺序号(Sequence Number)，简写为SEQ，
			- 32 位确认序号: 也称为应答号(Acknowledgment Number)，简写为ACK
				- 在握手阶段，确认序号将发送方的序号加1作为回答
		- 数据偏移
			- 4 位数据偏移／首部长度: TCP数据部分距TCP开头的偏移量(一个偏移量是4个字节)，亦即TCP首部的长度	·
			- 由于首部可能含有可选项内容，因此TCP报头的长度是不确定的
			- 报头不包含任何任选字段则长度为20字节，4位首部长度字段所能表示的最大值为1111，转化为10进制为15，`15*32/8 = 60`，故报头最大长度为60字节，TCP选项最多有40个字节
			- 首部长度也叫数据偏移，是因为首部长度实际上指示了数据区在报文段中的起始偏移值
		- 保留位: 为将来定义新的用途保留，现在一般置0
		- 6 位标志字段
			- URG: 紧急指针标志，为1时表示紧急指针有效，为0则忽略紧急指针
			- ACK: 确认序号标志，为1时表示确认号有效，为0表示报文中不含确认信息，忽略确认号字段；TCP协议规定在连接建立后所有传输的报文段都必须把ACK设置为1
			- PSH: push标志，为1表示是带有push标志的数据，指示接收方在接收到该报文段以后，应尽快将这个报文段交给应用程序，而不是在缓冲区排队
			- RST: 重置连接标志，用于重置由于主机崩溃或其他原因而出现错误的连接；或者用于拒绝非法的报文段和拒绝连接请求
			- SYN: 同步序号，用于建立连接过程，在连接请求中，SYN=1和ACK=0表示该数据段没有使用捎带的确认域；而连接应答捎带一个确认，即SYN=1和ACK=1。
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
			- 第三次握手: SYN=1，ACK=1，序号=J+1，确认号=K+1；此次一般会携带真正需要传输的数据
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
		- 客户端独有的: SYN_SENT 、FIN_WAIT1 、FIN_WAIT2 、CLOSING 、TIME_WAIT
		- 服务器独有的: LISTEN、SYN_RCVD 、CLOSE_WAIT、LAST_ACK
		- 共有的: CLOSED、ESTABLISHED
		- 描述
			- CLOSED: 起始点，tcp连接 超时或关闭时进入此状态
			- LISTEN: 服务端等待连接时的状态，调用Socket、bind、listen函数就能进入此状态，称为被动打开
			- SYN_SENT: 客户端发起连接，发送SYN给服务端，若不能连接进入CLOSED
			- SYN_RCVD: 服务端接受客户端的SYN，由LISTEN进入SYN_RCVD；同时返回一个ACK，并发送一个SYN给服务端
				- 特殊情况: 客户端发起SYN的同时接收到服务端的SYN，客户端会由SYN-sent转为SYN-rcvd状态
			- ESTABLISHED: 可以传输数据的状态
			- FIN_WAIT1: 主动关闭连接，发送FIN，由ESTABLISHED转为此状态
			- FIN_WAIT2: 主动关闭连接，接收到对方的FIN+ACK，由FIN_WAIT1转为此状态
			- CLOSE_WAIT: 收到FIN，发送ACK，被动关闭的一方关闭连接进入此状态
			- LAST_ACK: 发送FIN，同时在接受ACK时，由CLOSE_WAIT进入此状态；被关闭的一方发起关闭请求
			- CLOSING: 两边同时发送关闭请求，由FIN_WAIT1进入此状态，收到FIN请求，同时响应一个ACK
			- TIME_WAIT
				- FIN_WAIT2到此状态: 双方不同时发送FIN的情况下，主动关闭的一方在完成自己的请求后，收到对方的FIN后的状态
				- CLOSING到此状态: 双方同时发起关闭，都发送了FIN，同时接受FIN并发送ACK后的状态
				- FIN_WAIT2到此状态: 对方发来的FIN的ACK同时到达后的状态，与上一条的区别是本身发送的FIN回应的ACK先于对方的FIN到达，而上一条是FIN先到达

- Go语言中的TCP编程接口
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
		ssl_certificate "/path/to/server.crt";        // 填写刚生成的TLS证书（包含公钥）
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

### 3. Https
- Https
	- http:  (应用层)HTTP -> TCP -> IP
	- https: (应用层)HTTP -> SSL/TCP -> TCP -> IP
		- HTTP + 加密 + 认证 + 完整性保护 = HTTPS(HTTP Secure)

### 4. Http-Response
- Go语言中的标准库 `net/http`
```go
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

### 5. http-router
- Go语言中的第三方库 `github.com/julienschmidt/httprouter`
	- `go get -u github.com/julienschmidt/httprouter`
	- Router实现了`http.Handler`接口
	- 为各种request method提供了便捷的路由方式
	- 支持restful请求方式
	- 支持ServeFiles访问静态文件
	- 可以自定义捕获panic的方法

### 6. 关于请求校验的常见问题
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

### 7. http中间件
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

### 8. Go语言中的Http框架
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

```go
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
		c.TplName = "index.tpl”
		//Resquest和ResponseWriter都在beego.Controller.Ctx里
		fmt.Println("remote addr", c.Ctx.Request.RemoteAddr)
		//如果指定了response正文，就不会去渲染index.tpl了
		c.Ctx.WriteString("Hi boy") 
	}
``