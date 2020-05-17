package gnes

import "fmt"

const (
	MIRROR_MODE_SINGLE_LOWER = 0
	MIRROR_MODE_SINGLE_UPPER = 1
	MIRROR_MODE_VERITCAL     = 2
	MIRROR_MODE_HORIZONTAL   = 3
)

type ppu struct {
	mirroring uint8
}

func newPpu() (*ppu, error) {
	ppu := &ppu{}
	ppu.mirroring = MIRROR_MODE_SINGLE_LOWER
	return ppu, nil
}

func (ppu *ppu) setMirroring(mirrorMode uint8) error {
	if mirrorMode < MIRROR_MODE_SINGLE_LOWER || mirrorMode > MIRROR_MODE_HORIZONTAL {
		return fmt.Errorf("Invalid mirroring mode %d", mirrorMode)
	}
	ppu.mirroring = mirrorMode
	return nil
}
