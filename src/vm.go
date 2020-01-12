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
	assets      Assets
	variables   [VM_NUM_VARIABLES]int
	channelData [VM_NUM_THREADS]int
	gamePart    int
	sp          int
	stackCalls  [64]int
	pc          int
	bytecode    []uint8
}

func createNewState(assets Assets) VMState {
	state := VMState{gamePart: -1, assets: assets}
	//WTF? whats this? -> create const
	state.variables[0x54] = 0x81
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
	return state
}

func (state *VMState) saveCurrentSP() {
	if state.sp >= 0x40 {
		panic("Script::op_call() sp>0x40 stack overflow")
	}

	state.stackCalls[state.sp] = state.pc
	state.sp++
}

func (state *VMState) fetchByte() uint8 {
	result := state.bytecode[state.pc]
	state.pc++
	return result
}

func (state *VMState) fetchWord() uint16 {
	b1 := state.bytecode[state.pc]
	b2 := state.bytecode[state.pc+1]
	state.pc += 2
	return toUint16BE(b1, b2)
}

//TODO using a pointer - its own state can be modified. using the value, state does not change!
func (state *VMState) setupGamePart(newGamePart int) {
	if state.gamePart == newGamePart {
		return
	}
	if newGamePart < GAME_PART_FIRST || newGamePart > GAME_PART_LAST {
		panic("INVALID_GAME_PART")
	}

	//TODO get bytecode from current game part and add it to the VMstate
	gamePartAsset := state.assets.gameParts[newGamePart-GAME_PART_FIRST]
	fmt.Println("- load bytecode, resource", gamePartAsset.bytecode)
	state.bytecode = state.assets.loadEntryFromBank(gamePartAsset.bytecode)
	fmt.Println("- executeOp", state.bytecode[0:32])

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

func (state *VMState) executeOp() {
	opcode := state.bytecode[state.pc]
	fmt.Println("> step", opcode, state.pc)

	if opcode > 0x7F {
		offset := ((opcode << 8) | state.fetchByte()) << 1
		posX := state.fetchByte()
		posY := state.fetchByte()
		height := posY - 199
		if (height > 0) {
			posY = 199;
			posX += height;
		}
		fmt.Println("DRAW_POLY_BACKGROUND", posX, posY, offset)
		return
	}
	if opcode > 0x3F {
		offsetHi := state.fetchByte()
		offset := ((offsetHi << 8) | state.fetchByte()) << 1
		posX := state.fetchByte()
		posY := state.fetchByte()
		zoom := state.fetchByte()

		fmt.Println("DRAW_POLY_SPRITE", posX, posY, offset, zoom)
		return
	}

	//var offset uint16

	switch opcode {

	case 0x00:
		state.opMovConst()
	case 0x01:
		fmt.Println("op_mov")
		//uint8_t dstVariableId = _scriptPtr.fetchByte();
		//uint8_t srcVariableId = _scriptPtr.fetchByte();
	case 0x02:
		fmt.Println("op_add")
		//uint8_t dstVariableId = _scriptPtr.fetchByte();
		//uint8_t srcVariableId = _scriptPtr.fetchByte();
	case 0x03:
		fmt.Println("op_addConst")
		//uint8_t variableId = _scriptPtr.fetchByte();
		//int16_t value = _scriptPtr.fetchWord();

	case 0x04:
		state.opCall()
	case 0x05:
		fmt.Println("op_ret")
	case 0x06:
		fmt.Println("op_pauseThread")
	case 0x07:
		state.opJmp()
	case 0x08:
		fmt.Println("op_setVect")
		//		uint8_t threadId = _scriptPtr.fetchByte();
		//		uint16_t pcOffsetRequested = _scriptPtr.fetchWord();
	case 0x09:
		fmt.Println("op_jnz")
		//		uint8_t i = _scriptPtr.fetchByte();
		//	  _scriptPtr.fetchWord();
	case 0x0A:
		state.opCondJmp()
	case 0x0B:
		fmt.Println("op_setPalette")
		//		uint16_t paletteId = _scriptPtr.fetchWord();

	case 0x0C:
		fmt.Println("op_resetThread")
		//		uint8_t threadId = _scriptPtr.fetchByte();
		//		uint8_t i =        _scriptPtr.fetchByte();
	case 0x0D:
		fmt.Println("op_selectVideoPage")
		//		uint8_t frameBufferId = _scriptPtr.fetchByte();
	case 0x0E:
		fmt.Println("op_fillVideoPage")
		//		uint8_t pageId = _scriptPtr.fetchByte();
		//		uint8_t color = _scriptPtr.fetchByte()
	case 0x0F:
		fmt.Println("op_copyVideoPage")
		//		uint8_t srcPageId = _scriptPtr.fetchByte();
		//		uint8_t dstPageId = _scriptPtr.fetchByte();

	case 0x10:
		page := state.fetchByte()
		fmt.Println("Script::op_blitFramebuffer(), page", page)
	case 0x11:
		fmt.Println("op_killThread")
	case 0x12:
		fmt.Println("op_drawString")
		//		uint16_t stringId = _scriptPtr.fetchWord();
		//		uint16_t x = _scriptPtr.fetchByte();
		//		uint16_t y = _scriptPtr.fetchByte();
		//		uint16_t color = _scriptPtr.fetchByte();
	case 0x13:
		fmt.Println("op_sub")
		//		uint8_t i = _scriptPtr.fetchByte();
		//		uint8_t j = _scriptPtr.fetchByte();

	case 0x14:
		fmt.Println("op_and")
		//		uint8_t variableId = _scriptPtr.fetchByte();
		//		uint16_t n = _scriptPtr.fetchWord();
	case 0x15:
		fmt.Println("op_or")
		//		uint8_t variableId = _scriptPtr.fetchByte();
		//		uint16_t value = _scriptPtr.fetchWord();
	case 0x16:
		fmt.Println("op_shl")
		//		uint8_t variableId = _scriptPtr.fetchByte();
		//		uint16_t leftShiftValue = _scriptPtr.fetchWord();
	case 0x17:
		fmt.Println("op_shr")
		//		uint8_t variableId = _scriptPtr.fetchByte();
		//		uint16_t rightShiftValue = _scriptPtr.fetchWord();

	case 0x18:
		fmt.Println("op_playSound")
		//		uint16_t resourceId = _scriptPtr.fetchWord();
		//		uint8_t freq = _scriptPtr.fetchByte();
		//		uint8_t vol = _scriptPtr.fetchByte();
		//		uint8_t channel = _scriptPtr.fetchByte();
	case 0x19:
		fmt.Println("op_updateMemList aka load resource")
		//		uint16_t resourceId = _scriptPtr.fetchWord();
	case 0x1A:
		fmt.Println("op_playMusic")
		//		uint16_t resNum = _scriptPtr.fetchWord();
		//		uint16_t delay = _scriptPtr.fetchWord();
		//		uint8_t pos = _scriptPtr.fetchByte();
	default:
		fmt.Println("NO_OP", opcode)
	}
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
