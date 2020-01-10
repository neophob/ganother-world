package main

import (
	"fmt"
)

const (
	//TODO rename me to channel
	VM_NUM_THREADS   int = 64
	VM_NUM_VARIABLES int = 256

	GAME_PART_FIRST int = 0x3E80
	GAME_PART_LAST  int = 0x3E89
	GAME_PART1      int = 0x3E80
	GAME_PART2      int = 0x3E81 //Introductino
	GAME_PART3      int = 0x3E82
	GAME_PART4      int = 0x3E83 //Wake up in the suspended jail
	GAME_PART5      int = 0x3E84
	GAME_PART6      int = 0x3E85 //BattleChar sequence
	GAME_PART7      int = 0x3E86
	GAME_PART8      int = 0x3E87
	GAME_PART9      int = 0x3E88
	GAME_PART10     int = 0x3E89

	VM_NO_SETVEC_REQUESTED int = 0xFFFF
	VM_INACTIVE_THREAD     int = 0xFFFF

	VM_VARIABLE_RANDOM_SEED          int = 0x3C
	VM_VARIABLE_LAST_KEYCHAR         int = 0xDA
	VM_VARIABLE_HERO_POS_UP_DOWN     int = 0xE5
	VM_VARIABLE_MUS_MARK             int = 0xF4
	VM_VARIABLE_SCROLL_Y             int = 0xF9
	VM_VARIABLE_HERO_ACTION          int = 0xFA
	VM_VARIABLE_HERO_POS_JUMP_DOWN   int = 0xFB
	VM_VARIABLE_HERO_POS_LEFT_RIGHT  int = 0xFC
	VM_VARIABLE_HERO_POS_MASK        int = 0xFD
	VM_VARIABLE_HERO_ACTION_POS_MASK int = 0xFE
	VM_VARIABLE_PAUSE_SLICES         int = 0xFF
)

type VMState struct {
	variables   [VM_NUM_VARIABLES]int
	channelData [VM_NUM_THREADS]int
	gamePart    int
}

func createNewState() VMState {
	state := VMState{gamePart: -1}
	//WTF? whats this? -> create const
	state.variables[0x54] = 0x81
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
	return state
}

func setupGamePart(state VMState, newGamePart int) VMState {
	if state.gamePart == newGamePart {
		return state
	}
	if newGamePart < GAME_PART_FIRST || newGamePart > GAME_PART_LAST {
		panic("INVALID_GAME_PART")
	}

	state.gamePart = newGamePart
	//WTF? whats this? -> create const
	state.variables[0xE4] = 0x14

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.channelData {
		state.channelData[i] = 0xFF
	}

	//TODO WHY?
	state.channelData[0] = 0
	return state
}

// Run the Virtual Machine for every active threads
func mainLoop(state VMState) {
	for channelId := 0x00; channelId < VM_NUM_THREADS; channelId++ {
		n := state.channelData[channelId]
		fmt.Println("channel", channelId, n)
		if n != VM_INACTIVE_THREAD {

		}
	}

}
