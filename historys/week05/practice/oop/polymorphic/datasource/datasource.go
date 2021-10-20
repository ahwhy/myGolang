package main

import (
	"fmt"
	"log"
)
/*
1. 多个数据源有
2. query方法做查数据
3. pushdata方法做写入数据
*/
// 方法集合
type DataSource interface {
	PushData(data string)
	QueryData(name string) string
}

type redis struct {
	Name string
	Addr string
}

func (r *redis) PushData(data string) {
	log.Printf("[PushData][ds.name:%s][data:%s]", r.Name, data)
}
func (r *redis) QueryData(name string) string {
	log.Printf("[QueryData][ds.name:%s][data:%s]", r.Name, name)
	return name + "_redis"
}

type kafka struct {
	Name string
	Addr string
}

func (k *kafka) PushData(data string) {
	log.Printf("[PushData][ds.name:%s][data:%s]", k.Name, data)
}
func (k *kafka) QueryData(name string) string {
	log.Printf("[QueryData][ds.name:%s][data:%s]", k.Name, name)
	return name + "_kafka"
}

var Dm = make(map[string]DataSource)

func main() {
	r := redis{
		Name: "redis",
		Addr: "1.1",
	}
	k := kafka{
		Name: "kafka",
		Addr: "2.2",
	}
	// 注册数据源到承载的容器中
	Dm["redis"] = &r
	Dm["kafka"] = &k
	// 推送数据
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		for _, ds := range Dm {
			ds.PushData(key)
		}
	}
	// 查询数据
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		for _, ds := range Dm {
			res := ds.QueryData(key)
			log.Println(res)
		}
	}
}
