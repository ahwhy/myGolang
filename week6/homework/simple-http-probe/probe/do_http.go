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

//用 net/http/httptrace 写个http耗时探测的项目 simple-http-probe

func DoHttpProbe(url string) string {
	twSec := config.GlobalTwSec

	// 结果string
	dnsStr := ""
	targetAddr := ""
	// 提前定义好这些计算时间的对象
	var t0, t1, t2, t3, t4 time.Time
	// 全局或出错使用
	start := time.Now()
	// 初始化http req对象
	req, _ := http.NewRequest("GET", url, nil)

	//return fmt.Sprintf("%s + %d", url, twSec)

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			// 开始dns解析的时候我 赋值t0
			t0 = time.Now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			// dns解析完成时间为 t1
			t1 = time.Now()
			ips := make([]string, 0)
			for _, d := range dnsInfo.Addrs {
				ips = append(ips, d.IP.String())
			}
			dnsStr = strings.Join(ips, ",")
		},
		ConnectStart: func(network, addr string) {
			if t1.IsZero() {
				// 直接传ip没dns解析 开始连接算start
				t1 = time.Now()
			}
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("[无法建立和探测目标的连接][addr:%v][err:%v]", addr, err)
				return
			}
			targetAddr = addr
			t2 = time.Now()
		},

		GotConn: func(_ httptrace.GotConnInfo) {
			t3 = time.Now()
		},
		GotFirstResponseByte: func() {
			t4 = time.Now()
		},
	}
	// 标准用法
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	// 这里的目的就是造个有超时时间的 客户端
	client := http.Client{
		Timeout: time.Duration(twSec) * time.Second,
	}
	resp, err := client.Do(req)
	// 出错的情况
	if err != nil {
		msg := fmt.Sprintf("[http探测出错]\n"+
			"[http探测的目标:%s]\n"+
			"[错误详情:%v]\n"+
			"[总耗时:%s]\n",
			url,
			err,
			msDurationStr(time.Now().Sub(start)),
		)
		log.Printf(msg)
		return msg

	}
	defer resp.Body.Close()
	end := time.Now()
	// 没有dns
	if t0.IsZero() {
		t0 = t1
	}
	fmt.Println(t2)
	dnsLookup := msDurationStr(t1.Sub(t0))
	tcpConnection := msDurationStr(t3.Sub(t1))
	serverProcessing := msDurationStr(t4.Sub(t3))
	totoal := msDurationStr(end.Sub(t0))
	probeResStr := fmt.Sprintf(
		"[http探测的目标:%s]\n"+
			"[dns解析的结果:%s]\n"+
			"[连接的ip和端口：%s]\n"+
			"[状态码:%d]\n"+
			"[dns解析耗时：%s]\n"+
			"[tcp连接耗时：%s]\n"+
			"[服务端处理耗时：%s]\n"+
			"[总耗时：%s]\n",
		url,
		dnsStr,
		targetAddr,
		resp.StatusCode,
		dnsLookup,
		tcpConnection,
		serverProcessing,
		totoal,
	)
	return probeResStr
}

func msDurationStr(d time.Duration) string {
	return fmt.Sprintf("%dms", int(d/time.Millisecond))
}
