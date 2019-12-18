package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/sprite"
)

const (
	TILE_SIZE         = 32
	SCALE             = 3
	CAM_SPEED         = 160.0
	GRASS_TILES_START = 9
	GRASS_TILES_END   = 11
)

type World struct {
	spriteSheet pixel.Picture
	grassTiles  []*pixel.Sprite
	camPos      pixel.Vec
	Camera      pixel.Matrix
}

func NewWorld(s pixel.Picture) *World {
	world := World{
		spriteSheet: s,
		camPos:      pixel.ZV,
		Camera:      pixel.IM,
	}

	world.grassTiles = append(world.grassTiles, sprite.Make(world.spriteSheet, 1, 0))

	return &world
}

func (w *World) Update(gp *Gamepad, dt float64) {
	if speed, movingL := gp.MovingLeft(); movingL {
		w.camPos.X -= (CAM_SPEED * speed) * dt
	}

	if speed, movingR := gp.MovingRight(); movingR {
		w.camPos.X += (CAM_SPEED * speed) * dt
	}

	if speed, movingU := gp.MovingUp(); movingU {
		w.camPos.Y += (CAM_SPEED * speed) * dt
	}

	if speed, movingD := gp.MovingDown(); movingD {
		w.camPos.Y -= (CAM_SPEED * speed) * dt
	}
}

func (w *World) Draw(win *pixelgl.Window) {
	w.Camera = pixel.IM.Scaled(w.camPos, SCALE).Moved(win.Bounds().Center().Sub(w.camPos))
	win.SetMatrix(w.Camera)

	for x := 0; x < 15; x++ {
		for y := 0; y < 15; y++ {
			grass := w.grassTiles[0]
			grass.Draw(win, pixel.IM.Moved(pixel.V(float64(x*(32)), float64(y*(32)))))
		}
	}
}
