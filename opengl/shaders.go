package opengl

const (
	vertexShaderSource = `
		#version 460
		layout(location = 0) in vec2 position;
		layout(location = 1) in uint vertex_colour;
		uniform mat4 matrix;
		uniform vec3 colors[%d];
		out vec3 colour;
		void main() {
			colour = colors[vertex_colour];
			gl_Position =  matrix * vec4(position, 0.0, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 460
		out vec4 frag_colour;
		in vec3 colour;
		void main() {
			frag_colour = vec4(colour, 1.0);
		}
	` + "\x00"
)
