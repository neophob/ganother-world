//dummy video implementation, text output
package main

type DummyRenderer struct {
}

func (render DummyRenderer) blitPage(buffer [64000]Color) {
	Info(">VID: BLITPAGE")
}

func (render DummyRenderer) eventLoop(frameCount int) bool {
	Info(">VID: EVENTLOOP %i", frameCount)
	return frameCount > 128
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
