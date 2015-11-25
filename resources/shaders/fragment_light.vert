// Veretx shader: When completed should implement
// per vertex lighting

// Specify minimum OpenGL version
#version 330

// Define the vertex attributes
layout(location = 0) in vec3 position;
layout(location = 1) in vec4 colour;
layout(location = 2) in vec3 normal;

// Outputs to send to the fragment shader
out vec3 light_normal;
out vec3 light_direction, light_position;
out vec4 diffuse_color_mod;

// These are the uniforms that are defined in the application
uniform mat4 model, view, projection;
uniform mat3 normalmatrix;
uniform uint colourmode;
uniform vec4 lightpos;

void main()
{
	vec4 position_h = vec4(position, 1.0);	// Convert the (x,y,z) position to homogeneous coords (x,y,z,w)
	vec3 light_pos3 = lightpos.xyz;

	// Switch the vertex colour based on the colourmode
	if (colourmode == uint(1))
		diffuse_color_mod = colour;
	else
		diffuse_color_mod = vec4(1.0, 0, 0, 1.0);

	// Define our vectors for calculating diffuse and specular lighting
	mat4 mv_matrix  = view * model;						// Calculate the model-view transformation
	light_position  = (mv_matrix * position_h).xyz;		// Modify the vertex position (x, y, z, w) by the model-view transformation
	light_normal    = normalize(normalmatrix * normal);	// Modify the normals by the normal-matrix (i.e. to model-view (or eye) coordinates )
	light_direction = light_pos3 - light_position;	    // Calculate the vector from the light position to the vertex in eye space

	// Calculate the vertex position in projectin space and output to the pipleline using the reserved variable gl_Position
	gl_Position = (projection * view * model) * position_h;
}
