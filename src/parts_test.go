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

// TODO get game parts and then test validate palette entries ARE type RT_PALETTE and code is RT_BYTECODE
func TestPartsIntegrationValidateGamePartsPalettes(t *testing.T) {
	gamePartsIndex := getGameParts()
	gameData := readFile("../assets/memlist.bin")
	resourceMap, _ := unmarshallingMemlistBin(gameData)

	if len(resourceMap) == 0 {
		t.Errorf("Unexpected empty resourceMap: %d", resourceMap)
	}

	// TODO loop all parts
	partNumber := 1
	gamePart1 := gamePartsIndex[partNumber-1]
	part1Palette := resourceMap[gamePart1.palette]
	part1PaletteType := getResourceTypeName(int(part1Palette.resourceType))

	fmt.Printf("Test:%#02x, %-17s palette %d\n", gamePart1.palette,
		getResourceTypeName(int(part1Palette.resourceType)), part1Palette)

	if part1PaletteType != "RT_PALETTE" {
		t.Errorf("Expected GamePartContent %d palette (%#02x) to have type RT_PALETTE got: %s", partNumber, gamePart1.palette, part1PaletteType)
	}
}
