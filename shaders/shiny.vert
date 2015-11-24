//
// Adapted from the book (Phong light chapter)
// http://www.amazon.com/Antons-OpenGL-Tutorials-Anton-Gerdelan-ebook/dp/B00LAMQYF2
//

#version 330

// Define the vertex attributes
layout(location = 0) in vec3 position;
layout(location = 1) in vec4 colour;
layout(location = 2) in vec3 normal;

// These are the uniforms that are defined in the application
uniform mat4 model, view, projection;
uniform uint colourmode, emitmode;
uniform vec4 lightpos;

out vec3 position_eye, normal_eye;
out vec4 diffuse_color_mod;

void main () {
	// Switch the vertex colour based on the colourmode
	if (colourmode == uint(1))
		diffuse_color_mod = colour;
	else
		diffuse_color_mod = vec4(1.0, 0, 0, 1.0);

	position_eye = normalize(vec3 (lightpos - (view * model * vec4 (position, 1.0))));
	normal_eye = normalize(vec3 (lightpos - (view * model * vec4 (normal, 0.0))));
	gl_Position = projection * view * model * vec4 (position, 1.0);
}