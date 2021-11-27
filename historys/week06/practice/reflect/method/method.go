package main

import (
	"log"
	"reflect"
)

type Person struct {
	Name   string
	Age    int
	Gender string
}

func (p Person) ReflectCallFuncWithArgs(name string, age int) {
	log.Printf("[调用的是带参数的方法][args.name:%s][args.age:%d][[p.name:%s][p.age:%d]",
		name,
		age,
		p.Name,
		p.Age,
	)
}

func (p Person) ReflectCallFuncWithNoArgs() {
	log.Printf("[调用的是不带参数的方法]")
}

func main() {
	p1 := Person{
		Name:   "小乙",
		Age:    18,
		Gender: "男",
	}

	// 1. 首先通过 reflect.ValueOf(p1)获取 得到反射值类型
	getValue := reflect.ValueOf(p1)

	// 2. 带参数的方法调用
	methodValue := getValue.MethodByName("ReflectCallFuncWithArgs")
	// 参数是reflect.Value的切片
	args := []reflect.Value{reflect.ValueOf("李逵"), reflect.ValueOf(30)}
	methodValue.Call(args)

	// 3. 不带参数的方法调用
	methodValue = getValue.MethodByName("ReflectCallFuncWithNoArgs")
	// 参数是reflect.Value的切片
	args = make([]reflect.Value, 0)
	methodValue.Call(args)
}
