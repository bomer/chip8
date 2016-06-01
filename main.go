package main

import (
	"github.com/bomer/chip8/chip8"
)

var myChip8 chip8.Chip8

func main() {

	// Set up render system and register input callbacks
	// setupGraphics()
	// setupInput()

	// Initialize the Chip8 system and load the game into the memory
	myChip8.Init()
	// Doesnt exist yet
	myChip8.LoadGame("pong.c8")

	// fmt.Printf("Hello chip8.\n We will be using the memory range %d %d \n ", 0x000, 0xFFF)

	// Emulation loop
	for {
		// Emulate one cycle
		myChip8.EmulateCycle()

		// If the draw flag is set, update the screen
		// if(myChip8.drawFlag)
		// drawGraphics()

		// Store key press state (Press and Release)
		// myChip8.setKeys()
	}
}
