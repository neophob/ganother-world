package main

import (
	"github.com/neophob/ganother-world/logger"
)

func init() {
	logger.DisableColors()
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.Info("WASM Gother-World initializing...")
}

func main() {
	engine := buildEngine()
	RegisterCallbacks(&engine)
	logger.Info("Registered Callbacks")
	<-engine.shutdownChannel
}
