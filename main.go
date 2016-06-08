package main

import (
	"fmt"
	"github.com/bomer/chip8/chip8"
)

var myChip8 chip8.Chip8

//Temporarily draw straight to terminal, replce with a OPEN GL draw later. Pref with goMobile package.
func drawGraphics() {
	//y loop, 32 scan lines,x 64 pixels in each scan line
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if myChip8.Gfx[(y*64)+x] == 0 {
				//Black pixel
				fmt.Printf("x")
			} else {
				// White pixel
				fmt.Printf("o")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
}

func main() {

	// Set up render system and register input callbacks
	// setupGraphics()
	// setupInput()

	// Initialize the Chip8 system and load the game into the memory
	fmt.Printf("Chip 8 Initalising...\n")
	myChip8.Init()
	// Doesnt exist yet
	myChip8.LoadGame("pong.c8")

	// fmt.Printf("Hello chip8.\n We will be using the memory range %d %d \n ", 0x000, 0xFFF)

	// Emulation loop
	for {
		// Emulate one cycle
		myChip8.EmulateCycle()

		// If the draw flag is set, update the screen
		if myChip8.Draw_flag {
			drawGraphics()
			myChip8.Draw_flag = false
		}

		// Store key press state (Press and Release)
		// myChip8.setKeys()
	}
}
