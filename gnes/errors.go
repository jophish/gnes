package gnes

import "fmt"

var err_BAD_INES_HEADER = &gError{1, "ROM does not have iNES flag set"}
var err_NONZERO_INES_HEADER_BUFFER = &gError{2, "iNES ROM has nonzero header buffer"}
var err_BAD_MAGIC_CONSTANT = &gError{3, "ROM does not contain valid iNES file format magic constant"}
var err_MAPPER_UNSUPPORTED = &gError{4, "Mapper is currently unsupported"}
var err_INCONSISTENT_PRG_ROM_SIZE = &gError{5, "Cartridge does not contain amount of PRG ROM specified in header"}
var err_INCONSISTENT_PRG_RAM_SIZE = &gError{6, "Cartridge specifies inconsistent amount of PRG RAM for mapper type"}
var err_ADDR_OUT_OF_BOUNDS = &gError{7, "Address out of bounds "}

type gError struct {
	errType uint64
	info    string
}

func (e *gError) Error() string {
	return fmt.Sprintf("%s", e.info)
}
