//SDL video implementation
package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  int32 = 320 * 2
	windowHeight int32 = 200 * 2
)

//SDLHAL implements the HAL using ... SDL2
type SDLHAL struct {
	surface  *sdl.Surface
	renderer *sdl.Renderer
	window   *sdl.Window
	holdKeys map[uint32]bool
}

func buildSDLHAL() *SDLHAL {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		logger.Error("SDL INIT FAILED")
		panic(err)
	}

	window, err := sdl.CreateWindow("ganother world", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		windowWidth, windowHeight, sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		logger.Error("SDL CREATE WINDOW FAILED")
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		logger.Error("SDL CREATE RENDERER FAILED")
		panic(err)
	}
	renderer.SetLogicalSize(anotherworld.WIDTH, anotherworld.HEIGHT)
	renderer.Clear()
	renderer.Present()

	surface, err := window.GetSurface()
	if err != nil {
		logger.Error("SDL GET SURFACE FAILED")
		panic(err)
	}

	return &SDLHAL{
		surface:  surface,
		window:   window,
		renderer: renderer,
		holdKeys: make(map[uint32]bool),
	}
}

// blit image
func (render *SDLHAL) BlitPage(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]anotherworld.Color, posX, posY int) {
	lastSetColor := buffer[0]
	render.renderer.SetDrawColor(buffer[0].R, buffer[0].G, buffer[0].B, 255)
	offset := 0
	for y := 0; y < int(anotherworld.HEIGHT); y++ {
		for x := 0; x < int(anotherworld.WIDTH); x++ {
			if i := buffer[offset]; i != lastSetColor {
				render.renderer.SetDrawColor(i.R, i.G, i.B, 255)
				lastSetColor = i
			}
			render.renderer.DrawPoint(int32(x+posX), int32(y+posY))
			offset++
		}
	}
	render.renderer.Present()
}

// check keyboard input
func (render *SDLHAL) EventLoop(frameCount int) uint32 {
	keyPress := uint32(0x0)
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			keyPress |= anotherworld.KeyEsc
			logger.Debug(">ESC")
		case *sdl.KeyboardEvent:
			logger.Debug(">KeyboardEvent %v %v", t.State, t.Keysym.Scancode)
			isKeyPressed := t.State == sdl.PRESSED
			if t.Keysym.Sym == sdl.K_ESCAPE {
				keyPress |= anotherworld.KeyEsc
				render.holdKeys[anotherworld.KeyEsc] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_LEFT {
				keyPress |= anotherworld.KeyLeft
				render.holdKeys[anotherworld.KeyLeft] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_RIGHT {
				keyPress |= anotherworld.KeyRight
				render.holdKeys[anotherworld.KeyRight] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_UP {
				keyPress |= anotherworld.KeyUp
				render.holdKeys[anotherworld.KeyUp] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_DOWN {
				keyPress |= anotherworld.KeyDown
				render.holdKeys[anotherworld.KeyDown] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_SPACE {
				keyPress |= anotherworld.KeyFire
				render.holdKeys[anotherworld.KeyFire] = isKeyPressed
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_p {
				keyPress |= anotherworld.KeyPause
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_s {
				keyPress |= anotherworld.KeySave
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_l {
				keyPress |= anotherworld.KeyLoad
			}
		}
	}

	if render.holdKeys[anotherworld.KeyEsc] {
		keyPress |= anotherworld.KeyEsc
	}
	if render.holdKeys[anotherworld.KeyLeft] {
		keyPress |= anotherworld.KeyLeft
	}
	if render.holdKeys[anotherworld.KeyRight] {
		keyPress |= anotherworld.KeyRight
	}
	if render.holdKeys[anotherworld.KeyUp] {
		keyPress |= anotherworld.KeyUp
	}
	if render.holdKeys[anotherworld.KeyDown] {
		keyPress |= anotherworld.KeyDown
	}
	if render.holdKeys[anotherworld.KeyFire] {
		keyPress |= anotherworld.KeyFire
	}

	return keyPress
}

func (render *SDLHAL) PlayMusic(resNum, delay, pos int) {
	logger.Info(">SND: playMusic res:%d del:%d pos:%d", resNum, delay, pos)
}

func (render *SDLHAL) PlaySound(resNum, freq, vol, channel int) {
	logger.Info(">SND: playSound res:%d frq:%d vol:%d c:%d", resNum, freq, vol, channel)
}

// exit application, lets cleanup...
func (render *SDLHAL) Shutdown() {
	render.window.Destroy()
	sdl.Quit()
}
