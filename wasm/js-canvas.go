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
	canvasElement := js.Global().Get("document").Call("getElementById", domElementId)
	context2d := canvasElement.Call("getContext", "2d")

	return Canvas{
		context2d: context2d,
	}
}

func (c Canvas) SetColor(color anotherworld.Color) {
	c.context2d.Set("fillStyle", fmt.Sprintf("rgb(%d, %d, %d)", color.R, color.G, color.B))
}

func (c Canvas) DrawPoint(x, y int) {
	c.context2d.Call("fillRect", x, y, 1, 1)
}

func (c Canvas) fillStyle(style string) {
	c.context2d.Set("fillStyle", style)
}

func (c Canvas) fillRect(x, y, width, height int) {
	c.context2d.Call("fillRect", x, y, width, height)
}

func (c Canvas) blitIt(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]int) {
	var interfaceSlice []interface{} = make([]interface{}, len(buffer))
	for i, d := range buffer {
		interfaceSlice[i] = d
	}
	c.context2d.Call("blit", interfaceSlice)
}
