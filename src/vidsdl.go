//SDL video implementation
//TODO: implement multiple pages
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  int32 = 320 * 3
	WINDOW_HEIGHT int32 = 200 * 3

	WIDTH  int32 = 320
	HEIGHT int32 = 200
)

type SDLRenderer struct {
	surface  *sdl.Surface
	renderer *sdl.Renderer
	window   *sdl.Window
}

func buildSDLRenderer() SDLRenderer {
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

	return SDLRenderer{
		surface:  surface,
		window:   window,
		renderer: renderer,
	}
}

func (render SDLRenderer) setColor(col Color) {
	fmt.Println(">VID: SETCOLOR", col)
	render.renderer.SetDrawColor(col.r, col.g, col.g, 255)
}

//TODO move me to video?
func (render SDLRenderer) drawChar(posX, posY int32, char byte) {
	ofs := 8 * (int32(char) - 0x20)
	for j := int32(0); j < 8; j++ {
		ch := FONT[ofs+j]
		for i := int32(0); i < 8; i++ {
			if ch&(1<<(7-i)) > 0 {
				render.renderer.DrawPoint(posX+i, posY+j)
			}
		}
	}
}

func (render SDLRenderer) fillPage(page int) {
	fmt.Println("XXXLLKKKK", page)
	render.renderer.FillRect(nil)
}

func (render SDLRenderer) copyPage(src, dst int) {
	fmt.Println(">VID: COPYPAGE", src, dst)
}

// blit
func (render SDLRenderer) blitPage(page int) {
	fmt.Println(">VID: UPDATEDISPLAY", page)
	render.renderer.Present()
}

func (render SDLRenderer) eventLoop(frameCount int) bool {
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

func (render SDLRenderer) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}

func (render SDLRenderer) drawFilledPolygons(vx, vy []int16, col Color) {
	gfx.FilledPolygonColor(render.renderer, vx, vy, sdl.Color{col.r, col.g, col.g, 255})
}
