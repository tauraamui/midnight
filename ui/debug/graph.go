package debug

import (
	"fmt"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

type TimeSpent struct {
	DrawTime   time.Duration
	UpdateTime time.Duration
}

type Graph struct {
	TimesPerFrame []TimeSpent

	perfTimeSpentText *text.Text
	avgUpdateTime     int64
	avgDrawTime       int64
	imd               *imdraw.IMDraw
	initialWinW       float64
	barWidth          float64
	w, h              float64
}

func NewGraph(win *pixelgl.Window) *Graph {
	ttf := ttfFromBytesMust(goregular.TTF, SCALE*8)
	return &Graph{
		TimesPerFrame: []TimeSpent{},

		perfTimeSpentText: text.New(pixel.ZV, text.NewAtlas(ttf, text.ASCII, text.RangeTable(unicode.Latin))),
		imd:               imdraw.New(win),
		initialWinW:       win.Bounds().W(),
		barWidth:          (win.Bounds().W() * 0.45) / 60,
		w:                 win.Bounds().W() * 0.45,
		h:                 150,
	}
}

func (g *Graph) Draw(win *pixelgl.Canvas) {
	if win.Bounds().W() != g.initialWinW {
		// adding seemingly pointless 1 to the division in case list is ever 0
		// otherwise the whole program would come crashing down around our bun ears
		g.barWidth = (win.Bounds().W() * 0.30) / float64(len(g.TimesPerFrame)+1)
		g.w = win.Bounds().W() * 0.30
		g.initialWinW = win.Bounds().W()
	}
	g.imd.SetMatrix(pixel.IM.Moved(pixel.V(win.Bounds().W()-g.w, win.Bounds().H()-g.h)))
	g.imd.Clear()
	g.imd.Color = colornames.Gray
	g.imd.Push(pixel.ZV)
	g.imd.Push(pixel.V(g.w, g.h))
	g.imd.Rectangle(0)

	var updateTimeCount int64 = 0
	var drawTimeCount int64 = 0
	for i, ts := range g.TimesPerFrame {
		var x float64 = 0
		if i > 0 {
			x = float64(i)
			updateTimeCount += ts.UpdateTime.Microseconds()
			drawTimeCount += ts.DrawTime.Microseconds()
		}

		if i+1 == len(g.TimesPerFrame) {
			g.avgUpdateTime = updateTimeCount / int64(len(g.TimesPerFrame))
			g.avgDrawTime = drawTimeCount / int64(len(g.TimesPerFrame))
		}

		updateTimeBarHeight := g.h - float64(ts.UpdateTime.Microseconds())

		g.imd.Color = colornames.Blue
		g.imd.Push(pixel.V(float64((g.barWidth/2)+(g.barWidth*x)), g.h))
		g.imd.Push(pixel.V(float64((g.barWidth/2)+(g.barWidth*x)), updateTimeBarHeight))
		g.imd.Line(g.barWidth)

		g.imd.Color = colornames.Cyan
		g.imd.Push(pixel.V(float64((g.barWidth/2)+(g.barWidth*x)), updateTimeBarHeight))
		g.imd.Push(pixel.V(float64((g.barWidth/2)+(g.barWidth*x)), updateTimeBarHeight-float64(ts.DrawTime.Microseconds())))
		g.imd.Line(g.barWidth)
	}
	g.imd.Draw(win)

	win.SetMatrix(pixel.IM.Moved(pixel.V(win.Bounds().W()-g.w, win.Bounds().H()-g.h)))
	g.perfTimeSpentText.Clear()
	_, err := g.perfTimeSpentText.WriteString(fmt.Sprintf("AVG UPDATE+DRAW TIME: %d | %d MICRO", g.avgUpdateTime, g.avgDrawTime))
	if err != nil {
		panic(err)
	}
	g.perfTimeSpentText.Draw(win, pixel.IM.Scaled(
		pixel.ZV, 1.02,
	).Moved(pixel.ZV))

	win.SetMatrix(pixel.IM)
}
