package main

import (
	"fmt"
	"math"
	"sort"
)

//GotherWorld is the root object, it holds the whole world
type GotherWorld struct {
	video      Video
	vm         VMState
	gameState  GameState
	keyPresses uint32
}

//GameState used to save and load a game state
type GameState struct {
	vm    VMState
	video Video
}

func initGotherWorld(memlistData []byte, bankFilesMap map[int][]byte, noVideoOutput bool) GotherWorld {
	resourceMap, resourceStatistics := unmarshallingMemlistBin(memlistData)
	printResourceStats(resourceStatistics)

	gameParts := getGameParts()
	assets := Assets{
		memList:         resourceMap,
		gameParts:       gameParts,
		bank:            bankFilesMap,
		loadedResources: make(map[int][]uint8),
	}

	Info("- create state")
	vmState := createNewState(assets)
	video := initVideo(noVideoOutput)

	app := GotherWorld{video: video, vm: vmState}
	app.loadGamePart(GAME_PART_FIRST)
	return app
}

func (app *GotherWorld) exitRequested() bool {
	return app.keyPresses&KeyEsc > 0
}

func (app *GotherWorld) mainLoop(i int) {
	app.keyPresses = app.video.eventLoop(i)
	app.vm.mainLoop(app.keyPresses, &app.video)

	if app.keyPresses&KeySave > 0 {
		Info("SAVE STATE")
		app.gameState = GameState{app.vm, app.video}
	}
	if app.gameState.vm.gamePart > 0 && app.keyPresses&KeyLoad > 0 {
		Info("LOAD STATE")
		app.vm.loadGameParts(app.gameState.vm.gamePart)
		app.vm.variables = app.gameState.vm.variables
		app.vm.channelPC = app.gameState.vm.channelPC
		app.vm.nextLoopChannelPC = app.gameState.vm.nextLoopChannelPC
		app.vm.channelPaused = app.gameState.vm.channelPaused
		app.vm.stackCalls = app.gameState.vm.stackCalls
		app.video = app.gameState.video
	}

	if app.vm.loadNextPart > 0 {
		Info("- load next part %d", app.vm.loadNextPart)
		app.loadGamePart(app.vm.loadNextPart)
	}
}

func (app *GotherWorld) loadGamePart(partID int) {
	app.vm.setupGamePart(partID)
	//TODO rename videoAssets to game part assets
	//TODO rename video struct to game??
	//TODO add audio stuff
	videoAssets := app.vm.buildVideoAssets()
	app.video.updateGamePart(videoAssets)
}

func (app *GotherWorld) shutdown() {
	app.video.shutdown()
}

func printResourceStats(memlistStatistic MemlistStatistic) {
	Debug("Total # resources: %d", memlistStatistic.entryCount)
	Debug("Compressed       : %d", memlistStatistic.compressedEntries)
	Debug("Uncompressed     : %d", memlistStatistic.entryCount-memlistStatistic.compressedEntries)
	compressionRatio := 100 / float64(memlistStatistic.entryCount) * float64(memlistStatistic.compressedEntries)
	Debug("Note: %.0f%% of resources are compressed.", math.Round(compressionRatio))
	Debug("Total size (uncompressed) : %d bytes.", memlistStatistic.sizeUncompressed)
	Debug("Total size (compressed)   : %d bytes.", memlistStatistic.sizeCompressed)
	compressionGain := 100 * (1 - float64(memlistStatistic.sizeCompressed)/float64(memlistStatistic.sizeUncompressed))
	Debug("Note: Overall compression gain is : %.0f%%.", math.Round(compressionGain))

	sortedKeys := sortedKeys(memlistStatistic.resourceTypeMap)
	for i := 0; i < len(sortedKeys); i++ {
		k := sortedKeys[i]
		resourceName := getResourceTypeName(k)
		if len(resourceName) < 1 {
			resourceName = fmt.Sprintf("RT_UNKOWNN_%d", k)
		}
		Debug("Total %20s, files: %d", resourceName, memlistStatistic.resourceTypeMap[k])
	}
}

func sortedKeys(m map[int]int) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}
