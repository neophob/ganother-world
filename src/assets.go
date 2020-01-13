package main

import (
	"fmt"
)

type Assets struct {
	memList         map[int]MemlistEntry
	gameParts       map[int]GamePartContent
	bank            map[int][]uint8
	loadedResources map[int][]uint8
}

// this is a function for the Assets struct
func (assets Assets) loadEntryFromBank(index int) []uint8 {
	memlistEntry := assets.memList[index]
	bank := assets.bank[int(memlistEntry.bankId)]
	fmt.Printf("loadResource on bank %d size %d, offset %v\n", index, len(bank), memlistEntry)
	ofs := int(memlistEntry.bankOffset)
	result := bank[ofs : ofs+int(memlistEntry.packedSize)]
	if memlistEntry.packedSize == memlistEntry.unpackedSize {
		return result
	}
	fmt.Println("unpack data")
	returnValue, _ := unpack(result)
	return returnValue
}

func (assets *Assets) loadResource(id int) {
	if len(assets.loadedResources[id]) > 0 {
		fmt.Println("resource already loaded", id)
		return
	}
	assets.loadedResources[id] = assets.loadEntryFromBank(id)
}
