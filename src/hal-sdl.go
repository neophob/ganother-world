//SDL video implementation
package main

import (
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
		Error("SDL INIT FAILED")
		panic(err)
	}

	window, err := sdl.CreateWindow("ganother world", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		windowWidth, windowHeight, sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		Error("SDL CREATE WINDOW FAILED")
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		Error("SDL CREATE RENDERER FAILED")
		panic(err)
	}
	renderer.SetLogicalSize(WIDTH, HEIGHT)
	//renderer.SetLogicalSize(WIDTH*2, HEIGHT*2)
	renderer.Clear()
	renderer.Present()

	surface, err := window.GetSurface()
	if err != nil {
		Error("SDL GET SURFACE FAILED")
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
func (render *SDLHAL) blitPage(buffer [WIDTH * HEIGHT]Color, posX, posY int) {
	lastSetColor := buffer[0]
	render.renderer.SetDrawColor(buffer[0].r, buffer[0].g, buffer[0].b, 255)
	offset := 0
	for y := 0; y < int(HEIGHT); y++ {
		for x := 0; x < int(WIDTH); x++ {
			if i := buffer[offset]; i != lastSetColor {
				render.renderer.SetDrawColor(i.r, i.g, i.b, 255)
				lastSetColor = i
			}
			render.renderer.DrawPoint(int32(x+posX), int32(y+posY))
			offset++
		}
	}
	render.renderer.Present()
}

// check keyboard input
func (render *SDLHAL) eventLoop(frameCount int) uint32 {
	keyPress := uint32(0x0)
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			keyPress |= KeyEsc
			Debug(">ESC")
		case *sdl.KeyboardEvent:
			Debug(">KeyboardEvent %v %v", t.State, t.Keysym.Scancode)
			isKeyPressed := t.State == sdl.PRESSED
			if t.Keysym.Sym == sdl.K_ESCAPE {
				keyPress |= KeyEsc
				render.holdKeys[KeyEsc] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_LEFT {
				keyPress |= KeyLeft
				render.holdKeys[KeyLeft] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_RIGHT {
				keyPress |= KeyRight
				render.holdKeys[KeyRight] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_UP {
				keyPress |= KeyUp
				render.holdKeys[KeyUp] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_DOWN {
				keyPress |= KeyDown
				render.holdKeys[KeyDown] = isKeyPressed
			}
			if t.Keysym.Sym == sdl.K_SPACE {
				keyPress |= KeyFire
				render.holdKeys[KeyFire] = isKeyPressed
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_p {
				keyPress |= KeyPause
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_s {
				keyPress |= KeySave
			}
			if isKeyPressed && t.Keysym.Sym == sdl.K_l {
				keyPress |= KeyLoad
			}
		}
	}

	if render.holdKeys[KeyEsc] {
		keyPress |= KeyEsc
	}
	if render.holdKeys[KeyLeft] {
		keyPress |= KeyLeft
	}
	if render.holdKeys[KeyRight] {
		keyPress |= KeyRight
	}
	if render.holdKeys[KeyUp] {
		keyPress |= KeyUp
	}
	if render.holdKeys[KeyDown] {
		keyPress |= KeyDown
	}
	if render.holdKeys[KeyFire] {
		keyPress |= KeyFire
	}

	return keyPress
}

// exit application, lets cleanup...
func (render *SDLHAL) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}
