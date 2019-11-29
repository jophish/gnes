package gnes

const (
	vector_RESET_HI = 0xFFFD
	vector_RESET_LO = 0xFFFC
)
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

const (
	REGION_INTERNAL_RAM        = 1
	REGION_INTERNAL_RAM_MIRROR = 2
	REGION_PPU_REG             = 3
	REGION_PPU_REG_MIRROR      = 4
	REGION_APU_IO_REG          = 5
	REGION_APU_IO_TEST         = 6
	REGION_CART_SPACE          = 7
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

// mmu contains all memory accessible to all subsystems of the NES, and is the sole
// interface through which subsystems read and write memory.
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

func (mmu *mmu) getAddrPointer(addr uint16) (*uint8, error) {
	region, err := getAddrRegion(addr)
	if err != nil {
		return nil, err
	}

	var ptr *uint8
	switch region {
	case REGION_INTERNAL_RAM:
		ptr = &mmu.ram[addr]
	case REGION_INTERNAL_RAM_MIRROR:
		ptr = &mmu.ram[addr%INTERNAL_RAM_SIZE]
	//case REGION_PPU_REG:
	//case REGION_PPU_REG_MIRROR:
	//case REGION_APU_IO_REG:
	//case REGION_APU_IO_TEST:
	case REGION_CART_SPACE:
		ptr, err = mmu.mapper.getAddrPointer(addr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, &gError{err_ADDR_OUT_OF_BOUNDS}
	}

	return ptr, nil
}

func (mmu *mmu) read(addr uint16) (uint8, error) {
	ptr, err := mmu.getAddrPointer(addr)
	if err != nil {
		return 0, err
	}
	return *ptr, nil
}

func (mmu *mmu) read16(addr uint16) (uint16, error) {
	lowByte, err := mmu.read(addr)
	if err != nil {
		return 0, err
	}
	highByte, err := mmu.read(addr + 1)
	if err != nil {
		return 0, err
	}
	val := (uint16(highByte) << 8) & uint16(lowByte)
	return val, nil
}

func (mmu *mmu) write(val uint8, addr uint16) error {
	return nil
}

func getAddrRegion(addr uint16) (int, error) {
	if INTERNAL_RAM_ADDR <= addr && addr < INTERNAL_RAM_MIRROR_ADDR {
		return REGION_INTERNAL_RAM, nil
	} else if INTERNAL_RAM_MIRROR_ADDR <= addr && addr < PPU_REG_ADDR {
		return REGION_INTERNAL_RAM_MIRROR, nil
	} else if PPU_REG_ADDR <= addr && addr < PPU_REG_MIRROR_ADDR {
		return REGION_PPU_REG, nil
	} else if PPU_REG_MIRROR_ADDR <= addr && addr < APU_IO_REG_ADDR {
		return REGION_PPU_REG_MIRROR, nil
	} else if APU_IO_REG_ADDR <= addr && addr < APU_IO_TEST_ADDR {
		return REGION_APU_IO_REG, nil
	} else if APU_IO_TEST_ADDR <= addr && addr < CART_SPACE_ADDR {
		return REGION_APU_IO_TEST, nil
	} else if CART_SPACE_ADDR <= addr {
		return REGION_CART_SPACE, nil
	} else {
		return 0, &gError{err_ADDR_OUT_OF_BOUNDS}
	}
}
