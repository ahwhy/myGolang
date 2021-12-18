package channel_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/historys/week10/practice/channel"
)

func TestBasic(t *testing.T) {
	channel.Basic()
}

func TestBufferedChan(t *testing.T) {
	channel.BufferedChan()
}

func TestSyncAB(t *testing.T) {
	channel.SyncAB()
}

func TestDeadLockV1(t *testing.T) {
	ch := make(chan string)
	
	// send
	{
		ch <- "hello"
	}

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

func TestRunTaskWithPool(t *testing.T) {
	channel.RunTaskWithPool()
}