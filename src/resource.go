/*
This is the text used for go doc
*/
package main

import (
	"fmt"
	"math"
)

type ResourceType int

const (
	RT_SOUND          ResourceType = 0
	RT_MUSIC          ResourceType = 1
	RT_POLY_ANIM      ResourceType = 2
	RT_PALETTE        ResourceType = 3
	RT_BYTECODE       ResourceType = 4
	RT_POLY_CINEMATIC ResourceType = 5
	RT_TODO           ResourceType = 6
	RT_END            ResourceType = 255
	RT_UNKNOWN        ResourceType = -1
)

// TODO this is somehow not so DRY. maybe use iota?
func getResourceType(resourceType byte) ResourceType {
	switch resourceType {
	case 0:
		return RT_SOUND
	case 1:
		return RT_MUSIC
	case 2:
		return RT_POLY_ANIM
	case 3:
		return RT_PALETTE
	case 4:
		return RT_BYTECODE
	case 5:
		return RT_POLY_CINEMATIC
	case 6:
		return RT_TODO
	case 255:
		return RT_END
	}
	fmt.Println("unknown", resourceType)
	return RT_UNKNOWN
}

type MemlistEntry struct {
	state        uint8 //ofs: 0
	resourceType ResourceType
	bufPtr       uint8  //ofs: 2
	rankNum      uint8  //ofs: 6
	bankId       uint8  //ofs: 7
	bankOffset   uint32 //ofs: 8
	packedSize   uint16 //ofs: 14
	size         uint16 //ofs: 18
}

const MemlistEntrySize int = 20

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
			fmt.Printf("Total %d          files: %d\n", i, resourceTypeCount[i])
		}
	}
}

func unmarshallingMemlistBin(data []byte) map[int]MemlistEntry {
	resourceMap := make(map[int]MemlistEntry)
	resourceId := 0

	for i := 0; i < len(data); i += MemlistEntrySize {
		entry := MemlistEntry{
			state:        data[i],
			resourceType: getResourceType(data[i+1]),
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

func toUint16(lo, hi byte) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func toUint32(b1, b2, b3, b4 byte) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}
