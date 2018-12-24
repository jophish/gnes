package gnes

type registers struct {
	pc       uint16
	sp, x, y uint8
	// flags register
	n, v, b, d, i, z, c bool
}

type cpu struct {
	mmu  *mmu
	regs *registers
}

func newRegs() *registers {
	return &registers{}
}

// We pass an mmu instance to the cpu instead of creating one here,
// since multiple NES subsystems need to access the same MMU instance.
func newCpu(mmu *mmu) *cpu {
	cpu := &cpu{}
	cpu.mmu = mmu
	cpu.regs = newRegs()
	return cpu
}
