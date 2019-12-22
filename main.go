package main

import (
	"os"
	"strconv"
	"time"
	"unicode"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/tauraamui/midnight/game"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	WIN_WIDTH        = 800
	WIN_HEIGHT       = 600
	DEBUG_TEXT_SCALE = 3
)

var debugMode = false
var fullscreen = false
var worldUpdateDuration time.Duration = time.Second
var worldDrawDuration time.Duration = time.Second

func run() {
	var (
		win             = makeGLWindow()
		fps             = time.Tick(time.Second / 60)
		frames          = 0
		framesPerSecond = 0
		second          = time.Tick(time.Second)
		last            = time.Now()
		lastGamepadScan = time.Now()

		world   = game.NewWorld()
		gamepad = game.NewGamepad(win)

		fpsText = text.New(
			pixel.V(0, 0),
			text.NewAtlas(
				ttfFromBytesMust(goregular.TTF, game.SCALE*8), text.ASCII, text.RangeTable(unicode.Latin),
			),
		)

		gamepadText = text.New(
			pixel.V(0, 0),
			text.NewAtlas(
				ttfFromBytesMust(goregular.TTF, game.SCALE*8), text.ASCII, text.RangeTable(unicode.Latin),
			),
		)

		timeGraph = game.NewGraph(win)
	)

	for !win.Closed() {
		if win.JustReleased(pixelgl.KeyEscape) {
			win.Destroy()
			os.Exit(0)
		}

		if win.JustReleased(pixelgl.KeyF) {
			toggleFullscreen(win)
		}

		if win.JustReleased(pixelgl.KeyX) {
			debugMode = !debugMode
			fpsText.Clear()
		}

		deltaTime := time.Since(last).Seconds()
		last = time.Now()

		if time.Since(lastGamepadScan).Seconds() >= 3 {
			if !gamepad.JoystickAttached() {
				gamepad.AttachToJoystick()
			}
			lastGamepadScan = time.Now()
		}

		win.Clear(color.RGBA{R: 110, G: 201, B: 57, A: 255})

		beforeWorldUpdateTime := time.Now()
		world.Update(gamepad, deltaTime)
		worldUpdateDuration = time.Since(beforeWorldUpdateTime)

		beforeWorldDrawTime := time.Now()
		world.Draw(win)
		worldDrawDuration = time.Since(beforeWorldDrawTime)

		if debugMode {
			win.SetMatrix(pixel.IM)
			fpsText.Draw(
				win,
				pixel.IM.Scaled(
					pixel.ZV, DEBUG_TEXT_SCALE,
				).Moved(pixel.V(20, (win.Bounds().H()-(fpsText.LineHeight*DEBUG_TEXT_SCALE))-10)),
			)
			fpsText.Clear()
			_, err := fpsText.WriteString(strconv.Itoa(framesPerSecond))
			if err != nil {
				panic(err)
			}

			gamepadText.Clear()
			_, err = gamepadText.WriteString(gamepad.Debug())
			if err != nil {
				panic(err)
			}
			gamepadText.Draw(win, pixel.IM.Scaled(
				pixel.ZV, DEBUG_TEXT_SCALE,
			).Moved(pixel.V(15, 10*DEBUG_TEXT_SCALE)))

			timeGraph.Draw(win)
		}

		win.Update()
		<-fps

		frames++

		if frames < len(timeGraph.TimesPerFrame) {
			timeGraph.TimesPerFrame[frames] = game.TimeSpent{
				DrawTime:   worldDrawDuration,
				UpdateTime: worldUpdateDuration,
			}
		} else {
			timeGraph.TimesPerFrame = append(timeGraph.TimesPerFrame, game.TimeSpent{
				DrawTime:   worldDrawDuration,
				UpdateTime: worldUpdateDuration,
			})
		}
		select {
		case <-second:
			framesPerSecond = frames
			frames = 0
			timeGraph.TimesPerFrame = []game.TimeSpent{}
		default:
		}
	}
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

func toggleFullscreen(win *pixelgl.Window) {
	fullscreen = !fullscreen
	var mon *pixelgl.Monitor = nil
	if fullscreen {
		mon = getMonitor()
	}
	win.SetMonitor(mon)
}

func makeGLWindow() *pixelgl.Window {
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Midnight",
			Bounds: pixel.R(0, 0, float64(WIN_WIDTH), float64(WIN_HEIGHT)),
		},
	)
	if err != nil {
		panic(err)
	}
	return win
}

func getMonitor() *pixelgl.Monitor {
	for _, mon := range pixelgl.Monitors() {
		if mon.Name() == "UMC SHARP" {
			return mon
		}
	}
	return pixelgl.PrimaryMonitor()
}

func main() {
	pixelgl.Run(run)
}
