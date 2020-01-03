
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float ambientLightIntensity;

vec3 lightColor = vec3(0.9373, 0.7647, 0.4667);
vec3 ambientLightColor = vec3(1.0);

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	vec3 tColor = vec3(texture(uTexture, t).rgb);
	vec3 ambColor = (ambientLightIntensity * ambientLightColor) * tColor;
	vec3 spotlightColor = vec3(0.0);

	if (ambientLightIntensity < 0.9) {
		float pct = 0.0;
		float radius = 0.9;
		pct = distance(t, vec2(0.55));
		pct = smoothstep(radius-(radius*0.9), radius+(radius*0.9), dot(pct, pct)*4.0);
		spotlightColor = vec3(((1 - pct) * ambientLightColor * lightColor) * tColor);
		if (spotlightColor.x > ambColor.x && spotlightColor.y > ambColor.y && spotlightColor.z > ambColor.z) {
			ambColor = spotlightColor;
		}
	}

	fragColor = vec4(ambColor, 1.0);
}