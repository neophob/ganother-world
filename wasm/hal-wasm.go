package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const MAX_FRAME_COUNT = 32

type WASMHAL struct {
	keyCombo     uint32
	canvas       Canvas
	lastSetColor anotherworld.Color
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
	// logger.Debug(">VID: BLITPAGE %d %d", posX, posY)
	offset := 0
	for y := 0; y < int(anotherworld.HEIGHT); y++ {
		for x := 0; x < int(anotherworld.WIDTH); x++ {
			if color := buffer[offset]; color != render.lastSetColor {
				render.canvas.SetColor(color)
				render.lastSetColor = color
			}
			render.canvas.DrawPoint(x+posX, y+posY)
			offset += 1
		}
	}
}

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
