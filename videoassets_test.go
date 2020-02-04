package main

import (
	"testing"
)

func TestVideoAssetsColorInstance(t *testing.T) {
	color := Color{64, 128, 200}
	if color.r != 64 {
		t.Errorf("The world will end, color.r: %d", color.r)
	}
	if color.g != 128 {
		t.Errorf("The world will end, color.g: %d", color.g)
	}
	if color.b != 200 {
		t.Errorf("The world will end, color.b: %d", color.b)
	}
	if color.toUint32() != 0xFF4080C8 {
		t.Errorf("The world will end, color: %d", color.toUint32())
	}
}

func TestVideoAssetsEmpty(t *testing.T) {
	va := VideoAssets{}
	res := va.getPalette(4)
	if len(res) != 16 {
		t.Errorf("The world will end, len")
	}
}
