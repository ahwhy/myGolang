package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/x-mod/routine"
)

func prepare(ctx context.Context) error {
	fmt.Println("准备工作开始")
	fmt.Println("初始化全局变量")
	fmt.Println("创建数据库连接")
	fmt.Println("准备工作结束")
	return nil
}

func cleanup(ctx context.Context) error {
	fmt.Println("清理工作开始")
	fmt.Println("关闭各种文件句柄")
	fmt.Println("关闭数据库连接")
	fmt.Println("清理工作结束")
	return nil
}

func bizWork(ctx context.Context) error {
	fmt.Println("开始处理用户请求")
	defer fmt.Println("用户请求处理完毕")
	go func() {
		i := 0
		for {
			fmt.Printf("处理第%d个用户的请求\n", i)
			time.Sleep(1 * time.Second)
			i++
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func daemonWork(ctx context.Context) error {
	fmt.Println("开启一个守护协程")
	defer fmt.Println("守护协程结束")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return nil
	}
}

func main21() {
	ctx, cancle := context.WithCancel(context.Background())
	routine.Main(
		ctx,
		routine.ExecutorFunc(bizWork), //主要工作
		routine.Signal(syscall.SIGINT, routine.SigHandler(func() {
			fmt.Printf("got signal %d\n", syscall.SIGINT)
			cancle()
		})), //接收到终止信号时优雅地退出
		routine.Prepare(routine.ExecutorFunc(prepare)), //初始化
		routine.Cleanup(routine.ExecutorFunc(cleanup)), //收尾清理
		routine.Go(routine.ExecutorFunc(daemonWork)),   //跟主要工作并行的其他工作
	)
}
