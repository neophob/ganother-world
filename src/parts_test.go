package main

import (
	"testing"
)

func TestPartsAllDefined(t *testing.T) {
	result := getGameParts()
	if len(result) != GAME_PARTS_COUNT {
		t.Errorf("The world will end: %d", result)
	}
}

//TODO make integration test, load real assets, get game parts and then test validate palette entries ARE type RT_PALETTE and code is RT_BYTECODE
