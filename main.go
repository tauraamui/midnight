package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
	SCALE      = 6
)

func run() {
	win := makeGLWindow()
	fps := time.Tick(time.Second / 60)

	spritesheet, err := loadSpritesheet("assets/terrain.png")
	if err != nil {
		panic(err)
	}

	var terrainTiles []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			terrainTiles = append(terrainTiles, pixel.R(x, y, x+32, y+32))
		}
	}

	grass := pixel.NewSprite(spritesheet, terrainTiles[10])
	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		tilesToRenderX := int(win.Bounds().W()) / (32 * SCALE)
		tilesToRenderY := int(win.Bounds().H()) / (32 * SCALE)

		for x := 0; x < tilesToRenderX+5; x++ {
			for y := 0; y < tilesToRenderY+5; y++ {
				grass.Draw(win, pixel.IM.Scaled(pixel.ZV, SCALE).Moved(pixel.V(float64(x*(32*SCALE)), float64(y*(32*SCALE)))))
			}
		}

		win.Update()

		<-fps
	}
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
