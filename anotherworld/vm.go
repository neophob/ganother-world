package anotherworld

import (
	"github.com/neophob/ganother-world/logger"
)

const (
	//TODO rename me to channel
	VM_NUM_THREADS    int = 64
	VM_NUM_VARIABLES  int = 256
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

//VMState implements the state of the the VM
type VMState struct {
	assets            StaticGameAssets
	Variables         [VM_NUM_VARIABLES]int16
	ChannelPC         [VM_NUM_THREADS]uint16
	NextLoopChannelPC [VM_NUM_THREADS]uint16
	ChannelPaused     [VM_NUM_THREADS]bool
	StackCalls        [VM_MAX_STACK_SIZE]uint16
	GamePart          int
	LoadNextPart      int

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
	CountNoOps     int
	countOps       int
	CountSPNotZero int
}

func CreateNewState(assets StaticGameAssets) VMState {
	state := VMState{GamePart: -1, assets: assets}
	state.Variables[VM_VARIABLE_RANDOM_SEED] = 42
	//WTF? whats this? -> create const
	state.Variables[0xE4] = 0x14

	//BYPASS PROTECTION START
	state.Variables[0xBC] = 0x10
	state.Variables[0xC6] = 0x80
	state.Variables[0xDC] = 0x21
	// this variable explicit checked when level 2 starts - if content != 0xFA0 then a cond jump will be made to nowhere
	state.Variables[VM_VARIABLE_PROTECTION_CHECK] = 0xFA0
	//BYPASS END

	return state
}

//TODO integrate with SetupGamePart
func (state VMState) BuildVideoAssets() VideoAssets {
	return VideoAssets{
		Palette:   state.palette,
		Cinematic: state.cinematic,
		Video2:    state.video2,
	}
}

func (state *VMState) saveSP() {
	if state.sp >= VM_MAX_STACK_SIZE {
		panic("SaveSP, stack overflow")
	}
	state.StackCalls[state.sp] = state.pc
	state.sp++
}

func (state *VMState) restoreSP() {
	if state.sp == 0 {
		panic("restoreSP, stack underflow")
	}
	state.sp--
	state.pc = state.StackCalls[state.sp]
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
	return ToUint16BE(b1, b2)
}

func (state *VMState) SetupGamePart(newGamePart int) {
	if newGamePart < GAME_PART_FIRST || newGamePart > GAME_PART_LAST {
		panic("INVALID_GAME_PART")
	}
	if newGamePart == GAME_PART_FIRST {
		// VAR(0x54) indicates if the "Out of this World" title screen should be presented
		// language setting?
		state.Variables[VM_VARIABLE_TITLESCREEN] = 0x81
	}

	state.LoadGameParts(newGamePart - GAME_PART_FIRST)

	//Set all thread to inactive (pc at 0xFFFF or 0xFFFE )
	for i := range state.ChannelPC {
		state.ChannelPC[i] = VM_INACTIVE_THREAD
		state.NextLoopChannelPC[i] = VM_NO_TASK_OP
		state.ChannelPaused[i] = false
	}

	//activate first channel, set initial PC to 0
	state.pc = 0
	state.ChannelPC[0] = state.pc

	state.LoadNextPart = 0
}

// gamePart is the int between 0 and 10
func (state *VMState) LoadGameParts(gamePart int) {
	logger.Debug("LOAD GAME PART %d", gamePart)
	state.GamePart = gamePart

	gamePartAsset := state.assets.GameParts[gamePart]
	state.bytecode = state.assets.LoadEntryFromBank(gamePartAsset.Bytecode)
	state.palette = state.assets.LoadEntryFromBank(gamePartAsset.Palette)
	state.cinematic = state.assets.LoadEntryFromBank(gamePartAsset.Cinematic)
	state.video2 = state.assets.LoadEntryFromBank(gamePartAsset.Video2)
}

// Run the Virtual Machine for every active threads
func (state *VMState) MainLoop(keypresses uint32, video *Video) {
	state.handleKeypress(keypresses)

	//TODO check if next part needs to be loaded!
	state.setupChannels()
	for channelID := 0; channelID < VM_NUM_THREADS; channelID++ {
		channelPC := state.ChannelPC[channelID]
		channelPaused := state.ChannelPaused[channelID]
		// Inactive threads are marked with a thread instruction pointer set to 0xFFFF (VM_INACTIVE_THREAD).
		if channelPC != VM_INACTIVE_THREAD && !channelPaused {
			state.channelID = channelID
			state.paused = false
			state.pc = channelPC
			state.sp = 0
			//loop channel until finished
			for state.paused == false {
				state.executeOp(video)
			}
			logger.Debug("> step: PAUSED, pc[%5d], channel[%2d] >>> ", state.pc-1, channelID)
			if state.sp > 0 {
				state.CountSPNotZero++
			}
			state.ChannelPC[channelID] = state.pc
		}
	}
	logger.Debug("> --- MAINLOOP DONE")
}

func (state *VMState) handleKeypress(keypresses uint32) {
	leftRight := int16(0)
	upDown := int16(0)
	mask := int16(0)
	if keypresses&KeyRight > 0 {
		leftRight = 1
		mask |= 1
	}
	if keypresses&KeyLeft > 0 {
		leftRight = -1
		mask |= 2
	}
	if keypresses&KeyDown > 0 {
		upDown = 1
		mask |= 4
	}
	if keypresses&KeyUp > 0 {
		upDown = -1
		mask |= 8
	}
	state.Variables[VM_VARIABLE_HERO_POS_UP_DOWN] = upDown
	state.Variables[VM_VARIABLE_HERO_POS_JUMP_DOWN] = upDown
	state.Variables[VM_VARIABLE_HERO_POS_LEFT_RIGHT] = leftRight
	state.Variables[VM_VARIABLE_HERO_POS_MASK] = mask

	fireButton := int16(0)
	if keypresses&KeyFire > 0 {
		fireButton = 1
		mask |= 0x80
	}
	state.Variables[VM_VARIABLE_HERO_ACTION] = fireButton
	state.Variables[VM_VARIABLE_HERO_ACTION_POS_MASK] = mask
}

//no pending tasks when starting a loop
func (state *VMState) setupChannels() {
	for channelID := 0; channelID < VM_NUM_THREADS; channelID++ {
		if state.NextLoopChannelPC[channelID] != VM_NO_TASK_OP {
			state.ChannelPC[channelID] = state.NextLoopChannelPC[channelID]
			state.NextLoopChannelPC[channelID] = VM_NO_TASK_OP
		}
	}
}

func (state *VMState) executeOp(video *Video) {
	opcode := state.fetchByte()
	logger.Debug("> step: opcode[%2d], pc[%5d], channel[%2d] >>> ", opcode, state.pc-1, state.channelID)

	state.countOps++

	if opcode > 0x7F {
		state.opVidDrawPolyBackground(opcode, video)
		return
	}
	if opcode > 0x3F {
		state.opVidDrawPolySprite(opcode, video)
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
		state.opVidSetPalette(video)

	case 0x0C:
		state.opChangeTaskState()
	case 0x0D:
		state.opVidSelectPage(video)
	case 0x0E:
		state.opVidFillPage(video)
	case 0x0F:
		state.opVidCopyPage(video)

	case 0x10:
		state.opVidUpdatePage(video)
	case 0x11:
		state.opRemoveTask()
	case 0x12:
		state.opVidDrawString(video)
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
		state.opPlaySound(video)
	case 0x19:
		state.opUpdateResource()
	case 0x1A:
		state.opPlayMusic(video)
	default:
		state.CountNoOps++
		logger.Warn("NO_OP: %d", opcode)
	}
}
