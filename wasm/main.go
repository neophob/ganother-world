package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
	"syscall/js"
)

func init() {
	logger.SetLogLevel(logger.LEVEL_INFO)
	addTitle()
}

func main() {
	// TODO make logger wasm/browser console friendly
	logger.Info("WASM " + anotherworld.GetTitle())

	app := InitGame()
	defer app.Shutdown()

	channel := make(chan bool)
	<-channel
}

func addTitle() {
	document := js.Global().Get("document")
	h1 := document.Call("createElement", "h1")
	h1.Set("innerHTML", anotherworld.GetTitle())
	document.Get("body").Call("appendChild", h1)
}
