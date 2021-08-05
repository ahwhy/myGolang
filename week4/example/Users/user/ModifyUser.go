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

modify函数
	从命令行输入要修改的用户ID
	验证ID是否存在，如果存在，打印需要修改的用户信息
	并让用户输入y/n确认是否修改
	输入y修改用户信息，继续让用户从命令行分别输入用户名，联系方式，地址，进行更新

	fmt.Println(users)
*/

func ModifyUser(user []map[string]string) []map[string]string {
	tag := false

	usge := func() {
		fmt.Printf("请输入你要修改用户的%s: ", Uid)
	}
	usge2 := func() {
		fmt.Println("你要修改的用户不存在")
	}

	for {
		tag2 := true

		usge()
		fmt.Scanln(&uid)

		for _, u := range user {
			u_uid := u[Uid]
			if u_uid == strconv.Itoa(uid) {
				tag2 = false
				Inputuser(u)
				next := utils.Input("是否确认修改用户？确认请输入Y ")
				if next == "Y" || next == "y" || next == "Yes" || next == "yes" {
					u[Name] = utils.Input("修改用户的名称为: ")
					u[Tel] = utils.Input("修改用户的联系方式为: ")
					u[Addr] = utils.Input("修改用户的通信地址为: ")
					Inputuser(u)
				}

			}
		}

		if tag2 {
			usge2()
		}

		next := utils.Input("是否继续修改其他用户？退出请输入 N ")
		if next == "N" || next == "n" || next == "No" || next == "no" {
			tag = true
		}
		if tag {
			return user
		}
	}
}
