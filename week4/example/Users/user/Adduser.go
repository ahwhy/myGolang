package user

import (
	"fmt"
	"strconv"

	"github.com/ahwhy/myGolang/week4/example/Users/utils"
)

const (
	Uid  = "UID"
	Name = "名称"
	Tel  = "联系方式"
	Addr = "通信地址"
)

var (
	uid  int
	name string
	tel  int
	addr string
)

/*
定义全局变量
var users = []map[string]string{}

每个元素
	ID
	名称
	联系方式
	通信地址

add函数
	从命令行分别输入名称、联系方式、通信地址

	生成ID => 查找users中最大的Id+1 (无元素ID = 1)
	放入users

	fmt.Println(users)
*/

func Adduser(user []map[string]string) []map[string]string {
	tag := false

	usge := func() {
		fmt.Printf("请依次输入需要创建用户的%s、%s和%s: ", Name, Tel, Addr)
	}
	usge2 := func() {
		fmt.Println("正在退出用户创建")
	}

	for {
		usge()
		fmt.Scanln(&name, &tel, &addr)
		user = append(user, map[string]string{})
		a := 0
		for _, u := range user {
			u_uid, exist := u[Uid]
			if exist {
				b, _ := strconv.Atoi(u_uid)
				if b > a {
					a = b
				}
			} else {
				u[Uid] = strconv.Itoa(a + 1)
				uid = a + 1
				u[Name] = name
				u[Tel] = strconv.Itoa(tel)
				u[Addr] = addr
			}

		}
		Infouser()

		next := utils.Input("是否还需要继续创建其他用户？退出请输入 N/n ")
		if next == "N" || next == "n" || next == "No" || next == "no" {
			tag = true
		}
		if tag {
			usge2()
			return user
		}
	}
}
