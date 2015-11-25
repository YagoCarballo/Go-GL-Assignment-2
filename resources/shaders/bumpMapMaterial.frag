#version 330
uniform sampler2D DiffuseTextureSampler;
uniform sampler2D NormalTextureSampler;
uniform sampler2D SpecularTextureSampler;
in vec2 fragTexCoord;
in vec3 lightPoint;
out vec4 outputColor;
void main() {
    // Extract the normal from the normal map
    vec3 normal = normalize(texture(NormalTextureSampler, fragTexCoord.st).rgb * 2.0 - 1.0);

    // Determine where the light is positioned (this can be set however you like)
    vec3 light_pos = normalize(lightPoint);
//    vec3 light_pos = normalize(vec3(1.0, 1.0, 1.5));

    // Calculate the lighting diffuse value
    float diffuse = max(dot(normal, light_pos), 0.0);

    vec3 color = diffuse * texture(DiffuseTextureSampler, fragTexCoord.st).rgb;

    outputColor = vec4(color, 1.0);
}