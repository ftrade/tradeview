package scene

type Coordinate interface {
	float32 | int32
}

type View[T Coordinate] struct {
	X, Y          T
	Width, Height T
}

type ViewPort struct {
	Window View[int32]
	Scene  View[float32]
}

func (vp *ViewPort) WindowXToSceneX(x int32) float32 {
	return vp.Scene.X + vp.Scene.Width*(float32(x)/float32(vp.Window.Width))
}

// Rect represents a rectangle where [Rect.Top] >= [Rect.Bottom].
type Rect[T Coordinate] struct {
	Left   T
	Right  T
	Bottom T
	Top    T
}

func (r Rect[T]) Move(dx, dy T) Rect[T] {
	return Rect[T]{
		Left:   r.Left + dx,
		Right:  r.Right + dx,
		Bottom: r.Bottom + dy,
		Top:    r.Top + dy,
	}
}

func (r Rect[T]) Contains(p Point[T]) bool {
	if p.X < r.Left || p.X > r.Right {
		return false
	}
	if p.Y < r.Bottom || p.Y > r.Top {
		return false
	}
	return true
}

func (r Rect[T]) Width() T {
	return r.Right - r.Left
}

func (r Rect[T]) Height() T {
	return r.Top - r.Bottom
}

type Point[T Coordinate] struct {
	X T
	Y T
}

type Size[T Coordinate] struct {
	Width  T
	Height T
}
