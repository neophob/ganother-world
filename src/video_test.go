package main

import (
	"testing"
)

func TestVideoAssetFetch(t *testing.T) {
	va := VideoAssets{}
	if va.fetchByte() != 0 {
		t.Errorf("The world will end")
	}
	if va.fetchWord() != 0 {
		t.Errorf("The world will end")
	}
}
