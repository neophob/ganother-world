package main

type VideoAssets struct {
	palette   []uint8
	cinematic []uint8
	video2    []uint8

	videoPC int
}

func (asset *VideoAssets) fetchByte() uint8 {
	if len(asset.cinematic) == 0 {
		return 0
	}
	result := asset.cinematic[asset.videoPC]
	asset.videoPC++
	return result
}

func (asset *VideoAssets) fetchWord() uint16 {
	if len(asset.cinematic) == 0 {
		return 0
	}
	b1 := asset.cinematic[asset.videoPC]
	b2 := asset.cinematic[asset.videoPC+1]
	asset.videoPC += 2
	return toUint16BE(b1, b2)
}

type Renderer interface {
	drawString(color, posX, posY, stringId int)
	drawShape(color, zoom, posX, posY int)
	fillPage(page, color int)
	copyPage(src, dst, vscroll int)
	updateDisplay(page int)
	setDataBuffer(useSecondVideo bool, offset int)
	setWorkPagePtr(page int)
	setPalette(index int)
	mainLoop()
	shutdown()
	exitRequested(frameCount int) bool
	updateGamePart(videoAssets VideoAssets)
}
