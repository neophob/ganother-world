package main

import (
	"fmt"
	"syscall/js"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const (
	expectedBankAssets = 13
	defaultStartPart   = anotherworld.GAME_PART_ID_2
	maxStartPartOffset = 9
	minStartPartOffset = 0
)

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
		handleKeyEventWrapper(engine, inputs)
		return nil
	}))
	js.Global().Set("shutdown", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		shutdownJSWrapper(engine)
		return nil
	}))
	js.Global().Set("setLogLevel", js.FuncOf(setLogLevelWrapper))
}

func initGameJSWrapper(engine *Engine, inputs []js.Value) {
	expectedInputs := expectedBankAssets + 1
	if len(inputs) != expectedInputs {
		logger.Error("Expecting %v arguments, 1 for memlist the rest for banks got: %v", expectedInputs, len(inputs))
		return
	}
	if inputs[0].Type() != js.TypeObject {
		logger.Error("Argument 1 for InitGameWrapper(memlist, ...banks) must be a %v not %v", js.TypeObject, inputs[0].Type())
		return
	}

	jsMemlist := inputs[0]
	jsBankFiles := inputs[1:]
	memlist := copyBytesFromJS(jsMemlist)
	bankFilesMap := copyBankMap(jsBankFiles)
	engine.initGame(memlist, bankFilesMap)
}

func startGameFromPartWrapper(engine *Engine, inputs []js.Value) {
	startPartId := defaultStartPart
	if len(inputs) == 1 && inputs[0].Type() != js.TypeUndefined {
		startPartId += parseGamePartOffset(inputs[0])
	}

	logger.Info("Loading game from %v", startPartId)
	engine.app.LoadGamePart(startPartId)

	var renderFrame js.Func

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		engine.mainLoop()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)
}

func handleKeyEventWrapper(engine *Engine, inputs []js.Value) {
	if len(inputs) != 3 {
		logger.Error("Ignoring incomplete key event %v", inputs)
		return
	}
	event := KeyEvent{
		key:     inputs[0].String(),
		keyCode: inputs[1].Int(),
		pressed: inputs[2].String() == "keydown",
	}
	engine.setKeyState(&event)
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
	engine.app.Shutdown()
	engine.shutdownChannel <- true
}

func parseGamePartOffset(gamePartOffset js.Value) int {
	if gamePartOffset.Type() != js.TypeNumber {
		logger.Error("Invalid gamePart offset %v, using default.", gamePartOffset)
		return minStartPartOffset
	}
	parsedOffset := gamePartOffset.Int()
	if parsedOffset >= minStartPartOffset && parsedOffset <= maxStartPartOffset {
		logger.Info("Parsed gamePartOffset %d", parsedOffset)
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
