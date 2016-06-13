package main

import (
	"encoding/binary"
	"fmt"
	"github.com/bomer/chip8/chip8"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"image"
	"log"

	"os"

	"time"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
	img    glutil.Image
)

var myChip8 chip8.Chip8

func main() {
	myChip8.Init()
	// Doesnt exist yet
	myChip8.LoadGame("pong.c8")

	//Run emulator on another go-routine
	//Else emulator runs to slow on main thread.
	go func() {
		emuticker := time.NewTicker(time.Second / 360)
		for {
			myChip8.EmulateCycle()
			<-emuticker.C
		}
	}()

	app.Main(func(a app.App) {
		var glctx gl.Context
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
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}
				if myChip8.Draw_flag {
					// drawGraphics() //for debugging

					onPaint(glctx, sz)
					myChip8.Draw_flag = false
				}

				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case key.Event:
				if e.Code == key.CodeEscape {
					os.Exit(0)
					break
				}
				//Input for emu
				var keydown byte
				keydown = 0

				if e.Direction == key.DirPress {
					keydown = 1
				}
				if e.Code == key.Code1 {
					myChip8.Key[0x1] = keydown
				} else if e.Code == key.Code2 {
					myChip8.Key[0x2] = keydown
				} else if e.Code == key.Code3 {
					myChip8.Key[0x3] = keydown
				} else if e.Code == key.Code4 {
					myChip8.Key[0xC] = keydown
				} else if e.Code == key.CodeQ {
					myChip8.Key[0x4] = keydown
				} else if e.Code == key.CodeW {
					myChip8.Key[0x5] = keydown
				} else if e.Code == key.CodeE {
					myChip8.Key[0x6] = keydown
				} else if e.Code == key.CodeR {
					myChip8.Key[0xD] = keydown
				} else if e.Code == key.CodeA {
					myChip8.Key[0x7] = keydown
				} else if e.Code == key.CodeS {
					myChip8.Key[0x8] = keydown
				} else if e.Code == key.CodeD {
					myChip8.Key[0x9] = keydown
				} else if e.Code == key.CodeF {
					myChip8.Key[0xE] = keydown
				} else if e.Code == key.CodeZ {
					myChip8.Key[0xA] = keydown
				} else if e.Code == key.CodeX {
					myChip8.Key[0x0] = keydown
				} else if e.Code == key.CodeC {
					myChip8.Key[0xB] = keydown
				} else if e.Code == key.CodeV {
					myChip8.Key[0xF] = keydown
				}

			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)

	//Draw Buffer
	img = *images.NewImage(64, 32)

}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)

	//Draw Pixels onto screen
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			if myChip8.Gfx[(j*64)+i] == 0 {
				img.RGBA.Set(i, j, image.Black)
			}

		}
	}

	//Draw over whole screen
	tl := geom.Point{0, 0}
	tr := geom.Point{geom.Pt(sz.WidthPx / 4), 0}
	bl := geom.Point{0, geom.Pt(sz.HeightPx / 4)}
	img.Upload()

	// Set up the texture
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	img.Draw(sz, tl, tr, bl, img.RGBA.Bounds())
	fps.Draw(sz)

	//cleanup every  frame
	img.Release()
	img = *images.NewImage(64, 32)

}

const squareoffset = 0.057

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, squareoffset, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	squareoffset, 0.0, 0.0, // bottom right
	squareoffset, squareoffset, 0.0,
)

const (
	coordsPerVertex = 3
	vertexCount     = 4
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`

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

// // key events are a way to get input from GLFW.
// func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	//if u only want the on press, do = && && action == glfw.Press
// 	var keydown byte
// 	keydown = 0

// 	if key == glfw.KeyEscape && action == glfw.Press {
// 		w.SetShouldClose(true)
// 	}
// }

// // drawCircle draws a circle for the specified radius, rotation angle, and the specified number of sides
// func drawCircle(radius float64, sides int) {
// 	gl.Begin(gl.TRIANGLE_FAN)
// 	for a := 0.0; a < 2*math.Pi; a += (2 * math.Pi / float64(70)) {
// 		gl.Vertex2d(math.Sin(a)*radius, math.Cos(a)*radius)
// 	}
// 	gl.Vertex3f(0, 0, 0)
// 	gl.End()

// }

// // Old gfx code
// func drawPixel(x int, y int) {
// 	gl.Begin(gl.QUADS)
// 	gl.Vertex3f(float32(x*MULTIPLIER)+0.0, float32(y*MULTIPLIER)+0.0, 0.0)
// 	gl.Vertex3f(float32(x*MULTIPLIER)+0.0, float32(y*MULTIPLIER)+MULTIPLIER, 0.0)
// 	gl.Vertex3f(float32(x*MULTIPLIER)+MULTIPLIER, float32(y*MULTIPLIER)+MULTIPLIER, 0.0)
// 	gl.Vertex3f(float32(x*MULTIPLIER)+MULTIPLIER, float32(y*MULTIPLIER)+0.0, 0.0)
// 	gl.End()
// }

// func draw() {
// 	gl.Clear(gl.COLOR_BUFFER_BIT)
// 	gl.Enable(gl.BLEND)
// 	gl.Enable(gl.POINT_SMOOTH)
// 	gl.Enable(gl.LINE_SMOOTH)
// 	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
// 	gl.LoadIdentity()

// 	//Transform screen to keep player in middle. Added intentation to make obvious the push matrix is like a block
// 	gl.PushMatrix()

// 	// Draw
// 	for y := 0; y < 32; y++ {
// 		for x := 0; x < 64; x++ {
// 			if myChip8.Gfx[((31-y)*64)+x] == 0 {
// 				gl.Color3f(0.0, 0.0, 0.0)
// 			} else {
// 				gl.Color3f(1.0, 1.0, 1.0)
// 			}
// 			drawPixel(x, y)
// 		}

// 	}

// 	//Second Pop
// 	gl.PopMatrix()
// }

// // onResize sets up a simple 2d ortho context based on the window size
// func onResize(window *glfw.Window, w, h int) {
// 	w, h = window.GetSize() // query window to get screen pixels
// 	width, height := window.GetFramebufferSize()
// 	gl.Viewport(0, 0, int32(width), int32(height))
// 	gl.MatrixMode(gl.PROJECTION)
// 	gl.LoadIdentity()
// 	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
// 	gl.MatrixMode(gl.MODELVIEW)
// 	gl.LoadIdentity()
// 	gl.ClearColor(0, 0, 0, 0)
// }

// func main() {
// 	runtime.LockOSThread()

// 	// initialize glfw
// 	if err := glfw.Init(); err != nil {
// 		log.Fatalln("Failed to initialize GLFW: ", err)
// 	}
// 	defer glfw.Terminate()

// 	// create window
// 	window, err := glfw.CreateWindow(64*MULTIPLIER, 32*MULTIPLIER, os.Args[0], nil, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	window.SetFramebufferSizeCallback(onResize)
// 	window.MakeContextCurrent()

// 	if err := gl.Init(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// set up opengl context
// 	onResize(window, 64*MULTIPLIER, 32*MULTIPLIER)

// 	// glfw.KeyCallback(window)
// 	window.SetKeyCallback(keyCallback)

// 	runtime.LockOSThread()
// 	glfw.SwapInterval(1)

// 	ticker := time.NewTicker(time.Second / 30)
// 	for !window.ShouldClose() {
// 		// myChip8.EmulateCycle()

// 		// handleMouse(window)
// 		window.SwapBuffers()
// 		glfw.PollEvents()

// 		<-ticker.C // wait up to 1/60th of a second
// 	}
// }
