package chip8

//Struct for the Main Chip8 System
import (
	"fmt"
	"io/ioutil"
)

type Chip8 struct {
	//Current Op Code
	opcode uint16

	//Ram for the whole system. 4x1024 byes available
	memory [4096]byte

	// V is for the CPU Registers. v0,v1... v15. Last one is a carry flag
	V [16]byte

	//Index register + program counter,
	index uint16
	pc    uint16

	/* System Memory Map
	0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	0x200-0xFFF - Program ROM and work RAM
	*/

	//GPU Buffer
	gfx [64 * 32]byte

	//Sound Variables
	delay_timer byte
	sound_timer byte

	//Stack Position
	stack [16]uint16
	sp    uint16

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
	// Decode Opcode
	// Execute Opcode

	// Update timers

}
