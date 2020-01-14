package main

type GamePartContent struct {
	palette   int
	bytecode  int
	cinematic int
	video2    int
}

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

func getText(id int) string {
	return keys[id]
}

var keys = map[int]string{
	0x001:  "P E A N U T  3000",
	0x002:  "Copyright  } 1990 Peanut Computer, Inc.\nAll rights reserved.\n\nCDOS Version 5.01",
	0x003:  "2",
	0x004:  "3",
	0x005:  ".",
	0x006:  "A",
	0x007:  "@",
	0x008:  "PEANUT 3000",
	0x00A:  "R",
	0x00B:  "U",
	0x00C:  "N",
	0x00D:  "P",
	0x00E:  "R",
	0x00F:  "O",
	0x010:  "J",
	0x011:  "E",
	0x012:  "C",
	0x013:  "T",
	0x014:  "Shield 9A.5f Ok",
	0x015:  "Flux % 5.0177 Ok",
	0x016:  "CDI Vector ok",
	0x017:  " %%%ddd ok",
	0x018:  "Race-Track ok",
	0x019:  "SYNCHROTRON",
	0x01A:  "E: 23%\ng: .005\n\nRK: 77.2L\n\nopt: g+\n\n Shield:\n1: OFF\n2: ON\n3: ON\n\nP~: 1\n",
	0x01B:  "ON",
	0x01C:  "-",
	0x021:  "|",
	0x022:  "--- Theoretical study ---",
	0x023:  " THE EXPERIMENT WILL BEGIN IN    SECONDS",
	0x024:  "  20",
	0x025:  "  19",
	0x026:  "  18",
	0x027:  "  4",
	0x028:  "  3",
	0x029:  "  2",
	0x02A:  "  1",
	0x02B:  "  0",
	0x02C:  "L E T ' S   G O",
	0x031:  "- Phase 0:\nINJECTION of particles\ninto synchrotron",
	0x032:  "- Phase 1:\nParticle ACCELERATION.",
	0x033:  "- Phase 2:\nEJECTION of particles\non the shield.",
	0x034:  "A  N  A  L  Y  S  I  S",
	0x035:  "- RESULT:\nProbability of creating:\n ANTIMATTER: 91.V %\n NEUTRINO 27:  0.04 %\n NEUTRINO 424: 18 %\n",
	0x036:  "   Practical verification Y/N ?",
	0x037:  "SURE ?",
	0x038:  "MODIFICATION OF PARAMETERS\nRELATING TO PARTICLE\nACCELERATOR (SYNCHROTRON).",
	0x039:  "       RUN EXPERIMENT ?",
	0x03C:  "t---t",
	0x03D:  "000 ~",
	0x03E:  ".20x14dd",
	0x03F:  "gj5r5r",
	0x040:  "tilgor 25%",
	0x041:  "12% 33% checked",
	0x042:  "D=4.2158005584",
	0x043:  "d=10.00001",
	0x044:  "+",
	0x045:  "*",
	0x046:  "% 304",
	0x047:  "gurgle 21",
	0x048:  "{{{{",
	0x049:  "Delphine Software",
	0x04A:  "By Eric Chahi",
	0x04B:  "  5",
	0x04C:  "  17",
	0x12C:  "0",
	0x12D:  "1",
	0x12E:  "2",
	0x12F:  "3",
	0x130:  "4",
	0x131:  "5",
	0x132:  "6",
	0x133:  "7",
	0x134:  "8",
	0x135:  "9",
	0x136:  "A",
	0x137:  "B",
	0x138:  "C",
	0x139:  "D",
	0x13A:  "E",
	0x13B:  "F",
	0x13C:  "        ACCESS CODE:",
	0x13D:  "PRESS BUTTON OR RETURN TO CONTINUE",
	0x13E:  "   ENTER ACCESS CODE",
	0x13F:  "   INVALID PASSWORD !",
	0x140:  "ANNULER",
	0x141:  "      INSERT DISK ?\n\n\n\n\n\n\n\n\nPRESS ANY KEY TO CONTINUE",
	0x142:  " SELECT SYMBOLS CORRESPONDING TO\n THE POSITION\n ON THE CODE WHEEL",
	0x143:  "    LOADING...",
	0x144:  "              ERROR",
	0x15E:  "LDKD",
	0x15F:  "HTDC",
	0x160:  "CLLD",
	0x161:  "FXLC",
	0x162:  "KRFK",
	0x163:  "XDDJ",
	0x164:  "LBKG",
	0x165:  "KLFB",
	0x166:  "TTCT",
	0x167:  "DDRX",
	0x168:  "TBHK",
	0x169:  "BRTD",
	0x16A:  "CKJL",
	0x16B:  "LFCK",
	0x16C:  "BFLX",
	0x16D:  "XJRT",
	0x16E:  "HRTB",
	0x16F:  "HBHK",
	0x170:  "JCGB",
	0x171:  "HHFL",
	0x172:  "TFBB",
	0x173:  "TXHF",
	0x174:  "JHJL",
	0x181:  " BY",
	0x182:  "ERIC CHAHI",
	0x183:  "         MUSIC AND SOUND EFFECTS",
	0x184:  " ",
	0x185:  "JEAN-FRANCOIS FREITAS",
	0x186:  "IBM PC VERSION",
	0x187:  "      BY",
	0x188:  " DANIEL MORAIS",
	0x18B:  "       THEN PRESS FIRE",
	0x18C:  " PUT THE PADDLE ON THE UPPER LEFT CORNER",
	0x18D:  "PUT THE PADDLE IN CENTRAL POSITION",
	0x18E:  "PUT THE PADDLE ON THE LOWER RIGHT CORNER",
	0x258:  "      Designed by ..... Eric Chahi",
	0x259:  "    Programmed by...... Eric Chahi",
	0x25A:  "      Artwork ......... Eric Chahi",
	0x25B:  "Music by ........ Jean-francois Freitas",
	0x25C:  "            Sound effects",
	0x25D:  "        Jean-Francois Freitas\n             Eric Chahi",
	0x263:  "              Thanks To",
	0x264:  "           Jesus Martinez\n\n          Daniel Morais\n\n        Frederic Savoir\n\n      Cecile Chahi\n\n    Philippe Delamarre\n\n  Philippe Ulrich\n\nSebastien Berthet\n\nPierre Gousseau",
	0x265:  "Now Go Out Of This World",
	0x190:  "Good evening professor.",
	0x191:  "I see you have driven here in your\nFerrari.",
	0x192:  "IDENTIFICATION",
	0x193:  "Monsieur est en parfaite sante.",
	0x194:  "Y\n",
	0xFFFF: "",
}
