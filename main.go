package main

import "fmt"

//Current Op Code
var opcode uint16

//Ram for the whole system. 4x1024 byes available
var memory [4096]byte

// V is for the CPU Registers. v0,v1... v15. Last one is a carry flag
var V [16]byte

//Index register + program counter,
var index uint16
var programcounter uint16

/* System Memory Map
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

func main() {
	index = 0xFFF
	fmt.Printf("Hello chip8.\n We will be using the memory range %d %d \n ", 0x000, 0xFFF)
}
