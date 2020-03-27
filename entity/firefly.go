package entity

import (
	"math"

	"github.com/MichaelTJones/pcg"
	"github.com/go-gl/mathgl/mgl32"
)

var pcg32 = pcg.NewPCG32()

type Firefly struct {
	Render bool

	progressiveCurve        float32
	position                *mgl32.Vec2
	angleDec                float32
	angleIncr               float32
	angleDiffSinceDirChange float32
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
	f.updateDir()
	angleRad := decToRad(f.angleDec)
	dirVector := mgl32.Vec2{float32(math.Cos(angleRad)), float32(math.Sin(angleRad))}
	dirVector = dirVector.Normalize().Mul(0.0013)
	*f.position = f.position.Add(dirVector)
}

func (f *Firefly) Pos() *mgl32.Vec2 {
	return f.position
}

func (f *Firefly) UniformName() string {
	return "fireflyPositions"
}

func (f *Firefly) updateDir() {
	if f.angleDiffSinceDirChange > 50 {
		if randBool() {
			f.angleIncr *= -1
		}
		f.angleDiffSinceDirChange = 0
	}
	f.angleDiffSinceDirChange += func() float32 {
		if f.angleDec < 0 {
			return f.angleDec * -1
		}
		return f.angleDec
	}()
	f.angleDec += f.angleIncr
}

func decToRad(dec float32) float64 {
	return float64(dec * (math.Pi / 180))
}

func randBool() bool { return pcg32.Random()&0x01 == 0 }
