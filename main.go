package main

import "./gnes/lib"
import "fmt"

func main() {
	err := gneslib.RunCLIDebugger("roms/cpu.nes")
	if err != nil {
		fmt.Println(err)
	}

}
