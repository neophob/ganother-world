package main

const (
	WIDTH         int32 = 320
	HEIGHT        int32 = 200
	COLOR_ALPHA   int   = 16
	COLOR_BUFFER0 int   = 17
)

const (
	KEY_ESC   uint32 = 0x1
	KEY_LEFT  uint32 = 0x2
	KEY_RIGHT uint32 = 0x4
	KEY_UP    uint32 = 0x8
	KEY_DOWN  uint32 = 0x10
	KEY_FIRE  uint32 = 0x20
)

// implements actual rendering
type Renderer interface {
	blitPage(buffer [64000]Color, posX, posY int)
	eventLoop(frameCount int) uint32
	shutdown()
}

// implements buffer handling (4 buffers) and game specific alpha/buffer0 handling
type Video struct {
	renderer    Renderer
	videoAssets VideoAssets
	workerPage  int
	colors      [16]Color
	rawBuffer   [4][WIDTH * HEIGHT]uint8
	drawColor   uint8
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
	Warn(">VID: UNKNOWN PIXEL %d", colorIndex)
	return uint8(colorIndex)
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
	Debug(">VID: SETWORKPAGEPTR %d", page)
	video.workerPage = getWorkerPage(page)
}

// blit
// step 1 is to convert the indexed color (0..15 for "normal" colors, 16 for translucent and 17 for buffer0 value) to an rgb value
// step 2 is updating the sdl buffers
func (video *Video) updateDisplay(page int) {
	workerPage := getWorkerPage(page)
	Debug(">VID: UPDATEDISPLAY %d", workerPage)

	var outputBuffer [WIDTH * HEIGHT]Color
	for i := range video.rawBuffer[workerPage] {
		outputBuffer[i] = video.colors[video.rawBuffer[workerPage][i]]
	}
	video.renderer.blitPage(outputBuffer, 0, 0)

	//DEBUG OUTPUT
	for i := range video.rawBuffer[0] {
		outputBuffer[i] = video.colors[video.rawBuffer[0][i]]
	}
	video.renderer.blitPage(outputBuffer, 320, 0)
	for i := range video.rawBuffer[1] {
		outputBuffer[i] = video.colors[video.rawBuffer[1][i]]
	}
	video.renderer.blitPage(outputBuffer, 0, 200)
	for i := range video.rawBuffer[2] {
		outputBuffer[i] = video.colors[video.rawBuffer[2][i]]
	}
	video.renderer.blitPage(outputBuffer, 320, 200)
}

func getWorkerPage(page int) int {
	if page >= 0 && page < 4 {
		return page
	}
	switch page {
	case 0xFF:
		return 2
	case 0xFE:
		return 1
	default:
		Warn("updateWorkerPage != [0,1,2,3,0xFF,0xFE] == %d", page)
		return 0
	}
}

func (video *Video) drawString(color, posX, posY, stringID int) {
	text := getText(stringID)
	Debug(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s", color, posX, posY, text)
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
	Debug(">VID: DRAWCHAR char:%c, x:%d, y:%d", char, posX, posY)

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

func (video *Video) drawShape(color, offset, zoom, posX, posY int) {
	video.videoAssets.videoPC = offset
	i := video.videoAssets.fetchByte()

	Debug(">VID: DRAWSHAPE i:%d, color:%d, offset:%d, x:%d, y:%d, zoom:%d", i, color, offset, posX, posY, zoom)

	if i >= 0xC0 {
		if color&0x80 > 0 {
			color = int(i & 0x3F)
		}
		video.drawFilledPolygon(color, zoom, posX, posY)
	} else {
		i &= 0x3F
		if i == 2 {
			video.drawShapeParts(zoom, posX, posY)
		} else {
			Warn("drawShape INVALID! (%d != 2)\n", i)
		}
	}
}

func (video *Video) drawShapeParts(zoom, posX, posY int) {
	x := posX - int(video.videoAssets.fetchByte())*zoom/64
	y := posY - int(video.videoAssets.fetchByte())*zoom/64
	n := int16(video.videoAssets.fetchByte())
	Debug(">VID: DRAWSHAPEPARTS x:%d, y:%d, n:%d", x, y, n)

	for ; n >= 0; n-- {
		off := video.videoAssets.fetchWord()
		_x := x + int(video.videoAssets.fetchByte())*zoom/64
		_y := y + int(video.videoAssets.fetchByte())*zoom/64

		Debug(">VID: DRAWSHAPEPARTS off:%d at %d/%d", off, _x, _y)

		var color uint16 = 0xFF
		if off&0x8000 > 0 {
			color = uint16(video.videoAssets.cinematic[video.videoAssets.videoPC] & 0x7F)
			//TODO display head.. WTF is this?
			video.videoAssets.fetchWord()
		}
		off &= 0x7FFF

		oldVideoPc := video.videoAssets.videoPC
		video.drawShape(int(color), int(off*2), zoom, _x, _y)
		video.videoAssets.videoPC = oldVideoPc
	}
}

func (video *Video) drawFilledPolygon(col, zoom, posX, posY int) {
	Debug(">VID: FILLPOLYGON color:%d, x:%d, y:%d, zoom:%d", col, posX, posY, zoom)

	bbw := int(video.videoAssets.fetchByte()) * zoom / 64
	bbh := int(video.videoAssets.fetchByte()) * zoom / 64

	x1 := posX - bbw/2
	x2 := posX + bbw/2
	y1 := posY - bbh/2
	y2 := posY + bbh/2

	if x1 > 319 || x2 < 0 || y1 > 199 || y2 < 0 {
		//Warn(">VID: FILLPOLYGON INVALID")
		return
	}

	numVertices := int(video.videoAssets.fetchByte())

	if numVertices > 70 {
		Warn(">VID: TOOMANY %d", numVertices)
		panic("UNEXPECTED_AMOUNT_OF_VERTICES")
	}

	var vx, vy = make([]int16, numVertices), make([]int16, numVertices)
	for i := 0; i < numVertices; i++ {
		vx[i] = int16(x1 + int(video.videoAssets.fetchByte())*zoom/64)
		vy[i] = int16(y1 + int(video.videoAssets.fetchByte())*zoom/64)
	}

	Debug(">VID: FILLPOLYGON WorkerPage: %d, numVert: %d, col: %d, %v/%v", video.workerPage, numVertices, col, vx, vy)
	video.drawFilledPolygons(video.workerPage, vx, vy, col)
}

func (video *Video) setPalette(index int) {
	video.colors = video.videoAssets.getPalette(index >> 8)
	Debug(">VID: SETPALETTE %d", index>>8)
}

func (video *Video) eventLoop(frameCount int) uint32 {
	return video.renderer.eventLoop(frameCount)
}

func (video *Video) shutdown() {
	video.renderer.shutdown()
}
