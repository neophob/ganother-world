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

func (state *VMState) opCall() {
	offset := state.fetchWord()
	state.saveCurrentSP()
	state.pc = int(offset)
	fmt.Println("#op_call() jump to", state.pc)
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

func (state *VMState) opVidDrawPolyBackground(opcode uint8) {
	offset := ((int(opcode) << 8) | int(state.fetchByte())) << 1
	posX := state.fetchByte()
	posY := state.fetchByte()
	height := posY - 199
	if height > 0 {
		posY = 199
		posX += height
	}
	fmt.Println("DRAW_POLY_BACKGROUND", posX, posY, offset)
	//			_vid->setDataBuffer(_res->_segVideo1, off);
	//			_vid->drawShape(0xFF, 0x40, &pt);
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

	zoom := state.fetchByte()
	fmt.Printf("DRAW_POLY_SPRITE x:%d, y:%d, offset:%d, zoom:%d\n", posX, posY, offset, zoom)
}
