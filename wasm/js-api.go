package main

import (
	"fmt"
	"syscall/js"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const (
	expectedBankAssets = 13
	defaultStartPart   = anotherworld.GAME_PART_ID_1
	maxStartPartOffset = 9
	minStartPartOffset = 0
)

// TODO should be part of hal or engine...
var keyMap = make(map[uint32]bool)

func RegisterCallbacks(engine *Engine) {
	js.Global().Set("initGameFromURI", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		initGameJSWrapper(engine, inputs)
		return nil
	}))
	js.Global().Set("startGameFromPart", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		startGameFromPartWrapper(engine, inputs)
		return nil
	}))
	js.Global().Set("handleKeyEvent", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		startGameFromPartWrapper(engine, inputs)
		return nil
	}))
	js.Global().Set("shutdown", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		shutdownJSWrapper(engine)
		return nil
	}))
	js.Global().Set("setLogLevel", js.FuncOf(setLogLevelWrapper))
}

func initGameJSWrapper(engine *Engine, inputs []js.Value) interface{} {
	expectedInputs := expectedBankAssets + 1
	if len(inputs) != expectedInputs {
		logger.Error("Expecting %v arguments, 1 for memlist the rest for banks got: %v", expectedInputs, len(inputs))
		return nil
	}
	if inputs[0].Type() != js.TypeObject {
		logger.Error("Argument 1 for InitGameWrapper(memlist, ...banks) must be a %v not %v", js.TypeObject, inputs[0].Type())
		return nil
	}

	jsMemlist := inputs[0]
	jsBankFiles := inputs[1:]

	memlist := copyBytesFromJS(jsMemlist)
	bankFilesMap := copyBankMap(jsBankFiles)

	// TODO init should return app and set it on engine
	app = InitGame(memlist, bankFilesMap)

	return nil
}

func startGameFromPartWrapper(engine *Engine, inputs []js.Value) interface{} {
	startPartId := defaultStartPart
	if len(inputs) == 1 && inputs[0].Type() != js.TypeUndefined {
		startPartId += parseGamePartOffset(inputs[0])
	}

	nonDefaultPartLoadingIsNeeded := startPartId != defaultStartPart
	if nonDefaultPartLoadingIsNeeded {
		logger.Info("Loading game from %v", startPartId)
		// TODO call this from engine.app
		app.LoadGamePart(startPartId)
	}

	go startMainLoop()

	return nil
}

func handleKeyEventWrapper(engine *Engine, inputs []js.Value) interface{} {
	if len(inputs) != 3 {
		logger.Error("Ignoring incomplete key event", inputs)
	}
	event := KeyEvent{
		key:     inputs[0].String(),
		keyCode: inputs[1].Int(),
		pressed: inputs[2].String() == "keydown",
	}
	// TODO this should be bassed to engine
	setKeyState(&event)
	logger.Info("Updated KeyMap %v", keyMap)
	return nil
}

func setLogLevelWrapper(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) != 1 || inputs[0].Type() != js.TypeNumber {
		logger.Error("Invalid log level requested %v, number required", inputs[0])
	}
	logger.SetLogLevel(inputs[0].Int())
	return nil
}

func shutdownJSWrapper(engine *Engine) {
	logger.Info("Shutting down")
	// TODO use engine.app and engine.shutDown...
	app.Shutdown()
	shutdownChannel <- true
}

// TODO move this to a state shared with the hal
// pass it to Hal... some how figure out how to access the hal? app.video.Hal
func setKeyState(event *KeyEvent) {
	logger.Info("Key Event %v", event)
	/*
		TODO
		Where would the state be stored and how does the HAL access it?
		See hal-sdl for multiple key holds handling.
	*/
	switch event.key {
	case "Escape":
		keyMap[anotherworld.KeyEsc] = event.pressed
		return
	case "ArrowLeft":
		keyMap[anotherworld.KeyLeft] = event.pressed
		return
	case "ArrowRight":
		keyMap[anotherworld.KeyRight] = event.pressed
		return
	case "ArrowUp":
		keyMap[anotherworld.KeyUp] = event.pressed
		return
	case "ArrowDown":
		keyMap[anotherworld.KeyDown] = event.pressed
		return
	case " ":
		keyMap[anotherworld.KeyFire] = event.pressed
		return
	case "p", "P":
		keyMap[anotherworld.KeyPause] = event.pressed
		return
	case "s", "S":
		keyMap[anotherworld.KeySave] = event.pressed
		return
	case "l", "L":
		keyMap[anotherworld.KeyLoad] = event.pressed
		return
	}
}

func parseGamePartOffset(gamePartOffset js.Value) int {
	if gamePartOffset.Type() != js.TypeNumber {
		logger.Error("Invalid gamePart offset %v, using default.", gamePartOffset)
		return minStartPartOffset
	}
	parsedOffset := gamePartOffset.Int()
	if parsedOffset >= minStartPartOffset && parsedOffset <= maxStartPartOffset {
		logger.Info("Parsed gamePartOffset %v", gamePartOffset)
		return parsedOffset
	}
	logger.Error("Out of range gamePart offset %v, using default.", parsedOffset)
	return minStartPartOffset
}

func copyBankMap(bankInputs []js.Value) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("bank%02x", i)
		logger.Debug("- coping %s from JS", name)
		entry := copyBytesFromJS(bankInputs[i-1])
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}

func copyBytesFromJS(input js.Value) []byte {
	data := make([]uint8, input.Get("byteLength").Int())
	js.CopyBytesToGo(data, input)
	return data
}
