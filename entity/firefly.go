package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Firefly struct {
	Render bool

	progressiveCurve float32
	position         *mgl32.Vec2
	angleRad         float64
	angleIncrementor float64
}

func NewFirefly(x, y float32) *Firefly {
	f := &Firefly{
		Render:           true,
		position:         &mgl32.Vec2{x, y},
		angleRad:         0,
		angleIncrementor: 0,
	}
	return f
}

func (f *Firefly) Update() {
	dirVector := mgl32.Vec2{float32(math.Cos(f.angleRad)), float32(math.Sin(f.angleRad))}
	dirVector = dirVector.Normalize().Mul(0.001)
	*f.position = f.position.Add(dirVector)
	f.angleRad += f.angleIncrementor
}

func (f *Firefly) Pos() *mgl32.Vec2 {
	return f.position
}

func (f *Firefly) UniformName() string {
	return "fireflyPositions"
}
