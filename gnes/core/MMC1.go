package gnes

const (
	size_PRG_RAM  = 0x2000
	size_PRG_ROM1 = 0x4000
	size_PRG_ROM2 = 0x4000
)

const (
	addr_PRG_RAM  = 0x6000
	addr_PRG_ROM1 = 0x8000
	addr_PRG_ROM2 = 0xC000
	addr_END      = 0x10000
)

const (
	region_PRG_RAM  = 1
	region_PRG_ROM1 = 2
	region_PRG_ROM2 = 3
)

type mapper_MMC1 struct {
	prgRamSize,
	prgRomSize uint32

	prgRam []byte
	prgRom [][]byte

	romBank1,
	romBank2 int
}

func (*mapper_MMC1) write(val uint8, addr uint16) error {
	return &gError{err_ADDR_OUT_OF_BOUNDS}
}

func (mmu *mapper_MMC1) read(addr uint16) (uint8, error) {
	ptr, err := mmu.getAddrPointer(addr)
	if err != nil {
		return 0, err
	}
	return *ptr, nil
}

func (mmu *mapper_MMC1) getAddrPointer(addr uint16) (*uint8, error) {
	region, err := mmu.addrToRegion(addr)
	if err != nil {
		return nil, err
	}

	switch region {
	case region_PRG_RAM:
		return &mmu.prgRam[addr-addr_PRG_RAM], nil
	case region_PRG_ROM1:
		return &mmu.prgRom[mmu.romBank1][addr-addr_PRG_ROM1], nil
	case region_PRG_ROM2:
		return &mmu.prgRom[mmu.romBank2][addr-addr_PRG_ROM2], nil
	}
	return nil, &gError{errType: err_ADDR_OUT_OF_BOUNDS}
}

func (*mapper_MMC1) addrToRegion(addr uint16) (int, error) {
	if addr_PRG_RAM <= addr && addr < addr_PRG_ROM1 {
		return region_PRG_RAM, nil
	} else if addr_PRG_ROM1 <= addr && addr < addr_PRG_ROM2 {
		return region_PRG_ROM1, nil
	} else if addr_PRG_ROM2 <= addr {
		return region_PRG_ROM2, nil
	} else {
		return 0, &gError{err_ADDR_OUT_OF_BOUNDS}
	}

}
func newMapper_MMC1(info *cartInfo) (mapper, error) {
	mapper := &mapper_MMC1{}

	if uint32(len(info.data.prgRom))/PRG_ROM_SIZE != info.prgRomSize {
		return nil, &gError{err_INCONSISTENT_PRG_ROM_SIZE}
	}
	mapper.prgRomSize = info.prgRomSize
	mapper.prgRom = make([][]byte, mapper.prgRomSize)

	for i := uint32(0); i < mapper.prgRomSize; i++ {
		mapper.prgRom[i] = info.data.prgRom[i*PRG_ROM_SIZE : (i+1)*PRG_ROM_SIZE]
	}

	mapper.romBank1 = 0
	mapper.romBank2 = int(mapper.prgRomSize) - 1

	if info.prgRamSize > 1 {
		return nil, &gError{err_INCONSISTENT_PRG_RAM_SIZE}
	}
	mapper.prgRam = make([]byte, info.prgRamSize*PRG_RAM_SIZE)
	mapper.prgRamSize = info.prgRamSize

	return mapper, nil
}
