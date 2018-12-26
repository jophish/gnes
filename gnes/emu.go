package gnes

import "io/ioutil"
import "bytes"
import "fmt"

// file format section address enum
const (
	TRAINER_START_ADDR = 0x10
)

// file format section size enum
const (
	TRAINER_SIZE = 512
	PRG_ROM_SIZE = 16384
	CHR_ROM_SIZE = 8192
	PC_ROM_SIZE  = 8192
	PC_PROM_SIZE = 32
	PRG_RAM_SIZE = 8192
)

// system type enum
const (
	SYS_NTSC     = 0x1
	SYS_PAL      = 0x2
	SYS_NTSC_PAL = 0x3
)

// bit mask enum
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

// bit flag enum
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

// gameData contains the raw data sections from the rom file
type gameData struct {
	trainer []byte

	prgRom,
	chrRom []byte

	pcInstRom,
	pcProm []byte
}

// cartInfo contains iNES or NES 2.0 header information
type cartInfo struct {
	nes2 bool // If nes2, this is a NES2 format file, otherwise iNES

	prgRomSize,
	chrRomSize,
	prgRamSize,
	mapper,
	system uint32

	mirror,
	prgRamBatBacked,
	trainer,
	mirrorOverride,
	vs,
	pc10,
	prgRamPresent,
	busConflict bool

	data *gameData
}

// main emulator struct
type Emulator struct {
	cpu  *cpu
	mmu  *mmu
	info *cartInfo
}

func (emu *Emulator) ReadCpu(addr uint16) (uint8, error) {
	val, err := emu.mmu.read(addr)
	if err != nil {
		return 0, err
	}
	return val, nil
}
func NewEmulator(path string) (*Emulator, error) {
	emu := &Emulator{}
	emu.info = newCartInfo()
	err := emu.loadRom(path)
	if err != nil {
		return nil, err
	}
	return emu, nil
}

func newCartInfo() *cartInfo {
	info := &cartInfo{}
	info.data = &gameData{}
	return info
}

// Given a path to an iNES or NES2 format ROM,
// loads the ROM such that emulation is ready to begin
func (emu *Emulator) loadRom(path string) error {
	var err error
	rom, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = emu.info.loadCartInfo(rom)
	if err != nil {
		return err
	}

	// We can only initialize the mmu once we know which
	// mapper we need to use
	mmu, err := newMmu(emu.info.mapper, emu.info)
	if err != nil {
		return err
	}
	emu.mmu = mmu
	cpu, err := newCpu(mmu)
	if err != nil {
		return err
	}
	emu.cpu = cpu
	return nil
}

// Step is the main way for external callers to step through emulation. This takes care
// of fetching and executing opcodes/instructions, updating APU and PPU appropriately, etc.
func (emu *Emulator) Step() error {
	err := emu.cpu.stepInstruction()
	if err != nil {
		return err
	}
	return nil
}

// loadCartInfo loads a cartInfo struct with all the available data in the header
// of the given rom, which must be in either iNES or NES2.0 format.
func (info *cartInfo) loadCartInfo(rom []byte) error {
	// Check that the header magic constant is correct
	nesConstant := []byte{0x4e, 0x45, 0x53, 0x1a}
	if !bytes.Equal(rom[0:4], nesConstant) {
		return &gError{err_BAD_MAGIC_CONSTANT}
	}

	// Check if it's iNES or NES 2.0
	info.nes2 = (rom[7] & NES2_MASK) == NES2_FLAG

	var err error
	if info.nes2 {
		return nil
	} else {
		if err := info.loadINESData(rom); err != nil {
			return err
		}
		fmt.Printf("%+v\n", *info)
	}
	if err != nil {
		return err
	}

	return nil
}

// loadCommonData loads a cartInfo struct with the info common to both iNES
// and NES2.0 file formats.
func (info *cartInfo) loadCommonData(rom []byte) error {
	info.prgRomSize = uint32(rom[4])
	info.chrRomSize = uint32(rom[5])
	info.mirror = (rom[6] & MIRROR_MASK) == MIRROR_FLAG
	info.prgRamBatBacked = (rom[6] & BATBACK_MASK) == BATBACK_FLAG
	info.trainer = (rom[6] & TRAINER_MASK) == TRAINER_FLAG
	info.mirrorOverride = (rom[6] & MIRROR_OVERRIDE_MASK) == MIRROR_OVERRIDE_FLAG
	info.vs = (rom[7] & VS_MASK) == VS_FLAG
	info.pc10 = (rom[7] & PC10_MASK) == PC10_FLAG
	info.mapper = uint32(((rom[6] & MAPPER_LOW_NIB_MASK) >> 4) | (rom[7] & MAPPER_HI_NIB_MASK))
	return nil
}

// loadINESData loads a cartInfo struct with data under the assumption that the
// given byte array represents an iNES format ROM.
func (info *cartInfo) loadINESData(rom []byte) error {
	// Check to make sure the padding region isn't written over
	nes2 := (rom[7] & NES2_MASK) == NES2_FLAG
	zeroBytes := bytes.Equal(rom[11:15], []byte{0, 0, 0, 0})
	if nes2 {
		return &gError{err_BAD_INES_HEADER}
	}
	if !zeroBytes {
		return &gError{err_NONZERO_INES_HEADER_BUFFER}
	}
	// Get data that's common to both iNES and NES2.0
	if err := info.loadCommonData(rom); err != nil {
		return err
	}
	info.prgRamSize = uint32(rom[8])
	// We just ignore flag 9, since it's pretty outdated
	if (rom[10] & TV_SYS_MASK) == NTSC_FLAG {
		info.system = SYS_NTSC
	} else if (rom[10] & TV_SYS_MASK) == PAL_FLAG {
		info.system = SYS_PAL
	} else {
		info.system = SYS_NTSC_PAL
	}

	info.prgRamPresent = (rom[10] & PRG_RAM_PRESENT_MASK) == PRG_RAM_PRESENT_FLAG
	info.busConflict = (rom[10] & BUS_CONFLICT_MASK) == BUS_CONFLICT_FLAG

	var sectionStart uint32 = TRAINER_START_ADDR
	if info.trainer {
		info.data.trainer = rom[sectionStart : sectionStart+TRAINER_SIZE]
		sectionStart += TRAINER_SIZE
	} else {
		info.data.trainer = []byte{}
	}

	info.data.prgRom = rom[sectionStart : sectionStart+(PRG_ROM_SIZE*info.prgRomSize)]
	sectionStart += PRG_ROM_SIZE * info.prgRomSize
	info.data.chrRom = rom[sectionStart : sectionStart+(CHR_ROM_SIZE*info.chrRomSize)]
	sectionStart += CHR_ROM_SIZE * info.chrRomSize

	return nil
}
