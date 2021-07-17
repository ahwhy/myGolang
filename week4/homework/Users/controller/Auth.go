package controller

import (
	"fmt"
	"github.com/ahwhy/myGolang/week4/homework/Users/utils"
)

func Auth(password string) bool {
	for i := 0; i < 3; i++ {
		if utils.Md5text(utils.Input("请输入密码: ")) == password {
			return true
		} else {
			fmt.Println("输入密码错误！")
		}
	}
	return false
}
