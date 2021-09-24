package utils

import (
	"regexp"
	"strings"
)

// Snake 用于将驼峰式命名字符串替换为下划线连接
func Snake(txt string) string {
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := reg.ReplaceAllString(txt, "${1}_${2}")

	return strings.ToLower(snake)
}