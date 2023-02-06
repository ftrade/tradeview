package main

import (
	"fmt"
	"runtime"
	"time"
	"tradeview/geom"
	"tradeview/market"
	"tradeview/opengl"
	"tradeview/scene"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 1000
	height = 500
)

func main() {
	fmt.Println("App Started")
	report := market.LoadReport("/data/ws/data/candles.xml")
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := opengl.InitOpenGL()
	program.Validate()

	s := scene.New(report)
	s.Build(geom.NewRect(-1, -1, 1, 1))

	for !window.ShouldClose() {

		draw(s, window, program)
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Tradeview", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func draw(s *scene.Scene2D, window *glfw.Window, program opengl.Program) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(1, 1, 1, 1)

	gl.UseProgram(program.Id)

	s.Draw()

	fpsCalc()

	glfw.PollEvents()
	window.SwapBuffers()
}

var prevTime time.Time
var frameCount int

func fpsCalc() {
	frameCount++
	elapsed := time.Since(prevTime)
	if elapsed.Seconds() > 1 {
		fps := float64(frameCount) / elapsed.Seconds()
		fmt.Printf("FPS: %f\n", fps)
		prevTime = time.Now()
		frameCount = 0
	}
}
