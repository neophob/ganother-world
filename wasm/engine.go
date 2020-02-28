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
	hal             anotherworld.HAL
	keyMap          map[uint32]bool
	shutdownChannel chan bool
}

type KeyEvent struct {
	key     string
	keyCode int
	pressed bool
}

func buildEngine() Engine {
	return Engine{
		shutdownChannel: make(chan bool),
		keyMap:          make(map[uint32]bool),
		hal:             buildWASMHAL(),
	}
}

func (engine *Engine) initGame(memlist []byte, bankFilesMap map[int][]byte) {
	engine.app = anotherworld.InitGotherWorld(memlist, bankFilesMap, anotherworld.Video{Hal: engine.hal})
}

func (engine *Engine) startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; engine.app.ExitRequested() == false; i++ {
		engine.app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
	}
}

func (engine *Engine) setKeyState(event *KeyEvent) {
	anotherKey := mapKeyToAnotherworld(event.key)
	engine.keyMap[anotherKey] = event.pressed
	logger.Debug("Updated KeyMap with %v %v", event, engine.keyMap)
	engine.hal.(*WASMHAL).updateKeyStateFrom(&engine.keyMap)
}

func mapKeyToAnotherworld(key string) uint32 {
	switch key {
	case "Escape":
		return anotherworld.KeyEsc
	case "ArrowLeft":
		return anotherworld.KeyLeft
	case "ArrowRight":
		return anotherworld.KeyRight
	case "ArrowUp":
		return anotherworld.KeyUp
	case "ArrowDown":
		return anotherworld.KeyDown
	case " ":
		return anotherworld.KeyFire
	case "p", "P":
		return anotherworld.KeyPause
	case "s", "S":
		return anotherworld.KeySave
	case "l", "L":
		return anotherworld.KeyLoad
	}
	return anotherworld.KeyNone
}
