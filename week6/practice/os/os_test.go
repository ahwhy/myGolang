package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

func Test_Create(test *testing.T) {
	path := "test.txt"
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	name := "aa"
	fmt.Fprintf(file, "I am %s\n", name)

	file.Write([]byte("123456789\n"))
	fmt.Println(file.Write([]byte("bb")))
	fmt.Println(file.Write([]byte("qwertyuiopasdfghjkl\n"))) // 两个返回值 一个是 []byte 的长度，一个是 error
	file.WriteString("ccc联动")
}

func Test_Open(test *testing.T) {
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
			// EOF(End Of File) -> 标识文件读取结束
			// 非EOF
			if err != io.EOF {
				fmt.Print(err)
			} else {
				fmt.Print(err)
			}
			break
		}
		fmt.Println(string(content[:n]))
	}
}

func Test_Std(test *testing.T) {
	// - 标准输入  -> os.Stdin
	// - 标准输出  -> os.Stdout
	// - 标准错误  -> os.Stderr
	content := make([]byte, 3)

	fmt.Print("请输入内容: ")
	fmt.Println(os.Stdin.Read(content))
	fmt.Printf("%q\n", string(content))

	os.Stdout.WriteString("我是Stdout的输出")
	fmt.Fprintln(os.Stdout, "aaaaa")
	fmt.Fprintf(os.Stdout, "I am: %s", "aaaaa")

	// fmt.Scan -> Scanln, Scanf
	// fmt.Sscan -> 从字符串扫描到变量
	// fmt.Fscan -> 从文件扫描到变量
}

func Test_OpenFile(test *testing.T) {
	file, err := os.OpenFile("test.txt", os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	content := make([]byte, 2)
	fmt.Println(file.Read(content))
	fmt.Println(content)
	fmt.Println(file.Write([]byte("123456789")))
	fmt.Println(file.Read(content))
	fmt.Println(content)

}

func ScanInt() (int, error) {
	// 读取一行 进行转换
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strconv.Atoi(scanner.Text())
	}
	return 0, scanner.Err()
}

func Test_Scan(test *testing.T) {
	// os.Stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		break
	}

	num, err := ScanInt()
	fmt.Println(num, err)

}

func Test_Other(test *testing.T) {
	log.Printf("[获取命令行参数][res:%v]", os.Args)
	hn, _ := os.Hostname()
	log.Printf("[获取主机名][res:%v]", hn)
	log.Printf("[获取当前进程名][res:%v]", os.Getpid())
	log.Printf("[获取一条环境变量][res:%v]", os.Getenv("GOROOT"))
	// 获取所有环境变量
	env := os.Environ()
	for _, v := range env {
		log.Printf("[获取所有环境变量][res:%v]", v)
	}
	dir, _ := os.Getwd()
	log.Printf("[获取当前目录][res:%v]", dir)

	_ = os.Mkdir("config", 0755)
	log.Printf("[创建单一目录config目录]")

	// mkdir -p
	os.MkdirAll("config1/yaml/local", 0755)
	log.Printf("[递归创建目录config1目录]")
	// rm dir
	err := os.Remove("config")
	log.Printf("[删除单一目录config1目录][err:%v]", err)
	// rm -rf
	//err = os.RemoveAll("config1")
	//log.Printf("[全部删除目录config1目录][err:%v]", err)
}

func Test_IOCopy(test *testing.T) {
	file, err := os.Open("test.txt")
	fmt.Println(err)
	defer file.Close()

	io.Copy(os.Stdout, file)
}