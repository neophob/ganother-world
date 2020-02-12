package main

import (
	"fmt"
	"net/http"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

func InitGame(assetsURI string) anotherworld.GotherWorld {
	logger.Info("Loading assets over http from %s", assetsURI)
	data := fetchAssets(assetsURI + "memlist.bin")
	bankFilesMap := createBankMap(assetsURI)
	videoDriver := anotherworld.Video{Hal: anotherworld.DummyHAL{}}

	app := anotherworld.InitGotherWorld(data, bankFilesMap, videoDriver)
	return app
}

func fetchAssets(memlistURI string) []byte {
	// TODO update fetching to use HTTP:
	// https://github.com/golang/go/wiki/WebAssembly#configuring-fetch-options-while-using-nethttp
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	req.Header.Add("js.fetch:mode", "cors")
	if err != nil {
		// TODO handle error
		fmt.Println(err)
		return make([]byte, 0)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// TODO handle error
		fmt.Println(err)
		return make([]byte, 0)
	}
	defer resp.Body.Close()

	// TODO do something with resp...

	return make([]byte, 0)
}

func createBankMap(assetsURI string) map[int][]byte {
	// TODO update fetching to use HTTP:
	// https://github.com/golang/go/wiki/WebAssembly#configuring-fetch-options-while-using-nethttp
	bankFilesMap := make(map[int][]byte)
	return bankFilesMap
}
