package gui

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	Title       string
	View        mgl.Mat3
	window      *glfw.Window
	draw        func()
	viewChanged func(mgl.Mat3)
}

func InitWindow(width, height int, title string) *Window {
	w := initGlfw(width, height, title)
	glfw.GetCurrentContext().SetSizeCallback(resizeCb)
	glfw.GetCurrentContext().SetScrollCallback(scrollCb)
	return &Window{
		Title:  title,
		View:   mgl.Ident3(),
		window: w,
	}
}

func resizeCb(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func scrollCb(w *glfw.Window, xoff float64, yoff float64) {
	fmt.Printf("Scroll calb. xoff: %v, yoff: %v\n", xoff, yoff)
	// view = mgl.Scale2D(float32(yoff), float32(yoff))
}

func (w *Window) OnViewChange(cb func(mgl.Mat3)) {
	w.viewChanged = cb
}

func (w *Window) OnDraw(draw func()) {
	w.draw = draw
}

var prevTime time.Time
var frameCount int

func fpsCalc() {
	frameCount++
	elapsed := time.Since(prevTime)
	if elapsed.Seconds() > 1 {
		fps := float64(frameCount) / elapsed.Seconds()
		glfw.GetCurrentContext().SetTitle(fmt.Sprintf("Tradeview. FPS: %f", fps))
		prevTime = time.Now()
		frameCount = 0
	}
}

func (w *Window) RunRendering() {
	for !w.window.ShouldClose() {
		w.draw()

		fpsCalc()
		glfw.PollEvents()
		w.window.SwapBuffers()
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(width, height int, title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}
