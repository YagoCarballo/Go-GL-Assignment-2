#version 330

uniform sampler2D DiffuseTextureSampler;
uniform sampler2D NormalTextureSampler;
uniform sampler2D SpecularTextureSampler;

in vec4 lightPosition;
in vec3 lightNormal, lightDirection;
in vec2 fragTexCoord;

out vec4 outputColor;

// Global constants (for this vertex shader)
const vec4 colorAmbientGlobal   = vec4(0.05, 0.05, 0.05, 0.0);
const vec4 colorEmissive        = vec4(0.2, 0.2, 0.2, 0.0);
const int  shininess            = 9000;
const float radius              = 90.5;

void main() {
    vec4 colorDiffuse = texture(DiffuseTextureSampler, fragTexCoord);
    vec4 colorAmbient = vec4(colorDiffuse.xyz * 0.2, 0.0);
    vec4 colorSpecular =  vec4(1.0, 1.0, 0.5, 0.0);

    float lightDistance = length(lightDirection);

    // Normalise interpolated vectors
    vec3 L = normalize(lightDirection);
    vec3 N = normalize(lightNormal);

    // Calculate the diffuse component
    vec4 diffuse = max(dot(N, L), 0.0) * colorDiffuse;

    // Calculate the specular component using Phong specular reflection
    vec3 V = normalize(-lightPosition.xyz);
    vec3 R = reflect(-L, N);
    vec4 specular = pow(max(dot(R, V), 0.0), shininess) * colorSpecular;

    // Attenuation formula from:
    // http://gamedev.stackexchange.com/questions/56897/glsl-light-attenuation-color-and-intensity-formula

    // Calculate the attenuation factor;
    float attenuation = clamp(1.0 - lightDistance * lightDistance / (radius * radius), 0.0, 1.0);
    attenuation *= attenuation;

    // simple hack to make the light brighter, it would be better to change the attenuation equation!
    // attenuation *= 3.5;

    // Calculate the output colour, includung attenuation on the diffuse and specular components
    // Note that you may want to exclude the ambient form the attenuation factor so objects
    // are always visible, or include a global ambient
    outputColor = colorAmbient + (attenuation * (diffuse + specular)) + colorEmissive + colorAmbientGlobal;
}