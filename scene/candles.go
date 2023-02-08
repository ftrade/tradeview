package scene

import (
	"tradeview/market"
	"tradeview/opengl"

	"github.com/go-gl/gl/all-core/gl"
)

type Candles struct {
	rectsVerteces []float32
	rectsVao      uint32
	linesVerteces []float32
	linesVao      uint32
}

func BuildCandles(candles []market.Candle) *Candles {
	candleN := len(candles)
	rectsVerteces := make([]float32, 12*2*candleN)
	rectsColors := make([]uint32, 12*candleN)

	linesVerteces := make([]float32, 2*2*candleN)
	linesColors := make([]uint32, 2*candleN)

	barHalfWidth := float32(0.25)
	for i, bar := range candles {
		isRed := false
		color := Grean
		if bar.Open > bar.Close {
			isRed = true
			color = Red
		}
		x := float32(i)

		linePos := i * 4
		linesVerteces[linePos] = x
		linesVerteces[linePos+1] = bar.Low
		linesVerteces[linePos+2] = x
		linesVerteces[linePos+3] = bar.High
		linesColors[i*2] = color
		linesColors[i*2+1] = color

		barLeft := x - barHalfWidth
		barRight := x + barHalfWidth
		vy1, vy2 := bar.Open, bar.Close
		if isRed {
			vy1, vy2 = vy2, vy1
		}

		//2,5 _____ 6
		//   |\    |
		//   | \   |
		//   |  \  |
		//	 |	 \ |
		//  1|____\|3,4
		//

		rectPos := i * 12
		rectsVerteces[rectPos] = barLeft
		rectsVerteces[rectPos+1] = vy1

		rectsVerteces[rectPos+2] = barLeft
		rectsVerteces[rectPos+3] = vy2

		rectsVerteces[rectPos+4] = barRight
		rectsVerteces[rectPos+5] = vy1

		rectsVerteces[rectPos+6] = barRight
		rectsVerteces[rectPos+7] = vy1

		rectsVerteces[rectPos+8] = barLeft
		rectsVerteces[rectPos+9] = vy2

		rectsVerteces[rectPos+10] = barRight
		rectsVerteces[rectPos+11] = vy2

		rectsColors[i*6] = color
		rectsColors[i*6+1] = color
		rectsColors[i*6+2] = color
		rectsColors[i*6+3] = color
		rectsColors[i*6+4] = color
		rectsColors[i*6+5] = color
	}

	return &Candles{
		rectsVerteces: rectsVerteces,
		rectsVao:      opengl.MakeVao(rectsVerteces, rectsColors),
		linesVerteces: linesVerteces,
		linesVao:      opengl.MakeVao(linesVerteces, linesColors),
	}
}

func (cs *Candles) Draw() {
	gl.BindVertexArray(cs.rectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cs.rectsVerteces)/2))
	gl.BindVertexArray(cs.linesVao)
	gl.DrawArrays(gl.LINES, 0, int32(len(cs.linesVerteces)/2))
}
