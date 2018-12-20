package gnes

import "fmt"

var err_BAD_INES_HEADER = &gError{1, "ROM does not have iNES flag set"}
var err_NONZERO_INES_HEADER_BUFFER = &gError{2, "iNES ROM has nonzero header buffer"}

type gError struct {
	errType uint64
	info    string
}

func (e *gError) Error() string {
	return fmt.Sprintf("%s", e.info)
}
