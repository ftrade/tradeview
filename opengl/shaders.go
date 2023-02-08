package opengl

const (
	vertexShaderSource = `
		#version 460
		layout(location = 0) in vec2 position;
		layout(location = 1) in uint vertex_colour;
		uniform mat4 matrix;
		out vec3 colour;
		void main() {
			if (vertex_colour == 1) {
				colour = vec3(0.898, 0.274, 0.282);
			} else if (vertex_colour == 2) {
				colour = vec3(0.258, 0.650, 0.513);
			} else if (vertex_colour == 3) {
				colour = vec3(0.3, 0.3, 0.3);
			} else {
				colour = vec3(0.5, 0.5, 0.5);
			}
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
