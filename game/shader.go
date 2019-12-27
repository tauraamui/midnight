package game

import "os"

import "bytes"

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

type NightTimeShader struct {
	srcPath string
	src     string
}

func NewNightTimeShader() *NightTimeShader {
	return &NightTimeShader{
		srcPath: "./assets/shader/nighttime.glsl",
	}
}

func (s *NightTimeShader) Code() (string, error) {
	if len(s.src) > 0 {
		return s.src, nil
	}

	file, err := os.Open(s.srcPath)
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

func (s *NightTimeShader) Dirty() bool {
	return false
}
