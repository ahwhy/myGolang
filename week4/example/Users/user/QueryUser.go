package user

import (
	"strings"

	"github.com/ahwhy/myGolang/week4/example/Users/utils"
)

/*
定义全局变量
var users = []map[string]string{}

每个元素
	ID
	名称
	联系方式
	通信地址

query函数
	从命令行输入要查询的字符串
	遍历比较用户的名称，地址，联系方式，包含要查找的字符串就进行输出

	fmt.Println(users)
*/
func filer(user map[string]string, q string) bool {
	return strings.Contains(user[Name], q) ||
		strings.Contains(user[Tel], q) ||
		strings.Contains(user[Addr], q)
}

func QueryUser(user []map[string]string) []map[string]string {
	tag := false

	for {
		text := utils.Input("请输入你要查询的内容: ")
		for _, u := range user {
			if filer(u, text) {
				Inputuser(u)
			}
		}

		next := utils.Input("是否继续查询信息？退出请输入 N ")
		if next == "N" || next == "n" || next == "No" || next == "no" {
			tag = true
		}

		if tag {
			return user
		}
	}
}
