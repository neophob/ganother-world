//dummy video implementation, text output
package main

import (
	"fmt"
)

type DummyRenderer struct {
}

func (render DummyRenderer) blitPage(buffer [64000]Color) {
	fmt.Println(">VID: BLITPAGE")
}

func (render DummyRenderer) eventLoop(frameCount int) bool {
	fmt.Println(">VID: EVENTLOOP", frameCount)
	return frameCount > 128
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
