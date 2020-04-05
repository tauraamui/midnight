package entity

import (
	"math"

	"github.com/MichaelTJones/pcg"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
)

var pcg32 = pcg.NewPCG32()

type Firefly struct {
	Render bool

	imd                     *imdraw.IMDraw
	position                *mgl32.Vec2
	radius                  float64
	resolution              float32
	spread                  float32
	angleDec                float32
	angleIncr               float32
	angleDiffSinceDirChange float32
}

func NewFirefly(x, y float32) *Firefly {
	f := &Firefly{
		Render:     true,
		position:   &mgl32.Vec2{x, y},
		radius:     60,
		resolution: 64,
		spread:     2 * math.Pi,
		angleDec:   200,
		angleIncr:  1,
	}
	return f
}

func (f *Firefly) Update() {
	f.updateDir()
	angleRad := decToRad(f.angleDec)
	dirVector := mgl32.Vec2{float32(math.Cos(angleRad)), float32(math.Sin(angleRad))}
	dirVector = dirVector.Normalize().Mul(0.1)
	*f.position = f.position.Add(dirVector)
}

func (f *Firefly) Draw(win *pixelgl.Canvas, mat pixel.Matrix, center pixel.Vec, src *pixelgl.Canvas) {
	if f.imd == nil {
		imd := imdraw.New(nil)
		imd.Color = pixel.Alpha(1)
		imd.Push(pixel.ZV)
		imd.Color = pixel.Alpha(0)
		presicion := f.spread / f.resolution
		for angle := -f.spread / 2; angle <= f.spread/2+(presicion); angle += presicion {
			imd.Push(pixel.V(1, 0).Rotated(float64(angle)))
		}
		imd.Polygon(0)
		f.imd = imd
	}

	// draw the light arc (normally arc is full circle)
	p := pixel.V(float64(f.position.X()), float64(f.position.Y()))
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV, f.radius).Moved(mat.Project(p)))
	win.SetColorMask(pixel.RGB(.1, .1, 0.8).Mul(pixel.Alpha(0.9)))
	win.SetComposeMethod(pixel.ComposePlus)
	f.imd.Draw(win)

	// draw the bright world inside the light
	// win.SetMatrix(pixel.IM)
	// win.SetComposeMethod(pixel.ComposeIn)
	// src.Draw(win, pixel.IM.Moved(center))

	// win.SetMatrix(pixel.IM.Scaled(pixel.ZV, f.radius).Moved(mat.Project(p)))
	// win.SetColorMask(pixel.RGB(.1, .1, 0.8).Mul(pixel.Alpha(0.9)))
	// win.SetComposeMethod(pixel.ComposePlus)
	// f.imd.Draw(win)
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
