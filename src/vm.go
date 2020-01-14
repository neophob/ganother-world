package main

import (
	"fmt"
)

const (
	//TODO rename me to channel
	VM_NUM_THREADS   int = 64
	VM_NUM_VARIABLES int = 256

	VM_MAX_STACK_SIZE int = 64

	VM_NO_SETVEC_REQUESTED int = 0xFFFF
	VM_INACTIVE_THREAD     int = 0xFFFF

	VM_VARIABLE_RANDOM_SEED          int = 0x3C
	VM_VARIABLE_SCREEN_NUM           int = 0x67
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
	assets     Assets
	variables  [VM_NUM_VARIABLES]int
	channelPC  [VM_NUM_THREADS]int
	gamePart   int
	stackCalls [VM_MAX_STACK_SIZE]int
	bytecode   []uint8

	//TODO rename channel specific data
	sp        int
	pc        int
	channelId int
	paused    bool
}

func createNewState(assets Assets) VMState {
	state := VMState{gamePart: -1, assets: assets}
	//WTF? whats this? -> create const
	state.variables[0x54] = 0x81
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
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

	//TODO get bytecode from current game part and add it to the VMstate
	gamePartAsset := state.assets.gameParts[newGamePart-GAME_PART_FIRST]
	fmt.Println("- load bytecode, resource", gamePartAsset.bytecode)
	state.bytecode = state.assets.loadEntryFromBank(gamePartAsset.bytecode)
	fmt.Println("- executeOp", state.bytecode[0:32])

	state.gamePart = newGamePart
	//WTF? whats this? -> create const
	state.variables[0xE4] = 0x14

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.channelPC {
		state.channelPC[i] = VM_INACTIVE_THREAD
	}

	//activate first channel, set initial PC to 0
	state.channelPC[0] = 0
}

func (state *VMState) executeOp() {
	opcode := state.bytecode[state.pc]
	fmt.Printf("> step: opcode[%2d], pc[%5d], channel[%2d] >>> ", opcode, state.pc, state.channelId)
	state.pc++

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
		fmt.Println("NO_OP", opcode)
	}
}

// Run the Virtual Machine for every active threads
func (state *VMState) mainLoop() {
	for channelId := 0x00; channelId < VM_NUM_THREADS; channelId++ {
		channelPointerState := state.channelPC[channelId]

		// Inactive threads are marked with a thread instruction pointer set to 0xFFFF (VM_INACTIVE_THREAD).
		if channelPointerState != VM_INACTIVE_THREAD {
			state.channelId = channelId
			state.paused = false
			state.pc = channelPointerState
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
