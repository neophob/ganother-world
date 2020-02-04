package main

import (
	"github.com/neophob/ganother-world/anotherworld"
	"github.com/neophob/ganother-world/logger"
)

//TODO rename to staticGameAssets
type Assets struct {
	memList         map[int]anotherworld.MemlistEntry
	gameParts       map[int]GamePartContent
	bank            map[int][]uint8
	loadedResources map[int][]uint8
}

// this is a function for the Assets struct
func (assets Assets) loadEntryFromBank(index int) []uint8 {
	memlistEntry := assets.memList[index]
	bank := assets.bank[int(memlistEntry.BankID)]
	ofs := int(memlistEntry.BankOffset)
	result := bank[ofs : ofs+int(memlistEntry.PackedSize)]
	if memlistEntry.PackedSize == memlistEntry.UnpackedSize {
		logger.Info("loadResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
		return result
	}
	logger.Info("loadAndUnpackResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
	returnValue, _ := unpack(result)
	return returnValue
}

func (assets *Assets) loadResource(id int) {
	if len(assets.loadedResources[id]) > 0 {
		logger.Info("resource [%d] already loaded", id)
		return
	}
	assets.loadedResources[id] = assets.loadEntryFromBank(id)
}
