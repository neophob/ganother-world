package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
)

// renderer is a global variable that needs to implement the Renderer interface
var renderer Renderer = initVideo()

func initVideo() Renderer {
	// start with env VIDEO="SDL" ./main to enable SDL
	videoEnv := os.Getenv("VIDEO")
	if videoEnv == "SDL" {
		return buildSDLRenderer()
	}
	return DummyRenderer{}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("- load memlist.bin")
	data := readFile("./assets/memlist.bin")
	resourceMap, resourceStatistics := unmarshallingMemlistBin(data)
	printResourceStats(resourceStatistics)

	bankFilesMap := createBankMap("./assets/")
	gameParts := getGameParts()
	assets := Assets{
		memList:         resourceMap,
		gameParts:       gameParts,
		bank:            bankFilesMap,
		loadedResources: make(map[int][]uint8),
	}

	log.Println("- create state")
	vmState := createNewState(assets)

	log.Println("- setup game")
	vmState.setupGamePart(GAME_PART_ID_1 + 3)

	videoAssets := vmState.buildVideoAssets()
	renderer.updateGamePart(videoAssets)

	//start main loop
	exit := false
	for i := 0; exit == false; i++ {
		vmState.mainLoop()
		renderer.mainLoop()
		exit = renderer.exitRequested(i)
	}

	renderer.shutdown()
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("File reading error", err)
		os.Exit(1)
	}
	return data
}

func createBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("%sbank%02x", assetPath, i)
		log.Println("- load file", name)
		entry := readFile(name)
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}

func printResourceStats(memlistStatistic MemlistStatistic) {
	log.Println(memlistStatistic)
	fmt.Println("Total # resources:", memlistStatistic.entryCount)
	fmt.Println("Compressed       :", memlistStatistic.compressedEntries)
	fmt.Println("Uncompressed     :", memlistStatistic.entryCount-memlistStatistic.compressedEntries)
	var compressionRatio float64 = 100 / float64(memlistStatistic.entryCount) * float64(memlistStatistic.compressedEntries)
	fmt.Printf("Note: %.0f%% of resources are compressed.\n\n", math.Round(compressionRatio))
	fmt.Printf("Total size (uncompressed) : %d bytes.\n", memlistStatistic.sizeUncompressed)
	fmt.Printf("Total size (compressed)   : %d bytes.\n", memlistStatistic.sizeCompressed)
	var compressionGain float64 = 100 * (1 - float64(memlistStatistic.sizeCompressed)/float64(memlistStatistic.sizeUncompressed))
	fmt.Printf("Note: Overall compression gain is : %.0f%%.\n\n", math.Round(compressionGain))

	sortedKeys := sortedKeys(memlistStatistic.resourceTypeMap)
	for i := 0; i < len(sortedKeys); i++ {
		k := sortedKeys[i]
		resourceName := getResourceTypeName(k)
		if len(resourceName) < 1 {
			resourceName = fmt.Sprintf("RT_UNKOWNN_%d", k)
		}
		fmt.Printf("Total %20s, files: %d\n", resourceName, memlistStatistic.resourceTypeMap[k])
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
