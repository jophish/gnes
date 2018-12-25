package main

import "./gnes"
import "fmt"

func main() {
	var err error
	emu, err := gnes.NewEmulator("roms/cpu.nes")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 10; i++ {
		addr := uint16(0xffd8 + i)
		val, err := emu.ReadCpu(addr)
		if err != nil {
			fmt.Printf("%#x \n", addr)
			fmt.Println(err)
		} else {
			fmt.Printf("%#x: %#x\n", addr, val)
		}

	}

}
