package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Firefly struct {
	Render bool

	position  *mgl32.Vec2
	arcOrigin *mgl32.Vec2
	arcRadius float64
	angleDec  float64
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

	fx := float64(f.arcOrigin.X()) + f.arcRadius*math.Cos(0)
	fy := float64(f.arcOrigin.Y()) + f.arcRadius*math.Sin(0)

	*f.position = mgl32.Vec2{float32(fx), float32(fy)}
}
