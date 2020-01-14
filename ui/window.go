package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui/debug"
	"github.com/tauraamui/midnight/ui/input"
	"golang.org/x/image/colornames"
)

type Window struct {
	FPS int

	root         *pixelgl.Window
	input        *input.Gamepad
	fullscreen   bool
	closing      bool
	overlay      *debug.DebugOverlay
	worldCanvas  *pixelgl.Canvas
	debugCanvas  *pixelgl.Canvas
	shaderCanvas *pixelgl.Canvas
}

func NewWindow(win *pixelgl.Window) *Window {
	return &Window{
		root:    win,
		input:   input.NewGamepad(win),
		overlay: debug.NewDebugOverlay(win),
	}
}

func (w *Window) Update() *input.Gamepad {
	w.root.UpdateInput()
	if w.root.Pressed(pixelgl.KeyEscape) {
		w.closing = true
	}
	return w.input
}

func (w *Window) Draw(worldDraw func(*pixelgl.Canvas)) {
	w.root.Clear(colornames.Black)

	w.worldCanvas.Clear(colornames.Lightgray)
	w.debugCanvas.Clear(colornames.Lightgray)
	w.shaderCanvas.Clear(colornames.Lightgray)

	// render world onto own canvas
	worldDraw(w.worldCanvas)
	// paint world canvas onto shader canvas to apply current shader
	w.worldCanvas.Draw(w.shaderCanvas, pixel.IM.Moved(w.root.Bounds().Center()))

	// paint finished shader canvas onto debug canvas
	w.shaderCanvas.Draw(w.debugCanvas, pixel.IM.Moved(w.root.Bounds().Center()))

	// draw debug elements above everything else
	w.overlay.Draw(w.debugCanvas, w.input, w.FPS)

	// draw finished debug canvas onto window
	w.debugCanvas.Draw(w.root, pixel.IM.Moved(w.root.Bounds().Center()))
	w.root.Update()
}

func (w *Window) Fullscreen() bool { return w.fullscreen }

func (w *Window) Closing() bool { return w.closing || w.root.Closed() }

func (w *Window) toggleFullscreen() {
	defer func() {
		w.root.Update()
		w.worldCanvas.SetBounds(w.root.Canvas().Bounds())
		w.debugCanvas.SetBounds(w.root.Canvas().Bounds())
		w.shaderCanvas.SetBounds(w.root.Canvas().Bounds())
	}()
	w.fullscreen = !w.fullscreen
	var mon *pixelgl.Monitor = nil
	if w.fullscreen {
		mon = getMonitor()
	}
	w.root.SetMonitor(mon)
}

func getMonitor() *pixelgl.Monitor {
	for _, mon := range pixelgl.Monitors() {
		if mon.Name() == "UMC SHARP" {
			return mon
		}
	}
	return pixelgl.PrimaryMonitor()
}
