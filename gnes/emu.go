package gnes

import "io/ioutil"
import "errors"
import "bytes"
import "fmt"

const (
	NES2_MASK            = 0xc
	MIRROR_MASK          = 0x1
	BATBACK_MASK         = 0x2
	TRAINER_MASK         = 0x4
	MIRROR_OVERRIDE_MASK = 0x8
	MAPPER_LOW_NIB_MASK  = 0xf0
	MAPPER_HI_NIB_MASK   = 0xf0
	VS_MASK              = 0x1
	PC10_MASK            = 0x2
)

const (
	NES2_FLAG            = 0x8
	MIRROR_FLAG          = 0x1
	BATBACK_FLAG         = 0x2
	TRAINER_FLAG         = 0x4
	MIRROR_OVERRIDE_FLAG = 0x8
	VS_FLAG              = 0x1
	PC10_FLAG            = 0x2
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
	cpu  *cpu
	mmu  *mmu
	info *emuInfo
}

func NewEmulator() *Emulator {
	emu := &Emulator{}
	emu.info = newEmuInfo()
	emu.mmu = newMmu()
	return emu
}

func newEmuInfo() *emuInfo {
	emuInfo := &emuInfo{}
	emuInfo.cartInfo = &cartInfo{}
	return emuInfo
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
	cartInfo := info.cartInfo
	// Check if it's iNES or NES 2.0
	cartInfo.nes2 = (rom[7] & NES2_MASK) == NES2_FLAG

	// These flags are the same between file formats
	cartInfo.mirror = (rom[6] & MIRROR_MASK) == MIRROR_FLAG
	cartInfo.batBacked = (rom[6] & BATBACK_MASK) == BATBACK_FLAG
	cartInfo.trainerLoc = (rom[6] & TRAINER_MASK) == TRAINER_FLAG
	cartInfo.overrideMirror = (rom[6] & MIRROR_OVERRIDE_MASK) == MIRROR_OVERRIDE_FLAG
	// We append the upper 4 bits of the mapper later in the case of NES 2.0
	cartInfo.mapper = uint16(((rom[6] & MAPPER_LOW_NIB_MASK) >> 4) | (rom[7] & MAPPER_HI_NIB_MASK))
	// The remainder of the header processing is dependent on the file format
	if cartInfo.nes2 {

	} else {

	}
	fmt.Printf("%+v\n", *cartInfo)
	return nil
}
