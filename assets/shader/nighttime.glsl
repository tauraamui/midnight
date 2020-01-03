
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float ambientLightIntensity;

vec3 ambientLightColor = vec3(1.0);
vec3 spotLightColor = vec3(0.8431, 0.0588, 0.0588);

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	vec3 tColor = texture(uTexture, t).rgb;
	float distanceFromLight = 1-(distance(t, vec2(0.85)));
	vec3 ambColor = (ambientLightIntensity + (distanceFromLight * spotLightColor)) * tColor;
	fragColor = vec4(ambColor, 1.0);
}