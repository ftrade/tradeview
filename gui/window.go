package gui

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"time"

	"github.com/ftrade/tradeview/config"
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

	fontPath, ok := os.LookupEnv("TRUETYPE_FONT_PATH")
	if !ok {
		slog.Warn("missed TRUETYPE_FONT_PATH environment variable")
		os.Exit(ErrFontPathMissed)
	}
	font, err := glfont.LoadFont(fontPath, config.FontSize, width, height)
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

	w.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.Width = width
		window.Height = height
		font.UpdateResolution(width, height)
	})

	var prevX float64
	w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			prevX, _ = w.GetCursorPos()
		}
	})

	w.SetCursorPosCallback(func(w *glfw.Window, xpos, _ float64) {
		if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			dx := (prevX - xpos) / float64(window.Width)
			window.viewport.Move(float32(dx))
			window.ViewInfo = window.viewport.CalcView()
			prevX = xpos
		}
	})

	w.SetScrollCallback(func(w *glfw.Window, _, yoff float64) {
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

func fpsHander() func() {
	var prevTime time.Time
	var frameCount int
	frameCount++
	return func() {
		elapsed := time.Since(prevTime)
		if elapsed.Seconds() > 1 {
			fps := float64(frameCount) / elapsed.Seconds()
			glfw.GetCurrentContext().SetTitle(fmt.Sprintf("Tradeview. FPS: %0.2f", fps))
			prevTime = time.Now()
			frameCount = 0
		}
	}
}

func (w *Window) OnDraw(draw func()) {
	w.drawScene = draw
}

func (w *Window) RunRendering() {
	fpsCalc := fpsHander()
	for !w.window.ShouldClose() {
		w.drawScene()
		w.font.SetColor(FontColorR, FontColorG, FontColorB, 1)
		w.drawCrosslines()
		w.drawCandleInfo()
		w.drawTradeInfo()
		w.drawPricesLabels()

		fpsCalc()
		glfw.PollEvents()
		w.window.SwapBuffers()
	}
}

func (w *Window) drawPricesLabels() {
	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
	barsBottom := int(float32(w.Height) * config.BarsComponentHeight)
	x, y := w.window.GetCursorPos()
	barsHeight := float32(w.Height) * config.BarsComponentHeight
	if float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width) {
		dPrice := w.ViewInfo.MaxPrice - w.ViewInfo.MinPrice
		price := w.ViewInfo.MaxPrice - (float32(y)/barsHeight)*dPrice

		_ = w.font.Printf(0, float32(y), 1, fmt.Sprintf("%0.2f", price))
	}

	_ = w.font.Printf(0, float32(barsBottom), 1, fmt.Sprintf("%0.2f", w.ViewInfo.MinPrice))
	_ = w.font.Printf(0, config.FontSize, 1, fmt.Sprintf("%0.2f", w.ViewInfo.MaxPrice))
}

func (w *Window) drawCrosslines() {
	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
	x, y := w.window.GetCursorPos()
	barsHeight := float32(w.Height) * config.BarsComponentHeight
	if float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width) {
		w.crosslines.Update(float32(x), float32(y), float32(w.Width), float32(w.Height))
		w.crosslines.Draw()
	}
}

func (w *Window) drawCandleInfo() {
	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
	barsBottom := int(float32(w.Height) * config.BarsComponentHeight)
	x, y := w.window.GetCursorPos()
	barsHeight := float32(w.Height) * config.BarsComponentHeight
	onHover := float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width)
	if !onHover {
		return
	}
	const pad = 20
	viewX := w.viewport.WindowXToViewX(float32(x), float32(w.Width))
	bar, ok := w.viewport.FindBar(viewX)
	if ok {
		time := time.UnixMilli(bar.Timestampt)
		timeStr := time.Format("Mon Jan _2 15:04:05 2006")
		_ = w.font.Printf(float32(x)+pad, float32(barsBottom), 1, timeStr)

		_ = w.font.Printf(float32(x)+pad, float32(y)-pad*4-pad, 1, fmt.Sprintf("Open:  %0.2f", bar.Open))
		_ = w.font.Printf(float32(x)+pad, float32(y)-pad*3-pad, 1, fmt.Sprintf("High:  %0.2f", bar.High))
		_ = w.font.Printf(float32(x)+pad, float32(y)-pad*2-pad, 1, fmt.Sprintf("Low:   %0.2f", bar.Low))
		_ = w.font.Printf(float32(x)+pad, float32(y)-pad*1-pad, 1, fmt.Sprintf("Close: %0.2f", bar.Close))
		_ = w.font.Printf(float32(x)+pad, float32(y)-pad*0-pad, 1, fmt.Sprintf("Vol: %d", bar.Volume))
	}
}

func (w *Window) drawTradeInfo() {
	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
	x, y := w.window.GetCursorPos()
	barsHeight := float32(w.Height) * config.BarsComponentHeight
	onHover := float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width)
	if !onHover {
		return
	}
	const pad = 20
	viewX := w.viewport.WindowXToViewX(float32(x), float32(w.Width))

	// try make convient search distance depending scale and window width
	wdx := math.Pow(float64(w.viewport.ViewWidth()*0.1), 0.5) //nolint:mnd // just empirical numbers
	searchDistance := float32(wdx) * float32(math.Log10(wdx))
	const minDistance = 0.25
	if searchDistance < minDistance {
		searchDistance = minDistance
	}
	trades := w.tradeAxis.FindTrade(viewX, searchDistance)
	switch len(trades) {
	case 0:
		return
	case 1:
		trade := trades[0].Trade
		_ = w.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("Price: %0.2f", trade.Price))
		_ = w.font.Printf(float32(x)+pad, float32(y)+20*2+pad, 1, fmt.Sprintf("Vol: %d", trade.Volume))
		time := time.UnixMilli(trade.Timestampt)
		timeStr := time.Format("15:04:05")
		_ = w.font.Printf(float32(x)+pad, float32(y)+20*3+pad, 1, fmt.Sprintf("Time: %s", timeStr))
		if trade.Profit != 0 {
			_ = w.font.Printf(float32(x)+pad, float32(y)+20*4+pad, 1, fmt.Sprintf("Profit: %0.2f", trade.Profit))
		}
	default:
		_ = w.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("%d trades", len(trades)))
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
