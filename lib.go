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
	video      anotherworld.Video
	vm         anotherworld.VMState
	gameState  GameState
	keyPresses uint32
}

//GameState used to save and load a game state
type GameState struct {
	vm    anotherworld.VMState
	video anotherworld.Video
}

func initGotherWorld(memlistData []byte, bankFilesMap map[int][]byte, videoDriver anotherworld.Video) GotherWorld {
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
	vmState := anotherworld.CreateNewState(assets)
	app := GotherWorld{video: videoDriver, vm: vmState}
	app.loadGamePart(anotherworld.GAME_PART_FIRST)
	return app
}

func (app *GotherWorld) exitRequested() bool {
	return app.keyPresses&anotherworld.KeyEsc > 0
}

func (app *GotherWorld) MainLoop(i int) {
	app.keyPresses = app.video.EventLoop(i)
	app.vm.MainLoop(app.keyPresses, &app.video)

	if app.keyPresses&anotherworld.KeySave > 0 {
		logger.Info("SAVE STATE")
		app.gameState = GameState{app.vm, app.video}
	}
	if app.gameState.vm.GamePart > 0 && app.keyPresses&anotherworld.KeyLoad > 0 {
		logger.Info("LOAD STATE")
		app.vm.LoadGameParts(app.gameState.vm.GamePart)
		app.vm.Variables = app.gameState.vm.Variables
		app.vm.ChannelPC = app.gameState.vm.ChannelPC
		app.vm.NextLoopChannelPC = app.gameState.vm.NextLoopChannelPC
		app.vm.ChannelPaused = app.gameState.vm.ChannelPaused
		app.vm.StackCalls = app.gameState.vm.StackCalls
		app.video = app.gameState.video
	}

	if app.vm.LoadNextPart > 0 {
		logger.Info("- load next part %d", app.vm.LoadNextPart)
		app.loadGamePart(app.vm.LoadNextPart)
	}
}

func (app *GotherWorld) loadGamePart(partID int) {
	app.vm.SetupGamePart(partID)
	//TODO rename videoAssets to game part assets
	//TODO rename video struct to game??
	//TODO add audio stuff
	videoAssets := app.vm.BuildVideoAssets()
	app.video.UpdateGamePart(videoAssets)
}

func (app *GotherWorld) Shutdown() {
	app.video.Shutdown()
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
