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
}

func (*mapper_MMC1) write(val, addr uint8) error {
	return nil
}

func (*mapper_MMC1) read(addr uint8) (uint8, error) {
	return 0, nil
}

func newMapper_MMC1(info *cartInfo) mapper {

	return &mapper_MMC1{}
}
