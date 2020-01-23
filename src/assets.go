package main

//TODO rename to staticGameAssets
type Assets struct {
	memList         map[int]MemlistEntry
	gameParts       map[int]GamePartContent
	bank            map[int][]uint8
	loadedResources map[int][]uint8
}

// this is a function for the Assets struct
func (assets Assets) loadEntryFromBank(index int) []uint8 {
	memlistEntry := assets.memList[index]
	bank := assets.bank[int(memlistEntry.bankID)]
	ofs := int(memlistEntry.bankOffset)
	result := bank[ofs : ofs+int(memlistEntry.packedSize)]
	if memlistEntry.packedSize == memlistEntry.unpackedSize {
		Info("loadResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
		return result
	}
	Info("loadAndUnpackResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
	returnValue, _ := unpack(result)
	return returnValue
}

func (assets *Assets) loadResource(id int) {
	if len(assets.loadedResources[id]) > 0 {
		Info("resource already loaded", id)
		return
	}
	assets.loadedResources[id] = assets.loadEntryFromBank(id)
}
