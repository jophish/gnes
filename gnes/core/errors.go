package gnes

import "fmt"

const (
	err_BAD_INES_HEADER               = 1
	err_NONZERO_INES_HEADER_BUFFER    = 2
	err_BAD_MAGIC_CONSTANT            = 3
	err_MAPPER_UNSUPPORTED            = 4
	err_INCONSISTENT_PRG_ROM_SIZE     = 5
	err_INCONSISTENT_PRG_RAM_SIZE     = 6
	err_ADDR_OUT_OF_BOUNDS            = 7
	err_UNSUPPORTED_OPCODE            = 8
	err_UNIMPLEMENTED_OPCODE          = 9
	err_MMC1_PRG_RAM_DISABLED         = 10
	err_MMC1_INVALID_PRG_ROM_MODE_VAL = 11
	err_UNWRITEABLE_PPU_REG           = 12
	err_UNREADABLE_PPU_REG            = 13
)

var errToString = map[int]string{
	err_BAD_INES_HEADER:               "ROM does not have iNES flag set",
	err_NONZERO_INES_HEADER_BUFFER:    "iNES ROM has nonzero header buffer",
	err_BAD_MAGIC_CONSTANT:            "ROM does not contain valid iNES file format magic constant",
	err_MAPPER_UNSUPPORTED:            "Mapper %d is currently unsupported",
	err_INCONSISTENT_PRG_ROM_SIZE:     "Cartridge does not contain amount of PRG ROM specified in header",
	err_INCONSISTENT_PRG_RAM_SIZE:     "Cartridge specifies inconsistent amount of PRG RAM for mapper type",
	err_ADDR_OUT_OF_BOUNDS:            "Address out of bounds",
	err_UNSUPPORTED_OPCODE:            "Unsupported opcode %x at address %#x",
	err_UNIMPLEMENTED_OPCODE:          "Opcode unimplemented",
	err_MMC1_PRG_RAM_DISABLED:         "PRG RAM is disabled on this MMC1 mapper",
	err_MMC1_INVALID_PRG_ROM_MODE_VAL: "Invalid value %x for MMC1 PRG ROM mode",
	err_UNWRITEABLE_PPU_REG:           "Illegal PPU register to write",
	err_UNREADABLE_PPU_REG:            "Illegal PPU register to read",
}

type gError struct {
	errType int
}

func gErrorNew(errType int) *gError {
	return &gError{errType}
}

func (e *gError) Error() string {
	return fmt.Sprintf(errToString[e.errType])
}

type gError1 struct {
	errType int
	val1    uint64
}

func gError1New(errType int, val1 uint64) *gError1 {
	return &gError1{errType, val1}
}

func (e *gError1) Error() string {
	return fmt.Sprintf(errToString[e.errType], e.val1)
}

type gError2 struct {
	errType int
	val1,
	val2 uint64
}

func gError2New(errType int, val1, val2 uint64) *gError2 {
	return &gError2{errType, val1, val2}
}

func (e *gError2) Error() string {
	return fmt.Sprintf(errToString[e.errType], e.val1, e.val2)
}
