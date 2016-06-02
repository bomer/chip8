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
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct Stack Position")
	}

	//Set instructions to, 1, Move to, and 0x0226, pos 552
	myChip8.Memory[512] = 0x22
	myChip8.Memory[513] = 0x20
	myChip8.EmulateCycle()

	if myChip8.Sp != 1 {
		t.Error("Did not start in the correct Stack Position")
	}
	if myChip8.Pc != 544 {
		t.Error("Did not move the Program counter correctly")
	}
}

func TestOpCode3XNN(t *testing.T) {
	Prep()
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct Stack Position")
	}

	// Fail Case, does not match so only progress 2, stores 0 514
	myChip8.V[0] = 0x02
	//Set instructions to, 1, Move to, and 0x0226, pos 552
	myChip8.Memory[512] = 0x30
	myChip8.Memory[513] = 0x0f
	myChip8.EmulateCycle()

	if myChip8.Pc != 0x202 { //514
		t.Error("Did not Update the Stack Position correctly")
	}

	myChip8.Pc = 512
	myChip8.Memory[512] = 0x30
	myChip8.Memory[513] = 0x02
	myChip8.EmulateCycle()
	if myChip8.Pc != 0x0204 { //516
		t.Error("Did not Update the Stack Position correctly")
	}
}
