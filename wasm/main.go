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
	app := InitGame()
	defer app.Shutdown()

	channel := make(chan bool)
	<-channel
}
