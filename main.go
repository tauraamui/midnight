package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/game"
	"golang.org/x/image/colornames"
)

const (
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
)

func run() {
	win := makeGLWindow()
	fps := time.Tick(time.Second / 60)

	world := game.NewWorld(loadTerrainSprites())

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		world.Draw(win)
		win.Update()
		<-fps
	}
}

func loadTerrainSprites() (pixel.Picture, []pixel.Rect) {
	s, err := loadSpritesheet("assets/terrain.png")
	if err != nil {
		panic(err)
	}

	var terrainTiles []pixel.Rect
	for x := s.Bounds().Min.X; x < s.Bounds().Max.X; x += 32 {
		for y := s.Bounds().Min.Y; y < s.Bounds().Max.Y; y += 32 {
			terrainTiles = append(terrainTiles, pixel.R(x, y, x+32, y+32))
		}
	}

	return s, terrainTiles
}

func loadSpritesheet(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func makeGLWindow() *pixelgl.Window {
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Midnight",
			Bounds: pixel.R(0, 0, float64(WIN_WIDTH), float64(WIN_HEIGHT)),
		},
	)
	if err != nil {
		panic(err)
	}
	return win
}

func main() {
	pixelgl.Run(run)
}
