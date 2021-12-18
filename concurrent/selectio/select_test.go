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
