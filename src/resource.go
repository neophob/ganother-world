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
	bufPtr       uint8  //ofs: 2
	rankNum      uint8  //ofs: 6
	bankId       uint8  //ofs: 7
	bankOffset   uint32 //ofs: 8
	packedSize   uint16 //ofs: 14
	size         uint16 //ofs: 18
}

const MemlistEntrySize int = 20

func main() {
	data, err := ioutil.ReadFile("./assets/memlist.bin")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	var entryCount int = 0
	sizePacked, sizeUncompressed, compressedEntries := int(0), 0, 0

	for i := 0; i < len(data); i += MemlistEntrySize {
		var entry = MemlistEntry{
			state:        data[i],
			resourceType: getResourceType(data[i+1]),
			bufPtr:       data[i+2],
			rankNum:      data[i+6],
			bankId:       data[i+7],
			bankOffset:   toUint32(data[i+8], data[i+9], data[i+10], data[i+11]),
			packedSize:   toUint16(data[i+14], data[i+15]),
			size:         toUint16(data[i+18], data[i+19]),
		}
		fmt.Println("entry", entryCount, entry)
		entryCount++
		sizeUncompressed += int(entry.size)
		sizePacked += int(entry.packedSize)
		if entry.size != entry.packedSize {
			compressedEntries++
		}

	}
	fmt.Println("sizeUncompressed", sizeUncompressed)
	fmt.Println("sizePacked", sizePacked)
	fmt.Println("compressedEntries", compressedEntries)
}

func toUint16(lo, hi byte) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func toUint32(b1, b2, b3, b4 byte) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}
