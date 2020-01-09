/*
This is the text used for go doc
*/
package main

import (
	"fmt"
	"io/ioutil"
)

type ResourceType int

const (
	RT_SOUND          ResourceType = 0
	RT_MUSIC          ResourceType = 1
	RT_POLY_ANIM      ResourceType = 2
	RT_PALETTE        ResourceType = 3
	RT_BYTECODE       ResourceType = 4
	RT_POLY_CINEMATIC ResourceType = 5
	RT_TODO           ResourceType = 6
	RT_END            ResourceType = 255
	RT_UNKNOWN        ResourceType = -1
)

// TODO this is somehow not so DRY. maybe use iota?
func getResourceType(resourceType byte) ResourceType {
	switch resourceType {
	case 0:
		return RT_SOUND
	case 1:
		return RT_MUSIC
	case 2:
		return RT_POLY_ANIM
	case 3:
		return RT_PALETTE
	case 4:
		return RT_BYTECODE
	case 5:
		return RT_POLY_CINEMATIC
	case 6:
		return RT_TODO
	case 255:
		return RT_END
	}
	fmt.Println("unknown", resourceType)
	return RT_UNKNOWN
}

type MemlistEntry struct {
	state        uint8 //ofs: 0
	resourceType ResourceType
	bufPtr       uint8 //ofs: 2
	//unknownOffset4  uint16 //ofs: 3
	rankNum    uint8  //ofs: 6
	bankId     uint8  //ofs: 7
	bankOffset uint32 //ofs: 8
	//unknownOffset8  uint16 //ofs: 11
	packedSize uint16 //ofs: 14
	//unknownOffset10 uint16 //ofs: 15
	size uint16 //ofs: 18
}

const MemlistEntrySize = 20

func main() {
	data, err := ioutil.ReadFile("../assets/memlist.bin")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	var count = 0
	var sizePacked = 0
	var sizeUncompressed = 0

	for i := 0; i < len(data); i += MemlistEntrySize {
		var entry = MemlistEntry{
			state:        data[i],
			resourceType: getResourceType(data[i+1]),
			rankNum:      data[i+6],
			bankId:       data[i+7],
			//bankOffset:		data[i + 7],
			packedSize: toUint16(data[i+14], data[i+15]),
			size:       toUint16(data[i+18], data[i+19]),
		}
		fmt.Println("entry", count, entry)
		count++
		sizeUncompressed += int(entry.size)
		sizePacked += int(entry.packedSize)
	}
	fmt.Println("sizeUncompressed", sizeUncompressed)
	fmt.Println("sizePacked", sizePacked)
}

func toUint16(lo byte, hi byte) uint16 {
	return uint16(hi) | uint16(lo)<<8
}
