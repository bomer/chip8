package chip8

//Struct for the Main Chip8 System
import (
	"fmt"
	"io/ioutil"
)

type Chip8 struct {

	/* System Memory Map
	0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	0x200-0xFFF - Program ROM and work RAM
	*/

	//Ram for the whole system. 4x1024 byes available
	Memory [4096]byte
	// V is for the CPU Registers. v0,v1... v15. Last one is a carry flag
	V [16]byte

	//Index register + program counter,

	Pc     uint16 //Program Counter
	Opcode uint16 //Current Op Code
	Index  uint16 //Index Register
	Sp     uint16 //Stack Pointer

	//GPU Buffer
	Gfx       [64 * 32]byte
	Draw_flag bool

	//Sound Variables
	delay_timer byte
	sound_timer byte

	//Stack
	Stack [16]uint16

	//Keyboard
	key [16]byte
}

// Initialize registers and Memory once
func (self *Chip8) Init() {
	self.Pc = 0x200 // Program counter starts at 0x200, the Space of Memory after the interpreter
	self.Opcode = 0 // Reset current Opcode
	self.Index = 0  // Reset index register
	self.Sp = 0     // Reset stack pointer

	for x := 0; x < 16; x++ {
		self.V[x] = 0
	}

	// Clear display
	for i := 0; i < 64*32; i++ {
		self.Gfx[i] = 0
	}

}

//Read file in curent dir into Memory
func (self *Chip8) LoadGame(filename string) {
	rom, _ := ioutil.ReadFile(filename)
	rom_length := len(rom)
	if rom_length > 0 {
		// fmt.Printf("Rom Length = %d\n", rom_length)
	}

	//If room to store ROM in RAM
	if (4096 - 512) > rom_length {
		for i := 0; i < rom_length; i++ {
			self.Memory[i+512] = rom[i]
		}
	}

	// fmt.Printf("Rom %s loaded into Memory\n", filename)
}

//Tick to load next emulation cycle
func (self *Chip8) EmulateCycle() {
	// Fetch Opcode
	b1 := uint16(self.Memory[self.Pc])
	b2 := uint16(self.Memory[self.Pc+1])

	//Bitwise, add padding to end of first byte and append second byte to end
	self.Opcode = (b1 << 8) | b2
	// Decode Opcode
	// fmt.Printf("Processing Op Code %02x\n", self.Opcode)

	// 0x00E0 and 0x000E We have to do first because Golang seems to truncate 0x0000 into 0x00
	switch self.Opcode {
	case 0xE0: // 0x00E0: Clears the screen
		for i := 0; i < 64*32; i++ {
			self.Gfx[i] = 0
		}
		self.Draw_flag = true
		self.Pc += 2
		break

	case 0x0E: // 0x00EE: Returns from subroutine
		self.Sp--                     // 16 levels of stack, decrease stack pointer to prevent overwrite
		self.Pc = self.Stack[self.Sp] // Put the stored return address from the stack back into the program counter
		self.Pc += 2                  // Don't forget to increase the program counter!
		break

	}

	switch self.Opcode & 0xF000 {

	case 0xE000: // 0x00E0: Clears the screen
		for i := 0; i < 64*32; i++ {
			self.Gfx[i] = 0
		}
		self.Draw_flag = true
		self.Pc += 2
		break
	case 0xA000: // ANNN: Sets I to the address NNN
		// Execute Opcode
		self.Index = self.Opcode & 0x0FFF
		self.Pc += 2
		break

	//1 to 7, jump, call and skip instructions
	case 0x1000: // 0x1NNN: Jumps to address NNN
		self.Pc = self.Opcode & 0x0FFF
		break
	case 0x2000: // 0x2NNN: Calls subroutine at NNN.
		self.Stack[self.Sp] = self.Pc  // Store current address in stack
		self.Sp++                      // Increment stack pointer
		self.Pc = self.Opcode & 0x0FFF // Set the program counter to the address at NNN
		break
	case 0x3000: // 0x3XNN: Skips the next instruction if VX equals NN
		if uint16(self.V[(self.Opcode&0x0F00)>>8]) == self.Opcode&0x00FF {
			self.Pc += 4
		} else {
			self.Pc += 2
		}
		break
	case 0x4000: // 0x4XNN: Skips the next instruction if VX doesn't equal NN.
		if uint16(self.V[(self.Opcode&0x0F00)>>8]) != self.Opcode&0x00FF {
			self.Pc += 4
		} else {
			self.Pc += 2
		}
		break
	case 0x5000: // 0x5XY0: Skips the next instruction if VX equals VY.
		x := (self.Opcode & 0x0F00) >> 8
		y := (self.Opcode & 0x00F0) >> 4
		// fmt.Printf("x = %02x and y= %02x", x, y)
		// fmt.Printf("V0 = %02x v1= %02x", self.V[x], self.V[y])

		if uint16(self.V[x]) == uint16(self.V[y]) {
			self.Pc += 4
		} else {
			self.Pc += 2
		}
		break
	case 0x6000: //6XNN	Sets VX to NN.
		x := (self.Opcode & 0x0F00) >> 8
		NN := byte(self.Opcode & 0x00FF)
		self.V[x] = NN
		self.Pc += 2
		break
	case 0x7000: //0x7XNN	Adds NN to VX.
		x := (self.Opcode & 0xF00) >> 8
		NN := byte(self.Opcode & 0x00FF)
		self.V[x] += NN
		self.Pc += 2
		break

		//0X8000 - 8 CASES
		/*
			8XY0	Sets VX to the value of VY.
			8XY1	Sets VX to VX or VY.
			8XY2	Sets VX to VX and VY.
			8XY3	Sets VX to VX xor VY.
			8XY4	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		*/

	case 0x8000:
		switch self.Opcode & 0x000F { //	8XY0 Sets VX to the value of VY.
		case 0x0000: // 0x8XY0: Sets VX to the value of VY
			self.V[self.Opcode&0x0F00>>8] = self.V[self.Opcode&0x00F0>>4]
			self.Pc += 2
			break
		case 0x0001: // 0x8XY0: Sets VX to the value of VY
			self.V[self.Opcode&0x0F00>>8] |= self.V[self.Opcode&0x00F0>>4]
			self.Pc += 2
			break

		case 0x0002: // 0x8XY0: Sets VX to VX and VY.
			self.V[self.Opcode&0x0F00>>8] &= self.V[self.Opcode&0x00F0>>4]
			self.Pc += 2
			break

		case 0x0003: // 0x8XY3:	Sets VX to VX xor VY.
			self.V[self.Opcode&0x0F00>>8] ^= self.V[self.Opcode&0x00F0>>4]
			self.Pc += 2
			break
		case 0x0004: // 0x8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			self.V[self.Opcode&0x0F00>>8] = self.V[self.Opcode&0x00F0>>4]
			self.Pc += 2
			break
		}
	default:
		if self.Opcode != 0xE0 && self.Opcode != 0x0E {
			fmt.Println("Unknown Opcode!")
		}
		break

	}

	// Update timers
	if self.delay_timer > 0 {
		self.delay_timer--
	}

	if self.sound_timer > 0 {
		if self.sound_timer == 1 {
			// fmt.Printf("BEEP!\n")

		}
		self.sound_timer--
	}
}
