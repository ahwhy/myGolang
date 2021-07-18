# day07重点总结
- 工程组织很重要：多去看开源项目
- 单测比基准更重要
- 时间处理属于基础工具
- 日志处理非常重要

# day08讲什么
- 作业
- 日志logrotate 
- pprof
- 加解密
- 数据结构和算法

# 压测作业
- 字符串拼接benchmark
- 对比以下几个方法
```go
func xxx(n int ,str string) string{
}
```

- 使用+拼接
- 使用fmt.sprinf拼接
- 使用strings.builder拼接
- 使用bytes.buffer拼接
- 使用[]byte拼接
- 生成基本字符串的函数如下
```shell script
package main

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

```


# 发送钉钉这个logrus hook 
- 改造成异步发送的