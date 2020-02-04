package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
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
	resourceMap, resourceStatistics := anotherworld.UnmarshallingMemlistBin(memlistData)
	printResourceStats(resourceStatistics)

	gameParts := anotherworld.GetGameParts()
	assets := anotherworld.Assets{
		MemList:         resourceMap,
		GameParts:       gameParts,
		Bank:            bankFilesMap,
		LoadedResources: make(map[int][]uint8),
	}

	logger.Info("- create state")
	vmState := createNewState(assets)
	video := initVideo(noVideoOutput)

	app := GotherWorld{video: video, vm: vmState}
	app.loadGamePart(anotherworld.GAME_PART_FIRST)
	return app
}

func (app *GotherWorld) exitRequested() bool {
	return app.keyPresses&anotherworld.KeyEsc > 0
}

func (app *GotherWorld) mainLoop(i int) {
	app.keyPresses = app.video.EventLoop(i)
	app.vm.mainLoop(app.keyPresses, &app.video)

	if app.keyPresses&anotherworld.KeySave > 0 {
		logger.Info("SAVE STATE")
		app.gameState = GameState{app.vm, app.video}
	}
	if app.gameState.vm.gamePart > 0 && app.keyPresses&anotherworld.KeyLoad > 0 {
		logger.Info("LOAD STATE")
		app.vm.loadGameParts(app.gameState.vm.gamePart)
		app.vm.variables = app.gameState.vm.variables
		app.vm.channelPC = app.gameState.vm.channelPC
		app.vm.nextLoopChannelPC = app.gameState.vm.nextLoopChannelPC
		app.vm.channelPaused = app.gameState.vm.channelPaused
		app.vm.stackCalls = app.gameState.vm.stackCalls
		app.video = app.gameState.video
	}

	if app.vm.loadNextPart > 0 {
		logger.Info("- load next part %d", app.vm.loadNextPart)
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

func printResourceStats(memlistStatistic anotherworld.MemlistStatistic) {
	logger.Debug("Total # resources: %d", memlistStatistic.EntryCount)
	logger.Debug("Compressed       : %d", memlistStatistic.CompressedEntries)
	logger.Debug("Uncompressed     : %d", memlistStatistic.EntryCount-memlistStatistic.CompressedEntries)
	compressionRatio := 100 / float64(memlistStatistic.EntryCount) * float64(memlistStatistic.CompressedEntries)
	logger.Debug("Note: %.0f%% of resources are compressed.", math.Round(compressionRatio))
	logger.Debug("Total size (uncompressed) : %d bytes.", memlistStatistic.SizeUncompressed)
	logger.Debug("Total size (compressed)   : %d bytes.", memlistStatistic.SizeCompressed)
	compressionGain := 100 * (1 - float64(memlistStatistic.SizeCompressed)/float64(memlistStatistic.SizeUncompressed))
	logger.Debug("Note: Overall compression gain is : %.0f%%.", math.Round(compressionGain))

	sortedKeys := sortedKeys(memlistStatistic.ResourceTypeMap)
	for i := 0; i < len(sortedKeys); i++ {
		k := sortedKeys[i]
		resourceName := anotherworld.GetResourceTypeName(k)
		if len(resourceName) < 1 {
			resourceName = fmt.Sprintf("RT_UNKOWNN_%d", k)
		}
		logger.Debug("Total %20s, files: %d", resourceName, memlistStatistic.ResourceTypeMap[k])
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
