package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

type WASMHAL struct {
	keyCombo uint32
}

func buildWASMHAL() *WASMHAL {
	return &WASMHAL{
		keyCombo: anotherworld.KeyNone,
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
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render *WASMHAL) EventLoop(frameCount int) uint32 {
	logger.Info(">EVNT: EVENTLOOP %d, KeyCombo state: %v", frameCount, render.keyCombo)
	if frameCount >= MAX_FRAME_COUNT {
		logger.Info(">EVNT: Max frameCount reached (%d) triggering KeyEsc", frameCount)
		return anotherworld.KeyEsc
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
