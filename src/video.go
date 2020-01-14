//tripple buffering video buffers
//this mock implementation needs to be implemented for
//cli and web/wasm) - using an interface?
package main

import (
	"fmt"
)

//TODO where is the stringId defined?
func drawString(color, posX, posY, stringId int) {
	text := getText(stringId)
	fmt.Printf(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s\n", color, posX, posY, text)
}

func drawShape(color, zoom, posX, posY int) {
	fmt.Printf(">VID: DRAWSHAPE color:%d, x:%d, y:%d, zoom:%d\n", color, posX, posY, zoom)
}

func fillPage(page, color int) {
	fmt.Println(">VID: FILLPAGE", page, color)
}

func copyPage() {
	fmt.Println(">VID: COPYPAGE")
}

// blit
func updateDisplay(page int) {
	fmt.Println(">VID: UPDATEDISPLAY", page)
}

//TODO gimme a better name
func setDataBuffer(offset int) {
	fmt.Println(">VID: SETDATABUFFER", offset)
}

func setWorkPagePtr(page int) {
	fmt.Println(">VID: SETWORKPAGEPTR", page)
}
