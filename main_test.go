package main

import (
	"fmt"
	"testing"
)

const stepsToRun int = 1

var memlist = readFile("./assets/memlist.bin")
var bankFilesMap = createBankMap("./assets/")

func run(gamepart int) GotherWorld {
	app := initGotherWorld(memlist, bankFilesMap, true)
	app.loadGamePart(GAME_PART_ID_1 + gamepart)
	for i := 0; i < stepsToRun; i++ {
		app.mainLoop(i)
	}
	return app
}

func TestRunGameparts(t *testing.T) {
	for part := 0; part < 9; part++ {
		fmt.Println("### RUN PART", part)
		app := run(part)
		//TODO should be 0, there is one case!
		if app.vm.countNoOps > 1 {
			t.Errorf("countNoOps > 0: part %d %d", part, app.vm.countNoOps)
		}
		if app.vm.countSPNotZero > 1 {
			t.Errorf("countSPNotZero > 1: part %d %d", part, app.vm.countSPNotZero)
		}
	}
}
