package main

import (
	"fmt"
)

type UnpackCtx struct {
	size int
	crc  uint32
	bits uint32
}

/*struct UnpackCtx {
	int size;
	uint32_t crc;
	uint32_t bits;
	uint8_t *dst;
	const uint8_t *src;
};*/

// unpack crunched data
// last 4 bytes of the data junk contains the unpacked data lenght (uint32 BE)
func unpack(data []byte) {
	dataLen := len(data)
	unpackedSize := toUint32BE(data[dataLen-4], data[dataLen-3], data[dataLen-2], data[dataLen-1])
	crc := toUint32BE(data[dataLen-9], data[dataLen-8], data[dataLen-6], data[dataLen-5])
	bits := toUint32BE(data[dataLen-13], data[dataLen-12], data[dataLen-11], data[dataLen-10])

	unpackCtx := UnpackCtx{int(unpackedSize), crc, bits}

	fmt.Println("XXX", dataLen, unpackCtx)
}
