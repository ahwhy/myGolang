package utils

import(
	"crypto/md5"
	"fmt"
	"strings"
)

// Md5 md5加密
func Md5(txt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(txt)))
}

// Md5Salt md5加密(加盐)
func Md5Salt(txt, salt string) string {
	if salt == "" {
		salt = RandString(8)
	}

	return fmt.Sprintf("%s:%s", salt, Md5(fmt.Sprintf("%s:%s", salt, txt)))
}

// 切分Md5
func SplitMd5Salt(txt string) (string, string) {
	elements := strings.SplitN(txt, ":", 2)
	if len(elements) > 1 {
		return elements[0], elements[1]
	}

	return elements[0], ""
}