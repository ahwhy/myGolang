package loghook_test

import (
	"sync"
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

func TestSynchronize(t *testing.T) {
	dh := loghook.NewDingHook(url, app, lev, user)

	dh.Synchronize("测试 ---> 同步告警")
}

func TestAsynchronouse(t *testing.T) {
	dh := loghook.NewDingHook(url, app, lev, user)

	msg := []byte("测试 ---> 异步告警")
	dh.JsonBodies <- msg
	defer close(dh.JsonBodies)

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		dh.Asynchronous()
		wait.Done()
	}()
	wait.Wait()
}

func TestDingHook(t *testing.T) {
	dh := loghook.NewDingHook(url, app, lev, user)

	level := logger.InfoLevel
	logger.SetLevel(level)

	// 设置filename
	logger.SetReportCaller(true)
	logger.SetFormatter(&logger.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 添加hook
	logger.AddHook(dh)
	logger.Info("Hook logrus")
}
