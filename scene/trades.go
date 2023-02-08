package scene

import (
	"tradeview/market"
	"tradeview/opengl"

	"github.com/go-gl/gl/all-core/gl"
)

type Trades struct {
	linesVerteces []float32
	linesVao      uint32
}

func BuildTrades(trades []market.Trade, xAxis XAxis) *Trades {
	tradesN := len(trades)
	linesVerteces := make([]float32, 2*2*tradesN)
	linesColors := make([]uint32, 2*tradesN)

	minPrice, maxPrice, _ := xAxis.MinMaxPriceAndMaxVolume(0, xAxis.WidthX())
	for i, trade := range trades {
		t := trade.Timestampt
		x := xAxis.TimeToX(t)
		vertexOffset := i * 4
		linesVerteces[vertexOffset] = x
		linesVerteces[vertexOffset+1] = minPrice
		linesVerteces[vertexOffset+2] = x
		linesVerteces[vertexOffset+3] = maxPrice
		coloerOffset := i * 2
		color := Blue
		if trade.Profit > 0 {
			color = Grean
		} else if trade.Profit < 0 {
			color = Red
		}
		linesColors[coloerOffset] = color
		linesColors[coloerOffset+1] = color
	}
	return &Trades{
		linesVerteces: linesVerteces,
		linesVao:      opengl.MakeVao(linesVerteces, linesColors),
	}
}

func (t *Trades) Draw() {
	gl.BindVertexArray(t.linesVao)
	gl.DrawArrays(gl.LINES, 0, int32(len(t.linesVerteces)/2))
}
