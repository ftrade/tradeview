package opengl

import (
	"log/slog"
	"os"

	"github.com/go-gl/gl/all-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Program struct {
	ID       uint32
	matrixID int32
}

func (p *Program) Validate() {
	gl.ValidateProgram(p.ID)
	var params int32 = -1
	gl.GetProgramiv(p.ID, gl.VALIDATE_STATUS, &params)
	if params != gl.TRUE {
		log := GetInfoLog(p.ID, gl.GetProgramiv, gl.GetProgramInfoLog)
		slog.Info(log)
		os.Exit(1)
	}
}

func (p *Program) InitUniformMatrix() {
	matrixName := gl.Str("matrix\x00")
	p.matrixID = gl.GetUniformLocation(p.ID, matrixName)
}

func (p *Program) UpdateMatrix(m mgl.Mat4) {
	gl.UniformMatrix4fv(p.matrixID, 1, false, &m[0])
}
