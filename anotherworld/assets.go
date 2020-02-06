package anotherworld

import (
	"github.com/neophob/ganother-world/logger"
)

//TODO rename to staticGameAssets
type Assets struct {
	MemList         map[int]MemlistEntry
	GameParts       map[int]GamePartContent
	Bank            map[int][]uint8
	LoadedResources map[int][]uint8
}

// this is a function for the Assets struct
func (assets Assets) LoadEntryFromBank(index int) []uint8 {
	memlistEntry := assets.MemList[index]
	bank := assets.Bank[int(memlistEntry.BankID)]
	ofs := int(memlistEntry.BankOffset)
	result := bank[ofs : ofs+int(memlistEntry.PackedSize)]
	if memlistEntry.PackedSize == memlistEntry.UnpackedSize {
		logger.Info("LoadResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
		return result
	}
	logger.Info("loadAndUnpackResource on bank %d size %d, offset %v", index, len(bank), memlistEntry)
	returnValue, _ := unpack(result)
	return returnValue
}

func (assets *Assets) LoadResource(id int) {
	if len(assets.LoadedResources[id]) > 0 {
		logger.Info("resource [%d] already loaded", id)
		return
	}
	assets.LoadedResources[id] = assets.LoadEntryFromBank(id)
}
