package opengl

import (
	"fmt"
	"log/slog"

	"github.com/go-gl/gl/all-core/gl"
)

// MakeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32, colors []uint32) uint32 {
	var pointsVbo uint32
	const count = 1
	gl.GenBuffers(count, &pointsVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, pointsVbo)
	gl.BufferData(gl.ARRAY_BUFFER, Float32ByteSize*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(count, &vao)
	gl.BindVertexArray(vao)
	const _2d = 2
	gl.VertexAttribPointer(0, _2d, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	if len(colors) > 0 {
		var colorsVbo uint32
		gl.GenBuffers(count, &colorsVbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorsVbo)
		gl.BufferData(gl.ARRAY_BUFFER, Float32ByteSize*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)
		gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
		gl.EnableVertexAttribArray(1)
	}
	return vao
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

// InitOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info("OpenGL", "version", version)
}

func MakeProgram() Program {
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
	return prog
}
