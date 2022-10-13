package debug_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/debug"
)

func TestDebug(t *testing.T) {
	// default to setup dump stack
	debug.SetupDumpStackTrap()

	// ...
	fmt.Println("starting ...")
	time.Sleep( 600 * time.Second )
}
