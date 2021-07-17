package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func Add(a int, b int) int {
	return a + b
}

func Sub(a int, b int) int {
	return a - b
}

func Mul(a int, b int) int {
	return a * b
}

func Div(a int, b int) int {
	return a / b
}

func TestFuncslince(t *testing.T) {
	var fadd func(int, int) int = Add

	var fs []func(int, int) int
	fs = append(fs, fadd, Add, Sub, Mul, Div)
	fmt.Printf("%T;\n%#v\n", fs, fs)

	for _, f := range fs {
		fmt.Println(f(4, 2))
	}
}

func Print(pf func(a string) string, names ...string) {
	for i, v := range names {
		fmt.Println(i, pf(v))
	}
}

func printfmt(name string) string {
	return "*" + name + "*"
}

func TestFone(t *testing.T) {
	name := []string{"aa", "bb", "cc"}
	Print(printfmt, name...)
}

func sayHi() {
	fmt.Println("sayHi")
}

func sayHolle() {
	fmt.Println("sayHolle")
}

func genFunc() func() {
	if rand.Int()%2 == 0 {
		return sayHi
	} else {
		return sayHolle
	}
}

func aFields(split rune)bool{
	if  split == 'a' {
		return true
	}
	return false
}

func TestFreturn(t *testing.T) {
	rand.Seed(time.Now().Unix())
	f := genFunc()
	f()
	fmt.Printf("%q\n", strings.FieldsFunc("aasdfaasdfdsfsdafasdf",aFields))
}
