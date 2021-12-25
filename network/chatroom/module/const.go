package module

import "time"

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 9 * pongWait / 10
	maxMsgSize = 512 // 消息的长度不能超过512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)
