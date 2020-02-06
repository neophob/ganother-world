package anotherworld

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPartsAllDefined(t *testing.T) {
	result := GetGameParts()
	if len(result) != GAME_PARTS_COUNT {
		t.Errorf("The world will end: %d", result)
	}
}

func TestPartsIntegrationValidateGamePartTypes(t *testing.T) {
	gamePartsIndex := GetGameParts()
	gameData := readFile("../assets/memlist.bin")
	resourceMap, _ := UnmarshallingMemlistBin(gameData)

	if len(resourceMap) == 0 {
		t.Errorf("Unexpected empty resourceMap: %d", resourceMap)
	}

	for partIndex := 0; partIndex < GAME_PARTS_COUNT; partIndex++ {
		var partType string
		var resource MemlistEntry

		gamePart := gamePartsIndex[partIndex]

		resource = resourceMap[gamePart.Palette]
		partType = GetResourceTypeName(int(resource.ResourceType))
		if partType != "RT_PALETTE" {
			t.Errorf("Expected GamePartContent:%d palette (%#02x) to have type RT_PALETTE got: %s",
				partIndex, gamePart.Palette, partType)
		}

		resource = resourceMap[gamePart.Bytecode]
		partType = GetResourceTypeName(int(resource.ResourceType))
		if partType != "RT_BYTECODE" {
			t.Errorf("Expected GamePartContent %d bytecode (%#02x) to have type RT_BYTECODE got: %s",
				partIndex, gamePart.Bytecode, partType)
		}

		resource = resourceMap[gamePart.Cinematic]
		partType = GetResourceTypeName(int(resource.ResourceType))
		if partType != "RT_POLY_CINEMATIC" {
			if resource.UnpackedSize == 0 {
				fmt.Printf("Warning empty GamePartContent[%d].cinematic:%#02x, %-17s\n",
					partIndex, gamePart.Cinematic, partType)
			} else {
				t.Errorf("Expected GamePartContent %d cinematic (%#02x) to have type RT_COMMON_SHAPES got: %s",
					partIndex, gamePart.Cinematic, partType)
			}
		}

		resource = resourceMap[gamePart.Video2]
		partType = GetResourceTypeName(int(resource.ResourceType))
		if partType != "RT_COMMON_SHAPES" {
			if resource.UnpackedSize == 0 {
				fmt.Printf("Warning empty GamePartContent[%d].video2:%#02x, %-17s\n",
					partIndex, gamePart.Video2, partType)
			} else {
				t.Errorf("Expected GamePartContent %d video2 (%#02x) to have type RT_COMMON_SHAPES got: %s",
					partIndex, gamePart.Video2, partType)
			}
		}
	}
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
	}
	return data
}
