package game

import (
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
	SCALE             = 3
	CAM_SPEED         = 160.0
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

func NewWorld(s pixel.Picture, tt []pixel.Rect) *World {
	world := &World{
		spriteSheet: s,
		FPSText: text.New(
			pixel.V(0, 0),
			text.NewAtlas(
				ttfFromBytesMust(goregular.TTF, SCALE*8), text.ASCII, text.RangeTable(unicode.Latin),
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

func (w *World) Update(gp *Gamepad, dt float64) {
	if gp.MovingLeft() {
		w.camPos.X -= CAM_SPEED * dt
	}

	if gp.MovingRight() {
		w.camPos.X += CAM_SPEED * dt
	}

	if gp.MovingUp() {
		w.camPos.Y += CAM_SPEED * dt
	}

	if gp.MovingDown() {
		w.camPos.Y -= CAM_SPEED * dt
	}
}

func (w *World) Draw(win *pixelgl.Window) {
	cam := pixel.IM.Scaled(w.camPos, SCALE).Moved(win.Bounds().Center().Sub(w.camPos))
	win.SetMatrix(cam)

	for x := 0; x < 15; x++ {
		for y := 0; y < 15; y++ {
			grass := w.grassTiles[1]
			grass.Draw(win, pixel.IM.Moved(pixel.V(float64(x*(32)), float64(y*(32)))))
		}
	}

	w.FPSText.Draw(
		win, pixel.IM.Moved(cam.Unproject(pixel.V(w.FPSText.TabWidth, win.Bounds().H()-(w.FPSText.LineHeight+(win.Bounds().H()/13))))),
	)
}
