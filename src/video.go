package main

import "fmt"

type Renderer interface {
	setColor(col Color)
	fillPage(page int)
	blitPage(page int)
	copyPage(src, dst int)
	drawChar(posX, posY int32, char byte)
	drawFilledPolygons(vx, vy []int16, col Color)
	eventLoop(frameCount int) bool
	shutdown()
}

type Video struct {
	renderer    Renderer
	videoAssets VideoAssets
	workerPage  int
	loadPalette int
	colors      [16]Color
}

func (video *Video) updateGamePart(videoAssets VideoAssets) {
	video.videoAssets = videoAssets
	video.colors = videoAssets.getPalette(0)
}

func (video *Video) setColor(colorIndex int) {
	col := video.colors[colorIndex]
	video.renderer.setColor(col)
}

func (video *Video) fillPage(page, colorIndex int) {
	video.setColor(colorIndex)
	workerPage := getWorkerPage(page)
	video.renderer.fillPage(workerPage)
}

func (video *Video) copyPage(src, dst, vscroll int) {
	fmt.Println(">VID: COPYPAGE", vscroll)
	video.renderer.copyPage(src, dst)
}

func (video *Video) setWorkPagePtr(page int) {
	fmt.Println(">VID: SETWORKPAGEPTR", page)
	video.workerPage = getWorkerPage(page)
}

// blit
func (video *Video) updateDisplay(page int) {
	workerPage := getWorkerPage(page)
	fmt.Println(">VID: UPDATEDISPLAY", workerPage)

	if video.loadPalette != 0xFF {
		fmt.Println(">VID: UPDATEPAL", video.loadPalette)
		//render.colors = render.videoAssets.getPalette(render.loadPalette)
		video.loadPalette = 0xFF
	}
	video.renderer.blitPage(workerPage)
}

func getWorkerPage(page int) int {
	if page >= 0 && page <= 3 {
		return page
	}
	switch page {
	case 0xFF:
		return 2
	case 0xFE:
		return 1
	default:
		fmt.Println("updateWorkerPage != [0,1,2,3,0xFF,0xFE] ==", page)
		return 0
	}
}

func (video *Video) drawString(color, posX, posY, stringId int) {
	text := getText(stringId)
	fmt.Printf(">VID: DRAWSTRING color:%d, x:%d, y:%d, text:%s\n", color, posX, posY, text)
	//setWorkPagePtr(buffer);?

	video.setColor(color)
	charPosX := int32(posX)
	charPosY := int32(posY)
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			charPosY += int32(FONT_HEIGHT)
			charPosX = int32(posX)
		} else {
			video.renderer.drawChar(charPosX, charPosY, text[i])
			charPosX += 8
		}
	}
}

func (video *Video) drawShape(color, offset, zoom, posX, posY int) {
	video.videoAssets.videoPC = offset
	i := video.videoAssets.fetchByte()

	fmt.Printf(">VID: DRAWSHAPE i:%d, color:%d, offset:%d, x:%d, y:%d, zoom:%d\n", i, color, offset, posX, posY, zoom)

	if i >= 0xC0 {
		if color&0x80 > 0 {
			color = int(i & 0x3F)
		}
		video.drawFilledPolygon(color, zoom, posX, posY)
	} else {
		i &= 0x3F
		if i == 1 {
			fmt.Printf("drawShape INVALID! (1 != 2)\n")
		} else if i == 2 {
			video.drawShapeParts(zoom, posX, posY)
		} else {
			fmt.Printf("drawShape INVALID! (%d != 2)\n", i)
		}
	}
}

func (video *Video) drawShapeParts(zoom, posX, posY int) {
	x := posX - int(video.videoAssets.fetchByte())*zoom/64
	y := posY - int(video.videoAssets.fetchByte())*zoom/64
	n := int16(video.videoAssets.fetchByte())
	fmt.Printf(">VID: DRAWSHAPEPARTS x:%d, y:%d, n:%d\n", x, y, n)

	for ; n >= 0; n-- {
		off := video.videoAssets.fetchWord()
		_x := x + int(video.videoAssets.fetchByte())*zoom/64
		_y := y + int(video.videoAssets.fetchByte())*zoom/64

		fmt.Printf(">VID: DRAWSHAPEPARTS off:%d at %d/%d\n", off, _x, _y)

		var color uint16 = 0xFF
		if off&0x8000 > 0 {
			readOfs := video.videoAssets.videoPC & 0x7F
			b1 := video.videoAssets.cinematic[readOfs]
			color = uint16(b1)
			//TODO display head.. WTF is this?
			video.videoAssets.fetchWord()
		}
		off &= 0x7FFF

		oldVideoPc := video.videoAssets.videoPC
		video.drawShape(int(color), int(off*2), zoom, _x, _y)
		video.videoAssets.videoPC = oldVideoPc
	}
}

func (video *Video) drawFilledPolygon(color, zoom, posX, posY int) {
	fmt.Printf(">VID: FILLPOLYGON color:%d, x:%d, y:%d, zoom:%d\n", color, posX, posY, zoom)

	bbw := int(video.videoAssets.fetchByte()) * zoom / 64
	bbh := int(video.videoAssets.fetchByte()) * zoom / 64

	x1 := posX - bbw/2
	x2 := posX + bbw/2
	y1 := posY - bbh/2
	y2 := posY + bbh/2

	if x1 > 319 || x2 < 0 || y1 > 199 || y2 < 0 {
		fmt.Println(">VID: FILLPOLYGON INVALID")
		return
	}

	col := video.colors[color%16]
	numVertices := int(video.videoAssets.fetchByte())

	if numVertices > 70 {
		fmt.Println(">VID: TOOMANY", numVertices)
		panic("UNEXPECTED_AMOUNT_OF_VERTICES")
	}

	var vx, vy = make([]int16, numVertices), make([]int16, numVertices)
	for i := 0; i < numVertices; i++ {
		vx[i] = int16(x1 + int(video.videoAssets.fetchByte())*zoom/64)
		vy[i] = int16(y1 + int(video.videoAssets.fetchByte())*zoom/64)
	}

	fmt.Println(">VID: FILLPOLYGON", numVertices, vx, vy)
	video.renderer.drawFilledPolygons(vx, vy, col)
}

func (video *Video) setPalette(index int) {
	//render.loadPalette = index >> 8
	video.colors = video.videoAssets.getPalette(index >> 8)
	fmt.Println(">VID: SETPALETTE", index>>8)
}

func (video *Video) eventLoop(frameCount int) bool {
	return video.renderer.eventLoop(frameCount)
}

type VideoAssets struct {
	palette   []uint8
	cinematic []uint8
	video2    []uint8
	videoPC   int
}

type Color struct {
	r, g, b uint8
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

// each asset is 2048 bytes long, a palette stores 16 colors
// each palette file contains multiple palette files which can be selected using the index
func (asset VideoAssets) getPalette(index int) [16]Color {
	var palette [16]Color

	if len(asset.palette) == 0 {
		return palette
	}

	ofs := index * 32
	for i := 0; i < 16; i++ {
		color := toUint16BE(asset.palette[ofs], asset.palette[ofs+1])
		ofs += 2
		r := uint8(((color >> 8) & 0xF) << 4)
		g := uint8(((color >> 4) & 0xF) << 4)
		b := uint8((color & 0xF) << 4)
		palette[i] = Color{r, g, b}
	}
	return palette
}
