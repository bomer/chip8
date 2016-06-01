package main

import (
	"fmt"
	"github.com/bomer/chip8/chip8"
)

var system chip8.Chip8

func main() {

	system.Init()
	fmt.Printf("Hello chip8.\n We will be using the memory range %d %d \n ", 0x000, 0xFFF)
}
