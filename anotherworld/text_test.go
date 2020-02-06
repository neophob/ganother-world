package anotherworld

import (
	"testing"
)

func TestGetText(t *testing.T) {
	if GetText(1) != "P E A N U T  3000" {
		t.Errorf("The world will end")
	}
	if GetText(0) != "" {
		t.Errorf("The world will end")
	}
}
