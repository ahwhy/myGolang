package cspmodel

import (
	"fmt"
	"io"
	"time"
)

// 使用回调函数 代替 Goroutine
func CallBackMode() {
	respChan := make(chan SiteResp, len(endpoints))
	defer close(respChan)

	// 回调处理逻辑
	ret := []SiteResp{}
	cb := func(v SiteResp) {
		ret = append(ret, v)
	}

	// 并行爬取
	for _, endpoint := range endpoints {
		wg.Add(1)
		go doSiteRequestCallback(cb, endpoint)
	}

	// 等待结束
	wg.Wait()

	for _, v := range ret {
		fmt.Printf("【爬取网址】%s\n【花费时间】%dms\n【状态码】%d\n【网页内容】%v\n", v.Site, v.Cost, v.Status, v.Resp)
		fmt.Println("-----------------")
	}
}

type SiteRespCallback func(SiteResp)

// doSiteRequest 访问网站，response信息存入SiteResp
func doSiteRequestCallback(cb SiteRespCallback, url string) {
	res := SiteResp{
		Site: url,
	}
	startAt := time.Now()
	defer func() {
		res.Cost = time.Since(startAt).Milliseconds()
		cb(res)
		wg.Done()
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
