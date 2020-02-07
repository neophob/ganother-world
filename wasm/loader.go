package main

import (
	"github.com/neophob/ganother-world/anotherworld"
)

func InitGame() anotherworld.GotherWorld {
	data := fetchAssets("./assets/memlist.bin")
	bankFilesMap := createBankMap("./assets/")
	videoDriver := anotherworld.Video{Hal: anotherworld.DummyHAL{}}

	app := anotherworld.InitGotherWorld(data, bankFilesMap, videoDriver)
	return app
}

func fetchAssets(filename string) []byte {
	// TODO implement loading through CopyBytesToGo
	return make([]byte, 0)
}

func createBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	// TODO implement loading through CopyBytesToGo
	return bankFilesMap
}
