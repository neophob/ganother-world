package anotherworld

import (
	"fmt"
	"math"
	"sort"

	"github.com/neophob/ganother-world/logger"
)

//GotherWorld is the root object, it holds the whole world
type GotherWorld struct {
	video      Video
	Vm         VMState
	gameState  GameState
	keyPresses uint32
}

//GameState used to save and load a game state
type GameState struct {
	vm    VMState
	video Video
}

func InitGotherWorld(memlistData []byte, bankFilesMap map[int][]byte, videoDriver Video) GotherWorld {
	resourceMap, resourceStatistics := UnmarshallingMemlistBin(memlistData)
	printResourceStats(resourceStatistics)

	gameParts := GetGameParts()
	assets := StaticGameAssets{
		MemList:         resourceMap,
		GameParts:       gameParts,
		Bank:            bankFilesMap,
		LoadedResources: make(map[int][]uint8),
	}

	logger.Info("- create state")
	vmState := CreateNewState(assets)
	app := GotherWorld{video: videoDriver, Vm: vmState}
	app.LoadGamePart(GAME_PART_FIRST)
	return app
}

func (app *GotherWorld) ExitRequested() bool {
	return app.keyPresses&KeyEsc > 0
}

func (app *GotherWorld) MainLoop(i int) {
	app.keyPresses = app.video.EventLoop(i)
	app.Vm.MainLoop(app.keyPresses, &app.video)

	if app.keyPresses&KeySave > 0 {
		logger.Info("SAVE STATE")
		app.gameState = GameState{app.Vm, app.video}
	}
	if app.gameState.vm.GamePart > 0 && app.keyPresses&KeyLoad > 0 {
		logger.Info("LOAD STATE")
		app.Vm.LoadGameParts(app.gameState.vm.GamePart)
		app.Vm.Variables = app.gameState.vm.Variables
		app.Vm.ChannelPC = app.gameState.vm.ChannelPC
		app.Vm.NextLoopChannelPC = app.gameState.vm.NextLoopChannelPC
		app.Vm.ChannelPaused = app.gameState.vm.ChannelPaused
		app.Vm.StackCalls = app.gameState.vm.StackCalls
		app.video = app.gameState.video
	}

	if app.Vm.LoadNextPart > 0 {
		logger.Info("- load next part %d", app.Vm.LoadNextPart)
		app.LoadGamePart(app.Vm.LoadNextPart)
	}
}

func (app *GotherWorld) LoadGamePart(partID int) {
	app.Vm.SetupGamePart(partID)
	//TODO rename videoAssets to game part assets
	//TODO rename video struct to game??
	//TODO add audio stuff
	videoAssets := app.Vm.BuildVideoAssets()
	app.video.UpdateGamePart(videoAssets)
}

func (app *GotherWorld) Shutdown() {
	app.video.Shutdown()
}

func printResourceStats(memlistStatistic MemlistStatistic) {
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
		resourceName := GetResourceTypeName(k)
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
