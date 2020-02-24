package main

import (
	"time"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const (
	fixedLoopDelayFor25FPS = 20 * time.Millisecond
)

type Engine struct {
	app             anotherworld.GotherWorld
	shutdownChannel chan bool
}

func init() {
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.DisableColors()
	logger.Info("WASM Gother-World initializing...")
}

func main() {
	engine := Engine{
		shutdownChannel: make(chan bool),
	}
	RegisterCallbacks(&engine)
	<-engine.shutdownChannel
}

func (engine *Engine) startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; engine.app.ExitRequested() == false; i++ {
		engine.app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
	}
}
