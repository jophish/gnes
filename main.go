package main

import "./gnes"
import "fmt"

func main() {
	var err error
	_, err = gnes.NewEmulator("roms/cpu.nes")
	if err != nil {
		fmt.Println(err)
	}
}
