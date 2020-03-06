package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

type WASMHAL struct {
	keyCombo uint32
	canvas   Canvas
}

func buildWASMHAL() *WASMHAL {
	return &WASMHAL{
		keyCombo: anotherworld.KeyNone,
		canvas:   GetCanvas("gotherworld"),
	}
}

func (render *WASMHAL) updateKeyStateFrom(keyMap *map[uint32]bool) {
	render.keyCombo = anotherworld.KeyNone
	for key, pressed := range *keyMap {
		if pressed {
			render.keyCombo |= key
		}
	}
}

func (render *WASMHAL) BlitPage(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]anotherworld.Color, posX, posY int) {
	logger.Debug(">VID: BLITPAGE %d %d", posX, posY)
	// TODO implement color handling in js-canvas.go
	// TODO draw points passing x, y and color see hal-sdl.go
	x := posX
	y := posY
	render.canvas.DrawPoint(buffer[0], x, y)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render *WASMHAL) EventLoop(frameCount int) uint32 {
	if render.keyCombo != 0 {
		logger.Info(">EVNT: EVENTLOOP %d, KeyCombo state: %v", frameCount, render.keyCombo)
	}
	return render.keyCombo
}

func (render *WASMHAL) Shutdown() {
	//nothin to see here, move on!
}

func (render *WASMHAL) PlayMusic(resNum, delay, pos int) {
	logger.Info(">SND: playMusic res:%d del:%d pos:%d", resNum, delay, pos)
}

func (render *WASMHAL) PlaySound(resNum, freq, vol, channel int) {
	logger.Info(">SND: playSound res:%d frq:%d vol:%d c:%d", resNum, freq, vol, channel)
}
