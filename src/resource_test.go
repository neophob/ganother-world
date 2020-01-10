package main

// docs https://golang.org/pkg/testing/

import (
	"testing"
)

func TestConvertUin16FF00(t *testing.T) {
	result := toUint16BE(0xFF, 0)
	if result != 0xFF00 {
		t.Errorf("The world will end: %d", result)
	}
}

func TestConvertUin1600FF(t *testing.T) {
	result := toUint16BE(0, 0xFF)
	if result != 0x00FF {
		t.Errorf("The world will end: %d", result)
	}
}

func TestConvertUin3200FF00FF(t *testing.T) {
	result := toUint32BE(0, 0xFF, 0, 0xFF)
	if result != 0x00FF00FF {
		t.Errorf("The world will end: %d", result)
	}
}

func TestConvertUin32FF000000(t *testing.T) {
	result := toUint32BE(0xFF, 0, 0, 0)
	if result != 0xFF000000 {
		t.Errorf("The world will end: %d", result)
	}
}

func TestResourceTypeName0(t *testing.T) {
	result := getResourceTypeName(0)
	if result != "RT_SOUND" {
		t.Errorf("The world will end: %s", result)
	}
}

func TestResourceTypeName4(t *testing.T) {
	result := getResourceTypeName(4)
	if result != "RT_BYTECODE" {
		t.Errorf("The world will end: %s", result)
	}
}

func TestResourceTypeName44(t *testing.T) {
	result := getResourceTypeName(44)
	if result != "" {
		t.Errorf("The world will end: %s", result)
	}
}
