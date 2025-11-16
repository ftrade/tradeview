# About

Tradeview is an application for rendering static market data: candles and trades. It was created to analyze large amount of items with high performance.

It renders around 70k candles about 3000fps on AMD Ryzen™ 5 7640HS w/ Radeon™ 760M Graphics.

# Features

* High performance comparing other open source alternative (that are not a lot)
* Automatically scaling y axis to fit window height
* Scaling in/out around mouse x coordinate
* Show candle and trade info on mouse hover (without matching y)

# How to use

## Requirements

Tradeview uses OpenGL for rendering via [go-gl](https://github.com/go-gl/gl). go-gl uses cgo to invoke OpenGL API and your OS must have corresponding libraries. For debian/ubuntu you can install by:

```bash
sudo apt install xorg-dev fonts-freefont-ttf 
```
## Install binary

```bash
go install github.com/ftrade/tradeview@v0.1.0
```

## Run binary

Application require two inputs:
* TrueType font path -- for rendering a text. Passed via environment variable FONT_FILE_PATH
* path to report file that contains candles and trades. Passed via environment variable MARKET_FILE_PATH

```bash
TRUETYPE_FONT_PATH=path_to_tff MARKET_FILE_PATH=path_to_report tradeview
```
## Configuration

All configuration are passed as environment variables. They also can be passed using .env file.
Environment variables:
* FONT_FILE_PATH -- file path to TrueType font, required
* MARKET_FILE_PATH -- file path to market report (with candles and other staff), required. You can see example of a file in examples directory.
* VSYNC_ENABLED -- whether vsync is enabled, optional, default is "true".
* LOG_LEVEL -- log level for slog, default is "INFO".
* FONT_SIZE -- fon size for labels, default is 20.

You can use [report_sample.xml](https://github.com/ftrade/tradeview/blob/master/examples/report_sample.xml) as a small report file to test the application.
