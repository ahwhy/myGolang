package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
)

func TestString_ascll(t *testing.T) {
	// ascii码 明确给字符指定byte类型 代表 ascii码
	var ch1 byte = 'a'
	// 字符范围还在ascii码内，但是不明确byte，int32 =rune类型
	var ch2 = 'a'
	// int32 =rune类型
	var ch3 = '你'

	fmt.Printf("字符 指定ascii：%c id：%v 实际类型：%T \n", ch1, ch1, ch1)
	fmt.Printf("字符 ：%c id：%v 实际类型：%T \n", ch2, ch2, ch2)
	fmt.Printf("字符 中文：%c id：%v 实际类型：%T \n", ch3, ch3, ch3)
}

func TestString_Len(t *testing.T) {
	ch1 := "ma ge jiao yu"
	ch2 := "马哥教育"
	ch3 := "m哥"
	fmt.Printf("字符串%v 字节大小or字符长度%d  真实字符长度%d\n", ch1, len(ch1), utf8.RuneCountInString(ch1))
	fmt.Printf("字符串%v 字节大小or字符长度%d  真实字符长度%d\n", ch2, len(ch2), utf8.RuneCountInString(ch2))
	fmt.Printf("字符串%v 字节大小or字符长度%d  真实字符长度%d\n", ch3, len(ch3), utf8.RuneCountInString(ch3))
	/*
		字符串ma ge jiao yu 字节大小or字符长度13  真实字符长度13
		字符串马哥教育 字节大小or字符长度12  真实字符长度4
		字符串m哥 字节大小or字符长度4  真实字符长度2
	*/
}

func TestString_Range(t *testing.T) {
	// - 如果是ASCII字符：直接使用下标遍历
	// - 如果是unicode字符遍历：使用for range

	ch1 := "ma ge马哥教育"
	for i := 0; i < len(ch1); i++ {
		fmt.Printf("ascii: %c %d\n", ch1[i], ch1[i])
	}
	for _, i := range ch1 {
		fmt.Printf("unicode: %c %d \n", i, i)
	}
}

func TestString_Add(t *testing.T) {
	// - 使用+拼接多个字符串
	// - 支持换行
	s1 := "http://"
	s2 := "localhost:8080"
	s3 := s1 + s2
	fmt.Println(s3)
	s4 := "http://localhost:8080/api/v1" +
		"/login"
	fmt.Println(s4)
}

func TestString_Modify(t *testing.T) {
	// 字符串修改：通过 []byte和string转换 创建新的字符串达到
	// 举例 8080 改为8081
	s2 := "localhost:8080"
	fmt.Println(s2)
	sByte := []byte(s2)
	sByte[len(sByte)-1] = '1'
	s3 := string(sByte)
	fmt.Println(s3)
}

func TestString_Strings1(t *testing.T) {
	// Strings包
	s1 := "inf.bigdata.kafka"
	s2 := "localhost:8080/api/v1/host/1"

	serviceS1 := strings.Split(s1, ".")
	// SplitAfter 保留sep
	serviceS2 := strings.SplitAfter(s1, ".")

	pathS := strings.Split(s2, "/")
	// SplitN 结果长度为n，如果没切完就不切了，保留给最后一个
	pathSN := strings.SplitN(s2, "/", 2)
	fmt.Printf("[切割服务标识]%v\n", serviceS1)
	fmt.Printf("[切割服务标识][after]%v\n", serviceS2)
	fmt.Printf("[切割uri][]%v\n", pathS)
	fmt.Printf("[切割uri][SplitN]%v\n", pathSN)

	fmt.Println(strings.HasPrefix("inf.bigdata.kafka", "inf"))
	fmt.Println(strings.HasSuffix("inf.bigdata.kafka", "kafka"))
	fmt.Println(strings.HasSuffix("inf.bigdata.kafka", ""))
}

func TestString_Strings2(t *testing.T) {
	want := `
	[报警触发类型：%s]
	[报警名称：%s]
	[级别：%d级]	
	[机器ip列表：%s]	
	[表达式：%s]
	[最大报警次数：%d次]
	[触发时间：%s]
`
	alarmContent := fmt.Sprintf(
		want,
		"prometheus",
		"登录接口qps 大于100",
		3,
		"1.1.1.1, 2.2.2.2",
		`sum(rate(login_qps[1m])) >100`,
		2,
		time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
	)
	fmt.Println(alarmContent)

	ports := []int{8000, 8001, 8002}
	baseUri := "/api/v1/push"
	for _, p := range ports {
		uri := fmt.Sprintf("localhost:%d%s", p, baseUri)
		fmt.Println(uri)
	}
}

func TestString_Strings3(t *testing.T) {
	x := "@@@@@abchello_mage@"
	// 会去掉cutset连续的
	fmt.Println(strings.Trim(x, "@"))
	fmt.Println(strings.TrimLeft(x, "@"))
	fmt.Println(strings.TrimRight(x, "@"))
	fmt.Println(strings.TrimSpace(x))
	fmt.Println(strings.TrimPrefix(x, " abc"))
	fmt.Println(strings.TrimSuffix(x, " @"))
	f := func(r rune) bool {
		return unicode.Is(unicode.Han, r) //如果是汉字返回true
	}
	fmt.Println(strings.TrimFunc("你好啊abc", f))

	// TrimLeft会去掉连续的cutset bchello_mage@
	fmt.Println(strings.TrimLeft(x, "@a"))
	// TrimPrefix会去掉的单一的 @abchello_mage@
	fmt.Println(strings.TrimPrefix(x, "@a"))
}

func TestString_Join(t *testing.T) {
	baseUri := "http://localhost:8080/api/v1/query?"
	args := strings.Join([]string{"name=mage", "id=1", "env=online"}, "&")
	fullUri := baseUri + args
	fmt.Println(fullUri)
	// http://localhost:8080/api/v1/query?name=mage&id=1&env=online
}

func TestString_Add2(t *testing.T) {
	ss := []string{
		"A",
		"B",
		"C",
	}
	var b strings.Builder
	for _, s := range ss {
		b.WriteString(s)
	}
	fmt.Println(b.String())
}
