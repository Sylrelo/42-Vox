#version 400

out float 	FragColor;
in	vec2 	TexCoords;

uniform sampler2D 	texPosition;
uniform sampler2D 	texNormal;

float scale = 1;
float bias = 0.05;
float intensity = 1;
int iterations = 4;

vec2 vec[4] = vec2[](
    vec2(1.0, 0.0),
    vec2(-1.0, 0.0),
    vec2(0.0, 1.0),
    vec2(0.0, -1.0)
);

float rand(vec2 co){
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

float doAmbientOcclusion(vec2 tcoord, vec2 uv, vec3 p, vec3 cnorm)
{
    vec3 diff = texture(texPosition, tcoord + uv).xyz - p;

    vec3 v = normalize(diff);
    float d = length(diff) * scale;
    return max(0.0, dot(cnorm, v) - bias ) * (1.0 / (1.0 + d)) * intensity;
}

void	main() {
    float ao    = 0.0f;
    vec3 pos    = texture(texPosition, TexCoords).xyz;
    vec3 normal = texture(texNormal, TexCoords).xyz;
    vec2 rnd    = normalize(vec2(rand(pos.xy), rand(normal.xy)));
    float rad   = 1.0/  pos.z;

    for (int j = 0; j < iterations; ++j)
    {
      vec2 coord1 = reflect(vec[j], rnd) * rad;
      vec2 coord2 = vec2(coord1.x * 0.707 - coord1.y * 0.707, coord1.x * 0.707 + coord1.y * 0.707);
      
      ao += doAmbientOcclusion(TexCoords, coord1 * 0.25, pos, normal);
      ao += doAmbientOcclusion(TexCoords, coord2 * 0.5, pos, normal);
      ao += doAmbientOcclusion(TexCoords, coord1 * 0.75, pos, normal);
      ao += doAmbientOcclusion(TexCoords, coord2, pos, normal);
    }

    ao /= float(iterations) * 4.0;
    FragColor = 1 - ao;
}