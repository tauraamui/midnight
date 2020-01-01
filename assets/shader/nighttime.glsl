
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float ambientLightIntensity;

vec3 lightColor = vec3(1.0, 1.0, 1.0);

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	fragColor = vec4(((ambientLightIntensity * lightColor) * texture(uTexture, t).rgb), 1.0);
}