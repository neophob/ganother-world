package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

type WASMHAL struct {
	// TODO
}

func (render WASMHAL) BlitPage(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]anotherworld.Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render WASMHAL) EventLoop(frameCount int) uint32 {
	/*
		TODO Implement Key handling from JS, should js pass in a map
		can we fetch a map via callback's callback. Maybe update from js
		on event then use updated version here?
		Where would the state be stored and how does the HAL access it?
			K_ESCAPE - KeyEsc
			K_LEFT - KeyLeft
			K_RIGHT - KeyRight
			K_UP - KeyUp
			K_DOWN - KeyDown
			K_SPACE - KeyFire
			K_p - KeyPause
			K_s - KeySave
			K_l - KeyLoad
		See hal-sdl for multiple key holds handling.
	*/
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
