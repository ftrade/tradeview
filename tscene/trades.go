package tscene

import (
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/scene"
)

func BuildTrades(trades []market.Trade, xAxis *CandleAxis) (*scene.DrawItem, *TradeAxis) {
	tradeAxis := newTradeAxis(len(trades))
	tradesN := len(trades)

	verteces := make([]float32, 4*2*tradesN)
	colors := make([]uint32, 2*tradesN)

	halfWidth := float32(0.25)
	minPrice, maxPrice, _ := xAxis.MinMaxPriceAndMaxVolume(0, xAxis.WidthX())
	for i, trade := range trades {
		t := trade.Timestampt
		x := xAxis.TimeToX(t)
		tradeAxis.trades[i] = XTrade{x, trade}

		offset := i * 8
		verteces[offset] = x
		verteces[offset+1] = minPrice
		verteces[offset+2] = x
		verteces[offset+3] = maxPrice
		verteces[offset+4] = x - halfWidth
		verteces[offset+5] = trade.Price
		verteces[offset+6] = x + halfWidth
		verteces[offset+7] = trade.Price
		colorOffset := i * 2
		color := TradeOpen
		if trade.Profit > 0 {
			color = GoodTrade
		} else if trade.Profit < 0 {
			color = BadTrade
		}
		colors[colorOffset] = color
		colors[colorOffset+1] = color
	}
	return &scene.DrawItem{
		Verteces: verteces,
		Colors:   colors,
		Type:     scene.DrawLines,
		DrawMode: scene.StaticDrawMode,
	}, tradeAxis
}
