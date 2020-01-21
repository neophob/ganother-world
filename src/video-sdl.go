//SDL video implementation
//TODO: implement multiple pages
//TODO: split me away from the core lib here
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  int32 = 320 * 4
	WINDOW_HEIGHT int32 = 200 * 3
)

type SDLRenderer struct {
	surface  *sdl.Surface
	renderer *sdl.Renderer
	window   *sdl.Window
}

func buildSDLRenderer() *SDLRenderer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("ganother world", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}
	renderer.SetLogicalSize(WIDTH, HEIGHT)
	renderer.Clear()
	renderer.Present()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	return &SDLRenderer{
		surface:  surface,
		window:   window,
		renderer: renderer,
	}
}

// blit image
func (render *SDLRenderer) blitPage(buffer [64000]Color) {
	lastSetColor := buffer[0]
	render.renderer.SetDrawColor(buffer[0].r, buffer[0].g, buffer[0].b, 255)
	offset := 0
	for y := 0; y < int(HEIGHT); y++ {
		for x := 0; x < int(WIDTH); x++ {
			if i := buffer[offset]; i != lastSetColor {
				render.renderer.SetDrawColor(i.r, i.g, i.b, 255)
				lastSetColor = i
			}
			render.renderer.DrawPoint(int32(x), int32(y))
			offset++
		}
	}
	render.renderer.Present()
}

// check if app should exit, needs more functionality soon (key values..)
func (render *SDLRenderer) eventLoop(frameCount int) bool {
	exitRequested := false
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			exitRequested = true
			fmt.Println(">render.exitAppReq", exitRequested)
		case *sdl.KeyboardEvent:
			fmt.Println(">render.exitAppReq2", t.Keysym.Sym)
			if t.Keysym.Sym == sdl.K_ESCAPE && t.State == 1 {
				exitRequested = true
			}
		}
	}
	return exitRequested
}

// cleanup
func (render *SDLRenderer) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}
