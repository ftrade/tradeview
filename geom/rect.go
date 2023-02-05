package geom

type Rect struct {
	P1, P2 Point
}

func NewRect(x1 float32, y1 float32, x2 float32, y2 float32) Rect {
	return Rect{
		P1: Point{x1, y1},
		P2: Point{x2, y2},
	}
}

func (r *Rect) HSplit(y float32) (Rect, Rect) {
	lower := NewRect(r.P1.X, r.P1.Y, r.P2.X, y)
	upper := NewRect(r.P1.X, y, r.P2.X, r.P2.Y)
	return lower, upper
}

func (r *Rect) HSplitPer(percent float32) (Rect, Rect) {
	dy := r.P2.Y - r.P1.Y
	y := r.P1.Y + dy*percent
	return r.HSplit(y)
}

func (r *Rect) Width() float32 {
	return r.P2.X - r.P1.X
}

func (r *Rect) Height() float32 {
	return r.P2.Y - r.P1.Y
}
