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

	for partNumber := 1; partNumber < GAME_PARTS_COUNT; partNumber++ {
		gamePart := gamePartsIndex[partNumber-1]
		var partType string

		part1Palette := resourceMap[gamePart.palette]
		partType = getResourceTypeName(int(part1Palette.resourceType))
		if partType != "RT_PALETTE" {
			t.Errorf("Expected GamePartContent:%d palette (%#02x) to have type RT_PALETTE got: %s",
				partNumber, gamePart.palette, partType)
		}

		partByteCode := resourceMap[gamePart.bytecode]
		partType = getResourceTypeName(int(partByteCode.resourceType))
		if partType != "RT_BYTECODE" {
			t.Errorf("Expected GamePartContent %d bytecode (%#02x) to have type RT_BYTECODE got: %s",
				partNumber, gamePart.bytecode, partType)
		}

		partCinematic := resourceMap[gamePart.cinematic]
		partType = getResourceTypeName(int(partCinematic.resourceType))
		if partType != "RT_POLY_CINEMATIC" {
			t.Errorf("Expected GamePartContent %d cinematic (%#02x) to have type RT_POLY_CINEMATIC got: %s",
				partNumber, gamePart.cinematic, partType)
		}

		partVideo2 := resourceMap[gamePart.video2]
		partType = getResourceTypeName(int(partVideo2.resourceType))
		if partType != "RT_COMMON_SHAPES" {
			if partVideo2.unpackedSize == 0 {
				fmt.Printf("Warning empty GamePartContent[%d].video2:%#02x, %-17s video2 %d\n",
					partNumber, gamePart.video2, partType, partVideo2)
			} else {
				t.Errorf("Expected GamePartContent %d video2 (%#02x) to have type RT_COMMON_SHAPES got: %s",
					partNumber, gamePart.video2, partType)
			}
		}
	}
}
