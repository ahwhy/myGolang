package probe

import (
	"fmt"
	"log"
)

func NewHttpInfo(url string) *HttpInfo {
	return &HttpInfo{
		Url: url,
	}
}

type HttpInfo struct {
	Url              string
	Dnsinfo          string
	DnsStr           string
	TargetAddr       string
	Gotconinfo       string
	StatusCode       int
	DnsLookup        string
	FirstConnection  string
	TcpConnection    string
	ServerProcessing string
	Totoal           string

	Err error
}

func (info *HttpInfo) ProbeFail() string {
	msg := fmt.Sprintf("[HTTP探测出错]\n"+
		"[Http探测的目标: %s]\n"+
		"[状态码: %d]\n"+
		"[错误详情: %v]\n"+
		"[总耗时: %s]\n",
		info.Url,
		info.StatusCode,
		info.Err,
		info.Totoal,
	)

	log.Printf(msg)

	return msg
}

func (info *HttpInfo) ProbeOk() string {
	msg := fmt.Sprintf(
		"[Http探测的目标: %s]\n"+
			"[Dns解析的结果: %s]\n"+
			"[连接的Ip和端口: %s]\n"+
			"[状态码: %d]\n"+
			"[Dns解析耗时: %s]\n"+
			"[第一次连接耗时: %s]\n"+
			"[Tcp连接耗时: %s]\n"+
			"[服务端处理耗时: %s]\n"+
			"[总耗时: %s]\n"+
			"[Dns返回信息: %v]\n"+
			"[Got Conn返回信息: %v]\n",
		info.Url,
		info.DnsStr,
		info.TargetAddr,
		info.StatusCode,
		info.DnsLookup,
		info.FirstConnection,
		info.TcpConnection,
		info.ServerProcessing,
		info.Totoal,
		info.Dnsinfo,
		info.Gotconinfo,
	)

	log.Printf(msg)

	return msg
}
