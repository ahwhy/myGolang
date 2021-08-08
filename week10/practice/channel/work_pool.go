package chnanel

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Task struct {
	ID         int
	JobID      int
	Status     string
	CreateTime time.Time
}

func (t *Task) Run() {
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	t.Status = "Completed"
}

// worker的数量，即使用多少goroutine执行任务
const workerNum = 4

func RunTaskWithPool() {
	wg.Add(workerNum)

	// 创建容量为10的buffered channel
	taskQueue := make(chan *Task, 10)

	// 激活goroutine，执行任务
	for workID := 1; workID <= workerNum; workID++ {
		go worker(taskQueue, workID)
	}

	// 生成消息
	produceTask(taskQueue)

	// 5秒后 关闭管道，通知所有worker退出
	time.Sleep(5 * time.Second)
	close(taskQueue)

	wg.Wait()
}

func produceTask(out chan<- *Task) {
	// 将待执行任务放进buffered channel，共15个任务
	for i := 1; i <= 15; i++ {
		fmt.Println(i)
		out <- &Task{
			ID:         i,
			JobID:      100 + i,
			CreateTime: time.Now(),
		}

	}
}

// 从buffered channel中读取任务，并执行任务
func worker(in <-chan *Task, workID int) {
	defer wg.Done()
	for v := range in {
		fmt.Printf("Worker%d: recv a request: TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
		v.Run()
		fmt.Printf("Worker%d: Completed for TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
	}
}
