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
	WorldDrawTime   time.Duration
	WorldUpdateTime time.Duration
}

func NewGraph(win *pixelgl.Window, wdt, wut time.Duration) *Graph {
	return &Graph{
		imd:             imdraw.New(win),
		WorldDrawTime:   wdt,
		WorldUpdateTime: wut,
	}
}

func (g *Graph) Draw(win *pixelgl.Window) {
	g.imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))
	g.imd.Clear()
	g.imd.Color = colornames.Gray
	g.imd.Push(pixel.ZV)
	g.imd.Circle(500, 0)
	g.imd.Draw(win)
}
