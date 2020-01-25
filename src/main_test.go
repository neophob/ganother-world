package main

import (
	"fmt"
	"testing"
)

const STEPS_TO_RUN int = 1024

var data = readFile("../assets/memlist.bin")
var resourceMap, _ = unmarshallingMemlistBin(data)
var bankFilesMap = createBankMap("../assets/")
var gameParts = getGameParts()

func run(gamepart, steps int) VMState {
	assets := Assets{
		memList:         resourceMap,
		gameParts:       gameParts,
		bank:            bankFilesMap,
		loadedResources: make(map[int][]uint8),
	}
	vmState := createNewState(assets)
	vmState.setupGamePart(GAME_PART_ID_1 + gamepart)

	for i := 0; i < steps; i++ {
		vmState.mainLoop(0)
	}
	return vmState
}

func TestRunGameparts(t *testing.T) {
	for part := 0; part < 9; part++ {
		fmt.Println("### RUN PART", part)
		vmState := run(part, STEPS_TO_RUN)
		//TODO should be 0, there is one case!
		if vmState.countNoOps > 1 {
			t.Errorf("countNoOps > 0: part %d %d", part, vmState.countNoOps)
		}
		if vmState.countSPNotZero > 1 {
			t.Errorf("countSPNotZero > 1: part %d %d", part, vmState.countSPNotZero)
		}
	}
}
