package main

func calcStep(x1, y1, x2, y2 int) (int, int) {
	dy := y2 - y1
	delta := dy
	if delta == 0 {
		delta = 1
	}
	return ((x2 - x1) << 16) / delta, dy
}

//software version to draw a polygon
func (video *Video) drawFilledPolygons(page int, vx, vy []int16, col int) {
	/*renderer := render.screenRenderer[page]
	gfx.FilledPolygonColor(renderer, vx, vy, sdl.Color{col.r, col.g, col.g, 255})
	return/**/

	//	render.setColor(col)

	i := 0
	j := len(vx) - 1

	x2 := int(vx[i])
	x1 := int(vx[j])
	hliney := int(vy[i])
	if vy[j] < vy[i] {
		hliney = int(vy[j])
	}

	i++
	j--

	cpt1 := x1 << 16
	cpt2 := x2 << 16

	for numVertices := len(vx); numVertices > 0; numVertices -= 2 {
		step1, _ := calcStep(int(vx[j+1]), int(vy[j+1]), int(vx[j]), int(vy[j]))
		step2, h := calcStep(int(vx[i-1]), int(vy[i-1]), int(vx[i]), int(vy[i]))

		i++
		j--

		cpt1 = (cpt1 & 0xFFFF0000) | 0x7FFF
		cpt2 = (cpt2 & 0xFFFF0000) | 0x8000

		if h == 0 {
			cpt1 += step1
			cpt2 += step2
		} else {
			for ; h > 0; h-- {

				if hliney >= 0 {
					x1 = cpt1 >> 16
					x2 = cpt2 >> 16
					if x1 < int(WIDTH) && x2 >= 0 {
						if x1 < 0 {
							x1 = 0
						}
						if x2 >= int(WIDTH) {
							x2 = int(WIDTH) - 1
						}

						outputOffset := hliney * int(WIDTH)
						for x := x1; x <= x2; x++ {
							color := video.getColor(col, video.workerPage, outputOffset+x)
							video.rawBuffer[video.workerPage][outputOffset+x] = color
						}
					}
				}
				cpt1 += step1
				cpt2 += step2
				hliney++
				if hliney >= int(HEIGHT) {
					return
				}
			}
		}
	}
}
