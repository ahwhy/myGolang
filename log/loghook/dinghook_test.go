package loghook_test

import (
	"testing"

	"github.com/ahwhy/myGolang/log/loghook"

	logger "github.com/sirupsen/logrus"
)

var (
	url  = "https://oapi.dingtalk.com/robot/send?access_token=438a6ae04b3abf3cd7834f2d294c3c4cf6ffae2e956fabc64751f62b766e1e16"
	user = []string{"liangxiao", "yuyang"}
	lev  = []logger.Level{logger.WarnLevel, logger.InfoLevel}
	app  = "test"
)

func TestDingHook(t *testing.T) {
	dh := loghook.NewDingHook(url, app, lev, user)

	dh.DirectSend("测试 --> 直接发送信息")

	// level := logger.InfoLevel
	// logger.SetLevel(level)
	// // 设置filename
	// logger.SetReportCaller(true)
	// logger.SetFormatter(&logger.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
	// 添加hook
	logger.AddHook(dh)
	logger.Info("这是hook的logrus")
}
