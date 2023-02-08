package scene

import (
	"github.com/go-gl/gl/all-core/gl"
)

type CrossLines struct {
	vbo    uint32
	vao    uint32
	points [8]float32
}

func NewCrossLines() *CrossLines {
	cl := &CrossLines{}
	colors := [4]Color{
		DarkGrey, DarkGrey, DarkGrey, DarkGrey,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cl.points), gl.Ptr(&cl.points[0]), gl.DYNAMIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var colorVbo uint32
	gl.GenBuffers(1, &colorVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cl.points), gl.Ptr(&colors[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorVbo)
	gl.VertexAttribIPointer(1, 2, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	cl.vao = vao
	cl.vbo = vbo
	return cl
}

func (cl *CrossLines) Update(x, y, width, height float32) {
	cl.points = [8]float32{
		// horizontal line
		0, y, width, y,
		// vertical line
		x, 0, x, height,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, cl.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cl.points), gl.Ptr(&cl.points[0]), gl.STATIC_DRAW)
}

func (cl *CrossLines) Draw() {
	gl.BindVertexArray(cl.vao)
	gl.DrawArrays(gl.LINES, 0, 4)
}
