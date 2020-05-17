package gnes

type mapper_NROM struct {
	prgRom [][]byte
	prgRam []byte
	chrRom []byte

	prgRomSize uint32
	ppu        *ppu
}

func (mmu *mapper_NROM) write(val uint8, addr uint16) error {
	return nil
}

func (mmu *mapper_NROM) read(addr uint16) (uint8, error) {
	return 0, nil
}

func (mmu *mapper_NROM) getAddrPointer(addr uint16) (*uint8, error) {
	return nil, nil
}

func newMapper_NROM(info *cartInfo, ppu *ppu) (mapper, error) {
	mapper := &mapper_NROM{}

	if uint32(len(info.data.prgRom))/PRG_ROM_SIZE != info.prgRomSize {
		return nil, &gError{err_INCONSISTENT_PRG_ROM_SIZE}
	}

	mapper.prgRomSize = info.prgRomSize
	mapper.prgRom = make([][]byte, mapper.prgRomSize)

	for i := uint32(0); i < mapper.prgRomSize; i++ {
		mapper.prgRom[i] = info.data.prgRom[i*PRG_ROM_SIZE : (i+1)*PRG_ROM_SIZE]
	}
	mapper.ppu = ppu
	return nil, nil
}
