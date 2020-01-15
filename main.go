package main

import (
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/game"
)

const (
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
)

func run() {
	fps := time.Tick(time.Second / 60)
	second := time.Tick(time.Second)
	game := game.NewInstance(makeGLWindow())

	framesInSecond := 0
	for !game.Exiting() {
		game.Update()
		game.Draw()

		<-fps
		select {
		case <-second:
			game.SetCurrentFPS(framesInSecond + 1)
			framesInSecond = 0
		default:
			framesInSecond++
		}
	}
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
