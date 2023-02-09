# About

Tradeview is an application for rendering static market data: candles and trades. It was created to analyze large amount of items with high performance. It renders around 70k candles with 700fps on Nvidia 1050ti. 

# Features

* High performance comparing other open source alternative (that are not a lot)
* Automatically scaling y axis to fit window height
* Scaling in/out around mouse x coordinate
* Show candle and trade info on mouse hover (without matching y)

# How to use

## Requirements

Tradeview uses OpenGL for rendering via [go-gl](https://github.com/go-gl/gl). go-gl uses cgo to invoke OpenGL API and your OS must have corresponding libraries. For debian/ubuntu you can install by:

```bash
sudo apt install xorg-dev
```
## Install binary

```bash
go install github.com/ftrade/tradeview@v0.1
```

## Run binary

Application require two inputs:
* TrueType font path -- for rendering a text. Passed via environment variable TRUETYPE_FONT_PATH
* path to report file that contains candles and trades. Passed as first program argument

```bash
TRUETYPE_FONT_PATH=path_to_tff tradeview path_to_report
```

You can use [report_sample.xml](https://github.com/ftrade/tradeview/report_sample.xml) as a small report file to test the application.