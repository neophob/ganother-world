package main

const (
	//TODO rename me to channel
	VM_NUM_THREADS   int = 64
	VM_NUM_VARIABLES int = 256

	VM_MAX_STACK_SIZE int = 64

	VM_NO_SETVEC_REQUESTED uint16 = 0xFFFF
	VM_INACTIVE_THREAD     uint16 = 0xFFFF
	VM_NO_TASK_OP          uint16 = 0xFFFE

	VM_VARIABLE_RANDOM_SEED          int = 0x3C
	VM_VARIABLE_SCREEN_NUM           int = 0x67
	VM_VARIABLE_TITLESCREEN          int = 0x54
	VM_VARIABLE_LAST_KEYCHAR         int = 0xDA
	VM_VARIABLE_HERO_POS_UP_DOWN     int = 0xE5
	VM_VARIABLE_MUS_MARK             int = 0xF4
	VM_VARIABLE_SCROLL_Y             int = 0xF9
	VM_VARIABLE_PROTECTION_CHECK     int = 0xF2
	VM_VARIABLE_HERO_ACTION          int = 0xFA
	VM_VARIABLE_HERO_POS_JUMP_DOWN   int = 0xFB
	VM_VARIABLE_HERO_POS_LEFT_RIGHT  int = 0xFC
	VM_VARIABLE_HERO_POS_MASK        int = 0xFD
	VM_VARIABLE_HERO_ACTION_POS_MASK int = 0xFE
	VM_VARIABLE_PAUSE_SLICES         int = 0xFF
)

type VMState struct {
	assets            Assets
	variables         [VM_NUM_VARIABLES]int16
	channelPC         [VM_NUM_THREADS]uint16
	nextLoopChannelPC [VM_NUM_THREADS]uint16
	channelPaused     [VM_NUM_THREADS]bool
	stackCalls        [VM_MAX_STACK_SIZE]uint16
	gamePart          int
	loadNextPart      int

	palette   []uint8
	bytecode  []uint8
	cinematic []uint8
	video2    []uint8

	//context of the current channel
	sp        int
	pc        uint16
	channelID int
	paused    bool

	//statistics
	countNoOps     int
	countOps       int
	countSPNotZero int
}

func createNewState(assets Assets) VMState {
	state := VMState{gamePart: -1, assets: assets}
	state.variables[VM_VARIABLE_RANDOM_SEED] = 42
	//WTF? whats this? -> create const
	state.variables[0xE4] = 0x14

	//BYPASS PROTECTION START
	state.variables[0xBC] = 0x10
	state.variables[0xC6] = 0x80
	state.variables[0xDC] = 0x21
	// this variable explicit checked when level 2 starts - if content != 0xFA0 then a cond jump will be made to nowhere
	state.variables[VM_VARIABLE_PROTECTION_CHECK] = 0xFA0
	//BYPASS END

	return state
}

//TODO integrate with setupGamePart
func (state VMState) buildVideoAssets() VideoAssets {
	return VideoAssets{
		palette:   state.palette,
		cinematic: state.cinematic,
		video2:    state.video2,
	}
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

func (state *VMState) setupGamePart(newGamePart int) {
	if newGamePart < GAME_PART_FIRST || newGamePart > GAME_PART_LAST {
		panic("INVALID_GAME_PART")
	}
	if newGamePart == GAME_PART_FIRST {
		// VAR(0x54) indicates if the "Out of this World" title screen should be presented
		// language setting?
		state.variables[VM_VARIABLE_TITLESCREEN] = 0x81
	}

	state.loadGameParts(newGamePart - GAME_PART_FIRST)

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.channelPC {
		state.channelPC[i] = VM_INACTIVE_THREAD
		state.nextLoopChannelPC[i] = VM_NO_TASK_OP
		state.channelPaused[i] = false
	}

	//activate first channel, set initial PC to 0
	state.pc = 0
	state.channelPC[0] = state.pc

	state.loadNextPart = 0
}

// gamePart is the int between 0 and 10
func (state *VMState) loadGameParts(gamePart int) {
	Debug("LOAD GAME PART %d", gamePart)
	state.gamePart = gamePart

	gamePartAsset := state.assets.gameParts[gamePart]
	state.bytecode = state.assets.loadEntryFromBank(gamePartAsset.bytecode)
	state.palette = state.assets.loadEntryFromBank(gamePartAsset.palette)
	state.cinematic = state.assets.loadEntryFromBank(gamePartAsset.cinematic)
	state.video2 = state.assets.loadEntryFromBank(gamePartAsset.video2)
}

// Run the Virtual Machine for every active threads
func (state *VMState) mainLoop(keypresses uint32) {
	state.handleKeypress(keypresses)

	//TODO check if next part needs to be loaded!
	state.setupChannels()
	for channelID := 0; channelID < VM_NUM_THREADS; channelID++ {
		channelPC := state.channelPC[channelID]
		channelPaused := state.channelPaused[channelID]
		// Inactive threads are marked with a thread instruction pointer set to 0xFFFF (VM_INACTIVE_THREAD).
		if channelPC != VM_INACTIVE_THREAD && !channelPaused {
			state.channelID = channelID
			state.paused = false
			state.pc = channelPC
			state.sp = 0
			//loop channel until finished
			for state.paused == false {
				state.executeOp()
			}
			Debug("> step: PAUSED, pc[%5d], channel[%2d] >>> ", state.pc-1, channelID)
			if state.sp > 0 {
				state.countSPNotZero++
			}
			state.channelPC[channelID] = state.pc
		}
	}
	Debug("> --- MAINLOOP DONE")
}

func (state *VMState) handleKeypress(keypresses uint32) {
	leftRight := int16(0)
	upDown := int16(0)
	mask := int16(0)
	if keypresses&KEY_UP > 0 {
		upDown = -1
		mask |= 1
	}
	if keypresses&KEY_DOWN > 0 {
		upDown = 1
		mask |= 2
	}
	if keypresses&KEY_LEFT > 0 {
		leftRight = -1
		mask |= 4
	}
	if keypresses&KEY_RIGHT > 0 {
		leftRight = 1
		mask |= 8
	}

	state.variables[VM_VARIABLE_HERO_POS_UP_DOWN] = upDown
	state.variables[VM_VARIABLE_HERO_POS_JUMP_DOWN] = upDown
	state.variables[VM_VARIABLE_HERO_POS_LEFT_RIGHT] = leftRight
	state.variables[VM_VARIABLE_HERO_POS_MASK] = mask

	fireButton := int16(0)
	if keypresses&KEY_FIRE > 0 {
		fireButton = 1
		mask |= 80
	}

	state.variables[VM_VARIABLE_HERO_ACTION] = fireButton
	state.variables[VM_VARIABLE_HERO_POS_MASK] = mask
}

//no pending tasks when starting a loop
func (state *VMState) setupChannels() {
	for channelID := 0; channelID < VM_NUM_THREADS; channelID++ {
		if state.nextLoopChannelPC[channelID] != VM_NO_TASK_OP {
			state.channelPC[channelID] = state.nextLoopChannelPC[channelID]
			state.nextLoopChannelPC[channelID] = VM_NO_TASK_OP
		}
	}
}

func (state *VMState) executeOp() {
	opcode := state.fetchByte()
	Debug("> step: opcode[%2d], pc[%5d], channel[%2d] >>> ", opcode, state.pc-1, state.channelID)

	state.countOps++

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
		Warn("NO_OP: %d", opcode)
	}
}
