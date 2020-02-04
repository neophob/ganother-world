package main

import "github.com/neophob/ganother-world/logger"

//DummyHAL implementation, text output
type DummyHAL struct {
}

func (render DummyHAL) blitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

func (render DummyHAL) eventLoop(frameCount int) uint32 {
	logger.Info(">VID: EVENTLOOP %d", frameCount)
	if frameCount < 128 {
		return 0
	}
	return KeyEsc
}

func (render DummyHAL) shutdown() {
	//nothin to see here, move on!
}
