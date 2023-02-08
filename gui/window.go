package gui

import (
	"fmt"
	"time"
	"tradeview/config"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
)

type Window struct {
	Title         string
	ViewInfo      ViewInfo
	window        *glfw.Window
	drawScene     func()
	width, height int
	viewport      *Viewport
	font          *glfont.Font
}

func InitWindow(width, height int, title string, viewport *Viewport) *Window {
	w := initGlfw(width, height, title)

	font, err := glfont.LoadFont(config.FontPath, config.FontSize, width, height)
	if err != nil {
		panic(err)
	}

	window := &Window{
		Title:    title,
		window:   w,
		width:    width,
		height:   height,
		font:     font,
		viewport: viewport,
	}
	window.ViewInfo = viewport.CalcView()

	w.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.width = width
		window.height = height
		font.UpdateResolution(width, height)
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
			window.ViewInfo = window.viewport.CalcView()
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
	w.ViewInfo = w.viewport.CalcView()
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

func (w *Window) OnDraw(draw func()) {
	w.drawScene = draw
}

func (w *Window) RunRendering() {
	for !w.window.ShouldClose() {
		w.drawScene()
		w.drawLabels()

		fpsCalc()
		glfw.PollEvents()
		w.window.SwapBuffers()
	}
}

func (w *Window) drawLabels() {
	bottomY := int(float32(w.height) * config.BarsComponentHeight)

	w.font.SetColor(0.5, 0.5, 0.5, 1)
	w.font.Printf(0, float32(bottomY), 1, fmt.Sprintf("%0.2f", w.ViewInfo.MinPrice))
	w.font.Printf(0, config.FontSize, 1, fmt.Sprintf("%0.2f", w.ViewInfo.MaxPrice))
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
