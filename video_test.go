package main

import (
	"testing"
)

func TestVideoGetWorkerPage(t *testing.T) {
	res := getWorkerPage(0)
	if res != 0 {
		t.Errorf("The world will end")
	}
	res = getWorkerPage(3)
	if res != 3 {
		t.Errorf("The world will end")
	}
	res = getWorkerPage(0x40)
	if res != 0 {
		t.Errorf("The world will end")
	}
	res = getWorkerPage(0xFF)
	if res != 2 {
		t.Errorf("The world will end")
	}
	res = getWorkerPage(0xAA)
	if res != 0 {
		t.Errorf("The world will end")
	}
}
