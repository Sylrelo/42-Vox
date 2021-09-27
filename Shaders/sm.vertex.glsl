#version 400 core

layout (location = 0) in int position;
uniform mat4    light;

void main()
{
    vec3 pos = vec3(
		float((position >> 15) & 0xFF) - 100, 
		float((position >> 23) & 0x1FF), 
		float((position >> 7) & 0xFF) - 100
	);
    gl_Position = light * vec4(pos, 1.0);
}