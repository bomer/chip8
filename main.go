package main

import (
	"fmt"
	"github.com/bomer/chip8/chip8"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"log"
	"math"
	// "math/rand"
	"os"
	"runtime"
	"time"
)

const (
	BALL_RADIUS = 25
	MULTIPLIER  = 10
)

var myChip8 chip8.Chip8

//Temporarily draw straight to terminal, replce with a OPEN GL draw later. Pref with goMobile package.
func drawGraphics() {
	fmt.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
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

}

// key events are a way to get input from GLFW.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//if u only want the on press, do = && && action == glfw.Press
	var keydown byte
	keydown = 0

	if action == glfw.Press {
		keydown = 1
	}
	if key == glfw.Key1 {
		myChip8.Key[0x1] = keydown
	} else if key == glfw.Key2 {
		myChip8.Key[0x2] = keydown
	} else if key == glfw.Key3 {
		myChip8.Key[0x3] = keydown
	} else if key == glfw.Key4 {
		myChip8.Key[0xC] = keydown
	} else if key == glfw.KeyQ {
		myChip8.Key[0x4] = keydown
	} else if key == glfw.KeyW {
		myChip8.Key[0x5] = keydown
	} else if key == glfw.KeyE {
		myChip8.Key[0x6] = keydown
	} else if key == glfw.KeyR {
		myChip8.Key[0xD] = keydown
	} else if key == glfw.KeyA {
		myChip8.Key[0x7] = keydown
	} else if key == glfw.KeyS {
		myChip8.Key[0x8] = keydown
	} else if key == glfw.KeyD {
		myChip8.Key[0x9] = keydown
	} else if key == glfw.KeyF {
		myChip8.Key[0xE] = keydown
	} else if key == glfw.KeyZ {
		myChip8.Key[0xA] = keydown
	} else if key == glfw.KeyX {
		myChip8.Key[0x0] = keydown
	} else if key == glfw.KeyC {
		myChip8.Key[0xB] = keydown
	} else if key == glfw.KeyV {
		myChip8.Key[0xF] = keydown
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// drawCircle draws a circle for the specified radius, rotation angle, and the specified number of sides
func drawCircle(radius float64, sides int) {
	gl.Begin(gl.TRIANGLE_FAN)
	for a := 0.0; a < 2*math.Pi; a += (2 * math.Pi / float64(70)) {
		gl.Vertex2d(math.Sin(a)*radius, math.Cos(a)*radius)
	}
	gl.Vertex3f(0, 0, 0)
	gl.End()

}

// Old gfx code
func drawPixel(x int, y int) {
	gl.Begin(gl.QUADS)
	gl.Vertex3f(float32(x*MULTIPLIER)+0.0, float32(y*MULTIPLIER)+0.0, 0.0)
	gl.Vertex3f(float32(x*MULTIPLIER)+0.0, float32(y*MULTIPLIER)+MULTIPLIER, 0.0)
	gl.Vertex3f(float32(x*MULTIPLIER)+MULTIPLIER, float32(y*MULTIPLIER)+MULTIPLIER, 0.0)
	gl.Vertex3f(float32(x*MULTIPLIER)+MULTIPLIER, float32(y*MULTIPLIER)+0.0, 0.0)
	gl.End()
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	//Transform screen to keep player in middle. Added intentation to make obvious the push matrix is like a block
	gl.PushMatrix()

	// Draw
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if myChip8.Gfx[((31-y)*64)+x] == 0 {
				gl.Color3f(0.0, 0.0, 0.0)
			} else {
				gl.Color3f(1.0, 1.0, 1.0)
			}
			drawPixel(x, y)
		}

	}

	//Second Pop
	gl.PopMatrix()
}

// onResize sets up a simple 2d ortho context based on the window size
func onResize(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(0, 0, 0, 0)
}

func main() {
	runtime.LockOSThread()

	myChip8.Init()
	// Doesnt exist yet
	myChip8.LoadGame("pong.c8")

	// fmt.Printf("Hello chip8.\n We will be using the memory range %d %d \n ", 0x000, 0xFFF)

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize GLFW: ", err)
	}
	defer glfw.Terminate()

	// create window
	window, err := glfw.CreateWindow(64*MULTIPLIER, 32*MULTIPLIER, os.Args[0], nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.SetFramebufferSizeCallback(onResize)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// set up opengl context
	onResize(window, 64*MULTIPLIER, 32*MULTIPLIER)

	// glfw.KeyCallback(window)
	window.SetKeyCallback(keyCallback)

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	//Run emulator on another go-routine
	//Else emulator runs to slow on main thread.
	go func() {
		emuticker := time.NewTicker(time.Second / 360)
		for {
			myChip8.EmulateCycle()
			<-emuticker.C
		}
	}()

	ticker := time.NewTicker(time.Second / 30)
	for !window.ShouldClose() {
		// myChip8.EmulateCycle()
		if myChip8.Draw_flag {
			// drawGraphics() //for debugging
			draw()
			myChip8.Draw_flag = false
		}

		// handleMouse(window)
		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}
