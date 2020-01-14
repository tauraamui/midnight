package game

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui"
)

type Instance struct {
	window    *ui.Window
	world     *World
	lastDelta time.Time
}

func NewInstance(win *pixelgl.Window) *Instance {
	return &Instance{
		window: ui.NewWindow(win),
		world:  NewWorld(),

		lastDelta: time.Now(),
	}
}

func (i *Instance) Update() {
	dt := time.Since(i.lastDelta).Milliseconds()
	i.world.Update(i.window.Update(), float64(dt))
	i.lastDelta = time.Now()
}

func (i *Instance) Draw() {
	i.window.Draw(i.world.Draw)
}

func (i *Instance) Exiting() bool {
	return i.window.Closing()
}
