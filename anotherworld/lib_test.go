package anotherworld

import (
	"fmt"
	"testing"
)

const stepsToRun int = 10

var memlist = ReadFile("../assets/memlist.bin")
var bankFilesMap = CreateBankMap("../assets/")

func run(gamepart int) GotherWorld {
	videoDriver := Video{Hal: DummyHAL{}}
	app := InitGotherWorld(memlist, bankFilesMap, videoDriver)
	app.LoadGamePart(GAME_PART_ID_1 + gamepart)
	for i := 0; i < stepsToRun; i++ {
		app.MainLoop(i)
	}
	return app
}

func TestRunGameparts(t *testing.T) {
	for part := 0; part < 9; part++ {
		fmt.Println("### RUN PART", part)
		app := run(part)
		if app.Vm.CountNoOps > 0 {
			t.Errorf("CountNoOps > 0: part %d %d", part, app.Vm.CountNoOps)
		}
		if app.Vm.CountSPNotZero > 0 {
			t.Errorf("CountSPNotZero > 1: part %d %d", part, app.Vm.CountSPNotZero)
		}
	}
}
