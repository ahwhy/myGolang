package selectio_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/concurrent/selectio"
)

func TestBasic1(t *testing.T) {
	fmt.Println("don't cache")
	selectio.Basic1()
}

func TestBasic2(t *testing.T) {
	fmt.Println("don't cache")
	selectio.Basic2()
}

func TestSelectChannel(t *testing.T) {
	fmt.Println("don't cache")
	selectio.SelectChannel()
}

func TestSelectOrder(t *testing.T) {
	fmt.Println("don't cache")
	selectio.SelectOrder()
}

func TestTimeAfter(t *testing.T) {
	selectio.TimeAfter()
}

func TestSelectTimeout(t *testing.T) {
	selectio.SelectTimeout()
}

func TestSelectTimeoutV2(t *testing.T) {
	selectio.SelectTimeoutV2()
}

func TestCancelWithChannel(t *testing.T) {
	selectio.CancelWithChannel()
}

func TestCancelWithDown(t *testing.T) {
	selectio.CancelWithDown()
}

func TestCancelWithCtx(t *testing.T) {
	selectio.CancelWithCtx()
}

func TestGraceful_exit(t *testing.T) {
	selectio.Graceful_exit()
}
