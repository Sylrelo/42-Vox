#version 400 core

layout (location = 0) in int position;
layout (location = 1) in int faceData;

out vec4         		frag_position;
out vec2				frag_texcoords;
out vec3				frag_normal;
out float				billow_value;
out float				tex_id;

uniform					mat4 view;

const int V3_FACE_LEFT = 0;
const int V3_FACE_RIGHT = 1;
const int V3_FACE_TOP = 2;
const int V3_FACE_BOTTOM = 3;
const int V3_FACE_FRONT = 4;
const int V3_FACE_BACK = 5;
const int V3_SPRITE_RIGHT = 6;
const int V3_SPRITE_LEFT = 7;

vec3 getNormal(int face) {
	switch (face) {
		case V3_FACE_LEFT:
			return vec3(-1, 0, 0);
		case V3_FACE_RIGHT:
			return vec3(1, 0, 0);
		case V3_FACE_TOP:
			return vec3(0, 1, 0);
		case V3_FACE_BOTTOM:
			return vec3(0, -1, 0);
		case V3_FACE_FRONT:
			return vec3(0, 0, 1);
		case V3_FACE_BACK:
			return vec3(0, 0, -1);
		case V3_SPRITE_LEFT:
			return vec3(-1, 0, -1);
		case V3_SPRITE_RIGHT:
			return vec3(-1, 0, 1);
	}
}

void main()
{
	vec3 pos = vec3(
		float(((position) >> 15) & 0xFF) - 100, 
		float((position >> 23) & 0x1FF), 
		float(((position) >> 7) & 0xFF) - 100 
	);

	vec2 texCoords = vec2(
		(float((faceData >> 26) & 63)),
		(float((faceData >> 20) & 63))
	);
	
	tex_id 			= faceData >> 5 & 0x7F;
	billow_value 	= ((faceData >> 12) & 0xFF) / 255.0;
	
	if (((faceData >> 2) & 0x7) == V3_FACE_BOTTOM)
		billow_value -= .07;
	if (((faceData >> 2) & 0x7) == V3_FACE_LEFT)
		billow_value -= .03;
	if (((faceData >> 2) & 0x7) == V3_FACE_FRONT)
		billow_value -= .07;
	if (((faceData >> 2) & 0x7) == V3_FACE_BACK)
		billow_value -= .05;
	if (((faceData >> 2) & 0x7) == V3_FACE_RIGHT)
		billow_value -= .05;

	frag_texcoords 	= texCoords;
	frag_position 	= vec4(pos, 1);
	frag_normal 	= getNormal((faceData >> 2) & 0x7);
	gl_Position		= view * vec4(pos, 1);
}