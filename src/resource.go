/*
This is the text used for go doc
*/
package main

import (
	"fmt"
)

const (
	MemlistEntrySize int = 20
)

type MemlistEntry struct {
	state        uint8  //ofs: 0
	resourceType uint8  //ofs: 1
	rankNum      uint8  //ofs: 6
	bankId       uint8  //ofs: 7
	bankOffset   uint32 //ofs: 8
	packedSize   uint16 //ofs: 14
	size         uint16 //ofs: 18
}

type MemlistStatistic struct {
	resourceTypeCount map[int]int
	entryCount        int
	sizeCompressed    int
	sizeUncompressed  int
	compressedEntries int
}

type Assets struct {
	memList map[int]MemlistEntry
	bank    map[int][]byte
}

func unmarshallingMemlistBin(data []byte) (map[int]MemlistEntry, MemlistStatistic) {
	resourceMap := make(map[int]MemlistEntry)
	memlistStatistic := MemlistStatistic{ resourceTypeCount: make(map[int]int)}

	for i := 0; i < len(data); i += MemlistEntrySize {
		entry := MemlistEntry{
			state:        data[i],
			resourceType: data[i+1],
			rankNum:      data[i+6],
			bankId:       data[i+7],
			bankOffset:   toUint32BE(data[i+8], data[i+9], data[i+10], data[i+11]),
			packedSize:   toUint16BE(data[i+14], data[i+15]),
			size:         toUint16BE(data[i+18], data[i+19]),
		}
		resourceMap[memlistStatistic.entryCount] = entry
		memlistStatistic.entryCount++
		memlistStatistic.sizeCompressed += int(entry.packedSize)
		memlistStatistic.sizeUncompressed += int(entry.size)
		memlistStatistic.resourceTypeCount[int(entry.resourceType)]++
		if entry.size != entry.packedSize {
			memlistStatistic.compressedEntries++
		}
	}
	return resourceMap, memlistStatistic
}

func loadEntryFromBank(assets Assets, index int) {
	memlistEntry := assets.memList[index]
	bank := assets.bank[int(memlistEntry.bankId)]
	fmt.Printf("Bank %d size %d, offset %d\n", index, len(bank), memlistEntry.bankOffset)
	fmt.Println(memlistEntry)
}

func toUint16BE(lo, hi byte) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func toUint32BE(b1, b2, b3, b4 byte) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}

