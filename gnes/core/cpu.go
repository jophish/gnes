package gnes

import "errors"
import "fmt"

const (
	NEGATIVE_MASK = 0x80
	BIT_6_MASK    = 0x40
)

const (
	mode_IM     = 0
	mode_ZERO   = 1
	mode_ZERO_X = 2
	mode_ABS    = 3
	mode_ABS_X  = 4
	mode_ABS_Y  = 5
	mode_IND_X  = 6
	mode_IND_Y  = 7
	mode_NA     = 8
	mode_NI     = 9 // not implemented
	mode_IMP    = 10
	mode_A      = 11
	mode_REL    = 12
)

const (
	loc_A      = 0
	loc_IM     = 1
	loc_ZERO   = 2
	loc_ZERO_X = 3
	loc_ABS    = 4
	loc_ABS_X  = 5
	loc_ABS_Y  = 6
	loc_IND_X  = 7
	loc_IND_Y  = 8
	loc_NA     = 9
	loc_NI     = 10
	loc_X      = 11
	loc_Y      = 12
	loc_REL    = 13
)

var opText = []string{
	// 0   1    2     3     4     5    6      7     8     9     a     b     c     d     e      f
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 0
	"BPL", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 1
	"JSR", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "BIT", "NA", "NA", "NA", // 2
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 3
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "LSR", "NA", "JMP", "NA", "NA", "NA", // 4
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 5
	"RTS", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 6
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "SEI", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // 7
	"NA", "NA", "NA", "NA", "NA", "STA", "STX", "NA", "DEY", "NA", "NA", "NA", "NA", "STA", "STX", "NA", // 8
	"NA", "STA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "TXS", "NA", "NA", "NA", "NA", "NA", // 9
	"NA", "NA", "LDX", "NA", "NA", "NA", "NA", "NA", "TAY", "LDA", "NA", "NA", "NA", "NA", "NA", "NA", // a
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // b
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "CMP", "DEX", "NA", "NA", "NA", "NA", "NA", // c
	"BNE", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "CLD", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // d
	"NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // e
	"NA", "NA", "NA", "NA", "NA", "NA", "INC", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", "NA", // f
}

var opArray = []func(*cpu) error{
	//  0     1         2         3         4         5         6         7         8         9         a         b         c         d         e         f
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 0
	(*cpu).op_BPL, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 1
	(*cpu).op_JSR, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_BIT, (*cpu).z, (*cpu).z, (*cpu).z, // 2
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 3
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_LSR, (*cpu).z, (*cpu).op_JMP, (*cpu).z, (*cpu).z, (*cpu).z, // 4
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 5
	(*cpu).op_RTS, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 6
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_SEI, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 7
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_ST, (*cpu).op_ST, (*cpu).z, (*cpu).op_DEC, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_ST, (*cpu).op_ST, (*cpu).z, // 8
	(*cpu).z, (*cpu).op_ST, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_TXS, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // 9
	(*cpu).z, (*cpu).z, (*cpu).op_LD, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_TAY, (*cpu).op_LD, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // a
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // b
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_CMP, (*cpu).op_DEC, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // c
	(*cpu).op_BNE, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_CLD, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // d
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, // e
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).op_INC, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z} // f

var opMode = []int{
	// 0     1        2        3        4        5        6        7        8        9        a        b        c        d        e        f
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 0
	mode_REL, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 1
	mode_ABS, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_ABS, mode_NI, mode_NI, mode_NI, // 2
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 3
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_A, mode_NI, mode_ABS, mode_NI, mode_NI, mode_NI, // 4
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 5
	mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 6
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 7
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_ZERO, mode_ZERO, mode_NI, mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_ABS, mode_ABS, mode_NI, // 8
	mode_NI, mode_IND_Y, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // 9
	mode_NI, mode_NI, mode_IM, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_IMP, mode_IM, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // a
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // b
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_IM, mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // c
	mode_REL, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_IMP, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // d
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // e
	mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_ZERO_X, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, mode_NI, // f
}

var opCycles = []uint64{
	// 1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0
	2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 1
	6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, // 2
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 3
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 3, 0, 0, 0, // 4
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 5
	6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 6
	0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // 7
	0, 0, 0, 0, 0, 3, 3, 0, 2, 0, 0, 0, 0, 4, 4, 0, // 8
	0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, // 9
	0, 0, 2, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, // a
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // b
	0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, // c
	2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // d
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // e
	0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, // f
}

var opSrc = []int{
	// 0     1      2        3      4        5       6        7       8        9      a        b       c       d       e       f
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 0
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 1
	loc_ABS, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_ABS, loc_NI, loc_NI, loc_NI, // 2
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 3
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_A, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, // 4
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 5
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 6
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 7
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_A, loc_X, loc_NI, loc_Y, loc_NI, loc_NI, loc_NI, loc_NI, loc_A, loc_X, loc_NI, // 8
	loc_NI, loc_A, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 9
	loc_NI, loc_NI, loc_IM, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_IM, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // a
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // b
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_IM, loc_X, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // c
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // d
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // e
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_ZERO_X, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // f
}

var opDst = []int{
	// 0     1        2        3        4      5      6        7      8       9        a      b        c      d       e       f
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 0
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 1
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, // 2
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 3
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_A, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, // 4
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 5
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 6
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 7
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_ZERO, loc_ZERO, loc_NI, loc_Y, loc_NI, loc_NI, loc_NI, loc_NI, loc_ABS, loc_ABS, loc_NI, // 8
	loc_NI, loc_IND_Y, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // 9
	loc_NI, loc_NI, loc_X, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_A, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // a
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // b
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_X, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // c
	loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NA, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // d
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // e
	loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_ZERO_X, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, loc_NI, // f
}

type registers struct {
	pc          uint16
	sp, x, y, a uint8
	// flags register
	n, v, b, d, i, z, c bool
}

type Registers struct {
	PC                  uint16
	SP, X, Y, A         uint8
	N, V, B, D, I, Z, C bool
}

type cpu struct {
	cycles uint64
	mmu    *mmu
	regs   *registers
}

func (cpu *cpu) getPC() uint16 {
	return cpu.regs.pc
}

func (cpu *cpu) setPC(addr uint16) {
	cpu.regs.pc = addr
}

func newRegs() *registers {
	return &registers{}
}

func (cpu *cpu) newPublicRegs() *Registers {
	regs := cpu.regs
	return &Registers{
		regs.pc, regs.sp, regs.x, regs.y, regs.a, regs.n,
		regs.v, regs.b, regs.d, regs.i, regs.z, regs.c,
	}
}

// We pass an mmu instance to the cpu instead of creating one here,
// since multiple NES subsystems need to access the same MMU instance.
func newCpu(mmu *mmu) (*cpu, error) {
	cpu := &cpu{}
	cpu.mmu = mmu
	cpu.regs = newRegs()
	err := cpu.initPC()
	if err != nil {
		return nil, err
	}
	return cpu, nil
}

func (cpu *cpu) initPC() error {
	pc_hi, err := cpu.mmu.read(vector_RESET_HI)
	if err != nil {
		return err
	}
	pc_lo, err := cpu.mmu.read(vector_RESET_LO)
	if err != nil {
		return err
	}
	cpu.regs.pc = (uint16(pc_hi) << 8) | uint16(pc_lo)
	return nil
}

// stepInstruction is the main method of progressing emulation. It
// fetches the instruction at the current PC and executes it accordingly,
// and returns the number of cycles taken.
func (cpu *cpu) stepInstruction() (uint64, error) {
	previousCycles := cpu.cycles
	addr := cpu.regs.pc
	op, err := cpu.mmu.read(addr)
	if err != nil {
		return 0, err
	}
	opDispatcher := opArray[op]
	err = opDispatcher(cpu)
	if err != nil {
		return 0, err
	}
	newCycles := cpu.cycles
	return newCycles - previousCycles, nil
}

// getOpMnemonic returns the dissasembly of the opcode at address addr.
func (cpu *cpu) getOpMnemonic(addr uint16) (string, error) {
	op, err := cpu.mmu.read(addr)
	if err != nil {
		return "", err
	}

	mnemonic := opText[op]
	if mnemonic == "NA" {
		return "", errors.New("Unknown opcode mnemonic")
	}

	mode := opMode[op]

	switch mode {
	case mode_IM:
		val, err := cpu.mmu.read(addr + 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s #%02x", mnemonic, val), nil
	case mode_ABS:
		val, err := cpu.mmu.read16(addr + 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %04x", mnemonic, val), nil
	case mode_A:
		return fmt.Sprintf("%s A", mnemonic), nil
	case mode_REL:
		val, err := cpu.mmu.read(addr + 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s *%d", mnemonic, int8(val)), nil
	case mode_IMP:
		return mnemonic, nil
	case mode_ZERO:
		val, err := cpu.mmu.read(addr + 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s $%04x", mnemonic, val), nil
	case mode_IND_Y:
		val, err := cpu.mmu.read(addr + 1)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s ($%02x), Y", mnemonic, val), nil

	default:
		return "", errors.New("Unknown op mode")
	}
}

// getSourceValue gets the source value for the opcode referenced by
// the current PC and returns it. For mode_ABS_X, mode_ABS_Y, and
// mode_IND_Y opcodes, returns 1 if a page was crossed while retrieving
// the value, and 0 in all other cases.
func (cpu *cpu) getSourceValue() (uint8, uint64, error) {
	addr := cpu.regs.pc
	op, err := cpu.mmu.read(addr)
	if err != nil {
		return 0, 0, err
	}

	var val uint8
	var pgCross uint64 = 1

	switch opSrc[op] {
	case loc_ABS:
		absAddr, err := cpu.mmu.read16(addr + 1)
		if err != nil {
			break
		}
		val, err = cpu.mmu.read(absAddr)
		if err != nil {
			break
		}
	case loc_IM:
		val, err = cpu.mmu.read(addr + 1)
	case loc_A:
		val = cpu.regs.a
	case loc_X:
		val = cpu.regs.x
	case loc_Y:
		val = cpu.regs.y
	case loc_ZERO_X:
		zeroPageAddr, err := cpu.mmu.read(addr + 1)
		if err != nil {
			break
		}
		finalAddr := uint16(zeroPageAddr + cpu.regs.x)

		val, err = cpu.mmu.read(finalAddr)
	default:
		err = errors.New("Invalid source location")
	}

	return val, pgCross, err
}

// getBranchAddress returns the address a current branch instruction
// should branch to if the branch test succeeds. Since all branch
// instructions use relative addressing, the relative offset will
// be interpreted as a two's complement signed value. If a page boundary
// would be crossed, the second return value is a 1, else 0.
func (cpu *cpu) getBranchAddress() (uint16, uint64, error) {
	addr := cpu.regs.pc
	val, err := cpu.mmu.read(addr + 1)
	if err != nil {
		return 0, 0, err
	}

	opLength, err := cpu.getPCOpLength()
	if err != nil {
		return 0, 0, err
	}

	newPC := uint16(int16(cpu.regs.pc+opLength) + int16(int8(val)))

	pageCrossed := 0
	if (newPC & 0xF0) != (addr & 0xF0) {
		pageCrossed = 1
	}

	return newPC, uint64(pageCrossed), nil
}

// getJumpAddress returns the address the current JMP instruction
// should jump to.
func (cpu *cpu) getJumpAddress() (uint16, error) {
	op, err := cpu.getCurrentOp()
	if err != nil {
		return 0, err
	}

	switch opMode[op] {
	case mode_ABS:
		imAddr, err := cpu.mmu.read16(cpu.regs.pc + 1)
		if err != nil {
			return 0, err
		}
		return imAddr, nil
	default:
		return 0, errors.New("Invalid op mode for JMP instruction")

	}
}

// writeToDestination writes val to the location specified by the
// current LD or ST opcode.
func (cpu *cpu) writeToDestination(val uint8) error {
	op, err := cpu.getCurrentOp()
	if err != nil {
		return err
	}

	var imAddr uint16

	switch opDst[op] {
	case loc_A:
		cpu.regs.a = val
	case loc_ABS:
		imAddr, err = cpu.mmu.read16(cpu.regs.pc + 1)
		if err != nil {
			break
		}
		err = cpu.mmu.write(val, imAddr)
		if err != nil {
			break
		}
	case loc_X:
		cpu.regs.x = val
	case loc_Y:
		cpu.regs.y = val
	case loc_ZERO:
		zeroLowByte, err := cpu.mmu.read(cpu.regs.pc + 1)
		if err != nil {
			break
		}
		zeroAddr := uint16(zeroLowByte)
		err = cpu.mmu.write(val, zeroAddr)
		if err != nil {
			break
		}
	case loc_IND_Y:
		zeroPageAddr, err := cpu.mmu.read(cpu.regs.pc + 1)
		if err != nil {
			break
		}
		baseIndirectAddr, err := cpu.mmu.read16(uint16(zeroPageAddr))
		if err != nil {
			break
		}
		finalAddr := baseIndirectAddr + uint16(cpu.regs.y)
		err = cpu.mmu.write(val, finalAddr)
		if err != nil {
			break
		}
	default:
		err = errors.New("Invalid destination location")
	}
	return err
}

// getCurrentOp returns the opcode at the current PC. This is
// easy to do, but it's a very common task, so we provide this
// method for convenience.
func (cpu *cpu) getCurrentOp() (uint8, error) {
	addr := cpu.regs.pc
	return cpu.mmu.read(addr)
}

// getOpCycles returns the base number of cycles associated
// with the current opcode.
func (cpu *cpu) getOpCycles() (uint64, error) {
	op, err := cpu.getCurrentOp()
	if err != nil {
		return 0, err
	}
	return opCycles[op], nil
}

// getOpLength returns the length of the opcode op.
func (cpu *cpu) getOpLength(op uint8) (uint16, error) {
	switch opMode[op] {
	case mode_IM:
		return 2, nil
	case mode_ABS:
		return 3, nil
	case mode_A:
		return 1, nil
	case mode_REL:
		return 2, nil
	case mode_IMP:
		return 1, nil
	case mode_ZERO:
		return 2, nil
	case mode_IND_Y:
		return 2, nil
	case mode_ZERO_X:
		return 2, nil
	default:
		return 0, errors.New("Invalid opcode length")
	}
}

// getPCOpLength returns the length of the opcode referenced
// by the current PC, in bytes.
func (cpu *cpu) getPCOpLength() (uint16, error) {
	op, err := cpu.getCurrentOp()
	if err != nil {
		return 0, err
	}
	return cpu.getOpLength(op)
}

// getOpLengthAddr returns the length of the opcode at address addr.
func (cpu *cpu) getOpLengthAddr(addr uint16) (uint16, error) {
	op, err := cpu.mmu.read(addr)
	if err != nil {
		return 0, err
	}
	return cpu.getOpLength(op)
}

// incrementPC increments the PC to the next opcode, i.e.,
// by the number of bytes in the current opcode.
func (cpu *cpu) incrementPC() error {
	opLength, err := cpu.getPCOpLength()
	if err != nil {
		return err
	}
	cpu.regs.pc += opLength
	return nil
}

// getSP transforms the internal stack pointer register into its actual address
func (cpu *cpu) getSP() uint16 {
	return uint16(cpu.regs.sp) + 0x100
}

// push pushes a byte onto the stack
func (cpu *cpu) push(val uint8) error {
	err := cpu.mmu.write(val, cpu.getSP())
	if err != nil {
		return err
	}
	cpu.regs.sp--
	return nil
}

// pop pops a byte off the stack
func (cpu *cpu) pop() (uint8, error) {
	val, err := cpu.mmu.read(cpu.getSP() + 1)
	if err != nil {
		return 0, err
	}
	cpu.regs.sp++
	return val, nil
}

// push16 pushes a 16 bit value onto the stack and adjusts the stack pointer.
func (cpu *cpu) push16(val uint16) error {
	if err := cpu.push(msb(val)); err != nil {
		return err
	}
	if err := cpu.push(lsb(val)); err != nil {
		return err
	}
	return nil
}

// pop16 pops a 16 bit value off the stack and adjusts the stack pointer.
func (cpu *cpu) pop16() (uint16, error) {
	lsByte, err := cpu.pop()
	if err != nil {
		return 0, nil
	}
	msByte, err := cpu.pop()
	if err != nil {
		return 0, nil
	}
	return make16BitValue(msByte, lsByte), nil
}

/***********************************************/
/*             Opcode Functions                */
/***********************************************/

// op_RTS is responsible for returning from subroutines
func (cpu *cpu) op_RTS() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	newPC, err := cpu.pop16()
	if err != nil {
		return err
	}

	cpu.regs.pc = newPC + 1

	cpu.cycles += cycles

	return nil
}

// op_JSR is responsible for jumping to a subroutine
func (cpu *cpu) op_JSR() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	newPC, err := cpu.getJumpAddress()
	if err != nil {
		return err
	}

	cpu.push16(cpu.regs.pc + 2)

	cpu.regs.pc = newPC

	cpu.cycles += cycles

	return nil
}

// op_BIT is responsible for all bit test operations
func (cpu *cpu) op_BIT() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	finalVal := cpu.regs.a & val

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.z = (finalVal == 0)
	cpu.regs.v = (val & BIT_6_MASK) != 0
	cpu.regs.n = (val & NEGATIVE_MASK) != 0

	return nil
}

// op_INC is responsible for all increment operations
func (cpu *cpu) op_INC() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	finalVal := val + 1

	err = cpu.writeToDestination(finalVal)
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.z = false
	if finalVal == 0 {
		cpu.regs.z = true
	}

	cpu.regs.n = false
	if finalVal&NEGATIVE_MASK != 0 {
		cpu.regs.n = true
	}

	return nil
}

// op_DEC is responsible for all decrement operations
func (cpu *cpu) op_DEC() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	finalVal := val - 1

	err = cpu.writeToDestination(finalVal)
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.z = false
	if finalVal == 0 {
		cpu.regs.z = true
	}

	cpu.regs.n = false
	if finalVal&NEGATIVE_MASK != 0 {
		cpu.regs.n = true
	}

	return nil
}

// op_TAY transfers the contents of the accumulator to the Y register
func (cpu *cpu) op_TAY() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.y = cpu.regs.a

	cpu.regs.z = false
	if cpu.regs.y == 0 {
		cpu.regs.z = true
	}

	cpu.regs.n = false
	if cpu.regs.y&NEGATIVE_MASK != 0 {
		cpu.regs.n = true
	}

	return nil
}

// op_TXS transfers the contents of X to the SP register
func (cpu *cpu) op_TXS() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.sp = cpu.regs.x
	return nil
}

// op_CLD clears the decimal mode flag
func (cpu *cpu) op_CLD() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.d = false

	return nil
}

// op_SEI sets the interrupt disable flag
func (cpu *cpu) op_SEI() error {
	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	cpu.regs.i = true

	return nil
}

// op_JMP is responsible for all jump operations
func (cpu *cpu) op_JMP() error {
	addr, err := cpu.getJumpAddress()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	cpu.cycles += cycles
	cpu.setPC(addr)

	return nil
}

// op_LD is responsible for all load operations.
func (cpu *cpu) op_LD() error {
	val, pgCross, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.writeToDestination(val)
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	if val == 0 {
		cpu.regs.z = true
	} else {
		cpu.regs.z = false
	}

	if (val & 0x80) != 0 {
		cpu.regs.n = true
	} else {
		cpu.regs.n = false
	}

	cpu.cycles += cycles + pgCross

	return nil
}

// op_ST is responsible for all store operations
func (cpu *cpu) op_ST() error {
	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.writeToDestination(val)
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	cpu.cycles += cycles

	return nil
}

// op_LSR is responsible for all logical shift right operations
func (cpu *cpu) op_LSR() error {
	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	newVal := val >> 1

	err = cpu.writeToDestination(newVal)
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	if newVal == 0 {
		cpu.regs.z = true
	} else {
		cpu.regs.z = false
	}

	if (val & 0x1) == 0x1 {
		cpu.regs.c = true
	} else {
		cpu.regs.c = false
	}

	cpu.regs.n = false

	cpu.cycles += cycles

	return nil
}

// op_CMP is responsible for all compare operations
func (cpu *cpu) op_CMP() error {
	val, _, err := cpu.getSourceValue()
	if err != nil {
		return err
	}

	acc := cpu.regs.a

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	err = cpu.incrementPC()
	if err != nil {
		return err
	}

	if acc == val {
		cpu.regs.z = true
	} else {
		cpu.regs.z = false
	}

	if acc >= val {
		cpu.regs.c = true
	} else {
		cpu.regs.c = false
	}

	if ((acc - val) & 0x80) != 0 {
		cpu.regs.n = true
	} else {
		cpu.regs.n = false
	}

	cpu.cycles += cycles

	return nil
}

// op_BNE is responsible for the branch not equal operation
func (cpu *cpu) op_BNE() error {
	newPC, pageCrossed, err := cpu.getBranchAddress()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	if cpu.regs.z == false {
		cpu.cycles += 1 + pageCrossed
		cpu.regs.pc = newPC
	} else {
		cpu.cycles += cycles
		err := cpu.incrementPC()
		if err != nil {
			return err
		}
	}
	return nil
}

// op_BPL is responsible for the branch if positive operation
func (cpu *cpu) op_BPL() error {
	newPC, pageCrossed, err := cpu.getBranchAddress()
	if err != nil {
		return err
	}

	cycles, err := cpu.getOpCycles()
	if err != nil {
		return err
	}

	if cpu.regs.n == false {
		cpu.cycles += 1 + pageCrossed
		cpu.regs.pc = newPC
	} else {
		cpu.cycles += cycles
		err := cpu.incrementPC()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cpu *cpu) z() error {
	op, err := cpu.getCurrentOp()
	if err != nil {
		return nil
	}
	return gError2New(err_UNSUPPORTED_OPCODE, uint64(op), uint64(cpu.regs.pc))
}
