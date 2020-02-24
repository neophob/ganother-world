package main

import (
	"time"

	"github.com/neophob/ganother-world/anotherworld"
)

const (
	fixedLoopDelayFor25FPS = 20 * time.Millisecond
)

func InitEngine() Engine {
	return Engine{
		shutdownChannel: make(chan bool),
	}
}

type Engine struct {
	app             anotherworld.GotherWorld
	shutdownChannel chan bool
}

func (engine *Engine) startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; engine.app.ExitRequested() == false; i++ {
		engine.app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
	}
}
