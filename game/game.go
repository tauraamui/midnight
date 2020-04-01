package game

import (
	"time"

	"github.com/tauraamui/midnight/ui/shader"

	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/ui"
)

type Instance struct {
	currentFPS            int
	currentFramesInSecond int

	window *ui.Window
	world  *World
	shader *shader.Shader

	lastDelta               time.Time
	lastTimeBeforeAllUpdate time.Time
	updateDuration          time.Duration
}

func NewInstance(win *pixelgl.Window, fullscreen *bool) *Instance {
	return &Instance{
		window: ui.NewWindow(win, SCALE, *fullscreen),
		world:  NewWorld(),

		lastDelta: time.Now(),
	}
}

func (i *Instance) Update() {
	dt := time.Since(i.lastDelta).Seconds()
	gp := i.window.Update(i.currentFPS, i.currentFramesInSecond, i.updateDuration)
	i.lastTimeBeforeAllUpdate = time.Now()
	i.shader = i.world.Update(gp, dt)
	i.updateDuration = time.Since(i.lastTimeBeforeAllUpdate)
	i.lastDelta = time.Now()
}

func (i *Instance) Draw() {
	i.window.SetShader(i.shader)
	i.window.Draw(i.world.Draw)
}

func (i *Instance) Exiting() bool {
	return i.window.Closing()
}

func (i *Instance) SetCurrentFPS(fps int) {
	i.currentFPS = fps
}

func (i *Instance) SetCurrentFramesInSecond(frames int) { i.currentFramesInSecond = frames }
