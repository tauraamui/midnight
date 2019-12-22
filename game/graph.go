package game

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type TimeSpent struct {
	DrawTime   time.Duration
	UpdateTime time.Duration
}

type Graph struct {
	imd         *imdraw.IMDraw
	initialWinW float64
	w, h        float64

	TimesPerFrame []TimeSpent
}

func NewGraph(win *pixelgl.Window) *Graph {
	return &Graph{
		imd:           imdraw.New(win),
		initialWinW:   win.Bounds().W(),
		w:             win.Bounds().W() * 0.30,
		h:             100,
		TimesPerFrame: []TimeSpent{},
	}
}

func (g *Graph) Draw(win *pixelgl.Window) {
	if win.Bounds().W() != g.initialWinW {
		g.w = win.Bounds().W() * 0.30
		g.initialWinW = win.Bounds().W()
	}
	g.imd.SetMatrix(pixel.IM.Moved(pixel.V(win.Bounds().W()-g.w, win.Bounds().H()-g.h)))
	g.imd.Clear()
	g.imd.Color = colornames.Gray
	g.imd.Push(pixel.ZV)
	g.imd.Push(pixel.V(g.w, g.h))
	g.imd.Rectangle(0)
	g.imd.Color = colornames.Red
	for i, ts := range g.TimesPerFrame {
		var x float64 = 0
		if i > 0 {
			x = float64(i)
		}
		g.imd.Push(pixel.V(float64(5*x), g.h))
		g.imd.Push(pixel.V(float64(5*x), g.h-float64(ts.DrawTime.Microseconds()*10)))
		g.imd.Line(5)

		// g.imd.Push(pixel.V(float64(5*x), g.h))
		// g.imd.Push(pixel.V(float64(5*x), g.h-float64(ts.DrawTime.Microseconds()*10)))
		// g.imd.Line(5)
	}
	g.imd.Draw(win)
}
