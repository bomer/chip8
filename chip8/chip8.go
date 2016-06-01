package chip8

//Struct for the Main Chip8 System
import "fmt"

type Chip8 struct {
	//Current Op Code
	opcode uint16

	//Ram for the whole system. 4x1024 byes available
	memory [4096]byte

	// V is for the CPU Registers. v0,v1... v15. Last one is a carry flag
	V [16]byte

	//Index register + program counter,
	index          uint16
	programcounter uint16

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

func (self *Chip8) Init() {
	fmt.Printf("Chip 8 Initalising...\n")
}
