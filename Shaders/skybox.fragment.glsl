#version 400 core

out vec4 FragColor;
in vec3 TexCoords;

uniform samplerCube     texSkybox;

void main()
{    
    FragColor = texture(texSkybox, TexCoords);
}
