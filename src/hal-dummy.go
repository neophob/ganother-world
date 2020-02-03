package main

//DummyHAL implementation, text output
type DummyHAL struct {
}

func (render DummyHAL) blitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	Info(">VID: BLITPAGE %d %d", posX, posY)
}

func (render DummyHAL) eventLoop(frameCount int) uint32 {
	Info(">VID: EVENTLOOP %d", frameCount)
	if frameCount < 128 {
		return 0
	}
	return KEY_ESC
}

func (render DummyHAL) shutdown() {
	//nothin to see here, move on!
}