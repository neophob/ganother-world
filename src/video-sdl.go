//SDL video implementation
//TODO: implement multiple pages
//TODO: split me away from the core lib here
package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  int32 = 320 * 2
	WINDOW_HEIGHT int32 = 200 * 2
)

type SDLRenderer struct {
	surface  *sdl.Surface
	renderer *sdl.Renderer
	window   *sdl.Window
}

func buildSDLRenderer() *SDLRenderer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		Error("SDL INIT FAILED")
		panic(err)
	}

	window, err := sdl.CreateWindow("ganother world", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		Error("SDL CREATE WINDOW FAILED")
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		Error("SDL CREATE RENDERER FAILED")
		panic(err)
	}
	//renderer.SetLogicalSize(WIDTH, HEIGHT)
	renderer.SetLogicalSize(WIDTH*2, HEIGHT*2)
	renderer.Clear()
	renderer.Present()

	surface, err := window.GetSurface()
	if err != nil {
		Error("SDL GET SURFACE FAILED")
		panic(err)
	}

	return &SDLRenderer{
		surface:  surface,
		window:   window,
		renderer: renderer,
	}
}

// blit image
func (render *SDLRenderer) blitPage(buffer [64000]Color, posX, posY int) {
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
func (render *SDLRenderer) eventLoop(frameCount int) uint32 {
	keyPress := uint32(0x0)
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			keyPress |= KEY_ESC
			Debug(">ESC")
		case *sdl.KeyboardEvent:
			Debug(">KeyboardEvent %v", t.Keysym.Sym)
			if t.Keysym.Sym == sdl.K_ESCAPE && t.State == 1 {
				keyPress |= KEY_ESC
			}
			if t.Keysym.Sym == sdl.K_LEFT && t.State == 1 {
				keyPress |= KEY_LEFT
			}
			if t.Keysym.Sym == sdl.K_RIGHT && t.State == 1 {
				keyPress |= KEY_RIGHT
			}
			if t.Keysym.Sym == sdl.K_UP && t.State == 1 {
				keyPress |= KEY_UP
			}
			if t.Keysym.Sym == sdl.K_DOWN && t.State == 1 {
				keyPress |= KEY_DOWN
			}
			if t.Keysym.Sym == sdl.K_SPACE && t.State == 1 {
				keyPress |= KEY_FIRE
			}
		}
	}
	return keyPress
}

// exit application, lets cleanup...
func (render *SDLRenderer) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}
