package main

type Renderer interface {
	drawString(color, posX, posY, stringId int)
	drawShape(color, zoom, posX, posY int)
	fillPage(page, color int)
	copyPage(src, dst, vscroll int)
	updateDisplay(page int)
	setDataBuffer(offset int)
	setWorkPagePtr(page int)
	setPalette(index int)
}
