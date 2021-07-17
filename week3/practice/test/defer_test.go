package main

import (
	"fmt"
	"testing"
)

func TestDefer1(t *testing.T) {
	fmt.Println("func start")
	x := 10
	defer func(x int) {
		fmt.Println("in defer: ", x)
	}(x)
	x = 30
	fmt.Println("func end: ", x)
}

func TestDefer2(t *testing.T) {
	fmt.Println("func start")
	x := 10
	defer func(x *int) {
		fmt.Println("in defer: ", *x)
	}(&x)
	x = 30
	fmt.Println("func end: ", x)
}

func TestDefer3(t *testing.T) {
	fmt.Println("func start")
	x := 10
	defer func() {
		fmt.Println("in defer: ", x)
	}()
	x = 30
	fmt.Println("func end: ", x)
}

func TestDefer4(t *testing.T) {
	fmt.Println("start")
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)
	defer fmt.Println(4)
	fmt.Println("end")
}
