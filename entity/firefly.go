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
	f.angleDec += 1
	if f.angleDec >= 360 {
		f.angleDec = 0
	}

	fx := f.arcOrigin.X() + f.arcRadius*float32(math.Cos(float64(decToRad(f.angleDec))))
	fy := f.arcOrigin.Y() + f.arcRadius*float32(math.Sin(float64(decToRad(f.angleDec))))

	*f.position = mgl32.Vec2{fx, fy}
	// fmt.Printf("Firefly position: %f, %f\n", f.position.X(), f.position.Y())
}

func (f *Firefly) Pos() *mgl32.Vec2 {
	return f.position
}

func (f *Firefly) UniformName() string {
	return "fireflyPositions"
}

func decToRad(dec float32) float32 {
	return dec * (math.Pi / 180)
}
