package opengl

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	vertexShaderSource = `
		#version 460
		layout(location = 0) in vec2 position;
		layout(location = 1) in uint vertex_colour;
		out vec3 colour;
		void main() {
			if (vertex_colour == 1) {
				colour = vec3(0.898, 0.274, 0.282);
			} else if (vertex_colour == 2) {
				colour = vec3(0.258, 0.650, 0.513);
			} else {
				colour = vec3(0.5, 0.5, 0.5);
			}
			gl_Position = vec4(position, 0.0, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 460
		out vec4 frag_colour;
		in vec3 colour;
		void main() {
			frag_colour = vec4(colour, 1.0);
			// frag_colour = vec4(0.5, 0.5, 0.5, 1.0);
		}
	` + "\x00"
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
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}
