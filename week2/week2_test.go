package main

import (
	"fmt"
	"testing"
)

func Test_Unicode(t *testing.T) {
	// 1000 1000 0110 0011
	// 0110 0111 0000 1101
	// 0101 0101 1001 1100
	// 0110 1011 0010 0010
	// 0111 1010 0111 1111
	// 0100 1110 0010 1101
	// 0101 0110 1111 1101
	// 0111 1110 1010 0010
	nums := [...]int{
		0b1000100001100011,
		0b0110011100001101,
		0b0101010110011100,
		0b0110101100100010,
		0b0111101001111111,
		0b0100111000101101,
		0b0101011011111101,
		0b0111111010100010,
	}

	for i := 0; i < len(nums); i++ {
		fmt.Printf("二进制: %016b; Unicode: %U; 字符: %c.\n", nums[i], nums[i], nums[i])
	}
}

func Test_for(t *testing.T) {
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%-d * %-d = %-d\t", j, i, i*j)
		}
		fmt.Println()
	}
}

func Test_for2(t *testing.T) {
	for i := 1; i < 10; i++ {
		for j := 1; j < i; j++ {
			var n string
			fmt.Printf("%-2s   %-2s   %-2s\t", n, n, n)
		}
		for j := i; j < 10; j++ {
			fmt.Printf("%-2d * %-2d = %-2d\t", j, i, i*j)
		}
		fmt.Println()
	}
}

func Test_for3(t *testing.T) {
	var sum int
	i := 2
	var isP bool
	for i < 101 {
		isP = true
		j := 2
		for j <= (i / j) {
			if i%j == 0 {
				// fmt.Printf("%d不是素数\n",i)
				isP = false
				break
			}
			j++
		}
		if isP {
			fmt.Printf("%d是素数\n", i)
			sum += i
		}
		i++
	}
	fmt.Println(sum)
}

func Test_range(t *testing.T) {
	text := "我爱中国"
	for n, m := range text {
		fmt.Printf("%d %c\n", n, m)
	}
	for _, m := range text {
		fmt.Printf("%c\n", m)
	}
}

func Test_goto(t *testing.T) {
	/* 定义局部变量 */
	var a int = 10

	/* 循环 */
LOOP:
	for a < 20 {
		if a == 15 {
			/* 跳过迭代 */
			a = a + 1
			goto LOOP
		}
		fmt.Printf("a的值为 : %d\n", a)
		a++
	}
}

func Test_pointer(t *testing.T) {
	var pointer01 *int
	var pointer02 *float64
	var pointer03 *string
	fmt.Printf("%T, %T ,%T\n", pointer01, pointer02, pointer03)
	fmt.Println(pointer01, pointer02, pointer03)
	fmt.Printf("%v, %v ,%v\n", pointer01, pointer02 == nil, pointer03 == nil)

	var (
		age    int     = 25
		height float64 = 1.80
		motto  string  = "少年经不得顺境，中年经不得闲境，晚年经不得逆境"
	)
	// 指针变量初始化
	pointer01, pointer02, pointer03 = &age, &height, &motto
	pointer04, pointer05, pointer06 := &age, &height, &motto
	pointer07, pointer08, pointer09 := new(int), new(float64), new(string)
	// 打印变量地址
	fmt.Println(&age, &height, &motto)
	// 打印指针变量
	fmt.Println(pointer01, pointer02, pointer03)
	fmt.Println(pointer04, pointer05, pointer06)
	fmt.Println(pointer07, pointer08, pointer09)
	// 打印指针变量访问位置存储的值
	fmt.Println(age, height, motto)
	fmt.Println(*pointer01, *pointer02, *pointer03)
	fmt.Printf("%v, %v, %q\n", *pointer01, *pointer02, *pointer03)
	fmt.Printf("%v, %v, %q\n", *pointer04, *pointer05, *pointer06)
	fmt.Printf("%v, %v, %q\n", *pointer07, *pointer08, *pointer09)
	// 通过指针变量访问修改存储的值
	*pointer01 += 1
	*pointer02 = 1.81
	*pointer03 += "--曾国藩"
	fmt.Println(age, height, motto)
	fmt.Println(&age, &height, &motto)
	fmt.Println(pointer01, pointer02, pointer03)
	fmt.Println(*pointer01, *pointer02, *pointer03)
	// 与赋值新变量对比
	age2, height2, motto2 := age, height, motto
	// 修改新变量
	age2 += 1
	height2 = 1.82
	motto2 += "家书"
	fmt.Println(age, height, motto)
	fmt.Println(age2, height2, motto2)
	fmt.Println(&age, &height, &motto)
	fmt.Println(&age2, &height2, &motto2)
	// 定义声明指针的指针
	var ppointer01 **int
	var ppointer02 **float64 = &pointer02
	ppointer03 := &pointer03
	fmt.Println(ppointer01, ppointer02, ppointer03)
	fmt.Println(&ppointer01, &ppointer02, &ppointer03)
	ppointer01 = &pointer01
	fmt.Println(ppointer01, ppointer02, ppointer03)
	fmt.Println(&ppointer01, &ppointer02, &ppointer03)
	// 通过指针的指针访问变量地址和变量值
	fmt.Println(*pointer01, *pointer02, *pointer03)
	fmt.Println(**ppointer01, **ppointer02, **ppointer03)
	// 通过指针的指针修改和变量的值
	**ppointer01 += 1
	**ppointer02 = 1.83
	**ppointer03 += "家属"
	fmt.Println(**ppointer01, **ppointer02, **ppointer03)
	fmt.Println(*pointer01, *pointer02, *pointer03)
	fmt.Println(age, height, motto)
}
