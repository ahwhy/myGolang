# 简易版文件中转站

大家可能会遇到这种场景: 
  + 容器里面的日志帮忙拉下来下,我要定位下问题
  + 我有一个安装包下载不下来，我需要上传到服务器

如何解决这个问题:

+ 自己搭建一个文件服务器, 你帮他上传
+ 然用户使用系统上的工具自己做一个文件服务器
```
1. nohup python -m SimpleHTTPServer [port] & 快速搭建一个http服务
2. nc ...
```

## 解决方案

由于可能面临复杂网络下(多层跳板机)的上传和下载问题, 通过点对点传(比如scp之类)的很难行得通，所有选择用中转站的方式，比如oss

## 寻找现成工具
我们寻找有没有现有的客户端:
  + 阿里有oss browser, cli无法使用
  + 多个云商客户端不能通用
  
## 自己做一个简单的工具

1. 查看[阿里云 oss sdk使用样例](https://github.com/aliyun/aliyun-oss-go-sdk)

```go
client, err := oss.New("Endpoint", "AccessKeyId", "AccessKeySecret")
if err != nil {
    // HandleError(err)
}

bucket, err := client.Bucket("my-bucket")
if err != nil {
    // HandleError(err)
}

err = bucket.PutObjectFromFile("my-object", "LocalFile")
if err != nil {
    // HandleError(err)
}
```

2. 准备好你的阿里云bucket资源


3. 开始准备开发一个项目:
```
1. 初始化一个go mod工程 
go mod init "gitee.com/infraboard/go-course/day8/simple"
2. 编写main.go 引入sdk, 将基础样例封装为一个函数
3. 下载项目依赖
go mod tidy
```

4. 抽象变量与核心函数
```go
var (
	endpint    = ""
	acessKey   = ""
	secretKey  = ""
	bucketName = ""
	uploadFile = ""
)

func upload(filePaht string) error {
	client, err := oss.New(endpint, acessKey, secretKey)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(filePaht, filePaht)
	if err != nil {
		return err
	}

	return nil
}
```

5. 串联逻辑，编写入口函数

```go
func main() {
	if err := upload(uploadFile); err != nil {
		fmt.Printf("upload file error, %s", err)
		os.Exit(1)
	}
}
```

6. 调试: 直接运行会报错
```go
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
github.com/aliyun/aliyun-oss-go-sdk/oss.(*urlMaker).Init(0xc0001ba300, 0x0, 0x0, 0x1914c280000, 0x60, 0xc000024060)
	E:/Golang/pkg/mod/github.com/aliyun/aliyun-oss-go-sdk@v2.1.9+incompatible/oss/conn.go:722 +0x40c
github.com/aliyun/aliyun-oss-go-sdk/oss.New(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe1ff40, ...)
	E:/Golang/pkg/mod/github.com/aliyun/aliyun-oss-go-sdk@v2.1.9+incompatible/oss/client.go:52 +0xf2
main.upload(0x0, 0x0, 0x0, 0xc00004bf78)
	e:/Projects/Golang/go-course/day8/simple/simple.go:26 +0x93
main.main()
	e:/Projects/Golang/go-course/day8/simple/simple.go:19 +0x45
exit status 2
```

7. 校验程序入参
```go
func main() {
	if err := validate(); err != nil {
		fmt.Printf("validate paras error, %s\n", err)
		os.Exit(1)
	}

	if err := upload(uploadFile); err != nil {
		fmt.Printf("upload file error, %s\n", err)
		os.Exit(1)
	}
}

func validate() error {
	if endpint == "" {
		return fmt.Errorf("endpoint missed")
	}

	if acessKey == "" || secretKey == "" {
		return fmt.Errorf("access key or secret key missed")
	}

    if bucketName == "" {
		return fmt.Errorf("bucket name missed")
	}

	return nil
}
```

再次执行 验证参数检测逻辑是否生效
```go
validate paras error, access key or secret key missed
exit status 1
```

8. 传入测试参数后再次调试
```go
var (
	endpint    = "oss-cn-chengdu.aliyuncs.com"
	acessKey   = "ak"
	secretKey  = "sk"
	bucketName = "bucket"
	uploadFile = ""
)
```

结果如下:
```
upload file error, open : The system cannot find the file specified.
exit status 1
```

9. 配置好上传文件 再次调试
```go
var (
	endpint    = "oss-cn-chengdu.aliyuncs.com"
	acessKey   = "ak"
	secretKey  = "sk"
	bucketName = "cloud-station"
	uploadFile = "go.sum"
)
```

```
upload file error, oss: service returned error: StatusCode=403, ErrorCode=InvalidAccessKeyId, ErrorMessage="The OSS Access Key Id you provided does not exist in our records.", RequestId=60FB893504ACC03439856321
exit status 1
```

10. 配置上正确的key 再次调试
```go
var (
    endpint    = "oss-cn-chengdu.aliyuncs.com"
	acessKey   = "LTAI5tMvG5NA51eiH3ENZDaa"
	secretKey  = "找老师要"
	bucketName = "cloud-station"
	uploadFile = "go.sum"
)
```

```
upload file error, oss: service returned error: StatusCode=403, ErrorCode=AccessDenied, ErrorMessage="You have no right to access this object because of bucket acl.", RequestId=60FB89ACD9A90238360BCC5E
exit status 1
```
阿里云上配置该key的授权

11. 补充打印信息 再次调试
```go
upload file go.sum success
```

到此我们的程序可以正常运行，但是用户无法指定自己上传的文件

12. 为你的程序添加cli

使用flag从命令行接收用户输入的参数
```go
flag.StringVar(&uploadFile, "f", "", "指定本地文件")
```

定义一个函数用于读取用户从命令行传入的参数
```go
func loadParam() {
	flag.Parse()
}
```

这样就正常使用了
```go
$ go run simple.go -f go.mod 
upload file go.mod success
```

13. 为你的cli添加使用说明
```go
flag.BoolVar(&help, "h", false, "this help")
```
编写使用说明
```go
func usage() {
	fmt.Fprintf(os.Stderr, `cloud-station version: 0.0.1
Usage: cloud-station [-h] -f <uplaod_file_path>
Options:
`)
	flag.PrintDefaults()
}
```
最后解析参数的时候补充使用说明逻辑
```go
func loadParam() {
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}
}
```
测试
```
$ go run simple.go -h
cloud-station version: 0.0.1
Usage: cloud-station [-h] -f <uplaod_file_path>
Options:
  -h    this help
  -upload_file string
        指定本地文件
```

14. 打印下载链接
```go
signedURL, err := bucket.SignURL(filePaht, oss.HTTPGet, 60*60*24)
if err != nil {
    return fmt.Errorf("SignURL error, %s", err)
}
fmt.Printf("下载链接: %s\n", signedURL)
fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")
```

调试
```
下载链接: http://cloud-station.oss-cn-chengdu.aliyuncs.com/go.mod?Expires=1627186396&OSSAccessKeyId=LTAI5tMvG5NA51eiH3ENZDaa&Signature=TdwgxJVCOsmf84fqix8FKPmSLJc%3D

注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载
upload file go.mod success
```

15. 客户如何使用你提供的工具

+ 提供多个平台的版本(Linux, Windows, MacOS)
+ 用户如何获取到你的工具, 上传到一个工具仓库，提供下载链接
+ 为你的工具添加一个使用文档

## 来自用户的抱怨

+ 我上传大文件 没进度条, 不知道啥时候能传输完成
+ 没有传输速率, 不知道上传快慢

## 总结

我们以写脚本的一个模式写了一个工具, 他能最快的解决我们问题, 但是其中还是有一些简单的技巧

+ 面向过程开发的思维模式
+ 合理抽象函数
+ 合理使用变量
+ 注意校验用户输入

## 新的问题与需求

+ 中转站建在一个区域导致 其他区域必须外网访问速度慢，如何解决
+ API key 如何防止泄露
+ 如何避免业务方敏感文件泄密,是否需要控制下载时间段
+ 安全团队需要申请上传的东西是否有安全隐患
+ 安全团队需要确认下载的文件是否属于客户敏感数据
+ ...











