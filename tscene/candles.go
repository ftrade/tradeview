package tscene

import (
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/scene"
)

func BuildCandles(candles []market.Candle) (candleRect, minMaxLine *scene.DrawItem) {
	candleN := len(candles)
	rectsVerteces := make([]float32, 4*candleN)
	linesVerteces := make([]float32, 2*2*candleN)
	colors := make([]uint32, candleN)

	barHalfWidth := float32(0.25)
	for i, bar := range candles {
		color := Grean
		if bar.Open > bar.Close {
			color = Red
		}
		colors[i] = color
		x := float32(i)

		linePos := i * 4
		linesVerteces[linePos] = x
		linesVerteces[linePos+1] = bar.Low
		linesVerteces[linePos+2] = x
		linesVerteces[linePos+3] = bar.High

		barLeft := x - barHalfWidth
		barRight := x + barHalfWidth
		vy1, vy2 := bar.Open, bar.Close

		rectPos := i * 4
		rectsVerteces[rectPos] = barLeft
		rectsVerteces[rectPos+1] = barRight
		rectsVerteces[rectPos+2] = vy1
		rectsVerteces[rectPos+3] = vy2
	}

	candleRect = &scene.DrawItem{
		Verteces: rectsVerteces,
		Colors:   colors,
		DrawMode: scene.StaticDrawMode,
		Type:     scene.DrawRects,
	}
	minMaxLine = &scene.DrawItem{
		Verteces: linesVerteces,
		Colors:   colors,
		DrawMode: scene.StaticDrawMode,
		Type:     scene.DrawLines,
	}
	return candleRect, minMaxLine
}
