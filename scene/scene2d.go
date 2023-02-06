package scene

import (
	"tradeview/geom"
	"tradeview/market"
	"tradeview/opengl"

	"github.com/go-gl/gl/v4.6-core/gl"
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

func (s *Scene2D) Build(viewport geom.Rect) {
	stats := s.report.Stats()
	candleN := len(s.report.Candles.Items)
	barRectsVerteces := make([]float32, 12*2*candleN)
	barRectsColors := make([]uint32, 12*candleN)

	barLinesVerteces := make([]float32, 2*2*candleN)
	barLinesColors := make([]uint32, 2*candleN)

	volumeRectsVerteces := make([]float32, 12*2*candleN)
	volumeRectsColors := make([]uint32, 12*candleN)

	timeToX := func(t int64) float32 {
		return float32(t - stats.MinTime)
	}

	barCellWidth := float32(stats.MinTimeStep)
	modelLeft := timeToX(stats.MinTime) - barCellWidth
	modelRight := timeToX(stats.MaxTime) + barCellWidth

	volumeViewport, barViewport := viewport.HSplitPer(0.2)

	bm := geom.NewMapper(geom.NewRect(modelLeft, stats.MinPrice, modelRight, stats.MaxPrice), barViewport)
	vm := geom.NewMapper(geom.NewRect(modelLeft, 0, modelRight, float32(stats.MaxVolume)), volumeViewport)

	vBarWidth := bm.ViewportWithd(barCellWidth) / 2
	vBarHalfWidth := vBarWidth / 2
	for i, bar := range s.report.Candles.Items {
		isRed := false
		color := Grean
		if bar.Open > bar.Close {
			isRed = true
			color = Red
		}

		x := timeToX(bar.Timestampt)
		vX := bm.ViewportX(x)
		linePos := i * 4
		barLinesVerteces[linePos] = vX
		barLinesVerteces[linePos+1] = bm.ViewportY(bar.Low)
		barLinesVerteces[linePos+2] = vX
		barLinesVerteces[linePos+3] = bm.ViewportY(bar.High)

		barLinesColors[i*2] = color
		barLinesColors[i*2+1] = color

		vy1 := bm.ViewportY(bar.Open)
		vy2 := bm.ViewportY(bar.Close)
		if isRed {
			vy1, vy2 = vy2, vy1
		}
		vx1 := vX - vBarHalfWidth
		vx2 := vX + vBarHalfWidth

		//2,5 _____ 6
		//   |\    |
		//   | \   |
		//   |  \  |
		//	 |	 \ |
		//  1|____\|3,4
		//

		rectPos := i * 12
		barRectsVerteces[rectPos] = vx1
		barRectsVerteces[rectPos+1] = vy1

		barRectsVerteces[rectPos+2] = vx1
		barRectsVerteces[rectPos+3] = vy2

		barRectsVerteces[rectPos+4] = vx2
		barRectsVerteces[rectPos+5] = vy1

		barRectsVerteces[rectPos+6] = vx2
		barRectsVerteces[rectPos+7] = vy1

		barRectsVerteces[rectPos+8] = vx1
		barRectsVerteces[rectPos+9] = vy2

		barRectsVerteces[rectPos+10] = vx2
		barRectsVerteces[rectPos+11] = vy2

		barRectsColors[i*6] = color
		barRectsColors[i*6+1] = color
		barRectsColors[i*6+2] = color
		barRectsColors[i*6+3] = color
		barRectsColors[i*6+4] = color
		barRectsColors[i*6+5] = color

		vy1 = vm.ViewportY(0)
		vy2 = vm.ViewportY(float32(bar.Volume))

		volumeRectsVerteces[rectPos] = vx1
		volumeRectsVerteces[rectPos+1] = vy1

		volumeRectsVerteces[rectPos+2] = vx1
		volumeRectsVerteces[rectPos+3] = vy2

		volumeRectsVerteces[rectPos+4] = vx2
		volumeRectsVerteces[rectPos+5] = vy1

		volumeRectsVerteces[rectPos+6] = vx2
		volumeRectsVerteces[rectPos+7] = vy1

		volumeRectsVerteces[rectPos+8] = vx1
		volumeRectsVerteces[rectPos+9] = vy2

		volumeRectsVerteces[rectPos+10] = vx2
		volumeRectsVerteces[rectPos+11] = vy2

		volumeRectsColors[i*6] = color
		volumeRectsColors[i*6+1] = color
		volumeRectsColors[i*6+2] = color
		volumeRectsColors[i*6+3] = color
		volumeRectsColors[i*6+4] = color
		volumeRectsColors[i*6+5] = color
	}
	s.barLinesVerteces = barLinesVerteces
	s.barLinesVao = opengl.MakeVao(barLinesVerteces, barLinesColors)
	s.barRectsVerteces = barRectsVerteces
	s.barRectsVao = opengl.MakeVao(barRectsVerteces, barRectsColors)

	s.volumeRectsVerteces = volumeRectsVerteces
	s.volumeRectsVao = opengl.MakeVao(volumeRectsVerteces, nil)
}

func (s *Scene2D) Draw() {
	gl.BindVertexArray(s.barRectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.barRectsVerteces)/2))
	gl.BindVertexArray(s.barLinesVao)
	gl.DrawArrays(gl.LINES, 0, int32(len(s.barLinesVerteces)/2))
	gl.BindVertexArray(s.volumeRectsVao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.volumeRectsVerteces)/2))
}
