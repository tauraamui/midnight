package main

import (
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/game"
)

const (
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
)

func run() {
	fps := time.Tick(time.Second / 60)
	second := time.Tick(time.Second)
	game := game.NewInstance(makeGLWindow())
	game.SetShader(fragmentShader)
	for !game.Quitted() {
		game.Update()
		game.Draw()

		<-fps
		select {
		case <-second:
			game.FPS = game.CurrentFramesInSecond + 1
			game.CurrentFramesInSecond = 0
		default:
			game.CurrentFramesInSecond++
		}
	}
}
func makeGLWindow() *pixelgl.Window {
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Midnight",
			Bounds: pixel.R(0, 0, float64(WIN_WIDTH), float64(WIN_HEIGHT)),
		},
	)
	if err != nil {
		panic(err)
	}
	return win
}

func main() {
	pixelgl.Run(run)
}

var fragmentShader = `
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

float ambientStrength = 0.225f;
vec3 lightColor = vec3(1.0, 1.0, 1.0);

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	fragColor = vec4(((ambientStrength * lightColor) * texture(uTexture, t).rgb), 1.0);
}
`
