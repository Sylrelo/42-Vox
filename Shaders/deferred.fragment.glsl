#version 400 core

layout (location = 0) out vec3 gPosition;
layout (location = 1) out vec3 gNormal;
layout (location = 2) out vec4 gAlbedoSpec;

in vec4                 frag_position;
in vec2                 frag_texcoords;
in vec3                 frag_normal;
in float                tex_id;
in float                billow_value;

uniform sampler2DArray 	textures;
uniform sampler2D       texShadowmap;

uniform					mat4 light;
uniform					mat4 view;
uniform					mat4 matWorld;

float random(vec3 seed, int i){
	vec4 seed4 			= vec4(seed,i);
	float dot_product 	= dot(seed4, vec4(12.9898,78.233,45.164,94.673));
	return fract(sin(dot_product) * 43758.5453);
}

vec2 poissonDisk[16] = vec2[]( 
	vec2( -0.94201624, -0.39906216 ), 
	vec2( 0.94558609, -0.76890725 ), 
	vec2( -0.094184101, -0.92938870 ), 
	vec2( 0.34495938, 0.29387760 ), 
	vec2( -0.91588581, 0.45771432 ), 
	vec2( -0.81544232, -0.87912464 ), 
	vec2( -0.38277543, 0.27676845 ), 
	vec2( 0.97484398, 0.75648379 ), 
	vec2( 0.44323325, -0.97511554 ), 
	vec2( 0.53742981, -0.47373420 ), 
	vec2( -0.26496911, -0.41893023 ), 
	vec2( 0.79197514, 0.19090188 ), 
	vec2( -0.24188840, 0.99706507 ), 
	vec2( -0.81409955, 0.91437590 ), 
	vec2( 0.19984126, 0.78641367 ), 
	vec2( 0.14383161, -0.14100790 ) 
);

// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-16-shadow-mapping/#going-further
float   in_shadow(vec4 frag_light_pos)
{
    float shadow        = 0.0;
    float bias          = 0.003;
    vec3 coords         = (frag_light_pos.xyz / frag_light_pos.w) * 0.5 + 0.5;
	float visibility 	= 1.0;
    // shadow 			= coords.z - bias > texture(texShadowmap, coords.xy).r ? .4: 1.0;
	for (int i=0; i < 16; i++) {
		if ( texture( texShadowmap, coords.xy + poissonDisk[i] / 700.0 ).r < coords.z - bias ) {
			visibility -= 0.02;
		}
	}
	// if (visibility <= .7)
	// 	visibility = .7;
    return visibility;
}

void main()
{
    vec4 colorFromTexture = texture(textures, vec3(frag_texcoords.xy, tex_id));

 	if (colorFromTexture.a == 0.0) {
        discard;
    }

	float shadow    		= in_shadow(light * frag_position);
	vec3 viewSpacePosition 	= vec3(matWorld * frag_position);
    
    gPosition 		= viewSpacePosition;
    gNormal 		= mat3(matWorld) * frag_normal;
    gAlbedoSpec     = vec4(colorFromTexture.xyz * billow_value * shadow, 1.0);
}  