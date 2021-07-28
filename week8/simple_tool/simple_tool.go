package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	endpint    = "oss-cn-chengdu.aliyuncs.com"
	acessKey   = "LTAI5tPdJRm5driecxQFt6tK"
	secretKey  = "edsFTXwFh4bkodhiTdBF4mPVGSJaPb"
	bucketName = "cloud-station"
	uploadFile = "test.txt"
)

var (
	help bool
)

// 定义init函数
func init() {
	// 为cli添加使用说明
	flag.BoolVar(&help, "h", false, "this is help")
	// 使用flag从命令行接收用户输入的参数
	flag.StringVar(&uploadFile, "f", "test.txt", "指定本地文件路径")
}

// 使用说明
func usage() {
	fmt.Fprintf(os.Stderr, `cloud-station version: 0.0.1
Usage: cloud-station [-h] -f <uplaod_file_path>
Options:
`)
	flag.PrintDefaults()
}

// 定义loadParam函数 读取并分析用户从命令行传入的参数
func loadParam() {
	flag.Parse()

	// 使用说明逻辑
	if help {
		usage()
		os.Exit(0)
	}
}

// 定义validate函数 校验程序入参
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

// 阿里云SDK 抽象变量与核心函数
func upload(filePath string) error {
	client, err := oss.New(endpint, acessKey, secretKey)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(filePath, filePath)
	if err != nil {
		return err
	}

	signedURL, err := bucket.SignURL(filePath, oss.HTTPGet, 60*60*24)
	if err != nil {
		return fmt.Errorf("SignURL error, %s", err)
	}
	fmt.Printf("下载链接: %s\n", signedURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")

	return nil
}

func main() {
	loadParam()

	if err := validate(); err != nil {
		fmt.Printf("validate paras error, %s\n", err)
		os.Exit(1)
	}

	if err := upload(uploadFile); err != nil {
		fmt.Printf("upload file error, %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("upload file %s success\n", uploadFile)
}
