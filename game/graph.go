package game

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Graph struct {
	imd             *imdraw.IMDraw
	initialWinW     float64
	w, h            float64
	WorldDrawTime   time.Duration
	WorldUpdateTime time.Duration
}

func NewGraph(win *pixelgl.Window, wdt, wut time.Duration) *Graph {
	return &Graph{
		imd:             imdraw.New(win),
		initialWinW:     win.Bounds().W(),
		w:               win.Bounds().W() * 0.30,
		h:               100,
		WorldDrawTime:   wdt,
		WorldUpdateTime: wut,
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
	g.imd.Draw(win)
}
