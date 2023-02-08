package opengl

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/all-core/gl"
)

// makeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32, colors []uint32) uint32 {
	var pointsVbo uint32
	gl.GenBuffers(1, &pointsVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, pointsVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var colorsVbo uint32
	if len(colors) > 0 {
		gl.GenBuffers(1, &colorsVbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorsVbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, pointsVbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	if len(colors) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, colorsVbo)
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

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

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

	progId := gl.CreateProgram()
	gl.AttachShader(progId, vertexShader)
	gl.AttachShader(progId, fragmentShader)
	gl.LinkProgram(progId)

	prog := Program{
		Id: progId,
	}
	return prog
}
