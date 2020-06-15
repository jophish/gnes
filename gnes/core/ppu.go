package gnes

import "errors"
import "fmt"

const (
	MIRROR_MODE_SINGLE_LOWER = 0
	MIRROR_MODE_SINGLE_UPPER = 1
	MIRROR_MODE_VERITCAL     = 2
	MIRROR_MODE_HORIZONTAL   = 3
)

const (
	addr_PPU_START = 0x2000
	addr_PPU_END = 0x4000
)

// ppuRegisters represents the registers from PPU_REG_ADDR to PPU_REG_MIRROR
type ppuReginsters struct {
	ppuctrl,
	ppumask,
	ppustatus,
	oamaddr,
	oamdata,
	ppuscroll,
	ppuaddr,
	ppudata uint8
}

type ppu struct {
	mirroring uint8
	regs      *ppuRegisters
}

func newPpu() (*ppu, error) {
	ppu := &ppu{}
	ppu.mirroring = MIRROR_MODE_SINGLE_LOWER
	ppu.regs = &ppuRegisters{}
	return ppu, nil
}

func (ppu *ppu) setMirroring(mirrorMode uint8) error {
	if mirrorMode < MIRROR_MODE_SINGLE_LOWER || mirrorMode > MIRROR_MODE_HORIZONTAL {
		return fmt.Errorf("Invalid mirroring mode %d", mirrorMode)
	}
	ppu.mirroring = mirrorMode
	return nil
}

func (ppu *ppu) write(val uint8, addr uint16) error {
	if addr < addr_PPU_START || addr >= addr_PPU_END {
		return errors.New("Address out of bounds for PPU")
	}

	switch val % 8 {
	case 0:
		ppu.regs.ppuctrl = val
	case 1:
		ppu.regs.ppumask = val
	case 2:
		ppu.regs.ppustatus = val
	case 3:
		ppu.regs.oamaddr = val
	case 4:
		ppu.regs.oamdata = val
	case 5:
		ppu.regs.ppuscroll = val
	case 6:
		ppu.regs.ppuaddr = val
	case 7:
		ppu.regs.ppudata = val
	}

	return nil
}

func (ppu *ppu) getAddrPointer(addr uint16) (*uint8, error) {
	if addr < addr_PPU_START || addr >= addr_PPU_END {
		return nil, errors.New("Address out of bounds for PPU")
	}

	var ptr *uint8

	switch addr % 8 {
	case 0:
		ptr = &ppu.regs.ppuctrl
	case 1:
		ptr = &ppu.regs.ppumask
	case 2:
		ptr = &ppu.regs.ppustatus
	case 3:
		ptr = &ppu.regs.oamaddr
	case 4:
		ptr = &ppu.regs.oamdata
	case 5:
		ptr = &ppu.regs.ppuscroll
	case 6:
		ptr = &ppu.regs.ppuaddr
	case 7:
		ptr = &ppu.regs.ppudata
	}

	return ptr, nil
}

func (ppu *ppu) read(addr uint16) (uint8, error) {
	ptr, err := ppu.getAddrPointer(addr)
	if err != nil {
		return 0, err
	}
	return *ptr, nil
}
