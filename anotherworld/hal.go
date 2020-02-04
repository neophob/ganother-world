package anotherworld

//Virtual Key mapping
const (
	KeyEsc   uint32 = 0x1
	KeyLeft  uint32 = 0x2
	KeyRight uint32 = 0x4
	KeyUp    uint32 = 0x8
	KeyDown  uint32 = 0x10
	KeyFire  uint32 = 0x20
	KeyPause uint32 = 0x40
	KeySave  uint32 = 0x80
	KeyLoad  uint32 = 0x100
)

//HAL is the implementation abstraction (audio, video, io)
type HAL interface {
	BlitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int)
	EventLoop(frameCount int) uint32
	Shutdown()

	PlaySound(resNum, freq, vol, channel int)
	PlayMusic(resNum, delay, pos int)
}
