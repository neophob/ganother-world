package anotherworld

import "github.com/neophob/ganother-world/logger"

//DummyHAL implementation, text output
type DummyHAL struct {
}

func (render DummyHAL) BlitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

func (render DummyHAL) EventLoop(frameCount int) uint32 {
	logger.Info(">VID: EVENTLOOP %d", frameCount)
	if frameCount < 128 {
		return 0
	}
	return KeyEsc
}

func (render DummyHAL) Shutdown() {
	//nothin to see here, move on!
}

func (render DummyHAL) PlayMusic(resNum, delay, pos int) {
	logger.Info(">SND: playMusic res:%d del:%d pos:%d", resNum, delay, pos)
}

func (render DummyHAL) PlaySound(resNum, freq, vol, channel int) {
	logger.Info(">SND: playSound res:%d frq:%d vol:%d c:%d", resNum, freq, vol, channel)
}
