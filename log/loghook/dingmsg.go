package loghook

import (
	"bytes"
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

func NewDingMsg() *dingMsg {
	return &dingMsg{
		Msgtype: "text",
	}
}

type dingMsg struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Msgtype string `json:"msgtype"`
}

/*
- text类型消息要求字段
{
    "at": {
        "atMobiles":[
            "180xxxxxx"
        ],
        "atUserIds":[
            "user123"
        ],
        "isAtAll": false
    },
    "text": {
        "content":"我就是我, @XXX 是不一样的烟火"
    },
    "msgtype":"text"
}
*/

func (dm *dingMsg) DirectSend(url string) error {
	bs, err := json.Marshal(dm)
	if err != nil {
		logger.Errorf("[消息json.marshal失败][error:%v][msg:%v]", err, dm.Text.Content)
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		logger.Errorf("[消息发送失败][error:%v][msg:%v]", err, dm.Text.Content)
		return err
	}
	if res != nil || res.StatusCode != 200 {
		logger.Println(res, res.StatusCode)
		logger.Errorf("[钉钉返回错误][StatusCode:%v][msg:%v]", res.StatusCode, dm.Text.Content)
		return err
	}

	return nil
}
