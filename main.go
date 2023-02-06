package main

import (
	"fmt"
	"runtime"
	"tradeview/geom"
	"tradeview/gui"
	"tradeview/market"
	"tradeview/opengl"
	"tradeview/scene"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 1000
	height = 500
)

func main() {
	fmt.Println("App Started")
	report := market.LoadReport("/data/ws/data/candles_small.xml")
	runtime.LockOSThread()

	window := gui.InitWindow(width, height, "Tradeview")
	defer glfw.Terminate()

	program := opengl.InitOpenGL()
	program.Validate()

	s := scene.New(report)
	s.Build(geom.NewRect(-1, -1, 1, 1))

	program.InitMatrix(window.View)
	window.OnViewChange(program.UpdateMatrix)

	window.OnDraw(func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)

		gl.UseProgram(program.Id)
		s.Draw()
	})
	window.RunRendering()
}
