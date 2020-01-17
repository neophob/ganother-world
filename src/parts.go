package main

//AW has 10 different game parts. here we map which assets (code, video, palettes) belong to which level

type GamePartContent struct {
	palette   int
	bytecode  int
	cinematic int
	video2    int
}

const GAME_PARTS_COUNT = 10

const (
	GAME_PART_FIRST int = 0x3E80
	GAME_PART_LAST  int = 0x3E89
	GAME_PART_ID_1  int = 0x3E80
	GAME_PART_ID_2  int = 0x3E81 //Introductino
	GAME_PART_ID_3  int = 0x3E82
	GAME_PART_ID_4  int = 0x3E83 //Wake up in the suspended jail
	GAME_PART_ID_5  int = 0x3E84
	GAME_PART_ID_6  int = 0x3E85 //BattleChar sequence
	GAME_PART_ID_7  int = 0x3E86
	GAME_PART_ID_8  int = 0x3E87
	GAME_PART_ID_9  int = 0x3E88
	GAME_PART_ID_10 int = 0x3E89
)

// The game is divided in 10 parts - each part has its own code, palette and videos
// kPartCopyProtection 16000
// kPartIntro
// kPartWater
// kPartPrison
// kPartCite
// kPartArene
// kPartLuxe
// kPartFinal
// kPartPassword
func getGameParts() map[int]GamePartContent {
	GP_PALETTE := [GAME_PARTS_COUNT]int{0x14, 0x17, 0x1A, 0x1D, 0x20, 0x23, 0x26, 0x29, 0x7D, 0x7D}
	GP_BYTECODE := [GAME_PARTS_COUNT]int{0x15, 0x18, 0x1B, 0x1E, 0x21, 0x24, 0x27, 0x2A, 0x7E, 0x7E}
	GP_CINEMATIC := [GAME_PARTS_COUNT]int{0x16, 0x19, 0x1C, 0x1F, 0x22, 0x25, 0x28, 0x2B, 0x7f}
	GP_VIDEO2 := [GAME_PARTS_COUNT]int{0x00, 0x00, 0x11, 0x11, 0x11, 0x00, 0x11, 0x11, 0x00, 0x00}

	gamePartsMap := make(map[int]GamePartContent)

	for gamePart := 0; gamePart < GAME_PARTS_COUNT; gamePart++ {
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
