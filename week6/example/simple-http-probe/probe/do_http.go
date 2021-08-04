package probe

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/ahwhy/myGolang/week6/homework/simple-http-probe/config"
)

// 使用 net/http/httptrace 实现http耗时探测
func DoHttpProbe(url string) string {
	twSec := config.GlobalTwSec

	// test
	// return fmt.Sprintf("%s + %d", url, twSec)

	// 定义输出 String
	dnsStr := ""
	targetAddr := ""
	dnsinfo := ""
	gotconinfo := ""
	// 定义时间对象
	var t0, t1, t2, t3, t4 time.Time
	// 全局/报错使用
	start := time.Now()
	// 初始化http的req对象
	req, _ := http.NewRequest("GET", url, nil) // func NewRequest(method, url string, body io.Reader) (*Request, error)
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			// 开始解析DNS时间，t0
			t0 = time.Now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			// 	DNS解析完成时间，t1
			t1 = time.Now()
			dnsinfo = fmt.Sprintf("DNS Info: %+v\n", dnsInfo)
			ips := make([]string, 0)
			for _, d := range dnsInfo.Addrs {
				ips = append(ips, d.IP.String())
			}
			dnsStr = strings.Join(ips, ",")
		},
		ConnectStart: func(network, addr string) {
			if t1.IsZero() {
				// 传参时，直接传递IP地址，则没有DNS解析；ConnectStart为strat
				t1 = time.Now()
			}
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("[无法和探测目标建立连接][Addr:%v][Err:%v]", addr, err)
			}
			targetAddr = addr
			// 第一次成功建立连接的时间，t2
			t2 = time.Now()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			// 成功建立Tcp连接的时间，t3
			t3 = time.Now()
			gotconinfo = fmt.Sprintf("Got Conn: %+v\n", connInfo)
		},
		GotFirstResponseByte: func() {
			// 服务端第一次返回的时间，t
			t4 = time.Now()
		},
	}

	// 标准用法
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	// 目的：创建一个有超时时间的 客户端
	client := http.Client{
		Timeout: time.Duration(twSec) * time.Second,
	}
	resp, err := client.Do(req)
	// 出错的情况
	if err != nil {
		msg := fmt.Sprintf("[HTTP探测出错]\n"+
			"[Http探测的目标: %s]\n"+
			"[错误详情: %v]\n"+
			"[总耗时: %s]\n",
			url,
			err,
			msDurationStr(time.Now().Sub(start)),
		)
		log.Printf(msg)
		return msg
	}
	defer resp.Body.Close()
	end := time.Now()
	// 没有DNS
	if t0.IsZero() {
		t0 = t1
	}

	dnsLookup := msDurationStr(t1.Sub(t0))
	firstConnection := msDurationStr(t2.Sub(t1))
	tcpConnection := msDurationStr(t3.Sub(t1))
	serverProcessing := msDurationStr(t4.Sub(t3))
	totoal := msDurationStr(end.Sub(t0))
	probeResStr := fmt.Sprintf(
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
		url,
		dnsStr,
		targetAddr,
		resp.StatusCode,
		dnsLookup,
		firstConnection,
		tcpConnection,
		serverProcessing,
		totoal,
		dnsinfo,
		gotconinfo,
	)
	return probeResStr
}

// 传递 time.Duration类型，返回ms 单位的字符串
func msDurationStr(d time.Duration) string {
	return fmt.Sprintf("%dms", int(d/time.Microsecond)) // s/1000 = ms
}
