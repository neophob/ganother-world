package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const (
	expectedBankAssets     = 13
	defaultStartPart       = anotherworld.GAME_PART_ID_1
	maxStartPartOffset     = 9
	minStartPartOffset     = 0
	fixedLoopDelayFor25FPS = 20 * time.Millisecond
)

var app anotherworld.GotherWorld
var channel chan bool

func init() {
	// TODO ideally parse debug level from get parameter: ?logLevel=debug
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.DisableColors()
	logger.Info("WASM Gother-World initializing...")
}

func main() {
	RegisterCallbacks()

	channel = make(chan bool)
	<-channel
}

func RegisterCallbacks() {
	js.Global().Set("initGameFromURI", js.FuncOf(InitGameJSWrapper))
	js.Global().Set("startGameFromPart", js.FuncOf(startGameFromPart))
	js.Global().Set("shutdown", js.FuncOf(ShutdownJSWrapper))
	js.Global().Set("setLogLevel", js.FuncOf(SetLogLevelWrapper))
}

func InitGameJSWrapper(this js.Value, inputs []js.Value) interface{} {
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

	app = InitGame(memlist, bankFilesMap)

	return nil
}

func startGameFromPart(this js.Value, inputs []js.Value) interface{} {
	startPartId := defaultStartPart
	if len(inputs) == 1 && inputs[0].Type() != js.TypeUndefined {
		startPartId += parseGamePartOffset(inputs[0])
	}

	nonDefaultPartLoadingIsNeeded := startPartId != defaultStartPart
	if nonDefaultPartLoadingIsNeeded {
		logger.Info("Loading game from %v", startPartId)
		app.LoadGamePart(startPartId)
	}

	go startMainLoop()

	return nil
}

func startMainLoop() {
	// TODO sync this with request animation frame
	for i := 0; app.ExitRequested() == false; i++ {
		app.MainLoop(i)
		time.Sleep(fixedLoopDelayFor25FPS)
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

func SetLogLevelWrapper(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) != 1 || inputs[0].Type() != js.TypeNumber {
		logger.Error("Invalid log level requested %v, number required", inputs[0])
	}
	logger.SetLogLevel(inputs[0].Int())
	return nil
}

func ShutdownJSWrapper(this js.Value, inputs []js.Value) interface{} {
	logger.Info("Shutting down")
	app.Shutdown()
	channel <- true
	return nil
}
