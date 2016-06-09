package main

import (
	"fmt"
	"github.com/bomer/chip8/chip8"

	"math/rand"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	// "golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

var myChip8 chip8.Chip8

var screenData [64][32][3]byte
var SCREEN_WIDTH int
var SCREEN_HEIGHT int
var glctx gl.Context
var position gl.Attrib
var program gl.Program

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

func EmulationLoop() {

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

func SetupGraphics() {
	SCREEN_WIDTH = 640
	SCREEN_HEIGHT = 320
}

func main() {

	// Set up render system and register input callbacks
	// SetupGraphics()
	// setupInput()
	setupTexture()

	// Initialize the Chip8 system and load the game into the memory
	fmt.Printf("Chip 8 Initalising...\n")
	myChip8.Init()
	// Doesnt exist yet
	myChip8.LoadGame("pong.c8")

	go EmulationLoop()

	rand.Seed(time.Now().UnixNano())

	app.Main(func(a app.App) {

		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				a.Send(paint.Event{}) // keep animating
			case touch.Event:
				if down := e.Type == touch.TypeBegin; down || e.Type == touch.TypeEnd {
					// game.Press(down)
				}
			case key.Event:
				if e.Code != key.CodeSpacebar {
					break
				}
				if down := e.Direction == key.DirPress; down || e.Direction == key.DirRelease {
					// game.Press(down)
				}
			}
		}
	})
}

var (
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	// game      *Game
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	position = glctx.GetAttribLocation(program, "position")
	// game = NewGame()
	// scene = game.Scene(eng)
	triangleData := []byte{0.0, 4, 0.0, // top left
		0.0, 0.0, 0.0, // bottom left
		4, 0.0, 0.0, // bottom right
	}
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)
}

func onStop() {
	eng.Release()
	images.Release()
	// game = nil
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	// now := clock.Time(time.Since(startTime) * 60 / time.Second)
	// game.Update(now)
	// eng.Render(scene, now, sz)
	// render()

	// glctx.BindBuffer(gl.ARRAY_BUFFER, buf)

	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, 4, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	glctx.DisableVertexAttribArray(position)
}

// var triangleData [] byte
// triangleData := []byte { 0.0, 0.4, 0.0, // top left
// 	0.0, 0.0, 0.0, // bottom left
// 	0.4, 0.0, 0.0, // bottom right
// }

// OpenGL Crap
// Setup Texture
func setupTexture() {
	// Clear screen
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			screenData[y][x][0] = 0
			screenData[y][x][1] = 0
			screenData[y][x][2] = 0
		}
	}

	// // Create a texture

	// TexImage2D     (target Enum, level int, width, height int, format Enum, ty Enum, data []byte)
	// p := &screenData
	// var test []byte
	// test[1] = 0x01
	// test[2] = 0x01
	// test[3] = 0x01

	// test := make([]byte, 5, 5)
	// glctx.TexImage2D(gl.TEXTURE_2D, 0, SCREEN_WIDTH, SCREEN_HEIGHT, gl.RGB, gl.UNSIGNED_BYTE, (byte) * screenData))
	// glTexSubImage2D(GL_TEXTURE_2D, 0 ,0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, GL_RGB, GL_UNSIGNED_BYTE, (GLvoid*)screenData);
	// glctx.TexImage2D(gl.TEXTURE_2D, 0, SCREEN_WIDTH, SCREEN_HEIGHT, gl.RGB, gl.UNSIGNED_BYTE, test)

	// // Set up the texture
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP);
	// glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP);

	// // Enable textures
	// glEnable(GL_TEXTURE_2D);
}
