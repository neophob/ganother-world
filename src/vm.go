package main

import (
	"fmt"
)

const (
	//TODO rename me to channel
	VM_NUM_THREADS   int = 64
	VM_NUM_VARIABLES int = 256

	VM_MAX_STACK_SIZE int = 64

	VM_NO_SETVEC_REQUESTED uint16 = 0xFFFF
	VM_INACTIVE_THREAD     uint16 = 0xFFFF

	VM_VARIABLE_RANDOM_SEED          int = 0x3C
	VM_VARIABLE_SCREEN_NUM           int = 0x67
	VM_VARIABLE_TITLESCREEN          int = 0x54
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
	assets    Assets
	variables [VM_NUM_VARIABLES]int16
	channelPC [VM_NUM_THREADS]uint16
	//TODO not sure if this is used!
	channelPaused [VM_NUM_THREADS]bool
	gamePart      int
	stackCalls    [VM_MAX_STACK_SIZE]uint16
	bytecode      []uint8

	countNoOps int
	//TODO rename channel specific data, sp -> spIndex
	sp        int
	pc        uint16
	channelId int
	paused    bool
}

func createNewState(assets Assets) VMState {
	state := VMState{gamePart: -1, assets: assets}
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
	//WTF? whats this? -> create const
	state.variables[0xE4] = 0x14
	return state
}

func (state *VMState) saveSP() {
	if state.sp >= VM_MAX_STACK_SIZE {
		panic("SaveSP, stack overflow")
	}
	state.stackCalls[state.sp] = state.pc
	state.sp++
}

func (state *VMState) restoreSP() {
	if state.sp == 0 {
		panic("restoreSP, stack underflow")
	}
	state.sp--
	state.pc = state.stackCalls[state.sp]
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
	if newGamePart == GAME_PART_FIRST {
		// VAR(0x54) indicates if the "Out of this World" title screen should be presented
		state.variables[VM_VARIABLE_TITLESCREEN] = 0x81
	}

	gamePartAsset := state.assets.gameParts[newGamePart-GAME_PART_FIRST]
	fmt.Println("- load bytecode, resource", gamePartAsset.bytecode)
	state.bytecode = state.assets.loadEntryFromBank(gamePartAsset.bytecode)
	state.gamePart = newGamePart

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.channelPC {
		state.channelPC[i] = VM_INACTIVE_THREAD
	}

	//activate first channel, set initial PC to 0
	state.channelPC[0] = 0
}

func (state *VMState) executeOp() {
	lastPc := state.pc
	opcode := state.fetchByte()
	fmt.Printf("> step: opcode[%2d], pc[%5d], channel[%2d] >>> ", opcode, lastPc, state.channelId)

	if opcode > 0x7F {
		state.opVidDrawPolyBackground(opcode)
		return
	}
	if opcode > 0x3F {
		state.opVidDrawPolySprite(opcode)
		return
	}

	switch opcode {

	case 0x00:
		state.opMovConst()
	case 0x01:
		state.opMov()
	case 0x02:
		state.opAdd()
	case 0x03:
		state.opAddConst()

	case 0x04:
		state.opCall()
	case 0x05:
		state.opRet()
	case 0x06:
		state.opYieldTask()
	case 0x07:
		state.opJmp()

	case 0x08:
		state.opInstallTask()
	case 0x09:
		state.opJmpIfVar()
	case 0x0A:
		state.opCondJmp()
	case 0x0B:
		state.opVidSetPalette()

	case 0x0C:
		state.opChangeTaskState()
	case 0x0D:
		state.opVidSelectPage()
	case 0x0E:
		state.opVidFillPage()
	case 0x0F:
		state.opVidCopyPage()

	case 0x10:
		state.opVidUpdatePage()
	case 0x11:
		state.opRemoveTask()
	case 0x12:
		state.opVidDrawString()
	case 0x13:
		state.opSub()

	case 0x14:
		state.opAnd()
	case 0x15:
		state.opOr()
	case 0x16:
		state.opShl()
	case 0x17:
		state.opShr()

	case 0x18:
		state.opPlaySound()
	case 0x19:
		state.opUpdateResource()
	case 0x1A:
		state.opPlayMusic()
	default:
		state.countNoOps++
		fmt.Println("NO_OP", opcode)
	}
}

// Run the Virtual Machine for every active threads
func (state *VMState) mainLoop() {
	for channelId := 0x00; channelId < VM_NUM_THREADS; channelId++ {
		channelPC := state.channelPC[channelId]
		channelPaused := state.channelPaused[channelId]
		if channelPaused {
			fmt.Println("IIIKKKSSS!!!!!!!!!!", channelId, channelPC)
		}
		// Inactive threads are marked with a thread instruction pointer set to 0xFFFF (VM_INACTIVE_THREAD).
		if channelPC != VM_INACTIVE_THREAD && !channelPaused {
			state.channelId = channelId
			state.paused = false
			state.pc = channelPC
			state.sp = 0
			for state.paused == false {
				state.executeOp()
			}
			if state.sp > 0 {
				fmt.Println("WARNING, SP > 0", state.sp)
			}
			state.channelPC[channelId] = state.pc
		}
	}
}
