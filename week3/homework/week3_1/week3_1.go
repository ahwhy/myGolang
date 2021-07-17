package main

import (
	"fmt"
	"unsafe"
)

type Book struct {
	Title  string
	Author string
	Page   uint
	Tag    []string
}

const (
	A = unsafe.Offsetof(Book{}.Tag)
)

func main() {
	b := &Book{Tag: []string{"abc", "def", "hjk"}}
	// 根据结构体切片中第一个元素的内存地址, 计算出Tag的内存地址, 并访问
	Tag_array := unsafe.Pointer(&b.Tag[0])
	Tag_Sizeof := unsafe.Sizeof(b.Tag[0])
	b_Tag_0 := *(*string)(unsafe.Pointer(uintptr(Tag_array) + Tag_Sizeof*0))
	b_Tag_1 := *(*string)(unsafe.Pointer(uintptr(Tag_array) + Tag_Sizeof*1))
	b_Tag_2 := *(*string)(unsafe.Pointer(uintptr(Tag_array) + Tag_Sizeof*2))
	fmt.Println(b_Tag_0, b_Tag_1, b_Tag_2)

	// 根据结构体的内存地址, 计算出Tag的内存地址, 并访问
	b_addr := unsafe.Pointer(b)
	b_Tag := *(*[]string)(unsafe.Pointer(uintptr(b_addr) + A)) //  unsafe.Pointer(&b.Tag)
	fmt.Println(b_Tag)

	// 群里GO5095这位大佬教的骚操作
	b_Tag_0_addr := *(&(*(*[]string)(unsafe.Pointer(uintptr(b_addr) + A)))[0])
	b_Tag_1_addr := *(&(*(*[]string)(unsafe.Pointer(uintptr(b_addr) + A)))[1])
	b_Tag_2_addr := *(&(*(*[]string)(unsafe.Pointer(uintptr(b_addr) + A)))[2])
	fmt.Println(b_Tag_0_addr, b_Tag_1_addr, b_Tag_2_addr)
}
