package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	name string
	age  uint
)

func ScanfFromBufio() {
	fmt.Print("请输入你的姓名和年龄，以空格分隔：")
	stdin := bufio.NewReader(os.Stdin)
	line, _, err := stdin.ReadLine()
	if err != nil {
		panic(err)
	}

	n, err := fmt.Sscanln(string(line), &name, &age)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read number of items: %d", n)
	fmt.Println()
	
	fmt.Printf("姓名:%s 年龄:%d", name, age)
	fmt.Println()
}

func main() {
	ScanfFromBufio()
}
