package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui/debug"
	"golang.org/x/image/colornames"
)

type Window struct {
	root         *pixelgl.Window
	worldCanvas  *pixelgl.Canvas
	debugCanvas  *pixelgl.Canvas
	shaderCanvas *pixelgl.Canvas
}

func NewWindow() *Window {
	return &Window{}
}

func (w *Window) Draw(worldDraw func(*pixelgl.Canvas)) {
	w.worldCanvas.Clear(colornames.Lightgray)
	w.debugCanvas.Clear(colornames.Lightgray)
	w.shaderCanvas.Clear(colornames.Lightgray)

	debug.NewOverlay()

	// render world onto own canvas
	worldDraw(w.worldCanvas)
	// paint world canvas onto shader canvas to apply current shader
	w.worldCanvas.Draw(w.shaderCanvas, pixel.IM.Moved(w.root.Bounds().Center()))

	// paint finished shader canvas onto debug canvas
	w.shaderCanvas.Draw(w.debugCanvas, pixel.IM.Moved(w.root.Bounds().Center()))

	// draw finished debug canvas onto window
	w.debugCanvas.Draw(w.root, pixel.IM.Moved(w.root.Bounds().Center()))
}
