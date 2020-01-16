//dummy video implementation, text output
package main

import (
	"fmt"
)

type DummyRenderer struct {
}

func (render DummyRenderer) setColor(col Color) {
	fmt.Println(">VID: SETCOLOR", col)
}

func (render DummyRenderer) fillPage(page int) {
	fmt.Println(">VID: FILLPAGE", page)
}

func (render DummyRenderer) blitPage(page int) {
	fmt.Println(">VID: BLITPAGE", page)
}

func (render DummyRenderer) copyPage(src, dst int) {
	fmt.Println(">VID: COPYPAGE", src, dst)
}

func (render DummyRenderer) drawChar(posX, posY int32, char byte) {
	fmt.Printf(">VID: DRAWCHAR char:%s, x:%d, y:%d\n", char, posX, posY)
}

func (render DummyRenderer) drawFilledPolygons(vx, vy []int16, col Color) {
	fmt.Println(">VID: DRAWFILLEDPOLYGONS", col, vx, vy)
}

func (render DummyRenderer) eventLoop(frameCount int) bool {
	fmt.Println(">VID: EVENTLOOP", frameCount)
	return frameCount > 128
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
