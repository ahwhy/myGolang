package main

import (
	"fmt"
	"time"
)

func main() {
	//var num1 int
	//fmt.Print("请输入一个需要转换的B（字节）：")
	//fmt.Scanln(&num1)
	//fmt.Printf("已转换的值：%sKB\n",HumanBytesLoaded1(num1))
	//var num2 float64
	//fmt.Print("请输入一个需要转换的B（字节）：")
	//fmt.Scanln(&num2)
	//HumanBytesLoaded2(num2)
	var num3 float64
	var tag string
	array_Bytes := [4]string{"MB", "GB", "TB", "PB"}

	fmt.Print("请输入一个需要转换的B（字节）：")
	fmt.Scanln(&num3)
	
	// if num3 <= 0 || num3 > int64(PB) {
	// 	fmt.Println("输入的值错误")
	// }

	fmt.Printf("需要转换的值：%gB\n", num3)
	num3 = HumanBytesLoaded3(num3)
	fmt.Printf("已转换的值：%gKB\n", num3)

	for _, j := range array_Bytes {
		fmt.Printf("是否需要继续转换%s，如果需要请输入y:", j)
		fmt.Scanln(&tag)
		Next(tag)
		num3 = HumanBytesLoaded3(num3)
		fmt.Printf("已转换的值：%g%s \n", num3, j)
	}
}

// HumanBytesLoaded 单位转换(B KB MB GB TB EB)
// B——字节,KB——千比特,MB——兆比特,GB——吉比特
func HumanBytesLoaded3(bytesLength float64) float64 {
	resp := bytesLength / 1024
	
	return resp
}

//func HumanBytesLoaded2(bytesLength float64) {
//	// bytesLength
//	fmt.Printf("需要转换的值：%gB\n", bytesLength)
//	resp := bytesLength / 1024
//	fmt.Printf("已转换的值：%gKB\n", resp)
//	resp2 := resp / 1024
//	fmt.Printf("已转换的值：%gMB\n", resp2)
//	resp3 := resp2 / 1024
//	fmt.Printf("已转换的值：%gGB\n", resp3)
//	resp4 := resp3 / 1024
//	fmt.Printf("已转换的值：%gPB\n", resp4)
//}

//func HumanBytesLoaded1(bytesLength int) string {
//    var resp string
//
//    fmt.Printf("需要转换的值: %dB\n", bytesLength)
//
//    num_KB := bytesLength / 1024
//    resp = strconv.Itoa(num_KB)

//    return resp
//}

//func user_Num {
//	var iowrite float64
//	fmt.Print("请输入一个需要转换的B（字节）：")
//	fmt.Scanln(&iowrite)
//}

func Next(tag string) {
	str := "y"

	if tag == str {
		fmt.Println("程序继续")
	} else {
		fmt.Println("程序终止，感谢使用！")
		time.Sleep(99999999999999999) // tmp
	}
}
