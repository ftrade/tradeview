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
	View          mgl.Mat4
	window        *glfw.Window
	draw          func()
	viewChanged   func(mgl.Mat4)
	projection    mgl.Mat4
	width, height int
}

func InitWindow(width, height int, title string) *Window {
	w := initGlfw(width, height, title)

	window := &Window{
		Title:      title,
		View:       mgl.Ident4(),
		window:     w,
		projection: mgl.Ortho2D(-1, 1, -1, 1),
		width:      width,
		height:     height,
	}
	w.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.width = width
		window.height = height
	})

	var pressedX, pressedY float64
	wasPressed := false
	var movedView mgl.Mat4
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			pressedX, pressedY = w.GetCursorPos()
			wasPressed = true
		} else if button == glfw.MouseButtonLeft && action == glfw.Release && wasPressed {
			wasPressed = false
			window.View = movedView
		}
	})
	windowToView := func(wx, wy float64) (float32, float32) {
		wp := mgl.Vec3{float32(wx), float32(window.height) - float32(wy), 0}
		vp, _ := mgl.UnProject(wp, window.View, window.projection, 0, 0, window.width, window.height)
		return vp[0], vp[1]
	}
	w.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			pressedVX, pressedVY := windowToView(pressedX, pressedY)
			vx, vy := windowToView(xpos, ypos)
			move := mgl.Translate3D(vx-pressedVX, vy-pressedVY, 0)
			movedView = window.View.Mul4(move)
			window.viewChanged(movedView)
		}
	})

	w.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		var scaleFactor float32 = 0.8
		if yoff > 0 {
			scaleFactor = 1.2
		}
		wx, wy := w.GetCursorPos()
		vx, vy := windowToView(wx, wy)

		moveToMouse := mgl.Translate3D(vx, vy, 0)
		moveBack := mgl.Translate3D(-vx, -vy, 0)
		scale := mgl.Scale3D(scaleFactor, scaleFactor, 1)
		window.View = window.View.Mul4(moveToMouse).Mul4(scale).Mul4(moveBack)
		window.viewChanged(window.View)
	})
	w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if mods&glfw.ModControl != 0 {
			if key == glfw.KeyUp || key == glfw.KeyDown {
				scale := mgl.Scale3D(1, 1.2, 1)
				if key == glfw.KeyDown {
					scale = mgl.Scale3D(1, 0.8, 1)
				}
				window.View = window.View.Mul4(scale)
				window.viewChanged(window.View)
			}
		}

	})

	return window
}

func (w *Window) OnViewChange(cb func(mgl.Mat4)) {
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
