package main

import (
	"fmt"
	"log"
	"time"
	"github.com/patrickmn/go-cache"
)

var (
	defaultInterval = time.Second * 30
	UserCache       = cache.New(defaultInterval, defaultInterval)
)

type user struct {
	Name  string
	Email string
	Phone int64
}

func GetUser(name string) user {
	res, found := UserCache.Get(name)
	if found {
		u := res.(user)
		log.Printf("[found_user_in_cache][name:%s][value:%v]", name, u)
		return u
	} else {
		res := HttpGetUser(name)
		// 给每个key 单独设置超时时间
		UserCache.Set(name, res, defaultInterval)
		log.Printf("[not_found_user_in_cache][query_by_http][name:%s][value:%v]", name, res)
		return res
	}
}

func HttpGetUser(name string) user {
	//这里是去 接口中拿user
	u := user{
		Name:  name,
		Email: "qq.com",
		Phone: time.Now().Unix(),
	}
	return u
}

func queryUser() {
	for i := 0; i < 10; i++ {
		userName := fmt.Sprintf("user_name_%d", i)
		GetUser(userName)
	}
}

func main() {
	//先查一波
	log.Printf("[第一波查询][缓存数据应该都不存在][去db or http查询][设置缓存]")
	queryUser()
	// 再查一波
	log.Printf("[第二波查询][缓存数据存在][直接返回]")
	queryUser()
	time.Sleep(31 * time.Second)
	// 再查一波
	log.Printf("[第三波查询][缓存数据已失效][去db or http查询][设置缓存]")
	queryUser()
}
