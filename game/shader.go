package game

import (
	"bytes"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/markbates/pkger"
)

type Shader interface {
	Code() (string, error)
	Dirty() bool
}

type PassthroughShader struct {
	src string
}

func NewPassthroughShader() *PassthroughShader {
	return &PassthroughShader{
		src: `
			#version 330 core
			// The first line in glsl source code must always start with a version directive as seen above.

			// vTexCoords are the texture coordinates, provided by Pixel
			in vec2  vTexCoords;

			// fragColor is what the final result is and will be rendered to your screen.
			out vec4 fragColor;

			// uTexBounds is the texture's boundries, provided by Pixel.
			uniform vec4 uTexBounds;

			// uTexture is the actualy texture we are sampling from, also provided by Pixel.
			uniform sampler2D uTexture;

			void main() {
				// t is set to the screen coordinate for the current fragment.
				vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
				// And finally, we're telling the GPU that this fragment should be the color as sampled from our texture.
				fragColor = texture(uTexture, t);
			}
		`,
	}
}

func (s *PassthroughShader) Code() (string, error) { return s.src, nil }

func (s *PassthroughShader) Dirty() bool { return false }

type DayAndNightTimeShader struct {
	ambientLightIntensity *float32
	fireflyPos            *mgl32.Vec2
	src                   string
	dirty                 bool
}

func NewDayAndNightTimeShader() *DayAndNightTimeShader {
	fireflyPos := mgl32.Vec2{0.85, 0.85}
	shader := DayAndNightTimeShader{
		ambientLightIntensity: new(float32),
		fireflyPos:            &fireflyPos,
		dirty:                 true,
	}
	*shader.ambientLightIntensity = INTENSITY_PER_MINUTE
	return &shader
}

func (s *DayAndNightTimeShader) SetAmbientLightIntensity(ali float32) {
	*s.ambientLightIntensity = ali
	s.dirty = true
}

func (s *DayAndNightTimeShader) SetFireflyPos(pos mgl32.Vec2) {
	*s.fireflyPos = pos
	s.dirty = true
}

func (s *DayAndNightTimeShader) AmbientLightIntensity() *float32 { return s.ambientLightIntensity }
func (s *DayAndNightTimeShader) FireflyPos() *mgl32.Vec2         { return s.fireflyPos }

func (s *DayAndNightTimeShader) Code() (string, error) {
	if len(s.src) > 0 {
		return s.src, nil
	}

	file, err := pkger.Open("/assets/shader/nighttime.glsl")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var bytes bytes.Buffer
	_, err = bytes.ReadFrom(file)
	if err != nil {
		return "", err
	}

	s.src = bytes.String()
	return s.src, nil
}

func (s *DayAndNightTimeShader) Dirty() bool {
	wasDirty := s.dirty
	s.dirty = false
	return wasDirty
}
