package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

func init() {
	// TODO make logger wasm/browser console friendly
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.Info("WASM " + anotherworld.GetTitle() + " loading...")
}

func main() {
	// TODO get the assets from JS.
	// setup a struct which holds the call back functions (do a test with a print hello world)
	// then register them here (https://github.com/agnivade/shimmer/blob/d08fb873f760d93922e97085a792539e714df4a9/shimmer.go#L39)
	// so that from JS I can call a function with the image
	// then here using the callback we can read the data passed in from JS.

	// TODO Since init game needs the assets it should also be a callback triggered from JS once assets are loaded.
	app := InitGame()
	defer app.Shutdown()

	channel := make(chan bool)
	<-channel
}
