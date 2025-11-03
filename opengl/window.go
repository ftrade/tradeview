package opengl

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ftrade/tradeview/scene"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
)

const (
	ErrFontPathMissed  = 3
	OpenGLMajorVersion = 4
	OpenGLMinorVersion = 6
)

const (
	FontColorR = 0.3
	FontColorG = 0.3
	FontColorB = 0.3
)

type Window struct {
	Title         string
	tradeScene    *scene.Scene
	stagesDrawers []*StageDrawer
	// ViewInfo      ViewInfo
	Width, Height int
	window        *glfw.Window
	Font          *glfont.Font
	program       *Program
}

func InitWindow(width, height int, title string, fontSize int32) *Window {
	w := initGlfw(width, height, title)

	fontPath, ok := os.LookupEnv("TRUETYPE_FONT_PATH")
	if !ok {
		slog.Warn("missed TRUETYPE_FONT_PATH environment variable")
		os.Exit(ErrFontPathMissed)
	}
	font, err := glfont.LoadFont(fontPath, fontSize, width, height)
	if err != nil {
		panic(err)
	}

	window := &Window{
		Title:  title,
		Width:  width,
		Height: height,
		window: w,
		Font:   font,
	}

	buildWS := func(x, y float64) scene.WindowState {
		return scene.WindowState{
			MousePosition: scene.Point[int32]{
				X: int32(x),
				Y: int32(y),
			},
			WindowSize: scene.Size[int32]{
				Width:  int32(window.Width),
				Height: int32(window.Height),
			},
		}
	}

	w.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.Width = width
		window.Height = height
		font.UpdateResolution(width, height)
		window.tradeScene.WindowEvents.OnResize(buildWS(0, 0))
	})

	var prevX float64
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			prevX, _ = w.GetCursorPos()
		}
	})

	w.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		ws := buildWS(x, y)
		if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			ws.DragXDelta = int32(prevX - x)
			window.tradeScene.WindowEvents.OnDrag(ws)
			prevX = x
		}
		window.tradeScene.WindowEvents.OnHover(ws)
	})

	w.SetScrollCallback(func(w *glfw.Window, _, yoff float64) {
		x, _ := w.GetCursorPos()

		ws := buildWS(x, 0)
		if yoff < 0 {
			ws.ZoomIn = true
		}
		window.tradeScene.WindowEvents.OnZoom(ws)
	})

	return window
}

func fpsHander() func() {
	var prevTime time.Time
	var frameCount int
	return func() {
		frameCount++
		elapsed := time.Since(prevTime)
		if elapsed.Seconds() > 1 {
			fps := float64(frameCount) / elapsed.Seconds()
			glfw.GetCurrentContext().SetTitle(fmt.Sprintf("Tradeview. FPS: %0.2f", fps))
			prevTime = time.Now()
			frameCount = 0
		}
	}
}

func (w *Window) SetScene(tradeScene *scene.Scene) {
	w.tradeScene = tradeScene
}

func (w *Window) InitScene(program *Program) {
	w.program = program
	program.InitUniformMatrix()
	for _, stage := range w.tradeScene.Stages {
		stageDrawer := NewStageDrawer(stage, program)
		stageDrawer.Init()
		w.stagesDrawers = append(w.stagesDrawers, stageDrawer)
	}
}

func (w *Window) RunRendering() {
	fpsCalc := fpsHander()
	for !w.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)
		gl.UseProgram(w.program.ID)

		for _, sc := range w.stagesDrawers {
			sc.Draw()
		}
		w.Font.SetColor(FontColorR, FontColorG, FontColorB, 1)
		x, y := w.window.GetCursorPos()
		w.tradeScene.DrawDrawers(scene.DrawContext{
			MousePosition: scene.Point[int32]{
				X: int32(x),
				Y: int32(y),
			},
			WindowSize: scene.Size[int32]{
				Width:  int32(w.Width),
				Height: int32(w.Height),
			},
		})

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
	glfw.WindowHint(glfw.ContextVersionMajor, OpenGLMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, OpenGLMinorVersion)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}
