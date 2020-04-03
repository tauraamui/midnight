package entity

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
)

type Entity interface {
	Update()
	Draw(*pixelgl.Canvas, pixel.Matrix)
	Pos() *mgl32.Vec2
	UniformName() string
}
