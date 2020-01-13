package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Firefly struct {
	Render bool

	position  *mgl32.Vec2
	arcOrigin *mgl32.Vec2
	arcRadius float32
	angleDec  float32
}

func NewFirefly(x, y float32) *Firefly {
	return &Firefly{
		Render:    true,
		position:  &mgl32.Vec2{x, y},
		arcOrigin: &mgl32.Vec2{x - 1, y - 1},
		arcRadius: 5,
		angleDec:  0,
	}
}

func (f *Firefly) Update() {
	f.angleDec += 0.01
	if f.angleDec >= 360 {
		f.angleDec = 0
	}

	fx := f.arcOrigin.X() + f.arcRadius*float32(math.Cos(0))
	fy := f.arcOrigin.Y() + f.arcRadius*float32(math.Sin(0))

	*f.position = mgl32.Vec2{fx, fy}
}
