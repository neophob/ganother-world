/*
Resource handling: load files needed for game, unpack data
*/
package main

const (
	memlistEntrySize int = 20
)

//MemlistEntry is a pointer to game resources
type MemlistEntry struct {
	resourceType uint8  //ofs: 1
	rankNum      uint8  //ofs: 6
	bankID       uint8  //ofs: 7
	bankOffset   uint32 //ofs: 8
	packedSize   uint32 //ofs: 12
	unpackedSize uint32 //ofs: 16
}

//MemlistStatistic contains statistics about resources
type MemlistStatistic struct {
	resourceTypeMap   map[int]int
	entryCount        int
	sizeCompressed    int
	sizeUncompressed  int
	compressedEntries int
}

func unmarshallingMemlistBin(data []uint8) (map[int]MemlistEntry, MemlistStatistic) {
	resourceMap := make(map[int]MemlistEntry)
	memlistStatistic := MemlistStatistic{resourceTypeMap: make(map[int]int)}

	for i := 0; i < len(data); i += memlistEntrySize {
		entry := MemlistEntry{
			resourceType: data[i+1],
			rankNum:      data[i+6],
			bankID:       data[i+7],
			bankOffset:   toUint32BE(data[i+8], data[i+9], data[i+10], data[i+11]),
			packedSize:   toUint32BE(data[i+12], data[i+13], data[i+14], data[i+15]),
			unpackedSize: toUint32BE(data[i+16], data[i+17], data[i+18], data[i+19]),
		}
		// Bail out when last entry is found
		if entry.resourceType == 0xFF {
			break
		}
		Debug("R:%#02x, %-17s size=%5d (%5d)  bank=%2d  offset=%6d", memlistStatistic.entryCount,
			getResourceTypeName(int(entry.resourceType)), entry.unpackedSize, entry.packedSize, entry.bankID, entry.bankOffset)
		resourceMap[memlistStatistic.entryCount] = entry
		memlistStatistic.entryCount++
		memlistStatistic.sizeCompressed += int(entry.packedSize)
		memlistStatistic.sizeUncompressed += int(entry.unpackedSize)
		memlistStatistic.resourceTypeMap[int(entry.resourceType)]++
		if entry.unpackedSize != entry.packedSize {
			memlistStatistic.compressedEntries++
		}
	}
	return resourceMap, memlistStatistic
}

func toUint16BE(lo, hi uint8) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func toUint32BE(b1, b2, b3, b4 uint8) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}

func getResourceTypeName(id int) string {
	resourceNames := [...]string{"RT_SOUND", "RT_MUSIC", "RT_POLY_ANIM", "RT_PALETTE", "RT_BYTECODE", "RT_POLY_CINEMATIC", "RT_COMMON_SHAPES"}
	if id >= 0 && id < len(resourceNames) {
		return resourceNames[id]
	}
	return ""
}
