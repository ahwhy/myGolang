package main

import "fmt"

type Person struct {
	Name  string
	Email string
	Age   int
}

type Student struct {
	Person
	StudentId int
}

// 附属于person类的方法
// 指针相当于单实例绑定
func (p *Person) SayHello() {
	fmt.Printf("[Person.SayHello][name:%s]", p.Name)
}

func main() {
	p := Person{
		Name:  "xiaoyi",
		Email: "qq.com",
		Age:   18,
	}
	s := Student{
		Person:    p,
		StudentId: 123,
	}
	s.SayHello()
}
