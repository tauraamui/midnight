package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

const (
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
	TILE_SIZE  = 100
)

type Tile struct {
	width  int
	height int
}

func run() {
	fps := time.Tick(time.Second / 60)
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Midnight",
			Bounds: pixel.R(0, 0, float64(WIN_WIDTH), float64(WIN_HEIGHT)),
		},
	)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	tilesToRenderX := WIN_WIDTH / TILE_SIZE
	tilesToRenderY := WIN_HEIGHT / TILE_SIZE

	println("Tiles column:", tilesToRenderX, "Tiles row:", tilesToRenderY)

	for !win.Closed() {
		win.Clear(colornames.Black)

		for i := 0; i < tilesToRenderX; i++ {
			for j := 0; j < tilesToRenderY; j++ {
				tileX := i * TILE_SIZE
				tileY := j * TILE_SIZE
				imd.Push(pixel.V(float64(tileX), float64(tileY)), pixel.V(float64(tileX+TILE_SIZE), float64(tileY+TILE_SIZE)))
				imd.EndShape = imdraw.SharpEndShape
				imd.Color = colornames.Cyan
				// imd.Push(pixel.V(float64(i), float64(j)), pixel.V(float64(i+TILE_SIZE), float64(j+TILE_SIZE)))
				imd.Rectangle(1)
			}
		}
		imd.Draw(win)

		win.Update()
		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
