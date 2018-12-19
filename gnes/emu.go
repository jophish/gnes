package gnes

import "io/ioutil"
import "errors"
import "bytes"

const (
	NES2_MASK            = b1100
	MIRROR_MASK          = b1
	BATBACK_MASK         = b10
	TRAINER_MASK         = b100
	MIRROR_OVERRIDE_MASK = b1000
	MAPPER_LOW_NIB_MASK  = b11110000
	MAPPER_HI_NIB_MASK   = b11110000
	VS_MASK              = b1
	PC10_MASK            = b10
)

const (
	NES2_FLAG            = b1000
	MIRROR_FLAG          = b1
	BATBACK_FLAG         = b10
	TRAINER_FLAG         = b100
	MIRROR_OVERRIDE_FLAG = b1000
	VS_FLAG              = b1
	PC10_FLAG            = b10
)

// cartInfo contains iNES and NES 2.0 header information
type cartInfo struct {
	submapper,
	prgRamBacked,
	prgRamNonBacked,
	chrRamBacked,
	chrRamNonBacked uint8

	mapper, noPrgRom, noChrRom uint16 // These are word-length for NES 2.0 support

	batBacked,
	trainerLoc,
	mirror,
	overrideMirror,
	vs,
	pc10,
	nes2,
	palMode,
	ntscPal bool
}

type emuInfo struct {
	cartLoaded bool
	cartInfo   *cartInfo
}
type Emulator struct {
	cpu  *Cpu
	mmu  *Mmu
	info *emuInfo
}

func (emu *Emulator) LoadRom(path string) error {
	var err error
	rom, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = emu.info.getCartInfo(rom)
	if err != nil {
		return err
	}
	return nil
}

func (info *emuInfo) getCartInfo(rom []byte) error {
	// Check that the header magic constant is correct
	nesConstant := []byte{0x4e, 0x45, 0x53, 0x1a}
	if !bytes.Equal(rom[0:4], nesConstant) {
		return errors.New("Invalid .nes header")
	}

	// Check if it's iNES or NES 2.0
	info.nes2 = (rom[7] & NES2_MASK) == NES2_FLAG

	// These flags are the same between file formats
	info.mirror = (rom[6] & MIRROR_MASK) == MIRROR_FLAG
	info.batBacked = (rom[6] & BATBACK_MASK) == BATBACK_FLAG
	info.trainerLoc = (rom[6] & TRAINER_MASK) == TRAINER_FLAG
	info.overrideMirror = (rom[6] & MIRROR_OVERRIDE_MASK) == MIRROR_OVERRIDE_FLAG
	// We append the upper 4 bits of the mapper later in the case of NES 2.0
	info.mapper = ((rom[6] & MAPPER_LOW_NIB_MASK) >> 4) | (rom[7] & MAPPER_HI_NIB_MASK)
	// The remainder of the header processing is dependent on the file format
	if info.nes2 {

	} else {

	}

	return nil
}
