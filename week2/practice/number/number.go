// 猜数字小游戏
// 随机生成数字0-100
// 从控制台数据与生成数字比较
// 大 提示大了
// 小 提示小了
// 等于 成功，程序结束
// 最多猜五次，未猜对说太笨了，程序结束
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	num := rand.Intn(100)

	var (
		num2 int
		targ bool
	)

	for i := 1; i < 6; i++ {
		fmt.Println("Let play a game,please input a number：")
		fmt.Scan(&num2)
		if num == num2 {
			fmt.Println("You Win!!\nGame is over!!!")
			break
		} else if num > num2 {
			fmt.Println("The number is small")
		} else if num < num2 {
			fmt.Println("The number is big")
		}
		if i == 5 {
			targ = true
		}
	}

	if targ {
		fmt.Println("You don't have a chance!\nYou're too stupid!!\nGame is over!!!")
	}
	
	fmt.Println("Thank You!")
	os.Exit(1)
}
