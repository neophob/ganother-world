# ganother-world

Reference: https://fabiensanglard.net/anotherWorld_code_review/

# EPICS

- Resource Management / read MEMLIST.BIN
- Setup Video (SDL?)
- VM parser
- implement VM op codes (Polygons drawing)
- Audio / Sound FX

## VM parser

- IV's:

```
//The game is divided in 10 parts.
#define GAME_NUM_PARTS 10

#define GAME_PART_FIRST  0x3E80
#define GAME_PART1       0x3E80
#define GAME_PART2       0x3E81   //Introductino
#define GAME_PART3       0x3E82
#define GAME_PART4       0x3E83   //Wake up in the suspended jail
#define GAME_PART5       0x3E84
#define GAME_PART6       0x3E85   //BattleChar sequence
#define GAME_PART7       0x3E86
#define GAME_PART8       0x3E87
#define GAME_PART9       0x3E88
#define GAME_PART10      0x3E89
#define GAME_PART_LAST   0x3E89

```
