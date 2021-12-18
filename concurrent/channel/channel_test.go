package channel_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/concurrent/channel"
)

func TestBasic(t *testing.T) {
	channel.BasicChan()
}

func TestBufferedChan(t *testing.T) {
	channel.BufferedChan()
}

func TestSyncAB(t *testing.T) {
	channel.SyncAB()
}

func TestRunTaskWithPool(t *testing.T) {
	channel.RunTaskWithPool()
}

func TestDeadLockV1(t *testing.T) {
	ch := make(chan string)

	// send
	go func() {
		ch <- "hello"
	}()

	// receive
	{
		fmt.Println(<-ch)
	}
}

func TestDeadLockV2(t *testing.T) {
	ch := make(chan string, 1)

	// send
	{
		ch <- "hello"
	}

	// receive
	{
		fmt.Println(<-ch)
	}
}
