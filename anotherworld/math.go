package anotherworld

func ToUint16BE(lo, hi uint8) uint16 {
	return uint16(hi) | uint16(lo)<<8
}

func ToUint32BE(b1, b2, b3, b4 uint8) uint32 {
	return uint32(b4) | uint32(b3)<<8 | uint32(b2)<<16 | uint32(b1)<<24
}
