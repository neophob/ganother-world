package main

import (
	"fmt"
	"testing"
)

const STEPS_TO_RUN int = 2048

var data = readFile("../assets/memlist.bin")
var resourceMap, _ = unmarshallingMemlistBin(data)
var bankFilesMap = createBankMap("../assets/")
var gameParts = getGameParts()

func run(gamepart, steps int) {
	assets := Assets{
		memList:         resourceMap,
		gameParts:       gameParts,
		bank:            bankFilesMap,
		loadedResources: make(map[int][]uint8),
	}
	vmState := createNewState(assets)
	vmState.setupGamePart(GAME_PART_ID_1 + gamepart)

	for i := 0; i < steps; i++ {
		vmState.mainLoop()
	}
}

func TestRunGamepart2(t *testing.T) {
	for part := 1; part < 9; part++ {
		fmt.Println("### RUN PART", part)
		run(part, STEPS_TO_RUN)
	}
}

