#version 330

uniform sampler2D DiffuseTextureSampler;
uniform sampler2D NormalTextureSampler;
uniform sampler2D SpecularTextureSampler;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(DiffuseTextureSampler, fragTexCoord);
}