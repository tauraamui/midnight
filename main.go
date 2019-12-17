package main

import (
	"image"
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

	world := game.NewWorld(loadTerrainSprites())
	gamepad := game.NewGamepad(win)

	last := time.Now()
	lastFPSCheck := time.Now()
	frameCount := 0
	currentFPS := 0

	fpsText := text.New(
		pixel.V(0, 0),
		text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, game.SCALE*8), text.ASCII, text.RangeTable(unicode.Latin),
		),
	)

	gamepadText := text.New(
		pixel.V(0, 0),
		text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, game.SCALE*3), text.ASCII, text.RangeTable(unicode.Latin),
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
			fpsText.Draw(
				win, pixel.IM.Moved(world.Camera.Unproject(pixel.V(fpsText.TabWidth, win.Bounds().H()-(fpsText.LineHeight+(win.Bounds().H()/13))))),
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
			gamepadText.Draw(
				win, pixel.IM.Moved(world.Camera.Unproject(pixel.V(20, win.Bounds().H()/60))),
			)
		}

		win.Update()
		<-fps
	}
}

func loadTerrainSprites() (pixel.Picture, []pixel.Rect) {
	s, err := loadSpritesheet("assets/terrain.png")
	if err != nil {
		panic(err)
	}

	var terrainTiles []pixel.Rect
	for x := s.Bounds().Min.X; x < s.Bounds().Max.X; x += 32 {
		for y := s.Bounds().Min.Y; y < s.Bounds().Max.Y; y += 32 {
			terrainTiles = append(terrainTiles, pixel.R(x, y, x+32, y+32))
		}
	}

	return s, terrainTiles
}

func loadSpritesheet(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
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
