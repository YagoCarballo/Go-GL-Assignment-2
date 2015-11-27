// Minimal vertex shader

#version 330

// These are the vertex attributes
layout(location = 0) in vec3 position;
layout(location = 1) in vec4 colour;
layout(location = 2) in vec3 normal;

// Uniform variables are passed in from the application
uniform mat4 model, view, projection;
uniform uint colourmode;
uniform vec4 lightpos;
uniform vec4 tone;

// Outputs
out vec4 lightPosition;
out vec3 lightNormal, lightDirection;
out vec4 colorDiffuse;

// Color Constants
const vec4 toneModifier = vec4(0.662, 0.405, 0.022, 1);

void main() {
    vec3 lightPosV3 = lightpos.xyz;

    // Convert the (x,y,z) position to homogeneous coords (x,y,z,w)
	vec4 positionHomogeneus = vec4(position, 1.0);

    // Update the Diffuse Color
	if (colourmode == uint(0)) {
		colorDiffuse = vec4(0.8, 0.6, 0.2, 1.0);
	} else {
		colorDiffuse = colour + tone;
	}

    // Calculates the Transformations
	mat4 matrixModelView = view * model;
	mat3 matrixNormal = transpose(inverse(mat3(matrixModelView)));

    // Calculates the Lights
	lightPosition = matrixModelView * positionHomogeneus;
	lightNormal = normalize(matrixNormal * -normal);
	lightDirection = lightPosV3 - lightPosition.xyz;

	// Define the vertex position
	gl_Position = (projection * view * model) * positionHomogeneus;
}

