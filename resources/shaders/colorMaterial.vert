#version 330

layout(location = 0) in vec3 position;
layout(location = 2) in vec3 normal;

uniform mat4 model, view, projection;
uniform vec4 ambient, diffuse, specular, emissive;
uniform uint colourmode;
uniform vec4 lightpos;

out vec4 lightPosition;
out vec3 lightNormal, lightDirection;
out vec4 colorAmbient, colorDiffuse, colorSpecular, colorEmissive;

void main() {
    colorAmbient = ambient;
    colorDiffuse = diffuse;
    colorSpecular = specular;
    colorEmissive = emissive;

    vec3 lightPosV3 = lightpos.xyz;

    // Convert the (x,y,z) position to homogeneous coords (x,y,z,w)
    vec4 positionHomogeneus = vec4(position, 1.0);

    // Update the Diffuse Color
    if (colourmode == uint(0)) {
        colorDiffuse = vec4(0.8, 0.6, 0.2, 1.0);
    } else {
        colorDiffuse = diffuse;
    }

    // Calculates the Transformations
    mat4 matrixModelView = view * model;
    mat3 matrixNormal = transpose(inverse(mat3(matrixModelView)));

    // Calculates the Lights
    lightPosition = matrixModelView * positionHomogeneus;
    lightNormal = normalize(matrixNormal * normal);
    lightDirection = lightPosV3 - lightPosition.xyz;

    // Define the vertex position
    gl_Position = (projection * view * model) * positionHomogeneus;
}