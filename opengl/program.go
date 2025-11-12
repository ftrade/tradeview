package opengl

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-gl/gl/all-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Program struct {
	ID       uint32
	matrixID int32
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		log := GetInfoLog(shader, gl.GetShaderiv, gl.GetShaderInfoLog)

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func MakeProgram(colors []float32) Program {
	colorsN := len(colors) / 3
	vertexShaderSource := fmt.Sprintf(vertexShaderSource, colorsN)
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	progID := gl.CreateProgram()
	gl.AttachShader(progID, vertexShader)
	gl.AttachShader(progID, fragmentShader)
	gl.LinkProgram(progID)

	prog := Program{
		ID: progID,
	}
	prog.setupColors(colors, int32(colorsN))
	return prog
}

func (p *Program) setupColors(colors []float32, count int32) {
	colorsName := gl.Str("colors\x00")
	colorsId := gl.GetUniformLocation(p.ID, colorsName)
	gl.UseProgram(p.ID)
	gl.Uniform3fv(colorsId, count, &colors[0])
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
