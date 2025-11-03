package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/ftrade/tradeview/config"
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/opengl"
	"github.com/ftrade/tradeview/tscene"

	"github.com/go-gl/glfw/v3.3/glfw"
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

		// os.Exit(ErrNoCommandArgToReport)
	}
	report := market.LoadReport("/data/ws/data/candles.xml")
	runtime.LockOSThread()
	opengl.InitOpenGL()

	window := opengl.InitWindow(width, height, "Tradeview", config.FontSize)
	tradeScene := tscene.BuildScene(report, window.Font, height, width)
	// disable V-Sync (remove 60 FPS limit) — must be called after a current GL context exists
	glfw.SwapInterval(0)
	defer glfw.Terminate()

	program := opengl.MakeProgram()
	program.Validate()

	window.SetScene(tradeScene.Scene)
	window.InitScene(&program)
	window.RunRendering()
}
