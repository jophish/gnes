package main

import "./gnes"
import "fmt"

func main() {
	nesEmu := &gnes.Emulator{}
	err := nesEmu.LoadRom("mario.nes")
	if err != nil {
		fmt.Println(err)
	}
}
