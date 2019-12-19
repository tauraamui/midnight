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
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
)

var debugMode = false
var fullscreen = false

func run() {
	win := makeGLWindow()
	fps := time.Tick(time.Second / 60)

	world := game.NewWorld()
	gamepad := game.NewGamepad(win)

	last := time.Now()
	lastFPSCheck := time.Now()
	lastGamepadScan := time.Now()
	frameCount := 0
	currentFPS := 0

	fpsText := text.New(
		pixel.V(0, 0),
		text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, game.SCALE*4), text.ASCII, text.RangeTable(unicode.Latin),
		),
	)

	gamepadText := text.New(
		pixel.V(0, 0),
		text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, game.SCALE*1.5), text.ASCII, text.RangeTable(unicode.Latin),
		),
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
		world.Update(gamepad, deltaTime)
		world.Draw(win)

		if debugMode {
			frameCount++
			if time.Since(lastFPSCheck).Seconds() >= 1 {
				currentFPS = frameCount
				frameCount = 0
				lastFPSCheck = time.Now()
			}

			win.SetMatrix(pixel.IM)
			fpsText.Draw(
				win, pixel.IM.Moved(pixel.V(15, win.Bounds().H()-fpsText.LineHeight)),
			)
			fpsText.Clear()
			_, err := fpsText.WriteString(strconv.Itoa(currentFPS))
			if err != nil {
				panic(err)
			}

			gamepadText.Clear()
			_, err = gamepadText.WriteString(gamepad.Debug())
			if err != nil {
				panic(err)
			}
			gamepadText.Draw(win, pixel.IM.Moved(pixel.V(15, 10)))
		}

		win.Update()
		<-fps
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
