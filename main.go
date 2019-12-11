package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

func run() {
	fps := time.Tick(time.Second / 60)
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Midnight",
			Bounds: pixel.R(0, 0, float64(800), float64(600)),
		},
	)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Black)
		win.Update()
		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
