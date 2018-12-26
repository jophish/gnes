package gnes

const (
	mode_IM     = 0
	mode_ZERO   = 1
	mode_ZERO_X = 2
	mode_ABS    = 3
	mode_ABS_X  = 4
	mode_ABS_Y  = 5
)

var opArray = []func(*cpu) error{
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z,
	(*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z, (*cpu).z}

var opMode = []int{}

type registers struct {
	pc       uint16
	sp, x, y uint8
	// flags register
	n, v, b, d, i, z, c bool
}

type cpu struct {
	cycles uint64
	mmu    *mmu
	regs   *registers
}

func newRegs() *registers {
	return &registers{}
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
// fetches the instruction at the current PC and executes it accordingly.
func (cpu *cpu) stepInstruction() error {
	addr := cpu.regs.pc
	op, err := cpu.mmu.read(addr)
	if err != nil {
		return err
	}
	opDispatcher := opArray[op]
	err = opDispatcher(cpu)
	if err != nil {
		return &gError2{err_UNSUPPORTED_OPCODE, uint64(op), uint64(addr)}
	}
	return nil
}

func (cpu *cpu) op_BRK() error {
	return nil
}

func (cpu *cpu) z() error {
	return &gError{err_UNIMPLEMENTED_OPCODE}
}
