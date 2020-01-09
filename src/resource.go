/*
This is the text used for go doc
*/
package main

import (
	"fmt"
	"math"
)

const (
	MemlistEntrySize int = 20
)

type MemlistEntry struct {
	state        uint8  //ofs: 0
	resourceType uint8  //ofs: 1
	bufPtr       uint8  //ofs: 2
	rankNum      uint8  //ofs: 6
	bankId       uint8  //ofs: 7
	bankOffset   uint32 //ofs: 8
	packedSize   uint16 //ofs: 14
	size         uint16 //ofs: 18
}

func unmarshallingMemlistBin(data []byte) map[int]MemlistEntry {
	resourceMap := make(map[int]MemlistEntry)
	resourceId := 0

	for i := 0; i < len(data); i += MemlistEntrySize {
		entry := MemlistEntry{
			state:        data[i],
			resourceType: data[i+1],
			bufPtr:       data[i+2],
			rankNum:      data[i+6],
			bankId:       data[i+7],
			bankOffset:   toUint32(data[i+8], data[i+9], data[i+10], data[i+11]),
			packedSize:   toUint16(data[i+14], data[i+15]),
			size:         toUint16(data[i+18], data[i+19]),
		}
		resourceMap[resourceId] = entry
		resourceId++
	}
	return resourceMap
}

func loadEntryFromBank() {
	// TODO will have to read a MemlistEntry from the bank file 0x01 - 0x0d
}

func toUint16(lo, hi byte) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func toUint32(b1, b2, b3, b4 byte) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}

func printStatisticsForMemlistBin(resourceMap map[int]MemlistEntry) {
	entryCount := 0
	sizeCompressed, sizeUncompressed, compressedEntries := 0, 0, 0

	var resourceTypeCount [10]int

	for index, entry := range resourceMap {
		fmt.Println(" entry", index, entry)
		entryCount++
		sizeUncompressed += int(entry.size)
		sizeCompressed += int(entry.packedSize)
		if entry.size != entry.packedSize {
			compressedEntries++
		}
		entryType := entry.resourceType
		if int(entryType) < len(resourceTypeCount) {
			resourceTypeCount[entryType] = resourceTypeCount[entryType] + 1
		}
	}

	fmt.Println("===")
	fmt.Println("Total # resources:", len(resourceMap))
	fmt.Println("Compressed       :", compressedEntries)
	fmt.Println("Uncompressed     :", len(resourceMap)-compressedEntries)
	var compressionRatio float64 = 100 / float64(len(resourceMap)) * float64(compressedEntries)
	fmt.Printf("Note: %.0f%% of resources are compressed.\n\n\n", math.Round(compressionRatio))
	fmt.Printf("Total size (uncompressed) : %d bytes.\n", sizeUncompressed)
	fmt.Printf("Total size (compressed)   : %d bytes.\n", sizeCompressed)
	var compressionGain float64 = 100 * (1 - float64(sizeCompressed)/float64(sizeUncompressed))
	fmt.Printf("Note: Overall compression gain is : %.0f%%.\n\n\n", math.Round(compressionGain))
	for i := 0; i < len(resourceTypeCount); i++ {
		if resourceTypeCount[i] > 0 {
			// TODO name me:
			/* 	RT_SOUND          int = 0
			RT_MUSIC          int = 1
			RT_POLY_ANIM      int = 2
			RT_PALETTE        int = 3
			RT_BYTECODE       int = 4
			RT_POLY_CINEMATIC int = 5
			RT_TODO           int = 6
			RT_END            int = 255
			*/
			fmt.Printf("Total %d          files: %d\n", i, resourceTypeCount[i])
		}
	}
}
