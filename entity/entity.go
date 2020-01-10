package entity

import (
	"github.com/faiface/pixel/pixelgl"
)

type Entity interface {
	Draw(*pixelgl.Canvas, float64, bool, bool, bool, bool)
	loadSprites()
}
