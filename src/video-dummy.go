//dummy video implementation, text output
package main

type DummyRenderer struct {
}

func (render DummyRenderer) blitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	Info(">VID: BLITPAGE %d %d", posX, posY)
}

func (render DummyRenderer) eventLoop(frameCount int) uint32 {
	Info(">VID: EVENTLOOP %d", frameCount)
	if frameCount < 128 {
		return 0
	}
	return KEY_ESC
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
