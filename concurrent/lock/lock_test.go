package lock_test

import (
	"testing"

	"github.com/ahwhy/myGolang/concurrent/lock"
)

func TestProblemV1(t *testing.T) {
	lock.ProblemV1()
}

func TestProblemV2(t *testing.T) {
	lock.ProblemV2()
}

func TestRWLock(t *testing.T) {
	lock.RWLock()
}
