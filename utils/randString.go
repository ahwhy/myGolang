package utils

import (
	"math/rand"
	"time"
)

// RandString 生成随机字符串
func RandString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	len := len(letters)
	chars := make([]byte, length)

	for i := 0; i < length; i++ {
		chars[i] = letters[rand.Int()%len]
	}

	return string(chars)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
