package scene

type EventCallback func(WindowState)

type WindowEvents struct {
	OnClick  EventCallback
	OnHover  EventCallback
	OnDrag   EventCallback
	OnZoom   EventCallback
	OnResize EventCallback
}

type WindowState struct {
	MousePosition Point[int32]
	WindowSize    Size[int32]
	DragXDelta    int32
	ZoomIn        bool
}

func CallbackChain(chain, next EventCallback) EventCallback {
	if chain == nil {
		return next
	}
	return func(ws WindowState) {
		chain(ws)
		next(ws)
	}
}
