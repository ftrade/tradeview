package gui

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	Title         string
	BarMatrix     mgl.Mat4
	VolumeMatrix  mgl.Mat4
	window        *glfw.Window
	draw          func()
	width, height int
	viewport      *Viewport
}

func InitWindow(width, height int, title string) *Window {
	w := initGlfw(width, height, title)

	window := &Window{
		Title:  title,
		window: w,
		width:  width,
		height: height,
	}
	w.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.width = width
		window.height = height
	})

	var prevX float64
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			prevX, _ = w.GetCursorPos()
		}
	})

	w.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			dx := (prevX - xpos) / float64(window.width)
			window.viewport.Move(float32(dx))
			window.BarMatrix, window.VolumeMatrix = window.viewport.Matrices()
			prevX = xpos
		}
	})

	w.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		var scaleFactor float32 = 1.2
		if yoff > 0 {
			scaleFactor = 0.8
		}
		window.scale(scaleFactor)
	})

	return window
}

func (w *Window) scale(factor float32) {
	if factor < 0 {
		return
	}
	x, _ := w.window.GetCursorPos()
	whereScale := x / float64(w.width)
	w.viewport.Scale(factor, float32(whereScale))
	w.BarMatrix, w.VolumeMatrix = w.viewport.Matrices()
}

func (w *Window) SetViewport(viewport *Viewport) {
	w.viewport = viewport
	w.BarMatrix, w.VolumeMatrix = viewport.Matrices()
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
		glfw.GetCurrentContext().SetTitle(fmt.Sprintf("Tradeview. FPS: %0.2f", fps))
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
