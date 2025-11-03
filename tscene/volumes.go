package tscene

import (
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/scene"
)

func BuildVolumes(candles []market.Candle) *scene.DrawItem {
	candleN := len(candles)
	rectsVerteces := make([]float32, 4*candleN)

	barHalfWidth := float32(0.25)
	for i, bar := range candles {
		x := float32(i)

		barLeft := x - barHalfWidth
		barRight := x + barHalfWidth

		rectPos := i * 4
		rectsVerteces[rectPos] = barLeft
		rectsVerteces[rectPos+1] = barRight
		rectsVerteces[rectPos+2] = 0
		rectsVerteces[rectPos+3] = float32(bar.Volume)
	}

	return &scene.DrawItem{
		Verteces: rectsVerteces,
		DrawMode: scene.StaticDrawMode,
		Type:     scene.DrawRects,
	}
}
