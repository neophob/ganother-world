package main

import (
	"syscall/js"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

var app anotherworld.GotherWorld
var channel chan bool

func init() {
	// TODO make logger wasm/browser console friendly
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.Info("WASM " + anotherworld.GetTitle() + " initializing...")
}

func main() {
	// TODO get the assets from JS.
	// setup a struct which holds the call back functions (do a test with a print hello world)
	// then register them here (https://github.com/agnivade/shimmer/blob/d08fb873f760d93922e97085a792539e714df4a9/shimmer.go#L39)
	// so that from JS I can call a function with the image
	// then here using the callback we can read the data passed in from JS.
	RegisterCallbacks()

	channel = make(chan bool)
	<-channel
}

func RegisterCallbacks() {
	js.Global().Set("initGameFromURI", js.FuncOf(InitGameJSWrapper))
	js.Global().Set("shutdown", js.FuncOf(ShutdownJSWrapper))
}

func InitGameJSWrapper(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) < 1 || inputs[0].Type() != js.TypeString {
		logger.Error("Argument one for InitGameWrapper(assetsURI) must be a string")
		return nil
	}
	assetsURI := inputs[0].String()
	logger.Info("Initializing game from %s", assetsURI)
	app = InitGame(assetsURI)
	return nil
}

func ShutdownJSWrapper(this js.Value, inputs []js.Value) interface{} {
	logger.Info("Shutting down")
	app.Shutdown()
	channel <- true
	return nil
}
