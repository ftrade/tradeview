package opengl

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Program struct {
	Id       uint32
	matrixId int32
}

func (p *Program) Validate() {
	gl.ValidateProgram(p.Id)
	var params int32 = -1
	gl.GetProgramiv(p.Id, gl.VALIDATE_STATUS, &params)
	if params != gl.TRUE {
		log := GetInfoLog(p.Id, gl.GetProgramiv, gl.GetProgramInfoLog)
		fmt.Print(log)
		os.Exit(1)
	}
}

func (p *Program) InitMatrix(m mgl.Mat3) {
	matrixName := gl.Str("matrix\x00")
	p.matrixId = gl.GetUniformLocation(p.Id, matrixName)
	p.UpdateMatrix(m)
}

func (p *Program) UpdateMatrix(m mgl.Mat3) {
	gl.UseProgram(p.Id)
	gl.UniformMatrix3fv(p.matrixId, 1, false, &m[0])
}
