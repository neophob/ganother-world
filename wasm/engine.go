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
	videoDriver     anotherworld.Video
	keyMap          map[uint32]bool
	shutdownChannel chan bool
}

type KeyEvent struct {
	key     string
	keyCode int
	pressed bool
}

func InitEngine() Engine {
	return Engine{
		shutdownChannel: make(chan bool),
		keyMap:          make(map[uint32]bool),
		videoDriver:     anotherworld.Video{Hal: buildWASMHAL()},
	}
}

func (engine *Engine) initGame(memlist []byte, bankFilesMap map[int][]byte) {
	engine.app = anotherworld.InitGotherWorld(memlist, bankFilesMap, engine.videoDriver)
}

func (engine *Engine) startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; engine.app.ExitRequested() == false; i++ {
		engine.app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
	}
}

// pass it to Hal... some how figure out how to access the hal? app.video.Hal
func (engine *Engine) setKeyState(event *KeyEvent) {
	logger.Info("Key Event %v", event)
	/*
		TODO
		Where would the state be stored and how does the HAL access it?
		See hal-sdl for multiple key holds handling.
	*/
	switch event.key {
	case "Escape":
		engine.keyMap[anotherworld.KeyEsc] = event.pressed
		return
	case "ArrowLeft":
		engine.keyMap[anotherworld.KeyLeft] = event.pressed
		return
	case "ArrowRight":
		engine.keyMap[anotherworld.KeyRight] = event.pressed
		return
	case "ArrowUp":
		engine.keyMap[anotherworld.KeyUp] = event.pressed
		return
	case "ArrowDown":
		engine.keyMap[anotherworld.KeyDown] = event.pressed
		return
	case " ":
		engine.keyMap[anotherworld.KeyFire] = event.pressed
		return
	case "p", "P":
		engine.keyMap[anotherworld.KeyPause] = event.pressed
		return
	case "s", "S":
		engine.keyMap[anotherworld.KeySave] = event.pressed
		return
	case "l", "L":
		engine.keyMap[anotherworld.KeyLoad] = event.pressed
		return
	}
}
