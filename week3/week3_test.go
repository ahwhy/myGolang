package main

import (
	"fmt"
	"testing"
	"unsafe"
)

type Person struct {
	Name          string
	Age           int
	Gender        string
	Weight        uint
	FavoriteColor []string
	NewAttr       string
	Addr          Home
}

type Home struct {
	City string
	T1   T1
}

type T1 struct {
	T1 string
}

func (p Person) Add() int {
	return p.Age * 2
}

type Author struct {
	Name string
	Aage int
}

func (a *Author) GetName() string {
	return a.Name
}

type Titile struct {
	Main string
	Sub  string
}

type Book struct {
	Author *Author
	Titile *Titile
}

func (b *Book) GetName2() string {
	return b.Author.GetName() + "book"
}

func TestMain(t *testing.T) {
	b1 := Book{
		Author: &Author{
			Name: "laoyu",
		},
		Titile: &Titile{},
	}
	b2 := &b1

	// b2 := &Book{
	// 	Author: &Author{},
	// 	Titile: &Titile{},
	// }
	*b2.Author = *b1.Author
	*b2.Titile = *b1.Titile

	b2.Author.Name = "new author"
	fmt.Println(b1.Author.Name, b2.Author.Name)
	fmt.Printf("%v\n %v", b1, b2)
}
func FucnArgsForStruct(p Person) {
	p.Name = "func for struct"
}

func FucnArgsForStructP(p *Person) {
	p.Name = "func for struct"
}

func TestFunForArg(t *testing.T) {
	p := Person{Name: "person"}
	FucnArgsForStruct(p)
	fmt.Println(p.Name)
	FucnArgsForStructP(&p)
	fmt.Println(p.Name)
}

type A struct {
	a bool
	b int32
	c string
	d string
}

type B struct {
	b int32
	c string
	d string
	a bool
}

func TestStructSize(t *testing.T) {
	a := A{true, 16, "a", "aa"}
	b := B{32, "b", "bb", true}
	fmt.Println(unsafe.Sizeof(A{}), a)
	fmt.Println(unsafe.Sizeof(B{}), b)
}
