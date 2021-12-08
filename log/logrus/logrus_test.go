package logrus_test

import (
	"testing"

	"github.com/ahwhy/myGolang/log/logrus"
	logger "github.com/sirupsen/logrus"
)

type user struct {
	name string
	age  int
}

var a = user{
	name: "YY",
	age:  18,
}

func Test_LogrusMessage(t *testing.T) {
	logrus.InitLogger("info")
	logrus.Logger.Infoln("提示信息")
	logrus.Logger.Infoln(a)

	logrus.InitLogger("debug")
	logrus.Logger.Debugln("调试信息")
	logrus.Logger.Debugln(a)

	logrus.InitLogger("warning")
	logrus.Logger.Warnln("警告信息")
	logrus.Logger.Warnln(a)

	logrus.InitLogger("error")
	logrus.Logger.Errorln("错误信息")
	logrus.Logger.Errorln(a)

}

func Test_LogrusFields(t *testing.T) {
	logrus.InitLogger("debug")
	logrus.Logger.Infof("[格式化打印结构体:%+v]", a)
	logrus.Logger.WithFields(logger.Fields{
		"user_id":    123,
		"ip":         "1.1.1.1",
		"request_id": "asdwdadmaskmdlasmldkmwqlkdkmakldm",
	}).Info("用户登录成功")
}
