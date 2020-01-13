package main

import (
	"fmt"
)

type Assets struct {
	memList   map[int]MemlistEntry
	gameParts map[int]GamePartContent
	bank      map[int][]uint8
}

// this is a function for the Assets struct
func (assets Assets) loadEntryFromBank(index int) []uint8 {
	memlistEntry := assets.memList[index]
	bank := assets.bank[int(memlistEntry.bankId)]
	fmt.Printf("Bank %d size %d, offset %v\n", index, len(bank), memlistEntry)
	fmt.Println("slice", memlistEntry.bankOffset, memlistEntry.packedSize)
	ofs := int(memlistEntry.bankOffset)
	result := bank[ofs : ofs+int(memlistEntry.packedSize)]
	if memlistEntry.packedSize == memlistEntry.unpackedSize {
		return result
	}
	returnValue, _ := unpack(result)
	return returnValue
}
