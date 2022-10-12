package debug_test

import (
	"testing"
	"time"

	"github.com/ahwhy/myGolang/debug"
)

func TestDebug(t *testing.T) {
	// default to setup dump stack
	debug.SetupDumpStackTrap()

	// ...
	time.sleep( 5 * time.Second )
}