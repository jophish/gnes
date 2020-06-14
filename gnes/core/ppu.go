package gnes

import "fmt"

const (
	MIRROR_MODE_SINGLE_LOWER = 0
	MIRROR_MODE_SINGLE_UPPER = 1
	MIRROR_MODE_VERITCAL     = 2
	MIRROR_MODE_HORIZONTAL   = 3
)

// ppuRegisters represents the registers from PPU_REG_ADDR to PPU_REG_MIRROR
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
	return nil
}

func (ppu *ppu) getAddrPointer(addr uint16) (*uint8, error) {
	return nil, nil
}
