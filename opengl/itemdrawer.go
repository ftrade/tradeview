package opengl

import (
	"fmt"
	"log/slog"

	"github.com/ftrade/tradeview/scene"
	"github.com/go-gl/gl/all-core/gl"
)

type itemDrawer struct {
	drawItem *scene.DrawItem
	verteces []float32
	colors   []uint32
	vbo      uint32
	vao      uint32
}

func newItemDrawer(drawItem *scene.DrawItem) *itemDrawer {
	return &itemDrawer{
		drawItem: drawItem,
	}
}

func (id *itemDrawer) init() {
	id._recomputeVerteces()
	id._computeColors()

	gl.GenVertexArrays(1, &id.vao)
	gl.BindVertexArray(id.vao)

	gl.GenBuffers(1, &id.vbo)
	id._updateVertexBuffer()
	const _2d = 2
	gl.VertexAttribPointer(0, _2d, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	if len(id.colors) > 0 {
		var colorsVbo uint32 // as we won't update colors we don't save id of vbo
		gl.GenBuffers(1, &colorsVbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorsVbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(id.colors), gl.Ptr(id.colors), gl.STATIC_DRAW)

		const _1d = 1
		gl.VertexAttribIPointer(1, _1d, gl.UNSIGNED_INT, 0, nil)
		gl.EnableVertexAttribArray(1)
	}
}

func (id *itemDrawer) draw() {
	if id.drawItem.VertecesChanged {
		id._updateVertexBuffer()
		id._recomputeVerteces()
		id.drawItem.VertecesChanged = false
	}
	gl.BindVertexArray(id.vao)
	gl.DrawArrays(id._GLDrawType(), 0, int32(len(id.verteces)))
}

func (id *itemDrawer) _GLDrawType() uint32 {
	switch id.drawItem.Type {
	case scene.DrawLines:
		return gl.LINES
	case scene.DrawRects:
		return gl.TRIANGLES
	default:
		panic(fmt.Sprintf("Not supported draw mode: %d", id.drawItem.Type))
	}
}

func (id *itemDrawer) _drawMode() uint32 {
	switch id.drawItem.DrawMode {
	case scene.DynamicDrawMode:
		return gl.DYNAMIC_DRAW
	default:
		return gl.STATIC_DRAW
	}
}

func (id *itemDrawer) _recomputeVerteces() {
	div := id.drawItem.Verteces
	switch id.drawItem.Type {
	case scene.DrawLines:
		id.verteces = div
	case scene.DrawRects:
		if id.verteces == nil {
			id.verteces = make([]float32, 3*len(div))
		}
		for divI := 0; divI < len(div); divI += 4 {
			// 2,5 _____ 6
			//   |\    |
			//   | \   |
			//   |  \  |
			//	 |	 \ |
			//  1|____\|3,4
			//
			rectI := divI * 3
			left, right, bottom, top := div[divI], div[divI+1], div[divI+2], div[divI+3]
			id.verteces[rectI] = left
			id.verteces[rectI+1] = bottom

			id.verteces[rectI+2] = left
			id.verteces[rectI+3] = top

			id.verteces[rectI+4] = right
			id.verteces[rectI+5] = bottom

			id.verteces[rectI+6] = right
			id.verteces[rectI+7] = bottom

			id.verteces[rectI+8] = left
			id.verteces[rectI+9] = top

			id.verteces[rectI+10] = right
			id.verteces[rectI+11] = top
		}
	default:
		slog.Warn("Not supported draw type", "drawType", id.drawItem.Type)
	}
}

func (id *itemDrawer) _computeColors() {
	colors := id.drawItem.Colors
	if len(colors) == 0 {
		return
	}
	switch id.drawItem.Type {
	case scene.DrawLines:
		id.colors = make([]uint32, 0, 2*len(colors))
		for _, color := range colors {
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
		}
	case scene.DrawRects:
		id.colors = make([]uint32, 0, 6*len(colors))
		for _, color := range colors {
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
			id.colors = append(id.colors, color)
		}
	default:
		slog.Warn("Not supported draw type", "drawType", id.drawItem.Type)
	}
}

func (id *itemDrawer) _updateVertexBuffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, id.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, Float32ByteSize*len(id.verteces), gl.Ptr(&id.verteces[0]), id._drawMode())
}
