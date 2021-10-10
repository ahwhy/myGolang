// 抽卡程序, 中奖率万分之一
// 请回答如下问题, 比如: a + b =:  m
// 恭喜你 回答正确，随机获取n[1~10]的抽奖机会
// 当前是你第1次抽奖: 抽奖结果 未中奖
// 当前是你第2次抽奖. 抽奖结果 未中奖
// 当前是你第3次抽奖. 抽奖结果 未中奖
// ...

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	var num int
	a := rand.Intn(10)
	b := rand.Intn(10)
	c := rand.Intn(10) + 1

	fmt.Printf("Please answer a question:%d + %d = ", a, b)
	fmt.Scan(&num)
	if a+b == num {
		fmt.Println("恭喜你回答正确，将会随机获取[1~10]的抽奖机会")
		time.Sleep(1000)
		fmt.Printf("恭喜，你获得 %d 次抽奖机会\n", c)
		draw(c)
	} else {
		fmt.Println("回答错误，下次努力")
		os.Exit(1)
	}
}

func draw(num int) {
	lucky_Number := rand.Intn(10000)
	
	for i := 1; i < num+1; i++ {
		num1 := rand.Intn(10000)
		if lucky_Number == num1 {
			fmt.Println("恭喜你,中奖了！！！")
			os.Exit(1)
		} else {
			fmt.Printf("很遗憾你没有中奖，这是你第 %d 次机会，你还有 %d 次机会\n", i, num-i)
		}
	}
}
