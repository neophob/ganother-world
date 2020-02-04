package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"os"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

func main() {
	logger.Info("# GOTHER WORLD vDEV")

	noVideoOutput := flag.Bool("t", false, "Use Text only output (no SDL needed)")
	debug := flag.Bool("d", false, "Enable Debug Mode")
	startPart := flag.Int("p", 1, "Game part to start from (0-9)")
	flag.Parse()

	logger.Info("# KEYBOARD MAPPING:")
	logger.Info("- L: Load State")
	logger.Info("- S: Save State")

	if *debug == false {
		logger.SetLogLevel(logger.LEVEL_INFO)
	}

	logger.Info("- load memlist.bin")
	data := readFile("./assets/memlist.bin")
	bankFilesMap := createBankMap("./assets/")

	var videoDriver anotherworld.Video
	if *noVideoOutput == true {
		videoDriver = anotherworld.Video{Hal: anotherworld.DummyHAL{}}
	} else {
		videoDriver = anotherworld.Video{Hal: buildSDLHAL(), WorkerPage: 0xFE}
	}

	app := initGotherWorld(data, bankFilesMap, videoDriver)

	logger.Info("- setup game")
	app.loadGamePart(anotherworld.GAME_PART_ID_1 + *startPart)

	//start main loop
	for i := 0; app.exitRequested() == false; i++ {
		/*if i%30 == rand.Intn(30) {
			app.loadGamePart(GAME_PART_ID_1+rand.Intn(9))
		}*/
		app.mainLoop(i)

		//game run at approx 25 fps
		time.Sleep(20 * time.Millisecond)
	}

	app.Shutdown()
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("File reading error %v", err)
		os.Exit(1)
	}
	return data
}

func createBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("%sbank%02x", assetPath, i)
		logger.Debug("- load file %s", name)
		entry := readFile(name)
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}
