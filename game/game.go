package game

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui"
)

type Instance struct {
	currentFPS int
	window     *ui.Window
	world      *World
	lastDelta  time.Time
}

func NewInstance(win *pixelgl.Window) *Instance {
	return &Instance{
		window: ui.NewWindow(win, SCALE),
		world:  NewWorld(),

		lastDelta: time.Now(),
	}
}

func (i *Instance) Update() {
	dt := time.Since(i.lastDelta).Seconds()
	i.world.Update(i.window.Update(i.currentFPS), float64(dt))
	i.lastDelta = time.Now()
}

func (i *Instance) Draw() {
	i.window.Draw(i.world.Draw)
}

func (i *Instance) Exiting() bool {
	return i.window.Closing()
}

func (i *Instance) SetCurrentFPS(fps int) {
	i.currentFPS = fps
}
