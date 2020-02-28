package anotherworld

import "github.com/neophob/ganother-world/logger"

const MAX_FRAME_COUNT = 128

//DummyHAL implementation, text output
type DummyHAL struct {
}

func (render DummyHAL) BlitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	logger.Info(">VID: BLITPAGE %d %d", posX, posY)
}

//Outputs framecount, sends escape after MAX_FRAME_COUNT frames and ends the game
func (render DummyHAL) EventLoop(frameCount int) uint32 {
	logger.Info(">EVNT: EVENTLOOP %d", frameCount)
	if frameCount < MAX_FRAME_COUNT {
		return KeyNone
	}
	logger.Info(">EVNT: Max frameCount reached (%d) triggering KeyEsc", frameCount)
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
