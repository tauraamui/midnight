package game

import (
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	DEBUG_TEXT_SCALE = 3
)

type debugOverlay struct {
	perfGraph           *Graph
	fpsText             *text.Text
	gamepadAxisListText *text.Text

	enabled            bool
	drawTimeDuration   time.Duration
	updateTimeDuration time.Duration
}

func NewDebugOverlay(win *pixelgl.Window) *debugOverlay {
	ttf := ttfFromBytesMust(goregular.TTF, SCALE*8)
	return &debugOverlay{
		perfGraph: NewGraph(win),

		enabled:             false,
		fpsText:             text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		gamepadAxisListText: text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
	}
}

func (do *debugOverlay) update(win *pixelgl.Window, frames int) {
	if win.JustReleased(pixelgl.KeyX) {
		do.enabled = !do.enabled
		if !do.enabled {
			do.fpsText.Clear()
			do.perfGraph.TimesPerFrame = []TimeSpent{}
		}
	}

	if frames < len(do.perfGraph.TimesPerFrame) {
		do.perfGraph.TimesPerFrame[frames] = TimeSpent{
			DrawTime:   do.drawTimeDuration,
			UpdateTime: do.updateTimeDuration,
		}

		if frames+1 < len(do.perfGraph.TimesPerFrame) {
			do.perfGraph.TimesPerFrame[frames+1] = TimeSpent{}
		}
	} else {
		do.perfGraph.TimesPerFrame = append(do.perfGraph.TimesPerFrame, TimeSpent{
			DrawTime:   do.drawTimeDuration,
			UpdateTime: do.updateTimeDuration,
		})
	}
}

func (do *debugOverlay) draw(win *pixelgl.Window, gp *Gamepad, fps int) {
	if !do.enabled {
		return
	}

	win.SetMatrix(pixel.IM)

	do.perfGraph.Draw(win)

	do.fpsText.Draw(
		win,
		pixel.IM.Scaled(
			pixel.ZV, DEBUG_TEXT_SCALE,
		).Moved(pixel.V(20, (win.Bounds().H()-(do.fpsText.LineHeight*DEBUG_TEXT_SCALE))-10)),
	)
	do.fpsText.Clear()
	_, err := do.fpsText.WriteString(strconv.Itoa(fps))
	if err != nil {
		panic(err)
	}

	do.gamepadAxisListText.Clear()
	_, err = do.gamepadAxisListText.WriteString(gp.Debug())
	if err != nil {
		panic(err)
	}
	do.gamepadAxisListText.Draw(win, pixel.IM.Scaled(
		pixel.ZV, DEBUG_TEXT_SCALE,
	).Moved(pixel.V(15, 10*DEBUG_TEXT_SCALE)))

}

type Instance struct {
	FPS                   int
	CurrentFramesInSecond int

	rootWin    *pixelgl.Window
	fullscreen bool
	gamepad    *Gamepad
	world      *World

	debugOverlay    *debugOverlay
	lastGamepadScan time.Time
	lastDelta       time.Time
}

func NewInstance(win *pixelgl.Window) *Instance {
	return &Instance{
		FPS:                   0,
		CurrentFramesInSecond: 0,

		rootWin:      win,
		gamepad:      NewGamepad(win),
		world:        NewWorld(),
		debugOverlay: NewDebugOverlay(win),

		fullscreen:      false,
		lastGamepadScan: time.Now(),
		lastDelta:       time.Now(),
	}
}

func (i *Instance) Update() {
	beforeFinishedUpdate := time.Now()
	defer func() { i.debugOverlay.updateTimeDuration = time.Since(beforeFinishedUpdate) }()

	deltaTime := time.Since(i.lastDelta).Seconds()
	i.lastDelta = time.Now()

	i.attachToGamepad()
	win := i.rootWin

	i.debugOverlay.update(win, i.CurrentFramesInSecond)

	if win.JustReleased(pixelgl.KeyEscape) {
		win.Destroy()
		os.Exit(0)
	}

	if win.JustReleased(pixelgl.KeyF) {
		i.toggleFullscreen()
	}

	i.world.Update(i.gamepad, deltaTime)
}

func (i *Instance) Draw() {
	beforeFinishedDraw := time.Now()
	defer func() { i.debugOverlay.drawTimeDuration = time.Since(beforeFinishedDraw) }()
	win := i.rootWin
	win.Clear(colornames.Lightgray)
	i.world.Draw(win)
	i.debugOverlay.draw(win, i.gamepad, i.FPS)

	i.rootWin.Update()
}

func (i *Instance) Quitted() bool {
	return i.rootWin.Closed()
}

func (i *Instance) attachToGamepad() {
	if time.Since(i.lastGamepadScan).Seconds() >= 2 {
		if !i.gamepad.JoystickAttached() {
			i.gamepad.AttachToJoystick()
		}
		i.lastGamepadScan = time.Now()
	}
}

func (i *Instance) toggleFullscreen() {
	i.fullscreen = !i.fullscreen
	var mon *pixelgl.Monitor = nil
	if i.fullscreen {
		mon = getMonitor()
	}
	i.rootWin.SetMonitor(mon)
}

func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

func getMonitor() *pixelgl.Monitor {
	for _, mon := range pixelgl.Monitors() {
		if mon.Name() == "UMC SHARP" {
			return mon
		}
	}
	return pixelgl.PrimaryMonitor()
}
