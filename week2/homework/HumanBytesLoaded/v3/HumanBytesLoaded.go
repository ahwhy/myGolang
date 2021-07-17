package main

import (
	"fmt"
)

func main() {
	var num int64

	fmt.Print("请输入一个需要转换的B（字节）：")
	fmt.Scanln(&num)

	fmt.Printf("需要转换的值：%dB\n", num)

	fmt.Println(HumanBytesLoaded(num))

	fmt.Println("已完成所有转换，感谢使用！")
}

const (
	bu = 1 << 10
	kb = 1 << 20
	mb = 1 << 30
	gb = 1 << 40
	tb = 1 << 50
	eb = 1 << 60
)

// HumanBytesLoaded 单位转换
func HumanBytesLoaded(bytesLength int64) string {
	if bytesLength < bu {
		return fmt.Sprintf("%dB", bytesLength)
	} else if bytesLength < kb {
		return fmt.Sprintf("%.2fKB", float64(bytesLength)/float64(bu))
	} else if bytesLength < mb {
		return fmt.Sprintf("%.2fMB", float64(bytesLength)/float64(kb))
	} else if bytesLength < gb {
		return fmt.Sprintf("%.2fGB", float64(bytesLength)/float64(mb))
	} else if bytesLength < tb {
		return fmt.Sprintf("%.2fTB", float64(bytesLength)/float64(gb))
	} else {
		return fmt.Sprintf("%.2fEB", float64(bytesLength)/float64(tb))
	}
}
