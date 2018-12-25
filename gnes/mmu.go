package gnes

const (
	INTERNAL_RAM_SIZE        = 0x800
	INTERNAL_RAM_MIRROR_SIZE = 0x1800
	PPU_REG_SIZE             = 0x8
	PPU_REG_MIRROR_SIZE      = 0x1FF8
	APU_IO_REG_SIZE          = 0x18
	APU_IO_TEST_SIZE         = 0x8
	CART_SPACE_SIZE          = 0xBFE0
)

const (
	INTERNAL_RAM_ADDR        = 0x0
	INTERNAL_RAM_MIRROR_ADDR = 0x800
	PPU_REG_ADDR             = 0x2000
	PPU_REG_MIRROR_ADDR      = 0x2008
	APU_IO_REG_ADDR          = 0x4000
	APU_IO_TEST_ADDR         = 0x4018
	CART_SPACE_ADDR          = 0x4020
	MEM_END_ADDR             = 0x10000
)

const (
	PPUCTRL_ADDR   = 0x2000
	PPUMASK_ADDR   = 0x2001
	PPUSTATUS_ADDR = 0x2002
	OAMADDR_ADDR   = 0x2003
	OAMDATA_ADDR   = 0x2004
	PPUSCROLL_ADDR = 0x2005
	PPUADDR_ADDR   = 0x2006
	PPUDATA_ADDR   = 0x2007
)

// ppuRegisters represnts the registers from PPU_REG_ADDR to PPU_REG_MIRROR
type ppuRegisters struct {
	ppuctrl,
	ppumask,
	ppustatus,
	oamaddr,
	oamdata,
	ppuscroll,
	ppuaddr,
	ppudata uint8
}

type apuRegisters struct {
}

type mmu struct {
	mapper  mapper
	ram     [INTERNAL_RAM_SIZE]byte
	ppuRegs *ppuRegisters
	apuRegs *apuRegisters
}

func newMmu(mapperNum uint32, info *cartInfo) (*mmu, error) {
	mmu := &mmu{}
	mapper, err := numberToMapper(mapperNum, info)
	if err != nil {
		return nil, err
	}
	mmu.mapper = mapper
	mmu.ppuRegs = &ppuRegisters{}
	return mmu, nil
}

func (*mmu) read(addr uint8) (uint8, error) {
	return 0, nil
}

func (*mmu) write(val, addr uint8) error {
	return nil
}
