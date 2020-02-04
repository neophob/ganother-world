package anotherworld

import (
	"testing"
)

func TestVideoAssetsColorInstance(t *testing.T) {
	color := Color{64, 128, 200}
	if color.R != 64 {
		t.Errorf("The world will end, color.r: %d", color.R)
	}
	if color.G != 128 {
		t.Errorf("The world will end, color.g: %d", color.G)
	}
	if color.B != 200 {
		t.Errorf("The world will end, color.b: %d", color.B)
	}
	if color.toUint32() != 0xFF4080C8 {
		t.Errorf("The world will end, color: %d", color.toUint32())
	}
}

func TestVideoAssetsEmpty(t *testing.T) {
	va := VideoAssets{}
	res := va.GetPalette(4)
	if len(res) != 16 {
		t.Errorf("The world will end, len")
	}
}
