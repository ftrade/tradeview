package scene

import (
	"tradeview/market"
	"tradeview/opengl"

	"github.com/go-gl/gl/all-core/gl"
)

type Volumes struct {
	rectsVao      uint32
	rectsVerteces []float32
}

func BuildVolumes(candles []market.Candle) *Volumes {
	candleN := len(candles)
	rectsVerteces := make([]float32, 12*2*candleN)

	barHalfWidth := float32(0.25)
	for i, bar := range candles {
		x := float32(i)

		barLeft := x - barHalfWidth
		barRight := x + barHalfWidth

		//2,5 _____ 6
		//   |\    |
		//   | \   |
		//   |  \  |
		//	 |	 \ |
		//  1|____\|3,4
		//

		rectPos := i * 12
		vy1, vy2 := float32(0), float32(bar.Volume)

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
	}

	return &Volumes{
		rectsVerteces: rectsVerteces,
		rectsVao:      opengl.MakeVao(rectsVerteces, nil),
	}
}

func (v *Volumes) Draw() {
	gl.BindVertexArray(v.rectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(v.rectsVerteces)/2))
}
