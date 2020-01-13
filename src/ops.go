package main

import (
	"fmt"
)

func (state *VMState) opJmp() {
	offset := state.fetchWord()
	state.pc = int(offset)
	fmt.Println("#op_jmp() jump to", state.pc)
}

func (state *VMState) opMovConst() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_movConst", index, value)
	state.variables[index] = value
}

func (state *VMState) opMov() {
	dest := int(state.fetchByte())
	source := int(state.fetchByte())
	fmt.Println("#op_mov", source, dest, state.variables[source])
	state.variables[dest] = state.variables[source]
}

func (state *VMState) opAdd() {
	dest := state.fetchByte()
	source := state.fetchByte()
	fmt.Println("#op_add()", dest, state.variables[source])
	state.variables[dest] += state.variables[source]
}

func (state *VMState) opSub() {
	dest := state.fetchByte()
	source := state.fetchByte()
	fmt.Println("#op_sub()", dest, state.variables[source])
	state.variables[dest] -= state.variables[source]
}

func (state *VMState) opAnd() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_and()", index, value)
	state.variables[index] = state.variables[index] & value
}

func (state *VMState) opOr() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_or()", index, value)
	state.variables[index] = state.variables[index] | value
}

func (state *VMState) opShl() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_shl()", index, value)
	state.variables[index] = state.variables[index] << value
}

func (state *VMState) opShr() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_shr()", index, value)
	state.variables[index] = state.variables[index] >> value
}

func (state *VMState) opAddConst() {
	//TODO add workaround for vm bug
	//		if (_res->_currentPart == 16006 && _scriptPtr.pc == _res->_segCode + 0x6D48) {
	//	warning("Script::op_addConst() workaround for infinite looping gun sound");
	//snd_playSound(0x5B, 1, 64, 1);
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#op_addConst()", index, value)
	state.variables[index] += value
}

func (state *VMState) opCall() {
	state.saveSP()
	state.pc = int(state.fetchWord())
	fmt.Println("#op_call(), jump to pc:", state.pc)
}

func (state *VMState) opRet() {
	state.restoreSP()
	fmt.Println("#op_ret(), pc:", state.pc)
}

func (state *VMState) opInstallTask() {
	index := state.fetchByte()
	value := int(state.fetchWord())
	fmt.Println("#opInstallTask", index, value)
	//	assert(i < 0x40);
	// TODO validate me: 	_scriptTasks[1][i] = n;
	state.channelData[index] = value
}

func (state *VMState) opYieldTask() {
	//TODO 	_scriptPaused = true;
	fmt.Println("#opYieldTask TODO")
}

func (state *VMState) opRemoveTask() {
	fmt.Println("#opRemoveTask TODO")
	//TODO _scriptPtr.pc = _res->_segCode + 0xFFFF;
	//TODO _scriptPaused = true;
}

func (state *VMState) opChangeTaskState() {
	j := state.fetchByte()
	i := state.fetchByte()
	a := state.fetchByte()
	fmt.Println("#opChangeTaskState TODO", j, i, a)
	//TODO _scriptPtr.pc = _res->_segCode + 0xFFFF;
	//TODO _scriptPaused = true;
}

func (state *VMState) opJmpIfVar() {
	index := state.fetchByte()
	state.variables[index]--
	if state.variables[index] != 0 {
		state.opJmp()
	} else {
		state.fetchWord()
	}
}

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

func (state *VMState) opVidSetPalette() {
	index := state.fetchWord()
	fmt.Println("#opVidSetPalette", index)
	//TODO	_vid->_nextPal = num >> 8
}

func (state *VMState) opVidDrawString() {
	stringId := int(state.fetchWord())
	x := int(state.fetchByte())
	y := int(state.fetchByte())
	col := int(state.fetchByte())
	drawString(col, x, y, stringId)
}

func (state *VMState) opVidSelectPage() {
	page := int(state.fetchByte())
	setWorkPagePtr(page)
}

func (state *VMState) opVidFillPage() {
	page := int(state.fetchByte())
	color := int(state.fetchByte())
	fillPage(page, color)
}

func (state *VMState) opVidCopyPage() {
	source := state.fetchByte()
	dest := state.fetchByte()
	fmt.Println("#opVidCopyPage", source, dest)
	//TODO _vid->copyPage(i, j, _scriptVars[VAR_SCROLL_Y]);
}

func (state *VMState) opVidUpdatePage() {
	page := state.fetchByte()
	fmt.Println("#opVidUpdatePage", page)
	//TODO inp_handleSpecialKeys();
	//TODO bypass protection, handle pause
	//_vid->updateDisplay(page, _stub);
}

func (state *VMState) opVidDrawPolyBackground(opcode uint8) {
	offset := ((int(opcode) << 8) | int(state.fetchByte())) << 1
	posX := int(state.fetchByte())
	posY := int(state.fetchByte())
	height := posY - 199
	if height > 0 {
		posY = 199
		posX += height
	}
	setDataBuffer(offset)
	drawShape(0xFF, 0x40, posX, posY)
}

func (state *VMState) opVidDrawPolySprite(opcode uint8) {
	offsetHi := state.fetchByte()
	offset := ((int(offsetHi) << 8) | int(state.fetchByte())) << 1
	posX := int(state.fetchByte())

	if opcode&0x20 == 0 {
		if opcode&0x10 == 0 {
			posX = (posX << 8) | int(state.fetchByte())
		} else {
			posX = state.variables[posX]
		}
	} else {
		if opcode&0x10 > 0 {
			posX += 0x100
		}
	}
	posY := int(state.fetchByte())
	if opcode&8 == 0 {
		if opcode&4 == 0 {
			posY = (posY << 8) | int(state.fetchByte())
		} else {
			posY = state.variables[posY]
		}
	}

	zoom := int(state.fetchByte())
	if opcode&2 == 0 {
		if opcode&1 == 0 {
			//TODO hmm interesting...
			state.pc--
			fmt.Println("zoom decreased PC", state.pc)
			zoom = 0x40
		} else {
			zoom = state.variables[zoom]
		}
	} else {
		if opcode&1 > 0 {
			//_res->_useSegVideo2 = true;
			state.pc--
			fmt.Println("zoom decreased PC", state.pc)
			zoom = 0x40
		}
	}
	setDataBuffer(offset)
	drawShape(0xFF, zoom, posX, posY)
}

func (state *VMState) opPlayMusic() {
	resNum := int(state.fetchWord())
	delay := int(state.fetchWord())
	pos := int(state.fetchByte())
	fmt.Printf("op_playMusic(0x%X, %d, %d)\n", resNum, delay, pos)
	//TODO snd_playMusic(resNum, delay, pos);
}

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
