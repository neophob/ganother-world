package main

import (
	"flag"
	"time"

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
	data := anotherworld.ReadFile("./assets/memlist.bin")
	bankFilesMap := anotherworld.CreateBankMap("./assets/")

	var videoDriver anotherworld.Video
	if *noVideoOutput == true {
		videoDriver = anotherworld.Video{Hal: anotherworld.DummyHAL{}}
	} else {
		videoDriver = anotherworld.Video{Hal: buildSDLHAL(), WorkerPage: 0xFE}
	}

	app := anotherworld.InitGotherWorld(data, bankFilesMap, videoDriver)

	logger.Info("- setup game")
	app.LoadGamePart(anotherworld.GAME_PART_ID_1 + *startPart)

	//start main loop
	for i := 0; app.ExitRequested() == false; i++ {
		/*if i%30 == rand.Intn(30) {
			app.LoadGamePart(GAME_PART_ID_1+rand.Intn(9))
		}*/
		app.MainLoop(i)

		//game run at approx 25 fps
		time.Sleep(20 * time.Millisecond)
	}

	app.Shutdown()
}
