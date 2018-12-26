package main

import "./gnes"
import "fmt"

func main() {
	var err error
	emu, err := gnes.NewEmulator("roms/cpu.nes")
	if err != nil {
		fmt.Println(err)
	}
	err = emu.Step()
	if err != nil {
		fmt.Println(err)
	}

}
