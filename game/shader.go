package game

import (
	"bytes"
	"fmt"
	"strings"

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
	AmbientLightIntensity *float32
	FireflyPositions      []*mgl32.Vec2
	CamPos                *mgl32.Vec2

	src   string
	dirty bool
}

func NewDayAndNightTimeShader() *DayAndNightTimeShader {
	shader := DayAndNightTimeShader{
		AmbientLightIntensity: new(float32),
		CamPos:                &mgl32.Vec2{0, 0},
		FireflyPositions: []*mgl32.Vec2{
			&mgl32.Vec2{0.65, 0.65},
			&mgl32.Vec2{0.9, 0.9},
			&mgl32.Vec2{0.1, 0.9},
			&mgl32.Vec2{0.3, 0.4},
		},
		dirty: true,
	}
	*shader.AmbientLightIntensity = INTENSITY_PER_MINUTE
	return &shader
}

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
	s.src = strings.Replace(s.src, "//FIREFLY_POSITION_UNIFORMS", convFireflyPositions(s.FireflyPositions), -1)
	return s.src, nil
}

func (s *DayAndNightTimeShader) Dirty() bool {
	wasDirty := s.dirty
	s.dirty = false
	return wasDirty
}

func convFireflyPositions(list []*mgl32.Vec2) string {
	var buff strings.Builder
	buff.WriteString(fmt.Sprintf("const int FIREFLIES_COUNT = %d;\n", len(list)))
	buff.WriteString("uniform vec2[FIREFLIES_COUNT] fireflyPositions;\n")

	return buff.String()
}
