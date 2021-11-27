package main

import (
	"fmt"
)

type A struct {
	Name string
}

func (a *A) test1() {
	a.Name = "22" // 直接引用结构体字段，可以修改Name
}

func (a *A) test2() {
	a = &A{"33"} // 虽然 a 是指针，但这个表达式，修改的是指针的指向，无法修改Name
}

func (a *A) test3() {
	*a = A{"44"} // 修改指针指向的内存地址的值，可以修改Name
}

func main() {
	a := &A{"11"}
	fmt.Println(a.Name) // 11
	a.test1()
	fmt.Println(a.Name) // 22
	a.test2()
	fmt.Println(a.Name) // 22
	a.test3()
	fmt.Println(a.Name) // 44
}
