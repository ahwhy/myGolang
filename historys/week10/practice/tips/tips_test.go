package tips

import "testing"

func TestCancelWithChannel(t *testing.T) {
	CancelWithChannel()
}

func TestCancelWithDown(t *testing.T) {
	CancelWithDown()
}

func TestCancelWithCtx(t *testing.T) {
	CancelWithCtx()
}

func TestDealPanic(t *testing.T) {
	DealPanic()
}

func TestDealPanicInG(t *testing.T) {
	DealPanicInG()
}

func TestDealPanicInGV2(t *testing.T) {
	DealPanicInGV2()
}

func TestCallBackMode(t *testing.T) {
	CallBackMode()
}
