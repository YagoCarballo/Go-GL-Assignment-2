#version 330

layout(location = 0) in vec3 position;
layout(location = 2) in vec3 normal;

uniform mat4 model, view, projection;
uniform vec4 ambient, diffuse, specular, emissive;
uniform vec4 lightpos;

out vec4 fAmbient, fDiffuse, fSpecular, fEmissive;
out vec3 fNormals, fLightDirection;

void main() {
    fAmbient = ambient;
    fDiffuse = diffuse;
    fSpecular = specular;
    fEmissive = emissive;

    fNormals = normalize(mat3(model - view) * normal);
    fLightDirection = normalize(lightpos.xyz);

    gl_Position = projection * view * model * vec4(position, 1.0);
}