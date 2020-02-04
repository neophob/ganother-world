package main

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

//HAL implements HAL (audio, video, io)
type HAL interface {
	blitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int)
	eventLoop(frameCount int) uint32
	shutdown()
}
