package controller_test

import (
	"bytes"
	"testing"

	"github.com/ahwhy/myGolang/promethous/script_exportor/controller"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	sc *controller.ScriptCollector
)

func init() {
	zap.DevelopmentSetup()

	sc = controller.NewScriptCollector("../modules")
}

func TestScirptCollector(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	if err := sc.Exec("test.sh", "-f", b); err != nil {
		t.Fatal(err)
	}

	t.Log(b.String())
}
