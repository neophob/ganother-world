package main

import (
	"fmt"
	"syscall/js"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

const expectedBankAssets = 13

var app anotherworld.GotherWorld
var channel chan bool

func init() {
	// TODO make logger wasm/browser console friendly
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.Info("WASM " + anotherworld.GetTitle() + " initializing...")
}

func main() {
	// TODO get the assets from JS.
	// setup a struct which holds the call back functions (do a test with a print hello world)
	// then register them here (https://github.com/agnivade/shimmer/blob/d08fb873f760d93922e97085a792539e714df4a9/shimmer.go#L39)
	// so that from JS I can call a function with the image
	// then here using the callback we can read the data passed in from JS.
	RegisterCallbacks()

	channel = make(chan bool)
	<-channel
}

func RegisterCallbacks() {
	js.Global().Set("initGameFromURI", js.FuncOf(InitGameJSWrapper))
	js.Global().Set("shutdown", js.FuncOf(ShutdownJSWrapper))
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

func copyBankMap(bankInputs []js.Value) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("bank%02x", i)
		logger.Info("- coping %s from JS", name)
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

func ShutdownJSWrapper(this js.Value, inputs []js.Value) interface{} {
	logger.Info("Shutting down")
	app.Shutdown()
	channel <- true
	return nil
}
