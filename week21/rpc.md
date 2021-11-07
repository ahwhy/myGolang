# RPC入门

RPC是远程过程调用的简称，是分布式系统中不同节点间流行的通信方式。在互联网时代，RPC已经和IPC一样成为一个不可或缺的基础构件。

![](./images/rpc-model.png)

+ RPC传输协议
+ 消息序列化与反序列化

下面是一个基于HTTP的 JSON的 RPC:

![](./images/rpc-http-model.png)


## Go语言RPC

Go语言的标准库也提供了一个简单的RPC实现, 包的路径为net/rpc，也就是放在了net包目录下面。因此我们可以猜测该RPC包是建立在net包基础之上的

因此是之前着2中知识的集合:
+ Socket编程
+ 序列化(json/xml/...)

## "Hello, World"

由上图可知一个rpc服务由2个部分组成:
+ server
+ client

我们分别建2个项目

### RPC Server

下面是HelloService, 提供了一个 Hello 方法
```go
type HelloService struct {}

// Hello的逻辑 就是 将对方发送的消息前面添加一个Hello 然后返还给对方
// 由于我们是一个rpc服务, 因此参数上面还是有约束：
// 		第一个参数是请求
// 		第二个参数是响应
// 可以类比Http handler
func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}
```

如何把这个Hello方法, 变成一个RPC，直接供客户端调用喃?

Restful接口, 我们已经很熟悉了, 这也是一种RPC: JSON + HTTP

我们使用net/rpc包, 来实现一个rpc server:
```go
func main() {
	// 把我们的对象注册成一个rpc的 receiver
	// 其中rpc.Register函数调用会将对象类型中所有满足RPC规则的对象方法注册为RPC函数，
	// 所有注册的方法会放在“HelloService”服务空间之下
	rpc.RegisterName("HelloService", new(HelloService))

	// 然后我们建立一个唯一的TCP链接，
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务。
	// 没Accept一个请求，就创建一个goroutie进行处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// 前面都是tcp的知识, 到这个RPC就接管了
		// 因此 你可以认为 rpc 帮我们封装消息到函数调用的这个逻辑,
		// 提升了工作效率, 逻辑比较简洁，可以看看他代码
		go rpc.ServeConn(conn)
	}
}
```

### RPC Client

客户端如何调用我们server的Hello函数喃?

```go
func main() {
	// 首先是通过rpc.Dial拨号RPC服务, 建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 然后通过client.Call调用具体的RPC方法
	// 在调用client.Call时:
	// 		第一个参数是用点号链接的RPC服务名字和方法名字，
	// 		第二个参数是 请求参数
	//      第三个是请求响应, 必须是一个指针, 有底层rpc服务帮你赋值
	var reply string
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
```

rpc服务最多的优点就是 我们可以像使用本地函数一样使用 远程服务上的函数, 因此有几个关键点:

+ 远程连接: 类似于我们的pkg
+ 函数名称: 要表用的函数名称
+ 函数参数: 这个需要符合RPC服务的调用签名, 及第一个参数是请求，第二个参数是响应
+ 函数返回: rpc函数的返回是 连接异常信息, 真正的业务Response不能作为返回值

### 测试

1. 启动服务端
```sh
go run rpc/server/main.go
```

2. 客户端请求服务端
```sh
$ go run rpc/client/main.go
hello:hello
```

一个基础的rpc服务就是这么简单

## 基于接口的RPC服务

```go
// Call invokes the named function, waits for it to complete, and returns its error status.
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
```
上面是client call 方法, 里面3个参数2个interface{}, 你再使用的时候 可能真不知道要传入什么,
这就好像你写了一个HTTP的服务, 没有接口文档, 容易调用错误

如何避免这种情况喃? 我们可以对客户端进行一次封装, 使用接口当我们的 文档, 明确参数类型

定义hello service的接口
```go
package service

const HelloServiceName = "HelloService"

type HelloService interface {
	Hello(request string, reply *string) error
}
```


约束服务端:
```go
// 通过接口约束HelloService服务
var _ service.HelloService = (*HelloService)(nil)
```

封装客户端, 让其满足HelloService接口约束
```go
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
```

基于接口约束后的客户端使用就要容易很多了:
```go
func main() {
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
}
```

进行验证

着和我们之前写的pkg下的服务是不是很像, 一切都是有备而来


## gob编码

标准库的RPC默认采用Go语言特有的gob编码, 标准库gob是golang提供的“私有”的编解码方式，它的效率会比json，xml等更高，特别适合在Go语言程序间传递数据
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

gob的使用很简单, 和之前使用base64编码理念一样, 有 Encoder和Decoder
```go
func GobEncode(val interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(val); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func GobDecode(data []byte, value interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(value)
}
```

写个测试用例测测
```go
func TestGobCode(t *testing.T) {
	t1 := &TestStruct{"name", "value"}
	resp, err := service.GobEncode(t1)
	fmt.Println(resp, err)

	t2 := &TestStruct{}
	service.GobDecode(resp, t2)
	fmt.Println(t2, err)
}
```

## Json ON TCP

gob是golang提供的“私有”的编解码方式，因此从其它语言调用Go语言实现的RPC服务将比较困难

因此我们可以选用所有语言都支持的比较好的一些编码:
+ MessagePack: 高效的二进制序列化格式。它允许你在多种语言(如JSON)之间交换数据。但它更快更小
+ JSON: 文本编码
+ XML：文本编码
+ Protobuf 二进制编码

Go语言的RPC框架有两个比较有特色的设计：
+ RPC数据打包时可以通过插件实现自定义的编码和解码；
+ RPC建立在抽象的io.ReadWriteCloser接口之上的，我们可以将RPC架设在不同的通讯协议之上。

这里我们将尝试通过官方自带的net/rpc/jsonrpc扩展实现一个跨语言的RPC。

服务端:
```go
func main() {
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
}
```

客户端:

```go
func DialHelloService(network, address string) (*HelloServiceClient, error) {
	// c, err := rpc.Dial(network, address)
	// if err != nil {
	// 	return nil, err
	// }

	// 建立链接
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	// 采用Json编解码的客户端
	c := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{Client: c}, nil
}
```

验证功能是否正常

由于没有合适的tcp工具, 比如nc, 同学可以下来自己验证
```sh
$ echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
{"id":1,"result":"hello:hello","error":null}
```

## Json ON HTTP

Go语言内在的RPC框架已经支持在Http协议上提供RPC服务, 为了支持跨语言，编码我们依然使用Json

新的RPC服务其实是一个类似REST规范的接口，接收请求并采用相应处理流程

首先我们依然要解决JSON编解码的问题, 我们需要将HTTP接口的Handler参数传递给jsonrpc, 因此需要满足jsonrpc接口,
因此我们需要提前构建也给conn io.ReadWriteCloser, writer现成的 reader就是request的body, 直接内嵌就可以
```go
func NewRPCReadWriteCloserFromHTTP(w http.ResponseWriter, r *http.Request) *RPCReadWriteCloser {
	return &RPCReadWriteCloser{w, r.Body}
}

type RPCReadWriteCloser struct {
	io.Writer
	io.ReadCloser
}
```

服务端:
```go
func main() {
	rpc.RegisterName("HelloService", new(HelloService))

	// RPC的服务架设在“/jsonrpc”路径，
	// 在处理函数中基于http.ResponseWriter和http.Request类型的参数构造一个io.ReadWriteCloser类型的conn通道。
	// 然后基于conn构建针对服务端的json编码解码器。
	// 最后通过rpc.ServeRequest函数为每次请求处理一次RPC方法调用
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		conn := NewRPCReadWriteCloserFromHTTP(w, r)
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":1234", nil)
}
```

![](./images/http-json-rpc.png)


这种用法常见于你的rpc服务需要暴露多种协议的时候, 其他时候还是老老实实写Restful API

## 参考

+ [什么是RPC](https://www.jianshu.com/p/7d6853140e13)