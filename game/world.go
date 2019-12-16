package game

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	TILE_SIZE         = 32
	SCALE             = 3
	GRASS_TILES_START = 9
	GRASS_TILES_END   = 11
)

func computeTilesToRender(rec pixel.Rect) (int, int) {
	return int(rec.W() / (32 * SCALE)), int(rec.H() / (32 * SCALE))
}

func NewWorld(s pixel.Picture, tt []pixel.Rect) *World {
	world := &World{
		spriteSheet: s,
	}

	for i := GRASS_TILES_START; i < GRASS_TILES_END; i++ {
		world.grassTiles = append(world.grassTiles, pixel.NewSprite(world.spriteSheet, tt[i]))
	}

	return world
}

type World struct {
	spriteSheet pixel.Picture
	grassTiles  []*pixel.Sprite
}

func (w *World) Draw(win *pixelgl.Window) {
	maxX, maxY := computeTilesToRender(win.Bounds())
	for x := 0; x < maxX+5; x++ {
		for y := 0; y < maxY+5; y++ {
			rand.Seed(int64((x * 83) + (y * 385)))
			grass := w.grassTiles[ranRange(0, 2)]
			grass.Draw(win, pixel.IM.Scaled(pixel.ZV, SCALE).Moved(pixel.V(float64(x*(32*SCALE)), float64(y*(32*SCALE)))))
		}
	}
}

func ranRange(min, max int) int { return rand.Intn(max-min) + min }
