package debug

import (
	"strconv"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/tauraamui/midnight/ui/input"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var SCALE float64

const (
	DEBUG_TEXT_SCALE = 3
)

type DebugOverlay struct {
	perfGraph           *Graph
	fpsText             *text.Text
	gamepadAxisListText *text.Text
	perfNumbersText     *text.Text
	worldClockText      *text.Text

	enabled            bool
	drawTimeDuration   time.Duration
	updateTimeDuration time.Duration
}

func NewDebugOverlay(win *pixelgl.Window) *DebugOverlay {
	ttf := ttfFromBytesMust(goregular.TTF, SCALE*8)
	return &DebugOverlay{
		perfGraph: NewGraph(win),

		enabled:             false,
		fpsText:             text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		gamepadAxisListText: text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		perfNumbersText:     text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		worldClockText:      text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
	}
}

func (do *DebugOverlay) Update(win *pixelgl.Window, frames int, fsToggled bool) {
	if win.JustReleased(pixelgl.KeyX) {
		do.enabled = !do.enabled
		if !do.enabled {
			do.fpsText.Clear()
			do.perfGraph.TimesPerFrame = []TimeSpent{}
			return
		}
	}

	if !do.enabled {
		return
	}

	do.perfGraph.FullscreenToggled = fsToggled

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

func (do *DebugOverlay) Draw(win *pixelgl.Canvas, gp *input.Gamepad, fps int) {
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

	// do.worldClockText.Clear()
	// _, err = do.worldClockText.WriteString(c.Current.String())
	// if err != nil {
	// 	panic(err)
	// }
	// do.worldClockText.Draw(win, pixel.IM.Scaled(
	// 	pixel.ZV, 1.2,
	// ).Moved(pixel.V(win.Bounds().W()-200, 10*1.2)))
}

func (do *DebugOverlay) SetDrawTimeDuration(d time.Duration)   { do.drawTimeDuration = d }
func (do *DebugOverlay) SetUpdateTimeDuration(d time.Duration) { do.updateTimeDuration = d }

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
