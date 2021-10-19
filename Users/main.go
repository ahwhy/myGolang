package main

import (
	"fmt"
	"os"

	"github.com/ahwhy/myGolang/Users/controller"
	"github.com/ahwhy/myGolang/Users/user"
	"github.com/ahwhy/myGolang/Users/utils"
)

var users = []map[string]string{}

var password = "8a6f2805b4515ac12058e79e66539be9"

const (
	Uid  = "UID"
	Name = "名称"
	Tel  = "联系方式"
	Addr = "通信地址"
)

// var (
// 	uid  int8
// 	name string
// 	tel  int64
// 	addr string
// )

/*
用户增删改查
定义全局变量
var users = []map[string]string{}

每个元素
	ID
	名称
	联系方式
	通信地址

4个程序，每个程序写一个功能
1、增加
	add函数
	从命令行分别输入名称、联系方式、通信地址

	生成ID => 查找users中最大的Id+1 (无元素ID = 1)
	放入users

	fmt.Println(users)

2、删除
	del函数
	从命令行中输入要删除的用户ID
	验证ID是否存在，如果存在，打印需要删除的用户信息
	并让用户输入y/n 确认是否删除
	输入y删除用户信息

3、修改
	modify函数
	从命令行输入要修改的用户ID
	验证ID是否存在，如果存在，打印需要删除的用户信息
	并让用户输入y/n确认是否修改
	输入y修改用户信息，继续让用户从命令行分别输入用户名，联系方式，地址，进行更新

4、查找
	query函数
	从命令行输入要查询的字符串
	遍历比较用户的名称，地址，联系方式，包含要查找的字符串就进行输出

5、用户管理
	循环 让用户从控制台输入指令
		add => 执行 add功能
		delete => 执行delete功能
		modify => 执行modify功能
		query => 执行query功能
		exit => 退出
		help => 帮助信息

6、添加密码功能
	启动程序时让用户输入密码 毕竟对象
	在程序中内置一个md5值

	计算用户输入的密码MD5值，与程序中MD5比较
	输入失败3次退出程序，如果成功执行用户操作

7、用户输出 tablewiter

8、使用映射储存操作指令以及调用函数关系
	map[string]callback
	callback = get(input
	callback)
*/

func main() {
	if !controller.Auth(password) {
		fmt.Println("错误次数过多，程序退出！")
		os.Exit(1)
	}
	for {
		next := utils.Input("请输入要进行的操作: ")
		switch next {
		case "add":
			users = user.Adduser(users)
		case "delete":
			users = user.Deluser(users)
		case "modify":
			users = user.ModifyUser(users)
		case "query":
			users = user.QueryUser(users)
		case "all":
			fmt.Println(users)
		case "help":
			utils.Help()
		case "exit":
			fmt.Println("退出。")
			os.Exit(1)
		default:
			fmt.Println("输入指令错误！")
		}
	}
	// Run()
}

// func init() {
// 	Register("add", user.Adduser)
// 	Register("delete", user.Deluser)
// 	Register("modify", user.ModifyUser)
// 	Register("query", user.QueryUser)
// }

// var routers = map[string]func(u []map[string]string) []map[string]string{}

// func Register(op string, callback func(u []map[string]string) []map[string]string) {
// 	if _, ok := routers[op]; ok {
// 		panic(fmt.Sprintf("指令%s已经存在", op))
// 	}
// 	routers[op] = callback
// }

// func Run() {
// 	for {
// 		text := utils.Input("请输入要进行的操作: ")
// 		if text == "exit" {
// 			fmt.Println("退出程序！")
// 			break
// 		}
// 		if text == "help" {
// 			utils.Help()
// 		}
// 		if text == "all" {
// 			fmt.Println(users)
// 		}
// 		if action, ok := routers[text]; ok {
// 			action(users)
// 		} else {
// 			fmt.Println("输入指令错误！")
// 		}
// 	}
// }
