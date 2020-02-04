package anotherworld

import (
	"testing"
)

func TestResourceTypeName0(t *testing.T) {
	result := GetResourceTypeName(0)
	if result != "RT_SOUND" {
		t.Errorf("The world will end: %s", result)
	}
}

func TestResourceTypeName4(t *testing.T) {
	result := GetResourceTypeName(4)
	if result != "RT_BYTECODE" {
		t.Errorf("The world will end: %s", result)
	}
}

func TestResourceTypeName44(t *testing.T) {
	result := GetResourceTypeName(44)
	if result != "" {
		t.Errorf("The world will end: %s", result)
	}
}
