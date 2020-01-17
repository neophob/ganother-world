package main

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
