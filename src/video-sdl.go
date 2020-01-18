//SDL video implementation
//TODO: implement multiple pages
//TODO: split me away from the core lib here
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
	surface        *sdl.Surface
	screenBuffers  [4]*sdl.Surface
	renderer       *sdl.Renderer
	screenRenderer [4]*sdl.Renderer
	window         *sdl.Window
	drawColor      Color
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

	screenBuffers := [4]*sdl.Surface{
		buildBuffer(surface),
		buildBuffer(surface),
		buildBuffer(surface),
		buildBuffer(surface),
	}

	screenRenderer := [4]*sdl.Renderer{
		buildRenderer(screenBuffers[0]),
		buildRenderer(screenBuffers[1]),
		buildRenderer(screenBuffers[2]),
		buildRenderer(screenBuffers[3]),
	}

	return &SDLRenderer{
		surface:        surface,
		window:         window,
		renderer:       renderer,
		screenBuffers:  screenBuffers,
		screenRenderer: screenRenderer,
	}
}

func buildBuffer(surface *sdl.Surface) *sdl.Surface {
	buffer, _ := sdl.CreateRGBSurface(0, 320, 200, 32, surface.Format.Rmask, surface.Format.Gmask, surface.Format.Bmask, surface.Format.Amask)
	//	buffer, _ := sdl.CreateRGBSurface(0, WINDOW_WIDTH, WINDOW_HEIGHT, 32, surface.Format.Rmask, surface.Format.Gmask, surface.Format.Bmask, surface.Format.Amask)
	return buffer
}

func buildRenderer(surface *sdl.Surface) *sdl.Renderer {
	renderer, _ := sdl.CreateSoftwareRenderer(surface)
	renderer.SetLogicalSize(WIDTH, HEIGHT)
	renderer.Clear()
	return renderer
}

func (render *SDLRenderer) setColor(col Color) {
	fmt.Println(">VID: SETCOLOR", col)
	render.drawColor = col
	render.screenRenderer[0].SetDrawColor(col.r, col.g, col.b, 255)
	render.screenRenderer[1].SetDrawColor(col.r, col.g, col.b, 255)
	render.screenRenderer[2].SetDrawColor(col.r, col.g, col.b, 255)
	render.screenRenderer[3].SetDrawColor(col.r, col.g, col.b, 255)
}

func (render *SDLRenderer) fillPage(page int) {
	fmt.Println(">VID: FILLPAGE", page, render.drawColor)
	target := render.screenBuffers[page]
	err := target.FillRect(nil, render.drawColor.toUint32())
	if err != nil {
		panic(err)
	}
}

func (render *SDLRenderer) copyPage(src, dst int) {
	fmt.Println(">VID: COPYPAGE", src, dst)
	surfaceSrc := render.screenBuffers[src]
	surfaceDest := render.screenBuffers[dst]
	err := surfaceSrc.Blit(nil, surfaceDest, nil)
	if err != nil {
		panic(err)
	}
	renderer := render.screenRenderer[dst]
	renderer.Present()
}

// blit
func (render *SDLRenderer) blitPage(page int) {
	/*	fmt.Println(">VID: UPDATEDISPLAY", page)
		surface := render.screenBuffers[page]
		err := surface.Blit(nil, render.surface, nil)
		if err != nil {
			panic(err)
		}*/

	//rect := sdl.Rect{0, 0, 320, 200}
	render.screenBuffers[0].Blit(nil, render.surface, nil)

	r1 := &sdl.Rect{320, 0, 320, 200}
	render.screenBuffers[1].Blit(nil, render.surface, r1)

	r2 := &sdl.Rect{0, 200, 320, 200}
	render.screenBuffers[2].Blit(nil, render.surface, r2)

	r3 := &sdl.Rect{320, 200, 320, 200}
	render.screenBuffers[3].Blit(nil, render.surface, r3)

	r4 := &sdl.Rect{640, 100, 320, 200}
	render.screenBuffers[page].Blit(nil, render.surface, r4)

	render.renderer.Present()
}

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

func (render *SDLRenderer) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}

func (render *SDLRenderer) drawFilledPolygons(page int, vx, vy []int16, col Color) {
	renderer := render.screenRenderer[page]
	gfx.FilledPolygonColor(renderer, vx, vy, sdl.Color{col.r, col.g, col.g, 255})
}

func (render *SDLRenderer) drawPixel(page int, posX, posY int32) {
	renderer := render.screenRenderer[page]
	renderer.DrawPoint(posX, posY)
}
