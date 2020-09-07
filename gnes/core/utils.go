package gnes

func msb(val uint16) uint8 {
	return uint8((val & 0xFF00) >> 8)
}

func lsb(val uint16) uint8 {
	return uint8(val & 0xFF)
}

func make16BitValue(msByte uint8, lsByte uint8) uint16 {
	return (uint16(msByte) << 8) | uint16(lsByte)
}
