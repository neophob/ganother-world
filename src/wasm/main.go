package main

import (
	"syscall/js"
	"github.com/neophob/ganother-world/gamelib"
)

func addTitle() {
	gameTitle := gamelib.GetTitle()

	document := js.Global().Get("document")
	h1 := document.Call("createElement", "h1")
	h1.Set("innerHTML", gameTitle)
	document.Get("body").Call("appendChild", h1)
}

func main() {
	gameTitle := gamelib.GetTitle()

	println("Hello " + gameTitle)

	addTitle()

	channel := make(chan bool)
	<-channel
}
