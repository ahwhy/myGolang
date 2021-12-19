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
	endpoints = []string{
		"https://www.baidu.com",
		"https://segmentfault.com/",
		"https://blog.csdn.net/",
		"https://www.jd.com/",
	}
)

type SiteResp struct {
	Site   string
	Cost   int64
	Status int
	Resp   string
	Err    error
}

func BarrierMode() {
	respChan := make(chan SiteResp, len(endpoints))
	defer close(respChan)

	// 并行爬取
	for _, endpoint := range endpoints {
		go doSiteRequest(respChan, endpoint)
	}

	// 聚合结果
	down := make(chan struct{})
	ret := make([]SiteResp, 0, len(endpoints))
	go mergeResponse(respChan, &ret, down)

	<-down

	for _, v := range ret {
		fmt.Printf("【爬取网址】%s\n【花费时间】%dms\n【状态码】%d\n【网页内容】%v\n", v.Site, v.Cost, v.Status, v.Resp)
		fmt.Println("-----------------")
	}
}

// doSiteRequest 访问网站，response信息存入SiteResp
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

	byt, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Err = err
		return
	}
	res.Resp = string(byt)
}

// mergeResponse 汇总管道中的数据
func mergeResponse(resp <-chan SiteResp, ret *[]SiteResp, down chan struct{}) {
	defer func() {
		down <- struct{}{}
	}()

	count := 0
	for v := range resp {
		*ret = append(*ret, v)
		count++

		if count == cap(*ret) {
			return
		}
	}
}
