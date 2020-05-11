package main

import (
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

func (c Canvas) blitIt(buffer [anotherworld.WIDTH * anotherworld.HEIGHT]int) {
	var interfaceSlice []interface{} = make([]interface{}, len(buffer))
	for i, d := range buffer {
		interfaceSlice[i] = d
	}
	c.context2d.Call("blit", interfaceSlice)
}
