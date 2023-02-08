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

	viewport := gui.NewViewport(report)
	window := gui.InitWindow(width, height, "Tradeview", viewport)
	defer glfw.Terminate()

	program := opengl.MakeProgram()
	program.InitUniformMatrix()
	program.Validate()

	s := scene.New(report)
	s.Build()

	window.OnDraw(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)

		gl.UseProgram(program.Id)
		program.UpdateMatrix(window.ViewInfo.BarsMat)
		s.DrawBars()
		program.UpdateMatrix(window.ViewInfo.VolumesMat)
		s.DrawVolumes()

	})
	window.RunRendering()
}
