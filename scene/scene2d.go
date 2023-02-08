package scene

import (
	"tradeview/market"
	"tradeview/opengl"

	"github.com/go-gl/gl/all-core/gl"
)

type Scene2D struct {
	report              market.Report
	barRectsVao         uint32
	barLinesVao         uint32
	volumeRectsVao      uint32
	barRectsVerteces    []float32
	barLinesVerteces    []float32
	volumeRectsVerteces []float32
}

func New(report market.Report) *Scene2D {
	return &Scene2D{
		report: report,
	}
}

func (s *Scene2D) Build() {
	bars := s.report.Candles.Items
	candleN := len(bars)
	barRectsVerteces := make([]float32, 12*2*candleN)
	barRectsColors := make([]uint32, 12*candleN)

	barLinesVerteces := make([]float32, 2*2*candleN)
	barLinesColors := make([]uint32, 2*candleN)

	volumeRectsVerteces := make([]float32, 12*2*candleN)

	barHalfWidth := float32(0.25)
	for i, bar := range bars {
		isRed := false
		color := Grean
		if bar.Open > bar.Close {
			isRed = true
			color = Red
		}
		x := float32(i)

		linePos := i * 4
		barLinesVerteces[linePos] = x
		barLinesVerteces[linePos+1] = bar.Low
		barLinesVerteces[linePos+2] = x
		barLinesVerteces[linePos+3] = bar.High
		barLinesColors[i*2] = color
		barLinesColors[i*2+1] = color

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
		barRectsVerteces[rectPos] = barLeft
		barRectsVerteces[rectPos+1] = vy1

		barRectsVerteces[rectPos+2] = barLeft
		barRectsVerteces[rectPos+3] = vy2

		barRectsVerteces[rectPos+4] = barRight
		barRectsVerteces[rectPos+5] = vy1

		barRectsVerteces[rectPos+6] = barRight
		barRectsVerteces[rectPos+7] = vy1

		barRectsVerteces[rectPos+8] = barLeft
		barRectsVerteces[rectPos+9] = vy2

		barRectsVerteces[rectPos+10] = barRight
		barRectsVerteces[rectPos+11] = vy2

		barRectsColors[i*6] = color
		barRectsColors[i*6+1] = color
		barRectsColors[i*6+2] = color
		barRectsColors[i*6+3] = color
		barRectsColors[i*6+4] = color
		barRectsColors[i*6+5] = color

		//same order as bar's one
		vy1, vy2 = 0, float32(bar.Volume)

		volumeRectsVerteces[rectPos] = barLeft
		volumeRectsVerteces[rectPos+1] = vy1

		volumeRectsVerteces[rectPos+2] = barLeft
		volumeRectsVerteces[rectPos+3] = vy2

		volumeRectsVerteces[rectPos+4] = barRight
		volumeRectsVerteces[rectPos+5] = vy1

		volumeRectsVerteces[rectPos+6] = barRight
		volumeRectsVerteces[rectPos+7] = vy1

		volumeRectsVerteces[rectPos+8] = barLeft
		volumeRectsVerteces[rectPos+9] = vy2

		volumeRectsVerteces[rectPos+10] = barRight
		volumeRectsVerteces[rectPos+11] = vy2
	}

	s.barLinesVerteces = barLinesVerteces
	s.barLinesVao = opengl.MakeVao(barLinesVerteces, barLinesColors)
	s.barRectsVerteces = barRectsVerteces
	s.barRectsVao = opengl.MakeVao(barRectsVerteces, barRectsColors)

	s.volumeRectsVerteces = volumeRectsVerteces
	s.volumeRectsVao = opengl.MakeVao(volumeRectsVerteces, nil)
}

func (s *Scene2D) DrawBars() {
	gl.BindVertexArray(s.barRectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.barRectsVerteces)/2))
	gl.BindVertexArray(s.barLinesVao)
	gl.DrawArrays(gl.LINES, 0, int32(len(s.barLinesVerteces)/2))
}

func (s *Scene2D) DrawVolumes() {
	gl.BindVertexArray(s.volumeRectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.volumeRectsVerteces)/2))
}
