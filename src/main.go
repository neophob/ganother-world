package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"os"
)

func main() {
	Info("# GOTHER WORLD vDEV")

	noVideoOutput := flag.Bool("t", false, "Use Text only output (no SDL needed)")
	debug := flag.Bool("d", false, "Enable Debug Mode")
	startPart := flag.Int("p", 1, "Game part to start from (0-9)")
	flag.Parse()

	Info("# KEYBOARD MAPPING:")
	Info("- L: Load State")
	Info("- S: Save State")

	if *debug == false {
		SetLogLevel(LEVEL_INFO)
	}

	Info("- load memlist.bin")
	data := readFile("./assets/memlist.bin")
	bankFilesMap := createBankMap("./assets/")

	app := initGotherWorld(data, bankFilesMap, *noVideoOutput)

	Info("- setup game")
	app.loadGamePart(GAME_PART_ID_1 + *startPart)

	//start main loop
	for i := 0; app.exitRequested() == false; i++ {
		/*if i%30 == rand.Intn(30) {
			loadGamePart(&vmState, GAME_PART_ID_1+rand.Intn(9))
		}*/
		app.mainLoop(i)

		//game run at approx 25 fps
		time.Sleep(20 * time.Millisecond)
	}

	app.video.shutdown()
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		Error("File reading error %v", err)
		os.Exit(1)
	}
	return data
}

func createBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("%sbank%02x", assetPath, i)
		Debug("- load file %s", name)
		entry := readFile(name)
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}
