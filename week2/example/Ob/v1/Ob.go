package main

import (
	"fmt"
	"strconv"
	"strings"
)

//请将这段二进制翻译成中文(unicode编码), b并打印出翻译过程(比如 二进制: 1000 1000 0110 0011, unicode: U+8863, 字符: 衣)
// 1000 1000 0110 0011  衣
// textUnquoted
// 0101 0101 1001 1100  喜
// 0110 1011 0010 0010  欢
// 0111 1010 0111 1111  穿
// 0100 1110 0010 1101  中
// 0101 0110 1111 1101  国
// 0111 1110 1010 0010  红

func main() {
	sliceBin:=[]int{0b0101010110011100,0b0110101100100010,0b0111101001111111,0b0100111000101101,0b0101011011111101,0b0111111010100010}

	for _,u := range sliceBin{
		unicodeStr:=fmt.Sprintf("%U",u)
		toCN(unicodeStr)
	}
	fmt.Println()
}

//将unicode转换成中文
func toCN(u string) string {
	uSp := strings.Split(u, "U+") //通过U+来切割unicode
	var context string

	for _, v := range uSp {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32) //将遍历出的值通过16进制转换为10进制，返回结果为int32
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)  //通过字符连接起来
	}
	fmt.Printf("%v",context)
	
	return context
}
