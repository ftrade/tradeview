package scene

import (
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/opengl"

	"github.com/go-gl/gl/all-core/gl"
)

type Trades struct {
	TradeAxis     TradeAxis
	linesVerteces []float32
	linesVao      uint32
}

func BuildTrades(trades []market.Trade, xAxis BarAxis) *Trades {
	tradeAxis := newTradeAxis(len(trades))
	tradesN := len(trades)

	linesVerteces := make([]float32, 4*2*tradesN)
	linesColors := make([]uint32, 4*tradesN)

	halfWidth := float32(0.25)
	minPrice, maxPrice, _ := xAxis.MinMaxPriceAndMaxVolume(0, xAxis.WidthX())
	for i, trade := range trades {
		t := trade.Timestampt
		x := xAxis.TimeToX(t)
		tradeAxis.trades[i] = XTrade{x, trade}

		vertexOffset := i * 8
		linesVerteces[vertexOffset] = x
		linesVerteces[vertexOffset+1] = minPrice
		linesVerteces[vertexOffset+2] = x
		linesVerteces[vertexOffset+3] = maxPrice
		linesVerteces[vertexOffset+4] = x - halfWidth
		linesVerteces[vertexOffset+5] = trade.Price
		linesVerteces[vertexOffset+6] = x + halfWidth
		linesVerteces[vertexOffset+7] = trade.Price
		coloerOffset := i * 4
		color := TradeOpen
		if trade.Profit > 0 {
			color = GoodTrade
		} else if trade.Profit < 0 {
			color = BadTrade
		}
		linesColors[coloerOffset] = color
		linesColors[coloerOffset+1] = color
		linesColors[coloerOffset+2] = color
		linesColors[coloerOffset+3] = color
	}
	return &Trades{
		TradeAxis:     tradeAxis,
		linesVerteces: linesVerteces,
		linesVao:      opengl.MakeVao(linesVerteces, linesColors),
	}
}

func (t *Trades) Draw() {
	gl.BindVertexArray(t.linesVao)
	gl.DrawArrays(gl.LINES, 0, int32(len(t.linesVerteces)/2))
}
