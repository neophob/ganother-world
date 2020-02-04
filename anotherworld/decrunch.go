package anotherworld

//UnpackCtx encapsulate a decompress operation
type UnpackCtx struct {
	size         uint32
	crc          uint32
	bits         uint32
	source       []uint8
	sourceOffset int
	dest         []uint8
	destOffset   int
}

func (unpackCtx *UnpackCtx) readUInt32BE() uint32 {
	ret := ToUint32BE(
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

func (unpackCtx *UnpackCtx) getBits(count int) int {
	bits := 0
	for i := 0; i < count; i++ {
		bits <<= 1
		if unpackCtx.nextBit() {
			bits |= 1
		}
	}
	return bits
}

func (unpackCtx *UnpackCtx) copyLiteral(bitsCount int, len int) {
	count := unpackCtx.getBits(bitsCount) + len + 1
	unpackCtx.size -= uint32(count)
	if unpackCtx.size < 0 {
		count += int(unpackCtx.size)
		unpackCtx.size = 0
	}
	for i := 0; i < count; i++ {
		unpackCtx.dest[unpackCtx.destOffset-i] = uint8(unpackCtx.getBits(8))
	}
	unpackCtx.destOffset -= count
}

func (unpackCtx *UnpackCtx) copyReference(bitsCount int, count int) {
	unpackCtx.size -= uint32(count)
	if unpackCtx.size < 0 {
		count += int(unpackCtx.size)
		unpackCtx.size = 0
	}
	offset := unpackCtx.getBits(bitsCount)
	for i := 0; i < count; i++ {
		unpackCtx.dest[unpackCtx.destOffset-i] = unpackCtx.dest[unpackCtx.destOffset-i+offset]
	}
	unpackCtx.destOffset -= count
}

// unpack crunched data
// last 4 bytes of the data junk contains the unpacked data lenght (uint32 BE)
func unpack(data []uint8) ([]uint8, uint32) {
	dataLen := len(data)
	unpackCtx := UnpackCtx{source: data, sourceOffset: dataLen - 1}
	unpackCtx.size = unpackCtx.readUInt32BE()
	unpackCtx.crc = unpackCtx.readUInt32BE()
	unpackCtx.bits = unpackCtx.readUInt32BE()
	unpackCtx.dest = make([]uint8, unpackCtx.size)
	unpackCtx.destOffset = int(unpackCtx.size - 1)

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
	return unpackCtx.dest, unpackCtx.crc
}
