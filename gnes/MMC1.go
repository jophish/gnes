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
)

type mapper_MMC1 struct {
	prgRamSize,
	prgRomSize uint32

	prgRam []byte
	prgRom [][]byte
}

func (*mapper_MMC1) write(val, addr uint8) error {
	return nil
}

func (*mapper_MMC1) read(addr uint8) (uint8, error) {
	return 0, nil
}

func newMapper_MMC1(info *cartInfo) (mapper, error) {
	mapper := &mapper_MMC1{}

	if uint32(len(info.data.prgRom))/PRG_ROM_SIZE != info.prgRomSize {
		return nil, err_INCONSISTENT_PRG_ROM_SIZE
	}
	mapper.prgRomSize = info.prgRomSize
	mapper.prgRom = make([][]byte, mapper.prgRomSize)

	for i := uint32(0); i < mapper.prgRomSize; i++ {
		mapper.prgRom[i] = info.data.prgRom[i*PRG_ROM_SIZE : (i+1)*PRG_ROM_SIZE]
	}

	if info.prgRamSize > 1 {
		return nil, err_INCONSISTENT_PRG_RAM_SIZE
	}
	mapper.prgRam = make([]byte, info.prgRamSize*PRG_RAM_SIZE)
	mapper.prgRamSize = info.prgRamSize

	return mapper, nil
}
