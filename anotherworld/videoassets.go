package anotherworld

//VideoAssets holds the assets for the currently active game part
type VideoAssets struct {
	Palette   []uint8
	Cinematic []uint8
	Video2    []uint8
}

//Color represents an indexed color palette
type Color struct {
	R, G, B uint8
}

// each asset is 2048 bytes long, a palette stores 16 colors
// each palette file contains multiple palette files which can be selected using the index
func (asset VideoAssets) GetPalette(index int) [16]Color {
	var palette [16]Color

	if len(asset.Palette) == 0 {
		return palette
	}

	ofs := index * 32
	for i := 0; i < 16; i++ {
		color := ToUint16BE(asset.Palette[ofs], asset.Palette[ofs+1])
		ofs += 2
		r := uint8(((color >> 8) & 0xF) << 4)
		g := uint8(((color >> 4) & 0xF) << 4)
		b := uint8((color & 0xF) << 4)
		palette[i] = Color{r, g, b}
	}
	return palette
}

func (color Color) toUint32() uint32 {
	return uint32(0xFF000000) + uint32(color.R)<<16 + uint32(color.G)<<8 + uint32(color.B)
}
