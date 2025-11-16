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
	width             = 1000
	height            = 500
	ErrFailToLoadEnvs = 2
	ErrParseLogLevel  = 3
)

func main() {
	envs, err := config.LoadEnvs()
	if err != nil {
		slog.Error("Fail to load the environment file.", "error", err.Error())
		os.Exit(ErrFailToLoadEnvs)
	}

	err = config.InitLogger(envs.LogLevel)
	if err != nil {
		slog.Error("Fail to parse the log level environment parameter.", "error", err.Error())
		os.Exit(ErrParseLogLevel)
	}

	slog.Info("App Started")
	report := market.LoadReport(envs.MarketFilePath)
	runtime.LockOSThread()
	opengl.InitOpenGL()
	window := opengl.InitWindow(width, height, "Tradeview", envs.FontFilePath, envs.FontSize)
	tradeScene := tscene.BuildScene(report, window.Font, envs.FontSize, height, width)
	defer glfw.Terminate()

	if !envs.VsyncEnabled {
		glfw.SwapInterval(0)
	}

	program := opengl.MakeProgram(tscene.Colors())
	program.Validate()
	window.SetScene(tradeScene.Scene)
	window.InitScene(&program)
	window.RunRendering()
}
