package chip8_test

import (
	"fmt"
	"github.com/bomer/chip8/chip8"
	"testing"
)

var myChip8 chip8.Chip8

func Prep() {
	myChip8.Init()
	myChip8.LoadGame("../pong.c8")
}

func TestInit(t *testing.T) {
	Prep()
	// fmt.Println("mem = %s", myChip8.Memory)
	fmt.Println("Starting pc  = %s", myChip8.Pc)
	if myChip8.Pc != 512 {
		t.Error("Error -  Did not initalise correctly")
	}
}

func TestOpCode1NNN(t *testing.T) {

	Prep()
	//Set instructions to, 1, Move to, and 0x0226, pos 552
	myChip8.Memory[512] = 0x12
	myChip8.Memory[513] = 0x26
	myChip8.EmulateCycle()

	// fmt.Println("PC=== %s", myChip8.Pc)
	if myChip8.Pc != 550 {
		t.Error("Did not start in the correct Program counter")
	}
	myChip8.EmulateCycle()

	if myChip8.Pc != 552 {
		t.Error("Did not move the Program counter correctly")
	}

}

func TestOpCode2NNN(t *testing.T) {

	Prep()
	//Set instructions to, 1, Move to, and 0x0226, pos 552
	myChip8.Memory[512] = 0x22
	myChip8.Memory[513] = 0x26
	myChip8.EmulateCycle()

	fmt.Println("SP=== %s", myChip8.Sp)
	// if myChip8.Pc != 550 {
		t.Error("Did not start in the correct Program counter")
	}
	// myChip8.EmulateCycle()

	// if myChip8.Pc != 552 {
		// t.Error("Did not move the Program counter correctly")
	// }

}
