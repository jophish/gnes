package gnes

type cpuFlags struct {
	n, v, b, d, i, z, c bool
}

type registers struct {
	pc       uint16
	sp, x, y uint8
	f        *cpuFlags
}

type Cpu struct {
	nes_regs *registers
}
