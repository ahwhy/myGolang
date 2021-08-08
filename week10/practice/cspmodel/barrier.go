package cspmodel

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
)

type SiteResp struct {
	Site   string
	Err    error
	Resp   string
	Status int
	Cost   int64
}

func BarrierMode() {
	endpoints := []string{
		"https://www.baidu.com",
		"https://segmentfault.com/",
		"https://blog.csdn.net/",
		"https://www.jd.com/",
	}

	// 一个endpoints返回一个结果, 缓冲可以确定
	respChan := make(chan SiteResp, len(endpoints))
	defer close(respChan)

	// 并行爬取
	for _, endpoints := range endpoints {
		go doSiteRequest(respChan, endpoints)
	}

	// 聚合结果
	down := make(chan struct{})
	ret := make([]SiteResp, 0, len(endpoints))
	go mergeResponse(respChan, &ret, down)

	// 等待结束
	<-down

	for _, v := range ret {
		fmt.Println(v)
	}
}

// 1. 从返回管道里面取数据
// 2. 所有往返的返回了结果, 取了4次数据, 4个网站都返回了数据, 我往外发送一个信号, merge工作完成(down)
func mergeResponse(resp <-chan SiteResp, ret *[]SiteResp, down chan struct{}) {
	defer func() {
		down <- struct{}{}
	}()

	count := 0
	for v := range resp {
		*ret = append(*ret, v)
		count++

		// 填充完成,  返回
		if count == cap(*ret) {
			return
		}
	}

}

// 如何访问一个网站, go http
// 构造一个并行体
func doSiteRequest(out chan<- SiteResp, url string) {
	res := SiteResp{
		Site: url,
	}
	startAt := time.Now()

	defer func() {
		res.Cost = time.Since(startAt).Milliseconds()
		out <- res
	}()

	// 爬取网页
	resp, err := client.Get(url)
	if resp != nil {
		res.Status = resp.StatusCode
	}

	if err != nil {
		res.Err = err
		return
	}

	// 站不处理结果
	_, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Err = err
		return
	}

	// res.Resp = string(byt)
}
