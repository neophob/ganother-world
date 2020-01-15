//dummy video implementation, text output
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH  int32 = 800
	HEIGHT int32 = 600
)

type SDLRenderer struct {
	surface *sdl.Surface
	window  *sdl.Window
	exitReq bool
}

func buildSDLRenderer() *SDLRenderer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	//	defer sdl.Quit()

	window, err := sdl.CreateWindow("ganother world", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT, sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}
	//	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	rect := sdl.Rect{0, 0, WIDTH, HEIGHT}
	surface.FillRect(&rect, 0xff000000)
	window.UpdateSurface()

	return &SDLRenderer{surface: surface, window: window}
}

//TODO where is the stringId defined?
func (render SDLRenderer) drawString(color, posX, posY, stringId int) {
	text := getText(stringId)
	fmt.Printf(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s\n", color, posX, posY, text)
}

func (render SDLRenderer) drawShape(color, zoom, posX, posY int) {
	fmt.Printf(">VID: DRAWSHAPE color:%d, x:%d, y:%d, zoom:%d\n", color, posX, posY, zoom)
}

func (render SDLRenderer) fillPage(page, color int) {
	fmt.Println(">VID: FILLPAGE", page, color)
	//_graphics->clearBuffer(getPagePtr(page), color);
}

func (render SDLRenderer) copyPage(src, dst, vscroll int) {
	fmt.Println(">VID: COPYPAGE", src, dst, vscroll)
}

// blit
func (render SDLRenderer) updateDisplay(page int) {
	fmt.Println(">VID: UPDATEDISPLAY", page)
}

//TODO gimme a better name
func (render SDLRenderer) setDataBuffer(offset int) {
	fmt.Println(">VID: SETDATABUFFER", offset)
}

func (render SDLRenderer) setWorkPagePtr(page int) {
	fmt.Println(">VID: SETWORKPAGEPTR", page)
}

func (render SDLRenderer) setPalette(index int) {
	fmt.Println(">VID: SETPALETTE", index>>8)
	//TODO	_vid->_nextPal = num >> 8
}

func (render *SDLRenderer) mainLoop() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			render.exitReq = true
			fmt.Println(">render.exitReq", render.exitReq)
		default:
			//fmt.Println("SDL EVENT", event)

		}
	}
}

func (render SDLRenderer) shutdown() {
	render.window.Destroy()
	sdl.Quit()
}

func (render SDLRenderer) exitRequested(frameCount int) bool {
	return render.exitReq
}
