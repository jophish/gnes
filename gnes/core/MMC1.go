package gnes

import "errors"
import "fmt"

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
	addr_REG_CONTROL   = 0x8000
	addr_REG_CHR_BANK0 = 0xA000
	addr_REG_CHR_BANK1 = 0xC000
	addr_REG_PRG_BANK  = 0xE000
)

const (
	region_PRG_RAM       = 1
	region_PRG_ROM1      = 2
	region_PRG_ROM2      = 3
	region_REG_CONTROL   = 4
	region_REG_CHR_BANK0 = 5
	region_REG_CHR_BANK1 = 6
	region_REG_PRG_BANK  = 7
)

const (
	PRG_ROM_MODE_32K    = 0
	PRG_ROM_MODE_FIX_LO = 1
	PRG_ROM_MODE_FIX_HI = 2
)

const (
	CHR_ROM_MODE_8K = 0
	CHR_ROM_MODE_4K = 1
)

const (
	RESET_MASK             = 0x80
	NEW_WRITE_MASK         = 0x10
	DONE_WRITE_MASK        = 0x1
	DATA_BIT_MASK          = 0x1
	DATA_MASK              = 0x1F
	MIRRORING_MASK         = 0x3
	PRG_ROM_BANK_MODE_MASK = 0xC
	CHR_ROM_BANK_MODE_MASK = 0x10
)

type mapper_MMC1 struct {
	prgRomSize uint32

	prgRam []byte
	prgRom [][]byte

	prgRomBank uint8 // Select 16 KB PRG ROM bank (low bit ignored in PRG_ROM_MODE_32K)

	chrBank1, // Select 4 KB or 8 KB CHR bank at PPU $0000 (low bit ignored in CHR_ROM_MODE_8K)
	chrBank2 uint8 // Select 4 KB CHR bank at PPU $1000 (ignored in CHR_ROM_MODE_8K)

	prgRomMode,
	chrRomMode int

	prgRamEnable bool

	shiftReg uint8 // Internal shift register, used for holding temporary state
	ppu      *ppu
}

func (mmu *mapper_MMC1) write(val uint8, addr uint16) error {
	if mmu.addrIsLoadReg(addr) {
		if (val & RESET_MASK) != 0 {
			// If bit 7 of the value is 1, we just reset the shift register's contents
			mmu.shiftReg = NEW_WRITE_MASK
			fmt.Println("Reset MMC1 shift register")
		} else {
			// otherwise, write to the shift register and if full, update state
			newVal := ((val & DATA_BIT_MASK) << 4) | (mmu.shiftReg >> 1)
			if (mmu.shiftReg & DONE_WRITE_MASK) != 0 {
				// If bit 0 of shift register is 1, this is the last bit of data, so update
				// state and reset the shift register.

				register, err := mmu.addrToRegister(addr)
				if err != nil {
					return err
				}

				switch register {
				case region_REG_CONTROL:
					fmt.Println("Updating control register (MMC1)")
					prgRomMode, err := mmu.getPrgRomBankMode((newVal & PRG_ROM_BANK_MODE_MASK) >> 2)

					if err != nil {
						return err
					}
					mmu.prgRomMode = prgRomMode

					if (newVal & CHR_ROM_BANK_MODE_MASK) == 0 {
						mmu.chrRomMode = CHR_ROM_MODE_8K
					} else {
						mmu.chrRomMode = CHR_ROM_MODE_4K
					}

					mirrorMode := newVal & MIRRORING_MASK
					err = mmu.ppu.setMirroring(mirrorMode)
					if err != nil {
						return err
					}

				case region_REG_CHR_BANK0:

				case region_REG_CHR_BANK1:
				case region_REG_PRG_BANK:
				}

			} else {
				// otherwise, keep pushing data into shift reg
				fmt.Printf("New MMC1 shift register value: %0#8b\n", newVal)
				mmu.shiftReg = newVal
			}
		}
		return nil
	}

	return &gError{err_ADDR_OUT_OF_BOUNDS}
}

func (*mapper_MMC1) getPrgRomBankMode(val uint8) (int, error) {
	if val == 0x0 || val == 0x1 {
		return PRG_ROM_MODE_32K, nil
	} else if val == 0x2 {
		return PRG_ROM_MODE_FIX_LO, nil
	} else if val == 0x3 {
		return PRG_ROM_MODE_FIX_HI, nil
	} else {
		return 0, &gError1{err_MMC1_INVALID_PRG_ROM_MODE_VAL, uint64(val)}
	}
}
func (mmu *mapper_MMC1) addrIsLoadReg(addr uint16) bool {
	return (addr >= addr_REG_CONTROL)
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
		if mmu.prgRamEnable {
			return &mmu.prgRam[addr-addr_PRG_RAM], nil
		} else {
			return nil, &gError{err_MMC1_PRG_RAM_DISABLED}
		}
	case region_PRG_ROM1:
		switch mmu.prgRomMode {
		case PRG_ROM_MODE_32K:
			return &mmu.prgRom[mmu.prgRomBank][addr-addr_PRG_ROM1], nil
		case PRG_ROM_MODE_FIX_LO:
			return &mmu.prgRom[0][addr-addr_PRG_ROM1], nil
		case PRG_ROM_MODE_FIX_HI:
			return &mmu.prgRom[mmu.prgRomBank][addr-addr_PRG_ROM1], nil
		}
	case region_PRG_ROM2:
		switch mmu.prgRomMode {
		case PRG_ROM_MODE_32K:
			return &mmu.prgRom[mmu.prgRomBank+1][addr-addr_PRG_ROM2], nil
		case PRG_ROM_MODE_FIX_LO:
			return &mmu.prgRom[mmu.prgRomBank][addr-addr_PRG_ROM2], nil
		case PRG_ROM_MODE_FIX_HI:
			return &mmu.prgRom[mmu.prgRomSize-1][addr-addr_PRG_ROM2], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Address %#x out of bounds during MMC1 read", addr))
}

func (*mapper_MMC1) addrToRegister(addr uint16) (int, error) {
	if addr_REG_CONTROL <= addr && addr < addr_REG_CHR_BANK0 {
		return region_REG_CONTROL, nil
	} else if addr_REG_CHR_BANK0 <= addr && addr < addr_REG_CHR_BANK1 {
		return region_REG_CHR_BANK0, nil
	} else if addr_REG_CHR_BANK1 <= addr {
		return region_REG_PRG_BANK, nil
	} else {
		return 0, &gError{err_ADDR_OUT_OF_BOUNDS}
	}
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

func newMapper_MMC1(info *cartInfo, ppu *ppu) (mapper, error) {
	mapper := &mapper_MMC1{}

	if uint32(len(info.data.prgRom))/PRG_ROM_SIZE != info.prgRomSize {
		return nil, &gError{err_INCONSISTENT_PRG_ROM_SIZE}
	}
	mapper.prgRomSize = info.prgRomSize
	mapper.prgRom = make([][]byte, mapper.prgRomSize)

	for i := uint32(0); i < mapper.prgRomSize; i++ {
		mapper.prgRom[i] = info.data.prgRom[i*PRG_ROM_SIZE : (i+1)*PRG_ROM_SIZE]
	}

	mapper.prgRomMode = PRG_ROM_MODE_FIX_HI
	mapper.prgRomBank = 0

	mapper.chrRomMode = CHR_ROM_MODE_8K
	mapper.chrBank1 = 0
	mapper.chrBank2 = 0

	if info.prgRamSize > 1 {
		return nil, &gError{err_INCONSISTENT_PRG_RAM_SIZE}
	}
	mapper.prgRam = make([]byte, info.prgRamSize*PRG_RAM_SIZE)
	mapper.prgRamEnable = true

	mapper.shiftReg = 0

	mapper.ppu = ppu
	return mapper, nil
}
