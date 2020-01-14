package main

import (
	"fmt"
)
/*
# ranges - should be a int16 type!
> step: opcode[78], pc[20170], channel[16] >>> >VID: SETDATABUFFER 114708
>VID: DRAWSHAPE color:255, x:17193, y:32, zoom:0
> step: opcode[78], pc[20177], channel[16] >>> >VID: SETDATABUFFER 123440
>VID: DRAWSHAPE color:255, x:95, y:28, zoom:64

> step: opcode[ 4], pc[    0], channel[ 0] >>> #op_call(), jump to pc: 17119 62683
> step: opcode[14], pc[17119], channel[ 0] >>> >VID: FILLPAGE 1 0
> step: opcode[14], pc[17122], channel[ 0] >>> >VID: FILLPAGE 2 0
> step: opcode[13], pc[17125], channel[ 0] >>> >VID: SETWORKPAGEPTR 1
> step: opcode[16], pc[17127], channel[ 0] >>> >VID: UPDATEDISPLAY 1
> step: opcode[16], pc[17129], channel[ 0] >>> >VID: UPDATEDISPLAY 255
> step: opcode[16], pc[17131], channel[ 0] >>> >VID: UPDATEDISPLAY 2
> step: opcode[ 5], pc[17133], channel[ 0] >>> #op_ret(), pc: 1
> step: opcode[66], pc[    1], channel[ 0] >>> >VID: SETDATABUFFER 114192
>VID: DRAWSHAPE color:255, x:15423, y:23552, zoom:255
> step: opcode[ 0], pc[    9], channel[ 0] >>> #op_movConst 4 246
> step: opcode[ 0], pc[   13], channel[ 0] >>> #op_movConst 1 50
> step: opcode[ 0], pc[   17], channel[ 0] >>> #op_movConst 0 1560
> step: opcode[ 0], pc[   21], channel[ 0] >>> #op_movConst 106 5120
> step: opcode[ 2], pc[   25], channel[ 0] >>> #op_add() index=24, var1=0, var2=1560
> step: opcode[106], pc[   28], channel[ 0] >>> >VID: SETDATABUFFER 10240
>VID: DRAWSHAPE color:255, x:3, y:24, zoom:0
> step: opcode[106], pc[   34], channel[ 0] >>> >VID: SETDATABUFFER 10240
>VID: DRAWSHAPE color:255, x:1, y:24, zoom:0
> step: opcode[106], pc[   40], channel[ 0] >>> >VID: SETDATABUFFER 10240
>VID: DRAWSHAPE color:255, x:0, y:4, zoom:71
> step: opcode[56], pc[   46], channel[ 0] >>> NO_OP 56
*/

//Continues the code execution at the indicated address.
func (state *VMState) opJmp() {
	offset := state.fetchWord()
	state.pc = int(offset)
	fmt.Println("#op_jmp() jump to", state.pc)
}

//Set.i variable, value - Initialises the variable with an integer value from -32768 to 32767.
func (state *VMState) opMovConst() {
	index := state.fetchByte()
	value := int16(state.fetchWord())
	fmt.Println("#op_movConst", index, value)
	state.variables[index] = value
}

//Initialises variable 1 with variable 2.
func (state *VMState) opMov() {
	dest := int(state.fetchByte())
	source := int(state.fetchByte())
	fmt.Println("#op_mov", source, dest, state.variables[source])
	state.variables[dest] = state.variables[source]
}

//Variable = Variable + Integer value
func (state *VMState) opAddConst() {
	//TODO add workaround for vm bug
	//		if (_res->_currentPart == 16006 && _scriptPtr.pc == _res->_segCode + 0x6D48) {
	//	warning("Script::op_addConst() workaround for infinite looping gun sound");
	//snd_playSound(0x5B, 1, 64, 1);
	index := state.fetchByte()
	value := int16(state.fetchWord())
	fmt.Printf("#op_addConst() index=%d, value=%d, add=%d\n", index, state.variables[index], value)
	state.variables[index] += value
}

//Add Variable1, Variable2. Variable1 = Variable 1 + Variable2
func (state *VMState) opAdd() {
	dest := state.fetchByte()
	source := state.fetchByte()
	fmt.Printf("#op_add() index=%d, var1=%d, var2=%d\n", dest, state.variables[dest], state.variables[source])
	state.variables[dest] += state.variables[source]
}

//Sub Variable1, Variable2, Variable1 = Variable1 - Variable2
func (state *VMState) opSub() {
	dest := state.fetchByte()
	source := state.fetchByte()
	fmt.Println("#op_sub()", dest, state.variables[source])
	state.variables[dest] -= state.variables[source]
}

//Variable = Variable AND value
func (state *VMState) opAnd() {
	index := state.fetchByte()
	value := int16(state.fetchWord())
	fmt.Println("#op_and()", index, value)
	state.variables[index] &= value
}

//Variable = Variable OR value
func (state *VMState) opOr() {
	index := state.fetchByte()
	value := int16(state.fetchWord())
	fmt.Println("#op_or()", index, value)
	state.variables[index] |= value
}

//Makes a N bit rotation to the left on the variable. Zeros on the right.
func (state *VMState) opShl() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	state.variables[index] <<= uint(value)
	fmt.Println("#op_shl()", index, value, state.variables[index])
}

//Makes a N bit rotation to the right on the variable.
func (state *VMState) opShr() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	state.variables[index] >>= uint(value)
	fmt.Println("#op_shr()", index, value, state.variables[index])
}

//Jsr Adress - Executes the subroutine located at the indicated address.
func (state *VMState) opCall() {
	state.saveSP()
	state.pc = int(state.fetchWord())
	fmt.Println("#op_call(), jump to pc:", state.pc, len(state.bytecode))
}

//End of a subroutine.
func (state *VMState) opRet() {
	state.restoreSP()
	fmt.Println("#op_ret(), pc:", state.pc)
}

//Setvec "numéro de canal", address - Initialises a channel with a code address to execute
func (state *VMState) opInstallTask() {
	channelID := state.fetchByte()
	address := int(state.fetchWord())
	fmt.Println("#opInstallTask", channelID, address)
	// TODO validate me: 	_scriptTasks[1][i] = n;
	state.channelPC[channelID] = address
}

//Break - Temporarily stops the executing channel and goes to the next.
func (state *VMState) opYieldTask() {
	fmt.Println("#opYieldTask")
	state.paused = true
}

//Bigend - Permanently stops the executing channel and goes to the next.
func (state *VMState) opRemoveTask() {
	fmt.Println("#opRemoveTask", state.channelId)
	state.pc = VM_INACTIVE_THREAD
	state.paused = true
}

//Vec début, fin, type - Deletes, freezes or unfreezes a series of channels
func (state *VMState) opChangeTaskState() {
	channelIdStart := int(state.fetchByte())
	channelIdEnd := int(state.fetchByte())
	changeType := state.fetchByte()
	fmt.Println("#opChangeTaskState", channelIdStart, channelIdEnd, changeType)
	for i := channelIdStart; i <= channelIdEnd; i++ {
		if changeType == 2 {
			state.channelPaused[i] = false
			//fmt.Println("#opChangeTaskState UNBLOCK? TODO", i)
		} else {
			fmt.Println("#opChangeTaskState TODO", i, changeType)
		}
	}
}

//Dbra Variable, Adress - Decrements the variable, if the result is different from zero the execution continues at the indicated address.
func (state *VMState) opJmpIfVar() {
	index := state.fetchByte()
	state.variables[index]--
	fmt.Println("#opJmpIfVar", state.variables[index])
	if state.variables[index] != 0 {
		state.opJmp()
	} else {
		state.fetchWord()
	}
}

//Conditional branch, If (=Si) the comparison of the variables is right, the execution continues at the indicated address.
func (state *VMState) opCondJmp() {
	op := state.fetchByte()
	variableId := uint16(state.fetchByte())
	currentVariable := state.variables[variableId]
	var newVariable uint16
	if op&0x80 > 0 {
		newVariable = uint16(state.variables[state.fetchByte()])
	} else if op&0x40 > 0 {
		newVariable = state.fetchWord()
	} else {
		newVariable = uint16(state.fetchByte())
	}
	fmt.Printf("#op_condJmp op=%d, variableId=%d, currentVariable=%d, newVariable=%d\n", op, variableId, currentVariable, newVariable)
	expr := false
	switch op & 7 {
	case 0:
		expr = (variableId == newVariable)
		//TODO implement bypass protection
	case 1:
		expr = (variableId != newVariable)
	case 2:
		expr = (variableId > newVariable)
	case 3:
		expr = (variableId >= newVariable)
	case 4:
		expr = (variableId < newVariable)
	case 5:
		expr = (variableId <= newVariable)
	default:
		fmt.Println("#op_condJmp: Invalid condition!")
	}
	if expr {
		state.opJmp()
		//fixUpPalette_changeScreen(_res->_currentPart, _scriptVars[VAR_SCREEN_NUM]);
	} else {
		state.fetchWord()
	}
}

// Fade "palette number" - Changes of colour palette
func (state *VMState) opVidSetPalette() {
	index := state.fetchWord()
	fmt.Println("#opVidSetPalette", index)
	//TODO	_vid->_nextPal = num >> 8
}

//Text "text number", x, y, color - Displays in the work screen the specified text for the coordinates x,y.
func (state *VMState) opVidDrawString() {
	stringId := int(state.fetchWord())
	x := int(state.fetchByte())
	y := int(state.fetchByte())
	col := int(state.fetchByte())
	drawString(col, x, y, stringId)
}

//SetWS "Screen number" - Sets the work screen, which is where the polygons will be drawn by default.
func (state *VMState) opVidSelectPage() {
	page := int(state.fetchByte())
	setWorkPagePtr(page)
}

//Clr "Screen number", Color - Deletes a screen with one colour. Ingame, there are 4 screen buffers
func (state *VMState) opVidFillPage() {
	page := int(state.fetchByte())
	color := int(state.fetchByte())
	fillPage(page, color)
}

//Copy "Screen number A", "Screen number B" - Copies screen buffer A to screen buffer B.
func (state *VMState) opVidCopyPage() {
	source := state.fetchByte()
	dest := state.fetchByte()
	fmt.Println("#opVidCopyPage", source, dest)
	//TODO _vid->copyPage(i, j, _scriptVars[VAR_SCROLL_Y]);
}

//Show "Screen number" - Displays the screen buffer specified in the next video frame.
func (state *VMState) opVidUpdatePage() {
	page := int(state.fetchByte())
	//TODO inp_handleSpecialKeys();
	//TODO bypass protection, handle pause
	updateDisplay(page)
}

func (state *VMState) opVidDrawPolyBackground(opcode uint8) {
	offset := ((uint16(opcode) << 8) | uint16(state.fetchByte())) << 1
	//_res->_useSegVideo2 = false;
	posX := int(state.fetchByte())
	posY := int(state.fetchByte())
	height := posY - 199
	if height > 0 {
		posY = 199
		posX += height
	}
	fmt.Println("opVidUpdatePage", offset)
	setDataBuffer(int(offset))
	drawShape(0xFF, 0x40, posX, posY)
}

//Spr "'object name" , x, y, z - In the work screen, draws the graphics tool at the coordinates x,y and the zoom factor z. A polygon, a group of polygons...
func (state *VMState) opVidDrawPolySprite(opcode uint8) {
	offsetHi := state.fetchByte()
	offset := ((uint16(offsetHi) << 8) | uint16(state.fetchByte())) << 1
	posX := int16(state.fetchByte())
	//_res->_useSegVideo2 = false;
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
			//TODO hmm interesting...
			state.pc--
			fmt.Println("zoom decreased PC", state.pc)
			zoom = 0x40
		} else {
			zoom = uint16(state.variables[zoom])
		}
	} else {
		if opcode&1 > 0 {
			//_res->_useSegVideo2 = true;
			state.pc--
			fmt.Println("zoom decreased PC", state.pc)
			zoom = 0x40
		}
	}
	fmt.Printf("opVidDrawPolySprite")
	setDataBuffer(int(offset))
	drawShape(0xFF, int(zoom), int(posX), int(posY))
}

//Initialises a song.
func (state *VMState) opPlayMusic() {
	resNum := int(state.fetchWord())
	delay := int(state.fetchWord())
	pos := int(state.fetchByte())
	fmt.Printf("op_playMusic(0x%X, %d, %d)\n", resNum, delay, pos)
	//TODO snd_playMusic(resNum, delay, pos);
}

//Plays the sound file on one of the four game audio channels with specific height and volume.
func (state *VMState) opPlaySound() {
	resNum := int(state.fetchWord())
	freq := int(state.fetchByte())
	vol := int(state.fetchByte())
	channel := int(state.fetchByte())
	fmt.Printf("op_playSound(0x%X, %d, %d, %d)\n", resNum, freq, vol, channel)
	//TODO snd_playSound(resNum, freq, vol, channel);
}

func (state *VMState) opUpdateResource() {
	id := int(state.fetchWord())
	fmt.Println("opUpdateResource", id)
	if id >= GAME_PART_ID_1 {
		fmt.Println("should load next part", id)
		state.setupGamePart(id)
		return
	}
	state.assets.loadResource(id)
	/*
		if (num == 0) {
				_ply->stop();
				_mix->stopAll();
				_res->invalidateRes();
			} else {
				_res->update(num);
			}*/
	//TODO _res->update(num);
}
