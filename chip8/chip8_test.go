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
	fmt.Println("Running Unit Tests....")
	Prep()
	if myChip8.Pc != 512 {
		t.Error("Error -  Did not initalise correctly")
	}
}

//Clear Screen, put a byte in v memory on, then send memory to wipe memory.
func TestOpCode00E0(t *testing.T) {

	Prep()

	myChip8.Memory[512] = 0x00
	myChip8.Memory[513] = 0xE0

	//Check empty
	for i := 0; i < 64*32; i++ {
		if myChip8.Gfx[i] == 1 {
			t.Error("GFX Memory not initialised correctly")
		}
	}
	myChip8.Gfx[10] = 1
	found := false
	for i := 0; i < 64*32; i++ {
		if myChip8.Gfx[i] == 1 {
			found = true
		}
	}
	if !found {
		t.Error("GFX Memory did not have a byte turned on.")
	}
	if myChip8.Draw_flag != false {
		t.Error("Draw Flag not set when it should be")
	}
	myChip8.EmulateCycle()
	//Check empty
	for i := 0; i < 64*32; i++ {
		if myChip8.Gfx[i] == 1 {
			t.Error("GFX Memory has a byte on when all should be off")
		}
	}

	if myChip8.Draw_flag != true {
		t.Error("Draw Flag not set when it should be")
	}
}

// Return from Sub Routine, minus the pointer, move it back to where it was and then move the pointer + 2
func TestOpCode000E(t *testing.T) {
	Prep()

	myChip8.Sp = 2
	myChip8.Stack[1] = 0x210 //Tell it to goto 530

	//Set instructions to Reset Pointer,
	myChip8.Memory[512] = 0x00
	myChip8.Memory[513] = 0x0E

	myChip8.EmulateCycle()
	if myChip8.Sp != 1 {
		t.Error("Did not reduce the Stack pointer correctly")
	}
	if myChip8.Pc != 530 {
		t.Error("Did not reduce the Program Counter correctly")
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

func TestOpCode4XNN(t *testing.T) {
	Prep()
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct Stack Position")
	}

	// Success Case, does not match so  progress 4
	myChip8.V[0] = 0x02
	myChip8.Memory[512] = 0x40
	myChip8.Memory[513] = 0xff
	myChip8.EmulateCycle()

	if myChip8.Pc != 0x204 { //516
		t.Error("Did not Update the program counter correctly")
	}

	// Success Case, does not match so  progress 4
	myChip8.Pc = 512
	myChip8.Memory[513] = 0x02
	myChip8.EmulateCycle()

	// fmt.Printf("pc = %02x", myChip8.Pc)
	if myChip8.Pc != 0x202 { //514
		t.Error("Did not Update the program counter correctly")
	}
}

// 0x5XY0: Skips the next instruction if VX equals VY.
func TestOpCode5XY0(t *testing.T) {
	Prep()
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct program counter")
	}

	// Success Case, v0 and v1 are both the same memory.
	myChip8.V[0] = 0x02
	myChip8.V[1] = 0x02
	myChip8.Memory[512] = 0x50
	myChip8.Memory[513] = 0x11
	myChip8.EmulateCycle()

	if myChip8.Pc != 0x204 { //516
		t.Error("Did not Update the program counter correctly")
	}

	//Fail Case - Values are not equal
	myChip8.Pc = 512
	myChip8.V[0] = 0x02
	myChip8.V[1] = 0x03

	myChip8.EmulateCycle()

	if myChip8.Pc != 0x202 { //514
		t.Error("Did not Update the program counter correctly")
	}
}
