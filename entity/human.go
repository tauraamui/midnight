package entity

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Human struct {
	pos pixel.Vec
}

func (h *Human) Draw(
	win *pixelgl.Canvas,
	animSpeed float64,
	movingL, movingR, movingU, movingD bool,
) {

}
