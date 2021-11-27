package impl

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/ahwhy/myGolang/http_probe/config"
	"github.com/ahwhy/myGolang/http_probe/probe"
)

func DoHttpProbe(url string) string {
	twSec := config.GlobalTwSec
	httpInfo := probe.NewHttpInfo(url)

	// 定义时间对象
	var t0, t1, t2, t3, t4 time.Time
	// 全局/报错使用
	start := time.Now()

	// 初始化http的req对象
	// func NewRequest(method, url string, body io.Reader) (*Request, error)
	req, _ := http.NewRequest("GET", url, nil)
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			// 开始解析DNS时间，t0
			t0 = time.Now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			// 	DNS解析完成时间，t1
			t1 = time.Now()
			httpInfo.Dnsinfo = fmt.Sprintf("DNS Info: %+v\n", dnsInfo)
			ips := make([]string, 0)
			for _, d := range dnsInfo.Addrs {
				ips = append(ips, d.IP.String())
			}
			httpInfo.DnsStr = strings.Join(ips, ",")
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
			httpInfo.TargetAddr = addr
			// 第一次成功建立连接的时间，t2
			t2 = time.Now()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			// 成功建立Tcp连接的时间，t3
			t3 = time.Now()
			httpInfo.Gotconinfo = fmt.Sprintf("Got Conn: %+v\n", connInfo)
		},
		GotFirstResponseByte: func() {
			// 服务端第一次返回的时间，t4
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
		httpInfo.Totoal = probe.MsdurationStr(time.Now().Sub(start))
		httpInfo.Err = err

		return httpInfo.ProbeFail()
	}
	defer resp.Body.Close()

	end := time.Now()
	// 没有DNS
	if t0.IsZero() {
		t0 = t1
	}

	httpInfo.StatusCode = resp.StatusCode
	httpInfo.DnsLookup = probe.MsdurationStr(t1.Sub(t0))
	httpInfo.FirstConnection = probe.MsdurationStr(t2.Sub(t1))
	httpInfo.TcpConnection = probe.MsdurationStr(t3.Sub(t1))
	httpInfo.ServerProcessing = probe.MsdurationStr(t4.Sub(t3))
	httpInfo.Totoal = probe.MsdurationStr(end.Sub(t0))

	return httpInfo.ProbeOk()
}
