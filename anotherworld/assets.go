package anotherworld

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neophob/ganother-world/logger"
)

type StaticGameAssets struct {
	MemList         map[int]MemlistEntry
	GameParts       map[int]GamePartContent
	Bank            map[int][]uint8
	LoadedResources map[int][]uint8
}

//LoadEntryFromBank extract a resource from a serialized bank
func (assets StaticGameAssets) LoadEntryFromBank(index int) []uint8 {
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

//LoadResource return a resource
func (assets *StaticGameAssets) LoadResource(id int) {
	if len(assets.LoadedResources[id]) > 0 {
		logger.Info("resource [%d] already loaded", id)
		return
	}
	assets.LoadedResources[id] = assets.LoadEntryFromBank(id)
}

//ReadFile reads a file from disc
func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("File reading error %v", err)
		os.Exit(1)
	}
	return data
}

//CreateBankMap loads the bank asset files from disk and returns a map with its content
func CreateBankMap(assetPath string) map[int][]byte {
	bankFilesMap := make(map[int][]byte)
	for i := 0x01; i < 0x0e; i++ {
		name := fmt.Sprintf("%sbank%02x", assetPath, i)
		logger.Debug("- load file %s", name)
		entry := ReadFile(name)
		bankFilesMap[i] = entry
	}
	return bankFilesMap
}
