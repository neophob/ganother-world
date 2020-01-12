package main

type GamePartContent struct {
	palette   int
	bytecode  int
	cinematic int
	video2    int
}

// The game is divided in 10 parts - each part has its own code, palette and videos
// kPartCopyProtection
// kPartIntro
// kPartWater
// kPartPrison
// kPartCite
// kPartArene
// kPartLuxe
// kPartFinal
// kPartPassword
func getGameParts() map[int]GamePartContent {
	GP_PALETTE := [10]int{0x14, 0x17, 0x1A, 0x1D, 0x20, 0x23, 0x26, 0x29, 0x7D, 0x7D}
	GP_BYTECODE := [10]int{0x15, 0x18, 0x1B, 0x1E, 0x21, 0x24, 0x27, 0x2A, 0x7E, 0x7E}
	GP_CINEMATIC := [10]int{0x16, 0x19, 0x1C, 0x1F, 0x22, 0x25, 0x28, 0x2B, 0x7f}
	GP_VIDEO2 := [10]int{0x00, 0x00, 0x11, 0x11, 0x11, 0x00, 0x11, 0x11, 0x00, 0x00}

	gamePartsMap := make(map[int]GamePartContent)

	for gamePart := 0; gamePart < 10; gamePart++ {
		entry := GamePartContent{
			palette:   GP_PALETTE[gamePart],
			bytecode:  GP_BYTECODE[gamePart],
			cinematic: GP_CINEMATIC[gamePart],
			video2:    GP_VIDEO2[gamePart],
		}
		gamePartsMap[gamePart] = entry
	}

	return gamePartsMap
}
