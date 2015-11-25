// Minimal fragment shader

#version 330

// Global constants (for this vertex shader)
vec4 specular_colour = vec4(1.0, 0.8, 0.6, 1.0);
vec4 global_ambient = vec4(0.05, 0.05, 0.05, 1.0);
int  shininess = 15;

// Inputs from the vertex shader
in vec3 light_normal, light_direction, light_position;
in vec4 diffuse_color_mod;

uniform uint emitmode;

// Output pixel fragment colour
out vec4 outputColor;
void main()
{
	vec4 emissive = vec4(0); // Create a vec4(0, 0, 0) for our emmissive light
	vec4 color_ambient = diffuse_color_mod * 0.2;
	vec4 color_specular =  vec4(1.0, 1.0, 0.5, 1.0);
	float light_distance = length(light_direction);

	// Normalise interpolated vectors
	vec3 L = normalize(light_direction);
	vec3 N = normalize(light_normal);

	// Calculate the diffuse component
	vec4 diffuse = max(dot(N, L), 0.0) * diffuse_color_mod;

	// Calculate the specular component using Phong specular reflection
	vec3 V = normalize(-light_position.xyz);
	vec3 R = reflect(-L, N);
	vec4 specular = pow(max(dot(R, V), 0.0), shininess) * color_specular;

    // Attenuation formula from:
    // http://gamedev.stackexchange.com/questions/56897/glsl-light-attenuation-color-and-intensity-formula
    float radius = 10.5;

    // Calculate the attenuation factor;
    float attenuation = clamp(1.0 - light_distance * light_distance / (radius * radius), 0.0, 1.0);
    attenuation *= attenuation;
	
	// simple hack to make the light brighter, it would be better to change the attenuation equation!
	attenuation *= 3.5;

	// If emitmode is 1 then we enable emmissive lighting
	if (emitmode == uint(1)) emissive = vec4(1.0, 1.0, 0.8, 1.0);

	// Calculate the output colour, includung attenuation on the diffuse and specular components
	// Note that you may want to exclude the ambient form the attenuation factor so objects
	// are always visible, or include a global ambient
	outputColor = attenuation * (color_ambient + diffuse + specular) + emissive + global_ambient;
}
