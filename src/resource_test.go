package main

// docs https://golang.org/pkg/testing/

import (
	"testing"
)

func TestHello(t *testing.T) {
	if 2 == 1 {
		t.Errorf("The world will end: %d", 2)
	}
}
