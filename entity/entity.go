package entity

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Entity interface {
	Update()
	Pos() *mgl32.Vec2
	UniformName() string
}
