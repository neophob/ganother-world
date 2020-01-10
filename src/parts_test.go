package main

import (
	"testing"
)

func TestParts(t *testing.T) {
	result := getGameParts()
	if len(result) != 10 {
		t.Errorf("The world will end: %d", result)
	}
}

//TODO make integration test, load real assets, get game parts and then test validate palette entries ARE type RT_PALETTE and code is RT_BYTECODE
