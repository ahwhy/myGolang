# Protobuf 

## 一、Protobuf介绍

### 1. Protobuf
 - Protobuf是Protocol Buffers的简称
 	- Google公司开发的一种数据描述语言，并于2008年对外开源
 	- Protobuf刚开源时的定位类似于XML、JSON等数据描述语言，通过附带工具生成代码并实现将结构化数据序列化的功能
 	- Protobuf作为接口规范的描述语言，可以作为设计安全的跨语言PRC接口的基础工具
 
 - Protobuf特性
	- 编解码效率
	- 高压缩比
	- 多语言支持

 - 数据编码
	- 在XML或JSON等数据描述语言中，一般通过成员的名字来绑定对应的数据
	- 但是Protobuf编码却是通过成员的唯一编号来绑定对应的数据，因此Protobuf编码后数据的体积会比较小，但是也非常不便于人类查阅
	- 目前并不关注Protobuf的编码技术，最终生成的Go结构体可以自由采用JSON或gob等编码格式，因此可以暂时忽略Protobuf的成员编码部分

### 2. Protobuf安装流程
- 需要定义数据，通过编译器，来生成不同语言的代码

- 安装编译器
	- protobuf的编译器叫: `protoc(protobuf compiler)`
	- 下载地址:  [Github Protobuf](https://github.com/protocolbuffers/protobuf/releases)
	- 安装编译器二进制
		- linux/unix系统 `$ mv bin/protoc usr/bin`
		- windows系统
			- 注意: Windows 上的 git-bash 上默认的 `/usr/bin` 目录在: `C:\Program Files\Git\usr\bin\`
			- 首先将`bin`下的 `protoc`编译器 移动到`C:\Program Files\Git\usr\bin\`
	- 安装编译器库
		- include 下的库文件需要安装到 `/usr/local/include/`
		- linux/unix系统 `$ mv include/google /usr/local/include`
		- windows系统 `C:\Program Files\Git\usr\local\include`
	- 验证安装
```shell
		$ protoc --version
		libprotoc 3.19.1
```
	- 安装Go语言插件
		- Protobuf核心的工具集是C++语言开发的，在官方的protoc编译器中并不支持Go语言
		- 要想基于 `*.proto` 文件生成相应的Go代码，需要安装相应的插件
```shell
			$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### 3. Protobuf编解码
- 在rpc中使用probuf
	- 定义 `*.proto`文件
```
		// hello.proto
		// syntax: 表示采用proto3的语法；第三版的Protobuf对语言进行了提炼简化，所有成员均采用类似Go语言中的零值初始化(不再支持自定义默认值)，因此消息成员也不再需要支持required特性
		// package: 指明当前是main包(这样可以和Go的包名保持一致，简化例子代码)，当然用户也可以针对不同的语言定制对应的包路径和名称
		// option: protobuf的一些选项参数, 这里指定的是要生成的Go语言package路径, 其他语言参数各不相同
		// message: 关键字定义一个新的String类型，在最终生成的Go语言代码中对应一个String结构体。String类型中只有一个字符串类型的value成员，该成员编码时用1编号代替名字
		syntax = "proto3";
		package hello;
		option go_package="gitee.com/infraboard/go-course/day21/pb";
		message String {
			string value = 1;
		}
```
	- 编译 `*.proto`文件
```
		$ cd day21 ## 以day21作为编译的工作目录
		$ protoc -I=. --go_out=./pb --go_opt=module="gitee.com/infraboard/go-course/day21/pb" pb/hello.proto
		// -I: -IPATH, --proto_path=PATH, 指定proto文件搜索的路径, protoc会根据该路径来搜索导入的包, 如果有多个路径 可以多次使用-I 来指定, 如果不指定默认为当前目录
		// --go_out: --go指插件的名称, 安装的插件为: protoc-gen-go, 而protoc-gen是插件命名规范, go是插件名称, 因此这里是--go, 而--go_out 表示的是 go插件的 out参数, 这里指编译产物的存放目录
		// --go_opt: protoc-gen-go插件opt参数, 这里的module指定了go module, 生成的go pkg 会去除掉module路径，生成对应pkg
		// pb/hello.proto: proto文件路径
```
	- 在当前目录下生成Go语言对应的pkg

- 序列化与反序列化
	- 使用`google.golang.org/protobuf/proto`工具提供的API来进行序列化与反序列化:
		- Marshal: 序列化
		- Unmarshal: 反序列化
	- 客户端 ---> 服务端 基于protobuf的数据交互过程
```go
		clientObj := &pb.String{Value: "hello proto3"}
		// 序列化
		out, err := proto.Marshal(clientObj)
		if err != nil {
			log.Fatalln("Failed to encode obj:", err)
		}
		// 二进制编码
		// encode bytes:  [10 12 104 101 108 108 111 32 112 114 111 116 111 51]
		fmt.Println("encode bytes: ", out)

		// 反序列化
		serverObj := &pb.String{}
		err = proto.Unmarshal(out, serverObj)
		if err != nil {
			log.Fatalln("Failed to decode obj:", err)
		}
		// decode obj:  value:"hello proto3"
		fmt.Println("decode obj: ", serverObj)
```

- 基于protobuf的RPC `Protobuf ON TCP`

## 二、Protobuf语法

- 定义消息类型
	- 示例
```
		syntax = "proto3";
		
		/* SearchRequest represents a search query, with pagination options to
		* indicate which results to include in the response. */
		
		message SearchRequest {
		string query = 1;
		int32 page_number = 2; // Which page number do we want?
		int32 result_per_page = 3;
		}
```
	- 模板
```
		// comment: 注释 /* */或者 //
		// message_name: 同一个pkg内，必须唯一
		// filed_rule: 可以没有, 常用的有repeated, oneof
		// filed_type: 数据类型, protobuf定义的数据类型, 生产代码的会映射成对应语言的数据类型
		// filed_name: 字段名称, 同一个message 内必须唯一
		// field_number: 字段的编号, 序列化成二进制数据时的字段编号, 同一个message 内必须唯一, 1 ~ 15 使用1个Byte表示, 16 ～ 2047 使用2个Byte表示
		<comment>
		
		message  <message_name> {
		  <filed_rule>  <filed_type> <filed_name> = <field_number> 
							类型         名称           编号  
		}
```
	- 如果想保留一个编号，以备后来使用可以使用 reserved 关键字声明
```
		message Foo {
		  reserved 2, 15, 9 to 11;
		  reserved "foo", "bar";
		}
```

- Value(Filed) Types
	- .proto: double、float、int32、int64、uint32、uint64、sint32、sint64、fixed32、fixed64、sfixed32、sfixed64、bool、string、bytes
	- Go Type: float64、float32、int32、int64、uint32、uint64、int32、int64、uint32、uint64、int32、int64、bool、string、[]byte

- 枚举类型
	- 使用enum来声明枚举类型
```
		enum Corpus {
			UNIVERSAL = 0;
			WEB = 1;
			IMAGES = 2;
			LOCAL = 3;
			NEWS = 4;
			PRODUCTS = 5;
			VIDEO = 6;
		}

		// enum_name: 枚举名称
		// element_name: pkg内全局唯一, 很重要
		// element_name: 必须从0开始, 0表示类型的默认值, 32-bit integer
		enum <enum_name> {
			<element_name> = <element_number>
		}
```

- 别名
	- 如果有2个同名的枚举需求: 比如 TaskStatus 和 PipelineStatus 都需要Running，就可以添加一个: `option allow_alias = true`;
```
		message MyMessage1 {
		  enum EnumAllowingAlias {
		    option allow_alias = true;
		    UNKNOWN = 0;
		    STARTED = 1;
		    RUNNING = 1;
		  }
		}
		
		message MyMessage2 {
		  enum EnumNotAllowingAlias {
		    UNKNOWN = 0;
		    STARTED = 1;
		    // RUNNING = 1;  // Uncommenting this line will cause a compile error inside Google and a warning message outside.
		  }
		}
```
	- 同理枚举也支持预留值
```
		enum Foo {
		  reserved 2, 15, 9 to 11, 40 to max;
		  reserved "FOO", "BAR";
		}
```

- 数组类型
	- 如果想声明: []string, []Item 数组类型
	- 可以使用 `filed_rule: repeated`
```
		message SearchResponse {
		  repeated Result results = 1;
		}
		
		// 会编译为:
		// type SearchResponse SearchResponse {
		//   results []*Result
		// }
```

- Map
	- protobuf 声明map
```
		map<string, Project> projects = 3;
		// projects map[string, Project]
		
		map<key_type, value_type> map_field = N;
```

- Oneof
	- 很像范型 比如 test_oneof 字段的类型 必须是 string name 和 SubMessage sub_message 其中之一:
```
	message Sub1 {
	  string name = 1;
	}
	
	message Sub2 {
	  string name = 1;
	}
	
	message SampleMessage {
	  oneof test_oneof {
	    Sub1 sub1 = 1;
	    Sub2 sub2 = 2;
	  }
	}
	
	// 生成的结构体
	type SampleMessage struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields
	
		// Types that are assignable to TestOneof:
		//	*SampleMessage_Sub1
		//	*SampleMessage_Sub2
		TestOneof isSampleMessage_TestOneof `protobuf_oneof:"test_oneof"`
	}

	// 使用
	of := &pb.SampleMessage{}
	of.GetSub1()
	of.GetSub2()
```

- Any
	- 当无法明确定义数据类型的时候，可以使用Any表示
```
		// 应用其他的proto文件
		import "google/protobuf/any.proto";
		
		message ErrorStatus {
		string message = 1;
		repeated google.protobuf.Any details = 2;
		}
		
		// any本质上就是一个bytes数据结构
		type ErrorStatus struct {
			state         protoimpl.MessageState
			sizeCache     protoimpl.SizeCache
			unknownFields protoimpl.UnknownFields
		
			Message string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
			Details []*anypb.Any `protobuf:"bytes,2,rep,name=details,proto3" json:"details,omitempty"`
		}

		// Any的定义
		type Any struct {
			state         protoimpl.MessageState
			sizeCache     protoimpl.SizeCache
			unknownFields protoimpl.UnknownFields
		
			...
			// Note: this functionality is not currently available in the official
			// protobuf release, and it is not used for type URLs beginning with
			// type.googleapis.com.
			//
			// Schemes other than `http`, `https` (or the empty scheme) might be
			// used with implementation specific semantics.
			//
			TypeUrl string `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
			// Must be a valid serialized protocol buffer of the above specified type.
			Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
		}
```

- 类型嵌套
	- 在message里面嵌套message
	- 与Go结构体嵌套一样，但是不允许 匿名嵌套，必须指定字段名称
```
		message Outer {                  // Level 0
		  message MiddleAA {  // Level 1
		    message Inner {   // Level 2
		      int64 ival = 1;
		      bool  booly = 2;
		    }
		}
		  message MiddleBB {  // Level 1
		    message Inner {   // Level 2
		      int32 ival = 1;
		      bool  booly = 2;
		    }
		  }
		}
```

- 引用包
	- `import "google/protobuf/any.proto";` 
	- 这是读取的标准库，在安装protoc的时候，已经把改lib 挪到usr/local/include下面了，所以可以直接找到
	- 如果导入的proto文件并没有在 /usr/local/include 目录下，通过-I 可以添加搜索的路径，这样就编译器就可以找到引入的包

- 更新的规范
	- `Don't change the field numbers for any existing fields.`
	- [Updating A Message Type](https://developers.google.com/protocol-buffers/docs/proto3#updating)

## 三、参考文档

- [Protocol Buffers](https://developers.google.com/protocol-buffers) Google官网 
- [Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
- [Protocol Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)