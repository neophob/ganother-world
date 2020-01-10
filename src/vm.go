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
	sp          int
	pc          int
}

func createNewState() VMState {
	state := VMState{gamePart: -1}
	//WTF? whats this? -> create const
	state.variables[0x54] = 0x81
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
	return state
}

//TODO using a pointer - its own state can be modified. using the value, state does not change!
func (state *VMState) setupGamePart(newGamePart int) {
	if state.gamePart == newGamePart {
		return
	}
	if newGamePart < GAME_PART_FIRST || newGamePart > GAME_PART_LAST {
		panic("INVALID_GAME_PART")
	}

	state.gamePart = newGamePart
	//WTF? whats this? -> create const
	state.variables[0xE4] = 0x14

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.channelData {
		state.channelData[i] = VM_INACTIVE_THREAD
	}

	//TODO WHY?
	state.channelData[0] = 0
}

// Run the Virtual Machine for every active threads
func mainLoop(state VMState) {
	for channelId := 0x00; channelId < VM_NUM_THREADS; channelId++ {
		channelPointerState := state.channelData[channelId]

		// Inactive threads are marked with a thread instruction pointer set to 0xFFFF (VM_INACTIVE_THREAD).
		if channelPointerState != VM_INACTIVE_THREAD {
			fmt.Println("channel active!", channelId, channelPointerState)
			//TODO load resource!
			state.pc = 0 + channelPointerState
			//			_scriptPtr.pc = res->segBytecode + n;
			//		uint8_t opcode = _scriptPtr.fetchByte();
			//  execute
			//state.channelData[channelId] = _scriptPtr.pc - res->segBytecode;
		}
	}
}

func executeOp(code []byte) {
	pc := 0
	step(pc, code)
}

func step(pc int, code []byte) {
	opcode := code[pc]
	switch opcode {
	case 0x00:
		fmt.Println("op_movConst")
		//		uint8 variableId = _scriptPtr.fetchByte();
		//		int16 value = _scriptPtr.fetchWord();

	case 0x01:
		fmt.Println("op_mov")
	case 0x02:
		fmt.Println("op_add")
	case 0x03:
		fmt.Println("op_addConst")

	case 0x04:
		fmt.Println("op_call")
	case 0x05:
		fmt.Println("op_ret")
	case 0x06:
		fmt.Println("op_pauseThread")
	case 0x07:
		fmt.Println("op_jmp")

	case 0x08:
		fmt.Println("op_setSetVect")
	case 0x09:
		fmt.Println("op_jnz")
	case 0x0A:
		fmt.Println("op_condJmp")
	case 0x0B:
		fmt.Println("op_setPalette")

	case 0x0C:
		fmt.Println("op_resetThread")
	case 0x0D:
		fmt.Println("op_selectVideoPage")
	case 0x0E:
		fmt.Println("op_fillVideoPage")
	case 0x0F:
		fmt.Println("op_copyVideoPage")

	case 0x10:
		fmt.Println("op_blitFramebuffer")
	case 0x11:
		fmt.Println("op_killThread")
	case 0x12:
		fmt.Println("op_drawString")
	case 0x13:
		fmt.Println("op_sub")

	case 0x14:
		fmt.Println("op_and")
	case 0x15:
		fmt.Println("op_or")
	case 0x16:
		fmt.Println("op_shl")
	case 0x17:
		fmt.Println("op_shr")

	case 0x18:
		fmt.Println("op_playSound")
	case 0x19:
		fmt.Println("op_updateMemList")
	case 0x1A:
		fmt.Println("op_playMusic")

	}
}
