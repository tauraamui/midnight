package game

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui"
)

type Instance struct {
	currentFPS              int
	currentFramesInSecond   int
	window                  *ui.Window
	world                   *World
	lastDelta               time.Time
	lastTimeBeforeAllUpdate time.Time
}

func NewInstance(win *pixelgl.Window) *Instance {
	return &Instance{
		window: ui.NewWindow(win, SCALE),
		world:  NewWorld(),

		lastDelta: time.Now(),
	}
}

func (i *Instance) Update() {
	i.lastTimeBeforeAllUpdate = time.Now()
	dt := time.Since(i.lastDelta).Seconds()
	i.world.Update(i.window.Update(
		i.currentFPS, i.currentFramesInSecond, time.Since(i.lastTimeBeforeAllUpdate),
	), float64(dt))
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

func (i *Instance) SetCurrentFramesInSecond(frames int) { i.currentFramesInSecond = frames }
