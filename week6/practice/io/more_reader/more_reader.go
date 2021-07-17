package main

import (
	"io"
	"log"
	"strings"
)

type alphaReader struct {
	ioReader io.Reader
}

func alpha(r byte) byte {
	// r在 A-Z 或者 a-z
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func (a *alphaReader) Read(p []byte) (int, error) {
	// 复用io.reader的read方法
	n, err := a.ioReader.Read(p)
	if err != nil {
		return n, err
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}
	copy(p, buf)
	return n, nil
}

func main() {
	myReader := alphaReader{
		strings.NewReader("mage jiaoyu 2021 go !!!"),
	}
	p := make([]byte, 4)
	for {
		n, err := myReader.Read(p)
		if err == io.EOF {
			log.Printf("[EOF错误]")
			break
		}
		log.Printf("[读取到的长度%d 内容%s]", n, string(p[:n]))
	}
}
