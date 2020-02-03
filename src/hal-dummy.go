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
	return KeyEsc
}

func (render DummyHAL) shutdown() {
	//nothin to see here, move on!
}

func (render DummyHAL) playMusic(resNum, delay, pos int) {
	Info(">SND: playMusic res:%d del:%d pos:%d", resNum, delay, pos)
}

func (render DummyHAL) playSound(resNum, freq, vol, channel int) {
	Info(">SND: playSound res:%d frq:%d vol:%d c:%d", resNum, freq, vol, channel)
}
