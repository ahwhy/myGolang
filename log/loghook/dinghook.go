package loghook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

func NewDingHook(url, app string, lev []logger.Level, user []string) *dingHook {
	return &dingHook{
		apiUrl:     url,
		levels:     lev,
		atMobiles:  user,
		appName:    app,
		jsonBodies: make(chan []byte),
		closeChan:  make(chan bool),
	}
}

// 钉钉Hook 结构体
type dingHook struct {
	apiUrl     string // 钉钉群机器人 token url
	levels     []logger.Level
	atMobiles  []string    // 通知的对象
	appName    string      // 模块前缀
	jsonBodies chan []byte // 异步发送内容队列
	closeChan  chan bool   // 主进程关闭消息通道
}

// Levels Hook的日志级别
func (dh *dingHook) Levels() []logger.Level {
	return dh.levels
}

// Fire 发送日志的具体逻辑
func (dh *dingHook) Fire(e *logger.Entry) error {
	msg, _ := e.String()
	if err := dh.DirectSend(msg); err != nil {
		return err
	}

	return nil
}

// 设置报警主体函数  -> 同步模式 推送告警信息给钉钉
func (dh *dingHook) DirectSend(msg string) error {
	dm := NewDingMsg()
	dm.At.AtMobiles = dh.atMobiles
	dm.Text.Content = fmt.Sprintf("[日志告警log]\n[app=%s]\n[日志详情: %s]", dh.appName, msg)

	bs, err := json.Marshal(dm)
	if err != nil {
		logger.Errorf("[消息json.marshal失败][error:%v][msg:%v]", err, msg)
		return err
	}

	res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		logger.Errorf("[消息发送失败][error:%v][msg:%v]", err, msg)
		return err
	}
	if res != nil || res.StatusCode != 200 {
		logger.Println(res, res.StatusCode)
		logger.Errorf("[钉钉返回错误][StatusCode:%v][msg:%v]", res.StatusCode, msg)
		return err
	}

	return nil
}

// func init() {
// 	level := logger.InfoLevel
// 	logger.SetLevel(level)

// 	// 设置filename
// 	logger.SetReportCaller(true)
// 	logger.SetFormatter(&logger.JSONFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 	})
// }

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

/*
 - 改造成异步发送的

**调用者：各种服务进程**

**被调用者：内核**

**同步（synchronous）：**被调用者（内核）不提供通知机制，调用者（服务进程）需要主动反复询问，被调用者（内核）该事件是否处理完成

 **异步（asynchronous）：**被调用者（内核）通过某种机制告诉调用者（服务进程）事件的处理进度和运行状态

1、同步：进程发出请求调用后，等内核返回响应以后才继续下一个请求，即如果内核一直不返回数据，那么进程就一直等。

2、异步：进程发出请求调用后，不等内核返回响应，接着处理下一个请求,**Nginx是异步的**
*/
