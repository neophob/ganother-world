package main

import (
	"fmt"
	"testing"
)

func TestPartsAllDefined(t *testing.T) {
	result := getGameParts()
	if len(result) != GAME_PARTS_COUNT {
		t.Errorf("The world will end: %d", result)
	}
}

func TestPartsIntegrationValidateGamePartTypes(t *testing.T) {
	gamePartsIndex := getGameParts()
	gameData := readFile("../assets/memlist.bin")
	resourceMap, _ := unmarshallingMemlistBin(gameData)

	if len(resourceMap) == 0 {
		t.Errorf("Unexpected empty resourceMap: %d", resourceMap)
	}

	for partIndex := 0; partIndex < GAME_PARTS_COUNT; partIndex++ {
		var partType string
		var resource MemlistEntry

		gamePart := gamePartsIndex[partIndex]

		resource = resourceMap[gamePart.palette]
		partType = getResourceTypeName(int(resource.resourceType))
		if partType != "RT_PALETTE" {
			t.Errorf("Expected GamePartContent:%d palette (%#02x) to have type RT_PALETTE got: %s",
				partIndex, gamePart.palette, partType)
		}

		resource = resourceMap[gamePart.bytecode]
		partType = getResourceTypeName(int(resource.resourceType))
		if partType != "RT_BYTECODE" {
			t.Errorf("Expected GamePartContent %d bytecode (%#02x) to have type RT_BYTECODE got: %s",
				partIndex, gamePart.bytecode, partType)
		}

		resource = resourceMap[gamePart.cinematic]
		partType = getResourceTypeName(int(resource.resourceType))
		if partType != "RT_POLY_CINEMATIC" {
			if resource.unpackedSize == 0 {
				fmt.Printf("Warning empty GamePartContent[%d].cinematic:%#02x, %-17s\n",
					partIndex, gamePart.cinematic, partType)
			} else {
				t.Errorf("Expected GamePartContent %d cinematic (%#02x) to have type RT_COMMON_SHAPES got: %s",
					partIndex, gamePart.cinematic, partType)
			}
		}

		resource = resourceMap[gamePart.video2]
		partType = getResourceTypeName(int(resource.resourceType))
		if partType != "RT_COMMON_SHAPES" {
			if resource.unpackedSize == 0 {
				fmt.Printf("Warning empty GamePartContent[%d].video2:%#02x, %-17s\n",
					partIndex, gamePart.video2, partType)
			} else {
				t.Errorf("Expected GamePartContent %d video2 (%#02x) to have type RT_COMMON_SHAPES got: %s",
					partIndex, gamePart.video2, partType)
			}
		}
	}
}
