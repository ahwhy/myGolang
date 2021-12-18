package loghook

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
)

func NewDingHook(url, app string, lev []logger.Level, user []string) *dingHook {
	return &dingHook{
		apiUrl:     url,
		levels:     lev,
		atMobiles:  user,
		appName:    app,
		JsonBodies: make(chan []byte, 1),
		closeChan:  make(chan bool),
	}
}

// 钉钉Hook 结构体
type dingHook struct {
	apiUrl     string // 钉钉群机器人 token url
	levels     []logger.Level
	atMobiles  []string    // 通知的对象
	appName    string      // 模块前缀
	JsonBodies chan []byte // 异步发送内容队列
	closeChan  chan bool   // 主进程关闭消息通道
}

// Levels Hook的日志级别
func (dh *dingHook) Levels() []logger.Level {
	return dh.levels
}

// Fire 发送日志的具体逻辑 同步
// func (dh *dingHook) Fire(e *logger.Entry) error {
// 	msg, _ := e.String()
// 	if err := dh.Synchronize(msg); err != nil {
// 		return err
// 	}

// 	return nil
// }

// Fire 发送日志的具体逻辑 异步
func (dh *dingHook) Fire(e *logger.Entry) error {
	defer func() {
		<-dh.closeChan
	}()

	msg, _ := e.Bytes()
	dh.JsonBodies <- msg
	defer close(dh.JsonBodies)

	go dh.Asynchronous()

	return nil
}

// SynchronizeSend 报警主体函数  -> 同步模式 推送告警信息给钉钉
func (dh *dingHook) Synchronize(msg string) error {
	dm := NewDingMsg()
	dm.At.AtMobiles = dh.atMobiles
	dm.Text.Content = fmt.Sprintf("[日志告警]\n[app=%s]\n[日志详情: %s]", dh.appName, msg)

	if err := dm.DirectSend(dh.apiUrl); err != nil {
		return err
	}

	return nil
}

// AsynchronousSend 报警主体函数  -> 异步模式 推送告警信息给钉钉
func (dh *dingHook) Asynchronous() error {
	dm := NewDingMsg()
	dm.At.AtMobiles = dh.atMobiles
	dm.Text.Content = fmt.Sprintf("[日志告警]\n[app=%s]\n[日志详情: %s]", dh.appName, string(<-dh.JsonBodies))

	if err := dm.DirectSend(dh.apiUrl); err != nil {
		return err
	}

	dh.closeChan <- true
	defer close(dh.closeChan)

	return nil
}

/*
- github.com/sirupsen/logrus
 	- 首先定义相关结构体，然后实现Levels和Fire方法 --> 实现 Hook接口
		- Levels中定义日志等级
		- Fire中处理日志发送逻辑
			- 比如发送到redis、es、钉钉、logstash
 	- 调用AddHook()，直接打印日志并发送

func AddHook(hook Hook) {
	std.AddHook(hook)
}

type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}

func (hooks LevelHooks) Fire(level Level, entry *Entry) error {
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}
	return nil
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(InfoLevel, args...)
}
*/
