package main

import "github.com/neophob/ganother-world/logger"

const (
	WIDTH         int32 = 320
	HEIGHT        int32 = 200
	COLOR_ALPHA   int   = 16
	COLOR_BUFFER0 int   = 17
)

//Video implements buffer handling (4 buffers) and game specific alpha/buffer0 handling
type Video struct {
	hal         HAL
	videoAssets VideoAssets
	workerPage  int
	colors      [16]Color
	rawBuffer   [4][WIDTH * HEIGHT]uint8
	drawColor   uint8
}

func initVideo(noVideoOutput bool) Video {
	if noVideoOutput == false {
		return Video{hal: buildSDLHAL(), workerPage: 0xFE}
	}
	return Video{hal: DummyHAL{}}
}

func (video *Video) updateGamePart(videoAssets VideoAssets) {
	video.videoAssets = videoAssets
	video.colors = videoAssets.getPalette(0)
}

func (video *Video) getColor(colorIndex, page, ofs int) uint8 {
	if colorIndex < COLOR_ALPHA {
		return uint8(colorIndex)
	}
	if colorIndex == COLOR_ALPHA {
		//Alpha color, take current value and add 8 to the index, then map it
		i := video.rawBuffer[page][ofs]
		return (i | 0x08) & 0x0F
	}
	if colorIndex == COLOR_BUFFER0 {
		//Color from buffer 0 (aka transparent?)
		i := video.rawBuffer[0][ofs]
		return i & 0x0F
	}
	logger.Warn(">VID: UNKNOWN PIXEL COLOR %d", colorIndex)
	return uint8(colorIndex & 0x0F)
}

func (video *Video) setColor(colorIndex int) {
	video.drawColor = uint8(colorIndex)
}

func (video *Video) fillPage(page, colorIndex int) {
	workerPage := getWorkerPage(page)
	for i := range video.rawBuffer[workerPage] {
		video.rawBuffer[workerPage][i] = uint8(colorIndex)
	}
}

func (video *Video) copyPage(src, dst, vscroll int) {
	workerPageSrc := getWorkerPage(src)
	workerPageDst := getWorkerPage(dst)

	switch {
	case vscroll == 0:
		// copy full page
		for i := range video.rawBuffer[workerPageSrc] {
			video.rawBuffer[workerPageDst][i] = video.rawBuffer[workerPageSrc][i]
		}
	case vscroll < 0:
		// copy upper part of screen
		pixelToCopy := (int(HEIGHT) + vscroll) * int(WIDTH)
		verticalOffset := vscroll * int(WIDTH)
		destOffset := 0
		for i := 0; i < pixelToCopy; i++ {
			video.rawBuffer[workerPageDst][destOffset] = video.rawBuffer[workerPageSrc][i-verticalOffset]
			destOffset++
		}
	case vscroll > 0:
		// copy lower part of screen
		pixelToCopy := (int(HEIGHT) - vscroll) * int(WIDTH)
		verticalOffset := vscroll * int(WIDTH)
		sourceOffset := 0
		for i := 0; i < pixelToCopy; i++ {
			video.rawBuffer[workerPageDst][verticalOffset+i] = video.rawBuffer[workerPageSrc][sourceOffset]
			sourceOffset++
		}
	}
}

func (video *Video) setWorkPagePtr(page int) {
	logger.Debug(">VID: SETWORKPAGEPTR %d", page)
	video.workerPage = getWorkerPage(page)
}

// blit
// step 1 is to convert the indexed color (0..15 for "normal" colors, 16 for translucent and 17 for buffer0 value) to an rgb value
// step 2 is updating the sdl buffers
func (video *Video) updateDisplay(page int) {
	var workerPage int
	if page != 0xFE {
		if page == 0xFF {
			//SWAP buffer 1 and 2
			video.rawBuffer[1], video.rawBuffer[2] = video.rawBuffer[2], video.rawBuffer[1]
			workerPage = 1
		} else {
			workerPage = getWorkerPage(page)
		}
	}
	logger.Debug(">VID: UPDATEDISPLAY %d(%d)", page, workerPage)

	var outputBuffer [WIDTH * HEIGHT]Color
	for i := range video.rawBuffer[workerPage] {
		outputBuffer[i] = video.colors[video.rawBuffer[workerPage][i]]
	}
	video.hal.blitPage(outputBuffer, 0, 0)

	//DEBUG OUTPUT
	/*	for i := range video.rawBuffer[0] {
			outputBuffer[i] = video.colors[video.rawBuffer[0][i]]
		}
		video.hal.blitPage(outputBuffer, 320, 0)
		for i := range video.rawBuffer[1] {
			outputBuffer[i] = video.colors[video.rawBuffer[1][i]]
		}
		video.hal.blitPage(outputBuffer, 0, 200)
		for i := range video.rawBuffer[2] {
			outputBuffer[i] = video.colors[video.rawBuffer[2][i]]
		}
		video.hal.blitPage(outputBuffer, 320, 200)*/
}

func getWorkerPage(page int) int {
	if page >= 0 && page < 4 {
		return page
	}
	switch page {
	case 0x40:
		//this is a hack, rawgl does prevent this case
		return 0
	case 0xFF:
		return 2
	case 0xFE:
		return 1
	default:
		logger.Warn("updateWorkerPage != [0,1,2,3,0xFF,0xFE] == %d", page)
		return 0
	}
}

func (video *Video) drawString(color, posX, posY, stringID int) {
	text := getText(stringID)
	logger.Debug(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s", color, posX, posY, text)
	//setWorkPagePtr(buffer);?

	video.setColor(color)
	charPosX := int32(posX)
	charPosY := int32(posY)
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			charPosY += int32(FONT_HEIGHT)
			charPosX = int32(posX)
		} else {
			video.drawChar(charPosX*8, charPosY, text[i])
			charPosX += 1
		}
	}
}

func (video *Video) drawChar(posX, posY int32, char byte) {
	logger.Debug(">VID: DRAWCHAR char:%c, x:%d, y:%d", char, posX, posY)

	fontOffset := 8 * (int32(char) - 0x20)
	for j := int32(0); j < 8; j++ {
		ch := FONT[fontOffset+j]
		outputOffset := posX + (posY+j)*WIDTH
		for i := int32(0); i < 8; i++ {
			if ch&(1<<(7-i)) > 0 {
				video.rawBuffer[video.workerPage][outputOffset+i] = video.drawColor
			}
		}
	}
}

func (video *Video) drawShape(videoDataFetcher VideoDataFetcher, color, zoom, posX, posY int) {
	i := videoDataFetcher.fetchByte()
	logger.Debug(">VID: DRAWSHAPE i:%d, color:%d, fetcher:%v, x:%d, y:%d, zoom:%d",
		i, color, videoDataFetcher, posX, posY, zoom)

	if i >= 0xC0 {
		if color&0x80 > 0 {
			color = int(i & 0x3F)
		}
		video.drawFilledPolygon(videoDataFetcher, color, zoom, posX, posY)
	} else {
		i &= 0x3F
		if i == 2 {
			video.drawShapeParts(videoDataFetcher, zoom, posX, posY)
		} else {
			logger.Warn("drawShape INVALID! (%d != 2)\n", i)
		}
	}
}

func (video *Video) drawShapeParts(videoDataFetcher VideoDataFetcher, zoom, posX, posY int) {
	x := posX - int(videoDataFetcher.fetchByte())*zoom/64
	y := posY - int(videoDataFetcher.fetchByte())*zoom/64
	n := int16(videoDataFetcher.fetchByte())
	logger.Debug(">VID: DRAWSHAPEPARTS x:%d, y:%d, n:%d", x, y, n)

	for ; n >= 0; n-- {
		off := videoDataFetcher.fetchWord()
		_x := x + int(videoDataFetcher.fetchByte())*zoom/64
		_y := y + int(videoDataFetcher.fetchByte())*zoom/64

		logger.Debug(">VID: DRAWSHAPEPARTS off:%d at %d/%d", off, _x, _y)

		var color uint16 = 0xFF
		if off&0x8000 > 0 {
			color = uint16((*videoDataFetcher.asset)[videoDataFetcher.readOffset] & 0x7F)
			//TODO display head.. WTF is this?
			videoDataFetcher.fetchWord()
		}
		off &= 0x7FFF

		clonedVideoDataFetcher := videoDataFetcher.cloneWithUpdatedOffset(int(off * 2))
		video.drawShape(clonedVideoDataFetcher, int(color), zoom, _x, _y)
	}
}

func (video *Video) drawFilledPolygon(videoDataFetcher VideoDataFetcher, col, zoom, posX, posY int) {
	logger.Debug(">VID: FILLPOLYGON color:%d, x:%d, y:%d, zoom:%d", col, posX, posY, zoom)

	bbw := int(videoDataFetcher.fetchByte()) * zoom / 64
	bbh := int(videoDataFetcher.fetchByte()) * zoom / 64

	x1 := posX - bbw/2
	x2 := posX + bbw/2
	y1 := posY - bbh/2
	y2 := posY + bbh/2

	if x1 > 319 || x2 < 0 || y1 > 199 || y2 < 0 {
		//Warn(">VID: FILLPOLYGON INVALID")
		return
	}

	numVertices := int(videoDataFetcher.fetchByte())

	if numVertices > 70 {
		logger.Warn(">VID: TOOMANY %d", numVertices)
		panic("UNEXPECTED_AMOUNT_OF_VERTICES")
	}

	var vx, vy = make([]int16, numVertices), make([]int16, numVertices)
	for i := 0; i < numVertices; i++ {
		vx[i] = int16(x1 + int(videoDataFetcher.fetchByte())*zoom/64)
		vy[i] = int16(y1 + int(videoDataFetcher.fetchByte())*zoom/64)
	}

	logger.Debug(">VID: FILLPOLYGON WorkerPage: %d, numVert: %d, col: %d, %v/%v", video.workerPage, numVertices, col, vx, vy)
	video.drawFilledPolygons(video.workerPage, vx, vy, col)
}

func (video *Video) setPalette(index int) {
	video.colors = video.videoAssets.getPalette(index >> 8)
	//TODO fixup palette
	//part 16004 and palette 0x47 -> ret 8, part 16006 and palette 0x4a -> ret 1
	logger.Debug(">VID: SETPALETTE %d", index>>8)
}

func (video *Video) eventLoop(frameCount int) uint32 {
	return video.hal.eventLoop(frameCount)
}

func (video *Video) shutdown() {
	video.hal.shutdown()
}
