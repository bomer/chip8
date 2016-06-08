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

//0x6XNN: Sets VX to NN.
func TestOpCode6XNN(t *testing.T) {
	Prep()

	//ensure valid starting state
	if myChip8.V[1] != 0 { //514
		t.Error("Did not init CPU Registers correctly")
	}
	// fmt.Printf("v1 = %X", myChip8.V[1])
	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x61
	myChip8.Memory[513] = 0xaa
	myChip8.EmulateCycle()

	// fmt.Printf("v1 = %X", myChip8.V[1])

	if myChip8.V[1] != 0xaa { //514
		t.Error("Did not set CPU Register V1 correctly")
	}
	if myChip8.Pc != 0x202 {
		t.Error("Did not reduce the Program Counter correctly")
	}

}

//0x7XNN	Adds NN to VX.
func TestOpCode7XNN(t *testing.T) {
	Prep()

	myChip8.V[1] = 0xaa
	//ensure valid starting state 170 + 1
	if myChip8.V[1] != 0xaa { //514
		t.Error("Did not init CPU Registers correctly")
	}
	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x71
	myChip8.Memory[513] = 0x01
	myChip8.EmulateCycle()

	if myChip8.V[1] != 0xab { //171
		t.Error("Did not set CPU Register V1 correctly")
	}
	// Second succes state  - ab + 30 = db
	myChip8.Memory[514] = 0x71
	myChip8.Memory[515] = 0x30 //
	myChip8.EmulateCycle()

	// fmt.Printf("v1 = %X\n", myChip8.V[1])

	if myChip8.V[1] != 0xdb { //
		t.Error("Did not set CPU Register V1 correctly")
	}
}

//0x8XY0 Sets VX to the value of VY.
func TestOpCode8XY0(t *testing.T) {
	Prep()

	//Postive case, test settings v0 to v1, 1 to 2.
	myChip8.V[0] = 0x01
	myChip8.V[1] = 0x02

	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x10
	myChip8.EmulateCycle()

	if myChip8.V[0] != 0x02 {
		t.Error("Failed to set VY to VX")
	}
}

//0x8XY1  Sets VX to VX or VY.
func TestOpCode8XY1(t *testing.T) {
	Prep()

	//Postive case, test settings v0 to v1, 1 to 2.
	myChip8.V[0] = 0x01
	myChip8.V[1] = 0x05

	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x11
	myChip8.EmulateCycle()

	if myChip8.V[0] != 0x05 {
		t.Error("Failed to set VY to VX")
	}
}

//8XY2	Sets VX to VX and VY.

func TestOpCode8XY2(t *testing.T) {
	Prep()

	//Postive case, test settings v0 to v1, 1 to 2.
	myChip8.V[0] = 0x01
	myChip8.V[1] = 0x05

	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x12
	myChip8.EmulateCycle()

	if myChip8.V[0] != 0x01 {
		t.Error("Failed to set VY to VX")
	}
}

//8XY2	Sets VX to VX xor VY.

func TestOpCode8XY3(t *testing.T) {
	Prep()

	//Postive case, test settings v0 to v1, 1 to 2.
	myChip8.V[0] = 0x01
	myChip8.V[1] = 0x05

	// Success Case, v1 and v1 are both the same memory.
	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x13
	myChip8.EmulateCycle()

	if myChip8.V[0] != 0x04 {
		t.Error("Failed to set VY to VX")
	}
}

// 0x8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
// carr if if second mumber is greater than 255 - first
func TestOpCode8XY4(t *testing.T) {
	Prep()

	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x14
	// positive eg, 60 + 250, is 250 > (255 - 65 = 200) , YES = carry
	myChip8.V[0] = 60
	myChip8.V[1] = 200
	myChip8.EmulateCycle()
	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 0x04 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}

	//negative, 54 + 200, is 200 > (255-50=205), no
	myChip8.Pc = 512
	myChip8.V[0] = 54
	myChip8.V[1] = 200
	myChip8.EmulateCycle()
	if myChip8.V[0] != 0xfe {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 0 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}
	//Max, 255 + 255, is 255 > (255-255=0), yes
	myChip8.Pc = 512
	myChip8.V[0] = 0xff
	myChip8.V[1] = 0xff
	myChip8.EmulateCycle()
	if myChip8.V[0] != 0xfe {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}
}

// 0x8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
// carr if if second mumber is greater than the first
func TestOpCode8XY5(t *testing.T) {
	Prep()

	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x15
	// No Carry Case, 200-100
	myChip8.V[0] = 200
	myChip8.V[1] = 100
	myChip8.EmulateCycle()
	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 100 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}

	// Carry Case, 100-200
	myChip8.Pc = 512
	myChip8.V[0] = 100
	myChip8.V[1] = 200
	myChip8.EmulateCycle()
	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 156 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 0 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}
}

// 8XY6 Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.[2]
func TestOpCode8XY6(t *testing.T) {
	Prep()

	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x06
	// Simple case, 0x05
	//  5    to   2
	// 0101      0010 [1] < v[15]
	myChip8.V[0] = 0x5

	myChip8.EmulateCycle()
	if myChip8.V[0] != 0x02 {
		t.Error("Failed to bitshift VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to bitshift VX by 1 and get the correct least important bit flag")
	}

	// Harder case, 0x1101
	//  14  to   6
	// 1101     0110 [1] < v[15]
	myChip8.Pc = 512
	myChip8.V[0] = 0xD

	myChip8.EmulateCycle()
	if myChip8.V[0] != 0x06 {
		t.Error("Failed to bitshift VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to bitshift VX by 1 and get the correct least important bit flag")
	}
}

//8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
// carr if if second mumber is greater than the first
func TestOpCode8XY7(t *testing.T) {
	Prep()

	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x17
	// No Carry Case, vx=vy-vx, vx=200-100
	myChip8.V[0] = 100
	myChip8.V[1] = 200
	myChip8.EmulateCycle()
	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 100 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}

	// Carry Case, 100-200
	myChip8.Pc = 512
	myChip8.V[0] = 200
	myChip8.V[1] = 100
	myChip8.EmulateCycle()
	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 156 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.V[0xf] != 0 {
		t.Error("Failed to add VY to VX")
	}
	if myChip8.Pc != 514 {
		t.Error("Failed to update program counter")
	}
}

// 8XY6 Shifts VX Left by one. VF is set to the value of the least significant bit of VX before the shift.[2]
func TestOpCode8XYE(t *testing.T) {
	Prep()

	myChip8.Memory[512] = 0x80
	myChip8.Memory[513] = 0x0E
	// Simple case, 05
	//  5    to   10
	// 0101      1010 [1] < v[15]
	myChip8.V[0] = 0x5

	myChip8.EmulateCycle()

	if myChip8.V[0] != 0x0a {
		t.Error("Failed to bitshift VX")
	}
	if myChip8.V[0xf] != 0 {
		t.Error("Failed to bitshift VX by 1 and get the correct least important bit flag")
	}

	// Harder case, 1000 1101, testing the leading bit gets shifts and recored in v15
	//  142  to   26
	// 1000 1101      11010 [1] < v[15]
	myChip8.Pc = 512
	myChip8.V[0] = 142

	myChip8.EmulateCycle()

	// fmt.Printf("v0 = %d", myChip8.V[0])
	if myChip8.V[0] != 28 {
		t.Error("Failed to bitshift VX")
	}
	if myChip8.V[0xf] != 1 {
		t.Error("Failed to bitshift VX by 1 and get the correct least important bit flag")
	}

}

// 0x5XY0: Skips the next instruction if VX equals VY.
func TestOpCode9XY0(t *testing.T) {
	Prep()
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct program counter")
	}

	// Success Case, v0 and v1 are both the same memory.
	myChip8.V[0] = 0x02
	myChip8.V[1] = 0x02
	myChip8.Memory[512] = 0x90
	myChip8.Memory[513] = 0x11
	myChip8.EmulateCycle()

	if myChip8.Pc != 0x202 { //516
		t.Error("Did not Update the program counter correctly")
	}

	//Fail Case - Values are not equal
	myChip8.Pc = 512
	myChip8.V[0] = 0x02
	myChip8.V[1] = 0x03

	myChip8.EmulateCycle()

	if myChip8.Pc != 0x204 { //514
		t.Error("Did not Update the program counter correctly")
	}
}

////ANNN	Sets I to the address NNN.
func TestOpCodeANNN(t *testing.T) {
	Prep()
	if myChip8.Sp != 0 {
		t.Error("Did not start in the correct program counter")
	}

	// Success Case
	myChip8.Memory[512] = 0xA1
	myChip8.Memory[513] = 0x11
	myChip8.Index = 0x132
	myChip8.EmulateCycle()

	if myChip8.Index != 0x111 { //516
		t.Error("Did not Update the Index correctly")
	}
	if myChip8.Pc != 514 {
		t.Error("Did not Update the program counter correctly")
	}
}
