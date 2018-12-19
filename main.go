package main

import "./gnes"
import "fmt"

func main() {
	nesEmu := gnes.NewEmulator()
	err := nesEmu.LoadRom("roms/mario.nes")
	if err != nil {
		fmt.Println(err)
	}
}
