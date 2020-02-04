package main

import (
	"testing"
)

func TestGetText(t *testing.T) {
	if getText(1) != "P E A N U T  3000" {
		t.Errorf("The world will end")
	}
	if getText(0) != "" {
		t.Errorf("The world will end")
	}
}
