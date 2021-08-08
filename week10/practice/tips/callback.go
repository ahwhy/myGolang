package tips

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
)

type SiteResp struct {
	Err    error
	Resp   string
	Status int
	Cost   int64
}

func CallBackMode() {
	endpoints := []string{
		"https://www.baidu.com",
		"https://segmentfault.com/",
		"https://blog.csdn.net/",
		"https://www.jd.com/",
	}

	// 一个endpoints返回一个结果, 缓冲可以确定
	respChan := make(chan SiteResp, len(endpoints))
	defer close(respChan)

	// 回调处理逻辑
	ret := []SiteResp{}
	cb := func(v SiteResp) {
		ret = append(ret, v)
	}

	// 并行爬取
	for _, endpoints := range endpoints {
		wg.Add(1)
		go doSiteRequest(cb, endpoints)
	}

	// 等待结束
	wg.Wait()

	for _, v := range ret {
		fmt.Println(v)
	}
}

type SiteRespCallback func(SiteResp)

// 构造请求
func doSiteRequest(cb SiteRespCallback, url string) {
	res := SiteResp{}
	startAt := time.Now()
	defer func() {
		res.Cost = time.Since(startAt).Milliseconds()
		cb(res)
		wg.Done()
	}()

	resp, err := client.Get(url)
	if resp != nil {
		res.Status = resp.StatusCode
	}
	if err != nil {
		res.Err = err
		return
	}

	// 站不处理结果
	_, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Err = err
		return
	}

	// res.Resp = string(byt)
}
