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

func (render DummyRenderer) drawPixel(page int, posX, posY int32) {
	fmt.Printf(">VID: DRAWPIXEL x:%d, y:%d\n", page, posX, posY)
}

func (render DummyRenderer) drawFilledPolygons(page int, vx, vy []int16, col Color) {
	fmt.Println(">VID: DRAWFILLEDPOLYGONS", page, col, vx, vy)
}

func (render DummyRenderer) eventLoop(frameCount int) bool {
	fmt.Println(">VID: EVENTLOOP", frameCount)
	return frameCount > 128
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}
