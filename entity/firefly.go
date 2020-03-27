package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Firefly struct {
	Render bool

	progressiveCurve float32
	position         *mgl32.Vec2
	angleDec         float32
	angleIncr        float32
}

func NewFirefly(x, y float32) *Firefly {
	f := &Firefly{
		Render:    true,
		position:  &mgl32.Vec2{x, y},
		angleDec:  0,
		angleIncr: 1,
	}
	return f
}

func (f *Firefly) Update() {
	angleRad := decToRad(f.angleDec)
	dirVector := mgl32.Vec2{float32(math.Cos(angleRad)), float32(math.Sin(angleRad))}
	dirVector = dirVector.Normalize().Mul(0.0013)
	*f.position = f.position.Add(dirVector)

	if f.angleDec >= 0 && f.angleDec < 360 {
		f.angleIncr *= -1
	}

	if f.angleDec > 360 {
		f.angleDec = 0
	}

	f.angleDec += f.angleIncr
}

func (f *Firefly) Pos() *mgl32.Vec2 {
	return f.position
}

func (f *Firefly) UniformName() string {
	return "fireflyPositions"
}

func decToRad(dec float32) float64 {
	return float64(dec * (math.Pi / 180))
}
