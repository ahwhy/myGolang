package database

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	sqlInjectRegx *regexp.Regexp
)

func init() {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	var err error
	sqlInjectRegx, err = regexp.Compile(str)
	if err != nil {
		panic(err)
	}
}

// 对用户输入执行严格的校验，不能包含特殊符号和mysql保留字
func FilteredSQLInject(input string) bool {
	return sqlInjectRegx.MatchString(strings.ToLower(input))
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error()) // stdout是行缓冲的，其输出会放在一个buffer里面，只有到换行的时候，才会输出到屏幕；而stderr是无缓冲的，会直接输出
		os.Exit(1)
	}
}
