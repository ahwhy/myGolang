package user

import (
	"fmt"
	"strconv"

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

del函数
	从命令行中输入要删除的用户ID
	验证ID是否存在，如果存在，打印需要删除的用户信息
	并让用户输入y/n 确认是否删除
	输入y删除用户信息

	fmt.Println(users)
*/

func Deluser(user []map[string]string) []map[string]string {
	tag := false
	tag2 := false

	usge := func() {
		fmt.Printf("请输入你要删除用户的%s: ", Uid)
	}
	usge2 := func() {
		fmt.Println("你要删除的用户不存在")
	}

	for {
		num := 0
		tag3 := true

		usge()
		fmt.Scanln(&uid)

		for i, u := range user {
			u_uid := u[Uid]
			if u_uid == strconv.Itoa(uid) {
				tag3 = false
				num = i
				Inputuser(u)
				next := utils.Input("是否确认删除用户？确认请输入Y ")
				if next == "Y" || next == "y" || next == "Yes" || next == "yes" {
					tag = true
				}
			}
		}
		if tag {
			if num == 0 {
				user = user[num+1:]
			} else if num == len(user) {
				user = user[:num-1]
			} else {
				user = append(user[:num], user[num+1:]...)
			}
		}

		if tag3 {
			usge2()
		}

		next := utils.Input("是否继续删除其他用户？退出请输入 N ")
		if next == "N" || next == "n" || next == "No" || next == "no" {
			tag2 = true
		}
		if tag2 {
			return user
		}
	}
}
