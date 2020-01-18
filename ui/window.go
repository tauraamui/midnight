package ui

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui/debug"
	"github.com/tauraamui/midnight/ui/input"
	"golang.org/x/image/colornames"
)

var SCALE float64

type Window struct {
	FPS int

	fullscreen bool
	closing    bool

	root         *pixelgl.Window
	input        *input.Gamepad
	overlay      *debug.DebugOverlay
	worldCanvas  *pixelgl.Canvas
	debugCanvas  *pixelgl.Canvas
	shaderCanvas *pixelgl.Canvas
}

func NewWindow(win *pixelgl.Window, scale float64) *Window {
	SCALE = scale
	debug.SCALE = SCALE
	return &Window{
		root:    win,
		input:   input.NewGamepad(win),
		overlay: debug.NewDebugOverlay(win),

		worldCanvas:  pixelgl.NewCanvas(win.Bounds()),
		debugCanvas:  pixelgl.NewCanvas(win.Bounds()),
		shaderCanvas: pixelgl.NewCanvas(win.Bounds()),
	}
}

func (w *Window) Update(currentFPS, currentFramesInSecond int, updateDuration time.Duration) *input.Gamepad {
	w.FPS = currentFPS
	w.overlay.SetUpdateTimeDuration(updateDuration)
	w.overlay.Update(w.root, currentFramesInSecond)

	w.root.UpdateInput()
	if w.root.JustPressed(pixelgl.KeyEscape) {
		w.closing = true
	}

	if w.root.JustPressed(pixelgl.KeyF) {
		w.toggleFullscreen()
	}

	w.overlay.Update(w.root, w.FPS)

	return w.input
}

func (w *Window) Draw(worldDraw func(*pixelgl.Canvas)) {
	beforeWorldAndShaderDraw := time.Now()
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

	w.overlay.SetDrawTimeDuration(time.Since(beforeWorldAndShaderDraw))

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
