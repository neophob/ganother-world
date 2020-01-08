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
	RT_END            ResourceType = 255
	RT_UNKNOWN        ResourceType = -1
)

// this is somehow not so DRY
func getResourceType(resourceType int) ResourceType {
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
	case 255:
		return RT_END
	}
	return RT_UNKNOWN
}

type MemlistEntry struct {
	state           uint8
	resourceType    ResourceType
	bufPtr          uint8
	unknownOffset4  uint16
	rankNum         uint8
	bankId          uint8
	bankOffset      uint32
	unknownOffset8  uint16
	packedSize      uint16
	unknownOffset10 uint16
	size            uint16
}

const MemlistEntrySize = 20

func main() {
	data, err := ioutil.ReadFile("../assets/memlist.bin")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	for i := 0; i < len(data); i += MemlistEntrySize {
		var entry = MemlistEntry{
			state:        data[i],
			resourceType: getResourceType(data[i+1]),
		}
		fmt.Println("entry:", entry.state, entry.resourceType)
	}
}
