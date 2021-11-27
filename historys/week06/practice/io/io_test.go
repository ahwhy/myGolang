package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

// 实现一个reader 每次读取4个字节
func Test_Strings_NewReader(test *testing.T) {
	//
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

}

func Test_readFile(test *testing.T) {
	fileName := "test.txt"
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("[内容：%s]", bytes)
}

func Test_writeFile(test *testing.T) {
	fileName := "白日梦.txt"
	err := ioutil.WriteFile(fileName, []byte("升职加薪\n迎娶白富美"), 0644)
	fmt.Println(err)
}

func Test_readDir(test *testing.T) {
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
			f.ModTime())
	}
}
