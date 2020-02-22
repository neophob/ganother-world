package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

type WASMHAL struct {
	holdKeys map[uint32]bool
}

type KeyEvent struct {
	key     string
	keyCode int
	pressed bool
}

func buildWASMHAL() *WASMHAL {
	return &WASMHAL{
		holdKeys: make(map[uint32]bool),
	}
}

func (render WASMHAL) setKeyState(key *KeyEvent) {
	// TODO switch based on keyCode
	// then set render.holdKeys based on type
}

func (render WASMHAL) BlitPage(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]anotherworld.Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render WASMHAL) EventLoop(frameCount int) uint32 {
	keyPress := anotherworld.KeyNone
	logger.Info(">EVNT: EVENTLOOP %d", frameCount)
	if frameCount >= MAX_FRAME_COUNT {
		logger.Info(">EVNT: Max frameCount reached (%d) triggering KeyEsc", frameCount)
		return anotherworld.KeyEsc
	}
	// TODO get details from setKeyState and flatten onto 1 var.
	return keyPress
}

func (render WASMHAL) Shutdown() {
	//nothin to see here, move on!
}

func (render WASMHAL) PlayMusic(resNum, delay, pos int) {
	logger.Info(">SND: playMusic res:%d del:%d pos:%d", resNum, delay, pos)
}

func (render WASMHAL) PlaySound(resNum, freq, vol, channel int) {
	logger.Info(">SND: playSound res:%d frq:%d vol:%d c:%d", resNum, freq, vol, channel)
}
