package gnes

import "io/ioutil"
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
	TV_SYS_MASK          = 0x3
	PRG_RAM_PRESENT_MASK = 0x10
	BUS_CONFLICT_MASK    = 0x20
)

const (
	NES2_FLAG            = 0x8
	MIRROR_FLAG          = 0x1
	BATBACK_FLAG         = 0x2
	TRAINER_FLAG         = 0x4
	MIRROR_OVERRIDE_FLAG = 0x8
	VS_FLAG              = 0x1
	PC10_FLAG            = 0x2
	NTSC_FLAG            = 0x0
	PAL_FLAG             = 0x2
	PRG_RAM_PRESENT_FLAG = 0x10
	BUS_CONFLICT_FLAG    = 0x20
)

type iNESHeader struct {
	prgRomSize,
	chrRomSize,
	prgRamSize,
	mapper uint8

	mirror,
	prgRamBatBacked,
	trainer,
	mirrorOverride,
	vs,
	pc10,
	ntsc,
	pal,
	ntscPal,
	prgRamPresent,
	busConflict bool
}
type NES2Header struct {
}

// cartInfo contains iNES or NES 2.0 header information
type cartInfo struct {
	nes2, cartLoaded bool // If nes2, this is a NES2 format file, otherwise iNES
	// Only one of these will be populated at a time
	iNESInfo *iNESHeader
	NES2Info *NES2Header
}
type Emulator struct {
	cpu  *cpu
	mmu  *mmu
	info *cartInfo
}

func NewEmulator() *Emulator {
	emu := &Emulator{}
	emu.info = newCartInfo()
	emu.mmu = newMmu()
	return emu
}

func newCartInfo() *cartInfo {
	cartInfo := &cartInfo{}
	cartInfo.iNESInfo = &iNESHeader{}
	cartInfo.NES2Info = &NES2Header{}
	return cartInfo
}

// Given a path to an iNES or NES2 format ROM,
// loads the ROM such that emulation is ready to begin
func (emu *Emulator) LoadRom(path string) error {
	var err error
	rom, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = emu.info.loadCartInfo(rom)
	if err != nil {
		return err
	}
	return nil
}

// loadCartInfo
func (info *cartInfo) loadCartInfo(rom []byte) error {
	// Check that the header magic constant is correct
	nesConstant := []byte{0x4e, 0x45, 0x53, 0x1a}
	if !bytes.Equal(rom[0:4], nesConstant) {
		return nil
	}

	// Check if it's iNES or NES 2.0
	info.nes2 = (rom[7] & NES2_MASK) == NES2_FLAG

	var err error
	if info.nes2 {
		return nil
	} else {
		err = info.iNESInfo.loadHeader(rom)
		fmt.Printf("%+v\n", *info.iNESInfo)
	}
	if err != nil {
		return err
	}

	return nil
}

// loadHeader fills an iNESHeader with the appropriate data from
// a bytearray representing a rom. Returns an err_BAD_INES_HEADER
// if the rom isn't iNES format (NES2.0), or err_NONZERO_INES_HEADER_BUFFER
// if it is an iNES rom but has nonzero bits [11:15].
func (header *iNESHeader) loadHeader(rom []byte) error {
	// Check to make sure the padding region isn't written over
	nes2 := (rom[7] & NES2_MASK) == NES2_FLAG
	zeroBytes := bytes.Equal(rom[11:15], []byte{0, 0, 0, 0})
	if nes2 {
		return err_BAD_INES_HEADER
	}
	if !zeroBytes {
		return err_NONZERO_INES_HEADER_BUFFER
	}

	header.prgRomSize = rom[4]
	header.chrRomSize = rom[5]
	header.mirror = (rom[6] & MIRROR_MASK) == MIRROR_FLAG
	header.prgRamBatBacked = (rom[6] & BATBACK_MASK) == BATBACK_FLAG
	header.trainer = (rom[6] & TRAINER_MASK) == TRAINER_FLAG
	header.mirrorOverride = (rom[6] & MIRROR_OVERRIDE_MASK) == MIRROR_OVERRIDE_FLAG
	header.vs = (rom[7] & VS_MASK) == VS_FLAG
	header.pc10 = (rom[7] & PC10_MASK) == PC10_FLAG
	header.mapper = ((rom[6] & MAPPER_LOW_NIB_MASK) >> 4) | (rom[7] & MAPPER_HI_NIB_MASK)
	header.prgRamSize = rom[8]
	header.ntsc = (rom[10] & TV_SYS_MASK) == NTSC_FLAG
	header.pal = (rom[10] & TV_SYS_MASK) == PAL_FLAG
	header.ntscPal = !(header.ntsc || header.pal)
	header.prgRamPresent = (rom[10] & PRG_RAM_PRESENT_MASK) == PRG_RAM_PRESENT_FLAG
	header.busConflict = (rom[10] & BUS_CONFLICT_MASK) == BUS_CONFLICT_FLAG
	return nil
}
