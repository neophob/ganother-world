package main

import (
	"github.com/neophob/ganother-world/logger"
)

func init() {
	logger.SetLogLevel(logger.LEVEL_INFO)
	logger.DisableColors()
	logger.Info("WASM Gother-World initializing...")
}

func main() {
	engine := InitEngine()
	RegisterCallbacks(&engine)
	<-engine.shutdownChannel
}
