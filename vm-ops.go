package main

import "github.com/neophob/ganother-world/logger"

//Implementation of all VM ops

//Continues the code execution at the indicated address.
func (state *VMState) opJmp() {
	offset := state.fetchWord()
	state.pc = uint16(offset)
	logger.Debug("#op_jmp() jump to %d", state.pc)
}

//Set.i variable, value - Initialises the variable with an integer value from -32768 to 32767.
func (state *VMState) opMovConst() {
	index := state.fetchByte()
	value := int16(state.fetchWord())
	logger.Debug("#op_movConst %d %d", index, value)
	state.variables[index] = value
}

//Initialises variable 1 with variable 2.
func (state *VMState) opMov() {
	dest := state.fetchByte()
	source := state.fetchByte()
	logger.Debug("#op_mov %d %d %d", source, dest, state.variables[source])
	state.variables[dest] = state.variables[source]
}

//Variable = Variable + Integer value
func (state *VMState) opAddConst() {
	if state.gamePart == 5 && state.pc == 0x6D48 {
		logger.Warn("TODO Script::op_addConst() workaround for infinite looping gun sound")
		// The script 0x27 slot 0x17 doesn't stop the gun sound from looping.
		// This is a bug in the original game code, confirmed by Eric Chahi and
		// addressed with the anniversary editions.
		// For older releases (DOS, Amiga), we play the 'stop' sound like it is
		// done in other part of the game code.
		//snd_playSound(0x5B, 1, 64, 1);
	}
	index := state.fetchByte()
	value := int16(state.fetchWord())
	logger.Debug("#op_addConst() index=%d, value=%d, add=%d", index, state.variables[index], value)
	state.variables[index] += value
}

//Add Variable1, Variable2. Variable1 = Variable 1 + Variable2
func (state *VMState) opAdd() {
	dest := state.fetchByte()
	source := state.fetchByte()
	logger.Debug("#op_add() index=%d, var1=%d, var2=%d", dest, state.variables[dest], state.variables[source])
	state.variables[dest] += state.variables[source]
}

//Sub Variable1, Variable2, Variable1 = Variable1 - Variable2
func (state *VMState) opSub() {
	dest := state.fetchByte()
	source := state.fetchByte()
	logger.Debug("#op_sub() %d %d", dest, state.variables[source])
	state.variables[dest] -= state.variables[source]
}

//Variable = Variable AND value
func (state *VMState) opAnd() {
	index := state.fetchByte()
	value := state.fetchWord()
	logger.Debug("#op_and() %d %d", index, value)
	state.variables[index] &= int16(value)
}

//Variable = Variable OR value
func (state *VMState) opOr() {
	index := state.fetchByte()
	value := state.fetchWord()
	logger.Debug("#op_or() %d %d", index, value)
	state.variables[index] |= int16(value)
}

//Makes a N bit rotation to the left on the variable. Zeros on the right.
func (state *VMState) opShl() {
	index := state.fetchByte()
	value := state.fetchWord()
	state.variables[index] <<= uint(value)
	logger.Debug("#op_shl() %d %d %d", index, value, state.variables[index])
}

//Makes a N bit rotation to the right on the variable.
func (state *VMState) opShr() {
	index := state.fetchByte()
	value := state.fetchWord()
	state.variables[index] >>= uint(value)
	logger.Debug("#op_shr() %d %d %d", index, value, state.variables[index])
}

//Jsr Adress - Executes the subroutine located at the indicated address.
func (state *VMState) opCall() {
	newPC := state.fetchWord()
	state.saveSP()
	state.pc = newPC
	logger.Debug("#op_call(), jump to pc: %d", state.pc)
}

//End of a subroutine.
func (state *VMState) opRet() {
	state.restoreSP()
	logger.Debug("#op_ret(), pc: %d", state.pc)
}

//Setvec "numéro de canal", address - Initialises a channel with a code address to execute
//NOTE: if a channel is installed, e.g. Channel 0 installs Channel 1 - this new channel is only respected in the NEXT iteration!
func (state *VMState) opInstallTask() {
	channelID := state.fetchByte()
	address := state.fetchWord()
	logger.Debug("#opInstallTask %d %d", channelID, address)
	state.nextLoopChannelPC[channelID] = address
}

//Break - Temporarily stops the executing channel and goes to the next.
func (state *VMState) opYieldTask() {
	logger.Debug("#opYieldTask")
	state.paused = true
}

//Bigend - Permanently stops the executing channel and goes to the next.
func (state *VMState) opRemoveTask() {
	logger.Debug("#opRemoveTask %d", state.channelID)
	state.pc = VM_INACTIVE_THREAD
	state.paused = true
}

//Vec début, fin, type - Deletes, freezes or unfreezes a series of channels
func (state *VMState) opChangeTaskState() {
	channelIDStart := state.fetchByte()
	channelIDEnd := state.fetchByte()
	changeType := state.fetchByte()
	logger.Debug("#opChangeTaskState %d %d %d", channelIDStart, channelIDEnd, changeType)
	for i := channelIDStart; i <= channelIDEnd; i++ {
		switch changeType {
		case 0:
			state.channelPaused[i] = false
		case 1:
			state.channelPaused[i] = true
		case 2:
			state.nextLoopChannelPC[i] = VM_INACTIVE_THREAD
		}
	}
}

//Dbra Variable, Adress - Decrements the variable, if the result is different from zero the execution continues at the indicated address.
func (state *VMState) opJmpIfVar() {
	index := state.fetchByte()
	state.variables[index]--
	logger.Debug("#opJmpIfVar %d", state.variables[index])
	if state.variables[index] != 0 {
		state.opJmp()
	} else {
		state.fetchWord()
	}
}

//Conditional branch, If (=Si) the comparison of the variables is right, the execution continues at the indicated address.
func (state *VMState) opCondJmp() {
	op := state.fetchByte()
	variableID := uint16(state.fetchByte())
	currentVariable := state.variables[variableID]
	var newVariable int16
	if op&0x80 > 0 {
		newVariable = int16(state.variables[state.fetchByte()])
	} else if op&0x40 > 0 {
		newVariable = int16(state.fetchWord())
	} else {
		newVariable = int16(state.fetchByte())
	}
	logger.Debug("> step #op_condJmp (%d, 0x%02X, 0x%02X) var=0x%02X", op, currentVariable, newVariable, variableID)
	expr := false
	switch op & 7 {
	case 0:
		expr = (currentVariable == newVariable)
		if variableID == 0x29 && op&0x80 != 0 {
			logger.Warn("TODO BYPASS PROTECTION!")
			/*				// 4 symbols
							_scriptVars[0x29] = _scriptVars[0x1E];
							_scriptVars[0x2A] = _scriptVars[0x1F];
							_scriptVars[0x2B] = _scriptVars[0x20];
							_scriptVars[0x2C] = _scriptVars[0x21];
							// counters
							_scriptVars[0x32] = 6;
							_scriptVars[0x64] = 20;
							warning("Script::op_condJmp() bypassing protection");
							expr = true;*/
		}
	case 1:
		expr = (currentVariable != newVariable)
	case 2:
		expr = (currentVariable > newVariable)
	case 3:
		expr = (currentVariable >= newVariable)
	case 4:
		expr = (currentVariable < newVariable)
	case 5:
		expr = (currentVariable <= newVariable)
	default:
		logger.Debug("#op_condJmp: Invalid condition!")
	}
	if expr {
		logger.Debug("> step: TRUE!ILLJUMP")
		state.opJmp()
		//fixUpPalette_changeScreen(_res->_currentPart, _scriptVars[VAR_SCREEN_NUM]);
	} else {
		state.fetchWord()
	}
}

// Fade "palette number" - Changes of colour palette
func (state *VMState) opVidSetPalette(video *Video) {
	index := state.fetchWord()
	video.setPalette(int(index))
}

//Text "text number", x, y, color - Displays in the work screen the specified text for the coordinates x,y.
func (state *VMState) opVidDrawString(video *Video) {
	stringID := int(state.fetchWord())
	x := int(state.fetchByte())
	y := int(state.fetchByte())
	col := int(state.fetchByte())
	video.drawString(col, x, y, stringID)
}

//SetWS "Screen number" - Sets the work screen, which is where the polygons will be drawn by default.
func (state *VMState) opVidSelectPage(video *Video) {
	page := int(state.fetchByte())
	video.setWorkPagePtr(page)
}

//Clr "Screen number", Color - Deletes a screen with one colour. Ingame, there are 4 screen buffers
func (state *VMState) opVidFillPage(video *Video) {
	page := int(state.fetchByte())
	color := int(state.fetchByte())
	video.fillPage(page, color)
}

//Copy "Screen number A", "Screen number B" - Copies screen buffer A to screen buffer B.
func (state *VMState) opVidCopyPage(video *Video) {
	source := state.fetchByte()
	dest := state.fetchByte()

	not0x40 := ^0x40
	isVscrollEnabled := (source & uint8(not0x40)) & 0x80
	if source >= 0xFE || isVscrollEnabled == 0 {
		// no vscroll
		video.copyPage(int(source), int(dest), 0)
		//app.video.copyPage(int(source&uint8(not0x40)), int(dest), 0)
	} else {
		sourceTranslated := int(source & 3)
		vscroll := int(state.variables[VM_VARIABLE_SCROLL_Y])
		if vscroll >= -199 && vscroll <= 199 {
			video.copyPage(sourceTranslated, int(dest), vscroll)
		} else {
			logger.Warn("Invalid VSCROLL %d", vscroll)
		}
	}
}

//Show "Screen number" - Displays the screen buffer specified in the next video frame.
func (state *VMState) opVidUpdatePage(video *Video) {
	page := int(state.fetchByte())
	//TODO inp_handleSpecialKeys();
	if state.gamePart == 0 && state.variables[0x67] == 1 {
		logger.Debug("opVidUpdatePage: BYPASS PROTECTION %d", page)
		state.variables[0xDC] = 33
	}

	video.updateDisplay(page)
}

func (state *VMState) opVidDrawPolyBackground(opcode uint8, video *Video) {
	offset := ((uint16(opcode) << 8) | uint16(state.fetchByte())) << 1
	posX := int16(state.fetchByte())
	posY := int16(state.fetchByte())
	height := posY - 199
	if height > 0 {
		posY = 199
		posX += height
	}
	logger.Debug("opVidDrawPolyBackground %d %d", opcode, offset)
	videoDataFetcher := video.buildReader(false, int(offset))
	video.drawShape(videoDataFetcher, 0xFF, 0x40, int(posX), int(posY))
}

//Spr "'object name" , x, y, z - In the work screen, draws the graphics tool at the coordinates x,y and the zoom factor z. A polygon, a group of polygons...
func (state *VMState) opVidDrawPolySprite(opcode uint8, video *Video) {
	useSecondVideoResource := false
	offsetHi := state.fetchByte()
	offset := ((uint16(offsetHi) << 8) | uint16(state.fetchByte())) << 1
	posX := int16(state.fetchByte())
	if opcode&0x20 == 0 {
		if opcode&0x10 == 0 {
			posX = (posX << 8) | int16(state.fetchByte())
		} else {
			posX = state.variables[posX]
		}
	} else {
		if opcode&0x10 > 0 {
			posX += 0x100
		}
	}
	posY := int16(state.fetchByte())
	if opcode&8 == 0 {
		if opcode&4 == 0 {
			posY = (posY << 8) | int16(state.fetchByte())
		} else {
			posY = state.variables[posY]
		}
	}
	zoom := uint16(state.fetchByte())
	if opcode&2 == 0 {
		if opcode&1 == 0 {
			state.pc--
			logger.Debug("zoom decreased PC %d", state.pc)
			zoom = 0x40
		} else {
			zoom = uint16(state.variables[zoom])
		}
	} else {
		if opcode&1 > 0 {
			useSecondVideoResource = true
			state.pc--
			logger.Debug("useSecondVideoResource! zoom decreased PC: %d", state.pc)
			zoom = 0x40
		}
	}
	logger.Debug("opVidDrawPolySprite %d", offset)
	videoDataFetcher := video.buildReader(useSecondVideoResource, int(offset))
	video.drawShape(videoDataFetcher, 0xFF, int(zoom), int(posX), int(posY))
}

//Initialises a song.
func (state *VMState) opPlayMusic() {
	resNum := int(state.fetchWord())
	delay := int(state.fetchWord())
	pos := int(state.fetchByte())
	logger.Debug("op_playMusic(0x%X, %d, %d)", resNum, delay, pos)
	//TODO snd_playMusic(resNum, delay, pos);
}

//Plays the sound file on one of the four game audio channels with specific height and volume.
func (state *VMState) opPlaySound() {
	resNum := int(state.fetchWord())
	freq := int(state.fetchByte())
	vol := int(state.fetchByte())
	channel := int(state.fetchByte())
	logger.Debug("op_playSound(0x%X, %d, %d, %d)", resNum, freq, vol, channel)
	//TODO snd_playSound(resNum, freq, vol, channel);
}

func (state *VMState) opUpdateResource() {
	id := int(state.fetchWord())
	logger.Debug("opUpdateResource %d", id)
	if id >= GAME_PART_ID_1 {
		logger.Debug("should load next part %d", id)
		state.loadNextPart = id
		return
	}
	if id == 0 {
		logger.Warn("opUpdateResource TODO! INVALIDATE DATA %d", id)
		//_ply->stop();
		//_mix->stopAll();
		//_res->invalidateRes();
		return
	}
	state.assets.loadResource(id)
}
