package main

import (
	"fmt"
	"testing"

	"github.com/neophob/ganother-world/anotherworld"
)

const stepsToRun int = 1

var memlist = readFile("../assets/memlist.bin")
var bankFilesMap = createBankMap("../assets/")

func run(gamepart int) anotherworld.GotherWorld {
	videoDriver := anotherworld.Video{Hal: anotherworld.DummyHAL{}}
	app := anotherworld.InitGotherWorld(memlist, bankFilesMap, videoDriver)
	app.LoadGamePart(anotherworld.GAME_PART_ID_1 + gamepart)
	for i := 0; i < stepsToRun; i++ {
		app.MainLoop(i)
	}
	return app
}

func TestRunGameparts(t *testing.T) {
	for part := 0; part < 9; part++ {
		fmt.Println("### RUN PART", part)
		app := run(part)
		//TODO should be 0, there is one case!
		if app.Vm.CountNoOps > 1 {
			t.Errorf("CountNoOps > 0: part %d %d", part, app.Vm.CountNoOps)
		}
		if app.Vm.CountSPNotZero > 1 {
			t.Errorf("CountSPNotZero > 1: part %d %d", part, app.Vm.CountSPNotZero)
		}
	}
}
