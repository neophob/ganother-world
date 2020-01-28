package main

type VideoAssets struct {
	palette   []uint8
	cinematic []uint8
	video2    []uint8
}

type Color struct {
	r, g, b uint8
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

func (color Color) toUint32() uint32 {
	return uint32(0xFF000000) + uint32(color.r)<<16 + uint32(color.g)<<8 + uint32(color.b)
}
