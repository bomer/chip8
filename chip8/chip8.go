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
	memory [4096]byte
	// V is for the CPU Registers. v0,v1... v15. Last one is a carry flag
	V [16]byte

	//Index register + program counter,

	pc     uint16 //Program Counter
	opcode uint16 //Current Op Code
	index  uint16 //Index Register
	sp     uint16 //Stack Position

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

// Initialize registers and memory once
func (self *Chip8) Init() {
	fmt.Printf("Chip 8 Initalising...\n")

	self.pc = 0x200 // Program counter starts at 0x200, the space of memory after the interpreter
	self.opcode = 0 // Reset current opcode
	self.index = 0  // Reset index register
	self.sp = 0     // Reset stack pointer

	// Clear display

}

//Read file in curent dir into memory
func (self *Chip8) LoadGame(filename string) {
	rom, _ := ioutil.ReadFile(filename)
	rom_length := len(rom)
	if rom_length > 0 {
		fmt.Printf("Rom Length = %d\n", rom_length)
	}

	//If room to store ROM in RAM
	if (4096 - 512) > rom_length {
		for i := 0; i < rom_length; i++ {
			self.memory[i+512] = rom[i]
		}
	}

	fmt.Printf("Rom %s loaded into memory\n", filename)
}

//Tick to load next emulation cycle
func (self *Chip8) EmulateCycle() {
	// Fetch Opcode
	b1 := uint16(self.memory[self.pc])
	b2 := uint16(self.memory[self.pc+1])

	//Bitwise, add padding to end of first byte and append second byte to end
	self.opcode = (b1 << 8) | b2

	fmt.Printf("OpCode = %#02x\n", self.opcode)
	// Decode Opcode
	switch self.opcode & 0xF000 {
	case 0xA000: // ANNN: Sets I to the address NNN
		// Execute opcode
		self.index = self.opcode & 0x0FFF
		self.pc += 2
		break

	//1 to 7, jump, call and skip instructions
	case 0x1000: // 0x1NNN: Jumps to address NNN
		self.pc = self.opcode & 0x0FFF
		break
	case 0x2000: // 0x2NNN: Calls subroutine at NNN.
		self.stack[self.sp] = self.pc  // Store current address in stack
		self.sp++                      // Increment stack pointer
		self.pc = self.opcode & 0x0FFF // Set the program counter to the address at NNN
		break
	case 0x3000: // 0x3XNN: Skips the next instruction if VX equals NN
		if uint16(self.V[(self.opcode&0x0F00)>>8]) == self.opcode&0x00FF {
			self.pc += 4
		} else {
			self.pc += 2
		}
		break
	case 0x4000: // 0x4XNN: Skips the next instruction if VX doesn't equal NN.
		if uint16(self.V[(self.opcode&0x0F00)>>8]) != self.opcode&0x00FF {
			self.pc += 4
		} else {
			self.pc += 2
		}
	case 0x5000: // 0x5XY0: Skips the next instruction if VX equals VY.
		if uint16(self.V[(self.opcode&0x0F00)>>8]) != uint16(self.V[self.opcode&0x0F0]>>4) {
			self.pc += 4
		} else {
			self.pc += 2
		}

	default:
		fmt.Println("Unknown Opcode!")
		break

	}

	// Execute Opcode

	// Update timers

}
