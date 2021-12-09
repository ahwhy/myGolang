package main

import (
	"github.com/ahwhy/myGolang/historys/week07/practice/log/ilog"
)

type User struct {
	Name string
	Age  int
}

func main() {
	a := User{
		Name: "YY",
		Age:  18,
	}
	ilog.InitLogger("info")
	ilog.Logger.Infoln("提示信息")
	ilog.Logger.Infoln(a)
	ilog.InitLogger("debug")
	ilog.Logger.Debugln("调试信息")
	ilog.InitLogger("warn")
	ilog.Logger.Warnln("警告信息")
	ilog.InitLogger("error")
	ilog.Logger.Errorln("错误信息")

	// log.Infof("[格式化打印结构体:%+v]", a)
	// log.WithFields(log.Fields{
	// 	"user_id":    123,
	// 	"ip":         "1.1.1.1",
	// 	"request_id": "asdwdadmaskmdlasmldkmwqlkdkmakldm",
	// }).Info("用户登录成功")
}
