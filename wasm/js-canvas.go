package main

import (
	"syscall/js"
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

func (c Canvas) FillRect() {
	c.context2d.Call("fillRect", 130, 190, 40, 60)
}
