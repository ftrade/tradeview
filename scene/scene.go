package scene

type DrawContext struct {
	MousePosition Point[int32]
	WindowSize    Size[int32]
}

type Drawer interface {
	Draw(DrawContext)
}

type Scene struct {
	Stages       []*Stage
	WindowEvents WindowEvents
	Drawers      []Drawer
}

func NewScene() *Scene {
	s := &Scene{}
	passEventToStage := func(cbF func(stage *Stage) EventCallback) EventCallback {
		return func(ws WindowState) {
			for _, stage := range s.Stages {
				cbF(stage)(ws)
			}
		}
	}
	s.WindowEvents = WindowEvents{
		OnClick: passEventToStage(func(stage *Stage) EventCallback {
			return stage.WindowEvents.OnClick
		}),
		OnZoom: passEventToStage(func(stage *Stage) EventCallback {
			return stage.WindowEvents.OnZoom
		}),
		OnResize: passEventToStage(func(stage *Stage) EventCallback {
			return stage.WindowEvents.OnResize
		}),
		OnHover: passEventToStage(func(stage *Stage) EventCallback {
			return stage.WindowEvents.OnHover
		}),
		OnDrag: passEventToStage(func(stage *Stage) EventCallback {
			return stage.WindowEvents.OnDrag
		}),
	}
	return s
}

func (s *Scene) AddStage(stage *Stage) {
	s.Stages = append(s.Stages, stage)
}

func (s *Scene) DrawDrawers(dctx DrawContext) {
	for _, drawer := range s.Drawers {
		drawer.Draw(dctx)
	}
}

func (s *Scene) AddDrawer(drawer Drawer) {
	s.Drawers = append(s.Drawers, drawer)
}

type Stage struct {
	// WindowAnchor describe what part of the window will be used for rendering.
	//
	// Maximum value for top and right is 1.0
	WindowAnchor Rect[float32] //TODO not used at all
	// SceneView specify what part of scene need show
	SceneView        Rect[float32]
	SceneViewUpdated bool
	DrawItems        []*DrawItem
	WindowEvents     WindowEvents
}

func NewStage(windowAnchor Rect[float32]) *Stage {
	stage := &Stage{
		WindowAnchor:     windowAnchor,
		SceneViewUpdated: true,
	}

	passEventToDrawItems := func(cbF func(di *DrawItem) EventCallback) EventCallback {
		return func(ws WindowState) {
			for _, di := range stage.DrawItems {
				cb := cbF(di)
				if cb != nil {
					cb(ws)
				}
			}
		}
	}
	stage.WindowEvents = WindowEvents{
		OnHover: passEventToDrawItems(func(di *DrawItem) EventCallback {
			return di.WindowEvents.OnHover
		}),
		OnClick: passEventToDrawItems(func(di *DrawItem) EventCallback {
			return di.WindowEvents.OnClick
		}),
		OnDrag: passEventToDrawItems(func(di *DrawItem) EventCallback {
			return di.WindowEvents.OnDrag
		}),
		OnZoom: passEventToDrawItems(func(di *DrawItem) EventCallback {
			return di.WindowEvents.OnZoom
		}),
		OnResize: passEventToDrawItems(func(di *DrawItem) EventCallback {
			return di.WindowEvents.OnResize
		}),
	}
	return stage
}

func (s *Stage) AddDrawItem(drawItem *DrawItem) {
	s.DrawItems = append(s.DrawItems, drawItem)
}

type DrawType byte

const (
	DrawLines DrawType = iota
	DrawRects
)

type DrawItem struct {
	Verteces        []float32
	VertecesChanged bool
	Colors          []uint32
	Type            DrawType
	DrawMode        DrawMode
	WindowEvents    WindowEvents
}

type DrawMode byte

const (
	StaticDrawMode = iota
	DynamicDrawMode
)
