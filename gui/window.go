package gui

import (
	"fmt"
	"math"
	"time"
	"tradeview/config"
	"tradeview/scene"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nullboundary/glfont"
)

type Window struct {
	Title         string
	ViewInfo      ViewInfo
	Width, Height int
	crosslines    *scene.CrossLines
	window        *glfw.Window
	drawScene     func()
	viewport      *Viewport
	font          *glfont.Font
	tradeAxis     scene.TradeAxis
}

func InitWindow(width, height int, title string, viewport *Viewport) *Window {
	w := initGlfw(width, height, title)

	font, err := glfont.LoadFont(config.FontPath, config.FontSize, width, height)
	if err != nil {
		panic(err)
	}

	window := &Window{
		Title:      title,
		Width:      width,
		Height:     height,
		window:     w,
		crosslines: scene.NewCrossLines(),
		font:       font,
		viewport:   viewport,
	}
	window.ViewInfo = viewport.CalcView()

	w.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.Width = width
		window.Height = height
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
			dx := (prevX - xpos) / float64(window.Width)
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
		x, _ := w.GetCursorPos()
		whereScale := x / float64(window.Width)
		window.viewport.Scale(scaleFactor, float32(whereScale))
		window.ViewInfo = window.viewport.CalcView()
	})

	return window
}

func (w *Window) SetTradeAxis(ta scene.TradeAxis) {
	w.tradeAxis = ta
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
	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
	barsBottom := int(float32(w.Height) * config.BarsComponentHeight)
	x, y := w.window.GetCursorPos()
	barsHeight := float32(w.Height) * config.BarsComponentHeight
	if float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width) {
		dPrice := w.ViewInfo.MaxPrice - w.ViewInfo.MinPrice
		price := w.ViewInfo.MaxPrice - (float32(y)/barsHeight)*dPrice
		w.crosslines.Update(float32(x), float32(y), float32(w.Width), float32(w.Height))
		w.crosslines.Draw()

		pad := float32(20)
		w.font.SetColor(0.3, 0.3, 0.3, 1)
		w.font.Printf(0, float32(y), 1, fmt.Sprintf("%0.2f", price))
		viewX := w.viewport.WindowXToViewX(float32(x), float32(w.Width))
		bar, ok := w.viewport.FindBar(viewX)
		if ok {

			time := time.UnixMilli(bar.Timestampt)
			timeStr := time.Format("Mon Jan _2 15:04:05 2006")
			w.font.Printf(float32(x)+pad, float32(barsBottom), 1, timeStr)

			w.font.Printf(float32(x)+pad, float32(y)-20*4-pad, 1, fmt.Sprintf("Open:  %0.2f", bar.Open))
			w.font.Printf(float32(x)+pad, float32(y)-20*3-pad, 1, fmt.Sprintf("High:  %0.2f", bar.High))
			w.font.Printf(float32(x)+pad, float32(y)-20*2-pad, 1, fmt.Sprintf("Low:   %0.2f", bar.Low))
			w.font.Printf(float32(x)+pad, float32(y)-20*1-pad, 1, fmt.Sprintf("Close: %0.2f", bar.Close))
			w.font.Printf(float32(x)+pad, float32(y)-20*0-pad, 1, fmt.Sprintf("Vol: %d", bar.Volume))
		}

		// try make convient search distance depending scale and window width
		wdx := math.Pow(float64(w.viewport.ViewWidth()*0.1), 0.5)
		searchDistance := float32(wdx) * float32(math.Log10(wdx))
		if searchDistance < 0.25 {
			searchDistance = 0.25
		}
		trades := w.tradeAxis.FindTrade(viewX, searchDistance)
		if len(trades) > 0 {
			if len(trades) == 1 {
				trade := trades[0].Trade
				w.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("Price: %0.2f", trade.Price))
				w.font.Printf(float32(x)+pad, float32(y)+20*2+pad, 1, fmt.Sprintf("Vol: %d", trade.Volume))
				time := time.UnixMilli(trade.Timestampt)
				timeStr := time.Format("15:04:05")
				w.font.Printf(float32(x)+pad, float32(y)+20*3+pad, 1, fmt.Sprintf("Time: %s", timeStr))
				if trade.Profit != 0 {
					w.font.Printf(float32(x)+pad, float32(y)+20*4+pad, 1, fmt.Sprintf("Profit: %0.2f", trade.Profit))
				}

			} else {
				w.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("%d trades", len(trades)))
			}
		}
	}

	w.font.SetColor(0.3, 0.3, 0.3, 1)
	w.font.Printf(0, float32(barsBottom), 1, fmt.Sprintf("%0.2f", w.ViewInfo.MinPrice))
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
