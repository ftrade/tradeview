package opengl

import (
	"github.com/ftrade/tradeview/scene"
	"github.com/go-gl/mathgl/mgl32"
)

type StageDrawer struct {
	stage        *scene.Stage
	program      *Program
	itemsDrawers []*itemDrawer
	matrix       mgl32.Mat4
}

func NewStageDrawer(stage *scene.Stage, program *Program) *StageDrawer {
	return &StageDrawer{
		stage:   stage,
		program: program,
	}
}

func (sd *StageDrawer) Init() {
	for _, di := range sd.stage.DrawItems {
		itemDrawer := newItemDrawer(di)
		itemDrawer.init()
		sd.itemsDrawers = append(sd.itemsDrawers, itemDrawer)
	}
}

func (sd *StageDrawer) Draw() {
	if sd.stage.SceneViewUpdated {
		v := sd.stage.SceneView
		sd.matrix = mgl32.Ortho2D(v.Left, v.Right, v.Bottom, v.Top)
		sd.stage.SceneViewUpdated = false
	}
	sd.program.UpdateMatrix(sd.matrix)
	for _, itemDrawer := range sd.itemsDrawers {
		itemDrawer.draw()
	}
}
