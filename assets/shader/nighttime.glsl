
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float ambientLightIntensity;

//FIREFLY_POSITION_UNIFORMS

float attenConst = 1.0;
float attenLinear = 50.4;
float attenQuadratic = 129.6;

vec3 ambientLightColor = vec3(1.0, 1.0, 1.0);
vec3 spotLightColor = vec3(0.7137, 0.8431, 0.0588);

float getLightAtten(vec2 pos, vec2 t) {
	float distanceFromLight = distance(t, pos);
	return 1.0 / (attenConst + (attenLinear * distanceFromLight) + (attenQuadratic * pow(distanceFromLight, 2)));
}

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	vec3 tColor = texture(uTexture, t).rgb;

	vec3 ambientLight = (ambientLightIntensity * ambientLightColor);
	vec3 pointlightLight = (getLightAtten(fireflyPositions[0], t) * spotLightColor);

	vec3 ambColor = ambientLight;

	for (int i = 0; i < fireflyPositions.length(); i++) {
		ambColor += getLightAtten(fireflyPositions[i], t) * spotLightColor;
	}

	ambColor = ambColor * tColor;
	// vec3 ambColor = (ambientLight + pointlightLight) * tColor;
	fragColor = vec4(ambColor, 1.0);
}