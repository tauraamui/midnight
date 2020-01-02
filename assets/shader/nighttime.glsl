
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float ambientLightIntensity;

vec3 lightColor = vec3(0.9373, 0.7765, 0.502);

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	vec3 tColor = vec3(texture(uTexture, t).rgb);

	float pct = 0.0;
	pct = distance(t, vec2(0.65));
	tColor = vec3(((1 - pct) * lightColor) * tColor);

	fragColor = vec4(tColor, 1.0);
}