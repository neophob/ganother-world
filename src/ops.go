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
