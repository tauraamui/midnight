package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/sprite"
)

const (
	TILE_SIZE         = 32
	SCALE             = 2
	CAM_SPEED         = 160.0
	GRASS_TILES_START = 9
	GRASS_TILES_END   = 11
)

type World struct {
	Camera pixel.Matrix
	Bunny  *Bunny

	camPos                             pixel.Vec
	movingL, movingR, movingU, movingD bool
	spriteSheet                        pixel.Picture
	grassTiles                         []*pixel.Sprite
	currentVelocity                    float64
}

func NewWorld() *World {
	world := World{
		Camera: pixel.IM,
		Bunny:  NewBunny(),
		camPos: pixel.ZV,
	}
	world.loadSprites()

	return &world
}

func (w *World) Update(gp *Gamepad, dt float64) {
	w.currentVelocity = 0
	w.movingL, w.movingR, w.movingU, w.movingD = false, false, false, false

	if speed, movingL := gp.MovingLeft(); movingL {
		w.movingL = movingL
		w.currentVelocity = (CAM_SPEED * speed) * dt
		w.camPos.X -= w.currentVelocity
	}

	if speed, movingR := gp.MovingRight(); movingR {
		w.movingR = movingR
		w.currentVelocity = (CAM_SPEED * speed) * dt
		w.camPos.X += w.currentVelocity
	}

	if speed, movingU := gp.MovingUp(); movingU {
		w.movingU = movingU
		w.currentVelocity = (CAM_SPEED * speed) * dt
		w.camPos.Y += w.currentVelocity
	}

	if speed, movingD := gp.MovingDown(); movingD {
		w.movingD = movingD
		w.currentVelocity = (CAM_SPEED * speed) * dt
		w.camPos.Y -= w.currentVelocity
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

	w.Camera = pixel.IM
	win.SetMatrix(w.Camera)
	w.Bunny.Draw(
		win,
		w.currentVelocity,
		w.movingL, w.movingR, w.movingU, w.movingD,
	)
}

func (w *World) loadSprites() {
	s, err := sprite.LoadSpritesheet("./assets/terrain.png")
	if err != nil {
		panic(err)
	}

	w.spriteSheet = s
	w.grassTiles = append(w.grassTiles, sprite.Make(w.spriteSheet, 3, 0, 32))
}
