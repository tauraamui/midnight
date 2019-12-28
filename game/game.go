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
	perfNumbersText     *text.Text
	worldClockText      *text.Text

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
		perfNumbersText:     text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		worldClockText:      text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
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

func (do *debugOverlay) draw(win *pixelgl.Canvas, gp *Gamepad, c *WorldClock, fps int) {
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

	do.worldClockText.Clear()
	_, err = do.worldClockText.WriteString(c.Current.String())
	if err != nil {
		panic(err)
	}
	do.worldClockText.Draw(win, pixel.IM.Scaled(
		pixel.ZV, 1.2,
	).Moved(pixel.V(win.Bounds().W()-200, 10*1.2)))
}

type Instance struct {
	FPS                   int
	CurrentFramesInSecond int

	rootWin      *pixelgl.Window
	worldCanvas  *pixelgl.Canvas
	debugCanvas  *pixelgl.Canvas
	shaderCanvas *pixelgl.Canvas
	lastShader   Shader
	fullscreen   bool
	gamepad      *Gamepad
	world        *World

	debugOverlay    *debugOverlay
	lastGamepadScan time.Time
	lastDelta       time.Time
}

func NewInstance(win *pixelgl.Window) *Instance {
	return &Instance{
		FPS:                   0,
		CurrentFramesInSecond: 0,

		rootWin:      win,
		worldCanvas:  pixelgl.NewCanvas(win.Bounds()),
		debugCanvas:  pixelgl.NewCanvas(win.Bounds()),
		shaderCanvas: pixelgl.NewCanvas(win.Bounds()),
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

	shader := i.world.Update(i.gamepad, deltaTime)
	if i.lastShader != shader {
		shaderSrc, err := shader.Code()
		if err != nil {
			panic(err)
		}

		i.shaderCanvas.SetFragmentShader(shaderSrc)
		i.lastShader = shader
	}
}

func (i *Instance) Draw() {
	win := i.rootWin
	beforeFinishedDraw := time.Now()

	win.Clear(colornames.Black)

	i.worldCanvas.Clear(colornames.Lightgray)
	i.debugCanvas.Clear(colornames.Lightgray)
	i.shaderCanvas.Clear(colornames.Lightgray)

	// render world onto own canvas
	i.world.Draw(i.worldCanvas)
	// paint world canvas onto shader canvas to apply current shader
	i.worldCanvas.Draw(i.shaderCanvas, pixel.IM.Moved(win.Bounds().Center()))
	i.debugOverlay.drawTimeDuration = time.Since(beforeFinishedDraw)

	// paint finished shader canvas onto debug canvas
	i.shaderCanvas.Draw(i.debugCanvas, pixel.IM.Moved(win.Bounds().Center()))
	// draw debug elements above everything else
	i.debugOverlay.draw(i.debugCanvas, i.gamepad, i.world.Clock, i.FPS)
	// draw finished debug canvas onto window
	i.debugCanvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

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
	defer func() {
		i.rootWin.Update()
		i.worldCanvas.SetBounds(i.rootWin.Canvas().Bounds())
		i.debugCanvas.SetBounds(i.rootWin.Canvas().Bounds())
		i.shaderCanvas.SetBounds(i.rootWin.Canvas().Bounds())
	}()
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
