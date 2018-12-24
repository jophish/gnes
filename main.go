package main

import "./gnes"
import "fmt"

func main() {
	var err error
	nesEmu, err := gnes.NewEmulator()
	if err != nil {
		fmt.Println(err)
	}
	err = nesEmu.LoadRom("roms/cpu.nes")
	if err != nil {
		fmt.Println(err)
	}
}
