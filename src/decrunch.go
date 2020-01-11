package main

import (
	"fmt"
)

type UnpackCtx struct {
	size         uint32
	crc          uint32
	bits         uint32
	source       []byte
	sourceOffset int
}

func (unpackCtx *UnpackCtx) readUInt32BE() uint32 {
	ret := toUint32BE(
		unpackCtx.source[unpackCtx.sourceOffset-3],
		unpackCtx.source[unpackCtx.sourceOffset-2],
		unpackCtx.source[unpackCtx.sourceOffset-1],
		unpackCtx.source[unpackCtx.sourceOffset-0],
	)
	unpackCtx.sourceOffset -= 4
	return ret
}

func (unpackCtx *UnpackCtx) nextBit() bool {
	carry := (unpackCtx.bits & 1) != 0
	unpackCtx.bits >>= 1
	if unpackCtx.bits == 0 {
		unpackCtx.bits = unpackCtx.readUInt32BE()
		unpackCtx.crc ^= unpackCtx.bits
		carry = (unpackCtx.bits & 1) != 0
		unpackCtx.bits = (1 << 31) | (unpackCtx.bits >> 1)
	}
	return carry
}

func (unpackCtx *UnpackCtx) getBits(bitsCount int) int {
	return 0
}

func (unpackCtx *UnpackCtx) copyLiteral(bitsCount int, len int) {
}

func (unpackCtx *UnpackCtx) copyReference(bitsCount int, count int) {
}

// unpack crunched data
// last 4 bytes of the data junk contains the unpacked data lenght (uint32 BE)
func unpack(data []byte) {
	dataLen := len(data)
	unpackCtx := UnpackCtx{0, 0, 0, data, dataLen - 1}

	unpackCtx.size = unpackCtx.readUInt32BE()
	unpackCtx.crc = unpackCtx.readUInt32BE()
	unpackCtx.bits = unpackCtx.readUInt32BE()

	fmt.Println("XXX", dataLen, unpackCtx.size, unpackCtx.crc, unpackCtx.bits, unpackCtx.sourceOffset)

	// put current pointer to file end minus the 3 read uint32 values..
	unpackCtx.sourceOffset = dataLen - 1 - 4*3

	unpackCtx.crc ^= unpackCtx.bits
	for unpackCtx.size > 0 {
		if !unpackCtx.nextBit() {
			if !unpackCtx.nextBit() {
				unpackCtx.copyLiteral(3, 0)
			} else {
				unpackCtx.copyReference(8, 2)
			}
		} else {
			code := unpackCtx.getBits(2)
			switch code {
			case 3:
				unpackCtx.copyLiteral(8, 8)
			case 2:
				unpackCtx.copyReference(12, unpackCtx.getBits(8)+1)
			case 1:
				unpackCtx.copyReference(10, 4)
			case 0:
				unpackCtx.copyReference(9, 3)
			}
		}
	}
	//crc should be 0!
}
