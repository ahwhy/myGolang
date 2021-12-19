package main

import (
	"context"
	"fmt"
	"time"

	"github.com/x-mod/routine"
)

func mockWork(ctx context.Context) error {
	time.Sleep(200 * time.Millisecond)
	fmt.Println("execute mock work")
	return nil
}

func main22() {
	ctx := context.Background()
	exec := routine.ExecutorFunc(mockWork)

	timeout := routine.Timeout(100*time.Millisecond, exec)
	if err := timeout.Execute(ctx); err != nil {
		fmt.Println(err) //不符合预期，虽然会打印超时错误，但还是会等mockWork任务执行完毕
	}
	fmt.Println("=====================")
	retry := routine.Retry(3, exec) //最多重试3次，只要成功一次就不再重试了
	if err := retry.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	repeat := routine.Repeat(3, time.Second, exec) //重复执行3次，间隔1秒
	if err := repeat.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	concurrent := routine.Concurrent(4, exec) //同一个任务开多个协程并行执行
	if err := concurrent.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	parallel := routine.Parallel(exec, exec, exec) //开多个协程并行执行不同的任务
	if err := parallel.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	sequence := routine.Append(exec, exec, exec) //串行执行多个任务
	if err := sequence.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	command := routine.Command("echo", routine.ARG("hello MAGE")) //执行一个shell命令
	if err := command.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
	crontab := routine.Crontab("* * * * *", exec) //go的crontab有5位，依次表示：分、时、日、月、星期，5个*表示每分钟执行一次
	if err := crontab.Execute(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====================")
}
