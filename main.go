package main

import (
	"fmt"
	"runtime"
	"tradeview/config"
	"tradeview/gui"
	"tradeview/market"
	"tradeview/opengl"
	"tradeview/scene"

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
	report := market.LoadReport(config.ReportPath)
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
