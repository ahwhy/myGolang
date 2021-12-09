package main

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func numLen(n int64) int {
	return len(strconv.Itoa(int(n)))
}

func Test_Nowtime(t *testing.T) {
	now := time.Now()
	log.Printf("[当前时间对象为：%v]", now)
	log.Printf("[当前时间戳 秒级：%v][位数:%v]", now.Unix(), numLen(now.Unix()))
	log.Printf("[当前时间戳 毫秒级：%v][位数:%v]", now.UnixNano()/1e6, numLen(now.UnixNano()/1e6))
	log.Printf("[当前时间戳 纳秒级：%v][位数:%v]", now.UnixNano(), numLen(now.UnixNano()))
	log.Printf("[当前时间戳 纳秒小数部分：%v]", now.Nanosecond())
}

func Test_daytime(t *testing.T) {
	now := time.Now()
	//返回日期

	year, month, day := now.Date()
	log.Printf("[通过now.Data获取][年：%d 月:%d 日:%d]", year, month, day)
	log.Printf("[直接获取年 %d]", now.Year())
	log.Printf("[直接获取月 %d]", now.Month())
	log.Printf("[直接获取日 %d]", now.Day())

	// 分 时 秒

	hour, minute, second := now.Clock()
	log.Printf("[通过now.Clock获取][时：%d 分:%d 秒:%d]", hour, minute, second)
	log.Printf("[直接获取时 %d]", now.Hour())
	log.Printf("[直接获取分 %d]", now.Minute())
	log.Printf("[直接获取秒 %d]", now.Second())

	// 星期几
	log.Printf("[直接获取星期几 %d]", now.Weekday())
	// 时区
	zone, offset := now.Zone()
	log.Printf("[直接获取时区  %v ，和东utc时区差 几个小时：%d]", zone, offset/3600)

	log.Printf("[今天是 %d年 中的第 %d天]", now.Year(), now.YearDay())
}

func Test_formatime(t *testing.T) {
	now := time.Now()
	log.Printf("[全部 ：%v]", now.Format("2006-01-02 15:04:05"))
	log.Printf("[只有年 ：%v]", now.Format("2006"))
	log.Printf("[/分割 ：%v]", now.Format("2006/01/02 15:04"))
}

func Test_ParseInLocation(t *testing.T) {
	tStr := "2021-07-17 16:52:59"
	layout := "2006-01-02 15:04:05"

	t1, _ := time.ParseInLocation(layout, tStr, time.Local)
	t2, _ := time.ParseInLocation(layout, tStr, time.UTC)
	log.Printf("[ %s的 CST时区的时间戳为 ：%d]", tStr, t1.Unix())
	log.Printf("[ %s的 UTC时区的时间戳为 ：%d]", tStr, t2.Unix())
	log.Printf("[UTC - CST =%d 小时]", (t2.Unix()-t1.Unix())/3600)
}

var layout = "2006-01-02 15:04:05"

func tTostr(t time.Time) string {
	return time.Unix(t.Unix(), 0).Format(layout)
}

func Test_ParseDuration(t *testing.T) {
	now := time.Now()
	log.Printf("[当前时间为：%v]", tTostr(now))
	// 1小时1分1秒后
	t1, _ := time.ParseDuration("1h1m1s")
	m1 := now.Add(t1)
	log.Printf("[ 1小时1分1秒后时间为：%v]", tTostr(m1))

	// 1小时1分1秒前
	t2, _ := time.ParseDuration("-1h1m1s")
	m2 := now.Add(t2)
	log.Printf("[ 1小时1分1秒 前时间为：%v]", tTostr(m2))

	// sub计算两个时间差
	sub1 := now.Sub(m2)
	log.Printf("[ 时间差 ：%s 相差小时数：%v 相差分钟数：%v ]", sub1.String(), sub1.Hours(), sub1.Minutes())

	t3, _ := time.ParseDuration("-3h3m3s")
	m3 := now.Add(t3)
	log.Printf("[time.since 当前时间与t的时间差 ：%v]", time.Since(m3))
	log.Printf("[time.until t 当前时间的时间差 ：%v]", time.Until(m3))
	m4 := now.AddDate(0, 0, 5)
	log.Printf("[5天后 的时间 ：%v]", m4)
}

func Test_comparetime(t *testing.T) {
	now := time.Now()
	t1, _ := time.ParseDuration("1h")
	m1 := now.Add(t1)
	log.Printf("[a.after(b) a在b之后 ：%v]", m1.After(now))
	log.Printf("[a.Before(b) a在b之前 ：%v]", now.Before(m1))
	log.Printf("[a.Equal(b) a=b ：%v]", m1.Equal(now))
}
