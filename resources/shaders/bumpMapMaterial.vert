#version 330

// These are the vertex attributes
layout(location = 0) in vec3 position;
layout(location = 1) in vec2 texcoord;
layout(location = 2) in vec3 normal;
uniform sampler2D NormalTextureSampler;

// Uniform variables are passed in from the application
uniform mat4 model, view, projection;
uniform vec4 ambient, diffuse, specular, emissive;
uniform uint colourmode;
uniform vec4 lightpos;

// Outputs
out vec4 lightPosition;
out vec3 lightNormal, lightDirection;
out mat3 matrixNormal;
out vec2 textureCoordinates;
out vec4 ambientMaterial, diffuseMaterial, specularMaterial, emissiveMaterial;

void main() {
    ambientMaterial     = ambient;
    diffuseMaterial     = diffuse;
    specularMaterial    = specular;
    emissiveMaterial    = emissive;

    vec3 lightPosV3 = lightpos.xyz;

    // Convert the (x,y,z) position to homogeneous coords (x,y,z,w)
    vec4 positionHomogeneus = vec4(position, 1.0);

    // Calculates the Transformations
    mat4 matrixModelView = view * model;
    matrixNormal = transpose(inverse(mat3(matrixModelView)));

    // Calculates the Lights
    lightPosition = matrixModelView * positionHomogeneus;
    lightDirection = lightPosV3 - lightPosition.xyz;
    lightNormal = normalize(matrixNormal *  normal);

    // Define the vertex position
    gl_Position = (projection * view * model) * positionHomogeneus;


    // Sets the Texture coordinates
    textureCoordinates = texcoord;
}