package main

import (
	"testing"
)

func TestDataFetcher(t *testing.T) {
	asset := []uint8{1, 2, 3}
	vdf := VideoDataFetcher{asset: &asset}

	res := vdf.fetchByte()
	if res != 1 {
		t.Errorf("The world will end, res: %d", res)
	}
	if vdf.readOffset != 1 {
		t.Errorf("The world will end, offset: %d", vdf.readOffset)
	}

	resW := vdf.fetchWord()
	if resW != 0x203 {
		t.Errorf("The world will end, res: %d", resW)
	}
	if vdf.readOffset != 3 {
		t.Errorf("The world will end, offset: %d", vdf.readOffset)
	}
}

func TestDataFetcherClone(t *testing.T) {
	asset := []uint8{1, 2, 3}
	vdf := VideoDataFetcher{asset: &asset}
	vdfFinal := vdf.cloneWithUpdatedOffset(2)
	if vdfFinal.readOffset != 2 {
		t.Errorf("The world will end, offset: %d", vdfFinal.readOffset)
	}
	if vdfFinal.asset != vdf.asset {
		t.Errorf("The world will end, data")
	}
}
