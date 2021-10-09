package main

import (
	"fmt"
	"os"
)

func main() {
	var num int64
	var tag string
	array_Bytes := [4]string{"MB", "GB", "TB", "PB"}

	fmt.Print("请输入一个需要转换的B（字节）：")
	fmt.Scanln(&num)

	fmt.Printf("需要转换的值：%dB\n", num)
	num2 := float64(num)
	num2 = HumanBytesLoaded(num2)
	fmt.Printf("已转换的值：%.3fKB\n", num2)

	for _, j := range array_Bytes {
		fmt.Printf("是否需要继续转换%s，如果需要请输入y:", j)
		fmt.Scanln(&tag)
		Next(tag)
		num2 = HumanBytesLoaded(num2)
		fmt.Printf("已转换的值：%.3f%s \n", num2, j)
	}

	fmt.Println("已完成所有转换，感谢使用！")
}

func HumanBytesLoaded(bytesLength float64) float64 {
	resp := bytesLength / 1024
	
	return resp
}

func Next(tag string) {
	str := "y"

	if tag == str {
		fmt.Println("程序继续")
	} else {
		fmt.Println("程序终止，感谢使用！")
		os.Exit(1)
	}
}
