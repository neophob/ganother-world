package main

import (
	"fmt"
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

// Notes:
// Consider drawImage, if we can generate the correct format to output in 1 call.
// https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/drawImage
// OR
// access the canvas pixel data directly might be faster - at least for unscaled output
//  var id = ctx.getImageData(0, 0, canvasWidth, canvasHeight)
//  var pixels = id.data
//  pixels[0] = 0
//  pixels[1] = 255
//  pixels[2] = 255
//  pixels[3] = 255

func (c Canvas) SetColor(color anotherworld.Color) {
	c.fillStyle(fmt.Sprintf("#%X%X%X", color.R, color.G, color.B))
}

func (c Canvas) DrawPoint(x, y int) {
	c.context2d.Call("fillRect", x, y, 1, 1)
}

func (c Canvas) fillStyle(style string) {
	c.context2d.Set("fillStyle", style)
}

func (c Canvas) fillRect(x, y, length, width int) {
	c.context2d.Call("fillRect", x, y, length, width)
}
