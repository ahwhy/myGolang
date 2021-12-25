package socket

import (
	"fmt"
	"os"
)

var (
	ListenAddr = "0.0.0.0"
	IP         = "127.0.0.1" // ip也可设置成 0.0.0.0 和 空字符串
	Port       = 5656        // 改成 1023，会报错 bind: permission denied
)

type (
	Request struct {
		A int
		B int
	}
	Response struct {
		Sum int
	}
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error()) // stdout是行缓冲的，他的输出会放在一个buffer里面，只有到换行的时候，才会输出到屏幕；而stderr是无缓冲的，会直接输出
		os.Exit(1)
	}
}
