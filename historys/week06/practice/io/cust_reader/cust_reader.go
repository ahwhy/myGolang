package main

import (
	"io"
	"log"
)

type zimuguolv struct {
	src string
	cur int
}

func alpha(r byte) byte {
	// r在 A-Z 或者 a-z
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func (z *zimuguolv) Read(p []byte) (int, error) {
	// 当前位置 >= 字符串长度，说明已经读取到结尾了，返回 EOF
	if z.cur >= len(z.src) {
		return 0, io.EOF
	}
	// 定义一个剩余还没读到的长度
	x := len(z.src) - z.cur
	// bound叫做本次读取长度
	// n代表本次遍历 bound的索引
	n, bound := 0, 0
	if x >= len(p) {
		// 剩余长度超过缓冲区大小，说明本次可以完全填满换冲区
		bound = len(p)
	} else {
		// 剩余长度小于缓冲区大小，使用剩余长度输出，缓冲区填不满
		bound = x
	}

	buf := make([]byte, bound)

	for n < bound {
		if char := alpha(z.src[z.cur]); char != 0 {
			buf[n] = char
		}
		// 索引++
		n++
		z.cur++
	}
	copy(p, buf)
	return n, nil
}

func main() {
	zmreader := zimuguolv{
		src: "mage jiaoyu 2021 go !!!!",
	}

	p := make([]byte, 4)
	for {
		n, err := zmreader.Read(p)
		if err == io.EOF {
			log.Printf("[EOF错误]")
			break
		}
		log.Printf("[读取到的长度%d 内容%s]", n, string(p[:n]))
	}
}
