#version 330

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 position;
in vec2 texcoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = texcoord;
    gl_Position = projection * view * model * vec4(position, 1);
}