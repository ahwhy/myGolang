package chnanel

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	Basic()
}

func TestBufferedChan(t *testing.T) {
	BufferedChan()
}

func TestSyncAB(t *testing.T) {
	SyncAB()
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

func TestRunTaskWithPool(t *testing.T) {
	RunTaskWithPool()
}
