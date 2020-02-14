package main

import (
	"github.com/neophob/ganother-world/anotherworld"
)

func InitGame(memlist []byte) anotherworld.GotherWorld {
	data := memlist
	bankFilesMap := createBankMap("TODO")
	videoDriver := anotherworld.Video{Hal: anotherworld.DummyHAL{}}

	app := anotherworld.InitGotherWorld(data, bankFilesMap, videoDriver)
	return app
}

func createBankMap(assetsURI string) map[int][]byte {
	// TODO update fetching to use HTTP:
	// https://github.com/golang/go/wiki/WebAssembly#configuring-fetch-options-while-using-nethttp
	bankFilesMap := make(map[int][]byte)
	return bankFilesMap
}
