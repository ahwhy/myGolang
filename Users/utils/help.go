package utils

import (
	"fmt"
)

func Help() {
	fmt.Println("增加用户请输入: add")
	fmt.Println("删除用户请输入: delete")
	fmt.Println("修改用户请输入: modify")
	fmt.Println("查询用户请输入: query")
	fmt.Println("打印所用户请输入: all")
	fmt.Println("退出请输入: exit")
}
