/*
Resource handling: load files needed for game, unpack data
*/
package anotherworld

import "github.com/neophob/ganother-world/logger"

const (
	memlistEntrySize int = 20
)

//MemlistEntry is a pointer to game resources
type MemlistEntry struct {
	ResourceType uint8  //ofs: 1
	rankNum      uint8  //ofs: 6
	BankID       uint8  //ofs: 7
	BankOffset   uint32 //ofs: 8
	PackedSize   uint32 //ofs: 12
	UnpackedSize uint32 //ofs: 16
}

//MemlistStatistic contains statistics about resources
type MemlistStatistic struct {
	ResourceTypeMap   map[int]int
	EntryCount        int
	SizeCompressed    int
	SizeUncompressed  int
	CompressedEntries int
}

func UnmarshallingMemlistBin(data []uint8) (map[int]MemlistEntry, MemlistStatistic) {
	resourceMap := make(map[int]MemlistEntry)
	memlistStatistic := MemlistStatistic{ResourceTypeMap: make(map[int]int)}

	for i := 0; i < len(data); i += memlistEntrySize {
		entry := MemlistEntry{
			ResourceType: data[i+1],
			rankNum:      data[i+6],
			BankID:       data[i+7],
			BankOffset:   ToUint32BE(data[i+8], data[i+9], data[i+10], data[i+11]),
			PackedSize:   ToUint32BE(data[i+12], data[i+13], data[i+14], data[i+15]),
			UnpackedSize: ToUint32BE(data[i+16], data[i+17], data[i+18], data[i+19]),
		}
		// Bail out when last entry is found
		if entry.ResourceType == 0xFF {
			break
		}
		logger.Debug("R:%#02x, %-17s size=%5d (%5d)  bank=%2d  offset=%6d", memlistStatistic.EntryCount,
			GetResourceTypeName(int(entry.ResourceType)), entry.UnpackedSize, entry.PackedSize, entry.BankID, entry.BankOffset)
		resourceMap[memlistStatistic.EntryCount] = entry
		memlistStatistic.EntryCount++
		memlistStatistic.SizeCompressed += int(entry.PackedSize)
		memlistStatistic.SizeUncompressed += int(entry.UnpackedSize)
		memlistStatistic.ResourceTypeMap[int(entry.ResourceType)]++
		if entry.UnpackedSize != entry.PackedSize {
			memlistStatistic.CompressedEntries++
		}
	}
	return resourceMap, memlistStatistic
}

func GetResourceTypeName(id int) string {
	resourceNames := [...]string{"RT_SOUND", "RT_MUSIC", "RT_POLY_ANIM", "RT_PALETTE", "RT_BYTECODE", "RT_POLY_CINEMATIC", "RT_COMMON_SHAPES"}
	if id >= 0 && id < len(resourceNames) {
		return resourceNames[id]
	}
	return ""
}
