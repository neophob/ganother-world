package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

//WASMHAL implementation, text output
type WASMHAL struct {
}

func (render WASMHAL) BlitPage(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]anotherworld.Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render WASMHAL) EventLoop(frameCount int) uint32 {
	// TODO get key events from JS here
	logger.Info(">EVNT: EVENTLOOP %d", frameCount)
	if frameCount < MAX_FRAME_COUNT {
		return anotherworld.KeyNone
	}
	logger.Info(">EVNT: Max frameCount reached (%d) triggering KeyEsc", frameCount)
	return anotherworld.KeyEsc
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
