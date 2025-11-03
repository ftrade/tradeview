package tscene

import "github.com/ftrade/tradeview/scene"

func NewCrosslines() *scene.DrawItem {
	crosslines := &scene.DrawItem{
		Verteces: make([]float32, 8),
		Type:     scene.DrawLines,
	}
	crosslines.WindowEvents = scene.WindowEvents{
		OnHover: func(ws scene.WindowState) {
			crosslines.VertecesChanged = true
			// horizontal line
			crosslines.Verteces[1] = float32(ws.MousePosition.Y)
			crosslines.Verteces[2] = float32(0)
			crosslines.Verteces[3] = float32(ws.MousePosition.Y)
			crosslines.Verteces[2] = float32(ws.WindowSize.Width)
			// vertical line
			crosslines.Verteces[4] = float32(ws.MousePosition.X)
			crosslines.Verteces[5] = float32(0)
			crosslines.Verteces[6] = float32(ws.MousePosition.X)
			crosslines.Verteces[7] = float32(ws.WindowSize.Height)
		},
	}
	return crosslines
}
