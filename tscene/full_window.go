package tscene

import "github.com/ftrade/tradeview/scene"

func NewFullWindowStage(width, height int32) *scene.Stage {
	stage := scene.NewStage(scene.Rect[float32]{
		Left:   0,
		Right:  1,
		Bottom: 0,
		Top:    1,
	})
	stage.SceneView = scene.Rect[float32]{
		Right:  float32(width),
		Bottom: float32(height),
	}
	onResize := func(ws scene.WindowState) {
		stage.SceneView = scene.Rect[float32]{
			Left:   0,
			Right:  float32(ws.WindowSize.Width),
			Bottom: float32(ws.WindowSize.Height),
			Top:    0,
		}
		stage.SceneViewUpdated = true
	}
	stage.WindowEvents.OnResize = scene.CallbackChain(stage.WindowEvents.OnResize, onResize)

	return stage
}
