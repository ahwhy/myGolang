package main

import (
	"log"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person     //匿名结构体嵌套
	StudentId  int
	SchoolName string
	IsBaoSong  bool //是考上来的吗
	Hobbies    []string
	// panic: reflect.Value.Interface: cannot return value obtained from unexported field or method
	// hobbies    []string
	Labels map[string]string
}

//func (s *Student) GoHome() {
//	log.Printf("[回家了][sid:%d]", s.StudentId)
//}
func (s *Student) GoHome() {
	log.Printf("[回家了][sid:%d]", s.StudentId)
}

func (s Student) GotoSchool() {
	log.Printf("[去上学了][sid:%d]", s.StudentId)
}

func (s Student) Baosong() {
	log.Printf("[竞赛保送][sid:%d]", s.StudentId)
}

func main() {
	s := Student{
		Person:     Person{Name: "xiaoyi", Age: 9900},
		StudentId:  123,
		SchoolName: "五道口皇家男子职业技术学院",
		IsBaoSong:  true,
		Hobbies:    []string{"唱", "跳", "Rap"},
		//hobbies:    []string{"唱", "跳", "Rap"},
		Labels: map[string]string{"k1": "v1", "k2": "v2"},
	}
	p := Person{
		Name: "李逵",
		Age:  124,
	}

	reflectProbeStruct(s)
	reflectProbeStruct(p)
}

func reflectProbeStruct(s interface{}) {

	// 获取目标对象
	t := reflect.TypeOf(s)
	log.Printf("[对象的类型名称：%s]", t.Name())

	// 获取目标对象的值类型
	v := reflect.ValueOf(s)
	// 遍历获取成员变量
	for i := 0; i < t.NumField(); i++ {
		// Field 代表对象的字段名
		key := t.Field(i)
		value := v.Field(i).Interface()
		//
		if key.Anonymous {
			log.Printf("[匿名字段][第:%d个字段][字段名:%s][字段的类型:%v][字段的值:%v]", i+1, key.Name, key.Type, value)
		} else {
			log.Printf("[命名字段][第:%d个字段][字段名:%s][字段的类型:%v][字段的值:%v]", i+1, key.Name, key.Type, value)

		}

	}

	// 打印方法
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		log.Printf("[第:%d个方法][方法名称:%s][方法的类型:%v]", i+1, m.Name, m.Type)

	}
}
