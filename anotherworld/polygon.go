package anotherworld

func calcStep(x1, y1, x2, y2 int) (int, int) {
	dy := y2 - y1
	delta := dy
	if delta == 0 {
		delta = 1
	}
	return ((x2 - x1) << 16) / delta, dy
}

//software version to draw a polygon, draw polygon line by line
func (video *Video) drawFilledPolygons(page int, vx, vy []int16, col int) {
	i := 0
	j := len(vx) - 1

	xPositionEnd := int(vx[i])
	xPositionStart := int(vx[j])
	yPosition := int(vy[i])
	if vy[j] < vy[i] {
		yPosition = int(vy[j])
	}

	i++
	j--

	cpt1 := xPositionStart << 16
	cpt2 := xPositionEnd << 16

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

				if yPosition >= 0 {
					xPositionStart = int(int16(cpt1 >> 16))
					xPositionEnd = int(int16(cpt2 >> 16))
					if xPositionStart < int(WIDTH) && xPositionEnd >= 0 {
						if xPositionStart < 0 {
							xPositionStart = 0
						}
						if xPositionEnd >= int(WIDTH) {
							xPositionEnd = int(WIDTH) - 1
						}

						outputOffset := yPosition * int(WIDTH)
						for x := xPositionStart; x <= xPositionEnd; x++ {
							color := video.getColor(col, video.WorkerPage, outputOffset+x)
							video.rawBuffer[video.WorkerPage][outputOffset+x] = color
						}
					}
				}
				cpt1 += step1
				cpt2 += step2
				yPosition++
				if yPosition >= int(HEIGHT) {
					return
				}
			}
		}
	}
}
