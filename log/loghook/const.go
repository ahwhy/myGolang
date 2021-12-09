package loghook

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
	}
	Text struct {
		Content string `json:"content"`
	}
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
