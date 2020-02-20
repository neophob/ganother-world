package main

import (
	"github.com/neophob/ganother-world/anotherworld"
)

func InitGame(memlist []byte, bankFilesMap map[int][]byte) anotherworld.GotherWorld {
	videoDriver := anotherworld.Video{Hal: WASMHAL{}}
	app := anotherworld.InitGotherWorld(memlist, bankFilesMap, videoDriver)
	return app
}
