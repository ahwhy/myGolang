# Golang-gRPC  Golang的gRPC

## 一、Go语言中的RPC

### 1. RPC定义
- RPC是指远程过程调用，简单的理解是一个节点请求另一个节点提供的服务
	- 本地过程调用
		- 如果需要将本地student对象的age+1，可以实现一个addAge()方法，将student对象传入，对年龄进行更新之后返回即可，本地方法调用的函数体通过函数指针来指定
	- 远程过程调用
		- 上述操作的过程中，如果addAge()这个方法在服务端，执行函数的函数体在远程机器上
		- 首先客户端需要告诉服务器，需要调用的函数，这里函数和进程ID存在一个映射，客户端远程调用时，需要查一下函数，找到对应的ID，然后执行函数的代码
		- 客户端需要把本地参数传给远程函数，本地调用的过程中，直接压栈即可，但是在远程调用过程中不再同一个内存里，无法直接传递函数的参数，因此需要客户端把参数转换成字节流，传给服务端，然后服务端将字节流转换成自身能读取的格式，是一个序列化和反序列化的过程
		- 数据准备好了之后，网络传输层需要把调用的ID和序列化后的参数传给服务端，然后把计算好的结果序列化传给客户端，因此TCP层即可完成上述过程，gRPC中采用的是HTTP2协议

- 具体调用流程
```
	// Client端 
	// Student student = Call(ServerAddr, addAge, student)
	1. 将这个调用映射为Call ID
	2. 将Call ID，student（params）序列化，以二进制形式打包
	3. 把2中得到的数据包发送给ServerAddr，这需要使用网络传输层
	4. 等待服务器返回结果
	5. 如果服务器调用成功，那么就将结果反序列化，并赋给student，年龄更新
	
	// Server端
	1. 在本地维护一个Call ID到函数指针的映射call_id_map，可以用Map<String, Method> callIdMap
	2. 等待客户端请求
	3. 得到一个请求后，将其数据包反序列化，得到Call ID
	4. 通过在callIdMap中查找，得到相应的函数指针
	5. 将student（params）反序列化后，在本地调用addAge()函数，得到结果
	6. 将student结果序列化后通过网络返回给Client
```

- Call ID映射
	- 在本地调用中，函数体是直接通过函数指针来指定的，当调用函数时，编译器就自动调用它相应的函数指针
	- 但是在远程调用中，函数指针是不行的，因为两个进程的地址空间是完全不一样的
	- 在RPC中，所有的函数都必须有自己的一个ID，这个ID在所有进程中都是唯一确定的
	- 客户端在做远程过程调用时，必须附上这个ID
	- 然后还需要在客户端和服务端分别维护一个 {函数 <--> Call ID} 的对应表
	- 两者的表不一定需要完全相同，但相同的函数对应的Call ID必须相同
	- 当客户端需要进行远程调用时，查询该表找出相应的Call ID，然后把它传给服务端，服务端也通过查表，来确定客户端需要调用的函数，然后执行相应函数的代码

- 序列化和反序列化
	- 客户端需要参数值传给远程的函数
	- 在本地调用中，只需要把参数压到栈里，然后让函数自行去栈里读就行
	- 但是在远程过程调用时，客户端跟服务端是不同的进程，不能通过内存来传递参数
	- 甚至有时候客户端和服务端使用的都不是同一种语言(比如服务端用C++，客户端用Java或者Python)
	- 这时候就需要客户端把参数先转成一个字节流，传给服务端后，服务端再把字节流转成能读取的格式，这个过程叫序列化和反序列化
	- 同理，从服务端返回的值也需要序列化反序列化的过程

- 网络传输
	- 远程调用往往用在网络上，客户端和服务端是通过网络连接的
	- 所有的数据都需要通过网络传输，因此就需要有一个网络传输层
	- 网络传输层需要把Call ID和序列化后的参数字节流传给服务端，然后再把序列化后的调用结果传回客户端
	- 只要能完成这两者的，都可以作为传输层使用
	- 因此，它所使用的协议其实是不限的，能完成传输就行
	- 尽管大部分RPC框架都使用TCP协议，但其实UDP也可以，而gRPC干脆就用了HTTP2，同Java的Netty

- 而实现一个RPC框架，只需要完成以上三点的实现就可以完成基本的框架
	- Call ID映射可以直接使用函数字符串，也可以使用整数ID，映射表一般就是一个哈希表
	- 序列化反序列化可以自行编写，也可以使用Protobuf或者FlatBuffers
	- 网络传输可以使用socket，或者asio，ZeroMQ，Netty

### 2. Go语言中的RPC
- Go语言的标准库提供了一个简单的RPC实现 `net/rpc`

- RPC Server / Client
```go
	// RPC Server
	type HelloService struct {}

	// Hello的逻辑 将对方发送的消息前面添加一个Hello 然后返还给对方
	// 由于是一个rpc服务，参数上的约束为
	//   第一个参数是请求
	//   第二个参数是响应
	// 可以类比Http handler
	func (p *HelloService) Hello(request string, reply *string) error {
		*reply = "hello:" + request
		return nil
	}

	// 把HelloService对象注册成一个rpc的 receiver
	// 其中rpc.Register函数调用会将对象类型中所有满足RPC规则的对象方法注册为RPC函数
	// 所有注册的方法会放在"HelloService"服务空间之下
	rpc.RegisterName("HelloService", new(HelloService))

	// 然后建立一个唯一的TCP链接
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务
	// 每Accept一个请求，就创建一个goroutie进行处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		// 前面都是tcp的知识，这里之后就是RPC接管
		// 因此可以认为 rpc 封装消息到函数调用的这个逻辑
		// 提升了工作效率，逻辑比较简洁
		go rpc.ServeConn(conn)
	}

	// RPC Client
	// 首先通过rpc.Dial拨号RPC服务，建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 然后通过client.Call调用具体的RPC方法
	// 在调用client.Call时
	//   第一个参数是用点号链接的RPC服务名字和方法名字
	//   第二个参数是 请求参数
	//   第三个是请求响应，必须是一个指针，由底层rpc服务进行赋值
	var reply string
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
```

- 基于接口的RPC服务
	- `rpc.client.Call()`
	- 在client的call方法中，3个参数有2个interface{}，当使用的时候会有不知道要传入什么的情况产生
	- 可以对客户端进行一次封装，使用接口当作文档，明确参数类型

```go
	// rpc.client.Call()
	// Call invokes the named function, waits for it to complete, and returns its error status.
	func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
		call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
		return call.Error
	}
```

```go
	// 定义hello service的接口
	package service

	const HelloServiceName = "HelloService"

	type HelloService interface {
		Hello(request string, reply *string) error
	}
	
	// 约束服务端
	// 通过接口约束HelloService服务
	var _ service.HelloService = (*HelloService)(nil)
	
	// 封装客户端，让其满足HelloService接口约束
	// 约束客户端
	var _ service.HelloService = (*HelloServiceClient)(nil)

	type HelloServiceClient struct {
		*rpc.Client
	}

	func DialHelloService(network, address string) (*HelloServiceClient, error) {
		c, err := rpc.Dial(network, address)
		if err != nil {
			return nil, err
		}

		return &HelloServiceClient{Client: c}, nil
	}

	func (p *HelloServiceClient) Hello(request string, reply *string) error {
		return p.Client.Call(service.HelloServiceName+".Hello", request, reply)
	}

	client, err := DialHelloService("tcp", "localhost:1234")

	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string

	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
```

### 3. gob编码
- 标准库的RPC默认采用Go语言特有的gob编码，标准库gob是golang提供的"私有"的编解码方式，它的效率会比json，xml等更高，特别适合在Go语言程序间传递数据
```go
	// ServeConn runs the server on a single connection.
	// ServeConn blocks, serving the connection until the client hangs up.
	// The caller typically invokes ServeConn in a go statement.
	// ServeConn uses the gob wire format (see package gob) on the
	// connection. To use an alternate codec, use ServeCodec.
	// See NewClient's comment for information about concurrent access.
	func (server *Server) ServeConn(conn io.ReadWriteCloser) {
		buf := bufio.NewWriter(conn)
		srv := &gobServerCodec{
			rwc:    conn,
			dec:    gob.NewDecoder(conn),
			enc:    gob.NewEncoder(buf),
			encBuf: buf,
		}
		server.ServeCodec(srv)
	}
```

- gob的使用
```
	// GobEncode
	func GobEncode(val interface{}) ([]byte, error) {
		buf := bytes.NewBuffer([]byte{})
		encoder := gob.NewEncoder(buf)
		if err := encoder.Encode(val); err != nil {
			return []byte{}, err
		}
		return buf.Bytes(), nil
	}

	// GobDecode
	func GobDecode(data []byte, value interface{}) error {
		reader := bytes.NewReader(data)
		decoder := gob.NewDecoder(reader)
		return decoder.Decode(value)
	}

	// 测试用例
	func TestGobCode(t *testing.T) {
		t1 := &TestStruct{"name", "value"}
		resp, err := service.GobEncode(t1)
		fmt.Println(resp, err)
	
		t2 := &TestStruct{}
		service.GobDecode(resp, t2)
		fmt.Println(t2, err)
	}
```

### 4. Json ON TCP
- gob是golang提供的"私有"的编解码方式，因此从其它语言调用Go语言实现的RPC服务将比较困难
	- 可以选用所有语言都支持的比较好的一些编码
		- MessagePack: 高效的二进制序列化格式，它可以在多种语言(如JSON)之间交换数据，并且更快更小
		- JSON: 文本编码
		- XML: 文本编码
		- Protobuf 二进制编码
	- Go语言的RPC框架有两个比较有特色的设计
		- RPC数据打包时可以通过插件实现自定义的编码和解码
		- RPC建立在抽象的io.ReadWriteCloser接口之上的，可以将RPC架设在不同的通讯协议之上

- 使用Go语言标准库中的`net/rpc/jsonrpc` ，扩展实现一个跨语言的RPC
```go
	// Server
	rpc.RegisterName("HelloService", new(HelloService))
	
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
	
		// 代码中最大的变化是用rpc.ServeCodec函数替代了rpc.ServeConn函数，
		// 传入的参数是针对服务端的json编解码器
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
	
	// Client
	func DialHelloService(network, address string) (*HelloServiceClient, error) {
		// 建立链接
		conn, err := net.Dial(network, address)
		if err != nil {
			log.Fatal("net.Dial:", err)
			return nil, err
		}
	
		// 采用Json编解码的客户端
		c := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
		return &HelloServiceClient{Client: c}, nil
	}
	
	// $ echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
	// {"id":1,"result":"hello:hello","error":null}
```

- Json ON HTTP
	- Go语言中内在的RPC框架已经支持在Http协议上提供RPC服务，为了支持跨语言，编码依然需要使用Json
	- 新的RPC服务其实是一个类似REST规范的接口，接收请求并采用相应处理流程
	- 首先依然要解决JSON编解码的问题，需要将HTTP接口的Handler参数传递给jsonrpc
		- 满足jsonrpc接口，提前构建`io.ReadWriteCloser`类型的conn通道
		- Writer 是 ResponseWriter，ReadCloser 是 Request的Body，直接内嵌就可以
		- 服务端

```go
	func NewRPCReadWriteCloserFromHTTP(w http.ResponseWriter, r *http.Request) *RPCReadWriteCloser {
		return &RPCReadWriteCloser{w, r.Body}
	}
	
	type RPCReadWriteCloser struct {
		io.Writer
		io.ReadCloser
	}
```
		
```go
	rpc.RegisterName("HelloService", new(HelloService))
	
	// RPC的服务架设在"/jsonrpc"路径
	// 在处理函数中基于http.ResponseWriter和http.Request类型的参数构造一个io.ReadWriteCloser类型的conn通道
	// 然后基于conn构建针对服务端的json编码解码器
	// 最后通过rpc.ServeRequest函数为每次请求处理一次RPC方法调用
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		conn := NewRPCReadWriteCloserFromHTTP(w, r)
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	
	http.ListenAndServe(":1234", nil)

	// $ curl localhost:1234/jsonrpc -d'{"method":"HelloService.Hello","params":["hello"],"id":1}'
```

## 二、Go语言中的gRPC

### 1. gRPC简介
- gRPC是Google公司基于Protobuf开发的跨语言的开源RPC框架
	- gRPC基于HTTP/2协议设计，可以基于一个HTTP/2链接提供多个服务，对于移动设备更加友好

- gRPC技术栈
```
	// 数据交互格式: protobuf
	// 通信方式: 最底层为TCP或Unix Socket协议，在此之上是HTTP/2协议的实现
	// 核心库: 在HTTP/2协议之上又构建了针对Go语言的gRPC核心库
	// Stub: 应用程序通过gRPC插件生产的Stub代码和gRPC核心库通信，也可以直接和gRPC核心库通信
	Application
	Generated Stubs
	gRPC GO Core + Interceptors
	HTTP/2 (Security: TLS/SSL or ALTS, etc)
	Unix Domain Sockets; TCP
```

### 2. gRPC-demo
- 安装Protobuf
	- Protobuf grpc插件
	- 从Protobuf的角度看，gRPC只不过是一个针对service接口生成代码的生成器，需要提前安装grpc的代码生成插件
```shell
	# 安装Protobuf grpc插件
	# protoc-gen-go
	$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	
	# 安装protoc-gen-go-grpc插件
	$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$ protoc-gen-go-grpc --version                                   
	protoc-gen-go-grpc 1.1.0
```

- protobuf 定义接口的语法
```
	service <service_name> {
		rpc <function_name> (<request>) returns (<response>);
	}
	// service: 用于申明这是个服务的接口
	// service_name: 服务的名称,接口名称
	// function_name: 函数的名称
	// request: 函数参数， 必须的
	// response: 函数返回， 必须的, 不能没有
```

- grpc.proto
```proto
	syntax = "proto3";
	
	package hello;
	option go_package="gitee.com/infraboard/go-course/day21/pbrpc/service";
	
	// The HelloService service definition.
	service HelloService {
		rpc Hello (Request) returns (Response);
	}
	
	message Request {
		string value = 1;
	}
	
	message Response {
		string value = 1;
	}
	
	// $ protoc -I=. --go_out=./grpc/service --go_opt=module="gitee.com/infraboard/go-course/day21/grpc/service" \
	//  --go-grpc_out=./grpc/service --go-grpc_opt=module="gitee.com/infraboard/go-course/day21/grpc/service" \
	//  grpc/service/service.proto
```

- gRPC服务端
```go
	var _ service.HelloServiceServer = (*HelloService)(nil)
	
	type HelloService struct {
		service.UnimplementedHelloServiceServer
	}
	
	func (p *HelloService) Hello(ctx context.Context, req *service.Request) (*service.Response, error) {
		resp := &service.Response{}
		resp.Value = "hello:" + req.Value
		return resp, nil
	}
	
	// 首先是通过grpc.NewServer()构造一个gRPC服务对象
	grpcServer := grpc.NewServer()
	// 然后通过gRPC插件生成的RegisterHelloServiceServer函数注册实现的HelloServiceImpl服务
	service.RegisterHelloServiceServer(grpcServer, new(HelloService))
	
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	
	// 最后通过grpcServer.Serve(lis)在一个监听端口上提供gRPC服务
	grpcServer.Serve(lis)
```

- gRPC客户端
```go
	// grpc.Dial负责和gRPC服务建立链接
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	// NewHelloServiceClient函数基于已经建立的链接构造HelloServiceClient对象,
	// 返回的client其实是一个HelloServiceClient接口对象
	client := service.NewHelloServiceClient(conn)
	
	// 通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法
	req := &service.Request{Value: "hello"}
	reply, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
```

### 3. gRPC流
- RPC是远程函数调用，因此每次调用的函数参数和返回值不能太大，否则将严重影响每次调用的响应时间
	- 因此传统的RPC方法调用对于上传和下载较大数据量场景并不适合
	- 为此，gRPC框架针对服务器端和客户端分别提供了流特性

- 服务端或客户端的单向流是双向流的特例，在HelloService增加一个支持双向流的Channel方法
	- 定义streaming RPC 的语法
		- gRPC服务端
			- 逻辑
				- 接收一个Request
				- 响应一个Response
			- 双向流数据的发送和接收都是完全独立的行为
				- 需要注意的是，发送和接收的操作并不需要一一对应，用户可以根据真实场景进行组织代码
		- gRPC客户端
```go
	// 定义streaming RPC 的语法
	// The HelloService service definition.
	service HelloService {
		rpc Hello (Request) returns (Response) {}

		// rpc <function_name> (stream <type>) returns (stream <type>) {}
		// 关键字stream指定启用流特性，参数部分是接收客户端参数的流，返回值是返回给客户端的流
		rpc Channel (stream Request) returns (stream Response) {}
	}
	// $ protoc -I=. --go_out=./grpc/service --go_opt=module="gitee.com/infraboard/go-course/day21/grpc/service" --go-grpc_out=./grpc/service --go-grraboard/go-course/day21/grpc/service" grpc/service/service.proto

	// gRPC服务端
	func (p *HelloService) Channel(stream service.HelloService_ChannelServer) error {
		// 服务端在循环中接收客户端发来的数据
		for {
			// 接收一个请求
			args, err := stream.Recv()
			if err != nil {
				// 如果遇到io.EOF表示客户端流被关闭
				if err == io.EOF {
					return nil
				}
				return err
			}
	
			// 响应一个请求
			// 生成返回的数据通过流发送给客户端
			resp := &service.Response{Value: "hello:" + args.GetValue()}
			err = stream.Send(resp)
			if err != nil {
				// 服务端发送异常, 函数退出, 服务端流关闭
				return err
			}
		}
	}
	
	// gRPC客户端
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	client := service.NewHelloServiceClient(conn)
	
	// 客户端需要先调用Channel方法获取返回的流对象
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	
	// 在客户端我们将发送和接收操作放到两个独立的Goroutine。
	
	// 首先是向服务端发送数据
	go func() {
		for {
			if err := stream.Send(&service.Request{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()
	
	// 然后在循环中接收服务端返回的数据
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
```

### 4. 参考文档
- [GRPC Quick Start](https://grpc.io/docs/languages/go/quickstart/)

- [GRPC Examples](https://github.com/grpc/grpc-go/tree/master/examples)