package main

import (
	"fmt"
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
	width  = 1000
	height = 500
)

func main() {
	fmt.Println("App Started")
	if len(os.Args) < 2 {
		fmt.Println("missed path to report CLI argument")
		os.Exit(2)
	}
	report := market.LoadReport(os.Args[1])
	runtime.LockOSThread()
	opengl.InitOpenGL()

	xAxis := scene.NewXAxis(report.Candles.Items)
	viewport := gui.NewViewport(xAxis)
	window := gui.InitWindow(width, height, "Tradeview", viewport)
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

		gl.UseProgram(program.Id)
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
