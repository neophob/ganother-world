package main

import (
	"time"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const (
	fixedLoopDelayFor25FPS = 20 * time.Millisecond
)

// TODO avoid globals by creating an Engine struct and passing it to functions that need it.
var app anotherworld.GotherWorld
var shutdownChannel chan bool

func init() {
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.DisableColors()
	logger.Info("WASM Gother-World initializing...")
}

func main() {
	RegisterCallbacks()

	shutdownChannel = make(chan bool)
	<-shutdownChannel
}

// TODO part of engine?
func startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; app.ExitRequested() == false; i++ {
		app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
	}
}
