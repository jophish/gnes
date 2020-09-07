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
	VBLANK_BIT_MASK           uint8 = 0x80
	PPUSTATUS_UNUSED_BIT_MASK uint8 = 0x1F
)
const (
	size_PATTERN_TABLE_0    = 0x1000
	size_PATTERN_TABLE_1    = 0x1000
	size_NAMETABLE_0        = 0x400
	size_NAMETABLE_1        = 0x400
	size_NAMETABLE_2        = 0x400
	size_NAMETABLE_3        = 0x400
	size_PPU_VRAM           = 0x800
	size_NAMETABLE_MIRROR   = 0xF00
	size_PALETTE_RAM        = 0x20
	size_PALETTE_RAM_MIRROR = 0xE0
)

const (
	addr_PPU_START = 0x2000
	addr_PPU_END   = 0x4000
)

// ppuRegisters represents the raw registers from PPU_REG_ADDR to PPU_REG_MIRROR
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
	openLatch uint8

	cycles        uint64
	catchupCycles uint64

	regs *ppuRegisters
	vram [size_PPU_VRAM]byte

	currentScanline      uint16
	currentScanlineCycle uint16
	currentFrame         uint64
}

func newPpu() (*ppu, error) {
	ppu := &ppu{}

	ppu.cycles = 0
	ppu.catchupCycles = 0

	ppu.currentScanline = 261
	ppu.currentScanlineCycle = 0
	ppu.currentFrame = 0

	ppu.mirroring = MIRROR_MODE_SINGLE_LOWER
	ppu.regs = &ppuRegisters{}
	ppu.openLatch = 0

	return ppu, nil
}

// catchup catches the PPU up by ppu.catchupCycles. Since CPU/PPU execution is staggered
// and emulator execution is driven by the CPU, we need to catch the PPU up at each execution step
// in proportion with the number of cycles spent by the CPU.
func (ppu *ppu) catchup() error {

	for ppu.catchupCycles > 0 {
		// We've reached the end of the scanline. Increment scanline and reset scanline cycles
		if ppu.currentScanlineCycle == 341 {
			ppu.currentScanline++
			ppu.currentScanlineCycle = 0
		}

		// We've completed all scanlines. Increment the frame counter, and reset scanline.
		if ppu.currentScanline == 262 {
			ppu.currentFrame++
			ppu.currentScanline = 0
		}

		// Actually do things
		if ppu.currentScanline >= 0 && ppu.currentScanline <= 239 {
			// Visible scanlines
			ppu.currentScanlineCycle++
			ppu.catchupCycles--

		} else if ppu.currentScanline == 240 {
			// Post-render scanlines
			ppu.currentScanlineCycle++
			ppu.catchupCycles--

		} else if ppu.currentScanline >= 241 && ppu.currentScanline <= 260 {
			// Vertical blanking scanlines
			if ppu.currentScanline == 241 && ppu.currentScanlineCycle == 1 {
				ppu.regs.ppustatus |= VBLANK_BIT_MASK
			}

			ppu.currentScanlineCycle++
			ppu.catchupCycles--

		} else if ppu.currentScanline == 261 {
			// Pre-render scanline
			if ppu.currentScanlineCycle == 1 {
				ppu.regs.ppustatus &= ^VBLANK_BIT_MASK
			}
			ppu.currentScanlineCycle++
			ppu.catchupCycles--
		}
	}

	return nil
}

func (ppu *ppu) setMirroring(mirrorMode uint8) error {
	if mirrorMode < MIRROR_MODE_SINGLE_LOWER || mirrorMode > MIRROR_MODE_HORIZONTAL {
		return fmt.Errorf("Invalid mirroring mode %d", mirrorMode)
	}
	ppu.mirroring = mirrorMode
	return nil
}

// getCPUAddrPointer returns a pointer to a PPU-related register available in the CPU's memory map
func (ppu *ppu) getCPUAddrPointer(addr uint16) (*uint8, error) {
	if addr < addr_PPU_START || addr >= addr_PPU_END {
		return nil, errors.New("Address out of bounds for PPU")
	}

	var ptr *uint8

	switch addr % 8 {
	case 0:
		return nil, gErrorNew(err_UNWRITEABLE_PPU_REG)
	case 1:
		return nil, gErrorNew(err_UNWRITEABLE_PPU_REG)
	case 2:
		ptr = &ppu.regs.ppustatus
	case 3:
		return nil, gErrorNew(err_UNWRITEABLE_PPU_REG)
	case 4:
		ptr = &ppu.regs.oamdata
	case 5:
		return nil, gErrorNew(err_UNWRITEABLE_PPU_REG)
	case 6:
		return nil, gErrorNew(err_UNWRITEABLE_PPU_REG)
	case 7:
		ptr = &ppu.regs.ppudata
	}

	return ptr, nil
}

// writeCPUReg writes to a PPU-related register available on the CPU's memory map
func (ppu *ppu) writeCPU(val uint8, addr uint16) error {
	if addr < addr_PPU_START || addr >= addr_PPU_END {
		return errors.New("Address out of bounds for PPU")
	}

	switch val % 8 {
	case 0:
		ppu.regs.ppuctrl = val
	case 1:
		ppu.regs.ppumask = val
	case 2:
		break
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
	ppu.openLatch = val

	return nil
}

// readCPU reads a PPU-related value available on the CPU's memory map
func (ppu *ppu) readCPU(addr uint16) (uint8, error) {
	if addr < addr_PPU_START || addr >= addr_PPU_END {
		return 0, errors.New("Address out of bounds for PPU")
	}

	var val uint8

	switch addr % 8 {
	case 0:
		val = ppu.openLatch
	case 1:
		val = ppu.openLatch
	case 2:
		val = (ppu.regs.ppustatus & ^PPUSTATUS_UNUSED_BIT_MASK) | (ppu.openLatch & PPUSTATUS_UNUSED_BIT_MASK)
		ppu.regs.ppustatus &= ^VBLANK_BIT_MASK
		ppu.openLatch = val
	case 3:
		val = ppu.openLatch
	case 4:
		val = ppu.regs.oamdata
		ppu.openLatch = val
	case 5:
		val = ppu.openLatch
	case 6:
		val = ppu.openLatch
	case 7:
		val = ppu.regs.ppudata
		ppu.openLatch = val
	}

	return val, nil
}
