package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/ftrade/tradeview/gui"
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/opengl"
	"github.com/ftrade/tradeview/scene"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width                   = 1000
	height                  = 500
	ErrNoCommandArgToReport = 2
	MinCmdArgCount          = 2
)

func main() {
	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))

	slog.Info("App Started")
	if len(os.Args) < MinCmdArgCount {
		slog.Info("missed path to report CLI argument")
		os.Exit(ErrNoCommandArgToReport)
	}
	report := market.LoadReport(os.Args[1])
	runtime.LockOSThread()
	opengl.InitOpenGL()

	xAxis := scene.NewXAxis(report.Candles.Items)
	viewport := gui.NewViewport(xAxis)
	window := gui.InitWindow(width, height, "Tradeview", viewport)
	// disable V-Sync (remove 60 FPS limit) — must be called after a current GL context exists
	// glfw.SwapInterval(0)
	defer glfw.Terminate()

	program := opengl.MakeProgram()
	program.InitUniformMatrix()
	program.Validate()

	candles := scene.BuildCandles(report.Candles.Items)
	volumes := scene.BuildVolumes(report.Candles.Items)
	trades := scene.BuildTrades(report.Trades.Items, xAxis)
	window.SetTradeAxis(trades.TradeAxis)

	window.OnDraw(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)

		gl.UseProgram(program.ID)
		program.UpdateMatrix(window.ViewInfo.BarsMat)
		candles.Draw()
		trades.Draw()
		program.UpdateMatrix(window.ViewInfo.VolumesMat)
		volumes.Draw()
		// matrix for crosslines
		program.UpdateMatrix(mgl32.Ortho2D(0, float32(window.Width), float32(window.Height), 0))
	})
	window.RunRendering()
}
