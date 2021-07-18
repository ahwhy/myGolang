package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/week1/compile/pkg"
)

func Test_holle(t *testing.T) {
	pkg.Demo()
	fmt.Println("Holle World!\nThis is my first golang program!")
}

func Test_var(t *testing.T) {
	var name string = "atlantis"
	// 行注释
	/*
	   块注释
	*/

	// 全局变量
	// var定义变量
	// var 标识符（变量名称） 类型（变量的类型）
	// string 表示字符串
	// 若未设置则用零值进行初始化
	name = "atlantis2"
	/*
	   1、var flag type
	   2、var flag type = value
	   3、var flag = value
	*/
	// 局部变量
	// 标识符在局部内只能定义一次
	var (
		age    int = 25
		weight     = 145
	)
	fmt.Println(name, age)
	// 作用域 说明标识符的使用范围 {}
	{
		name = "atlantis3"
	}

	// height := 180 // 短声明 var height = 180
	// 短声明只能用在函数内部
	height, weight := 180, 145

	/*
		        var age, weight, height int
				var age, weight, height int = 1， 2， 3
				var age, weight, height  = 1， ""， 3

	*/

	fmt.Println(name, age, weight, height)                                    // 打印内容后会自动换行
	fmt.Print(name, age, weight, height)                                      // 打印内容后不会加换行符
	fmt.Printf("\n我叫%s,我的年龄是%d，我的体重是%T，我的身高是%d\n", name, age, weight, height) // 通过占位符进行标量填充
}

func Test_const(t *testing.T) {
	// 定义常量(常量需要初始化值)
	const statusNew int = 1
	const statusDeleted int = 2

	const (
		Monday = 10
		Tuesday
	)
	/*
		若未赋值，则使用最近的一个已使用的常量进行赋值
	*/
	fmt.Println(statusNew, statusDeleted)
	fmt.Println(Monday, Tuesday)

	/*
		枚举值
		iota 在一个小括号内，初始化为0，每调用一次+1
			statusa = iota // 0
			statusb = iota // 1
			statusc = iota // 2
			statusd = iota // 3
	*/
	const (
		statusa = iota
		statusb
		statusc
		statusd
	)
	fmt.Println(statusa, statusb, statusc, statusd)

	const (
		status1 = iota * 100
		status2
		status3
		status4
	)
	fmt.Println(status1, status2, status3, status4)
}

func Test_ioread(t *testing.T) {
	var name string
	fmt.Print("请输入你的名字:")
	fmt.Scan(&name)
	fmt.Println("你输入的内容是:", name)

	var age int
	fmt.Print("请输入你的年龄:")
	fmt.Scan(&age)
	fmt.Println("你输入的内容是:", age)
}

func Test_complex(t *testing.T) {
	var (
		c1 complex64 = 1 + 2i
		c2 complex64 = complex(3, 4)
	)
	fmt.Println(c1 + c2)
	fmt.Println(real(c1), imag(c1))
}

func Test_random(t *testing.T) {
	rand.Seed(time.Now().Unix()) // 设置随机数种子
	fmt.Println(rand.Intn(100))
}
