package main

import "fmt"

// 体现多态
// 告警通知的函数，根据不同的对象通知
// 有个共同的通知方法，每种对象自己实现

type notifer interface {
	Init() // 动作，定义的方法
	push()
	notify()
}

type user struct {
	name  string
	email string
}

type admin struct {
	name string
	age  int
}

func (u *user) Init()  {}
func (u *admin) Init() {}

func (u *user) push() {
	fmt.Printf("[普通用户][sendNotify to user %s]\n", u.name)
}
func (u *admin) push() {
	fmt.Printf("[管理员][sendNotify to user %s]\n", u.name)
}

func (u *user) notify() {
	fmt.Printf("[普通用户][sendNotify to user %s]\n", u.name)
}
func (u *admin) notify() {
	fmt.Printf("[管理员][sendNotify to user %s]\n", u.name)
}

// 多态的统一调用方法，入口
func sendNotify(n notifer) {
	n.notify()
	n.push()
}

func main() {
	u1 := user{
		name:  "aa",
		email: "qq.com",
	}
	a1 := admin{
		name: "bb",
		age:  18,
	}
	u1.push()
	a1.push()
	u1.notify()
	a1.notify()

	var n notifer
	n = &u1
	n.push()
	n.notify()
	n = &u1
	n.push()
	n.notify()

	ns := make([]notifer, 0)
	ns = append(ns, &u1, &a1)
	for _, n := range ns {
		sendNotify(n)
	}
}
