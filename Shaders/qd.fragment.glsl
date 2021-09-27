#version 400 core

out vec4    FragColor;
in vec2     TexCoords;

uniform sampler2D   texColor;
uniform sampler2D   texSsao;
uniform mat4        projection;

void main()
{
    vec3 frag_color         = texture(texColor, TexCoords).rgb;
    float occlusion_factor  = texture(texSsao, TexCoords).r;
    
    FragColor = vec4(frag_color * occlusion_factor , 1.0);
}