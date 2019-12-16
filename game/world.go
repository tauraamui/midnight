package game

import (
	"math/rand"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	TILE_SIZE         = 32
	SCALE             = 2
	CAM_SPEED         = 80.0
	GRASS_TILES_START = 9
	GRASS_TILES_END   = 11
)

func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

func computeTilesToRender(rec pixel.Rect) (int, int) {
	return int(rec.W() / 32), int(rec.H() / 32)
}

func NewWorld(s pixel.Picture, tt []pixel.Rect) *World {
	world := &World{
		spriteSheet: s,
		FPSText: text.New(
			pixel.V(0, 0),
			text.NewAtlas(
				ttfFromBytesMust(goregular.TTF, 48), text.ASCII, text.RangeTable(unicode.Latin),
			),
		),
		camPos: pixel.ZV,
	}

	for i := GRASS_TILES_START; i < GRASS_TILES_END; i++ {
		world.grassTiles = append(world.grassTiles, pixel.NewSprite(world.spriteSheet, tt[i]))
	}

	return world
}

type World struct {
	FPSText     *text.Text
	spriteSheet pixel.Picture
	grassTiles  []*pixel.Sprite
	camPos      pixel.Vec
}

func (w *World) Update(pressed func(pixelgl.Button) bool, dt float64) {
	if pressed(pixelgl.KeyA) {
		w.camPos.X -= CAM_SPEED * dt
	}

	if pressed(pixelgl.KeyD) {
		w.camPos.X += CAM_SPEED * dt
	}

	if pressed(pixelgl.KeyW) {
		w.camPos.Y += CAM_SPEED * dt
	}

	if pressed(pixelgl.KeyS) {
		w.camPos.Y -= CAM_SPEED * dt
	}
}

func (w *World) Draw(win *pixelgl.Window) {
	cam := pixel.IM.Scaled(w.camPos, SCALE).Moved(win.Bounds().Center().Sub(w.camPos))
	win.SetMatrix(cam)
	w.FPSText.Draw(
		win, pixel.IM.Moved(cam.Unproject(pixel.V(0, win.Bounds().H()-(w.FPSText.LineHeight+(w.FPSText.LineHeight/1.8))))),
	)
}

func ranRange(min, max int) int { return rand.Intn(max-min) + min }
