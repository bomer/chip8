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
	index  uint16 //Index Register
	Sp     uint16 //Stack Position

	//GPU Buffer
	Gfx [64 * 32]byte

	//Sound Variables
	delay_timer byte
	sound_timer byte

	//Stack Position
	stack [16]uint16

	//Keyboard
	key [16]byte
}

// Initialize registers and Memory once
func (self *Chip8) Init() {
	fmt.Printf("Chip 8 Initalising...\n")

	self.Pc = 0x200 // Program counter starts at 0x200, the Space of Memory after the interpreter
	self.Opcode = 0 // Reset current Opcode
	self.index = 0  // Reset index register
	self.Sp = 0     // Reset stack pointer

	for x := 0; x < 16; x++ {
		self.V[x] = 0
	}

	// Clear diSplay

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

	fmt.Printf("OpCode = %#02x\n", self.Opcode)
	// Decode Opcode
	switch self.Opcode & 0xF000 {
	case 0xA000: // ANNN: Sets I to the address NNN
		// Execute Opcode
		self.index = self.Opcode & 0x0FFF
		self.Pc += 2
		break

	//1 to 7, jump, call and skip instructions
	case 0x1000: // 0x1NNN: Jumps to address NNN
		self.Pc = self.Opcode & 0x0FFF
		break
	case 0x2000: // 0x2NNN: Calls subroutine at NNN.
		self.stack[self.Sp] = self.Pc  // Store current address in stack
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

	default:
		fmt.Println("Unknown Opcode!")
		break

	}

	// Execute Opcode

	// Update timers

}
