//dummy video implementation, text output
package main

import (
	"fmt"
)

type DummyRenderer struct {
}

func (render DummyRenderer) drawString(color, posX, posY, stringId int) {
	text := getText(stringId)
	fmt.Printf(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s\n", color, posX, posY, text)
}

func (render DummyRenderer) drawShape(color, offset, zoom, posX, posY int) {
	fmt.Printf(">VID: DRAWSHAPE offset:%d, x:%d, y:%d, zoom:%d\n", offset, posX, posY, zoom)
}

func (render DummyRenderer) fillPage(page, color int) {
	fmt.Println(">VID: FILLPAGE", page, color)
}

func (render DummyRenderer) copyPage(src, dst, vscroll int) {
	fmt.Println(">VID: COPYPAGE", src, dst, vscroll)
}

// blit
func (render DummyRenderer) updateDisplay(page int) {
	fmt.Println(">VID: UPDATEDISPLAY", page)
}

func (render DummyRenderer) setWorkPagePtr(page int) {
	fmt.Println(">VID: SETWORKPAGEPTR", page)
}

func (render DummyRenderer) setPalette(index int) {
	fmt.Println(">VID: SETPALETTE", index>>8)
	//TODO	_vid->_nextPal = num >> 8
}

func (render DummyRenderer) mainLoop() {
	//nothin to see here, move on!
}

func (render DummyRenderer) shutdown() {
	//nothin to see here, move on!
}

func (render DummyRenderer) updateGamePart(videoAssets VideoAssets) {
	//nothin to see here, move on!
}

func (render DummyRenderer) exitRequested(frameCount int) bool {
	return frameCount > 128
}
