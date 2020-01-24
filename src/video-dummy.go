//dummy video implementation, text output
package main

type DummyRenderer struct {
}

func (render DummyRenderer) blitPage(buffer [64000]Color) {
	Info(">VID: BLITPAGE")
}

func (render DummyRenderer) eventLoop(frameCount int) uint32 {
	Info(">VID: EVENTLOOP %i", frameCount)
	if frameCount < 128 {
		return 0
	}
	return KEY_ESC
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
