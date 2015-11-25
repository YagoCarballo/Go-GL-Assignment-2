#version 330
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform vec4 lightpos;
in vec3 position;
in vec2 texcoord;
out vec2 fragTexCoord;
out vec3 lightPoint;
void main() {
    lightPoint = lightpos.xyz;
    fragTexCoord = texcoord;
    gl_Position = projection * view * model * vec4(position, 1);
}