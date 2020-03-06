package main

import (
	"syscall/js"

	"github.com/neophob/ganother-world/anotherworld"
)

type Canvas struct {
	context2d js.Value
}

func GetCanvas(domElementId string) Canvas {
	canvasEl := js.Global().Get("document").Call("getElementById", domElementId)
	canvas2d := canvasEl.Call("getContext", "2d")

	canvas := Canvas{
		context2d: canvas2d,
	}

	return canvas
}

// Consider drawImage, if we can generate the correct format to output in 1 call.
// https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/drawImage

func (c Canvas) DrawPoint(color anotherworld.Color, x, y int) {
	// TODO convert color to hex RGB
	// c.fillStyle(#rgbColor)
	c.fillRect(x, y, 1, 1)
}

func (c Canvas) fillStyle(style string) {
	c.context2d.Set("fillStyle", style)
}

func (c Canvas) fillRect(x, y, length, width int) {
	c.context2d.Call("fillRect", x, y, length, width)
}
