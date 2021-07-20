package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 使用+拼接
func randomStringadd(n int) string {
	b := ""
	for i := 0; i < n; i++ {
		b += string(letterBytes[rand.Intn(len(letterBytes))])
	}
	return b
}

// 使用fmt.Sprinf拼接
func randomStringfmts(n int) string {
	b := ""
	for i := 0; i < n; i++ {
		b += fmt.Sprintf("%c", letterBytes[rand.Intn(len(letterBytes))])
	}
	return b
}

// 使用fmt.Sprinf拼接2
func randomStringfmts2(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return fmt.Sprintf("%s", b)
}

// 使用fmt.Prinf拼接
func randomStringfmtp(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("%c", letterBytes[rand.Intn(len(letterBytes))])
	}
	fmt.Println()
}

// 使用strings.Builder拼接
func randomStringbuilder(n int) string {
	b := &strings.Builder{}
	for i := 0; i < n; i++ {
		b.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}
	return b.String()
}

// 使用bytes.Buffer拼接
func randomStringbytes(n int) string {
	b := &bytes.Buffer{}
	for i := 0; i < n; i++ {
		b.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}
	return b.String()
}

// 使用[]byte拼接
func randomStringbyte(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println(randomStringadd(100))
	fmt.Println(randomStringfmts(100))
	fmt.Println(randomStringfmts2(100))
	randomStringfmtp(100)
	fmt.Println(randomStringbuilder(100))
	fmt.Println(randomStringbytes(100))
	fmt.Println(randomStringbyte(100))
}
