#version 330

in vec4 fAmbient, fDiffuse, fSpecular, fEmissive;
in vec3 fNormals, fLightDirection;

out vec4 outputColor;

int  shininess = 15;

void main() {
    vec3 L = normalize(fLightDirection);
    vec3 N = normalize(fNormals);

    vec3 V = normalize(-fLightDirection.xyz);
    vec3 R = reflect(-L, N);
    vec4 specular = pow(max(dot(R, V), 0.0), shininess) * fSpecular;

    outputColor = ((fDiffuse + specular) + fEmissive);
}